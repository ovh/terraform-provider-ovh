package ovh

type DomainZone struct {
	DnssecSupported    bool     `json:"dnssecSupported"`
	HasDnsAnycast      bool     `json:"hasDnsAnycast"`
	LastUpdate         string   `json:"lastUpdate"`
	Name               string   `json:"name"`
	NameServers        []string `json:"nameServers"`
	IamResourceDetails `json:"iam"`
}

func (v DomainZone) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["dnssec_supported"] = v.DnssecSupported
	obj["has_dns_anycast"] = v.HasDnsAnycast
	obj["last_update"] = v.LastUpdate
	obj["name"] = v.Name
	obj["urn"] = v.URN

	if v.NameServers != nil {
		obj["name_servers"] = v.NameServers
	}

	return obj
}

type DomainZoneConfirmTerminationOpts struct {
	Token string `json:"token"`
}

type DomainZoneTask struct {
	CanAccelerate bool   `json:"canAccelerate"`
	CanCancel     bool   `json:"canCancel"`
	CanRelaunch   bool   `json:"canRelaunch"`
	Comment       string `json:"comment"`
	CreationDate  string `json:"creationDate"`
	DoneDate      string `json:"doneDate"`
	Function      string `json:"function"`
	TaskID        int    `json:"id"`
	LastUpdate    string `json:"lastUpdate"`
	Status        string `json:"status"`
	TodoDate      string `json:"todoDate"`
}
