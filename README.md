# TailTunnel

<p align="center">
  <img src="frontend/static/logo.svg" alt="TailTunnel" width="128" height="128">
</p>

<p align="center">
  <strong>Your Tailscale Toolkit - SSH and Network Diagnostics in Your Browser</strong>
</p>

<p align="center">
  TailTunnel gives you a beautiful web interface to manage your Tailscale network.<br>
  Browser-based SSH, real-time network diagnostics, and connection monitoring - all in one place.
</p>

---

## What is TailTunnel?

TailTunnel is a web dashboard for your Tailscale network. Connect to machines via SSH directly in your browser, monitor network health with TailCanary, and see real-time connection quality across your entire tailnet.

### Key Features

- **Browser SSH** - Full terminal access without leaving your browser
- **TailCanary Network Diagnostics** - Real-time ping monitoring with latency graphs
- **Connection Type Detection** - See if peers are using direct, DERP relay, or peer relay connections
- **Auto-discovery** - Automatically finds all machines on your tailnet
- **Live Monitoring** - Continuous network health checks with automatic refresh
- **Device Info** - View NAT type, online status, and system information
- **Quick Search** - Filter machines by name, user, tag, or connection type
- **Responsive Design** - Works on desktop, tablet, and mobile
- **OAuth Login** - No auth keys required for the macOS app

### What You Get

- **SSH Machines**: Tiled grid view with one-click terminal access
- **TailCanary**: Network diagnostics with latency tracking and connection type visualization
- **Real-time Status**: Live online/offline status and connection quality
- **Historical Data**: Latency graphs showing connection performance over time
- **Smart Filtering**: Search across all attributes including connection types
- **HTTPS on Tailnet**: Automatic HTTPS with valid certificates when available

---

## Quick Start

**Two ways to run TailTunnel:**
- **macOS App** (`brew install --cask`) - GUI menu bar app with OAuth login
- **CLI Binary** (`brew install`) - Command-line tool for servers/headless systems
- **Docker** - Containerized deployment for always-on servers

Choose the method that best fits your needs below.

### macOS Menu Bar App (Recommended for Desktop)

The easiest way to run TailTunnel on macOS - no auth keys needed!

**Install via Homebrew:**
```bash
brew install --cask rajsinghtech/tap/tailtunnel
```

**Or download manually:**
1. Download `TailTunnel.zip` from [releases](https://github.com/rajsinghtech/tailtunnel/releases/latest)
2. Unzip and drag `TailTunnel.app` to `/Applications`
3. Launch the app
4. Click "Login to Tailscale" - your browser will open automatically
5. Authenticate in your browser
6. Click "Open Dashboard" when connected

**Features:**
- OAuth login (no auth keys!)
- Lives in your menu bar
- Auto-opens browser for authentication
- Serves on your tailnet with HTTPS
- Configurable hostname

**Configuration:**
Settings stored at `~/.tailtunnel/config.json`:
- Hostname (default: `tailtunnel`)
- State directory (default: `~/.tailtunnel/state`)

**Uninstall:**
```bash
brew uninstall --cask tailtunnel
# Or manually:
rm -rf /Applications/TailTunnel.app ~/.tailtunnel
```

### Docker

Perfect for servers and always-on deployments.

**Step 1: Get a Tailscale Auth Key**

1. Go to https://login.tailscale.com/admin/settings/keys
2. Click "Generate auth key"
3. Copy the key (starts with `tskey-auth-`)

**Step 2: Run with Docker**

```bash
docker run -d \
  --name tailtunnel \
  -e TS_AUTHKEY=your-key-here \
  -v tailtunnel-state:/var/lib/tailtunnel \
  ghcr.io/rajsinghtech/tailtunnel:latest
```

**Or with Docker Compose:**

Create `docker-compose.yml`:
```yaml
services:
  tailtunnel:
    image: ghcr.io/rajsinghtech/tailtunnel:latest
    container_name: tailtunnel
    restart: unless-stopped
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
    volumes:
      - tailtunnel-state:/var/lib/tailtunnel

volumes:
  tailtunnel-state:
```

Create `.env`:
```bash
TS_AUTHKEY=your-tailscale-auth-key
```

Start:
```bash
docker-compose up -d
```

**Step 3: Access the Dashboard**

TailTunnel serves on your tailnet with automatic HTTPS:
```
https://tailtunnel.your-tailnet.ts.net/
```

You'll see:
- **TailCanary**: Real-time network diagnostics with ping monitoring and latency graphs
- **SSH Machines**: All your SSH-enabled machines with one-click terminal access

### CLI Binary via Homebrew (macOS/Linux)

For headless servers, scripting, or if you prefer command-line tools:

```bash
brew install rajsinghtech/tap/tailtunnel
```

Run with auth key:
```bash
export TS_AUTHKEY=your-tailscale-auth-key
tailtunnel
```

Or with OAuth (opens browser):
```bash
tailtunnel  # Browser will open for authentication
```

### Build from Source

#### Prerequisites
- Go 1.25 or later
- Node.js 20 or later

#### Build Steps

**CLI binary:**
```bash
git clone https://github.com/rajsinghtech/tailtunnel.git
cd tailtunnel
make build
./tailtunnel
```

**macOS app:**
```bash
make build-macos-app
cp -r TailTunnel.app /Applications/
```

**Docker image:**
```bash
docker build -t tailtunnel .
```

---

## Configuration

### macOS App

Configure via the Settings menu:
- **Hostname**: Your tailnet hostname (default: `tailtunnel`)

Settings stored at `~/.tailtunnel/config.json`

### Docker / CLI

Configure using environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TS_AUTHKEY` | Tailscale auth key | - | Yes (Docker/headless) |
| `STATE_DIR` | Tailscale state directory | `/var/lib/tailtunnel` | No |

**Note:** The macOS app uses OAuth and doesn't need an auth key. For Docker/CLI, you can omit `TS_AUTHKEY` to use OAuth (opens browser).

### Getting a Tailscale Auth Key

1. Visit https://login.tailscale.com/admin/settings/keys
2. Click "Generate auth key"
3. Choose reusable and no expiration for servers
4. Copy the key

---

## License

MIT License - see [LICENSE](LICENSE) for details.
