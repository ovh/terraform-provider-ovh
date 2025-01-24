package ovh

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type DedicatedServer struct {
	AvailabilityZone   string `json:"availabilityZone"`
	Name               string `json:"name"`
	BootId             int    `json:"bootId"`
	BootScript         string `json:"bootScript"`
	CommercialRange    string `json:"commercialRange"`
	Datacenter         string `json:"datacenter"`
	Ip                 string `json:"ip"`
	LinkSpeed          int    `json:"linkSpeed"`
	Monitoring         bool   `json:"monitoring"`
	NewUpgradeSystem   bool   `json:"newUpgradeSystem"`
	NoIntervention     bool   `json:"noIntervention"`
	Os                 string `json:"os"`
	PowerState         string `json:"powerState"`
	ProfessionalUse    bool   `json:"professionalUse"`
	Rack               string `json:"rack"`
	Region             string `json:"region"`
	RescueMail         string `json:"rescueMail"`
	RescueSshKey       string `json:"rescueSshKey"`
	Reverse            string `json:"reverse"`
	RootDevice         string `json:"rootDevice"`
	ServerId           int    `json:"serverId"`
	State              string `json:"state"`
	SupportLevel       string `json:"supportLevel"`
	IamResourceDetails `json:"iam"`
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
	BootScript *string `json:"bootScript,omitempty"`
	Monitoring *bool   `json:"monitoring,omitempty"`
	State      *string `json:"state,omitempty"`
}

func (opts *DedicatedServerUpdateOpts) FromResource(d *schema.ResourceData) *DedicatedServerUpdateOpts {
	opts.BootId = helpers.GetNilInt64PointerFromData(d, "boot_id")
	opts.BootScript = helpers.GetNilStringPointerFromData(d, "boot_script")
	opts.Monitoring = helpers.GetNilBoolPointerFromData(d, "monitoring")
	opts.State = helpers.GetNilStringPointerFromData(d, "state")
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
	DoneDate   time.Time `json:"doneDate"`
	StartDate  time.Time `json:"startDate"`
}

type DedicatedServerInstallTaskCreateOpts struct {
	OperatingSystem string `json:"operatingSystem"`
}

func (opts *DedicatedServerInstallTaskCreateOpts) FromResource(d *schema.ResourceData) *DedicatedServerInstallTaskCreateOpts {
	opts.OperatingSystem = d.Get("operating_system").(string)

	return opts
}
