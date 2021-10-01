package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIPLoadbalancingTcpRouteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingTcpRouteRuleCreate,
		Read:   resourceIPLoadbalancingTcpRouteRuleRead,
		Update: resourceIPLoadbalancingTcpRouteRuleUpdate,
		Delete: resourceIPLoadbalancingTcpRouteRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpLoadbalancingTcpRouteRuleImportState,
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
				Computed: true,
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

func resourceIpLoadbalancingTcpRouteRuleImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceIPLoadbalancingTcpRouteRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	rule := (&IPLoadbalancingRouteRuleOpts{}).FromResource(d)

	serviceName := d.Get("service_name").(string)
	routeId := d.Get("route_id").(string)
	resp := &IPLoadbalancingRouteRule{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s/rule",
		url.PathEscape(serviceName),
		url.PathEscape(routeId),
	)

	if err := config.OVHClient.Post(endpoint, rule, resp); err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", resp.RuleId))

	return resourceIPLoadbalancingTcpRouteRuleRead(d, meta)
}

func resourceIPLoadbalancingTcpRouteRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	routeId := d.Get("route_id").(string)
	r := &IPLoadbalancingRouteRule{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s/rule/%s",
		url.PathEscape(serviceName),
		url.PathEscape(routeId),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		if k != "rule_id" {
			d.Set(k, v)
		}
	}

	return nil
}

func resourceIPLoadbalancingTcpRouteRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	rule := (&IPLoadbalancingRouteRuleOpts{}).FromResource(d)
	serviceName := d.Get("service_name").(string)
	routeId := d.Get("route_id").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s/rule/%s",
		url.PathEscape(serviceName),
		url.PathEscape(routeId),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Put(endpoint, rule, nil); err != nil {
		return fmt.Errorf("calling PUT %s:\n\t %s", endpoint, err.Error())
	}

	return resourceIPLoadbalancingTcpRouteRuleRead(d, meta)
}

func resourceIPLoadbalancingTcpRouteRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	routeId := d.Get("route_id").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s/rule/%s",
		url.PathEscape(serviceName),
		url.PathEscape(routeId),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId("")
	return nil
}
