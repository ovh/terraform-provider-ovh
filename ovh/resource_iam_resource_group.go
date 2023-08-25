package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIamResourceGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resources": {
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
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, rd *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// Before importing, check that the resource group is not read-only.
				config := meta.(*Config)
				id := rd.Id()

				var pol IamResourceGroup
				err := config.OVHClient.GetWithContext(ctx, fmt.Sprintf("/v2/iam/resourceGroup/%s?details=true", url.PathEscape(id)), &pol)
				if err != nil {
					return nil, err
				}

				if pol.ReadOnly {
					return nil, fmt.Errorf("resource_group '%s' is read-only", id)
				}

				return []*schema.ResourceData{rd}, nil
			},
		},
		ReadContext:   resourceIamResourceGroupRead,
		CreateContext: resourceIamResourceGroupCreate,
		UpdateContext: resourceIamResourceGroupUpdate,
		DeleteContext: resourceIamResourceGroupDelete,
	}
}

func resourceIamResourceGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Id()

	var pol IamResourceGroup
	err := config.OVHClient.GetWithContext(ctx, fmt.Sprintf("/v2/iam/resourceGroup/%s?details=true", url.PathEscape(id)), &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	var urns []string
	for _, r := range pol.Resources {
		urns = append(urns, r.URN)
	}
	d.Set("resources", urns)

	d.Set("name", pol.Name)
	d.Set("owner", pol.Owner)
	d.Set("created_at", pol.CreatedAt)
	d.Set("updated_at", pol.UpdatedAt)
	d.Set("read_only", pol.ReadOnly)
	d.Set("urn", pol.URN)

	return nil
}

func resourceIamResourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	var grp IamResourceGroup

	grp.Name = d.Get("name").(string)
	urns := d.Get("resources").(*schema.Set)
	for _, r := range urns.List() {
		urn := r.(string)
		grp.Resources = append(grp.Resources, IamResourceDetails{URN: urn})
	}

	var grpOut IamResourceGroup
	err := config.OVHClient.PostWithContext(ctx, "/v2/iam/resourceGroup", grp, &grpOut)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(grpOut.ID)

	return resourceIamResourceGroupRead(ctx, d, meta)
}

func resourceIamResourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Id()

	var pol IamResourceGroup

	pol.Name = d.Get("name").(string)
	urns := d.Get("resources").(*schema.Set)
	for _, r := range urns.List() {
		urn := r.(string)
		pol.Resources = append(pol.Resources, IamResourceDetails{URN: urn})
	}

	err := config.OVHClient.PutWithContext(ctx, "/v2/iam/resourceGroup/"+url.PathEscape(id), &pol, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIamResourceGroupRead(ctx, d, meta)
}

func resourceIamResourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Id()

	err := config.OVHClient.DeleteWithContext(ctx, "/v2/iam/resourceGroup/"+url.PathEscape(id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
