package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

// Opts
type CloudProjectNetworkPrivateCreateOpts struct {
	ServiceName string   `json:"serviceName"`
	VlanId      int      `json:"vlanId"`
	Name        string   `json:"name"`
	Regions     []string `json:"regions"`
}

func (p *CloudProjectNetworkPrivateCreateOpts) String() string {
	return fmt.Sprintf("projectId: %s, vlanId:%d, name: %s, regions: %s", p.ServiceName, p.VlanId, p.Name, p.Regions)
}

// Opts
type CloudProjectNetworkPrivateUpdateOpts struct {
	Name string `json:"name"`
}

type CloudProjectNetworkPrivateRegion struct {
	Status      string `json:"status"`
	Region      string `json:"region"`
	OpenStackId string `json:"openstackId"`
}

func (p *CloudProjectNetworkPrivateRegion) String() string {
	return fmt.Sprintf("Status:%s, Region: %s", p.Status, p.Region)
}

type CloudProjectNetworkPrivateResponse struct {
	Id      string                              `json:"id"`
	Status  string                              `json:"status"`
	Vlanid  int                                 `json:"vlanId"`
	Name    string                              `json:"name"`
	Type    string                              `json:"type"`
	Regions []*CloudProjectNetworkPrivateRegion `json:"regions"`
}

func (p *CloudProjectNetworkPrivateResponse) String() string {
	return fmt.Sprintf("Id: %s, Status: %s, Name: %s, Vlanid: %d, Type: %s, Regions: %s", p.Id, p.Status, p.Name, p.Vlanid, p.Type, p.Regions)
}

// Opts
type CloudProjectNetworkPrivatesCreateOpts struct {
	ServiceName string `json:"serviceName"`
	NetworkId   string `json:"networkId"`
	Dhcp        bool   `json:"dhcp"`
	NoGateway   bool   `json:"noGateway"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Network     string `json:"network"`
	Region      string `json:"region"`
}

func (p *CloudProjectNetworkPrivatesCreateOpts) String() string {
	return fmt.Sprintf("PCPNSCreateOpts[projectId: %s, networkId:%s, dhcp: %v, noGateway: %v, network: %s, start: %s, end: %s, region: %s]",
		p.ServiceName, p.NetworkId, p.Dhcp, p.NoGateway, p.Network, p.Start, p.End, p.Region)
}

type CloudProjectNetworkPrivatesResponse struct {
	Id        string    `json:"id"`
	GatewayIp string    `json:"gatewayIp"`
	Cidr      string    `json:"cidr"`
	IPPools   []*IPPool `json:"ipPools"`
}

func (p *CloudProjectNetworkPrivatesResponse) String() string {
	return fmt.Sprintf("PCPNSResponse[Id: %s, GatewayIp: %s, Cidr: %s, IPPools: %s]", p.Id, p.GatewayIp, p.Cidr, p.IPPools)
}

// Opts
type CloudProjectUserCreateOpts struct {
	Description *string  `json:"description,omitempty"`
	Role        *string  `json:"role,omitempty"`
	Roles       []string `json:"roles"`
}

func (p *CloudProjectUserCreateOpts) String() string {
	return fmt.Sprintf(
		"CloudProjectUserCreateOpts[description:%v, role:%v, roles:%s]",
		p.Description,
		p.Role,
		p.Roles,
	)
}

func (opts *CloudProjectUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectUserCreateOpts {
	opts.Description = helpers.GetNilStringPointerFromData(d, "description")
	opts.Role = helpers.GetNilStringPointerFromData(d, "role_name")
	opts.Roles, _ = helpers.StringsFromSchema(d, "role_names")
	return opts
}

type CloudProjectUser struct {
	CreationDate string                  `json:"creationDate"`
	Description  string                  `json:"description"`
	Id           int                     `json:"id"`
	Password     string                  `json:"password"`
	Roles        []*CloudProjectUserRole `json:"roles"`
	Status       string                  `json:"status"`
	Username     string                  `json:"username"`
}

func (u *CloudProjectUser) String() string {
	return fmt.Sprintf("UserResponse[Id: %v, Username: %s, Status: %s, Description: %s, CreationDate: %s]", u.Id, u.Username, u.Status, u.Description, u.CreationDate)
}

func (u CloudProjectUser) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["creation_date"] = u.CreationDate
	obj["description"] = u.Description
	//Dont set password as it must be set only at creation time
	obj["status"] = u.Status
	obj["username"] = u.Username

	// Set the roles
	var roles []map[string]interface{}
	for _, r := range u.Roles {
		roles = append(roles, r.ToMap())
	}
	obj["roles"] = roles
	return obj
}

type CloudProjectUserRole struct {
	Description string   `json:"description"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (r CloudProjectUserRole) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["description"] = r.Description
	obj["id"] = r.Id
	obj["name"] = r.Name
	obj["permissions"] = r.Permissions
	return obj
}

