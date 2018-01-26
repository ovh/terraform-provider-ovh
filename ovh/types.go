package ovh

import (
	"fmt"
	"time"
)

// Opts
type PublicCloudPrivateNetworkCreateOpts struct {
	ProjectId string   `json:"serviceName"`
	VlanId    int      `json:"vlanId"`
	Name      string   `json:"name"`
	Regions   []string `json:"regions"`
}

func (p *PublicCloudPrivateNetworkCreateOpts) String() string {
	return fmt.Sprintf("projectId: %s, vlanId:%d, name: %s, regions: %s", p.ProjectId, p.VlanId, p.Name, p.Regions)
}

// Opts
type PublicCloudPrivateNetworkUpdateOpts struct {
	Name string `json:"name"`
}

type PublicCloudPrivateNetworkRegion struct {
	Status string `json:"status"`
	Region string `json:"region"`
}

func (p *PublicCloudPrivateNetworkRegion) String() string {
	return fmt.Sprintf("Status:%s, Region: %s", p.Status, p.Region)
}

type PublicCloudPrivateNetworkResponse struct {
	Id      string                             `json:"id"`
	Status  string                             `json:"status"`
	Vlanid  int                                `json:"vlanId"`
	Name    string                             `json:"name"`
	Type    string                             `json:"type"`
	Regions []*PublicCloudPrivateNetworkRegion `json:"regions"`
}

func (p *PublicCloudPrivateNetworkResponse) String() string {
	return fmt.Sprintf("Id: %s, Status: %s, Name: %s, Vlanid: %d, Type: %s, Regions: %s", p.Id, p.Status, p.Name, p.Vlanid, p.Type, p.Regions)
}

// Opts
type PublicCloudPrivateNetworksCreateOpts struct {
	ProjectId string `json:"serviceName"`
	NetworkId string `json:"networkId"`
	Dhcp      bool   `json:"dhcp"`
	NoGateway bool   `json:"noGateway"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Network   string `json:"network"`
	Region    string `json:"region"`
}

func (p *PublicCloudPrivateNetworksCreateOpts) String() string {
	return fmt.Sprintf("PCPNSCreateOpts[projectId: %s, networkId:%s, dchp: %v, noGateway: %v, network: %s, start: %s, end: %s, region: %s]",
		p.ProjectId, p.NetworkId, p.Dhcp, p.NoGateway, p.Network, p.Start, p.End, p.Region)
}

type IPPool struct {
	Network string `json:"network"`
	Region  string `json:"region"`
	Dhcp    bool   `json:"dhcp"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

func (p *IPPool) String() string {
	return fmt.Sprintf("IPPool[Network: %s, Region: %s, Dhcp: %v, Start: %s, End: %s]", p.Network, p.Region, p.Dhcp, p.Start, p.End)
}

type PublicCloudPrivateNetworksResponse struct {
	Id        string    `json:"id"`
	GatewayIp string    `json:"gatewayIp"`
	Cidr      string    `json:"cidr"`
	IPPools   []*IPPool `json:"ipPools"`
}

func (p *PublicCloudPrivateNetworksResponse) String() string {
	return fmt.Sprintf("PCPNSResponse[Id: %s, GatewayIp: %s, Cidr: %s, IPPools: %s]", p.Id, p.GatewayIp, p.Cidr, p.IPPools)
}

// Opts
type PublicCloudUserCreateOpts struct {
	ProjectId   string `json:"serviceName"`
	Description string `json:"description"`
}

func (p *PublicCloudUserCreateOpts) String() string {
	return fmt.Sprintf("UserOpts[projectId: %s, description:%s]", p.ProjectId, p.Description)
}

type PublicCloudUserResponse struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Status       string `json:"status"`
	Description  string `json:"description"`
	Password     string `json:"password"`
	CreationDate string `json:"creationDate"`
}

func (p *PublicCloudUserResponse) String() string {
	return fmt.Sprintf("UserResponse[Id: %v, Username: %s, Status: %s, Description: %s, CreationDate: %s]", p.Id, p.Username, p.Status, p.Description, p.CreationDate)
}

type PublicCloudUserOpenstackRC struct {
	Content string `json:"content"`
}

// Opts
type VRackAttachOpts struct {
	Project string `json:"project"`
}

// Opts
type BillingMonthlyOpts struct {
	Project    string `json:"project"`
	InstanceId string `json:"instanceId"`
}

// Task Opts
type TaskOpts struct {
	ServiceName string `json:"serviceName"`
	TaskId      string `json:"taskId"`
}

type PublicCloudMonthlyBilling struct {
	Since  time.Time `json:"since"`
	Status string    `json:"status"`
}

func (p *PublicCloudMonthlyBilling) String() string {
	return fmt.Sprintf("Since:%s, Status: %s", p.Since, p.Status)
}

