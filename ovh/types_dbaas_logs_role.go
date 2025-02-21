package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DbaasLogsRoleCreateOpts defines the options for creating a role.
type DbaasLogsRoleCreateOpts struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (opts *DbaasLogsRoleCreateOpts) FromResource(d *schema.ResourceData) *DbaasLogsRoleCreateOpts {
	opts.Name = d.Get("name").(string)
	opts.Description = d.Get("description").(string)
	return opts
}

// DbaasLogsRoleUpdateOpts defines the options for updating a role.
type DbaasLogsRoleUpdateOpts struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (opts *DbaasLogsRoleUpdateOpts) FromResource(d *schema.ResourceData) *DbaasLogsRoleUpdateOpts {
	opts.Name = d.Get("name").(string)
	opts.Description = d.Get("description").(string)
	return opts
}

// DbaasLogsRole represents the role resource as returned by the API.
type DbaasLogsRole struct {
	CreatedAt    string `json:"createdAt"`
	Description  string `json:"description"`
	Name         string `json:"name"`
	NbMember     int    `json:"nbMember"`
	NbPermission int    `json:"nbPermission"`
	RoleId       string `json:"roleId"`
	UpdatedAt    string `json:"updatedAt"`
}

func (r *DbaasLogsRole) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"created_at":    r.CreatedAt,
		"description":   r.Description,
		"name":          r.Name,
		"nb_member":     r.NbMember,
		"nb_permission": r.NbPermission,
		"role_id":       r.RoleId,
		"updated_at":    r.UpdatedAt,
	}
}
