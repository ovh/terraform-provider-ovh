package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

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
	Ips       []*CloudProjectGatewayExternalIp `json:"ips"`
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
	Id            string                     `json:"id"`
	Action        string                     `json:"action"`
	CreateAt      string                     `json:"createdAt"`
	StartedAt     string                     `json:"startedAt"`
	CompletedAt   *string                    `json:"completedAt"`
	Progress      int                        `json:"progress"`
	Regions       []string                   `json:"regions"`
	ResourceId    *string                    `json:"resourceId"`
	Status        string                     `json:"status"`
	SubOperations []CloudProjectSubOperation `json:"subOperations"`
}

type CloudProjectSubOperation struct {
	ResourceId *string `json:"resourceId"`
	Action     string  `json:"action"`
}

func waitForCloudProjectOperation(ctx context.Context, c *ovh.Client, serviceName, operationId, actionType string) (string, error) {
	endpoint := fmt.Sprintf("/cloud/project/%s/operation/%s", url.PathEscape(serviceName), url.PathEscape(operationId))
	resourceID := ""
	err := retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		ro := &CloudProjectOperationResponse{}
		if err := c.GetWithContext(ctx, endpoint, ro); err != nil {
			return retry.NonRetryableError(err)
		}

		switch ro.Status {
		case "in-error":
			return retry.NonRetryableError(fmt.Errorf("operation %q ended in error", ro.Action))
		case "completed":
			if ro.ResourceId != nil {
				resourceID = *ro.ResourceId
			} else if len(ro.SubOperations) > 0 && actionType != "" {
				for _, subOp := range ro.SubOperations {
					if subOp.Action == actionType && subOp.ResourceId != nil {
						resourceID = *subOp.ResourceId
						break
					}
				}
			} else if len(ro.SubOperations) > 0 && ro.SubOperations[0].ResourceId != nil {
				resourceID = *ro.SubOperations[0].ResourceId
			}
			return nil
		default:
			return retry.RetryableError(fmt.Errorf("waiting for operation %s to be completed", ro.Id))
		}
	})

	return resourceID, err
}

// Opts
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

type CloudProjectNetworkPrivateV2CreateOpts struct {
	Name                        string           `json:"name"`
	Cidr                        string           `json:"cidr"`
	IpVersion                   int              `json:"ipVersion"`
	AllocationPools             []AllocationPool `json:"allocationPools"`
	DnsNameServers              []string         `json:"dnsNameServers"`
	HostRoutes                  []HostRoute      `json:"hostRoutes"`
	EnableDHCP                  bool             `json:"enableDhcp"`
	GatewayIp                   string           `json:"gatewayIp,omitempty"`
	EnableGatewayIP             bool             `json:"enableGatewayIp"`
	UseDefaultPublicDNSResolver *bool            `json:"useDefaultPublicDNSResolver"`
}

type HostRoute struct {
	Destination string `json:"destination"`
	Nexthop     string `json:"nextHop"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type CloudProjectNetworkPrivateV2Response struct {
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	Cidr            string           `json:"cidr"`
	IPVersion       int              `json:"ipVersion"`
	DHCPEnabled     bool             `json:"dhcpEnabled"`
	GatewayIp       *string          `json:"gatewayIp" description:"Gateway IP, null means no gateway"`
	AllocationPools []AllocationPool `json:"allocationPools"`
	HostRoutes      []HostRoute      `json:"hostRoutes"`
	DnsNameservers  []string         `json:"dnsNameServers"`
}

func (p *CloudProjectNetworkPrivateV2Response) String() string {
	return fmt.Sprintf("PCPNSResponse[Id: %s, Name: %s, Cidr: %s, IPVersion: %d, DHCPEnabled: %t, GatewayIp: %v, AllocationPools: %v, HostRoutes: %v, DnsNameservers: %s]",
		p.Id, p.Name, p.Cidr, p.IPVersion, p.DHCPEnabled, p.GatewayIp, p.AllocationPools, p.HostRoutes, p.DnsNameservers)
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
	ServiceName string `json:"tenantId"`
	UserId      string `json:"userId"`
}

type CloudProjectUserS3Secret struct {
	Secret string `json:"secret"`
}

type CloudProjectUserS3CredentialSecret struct {
	CloudProjectUserS3Credential
	CloudProjectUserS3Secret
}

func (u *CloudProjectUserS3Credential) String() string {
	return fmt.Sprintf("CloudProjectUserS3Credential[ServiceName:%s, UserId: %s, Access: %s]", u.ServiceName, u.UserId, u.Access)
}

func (u CloudProjectUserS3Credential) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["access_key_id"] = u.Access
	obj["service_name"] = u.ServiceName
	obj["internal_user_id"] = u.UserId
	return obj
}

func (u *CloudProjectUserS3Secret) String() string {
	return "CloudProjectUserS3Secret[Secret: ***]"
}

func (u CloudProjectUserS3Secret) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["secret_access_key"] = u.Secret
	return obj
}

func (u *CloudProjectUserS3CredentialSecret) String() string {
	return fmt.Sprintf("CloudProjectUserS3CredentialSecret[ServiceName:%s, UserId: %s, Access: %s, Secret: ***]", u.ServiceName, u.UserId, u.Access)
}

func (u CloudProjectUserS3CredentialSecret) ToMap() map[string]interface{} {
	obj := u.CloudProjectUserS3Credential.ToMap()
	for k, v := range u.CloudProjectUserS3Secret.ToMap() {
		obj[k] = v
	}
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
