package ovh

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type Vrack struct {
	Description        *string `json:"description"`
	Name               *string `json:"name"`
	IamResourceDetails `json:"iam"`
}

func (v Vrack) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	if v.Description != nil {
		obj["description"] = *v.Description
	}

	if v.Name != nil {
		obj["name"] = *v.Name
	}

	obj["urn"] = v.URN

	return obj
}

type VrackUpdateOpts struct {
	Description *string `json:"description,omitempty"`
	Name        *string `json:"name,omitempty"`
}

func (opts *VrackUpdateOpts) FromResource(d *schema.ResourceData) *VrackUpdateOpts {
	opts.Description = helpers.GetNilStringPointerFromData(d, "description")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")

	return opts
}

type VrackIp struct {
	Gateway string `json:"gateway"`
	Ip      string `json:"ip"`
	Zone    string `json:"zone"`
	Region  string `json:"region"`
}

func (v VrackIp) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["gateway"] = v.Gateway
	obj["ip"] = v.Ip
	obj["zone"] = v.Zone
	obj["region"] = v.Region

	return obj
}

// IPv4 Block Payload
type VrackIpCreateOpts struct {
	Block  string  `json:"block"`
	Region *string `json:"region,omitempty"`
}

func (opts *VrackIpCreateOpts) FromResource(d *schema.ResourceData) *VrackIpCreateOpts {
	opts.Block = d.Get("block").(string)
	opts.Region = helpers.GetNilStringPointerFromData(d, "region")
	return opts
}

type VrackIpV6 struct {
	Region string `json:"region"`
	IPv6   string `json:"ipv6"`
}

func (v VrackIpV6) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["region"] = v.Region
	obj["ipv6"] = v.IPv6
	return obj
}

// IPv6 Block Payload
type VrackIpV6CreateOpts struct {
	Block string `json:"block"`
}

func (opts *VrackIpV6CreateOpts) FromResource(d *schema.ResourceData) *VrackIpV6CreateOpts {
	opts.Block = d.Get("block").(string)
	return opts
}

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

type VrackVrackServicesCreateOpts struct {
	VrackServices string `json:"vrackServices"`
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

// IPv6 - Bridged subrange

type VrackIPv6BridgedSubrange struct {
	BridgedSubrange string `json:"bridgedSubrange"`
	Gateway         string `json:"gateway"`
	Slaac           string `json:"slaac"`
}

func (v VrackIPv6BridgedSubrange) ToMap() []map[string]interface{} {
	return []map[string]interface{}{{
		"subrange": v.BridgedSubrange,
		"gateway":  v.Gateway,
		"slaac":    v.Slaac,
	}}
}

type VrackIPv6BridgedSubrangeSlaacUpdateOpts struct {
	Slaac string `json:"slaac"`
}

func (opts *VrackIPv6BridgedSubrangeSlaacUpdateOpts) FromResource(d *schema.ResourceData) *VrackIPv6BridgedSubrangeSlaacUpdateOpts {
	slaac := "disabled"
	st := d.Get("bridged_subrange").(*schema.Set)
	if st != nil {
		if list := st.List(); len(list) == 1 {
			elems := list[0].(map[string]interface{})
			slaac = elems["slaac"].(string)
		}
	}
	opts.Slaac = slaac
	return opts
}
