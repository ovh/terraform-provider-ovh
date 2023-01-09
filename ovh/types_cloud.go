package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

// Opts
type CloudProjectGatewayInterfaceCreateOpts struct {
	Subnet string `json:"subnetId"`
}

func (p *CloudProjectGatewayInterfaceCreateOpts) String() string {
	return fmt.Sprintf("subnetId: %s", p.Subnet)
}

type CloudProjectGatewayInterfaceResponse struct {
	Id        string `json:"id"`
	Ip        string `json:"ip"`
	NetworkId string `json:"networkId"`
	SubnetId  string `json:"subnetId"`
}

// Opts
type CloudProjectGatewayCreateOpts struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

func (p *CloudProjectGatewayCreateOpts) String() string {
	return fmt.Sprintf("name: %s, model: %s", p.Name, p.Model)
}

type CloudProjectGatewayUpdateOpts struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

type CloudProjectGatewayInterface struct {
	Id        string `json:"id"`
	Ip        string `json:"ip"`
	SubnetId  string `json:"subnetId"`
	NetworkId string `json:"networkId"`
}

type CloudProjectGatewayExternalIp struct {
	Ip       string `json:"ip"`
	SubnetId string `json:"subnetId"`
}

type CloudProjectGatewayExternal struct {
	Ips       []*CloudProjectGatewayExternalIp `json:"ip"`
	NetworkId string                           `json:"networkId"`
}

type CloudProjectGatewayResponse struct {
	Id                  string                          `json:"id"`
	Name                string                          `json:"name"`
	Status              string                          `json:"status"`
	Interfaces          []*CloudProjectGatewayInterface `json:"interfaces"`
	ExternalInformation *CloudProjectGatewayExternal    `json:"externalInformation"`
	Region              string                          `json:"region"`
	Model               string                          `json:"model"`
}

type CloudProjectOperationResponse struct {
	Id          string   `json:"id"`
	Action      string   `json:"action"`
	CreateAt    string   `json:"createdAt"`
	StartedAt   string   `json:"startedAt"`
	CompletedAt *string  `json:"completedAt"`
	Progress    int      `json:"progress"`
	Regions     []string `json:"regions"`
	ResourceId  *string  `json:"resourceId"`
	Status      string   `json:"status"`
}

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
	// Don't set password as it must be set only at creation time
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

type CloudProjectUserS3Credential struct {
	Access      string `json:"access"`
	Secret      string `json:"secret"`
	ServiceName string `json:"tenantId"`
	UserId      string `json:"userId"`
}

func (u *CloudProjectUserS3Credential) String() string {
	return fmt.Sprintf("CloudProjectUserS3Credential[ServiceName:%s, UserId: %s, Access: %s]", u.ServiceName, u.UserId, u.Access)
}

func (u CloudProjectUserS3Credential) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["access_key_id"] = u.Access
	obj["secret_access_key"] = u.Secret
	obj["service_name"] = u.ServiceName
	obj["internal_user_id"] = u.UserId
	return obj
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
