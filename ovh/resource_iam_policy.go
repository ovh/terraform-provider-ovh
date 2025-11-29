package ovh

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIamPolicy() *schema.Resource {
	// Define the deepest level first (e.g., 3 levels deep)
	conditionLevel3Schema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Operator for this condition (MATCH, AND, OR, NOT)",
			},
			"values": {
				Type:        schema.TypeMap,
				Optional:    true,
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
				Required:    true,
				Description: "Operator for this condition (MATCH, AND, OR, NOT)",
			},
			"values": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Key-value pairs to match (e.g., resource.Tag(name), date(Europe/Paris).WeekDay, request.IP)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"condition": {
				Type:        schema.TypeList,
				Optional:    true,
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
				Required:    true,
				Description: "Operator for this condition (MATCH, AND, OR, NOT)",
			},
			"values": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Key-value pairs to match (e.g., resource.Tag(name), date(Europe/Paris).WeekDay, request.IP)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"condition": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of nested conditions. This is the recursive part.",
				Elem:        conditionLevel2Schema, // Points to the next level
			},
		},
	}

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
				Optional:    true,
				Description: "Expiration date of the policy, after this date it will no longer be applied",
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Conditions restrict permissions following resources, date or customer's information",
				Elem:        conditionLevel1Schema, // The top-level conditions use the first level schema
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

	if permGrps, ok := d.GetOk("permissions_groups"); ok {
		for _, e := range permGrps.(*schema.Set).List() {
			out.PermissionsGroups = append(out.PermissionsGroups, PermissionGroup{Urn: e.(string)})
		}
	}

	if expiredAt, ok := d.GetOk("expired_at"); ok {
		out.ExpiredAt = expiredAt.(string)
	}

	if conditions, ok := d.GetOk("conditions"); ok {
		out.Conditions = expandConditions(conditions.([]interface{}))
	}

	return out
}

// expandConditions converts Terraform schema data to IamConditions
func expandConditions(tfList []interface{}) *IamConditions {
	if len(tfList) == 0 {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})
	conditions := &IamConditions{
		Operator: tfMap["operator"].(string),
	}

	if values, ok := tfMap["values"].(map[string]interface{}); ok {
		conditions.Values = make(map[string]string)
		for k, v := range values {
			conditions.Values[k] = v.(string)
		}
	}

	if condList, ok := tfMap["condition"].([]interface{}); ok {
		for _, c := range condList {
			condMap := c.(map[string]interface{})
			condition := expandCondition(condMap)
			conditions.Conditions = append(conditions.Conditions, condition)
		}
	}
	return conditions
}

// expandCondition recursively expands a single condition
func expandCondition(tfMap map[string]interface{}) *IamCondition {
	condition := &IamCondition{
		Operator: tfMap["operator"].(string),
	}
	if values, ok := tfMap["values"].(map[string]interface{}); ok {
		condition.Values = make(map[string]string)
		for k, v := range values {
			condition.Values[k] = v.(string)
		}
	}

	// Handle nested conditions recursively
	if nestedList, ok := tfMap["condition"].([]interface{}); ok {
		for _, nested := range nestedList {
			nestedMap := nested.(map[string]interface{})
			condition.Conditions = append(condition.Conditions, expandCondition(nestedMap))
		}
	}
	return condition
}
