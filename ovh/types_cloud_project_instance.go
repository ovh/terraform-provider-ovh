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

type AutoBackup struct {
	Cron     string `json:"cron"`
	Rotation int    `json:"rotation"`
}

type Flavor struct {
	FlavorId string `json:"id"`
}

type BootFrom struct {
	ImageId  *string `json:"imageId,omitempty"`
	VolumeId *string `json:"volumeId,omitempty"`
}

type Group struct {
	GroupId string `json:"id"`
}

type SshKey struct {
	Name string `json:"name"`
}

type SshKeyCreate struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

type Network struct {
	Public  bool            `json:"public"`
	Private *PrivateNetwork `json:"private,omitempty"`
}

type PrivateNetwork struct {
	FloatingIP       *PrivateNetworkFloatingIP        `json:"floating_ip,omitempty"`
	FloatingIPCreate *PrivateNetworkFloatingIPCreate  `json:"floatingIpCreate,omitempty"`
	Gateway          *PrivateNetworkGateway           `json:"gateway,omitempty"`
	GatewayCreate    *PrivateNetworkGatewayCreate     `json:"gatewayCreate,omitempty"`
	IP               string                           `json:"ip,omitempty"`
	Network          *PrivateNetworkInformation       `json:"network,omitempty"`
	NetworkCreate    *PrivateNetworkInformationCreate `json:"networkCreate,omitempty"`
}

type PrivateNetworkFloatingIP struct {
	ID string `json:"id"`
}

type PrivateNetworkFloatingIPCreate struct {
	Description string `json:"description"`
}

type PrivateNetworkGateway struct {
	ID string `json:"id"`
}

type PrivateNetworkGatewayCreate struct {
	Model string `json:"model"`
	Name  string `json:"name"`
}

type PrivateNetworkInformation struct {
	ID       string `json:"id"`
	SubnetID string `json:"subnetId"`
}

type PrivateNetworkInformationCreate struct {
	Name   string                                `json:"name"`
	Subnet PrivateNetworkInformationSubnetCreate `json:"subnet"`
	VlanId *int                                  `json:"vlanId,omitempty"`
}

type PrivateNetworkInformationSubnetCreate struct {
	CIDR       string `json:"cidr,omitempty"`
	EnableDHCP bool   `json:"enableDhcp"`
	IPVersion  int    `json:"ipVersion,omitempty"`
}

type CloudProjectInstanceCreateOpts struct {
	AutoBackup       *AutoBackup   `json:"autobackup,omitempty"`
	AvailabilityZone *string       `json:"availabilityZone,omitempty"`
	BillingPeriod    string        `json:"billingPeriod"`
	BootFrom         *BootFrom     `json:"bootFrom,omitempty"`
	Bulk             int           `json:"bulk"`
	Flavor           *Flavor       `json:"flavor,omitempty"`
	Group            *Group        `json:"group,omitempty"`
	Name             string        `json:"name"`
	SshKey           *SshKey       `json:"sshKey,omitempty"`
	SshKeyCreate     *SshKeyCreate `json:"sshKeyCreate,omitempty"`
	UserData         *string       `json:"userData,omitempty"`
	Network          *Network      `json:"network,omitempty"`
}

type Address struct {
	Ip      *string `json:"ip"`
	Version *int    `json:"version"`
}

type AttachedVolume struct {
	Id string `json:"id"`
}

type CloudProjectInstanceResponse struct {
	Addresses        []Address        `json:"addresses"`
	AttachedVolumes  []AttachedVolume `json:"attachedVolumes"`
	AvailabilityZone string           `json:"availabilityZone"`
	FlavorId         string           `json:"flavorId"`
	FlavorName       string           `json:"flavorName"`
	Id               string           `json:"id"`
	ImageId          string           `json:"imageId"`
	Name             string           `json:"name"`
	Region           string           `json:"region"`
	SshKey           string           `json:"sshKey"`
	TaskState        string           `json:"taskState"`
	Status           string           `json:"status"`
}

