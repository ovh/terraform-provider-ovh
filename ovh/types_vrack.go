package ovh

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type VrackDedicatedServerInterface struct {
	Vrack                    string `json:"vrack"`
	DedicatedServerInterface string `json:"dedicatedServerInterface"`
}

type VrackDedicatedServerInterfaceCreateOpts struct {
	DedicatedServerInterface string `json:"dedicatedServerInterface"`
}

func (opts *VrackDedicatedServerInterfaceCreateOpts) FromResource(d *schema.ResourceData) *VrackDedicatedServerInterfaceCreateOpts {
	opts.DedicatedServerInterface = d.Get("interface_id").(string)
	return opts
}

type VrackDedicatedServer struct {
	Vrack           string `json:"vrack"`
	DedicatedServer string `json:"dedicatedServer"`
}

type VrackDedicatedServerCreateOpts struct {
	DedicatedServer string `json:"dedicatedServer"`
}

func (opts *VrackDedicatedServerCreateOpts) FromResource(d *schema.ResourceData) *VrackDedicatedServerCreateOpts {
	opts.DedicatedServer = d.Get("server_id").(string)
	return opts
}

type VrackCloudProject struct {
	Vrack   string `json:"vrack"`
	Project string `json:"project"`
}

type VrackCloudProjectCreateOpts struct {
	Project string `json:"project"`
}

func (opts *VrackCloudProjectCreateOpts) FromResource(d *schema.ResourceData) *VrackCloudProjectCreateOpts {
	opts.Project = d.Get("project_id").(string)
	return opts
}

type VrackIpLoadbalancing struct {
	Vrack           string `json:"vrack"`
	IpLoadbalancing string `json:"ipLoadbalancing"`
}

type VrackIpLoadbalancingCreateOpts struct {
	IpLoadbalancing string `json:"ipLoadbalancing"`
}

func (opts *VrackIpLoadbalancingCreateOpts) FromResource(d *schema.ResourceData) *VrackIpLoadbalancingCreateOpts {
	opts.IpLoadbalancing = d.Get("ip_loadbalancing").(string)
	return opts
}

type VrackTask struct {
	Id           int       `json:"id"`
	Function     string    `json:"function"`
	TargetDomain string    `json:"targetDomain"`
	Status       string    `json:"status"`
	ServiceName  string    `json:"serviceName"`
	OrderId      int       `json:"orderId"`
	LastUpdate   time.Time `json:"lastUpdate"`
	TodoDate     time.Time `json:"TodoDate"`
}
