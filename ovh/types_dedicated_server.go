package ovh

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type DedicatedServer struct {
	AvailabilityZone   string `json:"availabilityZone"`
	Name               string `json:"name"`
	BootId             int    `json:"bootId"`
	BootScript         string `json:"bootScript"`
	CommercialRange    string `json:"commercialRange"`
	Datacenter         string `json:"datacenter"`
	EfiBootloaderPath  string `json:"efiBootloaderPath"`
	Ip                 string `json:"ip"`
	LinkSpeed          int    `json:"linkSpeed"`
	Monitoring         bool   `json:"monitoring"`
	NewUpgradeSystem   bool   `json:"newUpgradeSystem"`
	NoIntervention     bool   `json:"noIntervention"`
	Os                 string `json:"os"`
	PowerState         string `json:"powerState"`
	ProfessionalUse    bool   `json:"professionalUse"`
	Rack               string `json:"rack"`
	Region             string `json:"region"`
	RescueMail         string `json:"rescueMail"`
	RescueSshKey       string `json:"rescueSshKey"`
	Reverse            string `json:"reverse"`
	RootDevice         string `json:"rootDevice"`
	ServerId           int    `json:"serverId"`
	State              string `json:"state"`
	SupportLevel       string `json:"supportLevel"`
	IamResourceDetails `json:"iam"`
}

func (ds DedicatedServer) String() string {
	return fmt.Sprintf(
		"name: %v, ip: %v, dc: %v, state: %v",
		ds.Name,
		ds.Ip,
		ds.Datacenter,
		ds.State,
	)
}

type DedicatedServerUpdateOpts struct {
	BootId            *int64  `json:"bootId,omitempty"`
	BootScript        *string `json:"bootScript,omitempty"`
	EfiBootloaderPath *string `json:"efiBootloaderPath,omitempty"`
	Monitoring        *bool   `json:"monitoring,omitempty"`
	State             *string `json:"state,omitempty"`
}

func (opts *DedicatedServerUpdateOpts) FromResource(d *schema.ResourceData) *DedicatedServerUpdateOpts {
	opts.BootId = helpers.GetNilInt64PointerFromData(d, "boot_id")
	opts.BootScript = helpers.GetNilStringPointerFromData(d, "boot_script")
	opts.EfiBootloaderPath = helpers.GetNilStringPointerFromData(d, "efi_bootloader_path")
	opts.Monitoring = helpers.GetNilBoolPointerFromData(d, "monitoring")
	opts.State = helpers.GetNilStringPointerFromData(d, "state")
	return opts
}

type DedicatedServerVNI struct {
	Enabled    bool     `json:"enabled"`
	Mode       string   `json:"mode"`
	Name       string   `json:"name"`
	NICs       []string `json:"networkInterfaceController"`
	ServerName string   `json:"serverName"`
	Uuid       string   `json:"uuid"`
	Vrack      *string  `json:"vrack,omitempty"`
}

func (vni DedicatedServerVNI) String() string {
	return fmt.Sprintf(
		"name: %v, uuid: %v, mode: %v, enabled: %v",
		vni.Name,
		vni.Uuid,
		vni.Mode,
		vni.Enabled,
	)
}

func (v DedicatedServerVNI) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["enabled"] = v.Enabled
	obj["mode"] = v.Mode
	obj["name"] = v.Name
	obj["nics"] = v.NICs
	obj["server_name"] = v.ServerName
	obj["uuid"] = v.Uuid

	if v.Vrack != nil {
		obj["vrack"] = *v.Vrack
	}

	return obj
}

type DedicatedServerBoot struct {
	BootId      int    `json:"bootId"`
	BootType    string `json:"bootType"`
	Description string `json:"description"`
	Kernel      string `json:"kernel"`
}

