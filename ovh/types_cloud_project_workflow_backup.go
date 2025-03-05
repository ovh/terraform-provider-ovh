package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type CloudProjectWorkflowBackupCreateOpts struct {
	Cron              *string `json:"cron"`
	InstanceId        *string `json:"instanceId"`
	MaxExecutionCount *int64  `json:"maxExecutionCount,omitempty"`
	Name              *string `json:"name"`
	Rotation          *int64  `json:"rotation"`
}

type CloudProjectWorkflowBackupResponse struct {
	BackupName string `json:"backupName"`
	CreatedAt  string `json:"createdAt"`
	Cron       string `json:"cron"`
	Id         string `json:"id"`
	InstanceId string `json:"instanceId"`
	Name       string `json:"name"`
}

func (opts *CloudProjectWorkflowBackupCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectWorkflowBackupCreateOpts {
	opts.Cron = helpers.GetNilStringPointerFromData(d, "cron")
	opts.InstanceId = helpers.GetNilStringPointerFromData(d, "instance_id")
	opts.MaxExecutionCount = helpers.GetNilInt64PointerFromData(d, "max_execution_count")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")
	opts.Rotation = helpers.GetNilInt64PointerFromData(d, "rotation")
	return opts
}

func (v CloudProjectWorkflowBackupResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["backup_name"] = v.BackupName
	obj["created_at"] = v.CreatedAt
	obj["cron"] = v.Cron
	obj["id"] = v.Id
	obj["instance_id"] = v.InstanceId
	obj["name"] = v.Name
	return obj
}
