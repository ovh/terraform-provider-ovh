package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type CloudProject struct {
	Access             string  `json:"access"`
	Description        *string `json:"description"`
	ProjectName        *string `json:"projectName"`
	ProjectId          string  `json:"project_id"`
	Status             string  `json:"status"`
	IamResourceDetails `json:"iam"`
}

func (v CloudProject) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["access"] = v.Access
	obj["project_id"] = v.ProjectId
	obj["status"] = v.Status
	obj["urn"] = v.URN

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

type CloudProjectrolesResponse struct {
	Roles []CloudProjectroleResponse `json:"roles"`
}

type CloudProjectroleResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type CloudProjectroleUpdate struct {
	RolesIds []string `json:"rolesIds"`
}
