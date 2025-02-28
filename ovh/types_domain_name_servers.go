package ovh

type DomainNameServers struct {
	Domain  string             `json:"domain"`
	Servers []DomainNameServer `json:"servers"`
}

func (v DomainNameServers) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["domain"] = v.Domain
	obj["servers"] = []map[string]interface{}{}

	for _, server := range v.Servers {
		obj["servers"] = append(obj["servers"].([]map[string]interface{}), server.ToMap())
	}

	return obj
}

type DomainNameServer struct {
	Id       int    `json:"id,omitempty"`
	Host     string `json:"host"`
	Ip       string `json:"ip,omitempty"`
	IsUsed   bool   `json:"isUsed,omitempty"`
	ToDelete bool   `json:"toDelete,omitempty"`
}

func (v DomainNameServer) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["host"] = v.Host
	obj["ip"] = v.Ip

	return obj
}

type DomainNameServerUpdateOpts struct {
	NameServers []DomainNameServer `json:"nameServers"`
}

type DomainNameServerTypeOpts struct {
	NameServerType string `json:"nameServerType"`
}
