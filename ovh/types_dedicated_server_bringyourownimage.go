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
	Enable       *bool              `json:enable`
	Hostname     *string            `json:hostname,omitempty`
	SSHKey       *string            `json:sshKey,omitempty`
	UserData     *string            `json:userData,omitempty`
	UserMetadata *map[string]string `json:userMetadata,omitempty`
}

func (opts *DedicatedServerBringYourOwnImageConfigDrive) FromResource(d *schema.ResourceData, parent string) (*DedicatedServerBringYourOwnImageConfigDrive, error) {
	/*configDriveSet := d.Get(parent).(*schema.Set).List()
	configDrive := configDriveSet[0]

	if len(configDriveSet) == 0 {
		return opts, nil
	}
	if len(configDriveSet) > 2 {
		return opts, errors.New("resource config_drive cannot have more than 2 elements")
	}*/

	//for _, to := range configDriveSet {
	/*

		Enable       *bool              `json:enable`
		Hostname     *string            `json:hostname`
		SSHKey       *string            `json:sshKey`
		UserData     *string            `json:userData`
		UserMetadata *map[string]string `json:userMetadata`
	*/
	/*metadata := to.(map[string]interface{})["metadata"].(*schema.Set).List()[0]
		annotations := metadata.(map[string]interface{})["annotations"].(map[string]interface{})
		labels := metadata.(map[string]interface{})["labels"].(map[string]interface{})
		finalizers := metadata.(map[string]interface{})["finalizers"].([]interface{})

		spec := configDrive.(map[string]interface{})["spec"].(*schema.Set).List()[0]
		taints := spec.(map[string]interface{})["taints"].([]interface{})
		unschedulable := spec.(map[string]interface{})["unschedulable"].(bool)

		if len(annotations) == 0 && len(labels) == 0 && len(finalizers) == 0 && len(taints) == 0 && unschedulable == false {
			// is empty
		} else {
			configDrive = to
			break
		}
	}*/

	opts.Enable = helpers.GetNilBoolPointerFromData(d, fmt.Sprintf("%s.enable", parent))
	opts.Hostname = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.hostname", parent))
	opts.SSHKey = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.ssh_key", parent))
	opts.UserData = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.user_data", parent))
	// TODO
	//opts.UserMetadata = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.user_metadata", parent))

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

	//opts.ConfigDrive = (&DedicatedServerBringYourOwnImageConfigDrive{}).FromResource(d, "config_drive")

	return opts, nil
}
