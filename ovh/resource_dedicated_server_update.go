package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDedicatedServerUpdate() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedServerUpdateCreateOrUpdate,
		Update: resourceDedicatedServerUpdateCreateOrUpdate,
		Read:   resourceDedicatedServerUpdateRead,
		Delete: resourceDedicatedServerUpdateDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your dedicated server.",
				Required:    true,
			},

			"boot_id": {
				Type:        schema.TypeInt,
				Description: "The boot id of your dedicated server.",
				Computed:    true,
				Optional:    true,
			},
			"monitoring": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Icmp monitoring state",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "error, hacked, hackedBlocked, ok",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(v.(string), []string{"error", "hacked", "hackedBlocked", "ok"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func resourceDedicatedServerUpdateCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	opts := (&DedicatedServerUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s:\n\t %q", endpoint, err)
	}

	//set fake id
	d.SetId(serviceName)

	return resourceDedicatedServerUpdateRead(d, meta)
}

func resourceDedicatedServerUpdateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	ds := &DedicatedServer{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s",
			url.PathEscape(serviceName),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"Error calling GET /dedicated/server/%s:\n\t %q",
			serviceName,
			err,
		)
	}

	d.Set("boot_id", ds.BootId)
	d.Set("monitoring", ds.Monitoring)
	d.Set("state", ds.State)

	//set fake id
	d.SetId(serviceName)
	return nil
}

func resourceDedicatedServerUpdateDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