type BillingMonthlyTaskResponse struct {
	Id             string                     `json:"id"`
	Status         string                     `json:"status"`
	Name           string                     `json:"name"`
	MonthlyBilling *PublicCloudMonthlyBilling `json:"monthlyBilling"`
	Flavor         *PublicCloudFlavor         `json:"flavor"`
}

type PublicCloudInstance struct {
	Id             string                     `json:"id"`
	Status         string                     `json:"status"`   // Instance status
	Name           string                     `json:"name"`     // Instance name
	Region         string                     `json:"region"`   // Instance region
	PlanCode       *string                    `json:"planCode"` // Order plan code
	ImageId        string                     `json:"imageId"`  // Instance image id
	Created        time.Time                  `json:"created"`
	FlavorId       string                     `json:"flavorId"` // Instance flavor id
	MonthlyBilling *PublicCloudMonthlyBilling `json:"monthlyBilling"`
	SSHKeyId       *string                    `json:"sshKeyId"` // Instance ssh key id
	IpAddresses    []*PublicCloudIpAddress    `json:"ipAddresses"`
}

type PublicCloudInstanceDetail struct {
	Id             string                     `json:"id"`
	Status         string                     `json:"status"`   // Instance status
	Name           string                     `json:"name"`     // Instance name
	Region         string                     `json:"region"`   // Instance region
	PlanCode       *string                    `json:"planCode"` // Order plan code
	Image          PublicCloudImage           `json:"image"`    // Instance image id
	Created        time.Time                  `json:"created"`
	SSHKey         string                     `json:"sshKey"` // Instance ssh key id
	MonthlyBilling *PublicCloudMonthlyBilling `json:"monthlyBilling"`
	IpAddresses    []PublicCloudIpAddress     `json:"ipAddresses"`
	Flavor         PublicCloudFlavor          `json:"flavor"`
}

type PublicCloudSSHKeyDetail struct {
	Id          string   `json:"id"`
	FingerPrint string   `json:"fingerPrint"`
	Name        string   `json:"name"`
	Regions     []string `json:"regions"`
	PublicKey   string   `json:"publicKey"`
}

type PublicCloudImage struct {
	Id           string    `json:"id"`
	Visibility   string    `json:"visibility"`
	FlavorType   *string   `json:"flavorType"`
	Status       string    `json:"status"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	PlanCode     *string   `json:"planCode"`
	MinDisk      int64     `json:"minDisk"`
	Size         float64   `json:"size"` // Image size (in GiB)
	Tags         []*string `json:"tags"`
	MinRam       int64     `json:"minRam"`
	CreationDate string    `json:"creationDate"`
	User         string    `json:"user"`
	Type         string    `json:"type"`
}

type PublicCloudIpAddress struct {
	GatewayIp *string `json:"gatewayIp"`
	NetworkId string  `json:"networkId"`
	Ip        string  `json:"ip"`
	Version   int64   `json:"version"` // IP version
	Type      string  `json:"type"`    // Instance IP address type
}

type PublicCloudFlavor struct {
	Id                string                     `json:"id"`
	OutboundBandwidth *int64                     `json:"outboundBandwidth"` // Max capacity of outbound traffic in Mbit/s
	Disk              int64                      `json:"disk"`              // Number of disks
	Name              string                     `json:"name"`              // Flavor name
	Region            string                     `json:"region"`            // Flavor region
	PlanCodes         *PublicCloudFlavorPlanCode `json:"planCodes"`
	OSType            string                     `json:"osType"`
	InboundBandwidth  *int64                     `json:"inboundBandwidth"`
	VCPUs             int32                      `json:"vcpus"` // Number of VCPUs
	Type              string                     `json:"type"`
	Ram               int64                      `json:"ram"`
	Available         bool                       `json:"available"`
}

type PublicCloudFlavorPlanCode struct {
	Hourly  string `json:"hourly"`  // Plan code to order hourly instance
	Monthly string `json:"monthly"` // Plan code to order monthly instance
}

type VRackAttachTaskResponse struct {
	Id           int       `json:"id"`
	Function     string    `json:"function"`
	TargetDomain string    `json:"targetDomain"`
	Status       string    `json:"status"`
	ServiceName  string    `json:"serviceName"`
	OrderId      int       `json:"orderId"`
	LastUpdate   time.Time `json:"lastUpdate"`
	TodoDate     time.Time `json:"TodoDate"`
}

type PublicCloudRegionResponse struct {
	ContinentCode      string                             `json:"continentCode"`
	DatacenterLocation string                             `json:"datacenterLocation"`
	Name               string                             `json:"name"`
	Services           []PublicCloudServiceStatusResponse `json:"services"`
}

func (r *PublicCloudRegionResponse) String() string {
	return fmt.Sprintf("Region: %s, Services: %s", r.Name, r.Services)
}

type PublicCloudServiceStatusResponse struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

func (s *PublicCloudServiceStatusResponse) String() string {
	return fmt.Sprintf("%s: %s", s.Name, s.Status)
}
