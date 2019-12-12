package ovh

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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

type DedicatedServerUpdateOpts struct {
	BootId     *int64  `json:"bootId,omitempty"`
	Monitoring *bool   `json:"monitoring,omitempty"`
	State      *string `json:"state,omitempty"`
}

func (opts *DedicatedServerUpdateOpts) FromResource(d *schema.ResourceData) *DedicatedServerUpdateOpts {
	opts.BootId = getNilInt64Pointer(d.Get("boot_id"))
	opts.Monitoring = getNilBoolPointer(d.Get("monitoring"))
	opts.State = getNilStringPointer(d.Get("state"))
	return opts
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

type DedicatedServerBoot struct {
	BootId      int    `json:"bootId"`
	BootType    string `json:"bootType"`
	Description string `json:"description"`
	Kernel      string `json:"kernel"`
}

type DedicatedServerTask struct {
	Id         int64     `json:"taskId"`
	Function   string    `json:"function"`
	Comment    string    `json:"comment"`
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"lastUpdate"`
	doneDate   time.Time `json:"doneDate"`
	startDate  time.Time `json:"startDate"`
}
