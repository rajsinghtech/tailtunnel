package canary

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"sync"
	"time"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
)

type Pinger struct {
	lc *tailscale.LocalClient
}

func NewPinger(lc *tailscale.LocalClient) *Pinger {
	return &Pinger{lc: lc}
}

func (p *Pinger) GetPeers(ctx context.Context) (*PeersResponse, error) {
	status, err := p.lc.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var peers []PeerInfo
	for _, peer := range status.Peer {
		tags := []string{}
		if peer.Tags != nil {
			tags = peer.Tags.AsSlice()
		}

		ip := ""
		if len(peer.TailscaleIPs) > 0 {
			ip = peer.TailscaleIPs[0].String()
		}

		peerRelay := ""
		if len(peer.PeerAPIURL) > 0 {
			peerRelay = peer.PeerAPIURL[0]
		}

		// Use DNS name to extract hostname if HostName is generic
		hostName := peer.HostName
		if hostName == "localhost" || hostName == "" {
			hostName = extractHostnameFromDNS(peer.DNSName)
		}

		// Get user information
		userLogin := ""
		userDisplay := ""
		if user, ok := status.User[peer.UserID]; ok {
			userLogin = user.LoginName
			userDisplay = user.DisplayName
		}

		peers = append(peers, PeerInfo{
			HostName:      hostName,
			DNSName:       peer.DNSName,
			IP:            ip,
			Online:        peer.Online,
			Active:        peer.Active,
			CurAddr:       peer.CurAddr,
			Relay:         peer.Relay,
			PeerRelay:     peerRelay,
			LastHandshake: peer.LastHandshake,
			RxBytes:       peer.RxBytes,
			TxBytes:       peer.TxBytes,
			OS:            peer.OS,
			UserLogin:     userLogin,
			UserDisplay:   userDisplay,
			Tags:          tags,
		})
	}

	return &PeersResponse{
		Peers:     peers,
		Timestamp: time.Now(),
	}, nil
}

// extractHostnameFromDNS extracts hostname from DNS name like "iphone172.keiretsu.ts.net."
func extractHostnameFromDNS(dnsName string) string {
	if dnsName == "" {
		return "unknown"
	}
	// Remove trailing dot
	dnsName = strings.TrimSuffix(dnsName, ".")
	// Split by dot and take first part
	parts := strings.Split(dnsName, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return dnsName
}

func (p *Pinger) Ping(ctx context.Context, ipStr string) (*PingResult, error) {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address: %w", err)
	}

	// Use "disco" ping type - fastest, doesn't involve IP layer
	result, err := p.lc.Ping(ctx, ip, "disco")
	if err != nil {
		return &PingResult{
			IP:      ipStr,
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	if result.Err != "" {
		return &PingResult{
			IP:       ipStr,
			NodeName: result.NodeName,
			Success:  false,
			Error:    result.Err,
		}, nil
	}

	connType := p.determineConnectionType(result)
	latencyMs := result.LatencySeconds * 1000

	return &PingResult{
		IP:             ipStr,
		NodeName:       result.NodeName,
		Success:        true,
		LatencyMs:      latencyMs,
		ConnectionType: connType,
		Endpoint:       result.Endpoint,
		DERPRegion:     result.DERPRegionCode,
		DERPRegionID:   result.DERPRegionID,
		PeerRelay:      result.PeerRelay,
	}, nil
}

func (p *Pinger) PingAll(ctx context.Context) (*PingAllResponse, error) {
	peers, err := p.GetPeers(ctx)
	if err != nil {
		return nil, err
	}

	var results []PingResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, peer := range peers.Peers {
		if !peer.Online || peer.IP == "" {
			continue
		}

		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			// Try with 10 second timeout first
			pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			result, err := p.Ping(pingCtx, ip)

			// If it timed out, try one more time with a shorter timeout (5s)
			if err != nil && (result == nil || !result.Success) {
				retryCtx, retryCancel := context.WithTimeout(ctx, 5*time.Second)
				defer retryCancel()

				retryResult, retryErr := p.Ping(retryCtx, ip)
				if retryErr == nil && retryResult != nil && retryResult.Success {
					result = retryResult
					err = nil
				}
			}

			if err != nil {
				result = &PingResult{
					IP:      ip,
					Success: false,
					Error:   err.Error(),
				}
			}

			mu.Lock()
			results = append(results, *result)
			mu.Unlock()
		}(peer.IP)
	}

	wg.Wait()

	return &PingAllResponse{
		Results:   results,
		Timestamp: time.Now(),
	}, nil
}

func (p *Pinger) determineConnectionType(result *ipnstate.PingResult) ConnectionType {
	// Check in order of preference: direct > peer-relay > derp
	if result.Endpoint != "" {
		return ConnectionDirect
	}
	if result.PeerRelay != "" {
		return ConnectionPeerRelay
	}
	if result.DERPRegionID != 0 {
		return ConnectionDERP
	}
	return ConnectionUnknown
}