type DedicatedServerTask struct {
	Id         int64     `json:"taskId"`
	Function   string    `json:"function"`
	Comment    string    `json:"comment"`
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"lastUpdate"`
	DoneDate   time.Time `json:"doneDate"`
	StartDate  time.Time `json:"startDate"`
}

type DedicatedServerReinstallTaskCreateOpts struct {
	Os             string                                      `json:"operatingSystem"`
	Customizations *DedicatedServerReinstallTaskCustomizations `json:"customizations,omitempty"`
	Properties     map[string]interface{}                      `json:"properties,omitempty"`
	Storage        []DedicatedServerReinstallTaskStorage       `json:"storage,omitempty"`
}

func (opts *DedicatedServerReinstallTaskCreateOpts) FromResource(d *schema.ResourceData) *DedicatedServerReinstallTaskCreateOpts {
	opts.Os = d.Get("os").(string)

	Customizations := d.Get("customizations").([]interface{})
	if len(Customizations) == 1 {
		opts.Customizations = (&DedicatedServerReinstallTaskCustomizations{}).FromResource(d, "customizations.0")
	}

	Properties := d.Get("properties").(map[string]interface{})
	if len(Properties) >= 1 {
		opts.Properties = d.Get("properties").(map[string]interface{})
	}

	Storage := d.Get("storage").([]interface{})
	if len(Storage) >= 1 {
		opts.Storage = make([]DedicatedServerReinstallTaskStorage, len(Storage))
		for i := range Storage {
			opts.Storage[i] = *(new(DedicatedServerReinstallTaskStorage).FromResource(d, fmt.Sprintf("storage.%d", i)))
		}
	}

	return opts
}

type DedicatedServerReinstallTaskStorage struct {
	DiskGroupId  int                   `json:"diskGroupId,omitempty"`
	HardwareRaid []HardwareRaidInstall `json:"hardwareRaid,omitempty"`
	Partitioning *Partitioning         `json:"partitioning,omitempty"`
}

type HardwareRaidInstall struct {
	Arrays    *int `json:"arrays,omitempty"`
	Disks     *int `json:"disks,omitempty"`
	RaidLevel *int `json:"raidLevel,omitempty"`
	Spares    *int `json:"spares,omitempty"`
}

type Partitioning struct {
	Disks      int      `json:"disks,omitempty"`
	Layout     []Layout `json:"layout,omitempty"`
	SchemeName string   `json:"schemeName,omitempty"`
}

type Layout struct {
	FileSystem string `json:"fileSystem,omitempty"`
	MountPoint string `json:"mountPoint,omitempty"`
	RaidLevel  int    `json:"raidLevel,omitempty"`
	Size       int    `json:"size,omitempty"`
	Extras     Extras `json:"extras,omitempty"`
}

type Extras struct {
	Lv ExtrasDetails `json:"lv,omitempty"`
	Zp ExtrasDetails `json:"zp,omitempty"`
}

type ExtrasDetails struct {
	Name string `json:"name,omitempty"`
}

type DedicatedServerReinstallTaskCustomizations struct {
	ConfigDriveUserData             *string                `json:"configDriveUserData,omitempty"`
	EfiBootloaderPath               *string                `json:"efiBootloaderPath,omitempty"`
	Hostname                        *string                `json:"hostname,omitempty"`
	HttpHeaders                     map[string]interface{} `json:"httpHeaders,omitempty"`
	ImageCheckSum                   *string                `json:"imageCheckSum,omitempty"`
	ImageCheckSumType               *string                `json:"imageCheckSumType,omitempty"`
	ImageType                       *string                `json:"imageType,omitempty"`
	ImageURL                        *string                `json:"imageURL,omitempty"`
	Language                        *string                `json:"language,omitempty"`
	PostInstallationScript          *string                `json:"postInstallationScript,omitempty"`
	PostInstallationScriptExtension *string                `json:"postInstallationScriptExtension,omitempty"`
	SshKey                          *string                `json:"sshKey,omitempty"`
}

func (opts *DedicatedServerReinstallTaskCustomizations) FromResource(d *schema.ResourceData, parent string) *DedicatedServerReinstallTaskCustomizations {
	opts.ConfigDriveUserData = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.config_drive_user_data", parent))
	opts.EfiBootloaderPath = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.efi_bootloader_path", parent))
	opts.Hostname = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.hostname", parent))
	opts.ImageCheckSum = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.image_check_sum", parent))
	opts.ImageCheckSumType = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.image_check_sum_type", parent))
	opts.ImageType = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.image_type", parent))
	opts.ImageURL = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.image_url", parent))
	opts.Language = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.language", parent))
	opts.PostInstallationScript = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script", parent))
	opts.PostInstallationScriptExtension = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script_extension", parent))
	opts.SshKey = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.ssh_key", parent))
	opts.HttpHeaders = helpers.GetMapFromData(d, fmt.Sprintf("%s.http_headers", parent))

	return opts
}

func (opts *DedicatedServerReinstallTaskStorage) FromResource(d *schema.ResourceData, parent string) *DedicatedServerReinstallTaskStorage {
	opts.DiskGroupId = d.Get(fmt.Sprintf("%s.disk_group_id", parent)).(int)

	opts.Partitioning = (&Partitioning{}).FromResource(d, fmt.Sprintf("%s.partitioning.0", parent))

	hardwareRaid := d.Get(fmt.Sprintf("%s.hardware_raid", parent)).([]interface{})
	if len(hardwareRaid) >= 1 {
		for i := 0; i < len(hardwareRaid); i++ {
			userHardwareRaid := (&HardwareRaidInstall{}).FromResource(d, fmt.Sprintf("%s.hardware_raid.%d", parent, i))
			opts.HardwareRaid = append(opts.HardwareRaid, userHardwareRaid)
		}
	}

	return opts
}

func (opts HardwareRaidInstall) FromResource(d *schema.ResourceData, parent string) HardwareRaidInstall {
	opts.Arrays = helpers.GetNilIntPointerFromData(d, fmt.Sprintf("%s.arrays", parent))
	opts.Disks = helpers.GetNilIntPointerFromData(d, fmt.Sprintf("%s.disks", parent))
	opts.RaidLevel = helpers.GetNilIntPointerFromData(d, fmt.Sprintf("%s.raid_level", parent))
	opts.Spares = helpers.GetNilIntPointerFromData(d, fmt.Sprintf("%s.spares", parent))
	return opts
}

func (opts *Partitioning) FromResource(d *schema.ResourceData, parent string) *Partitioning {
	opts.Disks = d.Get(fmt.Sprintf("%s.disks", parent)).(int)

	layout := d.Get(fmt.Sprintf("%s.layout", parent)).([]interface{})
	if len(layout) >= 1 {
		for i := 0; i < len(layout); i++ {
			userLayout := (&Layout{}).FromResource(d, fmt.Sprintf("%s.layout.%d", parent, i))
			opts.Layout = append(opts.Layout, userLayout)
		}
	}

	opts.SchemeName = d.Get(fmt.Sprintf("%s.scheme_name", parent)).(string)

	return opts
}

func (opts Layout) FromResource(d *schema.ResourceData, parent string) Layout {
	opts.FileSystem = d.Get(fmt.Sprintf("%s.file_system", parent)).(string)
	opts.MountPoint = d.Get(fmt.Sprintf("%s.mount_point", parent)).(string)
	opts.RaidLevel = d.Get(fmt.Sprintf("%s.raid_level", parent)).(int)
	opts.Size = d.Get(fmt.Sprintf("%s.size", parent)).(int)

	extras := d.Get(fmt.Sprintf("%s.extras", parent)).([]interface{})
	if len(extras) >= 1 {
		opts.Extras = (&Extras{}).FromResource(d, fmt.Sprintf("%s.extras.0", parent))
	}

	return opts
}

func (opts Extras) FromResource(d *schema.ResourceData, parent string) Extras {
	Lv := d.Get(fmt.Sprintf("%s.lv", parent)).([]interface{})
	if len(Lv) >= 1 {
		opts.Lv = (&ExtrasDetails{}).FromResource(d, fmt.Sprintf("%s.lv.0", parent))
	}
	Zp := d.Get(fmt.Sprintf("%s.zp", parent)).([]interface{})
	if len(Zp) >= 1 {
		opts.Zp = (&ExtrasDetails{}).FromResource(d, fmt.Sprintf("%s.zp.0", parent))
	}

	return opts
}

func (opts ExtrasDetails) FromResource(d *schema.ResourceData, parent string) ExtrasDetails {
	opts.Name = d.Get(fmt.Sprintf("%s.name", parent)).(string)

	return opts
}
