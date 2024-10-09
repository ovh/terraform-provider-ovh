package ovh

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type DbaasLogsOutputOpensearchAlias struct {
	AliasId     string `json:"aliasId"`
	CreatedAt   string `json:"createdAt"`
	CurrentSize int    `json:"currentSize"`
	Description string `json:"description"`
	IsEditable  bool   `json:"isEditable"`
	Name        string `json:"name"`
	NbIndex     int    `json:"nbIndex"`
	NbStream    int    `json:"nbStream"`
	UpdatedAt   string `json:"updatedAt"`
}

func (v DbaasLogsOutputOpensearchAlias) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["alias_id"] = v.AliasId
	obj["created_at"] = v.CreatedAt
	obj["current_size"] = v.CurrentSize
	obj["description"] = v.Description
	obj["is_editable"] = v.IsEditable
	obj["name"] = v.Name
	obj["nb_index"] = v.NbIndex
	obj["nb_stream"] = v.NbStream
	obj["updated_at"] = v.UpdatedAt

	return obj
}

type DbaasLogsOutputOpensearchAliasCreateOps struct {
	Description string `json:"description"`
	Suffix      string `json:"suffix"`
}

func (opts *DbaasLogsOutputOpensearchAliasCreateOps) FromResource(d *schema.ResourceData) *DbaasLogsOutputOpensearchAliasCreateOps {
	opts.Description = d.Get("description").(string)
	opts.Suffix = d.Get("suffix").(string)
	return opts
}

type DbaasLogsOutputOpensearchAliasUpdateOps struct {
	Description string `json:"description"`
}

func (opts *DbaasLogsOutputOpensearchAliasUpdateOps) FromResource(d *schema.ResourceData) *DbaasLogsOutputOpensearchAliasUpdateOps {
	opts.Description = d.Get("description").(string)
	return opts
}

type DbaasLogsOutputOpensearchAliasIndexCreate struct {
	IndexID string `json:"indexId"`
}

type DbaasLogsOutputOpensearchAliasStreamCreate struct {
	StreamID string `json:"streamId"`
}
