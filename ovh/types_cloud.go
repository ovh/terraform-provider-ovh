package ovh

import (
	"fmt"
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
	ProjectId   string `json:"serviceName"`
	Description string `json:"description"`
}

func (p *CloudUserCreateOpts) String() string {
	return fmt.Sprintf("UserOpts[projectId: %s, description:%s]", p.ProjectId, p.Description)
}

type CloudUserResponse struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Status       string `json:"status"`
	Description  string `json:"description"`
	Password     string `json:"password"`
	CreationDate string `json:"creationDate"`
}

func (p *CloudUserResponse) String() string {
	return fmt.Sprintf("UserResponse[Id: %v, Username: %s, Status: %s, Description: %s, CreationDate: %s]", p.Id, p.Username, p.Status, p.Description, p.CreationDate)
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
