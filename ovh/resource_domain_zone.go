package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDomainZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainZoneCreate,
		Read:   resourceDomainZoneRead,
		Delete: resourceDomainZoneDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: resourceDomainZoneSchema(),
	}
}

func resourceDomainZoneSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{

		// computed
		"dnssec_supported": {
			Type:        schema.TypeBool,
			Description: "Is DNSSEC supported by this zone",
			Computed:    true,
		},
		"has_dns_anycast": {
			Type:        schema.TypeBool,
			Description: "hasDnsAnycast flag of the DNS zone",
			Computed:    true,
		},
		"last_update": {
			Type:        schema.TypeString,
			Description: "Last update date of the DNS zone",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Zone name",
			Computed:    true,
		},
		"name_servers": {
			Type:        schema.TypeList,
			Description: "Name servers that host the DNS zone",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
		},
	}

	for k, v := range genericOrderSchema(false) {
		schema[k] = v
	}

	return schema
}

func resourceDomainZoneCreate(d *schema.ResourceData, meta interface{}) error {
	if err := orderCreate(d, meta, "dns"); err != nil {
		return fmt.Errorf("Could not order domain zone: %q", err)
	}

	return resourceDomainZoneRead(d, meta)
}

func resourceDomainZoneRead(d *schema.ResourceData, meta interface{}) error {
	_, details, err := orderRead(d, meta)
	if err != nil {
		return fmt.Errorf("Could not read domainZone order: %q", err)
	}

	config := meta.(*Config)
	zoneName := details[0].Domain

	log.Printf("[DEBUG] Will read domainZone: %s", zoneName)
	r := &DomainZone{}
	endpoint := fmt.Sprintf("/domain/zone/%s", zoneName)
	if err := config.OVHClient.Get(endpoint, &r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceDomainZoneDelete(d *schema.ResourceData, meta interface{}) error {
	_, details, err := orderRead(d, meta)
	if err != nil {
		return fmt.Errorf("Could not read domainZone order: %q", err)
	}

	config := meta.(*Config)
	zoneName := details[0].Domain

	id := d.Id()

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate domain zone %s for order %s", zoneName, id)
		endpoint := fmt.Sprintf(
			"/domain/zone/%s/terminate",
			url.PathEscape(zoneName),
		)
		if err := config.OVHClient.Post(endpoint, nil, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
				return "", nil
			}
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return zoneName, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of domain zone %s for order %s", zoneName, id)
		endpoint := fmt.Sprintf(
			"/domain/zone/%s/confirmTermination",
			url.PathEscape(zoneName),
		)
		if err := config.OVHClient.Post(endpoint, &DomainZoneConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDelete(d, meta, terminate, confirmTerminate); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
