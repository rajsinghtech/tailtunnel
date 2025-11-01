package tailscale

import (
	"context"
	"fmt"
	"net"
	"os"

	"tailscale.com/client/tailscale"
	"tailscale.com/tsnet"
)

type TailscaleClient struct {
	server *tsnet.Server
	lc     *tailscale.LocalClient
}

func NewTailscaleClient() (*TailscaleClient, error) {
	stateDir := os.Getenv("STATE_DIR")
	if stateDir == "" {
		stateDir = "/var/lib/tailtunnel"
	}

	srv := &tsnet.Server{
		Hostname: "tailtunnel",
		Dir:      stateDir,
		AuthKey:  os.Getenv("TS_AUTHKEY"),
	}

	if err := srv.Start(); err != nil {
		return nil, fmt.Errorf("failed to start tsnet server: %w", err)
	}

	lc, err := srv.LocalClient()
	if err != nil {
		srv.Close()
		return nil, fmt.Errorf("failed to get local client: %w", err)
	}

	return &TailscaleClient{
		server: srv,
		lc:     lc,
	}, nil
}

func (tc *TailscaleClient) DialSSH(ctx context.Context, machine string) (net.Conn, error) {
	return tc.server.Dial(ctx, "tcp", machine+":22")
}

func (tc *TailscaleClient) Close() error {
	return tc.server.Close()
}
