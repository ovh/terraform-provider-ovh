package ovh

// VPSSecondaryDNSDomain mirrors secondaryDns.SecondaryDNS from the OVH API.
type VPSSecondaryDNSDomain struct {
	CreationDate string `json:"creationDate"`
	Dns          string `json:"dns"`
	Domain       string `json:"domain"`
	IpMaster     string `json:"ipMaster"`
}

// VPSSecondaryDNSDomainCreateOpts is the body for
// POST /vps/{serviceName}/secondaryDnsDomains.
type VPSSecondaryDNSDomainCreateOpts struct {
	Domain string  `json:"domain"`
	Ip     *string `json:"ip,omitempty"`
}

// VPSSecondaryDNSNameServer mirrors secondaryDns.SecondaryDNSNameServer.
type VPSSecondaryDNSNameServer struct {
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	Ipv6     string `json:"ipv6,omitempty"`
}
