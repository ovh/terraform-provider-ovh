package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIPLoadbalancingRouteHTTP() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingRouteHTTPCreate,
		Read:   resourceIPLoadbalancingRouteHTTPRead,
		Update: resourceIPLoadbalancingRouteHTTPUpdate,
		Delete: resourceIPLoadbalancingRouteHTTPDelete,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"target": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"frontend_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// IPLoadbalancingRouteHTTPAction Action triggered when all rules match
type IPLoadbalancingRouteHTTPAction struct {
	Target string `json:"target,omitempty"` // Farm ID for "farm" action type or URL template for "redirect" action. You may use ${uri}, ${protocol}, ${host}, ${port} and ${path} variables in redirect target
	Status int    `json:"status,omitempty"` // HTTP status code for "redirect" and "reject" actions
	Type   string `json:"type,omitempty"`   // Action to trigger if all the rules of this route matches
}

//IPLoadbalancingRouteHTTP HTTP Route
type IPLoadbalancingRouteHTTP struct {
	Status      string                          `json:"status,omitempty"`      //Route status. Routes in "ok" state are ready to operate
	Weight      int                             `json:"weight,omitempty"`      //Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action
	Action      *IPLoadbalancingRouteHTTPAction `json:"action,omitempty"`      //Action triggered when all rules match
	RouteID     int                             `json:"routeId,omitempty"`     //Id of your route
	DisplayName string                          `json:"displayName,omitempty"` //Human readable name for your route, this field is for you
	FrontendID  int                             `json:"frontendId,omitempty"`  //Route traffic for this frontend
}

func resourceIPLoadbalancingRouteHTTPCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	action := &IPLoadbalancingRouteHTTPAction{}
	actionSet := d.Get("action").(*schema.Set).List()[0].(map[string]interface{})

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

	err := config.OVHClient.Get(endpoint, &r)
	if err != nil {
		return CheckDeleted(d, err, endpoint)
	}

	d.Set("status", r.Status)
	d.Set("weight", r.Weight)
	d.Set("display_name", r.DisplayName)
	d.Set("frontend_id", r.FrontendID)

	return nil
}

func resourceIPLoadbalancingRouteHTTPUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", service, d.Id())

	action := &IPLoadbalancingRouteHTTPAction{}
	actionSet := d.Get("action").(*schema.Set).List()[0].(map[string]interface{})

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

	return nil
}
