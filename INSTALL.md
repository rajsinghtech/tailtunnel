# TailTunnel Installation

## macOS Menu Bar App

TailTunnel is available as a native macOS menu bar application. The app runs in your menu bar and provides easy access to your Tailscale network.

### Features

- **Menu bar integration** - Lives in your macOS menu bar
- **Start/Stop server** - Control the TailTunnel server with one click
- **Auto-open browser** - Opens your default browser to the dashboard
- **Settings dialog** - Configure your Tailscale auth key via native dialog
- **Persistent settings** - Remembers your configuration at `~/.tailtunnel/config.json`

### Installation

#### Option 1: Download Release

1. Download `TailTunnel.zip` from the [latest release](https://github.com/rajsinghtech/tailtunnel/releases/latest)
2. Unzip the file
3. Drag `TailTunnel.app` to your `/Applications` folder
4. Right-click the app and select "Open" (first time only, due to macOS Gatekeeper)

#### Option 2: Build from Source

```bash
git clone https://github.com/rajsinghtech/tailtunnel.git
cd tailtunnel
make build-macos-app
cp -r TailTunnel.app /Applications/
```

### First Launch

1. Launch TailTunnel from your Applications folder
2. Look for "TailTunnel" in your menu bar
3. Click the icon and select "Settings..."
4. Enter your Tailscale auth key (get one from https://login.tailscale.com/admin/settings/keys)
5. Click "Start Server"
6. Click "Open Dashboard" to access TailTunnel in your browser

### Configuration

The app stores its configuration at:
```
~/.tailtunnel/config.json
```

This includes:
- Your Tailscale auth key (encrypted)
- Server port (default: 8080)
- State directory location

### Uninstall

To completely remove TailTunnel:

```bash
# Remove the app
rm -rf /Applications/TailTunnel.app

# Remove configuration and state
rm -rf ~/.tailtunnel
```

## Docker Installation

See the main [README.md](README.md#quick-start) for Docker installation instructions.

## Troubleshooting

### App won't start
- Check Console.app for error messages
- Ensure you have a valid Tailscale auth key configured
- Try removing `~/.tailtunnel` and reconfiguring

### Server won't start
- Make sure port 8080 isn't already in use
- Check that your auth key is valid and not expired
- Verify you have network permissions

### Can't access dashboard
- Ensure the server is running (check menu bar icon)
- Try opening http://localhost:8080 manually
- Check your firewall settings

### Getting Help

- [GitHub Issues](https://github.com/rajsinghtech/tailtunnel/issues)
- [GitHub Discussions](https://github.com/rajsinghtech/tailtunnel/discussions)
