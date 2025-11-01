export interface Machine {
	nodeKey: string;
	hostName: string;
	dnsName: string;
	tailscaleIPs: string[];
	os: string;
	online: boolean;
	sshHostKeys: string[];
	tags: string[];
	userLogin: string;
	userDisplay: string;
}

export interface MachineListResponse {
	machines: Machine[];
	self: Machine;
}
