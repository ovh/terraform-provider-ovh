package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type AutoBackup struct {
	Cron     *string `json:"cron",omitempty`
	Rotation *int    `json:"rotation",omitempty`
}

func (ab *AutoBackup) FromResource(d *schema.ResourceData, parent string) *AutoBackup {
	ab.Cron = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.cron", parent))
	ab.Rotation = helpers.GetNilIntPointerFromData(d, fmt.Sprintf("%s.rotation", parent))
	return ab
}

type Flavor struct {
	FlavorId *string `json:"id",omitempty`
}

func (fl *Flavor) FromResource(d *schema.ResourceData, parent string) *Flavor {
	fl.FlavorId = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.flavor_id", parent))
	return fl
}

type BootFrom struct {
	ImageId  *string `json:"imageId",omitempty`
	VolumeId *string `json:"volumeId",omitempty`
}

func (bf *BootFrom) FromResource(d *schema.ResourceData, parent string) *BootFrom {
	bf.ImageId = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.image_id", parent))
	bf.VolumeId = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.volume_id", parent))
	return bf
}

type Group struct {
	GroupId *string `json:"id";omitempty`
}

func (g *Group) FromResource(d *schema.ResourceData, parent string) *Group {
	g.GroupId = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.id", parent))
	return g
}

type SshKey struct {
	Name *string `json:"name",omitempty`
}

func (sk *SshKey) FromResource(d *schema.ResourceData, parent string) *SshKey {
	sk.Name = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.name", parent))
	return sk
}

type SshKeyCreate struct {
	Name      *string `json:"name",omitempty`
	PublicKey *string `json:"publicKey",omitempty`
}

func (skc *SshKeyCreate) FromResource(d *schema.ResourceData, parent string) *SshKeyCreate {
	skc.Name = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.name", parent))
	skc.PublicKey = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.public_key", parent))
	return skc
}

type Network struct {
	Public *bool `json:"public",omitempty`
}

func (n *Network) FromResource(d *schema.ResourceData, parent string) *Network {
	n.Public = helpers.GetNilBoolPointerFromData(d, fmt.Sprintf("%s.public", parent))
	return n
}

type CloudProjectInstanceCreateOpts struct {
	AutoBackup    *AutoBackup   `json:"autobackup"`
	BillingPeriod string        `json:"billingPeriod"`
	BootFrom      *BootFrom     `json:"bootFrom"`
	Bulk          int           `json:"bulk"`
	Flavor        *Flavor       `json:"flavor"`
	Group         *Group        `json:"group"`
	Name          string        `json:"name"`
	SshKey        *SshKey       `json:"sshKey"`
	SshKeyCreate  *SshKeyCreate `json:"sshKeyCreate"`
	UserData      string        `json:"userData"`
	Network       *Network      `json:"network"`
}

type Address struct {
	Ip      *string `json:"ip",omitempty`
	Version *int    `json:"version",omitempty`
}

type AttachedVolume struct {
	Id string `json:"id",omitempty`
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
	log.Printf("[DEBUG] titi operation %+v:", o)
	log.Printf("[DEBUG] titi operation %+v:", o.SubOperations)
	for _, subOperation := range o.SubOperations {
		log.Printf("[DEBUG] tutu operation %+v:", subOperation)
		subOperations = append(subOperations, subOperation.ToMap())
	}
	obj["subOperations"] = subOperations
	log.Printf("[DEBUG] titi operation %+v:", obj)
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

func (cpir CloudProjectInstanceResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["flavor_id"] = cpir.FlavorId
	obj["flavor_name"] = cpir.FlavorName
	obj["id"] = cpir.Id
	obj["image_id"] = cpir.ImageId
	obj["name"] = cpir.Name
	obj["ssh_key"] = cpir.SshKey
	obj["task_state"] = cpir.TaskState

	addresses := make([]map[string]interface{}, len(cpir.Addresses))
	for _, address := range cpir.Addresses {
		addresses = append(addresses, address.ToMap())
	}
	obj["addresses"] = addresses

	attachedVolumes := make([]map[string]interface{}, len(cpir.AttachedVolumes))
	for _, attachedVolume := range cpir.AttachedVolumes {
		attachedVolumes = append(addresses, attachedVolume.ToMap())
	}

	obj["attached_volumes"] = attachedVolumes
	return obj
}

type CloudProjectInstanceResponseList struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (cpir *CloudProjectInstanceCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectInstanceCreateOpts {

	bootFrom := d.Get("boot_from").([]interface{})
	flavor := d.Get("flavor").([]interface{})
	autoBackup := d.Get("auto_backup").([]interface{})
	group := d.Get("group").([]interface{})
	sshKey := d.Get("ssh_key").([]interface{})
	sshKeyCreate := d.Get("ssh_key_create").([]interface{})
	network := d.Get("network").([]interface{})

	cpir.BillingPeriod = d.Get("billing_period").(string)
	cpir.Name = d.Get("name").(string)
	cpir.Bulk = 1

	if len(bootFrom) == 1 {
		cpir.BootFrom = (&BootFrom{}).FromResource(d, "boot_from.0")
	}

	if len(flavor) == 1 {
		cpir.Flavor = (&Flavor{}).FromResource(d, "flavor.0")
	}

	if len(autoBackup) == 1 {
		cpir.AutoBackup = (&AutoBackup{}).FromResource(d, "auto_backup.0")
	}

	if len(group) == 1 {
		cpir.Group = (&Group{}).FromResource(d, "group.0")
	}

	if len(sshKey) == 1 {
		cpir.SshKey = (&SshKey{}).FromResource(d, "ssh_key.0")
	}

	if len(sshKeyCreate) == 1 {
		cpir.SshKeyCreate = (&SshKeyCreate{}).FromResource(d, "ssh_key_create.0")
	}

	if len(network) == 1 {
		cpir.Network = (&Network{}).FromResource(d, "network.0")
	}

	return cpir
}
