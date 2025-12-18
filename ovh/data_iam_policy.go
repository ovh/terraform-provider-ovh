package ovh

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIamPolicy() *schema.Resource {
	// Define the deepest level first (e.g., 3 levels deep)
	conditionLevel3Schema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operator for this condition (MATCH, AND, OR, NOT)",
			},
			"values": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Key-value pairs to match (e.g., resource.Tag(name), date(Europe/Paris).WeekDay, request.IP)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// No further "condition" Elem here to limit depth
		},
	}

	// Define the second level of conditions, pointing to the third level
	conditionLevel2Schema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operator for this condition (MATCH, AND, OR, NOT)",
			},
			"values": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Key-value pairs to match (e.g., resource.Tag(name), date(Europe/Paris).WeekDay, request.IP)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of nested conditions. This is the recursive part.",
				Elem:        conditionLevel3Schema, // Points to the next level
			},
		},
	}

	// Define the first level of conditions, pointing to the second level
	conditionLevel1Schema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operator for this condition (MATCH, AND, OR, NOT)",
			},
			"values": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Key-value pairs to match (e.g., resource.Tag(name), date(Europe/Paris).WeekDay, request.IP)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of nested conditions. This is the recursive part.",
				Elem:        conditionLevel2Schema, // Points to the next level
			},
		},
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"identities": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resources": {
				Type:     schema.TypeSet,
				Computed: true,
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
			"permissions_groups": {
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
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration date of the policy, after this date it will no longer be applied",
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Conditions restrict permissions following resources, date or customer's information",
				Elem:        conditionLevel1Schema, // The top-level conditions use the first level schema
			},
		},
		ReadContext: datasourceIamPolicyRead,
	}
}

func datasourceIamPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	id := d.Get("id").(string)

	var pol IamPolicy
	err := config.OVHClient.GetWithContext(ctx, "/v2/iam/policy/"+url.PathEscape(id), &pol)
	if err != nil {
		return diag.FromErr(err)
	}

	// Explicitly set the new attributes to ensure they're available
	if pol.ExpiredAt != "" {
		d.Set("expired_at", pol.ExpiredAt)
	}
	if pol.Conditions != nil {
		d.Set("conditions", []interface{}{conditionsToMap(pol.Conditions)})
	}

	d.SetId(id)
	return nil
}
