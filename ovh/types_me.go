package ovh

import (
	"fmt"
)

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

// MeIpxeScript Opts
type MeIpxeScriptCreateOpts struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Script      string `json:"script"`
}

type MeIpxeScriptResponse struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

func (s *MeIpxeScriptResponse) String() string {
	return fmt.Sprintf("IpxeScript: %s", s.Name)
}
