package ovh

import (
	"fmt"
)

type InstallationTemplate struct {
	BitFormat             int                          `json:"bitFormat,omitempty"`
	Category              string                       `json:"category,omitempty"`
	Description           string                       `json:"description,omitempty"`
	Distribution          string                       `json:"distribution,omitempty"`
	EndOfInstall          string                       `json:"endOfInstall,omitempty"`
	Family                string                       `json:"family,omitempty"`
	Filesystems           []string                     `json:"filesystems"`
	Inputs                []InstallationTemplateInputs `json:"inputs,omitempty"`
	License               *InstallationTemplateLicense `json:"license,omitempty"`
	LvmReady              *bool                        `json:"lvmReady,omitempty"`
	NoPartitioning        bool                         `json:"noPartitioning,omitempty"`
	Project               *InstallationTemplateProject `json:"project,omitempty"`
	SoftRaidOnlyMirroring bool                         `json:"soft_raid_only_mirroring,omitempty"`
	Subfamily             string                       `json:"subfamily,omitempty"`
	TemplateName          string                       `json:"templateName"`
}

func (v *InstallationTemplate) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["bit_format"] = v.BitFormat
	obj["category"] = v.Category

	obj["description"] = v.Description
	obj["distribution"] = v.Distribution
	obj["end_of_install"] = v.EndOfInstall
	obj["family"] = v.Family
	obj["filesystems"] = v.Filesystems

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
