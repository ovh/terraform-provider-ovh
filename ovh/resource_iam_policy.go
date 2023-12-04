package ovh

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIamPolicy() *schema.Resource {
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
				Optional: true,
			},
			"identities": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		ReadContext:   resourceIamPolicyRead,
		CreateContext: resourceIamPolicyCreate,
		UpdateContext: resourceIamPolicyUpdate,
		DeleteContext: resourceIamPolicyDelete,
	}
}

func resourceIamPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	var pol IamPolicy
	err := config.OVHClient.GetWithContext(ctx, "/v2/iam/policy/"+url.PathEscape(d.Id()), &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range pol.ToMap() {
		d.Set(k, v)
	}
	return nil
}

func resourceIamPolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	req := prepareIamPolicyCall(d)

	var pol IamPolicy
	err := config.OVHClient.PostWithContext(ctx, "/v2/iam/policy", req, &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range pol.ToMap() {
		d.Set(k, v)
	}

	d.SetId(pol.Id)
	return nil
}

func resourceIamPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	req := prepareIamPolicyCall(d)

	var pol IamPolicy
	err := config.OVHClient.PutWithContext(ctx, "/v2/iam/policy/"+url.PathEscape(d.Id()), req, &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range pol.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceIamPolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	err := config.OVHClient.DeleteWithContext(ctx, "/v2/iam/policy/"+url.PathEscape(d.Id()), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func prepareIamPolicyCall(d *schema.ResourceData) IamPolicy {
	var out IamPolicy

	out.Name = d.Get("name").(string)
	out.Description = d.Get("description").(string)
	ids := d.Get("identities").(*schema.Set)

	for _, id := range ids.List() {
		out.Identities = append(out.Identities, id.(string))
	}

	res := d.Get("resources").(*schema.Set)
	for _, r := range res.List() {
		out.Resources = append(out.Resources, IamResource{URN: r.(string)})
	}

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
