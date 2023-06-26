package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type DedicatedServerBringYourOwnImageStatus struct {
	CheckSum   string `json:"checksum"`
	Message    string `json:"message"`
	ServerName string `json:"servername"`
	Status     string `json:"status"`
}

type DedicatedServerBringYourOwnImageConfigDrive struct {
	Enable        *bool              `json:enable`
	Hostname      *string            `json:hostname,omitempty`
	SSHKey        *string            `json:sshKey,omitempty`
	UserData      *string            `json:userData,omitempty`
	UserMetadatas *map[string]string `json:userMetadatas,omitempty`
}

func (opts *DedicatedServerBringYourOwnImageConfigDrive) FromResource(d *schema.ResourceData, parent string) (*DedicatedServerBringYourOwnImageConfigDrive, error) {
	opts.Enable = helpers.GetNilBoolPointerFromData(d, fmt.Sprintf("%s.enable", parent))
	opts.Hostname = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.hostname", parent))
	opts.SSHKey = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.ssh_key", parent))
	opts.UserData = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.user_data", parent))

	userMetadatas := map[string]string{}
	providedUserMetadatas := d.Get(fmt.Sprintf("%s.user_metadatas", parent))
	if providedUserMetadatas != nil {
		for key, value := range providedUserMetadatas.(map[string]interface{}) {
			userMetadatas[key] = fmt.Sprintf("%v", value)
		}
	}
	opts.UserMetadatas = &userMetadatas

	return opts, nil
}

type DedicatedServerBringYourOwnImageCreateOpts struct {
	Url          string                                       `json:"URL"`
	CheckSum     string                                       `json:"checkSum,omitempty"`
	CheckSumType string                                       `json:"checkSumType,omitempty"`
	ConfigDrive  *DedicatedServerBringYourOwnImageConfigDrive `json:"configdrive"`
	Description  string                                       `json:"description,omitempty"`
	DiskGroupId  float64                                      `json:"diskGroupId,omitempty"`
	HttpHeader   map[string]string                            `json:"httpHeader,omitempty"`
	Type         string                                       `json:"type"`
}

func (opts *DedicatedServerBringYourOwnImageCreateOpts) FromResource(d *schema.ResourceData) (*DedicatedServerBringYourOwnImageCreateOpts, error) {
	opts.Url = d.Get("url").(string)
	opts.CheckSum = d.Get("check_sum").(string)
	opts.CheckSumType = d.Get("check_sum_type").(string)
	opts.Description = d.Get("description").(string)
	opts.DiskGroupId = d.Get("disk_group_id").(float64)
	opts.Type = d.Get("type").(string)

	httpHeaders := d.Get("http_headers").(map[string]interface{})
	for key, value := range httpHeaders {
		opts.HttpHeader[key] = fmt.Sprintf("%v", value)
	}

	configDrive, err := (&DedicatedServerBringYourOwnImageConfigDrive{}).FromResource(d, "config_drive.0")
	if err != nil {
		return nil, err
	}
	opts.ConfigDrive = configDrive

	return opts, nil
}
