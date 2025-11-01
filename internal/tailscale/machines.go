package tailscale

import (
	"context"
	"net/netip"
)

type Machine struct {
	NodeKey      string   `json:"nodeKey"`
	HostName     string   `json:"hostName"`
	DNSName      string   `json:"dnsName"`
	TailscaleIPs []string `json:"tailscaleIPs"`
	OS           string   `json:"os"`
	Online       bool     `json:"online"`
	SSHHostKeys  []string `json:"sshHostKeys"`
	Tags         []string `json:"tags"`
	UserLogin    string   `json:"userLogin"`
	UserDisplay  string   `json:"userDisplay"`
}

type MachineListResponse struct {
	Machines []Machine `json:"machines"`
	Self     Machine   `json:"self"`
}

func (tc *TailscaleClient) GetSSHMachines(ctx context.Context) (*MachineListResponse, error) {
	status, err := tc.lc.Status(ctx)
	if err != nil {
		return nil, err
	}

	var machines []Machine
	for _, peer := range status.Peer {
		if len(peer.SSH_HostKeys) > 0 {
			tags := []string{}
			if peer.Tags != nil {
				tags = peer.Tags.AsSlice()
			}

			userLogin := ""
			userDisplay := ""
			if userProfile, ok := status.User[peer.UserID]; ok {
				userLogin = userProfile.LoginName
				userDisplay = userProfile.DisplayName
			}

			machines = append(machines, Machine{
				NodeKey:      peer.PublicKey.String(),
				HostName:     peer.HostName,
				DNSName:      peer.DNSName,
				TailscaleIPs: formatIPs(peer.TailscaleIPs),
				OS:           peer.OS,
				Online:       peer.Online,
				SSHHostKeys:  peer.SSH_HostKeys,
				Tags:         tags,
				UserLogin:    userLogin,
				UserDisplay:  userDisplay,
			})
		}
	}

	self := Machine{}
	if status.Self != nil {
		selfTags := []string{}
		if status.Self.Tags != nil {
			selfTags = status.Self.Tags.AsSlice()
		}

		selfUserLogin := ""
		selfUserDisplay := ""
		if userProfile, ok := status.User[status.Self.UserID]; ok {
			selfUserLogin = userProfile.LoginName
			selfUserDisplay = userProfile.DisplayName
		}

		self = Machine{
			NodeKey:      status.Self.PublicKey.String(),
			HostName:     status.Self.HostName,
			DNSName:      status.Self.DNSName,
			TailscaleIPs: formatIPs(status.Self.TailscaleIPs),
			OS:           status.Self.OS,
			Online:       true,
			SSHHostKeys:  status.Self.SSH_HostKeys,
			Tags:         selfTags,
			UserLogin:    selfUserLogin,
			UserDisplay:  selfUserDisplay,
		}
	}

	return &MachineListResponse{
		Machines: machines,
		Self:     self,
	}, nil
}

func formatIPs(ips []netip.Addr) []string {
	result := make([]string, len(ips))
	for i, ip := range ips {
		result[i] = ip.String()
	}
	return result
}
