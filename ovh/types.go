package ovh

import (
	"fmt"
)

type IPPool struct {
	Network string `json:"network"`
	Region  string `json:"region"`
	Dhcp    bool   `json:"dhcp"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

func (p *IPPool) String() string {
	return fmt.Sprintf("IPPool[Network: %s, Region: %s, Dhcp: %v, Start: %s, End: %s]", p.Network, p.Region, p.Dhcp, p.Start, p.End)
}

// Task Opts
type TaskOpts struct {
	ServiceName string `json:"serviceName"`
	TaskId      string `json:"taskId"`
}

type UnitAndValue struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

func (v UnitAndValue) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["unit"] = v.Unit
	obj["value"] = v.Value

	return obj
}
