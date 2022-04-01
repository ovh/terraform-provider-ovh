package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type HostingPrivateDatabase struct {
	ServiceName    string        `json:"serviceName"`
	Cpu            int           `json:"cpu"`
	Datacenter     string        `json:"datacenter"`
	DisplayName    string        `json:"displayName"`
	Hostname       string        `json:"hostname"`
	HostnameFtp    string        `json:"hostnameFtp"`
	Infrastructure string        `json:"infrastructure"`
	Offer          string        `json:"offer"`
	Port           int           `json:"port"`
	PortFtp        int           `json:"portFtp"`
	QuotaSize      *UnitAndValue `json:"quotasize"`
	QuotaUsed      *UnitAndValue `json:"quotaUsed"`
	Ram            *UnitAndValue `json:"ram"`
	Server         string        `json:"server"`
	State          string        `json:"state"`
	Type           string        `json:"type"`
	Version        string        `json:"version"`
	VersionLabel   string        `json:"versionLabel"`
	VersionNumber  float64       `json:"versionNumber"`
}

func (v HostingPrivateDatabase) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["service_name"] = v.ServiceName
	obj["cpu"] = v.Cpu
	obj["datacenter"] = v.Datacenter
	obj["display_name"] = v.DisplayName
	obj["hostname"] = v.Hostname
	obj["hostname_ftp"] = v.HostnameFtp
	obj["infrastructure"] = v.Infrastructure
	obj["offer"] = v.Offer
	obj["port"] = v.Port
	obj["port_ftp"] = v.PortFtp
	obj["quota_size"] = v.QuotaSize.Value
	obj["quota_used"] = v.QuotaUsed.Value
	obj["ram"] = v.Ram.Value
	obj["state"] = v.State
	obj["type"] = v.Type
	obj["version"] = v.Version
	obj["version_label"] = v.VersionLabel
	obj["version_number"] = v.VersionNumber

	return obj
}

type HostingPrivateDatabaseOpts struct {
	DisplayName *string `json:"displayName"`
}

func (opts *HostingPrivateDatabaseOpts) FromResource(d *schema.ResourceData) *HostingPrivateDatabaseOpts {
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")

	return opts
}

type HostingPrivateDatabaseConfirmTerminationOpts struct {
	Token string `json:"token"`
}
