package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"
)

// Opts
type CloudNetworkPrivateCreateOpts struct {
	ProjectId string   `json:"serviceName"`
	VlanId    int      `json:"vlanId"`
	Name      string   `json:"name"`
	Regions   []string `json:"regions"`
}

func (p *CloudNetworkPrivateCreateOpts) String() string {
	return fmt.Sprintf("projectId: %s, vlanId:%d, name: %s, regions: %s", p.ProjectId, p.VlanId, p.Name, p.Regions)
}

// Opts
type CloudNetworkPrivateUpdateOpts struct {
	Name string `json:"name"`
}

type CloudNetworkPrivateRegion struct {
	Status string `json:"status"`
	Region string `json:"region"`
}

func (p *CloudNetworkPrivateRegion) String() string {
	return fmt.Sprintf("Status:%s, Region: %s", p.Status, p.Region)
}

type CloudNetworkPrivateResponse struct {
	Id      string                       `json:"id"`
	Status  string                       `json:"status"`
	Vlanid  int                          `json:"vlanId"`
	Name    string                       `json:"name"`
	Type    string                       `json:"type"`
	Regions []*CloudNetworkPrivateRegion `json:"regions"`
}

func (p *CloudNetworkPrivateResponse) String() string {
	return fmt.Sprintf("Id: %s, Status: %s, Name: %s, Vlanid: %d, Type: %s, Regions: %s", p.Id, p.Status, p.Name, p.Vlanid, p.Type, p.Regions)
}

// Opts
type CloudNetworkPrivatesCreateOpts struct {
	ProjectId string `json:"serviceName"`
	NetworkId string `json:"networkId"`
	Dhcp      bool   `json:"dhcp"`
	NoGateway bool   `json:"noGateway"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Network   string `json:"network"`
	Region    string `json:"region"`
}

func (p *CloudNetworkPrivatesCreateOpts) String() string {
	return fmt.Sprintf("PCPNSCreateOpts[projectId: %s, networkId:%s, dhcp: %v, noGateway: %v, network: %s, start: %s, end: %s, region: %s]",
		p.ProjectId, p.NetworkId, p.Dhcp, p.NoGateway, p.Network, p.Start, p.End, p.Region)
}

type CloudNetworkPrivatesResponse struct {
	Id        string    `json:"id"`
	GatewayIp string    `json:"gatewayIp"`
	Cidr      string    `json:"cidr"`
	IPPools   []*IPPool `json:"ipPools"`
}

func (p *CloudNetworkPrivatesResponse) String() string {
	return fmt.Sprintf("PCPNSResponse[Id: %s, GatewayIp: %s, Cidr: %s, IPPools: %s]", p.Id, p.GatewayIp, p.Cidr, p.IPPools)
}

// Opts
type CloudUserCreateOpts struct {
	Description *string  `json:"description,omitempty"`
	Role        *string  `json:"role,omitempty"`
	Roles       []string `json:"roles"`
}

func (p *CloudUserCreateOpts) String() string {
	return fmt.Sprintf(
		"CloudUserCreateOpts[description:%v, role:%v, roles:%s]",
		p.Description,
		p.Role,
		p.Roles,
	)
}

func (opts *CloudUserCreateOpts) FromResource(d *schema.ResourceData) *CloudUserCreateOpts {
	opts.Description = helpers.GetNilStringPointerFromData(d, "description")
	opts.Role = helpers.GetNilStringPointerFromData(d, "role_name")
	opts.Roles, _ = helpers.StringsFromSchema(d, "role_names")
	return opts
}

type CloudUser struct {
	CreationDate string           `json:"creationDate"`
	Description  string           `json:"description"`
	Id           int              `json:"id"`
	Password     string           `json:"password"`
	Roles        []*CloudUserRole `json:"roles"`
	Status       string           `json:"status"`
	Username     string           `json:"username"`
}

func (u *CloudUser) String() string {
	return fmt.Sprintf("UserResponse[Id: %v, Username: %s, Status: %s, Description: %s, CreationDate: %s]", u.Id, u.Username, u.Status, u.Description, u.CreationDate)
}

func (u CloudUser) ToMap() map[string]interface{} {
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

type CloudUserRole struct {
	Description string   `json:"description"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (r CloudUserRole) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["description"] = r.Description
	obj["id"] = r.Id
	obj["name"] = r.Name
	obj["permissions"] = r.Permissions
	return obj
}

