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

### What You Get

- **SSH Machines**: Tiled grid view with one-click terminal access
- **TailCanary**: Network diagnostics with latency tracking and connection type visualization
- **Real-time Status**: Live online/offline status and connection quality
- **Historical Data**: Latency graphs showing connection performance over time
- **Smart Filtering**: Search across all attributes including connection types

---

## Quick Start

### Homebrew (macOS/Linux)

Install TailTunnel CLI via Homebrew:

```bash
brew install rajsinghtech/tap/tailtunnel
```

Once installed, you can run:
```bash
tailtunnel
```

Configure with environment variables:
```bash
export TS_AUTHKEY=your-tailscale-auth-key
export PORT=8080
tailtunnel
```

### macOS Menu Bar App

The easiest way to run TailTunnel on macOS is with our native menu bar app:

**Option 1: Download Release**
1. Download `TailTunnel.zip` from the [latest release](https://github.com/rajsinghtech/tailtunnel/releases/latest)
2. Unzip and drag `TailTunnel.app` to your `/Applications` folder
3. Remove Gatekeeper quarantine (app is unsigned):
   ```bash
   xattr -cr /Applications/TailTunnel.app
   ```
4. Launch the app and configure your Tailscale auth key in Settings

> **Note:** If you see a "damaged" error, run the `xattr` command above to bypass macOS Gatekeeper.

**Option 2: Build from Source**
```bash
git clone https://github.com/rajsinghtech/tailtunnel.git
cd tailtunnel
make build-macos-app
cp -r TailTunnel.app /Applications/
```

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

### Docker

The fastest way to get started with Docker:

### Step 1: Get a Tailscale Auth Key

1. Go to https://login.tailscale.com/admin/settings/keys
2. Click "Generate auth key"
3. Copy the key (starts with `tskey-auth-`)

### Step 2: Run TailTunnel

```bash
docker run -d \
  --name tailtunnel \
  --cap-add=NET_ADMIN \
  -p 8080:8080 \
  -e TS_AUTHKEY=your-key-here \
  ghcr.io/rajsinghtech/tailtunnel:latest
```

Replace `your-key-here` with the auth key from Step 1.

### Step 3: Access the Dashboard

Open your browser and go to:
```
http://localhost:8080
```

You'll see two main features:
- **TailCanary**: Real-time network diagnostics with ping monitoring and latency graphs
- **SSH Machines**: All your SSH-enabled machines with one-click terminal access

---

## Alternative Installation Methods

### Using Docker Compose

1. Create a `.env` file:
```bash
TS_AUTHKEY=your-tailscale-auth-key
PORT=8080
```

2. Create `docker-compose.yml`:
```yaml
services:
  tailtunnel:
    image: ghcr.io/rajsinghtech/tailtunnel:latest
    container_name: tailtunnel
    restart: unless-stopped
    ports:
      - "${PORT:-8080}:8080"
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
    volumes:
      - tailtunnel-state:/var/lib/tailtunnel
    cap_add:
      - NET_ADMIN

volumes:
  tailtunnel-state:
```

3. Start the service:
```bash
docker-compose up -d
```

### Building from Source

#### Prerequisites
- Go 1.25 or later
- Node.js 20 or later

#### Build Steps

```bash
git clone https://github.com/rajsinghtech/tailtunnel.git
cd tailtunnel
cp .env.example .env
# Edit .env and add your Tailscale auth key
make install
make dev
```

The dashboard will be available at http://localhost:8080

---

## Configuration

TailTunnel is configured using environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TS_AUTHKEY` | Tailscale authentication key | - | Yes (first run) |
| `PORT` | HTTP server port | 8080 | No |
| `STATE_DIR` | Directory for Tailscale state | /var/lib/tailtunnel | No |

### Getting a Tailscale Auth Key

1. Visit https://login.tailscale.com/admin/settings/keys
2. Click "Generate auth key"
3. Choose your expiration preference
4. Copy the key

You only need the auth key on the first run. After that, TailTunnel remembers your connection.

