package ovh

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IpReverse struct {
	IpReverse string `json:"ipReverse"`
	Reverse   string `json:"reverse"`
}

func (v IpReverse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["ip_reverse"] = v.IpReverse
	obj["reverse"] = v.Reverse
	return obj
}

type IpReverseCreateOpts struct {
	IpReverse string `json:"ipReverse"`
	Reverse   string `json:"reverse"`
}

func (opts *IpReverseCreateOpts) FromResource(d *schema.ResourceData) *IpReverseCreateOpts {
	opts.IpReverse = d.Get("ip_reverse").(string)
	opts.Reverse = d.Get("reverse").(string)

	return opts
}

type IpTask struct {
	Comment     string    `json:"comment"`
	Destination string    `json:"destination"`
	DoneDate    time.Time `json:"doneDate"`
	Function    string    `json:"function"`
	LastUpdate  time.Time `json:"lastUpdate"`
	StartDate   time.Time `json:"startDate"`
	Status      string    `json:"status"`
	TaskID      int64     `json:"taskId"`
}

type IpMoveCreateOpts struct {
	Nexthop string `json:"nexthop"`
	To      string `json:"to"`
}

func (opts *IpMoveCreateOpts) FromResource(d *schema.ResourceData) *IpMoveCreateOpts {
	opts.Nexthop = d.Get("nexthop").(string)
	opts.To = d.Get("to").(string)

	return opts
}

type IpDestinations struct {
	CloudProject    []IpDestination `json:"cloudProject"`
	DedicatedCloud  []IpDestination `json:"dedicatedCloud"`
	DedicatedServer []IpDestination `json:"dedicatedServer"`
	HostingReseller []IpDestination `json:"hostingReseller"`
	IPLoadbalancing []IpDestination `json:"ipLoadbalancing"`
	VPS             []IpDestination `json:"vps"`
}

func (v IpDestinations) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["cloud_project"] = toList(v.CloudProject)
	obj["dedicated_cloud"] = toList(v.DedicatedCloud)
	obj["dedicated_server"] = toList(v.DedicatedServer)
	obj["hosting_reseller"] = toList(v.HostingReseller)
	obj["ip_loadbalancing"] = toList(v.IPLoadbalancing)
	obj["vps"] = toList(v.VPS)
	return obj
}

func toList(ipDestionations []IpDestination) []map[string]interface{} {
	obj := []map[string]interface{}{}
	for _, v := range ipDestionations {
		obj = append(obj, v.ToMap())
	}
	return obj
}

type IpDestination struct {
	Nexthop []string `json:"nexthop"`
	Service string   `json:"service"`
}

func (v IpDestination) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["nexthop"] = v.Nexthop
	obj["service"] = v.Service
	return obj
}