type CloudProjectUserOpenstackRC struct {
	Content string `json:"content"`
}

type CloudProjectRegionResponse struct {
	ContinentCode      string                       `json:"continentCode"`
	DatacenterLocation string                       `json:"datacenterLocation"`
	Name               string                       `json:"name"`
	Services           []CloudServiceStatusResponse `json:"services"`
}

func (r *CloudProjectRegionResponse) String() string {
	return fmt.Sprintf("Region: %s, Services: %s", r.Name, r.Services)
}

func (r *CloudProjectRegionResponse) HasServiceUp(service string) bool {
	for _, s := range r.Services {
		if s.Name == service && s.Status == "UP" {
			return true
		}
	}
	return false
}

type CloudServiceStatusResponse struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

func (s *CloudServiceStatusResponse) String() string {
	return fmt.Sprintf("%s: %s", s.Name, s.Status)
}

type CloudProjectKubeCreateOpts struct {
	Name             *string `json:"name,omitempty"`
	PrivateNetworkId *string `json:"privateNetworkId,omitempty"`
	Region           string  `json:"region"`
	Version          *string `json:"version,omitempty"`
}

func (opts *CloudProjectKubeCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeCreateOpts {
	opts.Region = d.Get("region").(string)
	opts.Version = helpers.GetNilStringPointerFromData(d, "version")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")
	opts.PrivateNetworkId = helpers.GetNilStringPointerFromData(d, "private_network_id")

	return opts
}

func (s *CloudProjectKubeCreateOpts) String() string {
	return fmt.Sprintf("%s(%s): %s", *s.Name, s.Region, *s.Version)
}

type CloudProjectKubeResponse struct {
	ControlPlaneIsUpToDate bool     `json:"controlPlaneIsUpToDate"`
	Id                     string   `json:"id"`
	IsUpToDate             bool     `json:"isUpToDate"`
	Name                   string   `json:"name"`
	NextUpgradeVersions    []string `json:"nextUpgradeVersions"`
	NodesUrl               string   `json:"nodesUrl"`
	PrivateNetworkId       string   `json:"privateNetworkId"`
	Region                 string   `json:"region"`
	Status                 string   `json:"status"`
	UpdatePolicy           string   `json:"updatePolicy"`
	Url                    string   `json:"url"`
	Version                string   `json:"version"`
}

func (v CloudProjectKubeResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["control_plane_is_up_to_date"] = v.ControlPlaneIsUpToDate
	obj["id"] = v.Id
	obj["is_up_to_date"] = v.IsUpToDate
	obj["name"] = v.Name
	obj["next_upgrade_versions"] = v.NextUpgradeVersions
	obj["nodes_url"] = v.NodesUrl
	obj["private_network_id"] = v.PrivateNetworkId
	obj["region"] = v.Region
	obj["status"] = v.Status
	obj["update_policy"] = v.UpdatePolicy
	obj["url"] = v.Url
	obj["version"] = v.Version[:strings.LastIndex(v.Version, ".")]

	return obj
}

func (s *CloudProjectKubeResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", s.Name, s.Id, s.Status)
}

type CloudProjectKubeKubeConfigResponse struct {
	Content string `json:"content"`
}

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
