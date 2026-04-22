package ovh

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

var (
	publicCloudProjectNameFormatRegex = regexp.MustCompile("^[0-9a-f]{12}4[0-9a-f]{19}$")
)

func resourceCloudProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectCreate,
		Update: resourceCloudProjectUpdate,
		Read:   resourceCloudProjectRead,
		Delete: resourceCloudProjectDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				d.Set("project_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: resourceCloudProjectSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultOrderTimeout),
		},
	}
}

func resourceCloudProjectSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"deletion_protection": {
			Type:        schema.TypeBool,
			Description: "Prevent the cloud project from being destroyed. Defaults to false.",
			Optional:    true,
			Default:     false,
		},

		// computed
		"urn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"access": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range genericOrderSchema(false) {
		schema[k] = v
	}

	return schema
}

func resourceCloudProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := orderCreateFromResource(d, meta, "cloud", true, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("could not order cloud project: %q", err)
	}

	order, details, err := orderReadInResource(d, meta)
	if err != nil {
		return fmt.Errorf("could not read cloud project order: %q", err)
	}

	serviceName, err := resourceCloudProjectGetServiceName(config, order, details)
	if err != nil {
		return err
	}

	d.SetId(serviceName)
	d.Set("project_id", serviceName)

	// After order delivery, the project may not be immediately available in the API.
	// Retry the GET for up to 10 minutes to avoid a spurious 404 error.
	endpoint := fmt.Sprintf("/cloud/project/%s", url.PathEscape(serviceName))
	if err := helpers.WaitAvailable(config.OVHClient, endpoint, 10*time.Minute); err != nil {
		return fmt.Errorf("waiting for cloud project %s to become available: %q", serviceName, err)
	}

	return resourceCloudProjectUpdate(d, meta)
}

func resourceCloudProjectGetServiceName(config *Config, order *MeOrder, details []*MeOrderDetail) (string, error) {
	// Looking for an order detail associated to a Public Cloud Project.
	// Cloud Project has a specific resource_name that we can grep through a Regexp
	for _, d := range details {
		domain := d.Domain
		if publicCloudProjectNameFormatRegex.MatchString(domain) {
			return domain, nil
		}
	}

	// For OVHcloud US, resource_name are not stored inside order detail, but inside the operation associated to the order detail.
	for _, orderDetail := range details {
		operations, err := orderDetailOperations(config.OVHClient, order.OrderId, orderDetail.OrderDetailId)
		if err != nil {
			return "", fmt.Errorf("could not read cloudProject order details operations: %q", err)
		}
		for _, operation := range operations {
			if publicCloudProjectNameFormatRegex.MatchString(operation.Resource.Name) {
				return operation.Resource.Name, nil
			}
		}
	}

	return "", fmt.Errorf("unknown service name")
}

func resourceCloudProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("project_id").(string)

	log.Printf("[DEBUG] Will update cloudProject: %s", serviceName)
	opts := (&CloudProjectUpdateOpts{}).FromResource(d)
	endpoint := fmt.Sprintf("/cloud/project/%s", url.PathEscape(serviceName))
	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("calling Put %s: %q", endpoint, err)
	}

	return resourceCloudProjectRead(d, meta)
}

func resourceCloudProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("project_id").(string)

	log.Printf("[DEBUG] Will read cloudProject: %s", serviceName)
	r := &CloudProject{}
	endpoint := fmt.Sprintf("/cloud/project/%s", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, &r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceCloudProjectDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("cloud project %s is protected from deletion. Set deletion_protection to false before destroying", d.Id())
	}

	config := meta.(*Config)
	serviceName := d.Get("project_id").(string)

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate cloud project %s", serviceName)
		endpoint := fmt.Sprintf(
			"/cloud/project/%s/terminate",
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
		log.Printf("[DEBUG] Will confirm termination of cloud project %s", serviceName)
		endpoint := fmt.Sprintf(
			"/cloud/project/%s/confirmTermination",
			url.PathEscape(serviceName),
		)
		if err := config.OVHClient.Post(endpoint, &CloudProjectConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDeleteFromResource(d, meta, terminate, confirmTerminate); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