func (v CloudProjectInstanceResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["flavor_id"] = v.FlavorId
	obj["availability_zone"] = v.AvailabilityZone
	obj["flavor_name"] = v.FlavorName
	obj["image_id"] = v.ImageId
	obj["id"] = v.Id
	obj["name"] = v.Name
	obj["ssh_key"] = v.SshKey
	obj["task_state"] = v.TaskState
	obj["status"] = v.Status

	addresses := make([]map[string]interface{}, 0, len(v.Addresses))
	for i := range v.Addresses {
		address := make(map[string]interface{})
		address["ip"] = v.Addresses[i].Ip
		address["version"] = v.Addresses[i].Version
		addresses = append(addresses, address)
	}
	obj["addresses"] = addresses

	attachedVolumes := make([]map[string]interface{}, 0, len(v.AttachedVolumes))
	for i := range v.AttachedVolumes {
		attachedVolume := make(map[string]interface{})
		attachedVolume["id"] = v.AttachedVolumes[i].Id
		attachedVolumes = append(attachedVolumes, attachedVolume)
	}

	obj["attached_volumes"] = attachedVolumes

	return obj
}

func (a Address) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["ip"] = a.Ip
	obj["version"] = a.Version
	return obj
}

func (a AttachedVolume) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["attached_volumes"] = a.Id
	return obj
}

type CloudProjectInstanceResponseList struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetFlavorId(i interface{}) *Flavor {
	if i == nil {
		return nil
	}
	flavorId := Flavor{}
	flavorSet := i.(*schema.Set).List()
	for _, flavor := range flavorSet {
		mapping := flavor.(map[string]interface{})
		flavorId.FlavorId = mapping["flavor_id"].(string)
	}
	return &flavorId
}

func GetAutoBackup(i interface{}) *AutoBackup {
	if i == nil {
		return nil
	}
	autoBackupOut := AutoBackup{}

	autoBackupSet := i.(*schema.Set).List()
	if len(autoBackupSet) == 0 {
		return nil
	}
	for _, autoBackup := range autoBackupSet {
		mapping := autoBackup.(map[string]interface{})
		autoBackupOut.Cron = mapping["cron"].(string)
		autoBackupOut.Rotation = mapping["rotation"].(int)
	}
	return &autoBackupOut
}

func GetBootFrom(i interface{}) *BootFrom {
	if i == nil {
		return nil
	}
	bootFromOutput := BootFrom{}

	bootFromSet := i.(*schema.Set).List()
	for _, bootFrom := range bootFromSet {
		mapping := bootFrom.(map[string]interface{})
		bootFromOutput.ImageId = helpers.GetNilStringPointerFromData(mapping, "image_id")
		bootFromOutput.VolumeId = helpers.GetNilStringPointerFromData(mapping, "volume_id")
	}

	return &bootFromOutput
}

func GetGroup(i interface{}) *Group {
	if i == nil {
		return nil
	}
	groupOut := Group{}

	groupSet := i.(*schema.Set).List()
	if len(groupSet) == 0 {
		return nil
	}
	for _, group := range groupSet {
		mapping := group.(map[string]interface{})
		groupOut.GroupId = mapping["id"].(string)
	}
	return &groupOut
}

func GetSshKey(i interface{}) *SshKey {
	if i == nil {
		return nil
	}
	sshOutput := SshKey{}

	sshSet := i.(*schema.Set).List()
	if len(sshSet) == 0 {
		return nil
	}
	for _, ssh := range sshSet {
		mapping := ssh.(map[string]interface{})
		sshOutput.Name = mapping["name"].(string)
	}

	return &sshOutput
}

func GetSshKeyCreate(i interface{}) *SshKeyCreate {
	if i == nil {
		return nil
	}
	sshCreateOutput := SshKeyCreate{}

	sshCreateSet := i.(*schema.Set).List()
	if len(sshCreateSet) == 0 {
		return nil
	}
	for _, ssh := range sshCreateSet {
		mapping := ssh.(map[string]interface{})
		sshCreateOutput.Name = mapping["name"].(string)
		sshCreateOutput.PublicKey = mapping["public_key"].(string)
	}

	return &sshCreateOutput
}

