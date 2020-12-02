package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIPLoadbalancingRouteHTTPRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingRouteHTTPRuleCreate,
		Read:   resourceIPLoadbalancingRouteHTTPRuleRead,
		Update: resourceIPLoadbalancingRouteHTTPRuleUpdate,
		Delete: resourceIPLoadbalancingRouteHTTPRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpLoadbalancingHttpRouteRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"match": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(v.(string), []string{"contains", "endswith", "exists", "in", "internal", "is", "matches", "startswith"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"negate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"pattern": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIpLoadbalancingHttpRouteRuleImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("Import Id is not service_name/route_id/rule id formatted")
	}
	serviceName := splitId[0]
	routeID := splitId[1]
	ruleID := splitId[2]

	d.SetId(ruleID)
	d.Set("route_id", routeID)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceIPLoadbalancingRouteHTTPRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	rule := &IPLoadbalancingRouteHTTPRule{
		DisplayName: d.Get("display_name").(string),
		Field:       d.Get("field").(string),
		Match:       d.Get("match").(string),
		Negate:      d.Get("negate").(bool),
		Pattern:     d.Get("pattern").(string),
		SubField:    d.Get("sub_field").(string),
	}

	service := d.Get("service_name").(string)
	routeID := d.Get("route_id").(string)
	resp := &IPLoadbalancingRouteHTTPRule{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s/rule", service, routeID)

	err := config.OVHClient.Post(endpoint, rule, resp)
	if err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", resp.RuleID))

	return resourceIPLoadbalancingRouteHTTPRuleRead(d, meta)
}

func resourceIPLoadbalancingRouteHTTPRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	routeID := d.Get("route_id").(string)
	r := &IPLoadbalancingRouteHTTPRule{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s/rule/%s", service, routeID, d.Id())

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("display_name", r.DisplayName)
	d.Set("field", r.Field)
	d.Set("match", r.Match)
	d.Set("negate", r.Negate)
	d.Set("pattern", r.Pattern)
	d.Set("sub_field", r.SubField)

	return nil
}

func resourceIPLoadbalancingRouteHTTPRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	routeID := d.Get("route_id").(string)

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s/rule/%s", service, routeID, d.Id())

	rule := &IPLoadbalancingRouteHTTPRule{
		DisplayName: d.Get("display_name").(string),
		Field:       d.Get("field").(string),
		Match:       d.Get("match").(string),
		Negate:      d.Get("negate").(bool),
		Pattern:     d.Get("pattern").(string),
		SubField:    d.Get("sub_field").(string),
	}

	err := config.OVHClient.Put(endpoint, rule, nil)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}

	return resourceIPLoadbalancingRouteHTTPRuleRead(d, meta)
}

func resourceIPLoadbalancingRouteHTTPRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	routeID := d.Get("route_id").(string)

	r := &IPLoadbalancingRouteHTTPRule{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s/rule/%s", service, routeID, d.Id())

	err := config.OVHClient.Delete(endpoint, &r)
	if err != nil {
		return fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	return nil
}
