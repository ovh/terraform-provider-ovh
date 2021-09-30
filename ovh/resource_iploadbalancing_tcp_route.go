package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIPLoadbalancingTcpRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingTcpRouteCreate,
		Read:   resourceIPLoadbalancingTcpRouteRead,
		Update: resourceIPLoadbalancingTcpRouteUpdate,
		Delete: resourceIPLoadbalancingTcpRouteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpLoadbalancingTcpRouteImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your IP load balancing",
				Required:    true,
				ForceNew:    true,
			},
			"action": {
				Type:        schema.TypeList,
				Description: "Action triggered when all rules match",
				Required:    true,
				ForceNew:    false,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:        schema.TypeString,
							Description: "Farm ID for \"farm\" action type, empty for others",
							Optional:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Action to trigger if all the rules of this route matches",
							Required:    true,
						},
					},
				},
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Human readable name for your route, this field is for you",
				Optional:    true,
			},
			"frontend_id": {
				Type:        schema.TypeInt,
				Description: "Route traffic for this frontend",
				Optional:    true,
				Computed:    true,
			},
			"weight": {
				Type:        schema.TypeInt,
				Description: "Route priority ([0..255]). 0 if null. Highest priority routes are evaluated last. Only the first matching route will trigger an action",
				Optional:    true,
				Computed:    true,
			},

			//computed
			"status": {
				Type:        schema.TypeString,
				Description: "Route status. Routes in \"ok\" state are ready to operate",
				Computed:    true,
			},
			"rules": {
				Type:        schema.TypeList,
				Description: "List of rules to match to trigger action",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:        schema.TypeString,
							Description: "Name of the field to match like \"protocol\" or \"host\". See \"/ipLoadbalancing/{serviceName}/route/availableRules\" for a list of available rules",
							Computed:    true,
						},
						"match": {
							Type:        schema.TypeString,
							Description: "Matching operator. Not all operators are available for all fields. See \"/availableRules\"",
							Computed:    true,
						},
						"negate": {
							Type:        schema.TypeBool,
							Description: "Invert the matching operator effect",
							Computed:    true,
						},
						"pattern": {
							Type:        schema.TypeString,
							Description: "Value to match against this match. Interpretation if this field depends on the match and field",
							Computed:    true,
						},
						"rule_id": {
							Type:        schema.TypeInt,
							Description: "Id of your rule",
							Computed:    true,
						},
						"sub_field": {
							Type:        schema.TypeString,
							Description: "Name of sub-field, if applicable. This may be a Cookie or Header name for instance",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceIpLoadbalancingTcpRouteImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/route id formatted")
	}
	serviceName := splitId[0]
	routeId := splitId[1]
	d.SetId(routeId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceIPLoadbalancingTcpRouteCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	route := (&IPLoadbalancingTcpRouteOpts{}).FromResource(d)
	resp := &IPLoadbalancingTcpRoute{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Post(endpoint, route, resp); err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", resp.RouteId))

	return resourceIPLoadbalancingTcpRouteRead(d, meta)
}

func resourceIPLoadbalancingTcpRouteRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	r := &IPLoadbalancingTcpRoute{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}
	// set resource attributes
	for k, v := range r.ToMap() {
		if k != "route_id" {
			d.Set(k, v)
		}
	}

	return nil
}

func resourceIPLoadbalancingTcpRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	route := (&IPLoadbalancingTcpRouteOpts{}).FromResource(d)

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Put(endpoint, route, nil); err != nil {
		return fmt.Errorf("calling PUT %s:\n\t %s", endpoint, err.Error())
	}

	return resourceIPLoadbalancingTcpRouteRead(d, meta)
}

func resourceIPLoadbalancingTcpRouteDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}
	d.SetId("")
	return nil
}
