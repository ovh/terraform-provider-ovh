package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectVrackResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (v *CloudProjectVrackResponse) ToMap(d *schema.ResourceData) map[string]interface{} {
	obj := make(map[string]interface{})
	obj["id"] = v.Id
	obj["name"] = v.Name
	obj["description"] = v.Description
	return obj
}
