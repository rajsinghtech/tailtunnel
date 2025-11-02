package canary

import "time"

type ConnectionType string

const (
	ConnectionDirect    ConnectionType = "direct"
	ConnectionDERP      ConnectionType = "derp"
	ConnectionPeerRelay ConnectionType = "peer-relay"
	ConnectionOffline   ConnectionType = "offline"
	ConnectionUnknown   ConnectionType = "unknown"
)

type PeerInfo struct {
	HostName      string         `json:"hostName"`
	DNSName       string         `json:"dnsName"`
	IP            string         `json:"ip"`
	Online        bool           `json:"online"`
	Active        bool           `json:"active"`
	CurAddr       string         `json:"curAddr,omitempty"`
	Relay         string         `json:"relay,omitempty"`
	PeerRelay     string         `json:"peerRelay,omitempty"`
	LastHandshake time.Time      `json:"lastHandshake"`
	RxBytes       int64          `json:"rxBytes"`
	TxBytes       int64          `json:"txBytes"`
	OS            string         `json:"os"`
	UserLogin     string         `json:"userLogin,omitempty"`
	UserDisplay   string         `json:"userDisplay,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
}

type PingResult struct {
	IP             string         `json:"ip"`
	NodeName       string         `json:"nodeName"`
	Success        bool           `json:"success"`
	Error          string         `json:"error,omitempty"`
	LatencyMs      float64        `json:"latencyMs"`
	ConnectionType ConnectionType `json:"connectionType"`
	Endpoint       string         `json:"endpoint,omitempty"`
	DERPRegion     string         `json:"derpRegion,omitempty"`
	DERPRegionID   int            `json:"derpRegionId,omitempty"`
	PeerRelay      string         `json:"peerRelay,omitempty"`
}

type PeersResponse struct {
	Peers     []PeerInfo `json:"peers"`
	Timestamp time.Time  `json:"timestamp"`
}

type PingRequest struct {
	IP   string `json:"ip"`
	Type string `json:"type,omitempty"`
}

type PingAllResponse struct {
	Results   []PingResult `json:"results"`
	Timestamp time.Time    `json:"timestamp"`
}
