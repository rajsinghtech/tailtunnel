package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/rajsinghtech/tailtunnel"
	"github.com/rajsinghtech/tailtunnel/internal/api"
	"github.com/rajsinghtech/tailtunnel/internal/tailscale"
)

type Settings struct {
	Hostname string `json:"hostname"`
	StateDir string `json:"stateDir"`
	AuthKey  string `json:"authKey,omitempty"` // Optional - for backwards compat
}

var (
	settings        *Settings
	ts              *tailscale.TailscaleClient
	serverRunning   bool
	loginInProgress bool
	authURLOpened   bool
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("TailTunnel")
	systray.SetTooltip("TailTunnel - Not Connected")

	loadSettings()

	mStatus := systray.AddMenuItem("Status: Checking...", "Connection status")
	mStatus.Disable()
	systray.AddSeparator()

	mLogin := systray.AddMenuItem("Login to Tailscale", "Authenticate with Tailscale")
	mLogout := systray.AddMenuItem("Logout", "Disconnect from Tailscale")
	mLogout.Hide()

	systray.AddSeparator()
	mOpen := systray.AddMenuItem("Open Dashboard", "Open TailTunnel in browser")
	mOpen.Disable()

	systray.AddSeparator()
	mSettings := systray.AddMenuItem("Settings...", "Configure TailTunnel")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit TailTunnel")

	// Auto-start server
	go func() {
		if err := startServer(); err != nil {
			log.Printf("Failed to start server: %v", err)
			mStatus.SetTitle("Status: Failed to start")
		}
	}()

	go func() {
		for {
			select {
			case <-mLogin.ClickedCh:
				if !loginInProgress {
					go handleLogin()
				}
			case <-mLogout.ClickedCh:
				handleLogout()
				mLogin.Show()
				mLogout.Hide()
				mOpen.Disable()
				mStatus.SetTitle("Status: Logged Out")
			case <-mOpen.ClickedCh:
				openDashboard()
			case <-mSettings.ClickedCh:
				openSettings()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

	// Monitor status
	go func() {
		for {
			time.Sleep(3 * time.Second)
			if ts != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				status, err := ts.Status(ctx)
				cancel()

				if err == nil && status != nil {
					if status.BackendState == "Running" {
						fqdn := strings.TrimSuffix(status.Self.DNSName, ".")
						mStatus.SetTitle(fmt.Sprintf("Status: Connected (%s)", fqdn))
						systray.SetTooltip(fmt.Sprintf("TailTunnel - Connected to %s", fqdn))
						mOpen.Enable()
						mLogin.Hide()
						mLogout.Show()
						serverRunning = true
					} else {
						mStatus.SetTitle(fmt.Sprintf("Status: %s", status.BackendState))
					}
				}
			}
		}
	}()

}

func onExit() {
	stopServer()
}

func loadSettings() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get home directory: %v", err)
		settings = getDefaultSettings()
		return
	}

	configPath := filepath.Join(homeDir, ".tailtunnel", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			settings = getDefaultSettings()
			saveSettings()
			return
		}
		log.Printf("Failed to read config: %v", err)
		settings = getDefaultSettings()
		return
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		log.Printf("Failed to parse config: %v", err)
		settings = getDefaultSettings()
		return
	}
}

func saveSettings() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".tailtunnel")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func getDefaultSettings() *Settings {
	homeDir, _ := os.UserHomeDir()
	return &Settings{
		Hostname: "tailtunnel",
		StateDir: filepath.Join(homeDir, ".tailtunnel", "state"),
	}
}

