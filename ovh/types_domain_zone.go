package ovh

import ()

type DomainZone struct {
	DnssecSupported bool     `json:"dnssecSupported"`
	HasDnsAnycast   bool     `json:"hasDnsAnycast"`
	LastUpdate      string   `json:"lastUpdate"`
	Name            string   `json:"name"`
	NameServers     []string `json:"nameServers"`
}

func (v DomainZone) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["dnssec_supported"] = v.DnssecSupported
	obj["has_dns_anycast"] = v.HasDnsAnycast
	obj["last_update"] = v.LastUpdate
	obj["name"] = v.Name

	if v.NameServers != nil {
		obj["name_servers"] = v.NameServers
	}

	return obj
}

type DomainZoneConfirmTerminationOpts struct {
	Token string `json:"token"`
}
