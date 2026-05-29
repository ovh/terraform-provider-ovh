package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPSConsoleURL() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSConsoleURLCreate,
		Read:   resourceVPSConsoleURLRead,
		Delete: resourceVPSConsoleURLDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service_name of your VPS.",
			},
			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Arbitrary map of values. Changing any value forces a new resource, " +
					"which re-issues a fresh console URL via POST /vps/{serviceName}/getConsoleUrl.",
			},

			// Computed
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "A fresh, single-use signed console URL. Issued on create and expires quickly server-side; do not persist long-term.",
			},
		},
	}
}

func resourceVPSConsoleURLCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/getConsoleUrl", url.PathEscape(serviceName))

	log.Printf("[DEBUG] Requesting fresh VPS console URL for %s", serviceName)

	var consoleURL string
	if err := config.OVHClient.Post(endpoint, nil, &consoleURL); err != nil {
		return fmt.Errorf("calling POST %s: %w", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("url", consoleURL)

	return nil
}

func resourceVPSConsoleURLRead(d *schema.ResourceData, meta interface{}) error {
	// The OVH API exposes no GET for a previously issued console URL: the POST
	// returns a one-shot signed URL intended for immediate use. We keep whatever
	// the create call stored in state and treat the URL as opaque. Any change to
	// service_name or triggers forces a new resource which re-issues a fresh URL.
	return nil
}

func resourceVPSConsoleURLDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no revoke endpoint; the URL auto-expires server-side.
	// Removing the resource from state is sufficient.
	d.SetId("")
	return nil
}
