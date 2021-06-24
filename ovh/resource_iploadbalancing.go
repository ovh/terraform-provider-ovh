package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func resourceIpLoadbalancing() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpLoadbalancingCreate,
		Update: resourceIpLoadbalancingUpdate,
		Read:   resourceIpLoadbalancingRead,
		Delete: resourceIpLoadbalancingDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: resourceIpLoadbalancingSchema(),
	}
}

func resourceIpLoadbalancingSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Description: "Set the name displayed in ManagerV6 for your iplb (max 50 chars)",
			Optional:    true,
			Computed:    true,
		},
		"ssl_configuration": {
			Type:        schema.TypeString,
			Description: "Modern oldest compatible clients : Firefox 27, Chrome 30, IE 11 on Windows 7, Edge, Opera 17, Safari 9, Android 5.0, and Java 8. Intermediate oldest compatible clients : Firefox 1, Chrome 1, IE 7, Opera 5, Safari 1, Windows XP IE8, Android 2.3, Java 7. Intermediate if null.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				err := helpers.ValidateStringEnum(strings.ToLower(v.(string)), []string{
					"intermediate",
					"modern",
				})
				if err != nil {
					errors = append(errors, err)
				}
				return
			},
		},

		//computed
		"ipv6": {
			Type:        schema.TypeString,
			Description: "The IPV6 associated to your IP load balancing. DEPRECATED.",
			Computed:    true,
		},
		"ipv4": {
			Type:        schema.TypeString,
			Description: "The IPV4 associated to your IP load balancing",
			Computed:    true,
		},
		"zone": {
			Type:        schema.TypeList,
			Description: "Location where your service is",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
		},
		"service_name": {
			Type:        schema.TypeString,
			Description: "The internal name of your IP load balancing",
			Computed:    true,
		},
		"ip_loadbalancing": {
			Type:        schema.TypeString,
			Description: "Your IP load balancing",
			Computed:    true,
		},
		"state": {
			Type:        schema.TypeString,
			Description: "Current state of your IP",
			Computed:    true,
		},
		"offer": {
			Type:        schema.TypeString,
			Description: "The offer of your IP load balancing",
			Computed:    true,
		},

		"vrack_eligibility": {
			Type:        schema.TypeBool,
			Description: "Vrack eligibility",
			Computed:    true,
		},
		"vrack_name": {
			Type:        schema.TypeString,
			Description: "Name of the vRack on which the current Load Balancer is attached to, as it is named on vRack product",
			Computed:    true,
		},

		// additional exported attributes
		"metrics_token": {
			Type:        schema.TypeString,
			Description: "The metrics token associated with your IP load balancing",
			Sensitive:   true,
			Computed:    true,
		},
		"orderable_zone": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Available additional zone for your Load Balancer",
			Set: func(v interface{}) int {
				r := v.(map[string]interface{})
				return hashcode.String(r["name"].(string))
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The zone three letter code",
						Computed:    true,
					},

					"plan_code": {
						Type:        schema.TypeString,
						Description: "The billing planCode for this zone",
						Computed:    true,
					},
				},
			},
		},
	}

	for k, v := range genericOrderSchema(true) {
		schema[k] = v
	}

	return schema
}

func resourceIpLoadbalancingCreate(d *schema.ResourceData, meta interface{}) error {
	if err := orderCreate(d, meta, "ipLoadbalancing"); err != nil {
		return fmt.Errorf("Could not order ipLoadbalancing: %q", err)
	}

	return resourceIpLoadbalancingUpdate(d, meta)
}

func resourceIpLoadbalancingUpdate(d *schema.ResourceData, meta interface{}) error {
	_, details, err := orderRead(d, meta)
	if err != nil {
		return fmt.Errorf("Could not read ipLoadbalancing order: %q", err)
	}

	var serviceName string
	if strings.Contains(details[0].Domain, "-zone-") {
		serviceName = strings.Split(details[0].Domain, "-zone-")[0]
	} else {
		serviceName = details[0].Domain
	}

	config := meta.(*Config)

	log.Printf("[DEBUG] Will update ipLoadbalancing: %s", serviceName)
	opts := (&IpLoadbalancingUpdateOpts{}).FromResource(d)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s", serviceName)
	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("calling Put %s: %q", endpoint, err)
	}

	return resourceIpLoadbalancingRead(d, meta)
}

func resourceIpLoadbalancingRead(d *schema.ResourceData, meta interface{}) error {
	_, details, err := orderRead(d, meta)
	if err != nil {
		return fmt.Errorf("Could not read ipLoadbalancing order: %q", err)
	}

	var serviceName string
	if strings.Contains(details[0].Domain, "-zone-") {
		serviceName = strings.Split(details[0].Domain, "-zone-")[0]
	} else {
		serviceName = details[0].Domain
	}

	config := meta.(*Config)
	log.Printf("[DEBUG] Will read ipLoadbalancing: %s", serviceName)

	r := &IpLoadbalancing{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s", serviceName)
	if err := config.OVHClient.Get(endpoint, &r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceIpLoadbalancingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()
	serviceName := d.Get("service_name").(string)

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate ipLoadbalancing %s for order %s", serviceName, id)
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/terminate",
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
		log.Printf("[DEBUG] Will confirm termination of ipLoadbalancing %s for order %s", serviceName, id)
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/confirmTermination",
			url.PathEscape(serviceName),
		)
		if err := config.OVHClient.Post(endpoint, &IpLoadbalancingConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDelete(d, meta, terminate, confirmTerminate); err != nil {
		return err
	}

	return nil
}
