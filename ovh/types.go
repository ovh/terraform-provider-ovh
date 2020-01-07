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

// MeSshKey Opts
type MeSshKeyCreateOpts struct {
	KeyName string `json:"keyName"`
	Key     string `json:"key"`
}

type MeSshKeyResponse struct {
	KeyName string `json:"keyName"`
	Key     string `json:"key"`
	Default bool   `json:"default"`
}

func (s *MeSshKeyResponse) String() string {
	return fmt.Sprintf("SSH Key: %s, key:%s, default:%t",
		s.Key, s.KeyName, s.Default)
}

type MeSshKeyUpdateOpts struct {
	Default bool `json:"default"`
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
