package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/rajsinghtech/tailtunnel"
	"github.com/rajsinghtech/tailtunnel/internal/api"
	"github.com/rajsinghtech/tailtunnel/internal/tailscale"
)

type Settings struct {
	AuthKey  string `json:"authKey"`
	StateDir string `json:"stateDir"`
	Port     int    `json:"port"`
}

var (
	settings     *Settings
	srv          *http.Server
	ts           *tailscale.TailscaleClient
	serverRunning bool
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("TailTunnel")
	systray.SetTooltip("TailTunnel - Tailscale Toolkit")

	loadSettings()

	mOpen := systray.AddMenuItem("Open Dashboard", "Open TailTunnel in browser")
	systray.AddSeparator()
	mStart := systray.AddMenuItem("Start Server", "Start the TailTunnel server")
	mStop := systray.AddMenuItem("Stop Server", "Stop the TailTunnel server")
	mStop.Disable()
	systray.AddSeparator()
	mSettings := systray.AddMenuItem("Settings...", "Configure TailTunnel")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit TailTunnel")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openDashboard()
			case <-mStart.ClickedCh:
				if startServer() {
					mStart.Disable()
					mStop.Enable()
					systray.SetTooltip("TailTunnel - Running")
				}
			case <-mStop.ClickedCh:
				stopServer()
				mStart.Enable()
				mStop.Disable()
				systray.SetTooltip("TailTunnel - Stopped")
			case <-mSettings.ClickedCh:
				openSettings()
			case <-mQuit.ClickedCh:
				systray.Quit()
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
		StateDir: filepath.Join(homeDir, ".tailtunnel", "state"),
		Port:     8080,
	}
}

func startServer() bool {
	if serverRunning {
		return true
	}

	if settings.AuthKey == "" {
		showError("Auth Key Required", "Please configure your Tailscale auth key in Settings.\n\nGet one from: https://login.tailscale.com/admin/settings/keys")
		return false
	}

	os.Setenv("TS_AUTHKEY", settings.AuthKey)
	os.Setenv("STATE_DIR", settings.StateDir)
	os.Setenv("PORT", fmt.Sprintf("%d", settings.Port))

	var err error
	ts, err = tailscale.NewTailscaleClient()
	if err != nil {
		showError("Failed to Start", fmt.Sprintf("Failed to initialize Tailscale client: %v", err))
		return false
	}

	handler := api.NewHandler(ts)
	router := api.NewRouter(handler, tailtunnel.FrontendFS)

	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on http://localhost:%d", settings.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
			serverRunning = false
		}
	}()

	time.Sleep(500 * time.Millisecond)
	serverRunning = true
	return true
}

func stopServer() {
	if !serverRunning {
		return
	}

	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		srv = nil
	}

	if ts != nil {
		ts.Close()
		ts = nil
	}

	serverRunning = false
}

func openDashboard() {
	if !serverRunning {
		if !startServer() {
			return
		}
	}

	url := fmt.Sprintf("http://localhost:%d", settings.Port)
	if err := browser.OpenURL(url); err != nil {
		showError("Failed to Open Browser", fmt.Sprintf("Could not open browser: %v\n\nPlease visit: %s", err, url))
	}
}

func openSettings() {
	script := fmt.Sprintf(`
tell application "System Events"
	activate
	set authKey to text returned of (display dialog "Enter your Tailscale Auth Key:" default answer "%s" with title "TailTunnel Settings")
	return authKey
end tell
`, settings.AuthKey)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	newAuthKey := string(output)
	if newAuthKey != "" && newAuthKey != "\n" {
		newAuthKey = newAuthKey[:len(newAuthKey)-1] // Remove trailing newline
		settings.AuthKey = newAuthKey
		if err := saveSettings(); err != nil {
			showError("Save Failed", fmt.Sprintf("Failed to save settings: %v", err))
		} else {
			showInfo("Settings Saved", "Your auth key has been saved. Restart the server for changes to take effect.")
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

func init() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		onExit()
		os.Exit(0)
	}()
}
