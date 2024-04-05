package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DbaasLogsInputEngine struct {
	Id           string `json:"engineId"`
	IsDeprecated bool   `json:"isDeprecated"`
	Name         string `json:"name"`
	Version      string `json:"version"`
}

func (v DbaasLogsInputEngine) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["id"] = v.Id
	obj["is_deprecated"] = v.IsDeprecated
	obj["name"] = v.Name
	obj["version"] = v.Version
	return obj
}

type DbaasLogsOperation struct {
	AliasId        *string `json:"aliasId"`
	CreatedAt      string  `json:"createdAt"`
	DashboardId    *string `json:"dashboardId"`
	IndexId        *string `json:"indexId"`
	InputId        *string `json:"inputId"`
	KibanaId       *string `json:"kibanaId"`
	OperationId    string  `json:"operationId"`
	OsdId          *string `json:"osdId"`
	RoleId         *string `json:"roleId"`
	State          string  `json:"state"`
	StreamId       *string `json:"streamId"`
	SubscriptionID *string `json:"subscriptionId"`
	UpdatedAt      string  `json:"updatedAt"`
}

type DbaasLogsOpts struct {
	ArchiveAllowedNetworks     []string `json:"archiveAllowedNetworks"`
	DirectInputAllowedNetworks []string `json:"directInputAllowedNetworks"`
	QueryAllowedNetworks       []string `json:"queryAllowedNetworks"`
}

func convertNetworks(networks []interface{}) []string {
	if networks == nil {
		return nil
	}
	networksString := make([]string, len(networks))
	for i, net := range networks {
		networksString[i] = net.(string)
	}
	return networksString
}

func (opts *DbaasLogsOpts) FromResource(d *schema.ResourceData) *DbaasLogsOpts {
	opts.ArchiveAllowedNetworks = convertNetworks(d.Get("archive_allowed_networks").(*schema.Set).List())
	opts.DirectInputAllowedNetworks = convertNetworks(d.Get("direct_input_allowed_networks").(*schema.Set).List())
	opts.QueryAllowedNetworks = convertNetworks(d.Get("query_allowed_networks").(*schema.Set).List())
	return opts
}
