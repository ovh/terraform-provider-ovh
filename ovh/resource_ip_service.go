package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIpService() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpServiceCreate,
		Update: resourceIpServiceUpdate,
		Read:   resourceIpServiceRead,
		Delete: resourceIpServiceDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: resourceIpServiceSchema(),
	}
}

func resourceIpServiceSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Custom description on your ip",
			Optional:    true,
			Computed:    true,
		},

		//computed
		"can_be_terminated": {
			Type:     schema.TypeBool,
			Computed: true,
		},

		"country": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"organisation_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"routed_to": {
			Type:        schema.TypeList,
			Description: "Routage information",
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_name": {
						Type:        schema.TypeString,
						Description: "Service where ip is routed to",
						Computed:    true,
					},
				},
			},
		},
		"service_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible values for ip type",
			Computed:    true,
		},
	}

	for k, v := range genericOrderSchema(true) {
		schema[k] = v
	}

	return schema
}

func resourceIpServiceCreate(d *schema.ResourceData, meta interface{}) error {
	if err := orderCreateFromResource(d, meta, "ip", true); err != nil {
		return fmt.Errorf("could not order ip: %q", err)
	}

	config := meta.(*Config)

	orderIdInt, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to convert orderID to int: %w", err)
	}

	serviceName, err := serviceNameFromOrder(config.OVHClient, int64(orderIdInt), d.Get("plan.0.plan_code").(string))
	if err != nil {
		return fmt.Errorf("failed to get service name from order ID: %w", err)
	}

	d.Set("service_name", serviceName)

	return resourceIpServiceUpdate(d, meta)
}

func resourceIpServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will update ip: %s", serviceName)
	opts := (&IpServiceUpdateOpts{}).FromResource(d)
	endpoint := fmt.Sprintf("/ip/service/%s",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("calling Put %s: %q", endpoint, err)
	}

	return resourceIpServiceRead(d, meta)
}

func resourceIpServiceRead(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("service_name").(string)

	config := meta.(*Config)
	log.Printf("[DEBUG] Will read ip: %s", serviceName)

	r := &IpService{}
	endpoint := fmt.Sprintf("/ip/service/%s",
		url.PathEscape(serviceName),
	)

	// This retry logic is there to handle a known API bug
	// which happens while an ipblock is attached/detached from
	// a Vrack
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := config.OVHClient.Get(endpoint, &r); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 400 {
				log.Printf("[DEBUG] known API bug when attaching/detaching vrack")
				return resource.RetryableError(err)
			}

			err = helpers.CheckDeleted(d, err, endpoint)
			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		}

		// Successful Get
		return nil
	})

	if err != nil {
		return err
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceIpServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	id := d.Id()
	serviceName := d.Get("service_name").(string)

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate ip %s for order %s", serviceName, id)
		endpoint := fmt.Sprintf(
			"/ip/service/%s/terminate",
			url.PathEscape(serviceName),
		)
		if err := config.OVHClient.Post(endpoint, nil, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
				return "", nil
			}
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return serviceName, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of ip %s for order %s", serviceName, id)
		endpoint := fmt.Sprintf(
			"/ip/service/%s/confirmTermination",
			url.PathEscape(serviceName),
		)
		if err := config.OVHClient.Post(endpoint, &IpServiceConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDeleteFromResource(d, meta, terminate, confirmTerminate); err != nil {
		return err
	}

	return nil
}
