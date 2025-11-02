package tailscale

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tailcfg"
	"tailscale.com/tsnet"
)

type TailscaleClient struct {
	server      *tsnet.Server
	lc          *tailscale.LocalClient
	authURLChan chan string
	statusChan  chan *ipnstate.Status
}

type Config struct {
	Hostname string
	StateDir string
	AuthKey  string
	Verbose  bool
}

func NewTailscaleClient() (*TailscaleClient, error) {
	return NewTailscaleClientWithConfig(Config{
		Hostname: "tailtunnel",
		StateDir: getStateDir(),
		AuthKey:  os.Getenv("TS_AUTHKEY"),
		Verbose:  false,
	})
}

func NewTailscaleClientWithConfig(cfg Config) (*TailscaleClient, error) {
	tc := &TailscaleClient{
		authURLChan: make(chan string, 1),
		statusChan:  make(chan *ipnstate.Status, 1),
	}

	srv := &tsnet.Server{
		Hostname:     cfg.Hostname,
		Dir:          cfg.StateDir,
		AuthKey:      cfg.AuthKey,
		RunWebClient: true,
		Logf:         tc.logFunc(cfg.Verbose),
	}

	if err := srv.Start(); err != nil {
		return nil, fmt.Errorf("failed to start tsnet server: %w", err)
	}

	lc, err := srv.LocalClient()
	if err != nil {
		srv.Close()
		return nil, fmt.Errorf("failed to get local client: %w", err)
	}

	tc.server = srv
	tc.lc = lc

	// Start monitoring for auth URL and status
	go tc.monitorStatus()

	return tc, nil
}

func (tc *TailscaleClient) logFunc(verbose bool) func(string, ...any) {
	return func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)

		// Extract auth URL from log messages
		if strings.Contains(msg, "go to:") {
			parts := strings.Split(msg, "go to:")
			if len(parts) > 1 {
				url := strings.TrimSpace(parts[1])
				select {
				case tc.authURLChan <- url:
				default:
				}
			}
		}

		if verbose {
			log.Printf("[tsnet] %s", msg)
		}
	}
}

func (tc *TailscaleClient) monitorStatus() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		status, err := tc.lc.StatusWithoutPeers(context.Background())
		if err != nil {
			continue
		}

		// Send auth URL if available
		if status.AuthURL != "" {
			select {
			case tc.authURLChan <- status.AuthURL:
			default:
			}
		}

		// Send status updates
		select {
		case tc.statusChan <- status:
		default:
		}

		// Stop monitoring once we're running
		if status.BackendState == ipn.Running.String() {
			return
		}
	}
}

func (tc *TailscaleClient) WaitForLogin(ctx context.Context) error {
	watcher, err := tc.lc.WatchIPNBus(ctx, ipn.NotifyInitialState|ipn.NotifyNoPrivateKeys)
	if err != nil {
		return fmt.Errorf("failed to watch IPN bus: %w", err)
	}
	defer watcher.Close()

	for {
		n, err := watcher.Next()
		if err != nil {
			return fmt.Errorf("watch error: %w", err)
		}

		if n.ErrMessage != nil {
			return fmt.Errorf("backend error: %s", *n.ErrMessage)
		}

		if n.BrowseToURL != nil && *n.BrowseToURL != "" {
			select {
			case tc.authURLChan <- *n.BrowseToURL:
			default:
			}
		}

		if n.State != nil && *n.State == ipn.Running {
			return nil
		}
	}
}

func (tc *TailscaleClient) Up(ctx context.Context) (*ipnstate.Status, error) {
	status, err := tc.server.Up(ctx)
	if err != nil {
		return nil, err
	}

	select {
	case tc.statusChan <- status:
	default:
	}

	return status, nil
}

func (tc *TailscaleClient) Status(ctx context.Context) (*ipnstate.Status, error) {
	return tc.lc.Status(ctx)
}

func (tc *TailscaleClient) Logout(ctx context.Context) error {
	return tc.lc.Logout(ctx)
}

func (tc *TailscaleClient) AuthURL() <-chan string {
	return tc.authURLChan
}

func (tc *TailscaleClient) StatusUpdates() <-chan *ipnstate.Status {
	return tc.statusChan
}

func (tc *TailscaleClient) ListenHTTPS(handler http.Handler) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, err := tc.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to bring up tailscale: %w", err)
	}

	enableHTTPS := status.Self.HasCap(tailcfg.CapabilityHTTPS) && len(tc.server.CertDomains()) > 0
	fqdn := strings.TrimSuffix(status.Self.DNSName, ".")

	if enableHTTPS {
		httpsListener, err := tc.server.ListenTLS("tcp", ":443")
		if err != nil {
			return fmt.Errorf("failed to listen on :443: %w", err)
		}

		httpListener, err := tc.server.Listen("tcp", ":80")
		if err != nil {
			return fmt.Errorf("failed to listen on :80: %w", err)
		}

		log.Printf("Serving HTTPS on https://%s/", fqdn)

		// Serve HTTPS
		go func() {
			if err := http.Serve(httpsListener, handler); err != nil {
				log.Printf("HTTPS server error: %v", err)
			}
		}()

		// Redirect HTTP to HTTPS
		redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+fqdn+r.RequestURI, http.StatusMovedPermanently)
		})

		log.Printf("Redirecting HTTP to HTTPS")
		return http.Serve(httpListener, redirectHandler)
	}

	// Fallback to HTTP only
	httpListener, err := tc.server.Listen("tcp", ":80")
	if err != nil {
		return fmt.Errorf("failed to listen on :80: %w", err)
	}

	log.Printf("Serving HTTP on http://%s/", fqdn)
	return http.Serve(httpListener, handler)
}

func (tc *TailscaleClient) DialSSH(ctx context.Context, machine string) (net.Conn, error) {
	return tc.server.Dial(ctx, "tcp", machine+":22")
}

func (tc *TailscaleClient) LocalClient() *tailscale.LocalClient {
	return tc.lc
}

func (tc *TailscaleClient) Close() error {
	return tc.server.Close()
}

func getStateDir() string {
	stateDir := os.Getenv("STATE_DIR")
	if stateDir == "" {
		// Use user's home directory by default for better UX
		if homeDir, err := os.UserHomeDir(); err == nil {
			stateDir = homeDir + "/.tailtunnel/state"
			// Create directory if it doesn't exist
			os.MkdirAll(stateDir, 0700)
		} else {
			// Fallback to /var/lib/tailtunnel (requires elevated permissions)
			stateDir = "/var/lib/tailtunnel"
		}
	}
	return stateDir
}
