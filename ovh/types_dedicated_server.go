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
	opts.BootId = getNilInt64PointerFromData(d, "boot_id")
	opts.Monitoring = getNilBoolPointerFromData(d, "monitoring")
	opts.State = getNilStringPointerFromData(d, "state")
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
	TemplateName        string                             `json:"templateName"`
	PartitionSchemeName *string                            `json:"partitionSchemeName,omitempty"`
	Details             *DedicatedServerInstallTaskDetails `json:"details"`
}

func (opts *DedicatedServerInstallTaskCreateOpts) FromResource(d *schema.ResourceData) *DedicatedServerInstallTaskCreateOpts {
	opts.TemplateName = d.Get("template_name").(string)
	opts.PartitionSchemeName = getNilStringPointerFromData(d, "partition_scheme_name")

	details := d.Get("details").([]interface{})
	if details != nil && len(details) == 1 {
		opts.Details = (&DedicatedServerInstallTaskDetails{}).FromResource(d, "details.0")

	}
	return opts
}

type DedicatedServerInstallTaskDetails struct {
	CustomHostname               *string `json:"customHostname,omitempty"`
	DiskGroupId                  *int64  `json:"diskGroupId,omitempty"`
	InstallRTM                   *bool   `json:"installRTM,omitempty"`
	InstallSqlServer             *bool   `json:"installSqlServer,omitempty"`
	Language                     *string `json:"language,omitempty"`
	NoRaid                       *bool   `json:"noRaid,omitempty"`
	PostInstallationScriptLink   *string `json:"postInstallationScriptLink,omitempty"`
	PostInstallationScriptReturn *string `json:"postInstallationScriptReturn,omitempty"`
	ResetHwRaid                  *bool   `json:"resetHwRaid,omitempty"`
	SoftRaidDevices              *int64  `json:"softRaidDevices,omitempty"`
	SshKeyName                   *string `json:"sshKeyName,omitempty"`
	UseDistribKernel             *bool   `json:"useDistribKernel,omitempty"`
	UseSpla                      *bool   `json:"useSpla,omitempty"`
}

func (opts *DedicatedServerInstallTaskDetails) FromResource(d *schema.ResourceData, parent string) *DedicatedServerInstallTaskDetails {
	opts.CustomHostname = getNilStringPointerFromData(d, fmt.Sprintf("%s.custom_hostname", parent))
	opts.DiskGroupId = getNilInt64PointerFromData(d, fmt.Sprintf("%s.disk_group_id", parent))
	opts.InstallRTM = getNilBoolPointerFromData(d, fmt.Sprintf("%s.install_rtm", parent))
	opts.InstallSqlServer = getNilBoolPointerFromData(d, fmt.Sprintf("%s.install_sql_server", parent))
	opts.Language = getNilStringPointerFromData(d, fmt.Sprintf("%s.language", parent))
	opts.NoRaid = getNilBoolPointerFromData(d, fmt.Sprintf("%s.no_raid", parent))
	opts.PostInstallationScriptLink = getNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script_link", parent))
	opts.PostInstallationScriptReturn = getNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script_return", parent))
	opts.ResetHwRaid = getNilBoolPointerFromData(d, fmt.Sprintf("%s.reset_hw_raid", parent))
	opts.SoftRaidDevices = getNilInt64PointerFromData(d, fmt.Sprintf("%s.soft_raid_devices", parent))
	opts.SshKeyName = getNilStringPointerFromData(d, fmt.Sprintf("%s.ssh_key_name", parent))
	opts.UseDistribKernel = getNilBoolPointerFromData(d, fmt.Sprintf("%s.use_distrib_kernel", parent))
	opts.UseSpla = getNilBoolPointerFromData(d, fmt.Sprintf("%s.use_spla", parent))

	return opts
}
