package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIPLoadbalancingRouteHTTPRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingRouteHTTPRuleCreate,
		Read:   resourceIPLoadbalancingRouteHTTPRuleRead,
		Update: resourceIPLoadbalancingRouteHTTPRuleUpdate,
		Delete: resourceIPLoadbalancingRouteHTTPRuleDelete,

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
					err := validateStringEnum(v.(string), []string{"contains", "endswith", "exists", "in", "internal", "is", "matches", "startswith"})
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

//IPLoadbalancingRouteHTTPRule HTTP Route Rule
type IPLoadbalancingRouteHTTPRule struct {
	RuleID      int    `json:"ruleId,omitempty"`      //Id of your rule
	RouteID     int    `json:"routeId,omitempty"`     //Id of your route
	DisplayName string `json:"displayName,omitempty"` //Human readable name for your rule
	Field       string `json:"field,omitempty"`       //Name of the field to match like "protocol" or "host". See "/ipLoadbalancing/{serviceName}/availableRouteRules" for a list of available rules
	Match       string `json:"match,omitempty"`       //Matching operator. Not all operators are available for all fields. See "/ipLoadbalancing/{serviceName}/availableRouteRules"
	Negate      bool   `json:"negate,omitempty"`      //Invert the matching operator effect
	Pattern     string `json:"pattern,omitempty"`     //Value to match against this match. Interpretation if this field depends on the match and field
	SubField    string `json:"subField,omitempty"`    //Name of sub-field, if applicable. This may be a Cookie or Header name for instance
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

	err := config.OVHClient.Get(endpoint, &r)
	if err != nil {
		return CheckDeleted(d, err, endpoint)
	}

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