func startServer() error {
	if serverRunning {
		return nil
	}

	var err error
	ts, err = tailscale.NewTailscaleClientWithConfig(tailscale.Config{
		Hostname: settings.Hostname,
		StateDir: settings.StateDir,
		AuthKey:  settings.AuthKey,
		Verbose:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize Tailscale client: %w", err)
	}

	handler := api.NewHandler(ts)
	router := api.NewRouter(handler, tailtunnel.FrontendFS)

	// Monitor for auth URLs (only open once)
	go func() {
		for authURL := range ts.AuthURL() {
			log.Printf("Auth URL received: %s", authURL)
			if !authURLOpened {
				authURLOpened = true
				log.Printf("Opening browser for first time: %s", authURL)
				// Silently open browser without blocking dialog
				if err := browser.OpenURL(authURL); err != nil {
					log.Printf("Failed to open browser: %v", err)
					// Only show error if browser fails to open
					showError("Browser Error", fmt.Sprintf("Could not open browser. Please visit:\n\n%s", authURL))
				}
			} else {
				log.Printf("Auth URL already opened, skipping duplicate")
			}
		}
	}()

	// Start server on tailnet (with HTTPS if available)
	go func() {
		log.Println("Starting server on tailnet...")
		if err := ts.ListenHTTPS(router); err != nil {
			log.Printf("Server error: %v", err)
			serverRunning = false
		}
	}()

	return nil
}

func stopServer() {
	if ts != nil {
		ts.Close()
		ts = nil
	}
	serverRunning = false
}

func handleLogin() {
	log.Println("=== handleLogin called ===")
	loginInProgress = true
	defer func() {
		loginInProgress = false
		log.Println("=== handleLogin finished ===")
	}()

	if ts == nil {
		log.Println("ts is nil, starting server...")
		if err := startServer(); err != nil {
			log.Printf("Failed to start server: %v", err)
			showError("Login Failed", fmt.Sprintf("Failed to start: %v", err))
			return
		}
		log.Println("Server started successfully")
	}

	// Don't block with a dialog - just wait in background
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("Calling WaitForLogin...")
	if err := ts.WaitForLogin(ctx); err != nil {
		log.Printf("WaitForLogin failed: %v", err)
		showError("Login Failed", fmt.Sprintf("Authentication failed: %v", err))
		return
	}

	log.Println("Login successful!")
	// Status monitoring will automatically update the menu when connected
}

func handleLogout() {
	if ts == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ts.Logout(ctx); err != nil {
		showError("Logout Failed", fmt.Sprintf("Failed to logout: %v", err))
		return
	}

	stopServer()
	authURLOpened = false // Reset for next login
	showInfo("Logged Out", "You have been logged out of Tailscale.")
}

func openDashboard() {
	if ts == nil || !serverRunning {
		showError("Not Connected", "Please login to Tailscale first.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	status, err := ts.Status(ctx)
	cancel()

	if err != nil || status == nil {
		showError("Error", "Failed to get connection status")
		return
	}

	// Try HTTPS first
	fqdn := strings.TrimSuffix(status.Self.DNSName, ".")
	url := fmt.Sprintf("https://%s/", fqdn)
	if err := browser.OpenURL(url); err != nil {
		// Fallback to HTTP
		url = fmt.Sprintf("http://%s/", fqdn)
		if err := browser.OpenURL(url); err != nil {
			showError("Failed to Open Browser", fmt.Sprintf("Could not open browser: %v\n\nPlease visit: %s", err, url))
		}
	}
}

func openSettings() {
	script := fmt.Sprintf(`
tell application "System Events"
	activate
	set hostname to text returned of (display dialog "Enter TailTunnel hostname:" default answer "%s" with title "TailTunnel Settings")
	return hostname
end tell
`, settings.Hostname)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	newHostname := string(output)
	if newHostname != "" && newHostname != "\n" {
		newHostname = newHostname[:len(newHostname)-1] // Remove trailing newline
		settings.Hostname = newHostname
		if err := saveSettings(); err != nil {
			showError("Save Failed", fmt.Sprintf("Failed to save settings: %v", err))
		} else {
			showInfo("Settings Saved", "Hostname updated. Please logout and login again for changes to take effect.")
		}
	}
}

func showError(title, message string) {
	script := fmt.Sprintf(`
tell application "System Events"
	display dialog "%s" with title "%s" buttons {"OK"} default button 1 with icon stop
end tell
`, message, title)
	exec.Command("osascript", "-e", script).Run()
}

func showInfo(title, message string) {
	script := fmt.Sprintf(`
tell application "System Events"
	display dialog "%s" with title "%s" buttons {"OK"} default button 1
end tell
`, message, title)
	exec.Command("osascript", "-e", script).Run()
}
