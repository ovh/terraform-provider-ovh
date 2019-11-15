package ovh

import "fmt"

type DedicatedServer struct {
	Name            string `json:"name"`
	BootId          int    `json:"bootId"`
	CommercialRange string `json:"commercialRange"`
	Datacenter      string `json:"datacenter"`
	Ip              string `json:"ip"`
	LinkSpeed       int    `json:"linkSpeed"`
	Monitoring      bool   `json:"monitoring"`
	Os              string `json:"os"`
	ProfessionalUse bool   `json:"professionalUse"`
	Rack            string `json:"rack"`
	RescueMail      string `json:"rescueMail"`
	Reverse         string `json:"reverse"`
	RootDevice      string `json:"rootDevice"`
	ServerId        int    `json:"serverId"`
	State           string `json:"state"`
	SupportLevel    string `json:"supportLevel"`
}

func (ds DedicatedServer) String() string {
	return fmt.Sprintf(
		"name: %v, ip: %v, dc: %v, state: %v",
		ds.Name,
		ds.Ip,
		ds.Datacenter,
		ds.State,
	)
}

type DedicatedServerVNI struct {
	Enabled    bool     `json:"enabled"`
	Mode       string   `json:"mode"`
	Name       string   `json:"name"`
	NICs       []string `json:"networkInterfaceController"`
	ServerName string   `json:"serverName"`
	Uuid       string   `json:"uuid"`
	Vrack      *string  `json:"vrack,omitempty"`
}

func (vni DedicatedServerVNI) String() string {
	return fmt.Sprintf(
		"name: %v, uuid: %v, mode: %v, enabled: %v",
		vni.Name,
		vni.Uuid,
		vni.Mode,
		vni.Enabled,
	)
}
func (v DedicatedServerVNI) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["enabled"] = v.Enabled
	obj["mode"] = v.Mode
	obj["name"] = v.Name
	obj["nics"] = v.NICs
	obj["server_name"] = v.ServerName
	obj["uuid"] = v.Uuid

	if v.Vrack != nil {
		obj["vrack"] = *v.Vrack
	}

	return obj
}
