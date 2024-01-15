package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIamPermissionsGroup() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: func(rd *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{rd}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allow": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"except": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deny": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext:   resourceIamPermissionsGroupsRead,
		CreateContext: resourceIamPermissionsGroupCreate,
		UpdateContext: resourceIamPermissionsGroupUpdate,
		DeleteContext: resourceIamPermissionsGroupDelete,
	}
}

func resourceIamPermissionsGroupsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	var pol IamPermissionsGroup
	err := config.OVHClient.GetWithContext(ctx, "/v2/iam/permissionsGroup/"+url.PathEscape(d.Id()), &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range pol.ToMap() {
		d.Set(k, v)
	}
	return nil
}

func resourceIamPermissionsGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	req := prepareIamPermissionGroupCall(d)

	var pol IamPermissionsGroup
	err := config.OVHClient.PostWithContext(ctx, "/v2/iam/permissionsGroup", req, &pol)
	if err != nil {
		warn := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "requestDetails",
			Detail:   fmt.Sprint(req),
		}
		diag := diag.FromErr(err)
		return append(diag, warn)
	}

	for k, v := range pol.ToMap() {
		d.Set(k, v)
	}

	d.SetId(pol.Urn)
	return nil
}

func resourceIamPermissionsGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	req := prepareIamPermissionGroupCall(d)

	var pol IamPermissionsGroup
	err := config.OVHClient.PutWithContext(ctx, "/v2/iam/permissionsGroup/"+url.PathEscape(d.Id()), req, &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range pol.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceIamPermissionsGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	err := config.OVHClient.DeleteWithContext(ctx, "/v2/iam/permissionsGroup/"+url.PathEscape(d.Id()), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func prepareIamPermissionGroupCall(d *schema.ResourceData) IamPermissionsGroup {
	var out IamPermissionsGroup

	out.Name = d.Get("name").(string)
	out.Description = d.Get("description").(string)

	if allows, ok := d.GetOk("allow"); ok {
		for _, a := range allows.(*schema.Set).List() {
			out.Permissions.Allow = append(out.Permissions.Allow, IamAction{Action: a.(string)})
		}
	}

	if except, ok := d.GetOk("except"); ok {
		for _, e := range except.(*schema.Set).List() {
			out.Permissions.Except = append(out.Permissions.Except, IamAction{Action: e.(string)})
		}
	}

	if deny, ok := d.GetOk("deny"); ok {
		for _, e := range deny.(*schema.Set).List() {
			out.Permissions.Deny = append(out.Permissions.Deny, IamAction{Action: e.(string)})
		}
	}
	return out
}
