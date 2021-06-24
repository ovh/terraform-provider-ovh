package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type CloudProject struct {
	Access      string  `json:"access"`
	Description *string `json:"description"`
	ProjectName *string `json:"projectName"`
	ProjectId   string  `json:"project_id"`
	Status      string  `json:"status"`
}

func (v CloudProject) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["access"] = v.Access
	obj["project_id"] = v.ProjectId
	obj["status"] = v.Status

	if v.Description != nil {
		obj["description"] = *v.Description
	}

	if v.ProjectName != nil {
		obj["project_name"] = *v.ProjectName
	}

	return obj
}

type CloudProjectUpdateOpts struct {
	Description *string `json:"description,omitempty"`
}

func (opts *CloudProjectUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectUpdateOpts {
	opts.Description = helpers.GetNilStringPointerFromData(d, "description")

	return opts
}

type CloudProjectConfirmTerminationOpts struct {
	Token string `json:"token"`
}
