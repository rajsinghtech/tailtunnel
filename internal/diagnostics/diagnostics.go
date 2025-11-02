package diagnostics

import (
	"context"
	"fmt"
	"net"

	"tailscale.com/client/tailscale"
)

type DiagnosticsInfo struct {
	HostName         string `json:"hostName"`
	TailscaleIP      string `json:"tailscaleIP"`
	DNSName          string `json:"dnsName"`
	OS               string `json:"os"`
	Online           bool   `json:"online"`
	ExitNodeID       string `json:"exitNodeId,omitempty"`
	MagicDNSSuffix   string `json:"magicDnsSuffix"`
	NATType          string `json:"natType,omitempty"`
	PortMapProtocol  string `json:"portMapProtocol,omitempty"`
}

func GetDiagnostics(ctx context.Context, lc *tailscale.LocalClient) (*DiagnosticsInfo, error) {
	status, err := lc.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	info := &DiagnosticsInfo{
		HostName:       status.Self.HostName,
		DNSName:        status.Self.DNSName,
		OS:             status.Self.OS,
		Online:         status.Self.Online,
		MagicDNSSuffix: status.MagicDNSSuffix,
	}

	// Get Tailscale IP
	if len(status.Self.TailscaleIPs) > 0 {
		info.TailscaleIP = status.Self.TailscaleIPs[0].String()
	}

	// Get exit node if set
	if status.ExitNodeStatus != nil {
		info.ExitNodeID = string(status.ExitNodeStatus.ID)
	}

	// Detect NAT type and port mapping
	natType, portMapProto := detectNATWithPortMapper()
	info.NATType = natType
	info.PortMapProtocol = portMapProto

	return info, nil
}

func detectNATWithPortMapper() (string, string) {
	// Check if we're behind NAT using basic IP check
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Unknown", ""
	}

	hasPrivate := false
	hasPublic := false

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if isPrivateIP(ipnet.IP) {
					hasPrivate = true
				} else {
					hasPublic = true
				}
			}
		}
	}

	// If we have public IPs, we're not behind NAT
	if hasPublic && !hasPrivate {
		return "No NAT", ""
	}

	// If we're not behind NAT, return early
	if !hasPrivate {
		return "Unknown", ""
	}

	// We're behind NAT
	// For now, we'll classify as "EZ NAT" which is more common with modern routers
	// A more sophisticated implementation would use STUN to determine if it's symmetric (Hard NAT)
	// Most consumer routers support endpoint-independent mapping (EZ NAT)
	return "EZ NAT", ""
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	private := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range private {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
