package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type InstallationTemplate struct {
	BitFormat             int                                `json:"bitFormat,omitempty"`
	Category              string                             `json:"category,omitempty"`
	Customization         *InstallationTemplateCustomization `json:"customization,omitempty"`
	DefaultLanguage       string                             `json:"defaultLanguage,omitempty"`
	Description           string                             `json:"description,omitempty"`
	Distribution          string                             `json:"distribution,omitempty"`
	EndOfInstall          string                             `json:"endOfInstall,omitempty"`
	Family                string                             `json:"family,omitempty"`
	Filesystems           []string                           `json:"filesystems"`
	HardRaidConfiguration bool                               `json:"hardRaidConfiguration,omitempty"`
	Inputs                []InstallationTemplateInputs       `json:"inputs,omitempty"`
	License               *InstallationTemplateLicense       `json:"license,omitempty"`
	LvmReady              *bool                              `json:"lvmReady,omitempty"`
	NoPartitioning        bool                               `json:"noPartitioning,omitempty"`
	Project               *InstallationTemplateProject       `json:"project,omitempty"`
	SoftRaidOnlyMirroring bool                               `json:"soft_raid_only_mirroring,omitempty"`
	Subfamily             string                             `json:"subfamily,omitempty"`
	TemplateName          string                             `json:"templateName"`
}

func (v InstallationTemplate) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["bit_format"] = v.BitFormat
	obj["category"] = v.Category

	if v.Customization != nil {
		customization := v.Customization.ToMap()
		if customization != nil {
			obj["customization"] = []interface{}{customization}
		}
	}

	obj["description"] = v.Description
	obj["distribution"] = v.Distribution
	obj["end_of_install"] = v.EndOfInstall
	obj["family"] = v.Family
	obj["filesystems"] = v.Filesystems

	obj["hard_raid_configuration"] = v.HardRaidConfiguration

	if v.Inputs != nil {
		inputs := make([]interface{}, len(v.Inputs))
		for i, input := range v.Inputs {
			inputs[i] = input.ToMap()
		}
		obj["inputs"] = inputs
	}

	if v.License != nil {
		obj["license"] = []map[string]interface{}{v.License.ToMap()}
	}

	if v.LvmReady != nil {
		obj["lvm_ready"] = *v.LvmReady
	}

	if v.Project != nil {
		obj["project"] = []map[string]interface{}{v.Project.ToMap()}
	}

	obj["no_partitioning"] = v.NoPartitioning
	obj["soft_raid_only_mirroring"] = v.SoftRaidOnlyMirroring
	obj["subfamily"] = v.Subfamily
	obj["template_name"] = v.TemplateName

	return obj
}

type InstallationTemplateCreateOpts struct {
	BaseTemplateName string `json:"baseTemplateName"`
	Name             string `json:"name"`
	DefaultLanguage  string `json:"defaultLanguage,omitempty"`
}

func (opts *InstallationTemplateCreateOpts) FromResource(d *schema.ResourceData) *InstallationTemplateCreateOpts {
	opts.BaseTemplateName = d.Get("base_template_name").(string)
	opts.Name = d.Get("template_name").(string)
	return opts
}

type InstallationTemplateUpdateOpts struct {
	DefaultLanguage string                             `json:"defaultLanguage,omitempty"`
	Customization   *InstallationTemplateCustomization `json:"customization"`
	TemplateName    string                             `json:"templateName"`
}

func (opts *InstallationTemplateUpdateOpts) FromResource(d *schema.ResourceData) *InstallationTemplateUpdateOpts {
	opts.TemplateName = d.Get("template_name").(string)
	customizations := d.Get("customization").([]interface{})
	if customizations != nil && len(customizations) == 1 {
		opts.Customization = (&InstallationTemplateCustomization{}).FromResource(d, "customization.0")
	}

	return opts
}

type InstallationTemplateCustomization struct {
	CustomHostname               *string `json:"customHostname,omitempty"`
	PostInstallationScriptLink   *string `json:"postInstallationScriptLink,omitempty"`
	PostInstallationScriptReturn *string `json:"postInstallationScriptReturn,omitempty"`
	SshKeyName                   *string `json:"sshKeyName,omitempty"`
}

func (v InstallationTemplateCustomization) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	custom_attr_set := false

	if v.CustomHostname != nil {
		obj["custom_hostname"] = *v.CustomHostname
		custom_attr_set = true
	}

	if v.PostInstallationScriptLink != nil {
		obj["post_installation_script_link"] = *v.PostInstallationScriptLink
		custom_attr_set = true
	}

	if v.PostInstallationScriptReturn != nil {
		obj["post_installation_script_return"] = *v.PostInstallationScriptReturn
		custom_attr_set = true
	}

	// dont return an object if nothing is set
	if custom_attr_set {
		return obj
	}

	return nil
}

