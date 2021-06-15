package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIPLoadbalancingRouteHTTP() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingRouteHTTPCreate,
		Read:   resourceIPLoadbalancingRouteHTTPRead,
		Update: resourceIPLoadbalancingRouteHTTPUpdate,
		Delete: resourceIPLoadbalancingRouteHTTPDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpLoadbalancingHttpRouteImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"target": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frontend_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceIpLoadbalancingHttpRouteImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceIPLoadbalancingRouteHTTPCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	action := &IPLoadbalancingRouteHTTPAction{}
	actionSet := d.Get("action").([]interface{})[0].(map[string]interface{})

	action.Status = actionSet["status"].(int)
	action.Target = actionSet["target"].(string)
	action.Type = actionSet["type"].(string)

	route := &IPLoadbalancingRouteHTTP{
		Action:      action,
		DisplayName: d.Get("display_name").(string),
		FrontendID:  d.Get("frontend_id").(int),
		Weight:      d.Get("weight").(int),
	}

	service := d.Get("service_name").(string)
	resp := &IPLoadbalancingRouteHTTP{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route", service)

	err := config.OVHClient.Post(endpoint, route, resp)
	if err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", resp.RouteID))

	return resourceIPLoadbalancingRouteHTTPRead(d, meta)
}

func resourceIPLoadbalancingRouteHTTPRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	r := &IPLoadbalancingRouteHTTP{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", service, d.Id())

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(fmt.Sprintf("%d", r.RouteID))

	actions := make([]map[string]interface{}, 0)
	action := make(map[string]interface{})
	action["status"] = r.Action.Status
	action["target"] = r.Action.Target
	action["type"] = r.Action.Type
	actions = append(actions, action)

	d.Set("weight", r.Weight)
	d.Set("display_name", r.DisplayName)
	d.Set("frontend_id", r.FrontendID)
	d.Set("action", actions)

	return nil
}

func resourceIPLoadbalancingRouteHTTPUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", service, d.Id())

	action := &IPLoadbalancingRouteHTTPAction{}
	actionSet := d.Get("action").([]interface{})[0].(map[string]interface{})

	action.Status = actionSet["status"].(int)
	action.Target = actionSet["target"].(string)
	action.Type = actionSet["type"].(string)

	route := &IPLoadbalancingRouteHTTP{
		Action:      action,
		DisplayName: d.Get("display_name").(string),
		FrontendID:  d.Get("frontend_id").(int),
		Weight:      d.Get("weight").(int),
	}

	err := config.OVHClient.Put(endpoint, route, nil)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}

	return resourceIPLoadbalancingRouteHTTPRead(d, meta)
}

func resourceIPLoadbalancingRouteHTTPDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	r := &IPLoadbalancingRouteHTTP{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", service, d.Id())

	err := config.OVHClient.Delete(endpoint, &r)
	if err != nil {
		return fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	d.SetId("")
	return nil
}
