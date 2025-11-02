export type ConnectionType = 'direct' | 'derp' | 'peer-relay' | 'offline' | 'unknown';

export interface PeerInfo {
	hostName: string;
	dnsName: string;
	ip: string;
	online: boolean;
	active: boolean;
	curAddr?: string;
	relay?: string;
	peerRelay?: string;
	lastHandshake: string;
	rxBytes: number;
	txBytes: number;
	os: string;
	userLogin?: string;
	userDisplay?: string;
	tags?: string[];
}

export interface PingResult {
	ip: string;
	nodeName: string;
	success: boolean;
	error?: string;
	latencyMs: number;
	connectionType: ConnectionType;
	endpoint?: string;
	derpRegion?: string;
	derpRegionId?: number;
	peerRelay?: string;
}

export interface PeersResponse {
	peers: PeerInfo[];
	timestamp: string;
}

export interface PingRequest {
	ip: string;
	type?: string;
}

export interface PingAllResponse {
	results: PingResult[];
	timestamp: string;
}
