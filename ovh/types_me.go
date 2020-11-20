package ovh

import (
	"fmt"
)

type MeIdentityUserResponse struct {
	Creation           string `json:"creation"`
	Description        string `json:"description"`
	Email              string `json:"email"`
	Group              string `json:"group"`
	LastUpdate         string `json:"lastUpdate"`
	Login              string `json:"login"`
	PasswordLastUpdate string `json:"passwordLastUpdate"`
	Status             string `json:"status"`
}

// MeIdentityUser Opts
type MeIdentityUserCreateOpts struct {
	Description string `json:"description"`
	Email       string `json:"email"`
	Group       string `json:"group"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

type MeIdentityUserUpdateOpts struct {
	Login       string `json:"user"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Group       string `json:"group"`
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
