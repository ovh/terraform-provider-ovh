package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectVrackResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (v *CloudProjectVrackResponse) ToMap(d *schema.ResourceData) map[string]interface{} {
	obj := make(map[string]interface{})
	obj["id"] = v.Id
	obj["name"] = v.Name
	obj["description"] = v.Description
	return obj
}

// Opts
type CloudProjectNetworkPrivateCreateOpts struct {
	ServiceName string   `json:"serviceName"`
	VlanId      int      `json:"vlanId"`
	Name        string   `json:"name"`
	Regions     []string `json:"regions"`
}

type CloudProjectNetworkPrivateUpdateOptsAlone struct {
	Region string `json:"region"`
}

func (p *CloudProjectNetworkPrivateCreateOpts) String() string {
	return fmt.Sprintf("projectId: %s, vlanId:%d, name: %s, regions: %s", p.ServiceName, p.VlanId, p.Name, p.Regions)
}

// Opts
type CloudProjectNetworkPrivateUpdateOpts struct {
	Name    string   `json:"name"`
	Regions []string `json:"regions"`
}

type CloudProjectNetworkPrivateRegion struct {
	Status      string `json:"status"`
	Region      string `json:"region"`
	OpenStackId string `json:"openstackId"`
}
