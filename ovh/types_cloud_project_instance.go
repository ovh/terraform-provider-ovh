package ovh

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
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
	Public bool `json:"public"`
}

type CloudProjectInstanceCreateOpts struct {
	AutoBackup    *AutoBackup   `json:"autobackup,omitempty"`
	BillingPeriod string        `json:"billingPeriod"`
	BootFrom      *BootFrom     `json:"bootFrom,omitempty"`
	Bulk          int           `json:"bulk"`
	Flavor        *Flavor       `json:"flavor,omitempty"`
	Group         *Group        `json:"group,omitempty"`
	Name          string        `json:"name"`
	SshKey        *SshKey       `json:"sshKey,omitempty"`
	SshKeyCreate  *SshKeyCreate `json:"sshKeyCreate,omitempty"`
	UserData      *string       `json:"userData,omitempty"`
	Network       *Network      `json:"network,omitempty"`
}

type Address struct {
	Ip      *string `json:"ip"`
	Version *int    `json:"version"`
}

type AttachedVolume struct {
	Id string `json:"id"`
}

type CloudProjectInstanceResponse struct {
	Addresses       []Address        `json:"addresses"`
	AttachedVolumes []AttachedVolume `json:"attachedVolumes"`
	FlavorId        string           `json:"flavorId"`
	FlavorName      string           `json:"flavorName"`
	Id              string           `json:"id"`
	ImageId         string           `json:"imageId"`
	Name            string           `json:"name"`
	Region          string           `json:"region"`
	SshKey          string           `json:"sshKey"`
	TaskState       string           `json:"taskState"`
}

func (v CloudProjectInstanceResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["flavor_id"] = v.FlavorId
	obj["flavor_name"] = v.FlavorName
	obj["image_id"] = v.ImageId
	obj["id"] = v.Id
	obj["name"] = v.Name
	obj["ssh_key"] = v.SshKey
	obj["task_state"] = v.TaskState

	addresses := make([]map[string]interface{}, 0)
	for i := range v.Addresses {
		address := make(map[string]interface{})
		address["ip"] = v.Addresses[i].Ip
		address["version"] = v.Addresses[i].Version
		addresses = append(addresses, address)
	}
	obj["addresses"] = addresses

	attachedVolumes := make([]map[string]interface{}, 0)
	for i := range v.AttachedVolumes {
		attachedVolume := make(map[string]interface{})
		attachedVolume["id"] = v.AttachedVolumes[i].Id
		attachedVolumes = append(attachedVolumes, attachedVolume)
	}

	obj["attached_volumes"] = attachedVolumes

	return obj
}

type CloudProjectOperation struct {
	Id            string         `json:"id"`
	Status        string         `json:"status"`
	SubOperations []SubOperation `json:"subOperations"`
}

type SubOperation struct {
	Id         string `json:"id"`
	ResourceId string `json:"resourceId"`
	Status     string `json:"status"`
}

func (sb SubOperation) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["id"] = sb.Id
	obj["resourceId"] = sb.ResourceId
	obj["status"] = sb.Status

	log.Printf("[DEBUG] tata suboperation %+v:", obj)
	return obj
}

func (o CloudProjectOperation) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["id"] = o.Id
	obj["status"] = o.Status
	subOperations := make([]map[string]interface{}, len(o.SubOperations))
	for _, subOperation := range o.SubOperations {
		subOperations = append(subOperations, subOperation.ToMap())
	}
	obj["subOperations"] = subOperations
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

func GetFlaorId(i interface{}) *Flavor {
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
	log.Printf("[DEBUG] BootFrom ------- %v", i)
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
	log.Printf("[DEBUG] Group ------- %v", i)
	if i == nil {
		return nil
	}
	groupOut := Group{}

	groupSet := i.(*schema.Set).List()
	for _, group := range groupSet {
		mapping := group.(map[string]interface{})
		groupOut.GroupId = mapping["id"].(string)
	}
	return &groupOut
}

func GetSshKey(i interface{}) *SshKey {
	log.Printf("[DEBUG] SshKey ------- %v", i)
	if i == nil {
		return nil
	}
	sshOutput := SshKey{}

	sshSet := i.(*schema.Set).List()
	for _, ssh := range sshSet {
		mapping := ssh.(map[string]interface{})
		sshOutput.Name = mapping["name"].(string)
	}

	return &sshOutput
}

func GetSshKeyCreate(i interface{}) *SshKeyCreate {
	log.Printf("[DEBUG] SshKeyCreate ------- %v", i)
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
		sshCreateOutput.Name = mapping["public_key"].(string)
	}

	return &sshCreateOutput
}

func GetNetwork(i interface{}) *Network {
	log.Printf("[DEBUG] Network ------- %v", i)
	if i == nil {
		return nil
	}
	networkOutput := Network{}

	networkSet := i.(*schema.Set).List()
	for _, network := range networkSet {
		mapping := network.(map[string]interface{})
		networkOutput.Public = mapping["public"].(bool)
	}
	return &networkOutput
}

func (cpir *CloudProjectInstanceCreateOpts) FromResource(d *schema.ResourceData) {
	cpir.Flavor = GetFlaorId(d.Get("flavor"))
	log.Printf("[DEBUG] flavor ------- %v", cpir.Flavor)
	cpir.AutoBackup = GetAutoBackup(d.Get("auto_backup"))
	log.Printf("[DEBUG] auto -------  %v", cpir.AutoBackup)
	cpir.BootFrom = GetBootFrom(d.Get("boot_from"))
	log.Printf("[DEBUG] boot_from -------  %v", cpir.BootFrom)
	cpir.Group = GetGroup(d.Get("group"))
	log.Printf("[DEBUG] group -------  %v", cpir.Group)
	cpir.SshKey = GetSshKey(d.Get("ssh_key"))
	log.Printf("[DEBUG] ssh_key -------  %v", cpir.SshKey)
	cpir.SshKeyCreate = GetSshKeyCreate(d.Get("ssh_key_create"))
	log.Printf("[DEBUG] ssh_key_create -------  %v", cpir.SshKeyCreate)
	cpir.Network = GetNetwork(d.Get("network"))
	log.Printf("[DEBUG] network -------  %v", cpir.Network)
	cpir.BillingPeriod = d.Get("billing_period").(string)
	log.Printf("[DEBUG] billing_period -------  %v", cpir.BillingPeriod)
	cpir.Name = d.Get("name").(string)
	log.Printf("[DEBUG] name -------  %v", cpir.Name)
	cpir.UserData = helpers.GetNilStringPointerFromData(d, "user_data")
	log.Printf("[DEBUG] user_data -------  %v", cpir.UserData)
}
