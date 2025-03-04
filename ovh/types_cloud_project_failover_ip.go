package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type FailoverIp struct {
	Block         string `json:"block"`
	ContinentCode string `json:"continentCode"`
	GeoLoc        string `json:"geoLoc"`
	Id            string `json:"id"`
	Ip            string `json:"ip"`
	Progress      *int64 `json:"progress"`
	RoutedTo      string `json:"routedTo"`
	Status        string `json:"status"`
	SubType       string `json:"subType"`
}

func (v FailoverIp) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["block"] = v.Block
	obj["continent_code"] = v.ContinentCode
	obj["geo_loc"] = v.GeoLoc
	obj["id"] = v.Id
	obj["ip"] = v.Ip
	obj["progress"] = v.Progress
	obj["routed_to"] = v.RoutedTo
	obj["status"] = v.Status
	obj["sub_type"] = v.SubType

	return obj
}

type ProjectIpFailoverAttachCreation struct {
	InstanceId *string `json:"instanceId,omitempty"`
}

func (opts *ProjectIpFailoverAttachCreation) FromResource(d *schema.ResourceData) *ProjectIpFailoverAttachCreation {
	opts.InstanceId = helpers.GetNilStringPointerFromData(d, "routed_to")
	return opts
}
