package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDedicatedServerOlaAggregation() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedServerOlaAggregationCreate,
		Read:   resourceDedicatedServerOlaAggregationRead,
		Delete: resourceDedicatedServerOlaAggregationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your dedicated server.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resulting VirtualNetworkInterface name.",
			},
			"virtual_network_interfaces": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of interfaces to aggregate.",
			},

			//Computed
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Details of this task",
			},
			"done_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Completion date",
			},
			"function": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Function name",
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update",
			},
			"start_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task Creation date",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task status",
			},
			"enabled_vrack_aggregation_vni": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vrack_aggregation VNI uuid",
			},
			"enabled_public_vni": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "public VNI uuid",
			},
		},
	}
}

func resourceDedicatedServerOlaAggregationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/ola/aggregation",
		url.PathEscape(serviceName),
	)
	opts := (&DedicatedServerOlaAggregationCreateOpts{}).FromResource(d)
	task := &DedicatedServerTask{}

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	if err := waitForDedicatedServerTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", task.Id))

	return dedicatedServerOlaAggregationRead(d, meta)
}

func dedicatedServerOlaAggregationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf(
			"Could not parse install task id %s,%s:\n\t %q",
			serviceName,
			d.Id(),
			err,
		)
	}

	task, err := getDedicatedServerTask(serviceName, id, config.OVHClient)
	if err != nil {
		return helpers.CheckDeleted(d, err, fmt.Sprintf(
			"dedicated server task %s/%s",
			serviceName,
			d.Id(),
		))
	}

	// Set VNIs attributes
	vnis, err := getDedicatedServerVNIs(d, meta)

	if err != nil {
		return fmt.Errorf("Error reading Dedicated Server VNIs: %s", err)
	}

	mapvnis := make([]map[string]interface{}, len(vnis))
	enabledVrackAggregationVnis := []string{}
	enabledPublicVnis := []string{}

	for i, vni := range vnis {
		mapvnis[i] = vni.ToMap()

		if vni.Enabled {
			switch vni.Mode {
			case "vrack_aggregation":
				enabledVrackAggregationVnis = append(enabledVrackAggregationVnis, vni.Uuid)
			case "public":
				enabledPublicVnis = append(enabledPublicVnis, vni.Uuid)
			default:
				log.Printf("[WARN] unknown VNI mode. DS {%v} VNI {%v}", serviceName, vni)
			}
		}
	}

	if len(enabledVrackAggregationVnis) > 0 {
		d.Set("enabled_vrack_aggregation_vni", enabledVrackAggregationVnis[0])
	} else {
		d.Set("enabled_vrack_aggregation_vni", []string{})
	}

	if len(enabledPublicVnis) > 0 {
		d.Set("enabled_public_vni", enabledPublicVnis[0])
	} else {
		d.Set("enabled_public_vni", []string{})
	}

	d.Set("comment", task.Comment)
	d.Set("function", task.Function)
	d.Set("status", task.Status)
	d.Set("last_update", task.LastUpdate.Format(time.RFC3339))
	d.Set("done_date", task.DoneDate.Format(time.RFC3339))
	d.Set("start_date", task.StartDate.Format(time.RFC3339))

	return nil
}

func resourceDedicatedServerOlaAggregationRead(d *schema.ResourceData, meta interface{}) error {
	// Nothing to do on READ
	//
	// IMPORTANT: This resource doesn't represent a real resource
	// but instead a task on a dedicated server. OVH may clean its tasks database after a while
	// so that the API may return a 404 on a task id. If we hit a 404 on a READ, then
	// terraform will understand that it has to recreate the resource, and consequently
	// will trigger new install task on the dedicated server.
	// This is something we must avoid!
	//
	return nil
}

func resourceDedicatedServerOlaAggregationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/ola/reset",
		url.PathEscape(serviceName),
	)
	// opts := (&DedicatedServerOlaAggregationDeleteOpts{}).FromResource(d)
	task := &DedicatedServerTask{}

	// The uuid of the aggregated interfaces is not returned in the API, as such we need to retrieve it on our own

	vnis := []string{
		d.Get("enabled_vrack_aggregation_vni").(string),
		d.Get("enabled_public_vni").(string),
	}
	for _, v := range vnis {
		if v == "" {
			continue
		}
		singleOpts := DedicatedServerOlaAggregationSingleDeleteOpts{
			VirtualNetworkInterface: v,
		}
		if err := config.OVHClient.Post(endpoint, singleOpts, task); err != nil {
			return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
		}

		if err := waitForDedicatedServerTask(serviceName, task, config.OVHClient); err != nil {
			return err
		}
	}

	// we cant delete the task through the API, just forget about its Id
	d.SetId("")
	return nil
}
