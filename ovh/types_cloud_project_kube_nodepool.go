package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type CloudProjectKubeNodePoolCreateOpts struct {
	AntiAffinity  *bool   `json:"antiAffinity,omitempty"`
	Autoscale     *bool   `json:"autoscale,omitempty"`
	DesiredNodes  *int    `json:"desiredNodes,omitempty"`
	FlavorName    string  `json:"flavorName"`
	MaxNodes      *int    `json:"maxNodes,omitempty"`
	MinNodes      *int    `json:"minNodes,omitempty"`
	MonthlyBilled *bool   `json:"monthlyBilled,omitempty"`
	Name          *string `json:"name,omitempty"`
}

func (opts *CloudProjectKubeNodePoolCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeNodePoolCreateOpts {
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, "autoscale")
	opts.AntiAffinity = helpers.GetNilBoolPointerFromData(d, "anti_affinity")
	opts.DesiredNodes = helpers.GetNilIntPointerFromData(d, "desired_nodes")
	opts.FlavorName = d.Get("flavor_name").(string)
	opts.MaxNodes = helpers.GetNilIntPointerFromData(d, "max_nodes")
	opts.MinNodes = helpers.GetNilIntPointerFromData(d, "min_nodes")
	opts.MonthlyBilled = helpers.GetNilBoolPointerFromData(d, "monthly_billed")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")

	return opts
}

func (s *CloudProjectKubeNodePoolCreateOpts) String() string {
	return fmt.Sprintf("%s(%s): %d/%d/%d", *s.Name, s.FlavorName, *s.DesiredNodes, *s.MinNodes, *s.MaxNodes)
}

type CloudProjectKubeNodePoolUpdateOpts struct {
	Autoscale    *bool `json:"autoscale,omitempty"`
	DesiredNodes *int  `json:"desiredNodes,omitempty"`
	MaxNodes     *int  `json:"maxNodes,omitempty"`
	MinNodes     *int  `json:"minNodes,omitempty"`
}

func (opts *CloudProjectKubeNodePoolUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeNodePoolUpdateOpts {
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, "autoscale")
	opts.DesiredNodes = helpers.GetNilIntPointerFromData(d, "desired_nodes")
	opts.MaxNodes = helpers.GetNilIntPointerFromData(d, "max_nodes")
	opts.MinNodes = helpers.GetNilIntPointerFromData(d, "min_nodes")
	return opts
}

func (s *CloudProjectKubeNodePoolUpdateOpts) String() string {
	return fmt.Sprintf("%d/%d/%d", *s.DesiredNodes, *s.MinNodes, *s.MaxNodes)
}

type CloudProjectKubeNodePoolResponse struct {
	Autoscale      bool   `json:"autoscale"`
	AntiAffinity   bool   `json:"antiAffinity"`
	AvailableNodes int    `json:"availableNodes"`
	CreatedAt      string `json:"createdAt"`
	CurrentNodes   int    `json:"currentNodes"`
	DesiredNodes   int    `json:"desiredNodes"`
	Flavor         string `json:"flavor"`
	Id             string `json:"id"`
	MaxNodes       int    `json:"maxNodes"`
	MinNodes       int    `json:"minNodes"`
	MonthlyBilled  bool   `json:"monthlyBilled"`
	Name           string `json:"name"`
	ProjectId      string `json:"projectId"`
	SizeStatus     string `json:"sizeStatus"`
	Status         string `json:"status"`
	UpToDateNodes  int    `json:"upToDateNodes"`
	UpdatedAt      string `json:"updatedAt"`
}

func (v CloudProjectKubeNodePoolResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["anti_affinity"] = v.AntiAffinity
	obj["autoscale"] = v.Autoscale
	obj["available_nodes"] = v.AvailableNodes
	obj["created_at"] = v.CreatedAt
	obj["current_nodes"] = v.CurrentNodes
	obj["desired_nodes"] = v.DesiredNodes
	obj["flavor"] = v.Flavor
	obj["flavor_name"] = v.Flavor
	obj["id"] = v.Id
	obj["max_nodes"] = v.MaxNodes
	obj["min_nodes"] = v.MinNodes
	obj["monthly_billed"] = v.MonthlyBilled
	obj["name"] = v.Name
	obj["project_id"] = v.ProjectId
	obj["size_status"] = v.SizeStatus
	obj["status"] = v.Status
	obj["up_to_date_nodes"] = v.UpToDateNodes
	obj["updated_at"] = v.UpdatedAt

	return obj
}

func (n *CloudProjectKubeNodePoolResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}