func GetNetwork(i interface{}) *Network {
	if i == nil {
		return nil
	}
	networkOutput := Network{}

	networkSet := i.(*schema.Set).List()
	for _, network := range networkSet {
		mapping := network.(map[string]interface{})
		networkOutput.Public = mapping["public"].(bool)

		privateNet, ok := mapping["private"]
		if !ok || privateNet == nil || privateNet.(*schema.Set).Len() == 0 {
			continue
		}

		for _, priv := range privateNet.(*schema.Set).List() {
			networkOutput.Private = &PrivateNetwork{}
			private := priv.(map[string]interface{})

			if floatingIP, ok := private["floating_ip"]; ok {
				for _, float := range floatingIP.(*schema.Set).List() {
					params := float.(map[string]interface{})
					networkOutput.Private.FloatingIP = &PrivateNetworkFloatingIP{
						ID: params["id"].(string),
					}
				}
			}

			if floatingIPCreate, ok := private["floating_ip_create"]; ok {
				for _, float := range floatingIPCreate.(*schema.Set).List() {
					params := float.(map[string]interface{})
					networkOutput.Private.FloatingIPCreate = &PrivateNetworkFloatingIPCreate{
						Description: params["description"].(string),
					}
				}
			}

			if gateway, ok := private["gateway"]; ok {
				for _, gateway := range gateway.(*schema.Set).List() {
					params := gateway.(map[string]interface{})
					networkOutput.Private.Gateway = &PrivateNetworkGateway{
						ID: params["id"].(string),
					}
				}
			}

			if gatewayCreate, ok := private["gateway_create"]; ok {
				for _, gateway := range gatewayCreate.(*schema.Set).List() {
					params := gateway.(map[string]interface{})
					networkOutput.Private.GatewayCreate = &PrivateNetworkGatewayCreate{
						Model: params["model"].(string),
						Name:  params["name"].(string),
					}
				}
			}

			networkOutput.Private.IP = private["ip"].(string)

			if network, ok := private["network"]; ok {
				for _, net := range network.(*schema.Set).List() {
					params := net.(map[string]interface{})
					networkOutput.Private.Network = &PrivateNetworkInformation{
						ID:       params["id"].(string),
						SubnetID: params["subnet_id"].(string),
					}
				}
			}

			if networkCreate, ok := private["network_create"]; ok {
				for _, net := range networkCreate.(*schema.Set).List() {
					params := net.(map[string]interface{})
					networkOutput.Private.NetworkCreate = &PrivateNetworkInformationCreate{
						Name: params["name"].(string),
					}

					if vlanID, ok := params["vlan_id"]; ok {
						intVlanID := vlanID.(int)
						networkOutput.Private.NetworkCreate.VlanId = &intVlanID
					}

					for _, subnet := range params["subnet"].(*schema.Set).List() {
						subnetParams := subnet.(map[string]interface{})
						networkOutput.Private.NetworkCreate.Subnet = PrivateNetworkInformationSubnetCreate{
							CIDR:       subnetParams["cidr"].(string),
							EnableDHCP: subnetParams["enable_dhcp"].(bool),
							IPVersion:  subnetParams["ip_version"].(int),
						}
					}
				}
			}
		}

	}

	return &networkOutput
}

func (cpir *CloudProjectInstanceCreateOpts) FromResource(d *schema.ResourceData) {
	cpir.Flavor = GetFlavorId(d.Get("flavor"))
	cpir.AutoBackup = GetAutoBackup(d.Get("auto_backup"))
	cpir.BootFrom = GetBootFrom(d.Get("boot_from"))
	cpir.Group = GetGroup(d.Get("group"))
	cpir.SshKey = GetSshKey(d.Get("ssh_key"))
	cpir.SshKeyCreate = GetSshKeyCreate(d.Get("ssh_key_create"))
	cpir.Network = GetNetwork(d.Get("network"))
	cpir.BillingPeriod = d.Get("billing_period").(string)
	cpir.Name = d.Get("name").(string)
	cpir.AvailabilityZone = helpers.GetNilStringPointerFromData(d, "availability_zone")
	cpir.UserData = helpers.GetNilStringPointerFromData(d, "user_data")
}

func waitForCloudProjectInstance(ctx context.Context, c *ovh.Client, serviceName, region, instance string, timeout time.Duration) error {
	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/instance/%s", url.PathEscape(serviceName), url.PathEscape(region), url.PathEscape(instance))

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {

		ro := &CloudProjectInstanceResponse{}
		if err := c.GetWithContext(ctx, endpoint, ro); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error calling Get %s:\n\t %q", endpoint, err))
		}

		switch ro.Status {
		case "ERROR", "UNKNOWN":
			return retry.NonRetryableError(fmt.Errorf("invalid status for instance: %q", ro.Status))
		case "ACTIVE":
			return nil
		default:
			return retry.RetryableError(fmt.Errorf("waiting for instance %s to be in state ACTIVE", ro.Id))
		}
	})

	return err
}
