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

type DataSourceHostingPrivateDatabaseDatabase struct {
	BackupTime   string                                           `json:"backupTime"`
	QuotaUsed    *UnitAndValue                                    `json:"quotaUsed"`
	CreationDate string                                           `json:"creationDate"`
	Users        []*DataSourceHostingPrivateDatabaseDatabaseUsers `json:"users"`
}

func (v DataSourceHostingPrivateDatabaseDatabase) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["backup_time"] = v.BackupTime
	obj["quota_used"] = v.QuotaUsed.Value
	obj["creation_date"] = v.CreationDate

	var users []map[string]interface{}
	for _, r := range v.Users {
		users = append(users, r.ToMap())
	}
	obj["users"] = users
	return obj
}

type DataSourceHostingPrivateDatabaseDatabaseUsers struct {
	UserName  string `json:"userName"`
	GrantType string `json:"grantType"`
}

func (v DataSourceHostingPrivateDatabaseDatabaseUsers) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["user_name"] = v.UserName
	obj["grant_type"] = v.GrantType
	return obj
}

type HostingPrivateDatabaseDatabase struct {
	DoneDate     string `json:"doneDate"`
	LastUpdate   string `json:"lastUpdate"`
	UserName     string `json:"userName"`
	DumpId       string `json:"dumpId"`
	DatabaseName string `json:"databaseName"`
	TaskId       int    `json:"id"`
	StartDate    string `json:"startDate"`
	Status       string `json:"status"`
}

func (v HostingPrivateDatabaseDatabase) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["database_name"] = v.DatabaseName
	return obj
}

type HostingPrivateDatabaseDatabaseCreateOpts struct {
	DatabaseName string `json:"databaseName"`
}

func (opts *HostingPrivateDatabaseDatabaseCreateOpts) FromResource(d *schema.ResourceData) *HostingPrivateDatabaseDatabaseCreateOpts {
	opts.DatabaseName = d.Get("database_name").(string)

	return opts
}

type DataSourceHostingPrivateDatabaseUser struct {
	CreationDate string                                           `json:"creationDate"`
	Databases    []*DataSourceHostingPrivateDatabaseUserDatabases `json:"databases"`
}

func (v DataSourceHostingPrivateDatabaseUser) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["creation_date"] = v.CreationDate

	var databases []map[string]interface{}
	for _, r := range v.Databases {
		databases = append(databases, r.ToMap())
	}
	obj["databases"] = databases
	return obj
}

type DataSourceHostingPrivateDatabaseUserDatabases struct {
	DatabaseName string `json:"databaseName"`
	GrantType    string `json:"grantType"`
}

func (v DataSourceHostingPrivateDatabaseUserDatabases) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["database_name"] = v.DatabaseName
	obj["grant_type"] = v.GrantType
	return obj
}

type HostingPrivateDatabaseUser struct {
	LastUpdate   string `json:"lastUpdate"`
	DoneDate     string `json:"doneDate"`
	Status       string `json:"status"`
	StartDate    string `json:"startDate"`
	DatabaseName string `json:"databaseName"`
	UserName     string `json:"userName"`
	TaskId       int    `json:"id"`
}

func (v HostingPrivateDatabaseUser) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["user_name"] = v.UserName
	return obj
}

type HostingPrivateDatabaseUserCreateOpts struct {
	Password string `json:"password"`
	UserName string `json:"userName"`
}

func (opts *HostingPrivateDatabaseUserCreateOpts) FromResource(d *schema.ResourceData) *HostingPrivateDatabaseUserCreateOpts {
	opts.Password = d.Get("password").(string)
	opts.UserName = d.Get("user_name").(string)

	return opts
}

type HostingPrivateDatabaseUserGrant struct {
	LastUpdate   string `json:"lastUpdate"`
	DoneDate     string `json:"doneDate"`
	Status       string `json:"status"`
	StartDate    string `json:"startDate"`
	DatabaseName string `json:"databaseName"`
	UserName     string `json:"userName"`
	TaskId       int    `json:"id"`
}

type HostingPrivateDatabaseUserGrantCreateOpts struct {
	DatabaseName string `json:"databaseName"`
	Grant        string `json:"grant"`
}

func (opts *HostingPrivateDatabaseUserGrantCreateOpts) FromResource(d *schema.ResourceData) *HostingPrivateDatabaseUserGrantCreateOpts {
	opts.DatabaseName = d.Get("database_name").(string)
	opts.Grant = d.Get("grant").(string)

	return opts
}

func (v HostingPrivateDatabaseUserGrantCreateOpts) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["grant"] = v.Grant
	obj["database_name"] = v.DatabaseName
	return obj
}

type DataSourceHostingPrivateDatabaseUserGrant struct {
	CreationDate string `json:"creationDate"`
	Grant        string `json:"grant"`
}

func (v DataSourceHostingPrivateDatabaseUserGrant) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["creation_date"] = v.CreationDate
	obj["grant"] = v.Grant

	return obj
}

type HostingPrivateDatabaseWhitelist struct {
	CreationDate string `json:"creationDate"`
	LastUpdate   string `json:"lastUpdate"`
	Name         string `json:"name"`
	Service      bool   `json:"service"`
	Sftp         bool   `json:"sftp"`
	Status       string `json:"status"`
	TaskId       int    `json:"id"`
}

func (v HostingPrivateDatabaseWhitelist) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = v.Name
	obj["service"] = v.Service
	obj["sftp"] = v.Sftp

	return obj
}

func (v HostingPrivateDatabaseWhitelist) DataSourceToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = v.Name
	obj["service"] = v.Service
	obj["sftp"] = v.Sftp
	obj["creation_date"] = v.CreationDate
	obj["last_update"] = v.LastUpdate
	obj["status"] = v.Status

	return obj
}

type HostingPrivateDatabaseWhitelistCreateOpts struct {
	Ip      string `json:"ip"`
	Name    string `json:"name"`
	Service bool   `json:"service"`
	Sftp    bool   `json:"sftp"`
}

func (opts *HostingPrivateDatabaseWhitelistCreateOpts) FromResource(d *schema.ResourceData) *HostingPrivateDatabaseWhitelistCreateOpts {
	opts.Ip = d.Get("ip").(string)
	opts.Name = d.Get("name").(string)
	opts.Service = d.Get("service").(bool)
	opts.Sftp = d.Get("sftp").(bool)

	return opts
}

type HostingPrivateDatabaseWhitelistUpdateOpts struct {
	Name    string `json:"name"`
	Service bool   `json:"service"`
	Sftp    bool   `json:"sftp"`
}

func (opts *HostingPrivateDatabaseWhitelistUpdateOpts) FromResource(d *schema.ResourceData) *HostingPrivateDatabaseWhitelistUpdateOpts {
	opts.Name = d.Get("name").(string)
	opts.Service = d.Get("service").(bool)
	opts.Sftp = d.Get("sftp").(bool)

	return opts
}
