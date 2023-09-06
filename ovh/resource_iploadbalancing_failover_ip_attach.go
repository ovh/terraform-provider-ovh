package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIpLoadbalancingFailoverIpAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpLoadbalancingFailoverIpAttachCreate,
		ReadContext:   resourceIpLoadbalancingFailoverIpAttachRead,
		DeleteContext: resourceIpLoadbalancingFailoverIpAttachDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: resourceIpLoadbalancingFailoverIpAttachSchema(),
	}
}

func resourceIpLoadbalancingFailoverIpAttachSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"ip": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				err := helpers.ValidateIpBlock(v.(string))
				if err != nil {
					errors = append(errors, err)
				}
				return
			},
		},
		"nexthop": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "",
		},
		"to": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func resourceIpLoadbalancingFailoverIpAttachRead(tx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/failover", url.PathEscape(serviceName))

	ipBlocks := []string{}
	if err := config.OVHClient.Get(endpoint, &ipBlocks); err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	match := false
	for _, ip := range ipBlocks {
		if ip == d.Get("ip").(string) {
			match = true
			d.SetId(serviceName)
		}
	}

	if !match {
		return diag.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return nil
}

func resourceIpLoadbalancingFailoverIpAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	ip := d.Get("ip").(string)
	payload := (&IpMoveCreateOpts{}).FromResource(d)
	ipTask := IpTask{}

	err := config.OVHClient.Post(
		fmt.Sprintf("/ip/%s/move", url.PathEscape(ip)), payload, &ipTask,
	)
	if err != nil {
		return diag.Errorf("Failed to move IP: %s", err)
	}

	err = waitForIpLoadbalancingFailoverIpAttachDone(ctx, config.OVHClient, ip, ipTask.TaskID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting ipTask %d to be DONE: %s", ipTask.TaskID, err.Error())
	}

	return resourceIpLoadbalancingFailoverIpAttachRead(ctx, d, meta)
}

func resourceIpLoadbalancingFailoverIpAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	ip := d.Get("ip").(string)
	ipTask := IpTask{}

	err := config.OVHClient.Post(
		fmt.Sprintf("/ip/%s/park", url.PathEscape(ip)), nil, &ipTask,
	)
	if err != nil {
		return diag.Errorf("Failed to park IP: %s", err)
	}

	err = waitForIpLoadbalancingFailoverIpAttachDone(ctx, config.OVHClient, ip, ipTask.TaskID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting ipTask %d to be DONE: %s", ipTask.TaskID, err.Error())
	}

	return nil
}
