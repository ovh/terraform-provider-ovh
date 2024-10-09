package ovh

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type DbaasLogsOutputOpensearchIndex struct {
	AlertNotifyEnabled *bool  `json:"alertNotifyEnabled"`
	CreatedAt          string `json:"createdAt"`
	CurrentSize        int    `json:"currentSize"`
	Description        string `json:"description"`
	IndexId            string `json:"indexId"`
	IsEditable         bool   `json:"isEditable"`
	MaxSize            int    `json:"maxSize"`
	Name               string `json:"name"`
	NbShard            int    `json:"nbShard"`
	UpdatedAt          string `json:"updatedAt"`
}

func (v DbaasLogsOutputOpensearchIndex) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["alert_notify_enabled"] = v.AlertNotifyEnabled
	obj["created_at"] = v.CreatedAt
	obj["current_size"] = v.CurrentSize
	obj["description"] = v.Description
	obj["index_id"] = v.IndexId
	obj["is_editable"] = v.IsEditable
	obj["name"] = v.Name
	obj["nb_shard"] = v.NbShard
	obj["updated_at"] = v.UpdatedAt

	return obj
}

type DbaasLogsOutputOpensearchIndexCreateOps struct {
	Description string `json:"description"`
	NbShard     int    `json:"nbShard"`
	Suffix      string `json:"suffix"`
}

func (opts *DbaasLogsOutputOpensearchIndexCreateOps) FromResource(d *schema.ResourceData) *DbaasLogsOutputOpensearchIndexCreateOps {
	opts.Description = d.Get("description").(string)
	opts.NbShard = d.Get("nb_shard").(int)
	opts.Suffix = d.Get("suffix").(string)
	return opts
}

type DbaasLogsOutputOpensearchIndexUpdateOps struct {
	AlertNotifyEnabled bool   `json:"alertNotifyEnabled"`
	Description        string `json:"description"`
}

func (opts *DbaasLogsOutputOpensearchIndexUpdateOps) FromResource(d *schema.ResourceData) *DbaasLogsOutputOpensearchIndexUpdateOps {
	opts.Description = d.Get("description").(string)
	return opts
}
