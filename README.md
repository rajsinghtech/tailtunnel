# TailTunnel

<p align="center">
  <img src="frontend/static/logo.svg" alt="TailTunnel" width="128" height="128">
</p>

<p align="center">
  <strong>SSH into your Tailscale machines from your browser</strong>
</p>

<p align="center">
  TailTunnel gives you a beautiful web interface to manage and connect to all your Tailscale machines with SSH enabled.<br>
  No more remembering hostnames or IP addresses - just click and connect.
</p>

---

## What is TailTunnel?

TailTunnel is a web dashboard that shows all your SSH-enabled Tailscale machines in one place. Click on any machine to open a terminal session directly in your browser.

### Key Features

- **Auto-discovery** - Automatically finds all SSH-enabled machines on your tailnet
- **Visual dashboard** - See all your machines with their current status
- **Quick search** - Filter machines by name, user, or tag
- **Browser-based terminal** - Full terminal access without leaving your browser
- **Responsive design** - Works on desktop, tablet, and mobile
- **Secure** - Uses Tailscale's zero-trust security model

### What You Get

- Tiled grid view of all your machines
- Real-time online/offline status
- Machine tags for organization
- User ownership information
- One-click SSH connections
- Full terminal functionality in your browser

---

## Quick Start

The fastest way to get started is with Docker.

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

You'll see all your SSH-enabled machines. Click "Connect SSH" on any machine to start a terminal session.

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

---

## Requirements

### For Running (Docker)
- Docker installed
- Tailscale account
- At least one SSH-enabled machine on your tailnet

### For Building from Source
- Go 1.25 or later
- Node.js 20 or later
- npm

### Supported Platforms
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Runs anywhere Docker runs

---

## Troubleshooting

### TailTunnel won't start

Check the logs:
```bash
docker logs tailtunnel
```

Common issues:
- **"Invalid auth key"** - Your auth key expired or is incorrect
- **"Permission denied"** - Add `--cap-add=NET_ADMIN` to docker run command
- **"Port already in use"** - Change the port: `-p 8081:8080`

### No machines showing up

Make sure your machines:
1. Are online in Tailscale
2. Have SSH enabled
3. Are on the same tailnet as TailTunnel

Check if SSH is enabled on a machine:
```bash
tailscale status
```

Look for machines with "ssh" in their hostinfo.

### Can't connect to a machine

Verify SSH is working:
```bash
tailscale ssh username@machine-name
```

If this doesn't work, TailTunnel won't be able to connect either.

### Container keeps restarting

Check if you're running as the correct user:
```bash
docker exec tailtunnel id
```

The container runs as user `tailtunnel` (uid 1000) for security.

---

## Security Considerations

### How TailTunnel Connects

- Uses Tailscale's built-in authentication
- All traffic encrypted through Tailscale
- SSH connections proxied through WebSocket
- No external access - only accessible on your tailnet

### Default Connection Settings

- Connects as `root` user by default
- SSH host key verification disabled (Tailscale provides authentication)
- Only accessible from within your Tailscale network

### Best Practices

1. Use reusable auth keys with short expiration
2. Run TailTunnel on a trusted machine
3. Limit SSH access using Tailscale ACLs
4. Regularly update to the latest version

---

## For Developers

### Architecture

- **Backend**: Go with tsnet (embedded Tailscale node)
- **Frontend**: SvelteKit with Tailwind CSS
- **Terminal**: xterm.js with WebSocket
- **Build**: Multi-stage Docker build

### Project Structure

```
tailtunnel/
├── cmd/tailtunnel/       # Application entry point
├── internal/
│   ├── api/              # HTTP handlers and routing
│   ├── ssh/              # SSH WebSocket proxy
│   └── tailscale/        # Tailscale client wrapper
├── frontend/             # SvelteKit web interface
├── .github/workflows/    # CI/CD pipelines
└── Dockerfile            # Multi-arch container build
```

### CI/CD

TailTunnel uses GitHub Actions for:
- Automated testing on pull requests
- Multi-architecture builds (amd64, arm64)
- Container image signing with Cosign
- SBOM generation
- Automated releases

See `.github/workflows/` for details.

### Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

---

## License

MIT License - see LICENSE file for details

## Author

Created by Raj Singh ([@rajsinghtech](https://github.com/rajsinghtech))

---

## Links

- **GitHub**: https://github.com/rajsinghtech/tailtunnel
- **Issues**: https://github.com/rajsinghtech/tailtunnel/issues
- **Tailscale**: https://tailscale.com
