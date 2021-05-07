package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type InstallationTemplate struct {
	AvailableLanguages         []string                           `json:"available_languages"`
	Beta                       *bool                              `json:"beta,omitempty"`
	BitFormat                  int                                `json:"bitFormat"`
	Category                   string                             `json:"category"`
	Customization              *InstallationTemplateCustomization `json:"customization,omitempty"`
	DefaultLanguage            string                             `json:"defaultLanguage"`
	Deprecated                 *bool                              `json:"deprecated,omitempty"`
	Description                string                             `json:"description"`
	Distribution               string                             `json:"distribution"`
	Family                     string                             `json:"family"`
	Filesystems                []string                           `json:"filesystems"`
	HardRaidConfiguration      *bool                              `json:"hardRaidConfigurtion,omitempty"`
	LastModification           *string                            `json:"last_modification"`
	LvmReady                   *bool                              `json:"lvmReady,omitempty"`
	SupportsDistributionKernel *bool                              `json:"supportsDistributionKernel,omitempty"`
	SupportsGptLabel           *bool                              `json:"supportsGptLabel,omitempty"`
	SupportsRTM                bool                               `json:"supportsRTM"`
	SupportsSqlServer          *bool                              `json:"supportsSqlServer,omitempty"`
	SupportsUEFI               *string                            `json:"supportsUEFI,omitempty"`
	TemplateName               string                             `json:"templateName"`
}

func (v InstallationTemplate) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["available_languages"] = v.AvailableLanguages

	if v.Beta != nil {
		obj["beta"] = *v.Beta
	}

	obj["bit_format"] = v.BitFormat
	obj["category"] = v.Category

	if v.Customization != nil {
		customization := v.Customization.ToMap()
		if customization != nil {
			obj["customization"] = []interface{}{customization}
		}
	}

	obj["default_language"] = v.DefaultLanguage

	if v.Deprecated != nil {
		obj["deprecated"] = *v.Deprecated
	}

	obj["description"] = v.Description
	obj["distribution"] = v.Distribution
	obj["family"] = v.Family
	obj["filesystems"] = v.Filesystems

	if v.HardRaidConfiguration != nil {
		obj["hard_raid_configuration"] = *v.HardRaidConfiguration
	}

	if v.LastModification != nil {
		obj["last_modification"] = *v.LastModification
	}

	if v.LvmReady != nil {
		obj["lvm_ready"] = *v.LvmReady
	}

	if v.SupportsDistributionKernel != nil {
		obj["supports_distribution_kernel"] = *v.SupportsDistributionKernel
	}

	if v.SupportsGptLabel != nil {
		obj["supports_gpt_label"] = *v.SupportsGptLabel
	}

	obj["supports_rtm"] = v.SupportsRTM

	if v.SupportsSqlServer != nil {
		obj["supports_sql_server"] = *v.SupportsSqlServer
	}

	if v.SupportsUEFI != nil {
		obj["supports_uefi"] = *v.SupportsUEFI
	}

	obj["template_name"] = v.TemplateName

	return obj
}

type InstallationTemplateCreateOpts struct {
	BaseTemplateName string `json:"baseTemplateName"`
	Name             string `json:"name"`
	DefaultLanguage  string `json:"defaultLanguage"`
}

func (opts *InstallationTemplateCreateOpts) FromResource(d *schema.ResourceData) *InstallationTemplateCreateOpts {
	opts.BaseTemplateName = d.Get("base_template_name").(string)
	opts.Name = d.Get("template_name").(string)
	opts.DefaultLanguage = d.Get("default_language").(string)
	return opts
}

type InstallationTemplateUpdateOpts struct {
	DefaultLanguage string                             `json:"defaultLanguage"`
	Customization   *InstallationTemplateCustomization `json:"customization"`
	TemplateName    string                             `json:"templateName"`
}

func (opts *InstallationTemplateUpdateOpts) FromResource(d *schema.ResourceData) *InstallationTemplateUpdateOpts {
	opts.TemplateName = d.Get("template_name").(string)
	opts.DefaultLanguage = d.Get("default_language").(string)

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
	UseDistributionKernel        *bool   `json:"useDistributionKernel,omitempty"`
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

	if v.SshKeyName != nil {
		obj["ssh_key_name"] = *v.SshKeyName
		custom_attr_set = true
	}

	if v.UseDistributionKernel != nil {
		obj["use_distribution_kernel"] = *v.UseDistributionKernel
		custom_attr_set = true
	}

	// dont return an object if nothing is set
	if custom_attr_set {
		return obj
	}

	return nil
}

func (opts *InstallationTemplateCustomization) FromResource(d *schema.ResourceData, parent string) *InstallationTemplateCustomization {
	opts.CustomHostname = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.custom_hostname", parent))
	opts.PostInstallationScriptLink = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script_link", parent))
	opts.PostInstallationScriptReturn = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.post_installation_script_return", parent))
	opts.SshKeyName = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.ssh_key_name", parent))
	opts.UseDistributionKernel = helpers.GetNilBoolPointerFromData(d, fmt.Sprintf("%s.use_distribution_kernel", parent))

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