type InstallationTemplateInputs struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Mandatory   bool     `json:"mandatory"`
	Default     *string  `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

func (v InstallationTemplateInputs) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = v.Name
	obj["description"] = v.Description
	obj["type"] = v.Type
	obj["mandatory"] = v.Mandatory
	obj["default"] = *v.Default
	obj["enum"] = &v.Enum
	return obj
}

type InstallTemplateLicenseItem struct {
	Name []string `json:"name,omitempty"`
	Url  string   `json:"url,omitempty"`
}

type InstallationTemplateLicense struct {
	Usage InstallTemplateLicenseItem `json:"usage,omitempty"`
	Os    InstallTemplateLicenseItem `json:"os,omitempty"`
}

func (v InstallationTemplateLicense) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["usage"] = []map[string]interface{}{v.Usage.ToMap()}
	obj["os"] = []map[string]interface{}{v.Os.ToMap()}
	return obj
}

func (v InstallTemplateLicenseItem) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = &v.Name
	obj["url"] = v.Url
	return obj
}

type InstallationTemplateProject struct {
	Usage InstallationTemplateProjectItem `json:"usage,omitempty"`
	Os    InstallationTemplateProjectItem `json:"os,omitempty"`
}

type InstallationTemplateProjectItem struct {
	Version      string   `json:"version,omitempty"`
	Url          string   `json:"url,omitempty"`
	ReleaseNotes string   `json:"release_notes,omitempty"`
	Name         string   `json:"name,omitempty"`
	Governance   []string `json:"governance,omitempty"`
}

func (v InstallationTemplateProject) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["usage"] = []map[string]interface{}{v.Usage.ToMap()}
	obj["os"] = []map[string]interface{}{v.Os.ToMap()}
	return obj
}

func (v InstallationTemplateProjectItem) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["version"] = v.Version
	obj["url"] = v.Url
	obj["release_notes"] = v.ReleaseNotes
	obj["name"] = v.Name
	obj["governance"] = &v.Governance
	return obj
}

func (opts *InstallationTemplateCustomization) FromResource(d *schema.ResourceData, parent string) *InstallationTemplateCustomization {
	opts.CustomHostname = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.custom_hostname", parent))
	opts.PostInstallationScriptLink = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script_link", parent))
	opts.PostInstallationScriptReturn = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script_return", parent))
	return opts
}

type Partition struct {
	Filesystem string       `json:"filesystem"`
	Mountpoint string       `json:"mountpoint"`
	Order      int          `json:"order"`
	Raid       *string      `json:"raid,omitempty"`
	Size       UnitAndValue `json:"size"`
	Type       string       `json:"type"`
	VolumeName *string      `json:"volumeName,omitempty"`
}

func (v Partition) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["filesystem"] = v.Filesystem
	obj["mountpoint"] = v.Mountpoint
	obj["order"] = v.Order

	if v.Raid != nil {
		obj["raid"] = fmt.Sprintf("raid%s", *v.Raid)
	}

	// always return size in MB
	obj["size"] = v.Size.Value
	obj["type"] = v.Type

	if v.VolumeName != nil {
		obj["volume_name"] = *v.VolumeName
	}

	return obj
}

type PartitionCreateOpts struct {
	Filesystem string  `json:"filesystem"`
	Mountpoint string  `json:"mountpoint"`
	Step       int     `json:"step"`
	Raid       *string `json:"raid,omitempty"`
	Size       int     `json:"size"`
	Type       string  `json:"type"`
	VolumeName *string `json:"volumeName,omitempty"`
}

func (opts *PartitionCreateOpts) FromResource(d *schema.ResourceData) *PartitionCreateOpts {
	opts.Filesystem = d.Get("filesystem").(string)
	opts.Mountpoint = d.Get("mountpoint").(string)
	opts.Step = d.Get("order").(int)

	if raid := helpers.GetNilStringPointerFromData(d, "raid"); raid != nil {
		raidValue := strings.ReplaceAll(*raid, "raid", "")
		opts.Raid = &raidValue
	}

	opts.Size = d.Get("size").(int)
	opts.Type = d.Get("type").(string)
	opts.VolumeName = helpers.GetNilStringPointerFromData(d, "volume_name")
	return opts
}

type PartitionUpdateOpts struct {
	Partition
}

func (opts *PartitionUpdateOpts) FromResource(d *schema.ResourceData) *PartitionUpdateOpts {
	opts.Filesystem = d.Get("filesystem").(string)
	opts.Mountpoint = d.Get("mountpoint").(string)
	opts.Order = d.Get("order").(int)

	if raid := helpers.GetNilStringPointerFromData(d, "raid"); raid != nil {
		raidValue := strings.ReplaceAll(*raid, "raid", "")
		opts.Raid = &raidValue
	}

	opts.Size = UnitAndValue{
		Unit:  "M",
		Value: d.Get("size").(int),
	}

	opts.Type = d.Get("type").(string)
	opts.VolumeName = helpers.GetNilStringPointerFromData(d, "volume_name")
	return opts
}

type HardwareRaid struct {
	Disks []string `json:"disks"`
	Mode  string   `json:"mode"`
	Name  string   `json:"name"`
	Step  int      `json:"step"`
}

func (v HardwareRaid) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["disks"] = v.Disks
	obj["mode"] = v.Mode
	obj["name"] = v.Name
	obj["step"] = v.Step

	return obj
}

type HardwareRaidCreateOrUpdateOpts struct {
	HardwareRaid
}

func (opts *HardwareRaidCreateOrUpdateOpts) FromResource(d *schema.ResourceData) *HardwareRaidCreateOrUpdateOpts {
	disks := d.Get("disks").([]interface{})
	opts.Disks = make([]string, len(disks))
	for i, disk := range disks {
		opts.Disks[i] = disk.(string)
	}

	opts.Mode = d.Get("mode").(string)
	opts.Name = d.Get("name").(string)
	opts.Step = d.Get("step").(int)
	return opts
}

type PartitionScheme struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

func (v PartitionScheme) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = v.Name
	obj["priority"] = v.Priority
	return obj
}

type PartitionSchemeCreateOrUpdateOpts struct {
	PartitionScheme
}

func (opts *PartitionSchemeCreateOrUpdateOpts) FromResource(d *schema.ResourceData) *PartitionSchemeCreateOrUpdateOpts {
	opts.Name = d.Get("name").(string)
	opts.Priority = d.Get("priority").(int)
	return opts
}