type CloudUserOpenstackRC struct {
	Content string `json:"content"`
}

type CloudRegionResponse struct {
	ContinentCode      string                       `json:"continentCode"`
	DatacenterLocation string                       `json:"datacenterLocation"`
	Name               string                       `json:"name"`
	Services           []CloudServiceStatusResponse `json:"services"`
}

func (r *CloudRegionResponse) String() string {
	return fmt.Sprintf("Region: %s, Services: %s", r.Name, r.Services)
}

func (r *CloudRegionResponse) HasServiceUp(service string) bool {
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

type PublicCloudKubernetesClusterResponse struct {
	ControlPlaneIsUpToDate bool     `json:"controlPlaneIsUpToDate"`
	Id                     string   `json:"id"`
	IsUpToDate             bool     `json:"isUpToDate"`
	Name                   string   `json:"name"`
	NextUpgradeVersions    []string `json:"nextUpgradeVersions"`
	NodesUrl               string   `json:"nodesUrl"`
	Region                 string   `json:"region"`
	Status                 string   `json:"status"`
	UpdatePolicy           string   `json:"updatePolicy"`
	Url                    string   `json:"url"`
	Version                string   `json:"version"`
}

func (s *PublicCloudKubernetesClusterResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", s.Name, s.Id, s.Status)
}

type CloudKubernetesKubeConfigResponse struct {
	Content string `json:"content"`
}

type CloudKubernetesNodeResponse struct {
	Id         string `json:"id"`
	ProjectId  string `json:"projectId"`
	InstanceId string `json:"instanceId"`
	IsUpToDate bool   `json:"isUpToDate"`
	Name       string `json:"name"`
	Flavor     string `json:"flavor"`
	Status     string `json:"status"`
	Version    string `json:"version"`
}

func (n *CloudKubernetesNodeResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}

type CloudKubernetesNodeCreationRequest struct {
	FlavorName string `json:"flavorName"`
	Name       string `json:"name"`
}

type CloudKubernetesNodePoolCreationRequest struct {
	DesiredNodes  int    `json:"desiredNodes"`
	MaxNodes      int    `json:"maxNodes"`
	MinNodes      int    `json:"minNodes"`
	FlavorName    string `json:"flavorName"`
	Name          string `json:"name"`
	MonthlyBilled bool   `json:"monthlyBilled"`
	AntiAffinity  bool   `json:"antiAffinity"`
}

type CloudKubernetesNodePoolUpdateRequest struct {
	DesiredNodes int `json:"desiredNodes"`
	MaxNodes     int `json:"maxNodes"`
	MinNodes     int `json:"minNodes"`
}

type CloudKubernetesNodePoolResponse struct {
	Id             string `json:"id"`
	ProjectId      string `json:"projectId"`
	Name           string `json:"name"`
	AntiAffinity   bool   `json:"antiAffinity"`
	AvailableNodes int    `json:"availableNodes"`
	CreatedAt      string `json:"createdAt"`
	CurrentNodes   int    `json:"currentNodes"`
	DesiredNodes   int    `json:"desiredNodes"`
	Flavor         string `json:"flavor"`
	MaxNodes       int    `json:"maxNodes"`
	MinNodes       int    `json:"minNodes"`
	MonthlyBilled  bool   `json:"monthlyBilled"`
	SizeStatus     string `json:"sizeStatus"`
	Status         string `json:"status"`
	UpToDateNodes  int    `json:"upToDateNodes"`
	UpdatedAt      string `json:"updatedAt"`
}
