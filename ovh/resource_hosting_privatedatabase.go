package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceHostingPrivateDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostingPrivateDatabaseCreate,
		Update: resourceHostingPrivateDatabaseUpdate,
		Read:   resourceHostingPrivateDatabaseRead,
		Delete: resourceHostingPrivateDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				d.Set("service_name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: resourceHostingPrivateDatabaseSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultOrderTimeout),
		},
	}
}

func resourceHostingPrivateDatabaseSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"service_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		// Computed
		"advanced_configuration": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Advanced configuration key / value",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"urn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of CPU on your private database",
		},
		"datacenter": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Datacenter where this private database is located",
		},
		"display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Name displayed in customer panel for your private database",
		},
		"hostname": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private database hostname",
		},
		"hostname_ftp": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private database FTP hostname",
		},
		"infrastructure": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Infrastructure where service was stored",
		},
		"offer": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of the private database offer",
		},
		"port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Private database service port",
		},
		"port_ftp": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Private database FTP port",
		},
		"quota_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Space allowed (in MB) on your private database",
		},
		"quota_used": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Sapce used (in MB) on your private database",
		},
		"ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Amount of ram (in MB) on your private database",
		},
		"server": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private database server name",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private database state",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private database type",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private database available versions",
		},
		"version_label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private database version label",
		},
		"version_number": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "Private database version number",
		},
	}

	for k, v := range genericOrderSchema(false) {
		schema[k] = v
	}

	return schema
}

func resourceHostingPrivateDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := orderCreateFromResource(d, meta, "privateSQL", true, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("could not order privateDatabase: %q", err)
	}

	orderIdInt, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to convert orderID to int: %w", err)
	}

	serviceName, err := serviceNameFromOrder(config.OVHClient, int64(orderIdInt), d.Get("plan.0.plan_code").(string))
	if err != nil {
		return fmt.Errorf("could not retrieve service name from order: %w", err)
	}

	d.SetId(serviceName)
	d.Set("service_name", serviceName)

	return resourceHostingPrivateDatabaseUpdate(d, meta)
}

func resourceHostingPrivateDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	if d.HasChange("display_name") {
		log.Printf("[DEBUG] Will update privateDatabase: %s", serviceName)
		opts := (&HostingPrivateDatabaseOpts{}).FromResource(d)
		endpoint := fmt.Sprintf("/hosting/privateDatabase/%s", url.PathEscape(serviceName))
		if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
			return fmt.Errorf("calling Put %s: %q", endpoint, err)
		}
	}

	if d.HasChange("advanced_configuration") {
		log.Printf("[DEBUG] Will update privateDatabase config: %s", serviceName)
		configEndpoint := fmt.Sprintf("/hosting/privateDatabase/%s/config/update", url.PathEscape(serviceName))
		acParams := d.Get("advanced_configuration").(map[string]interface{})
		parameters := make([]map[string]interface{}, 0, len(acParams))
		for k, v := range acParams {
			parameters = append(parameters, map[string]interface{}{
				"key":   k,
				"value": v,
			})
		}
		payload := map[string]interface{}{
			"parameters": parameters,
		}
		if err := config.OVHClient.Post(configEndpoint, payload, nil); err != nil {
			return fmt.Errorf("calling Post %s with params %+v:\n\t %q", configEndpoint, payload, err)
		}
	}

	return resourceHostingPrivateDatabaseRead(d, meta)

}

func resourceHostingPrivateDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read privateDatabase: %s", serviceName)
	ds := &HostingPrivateDatabase{}
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s", url.PathEscape(serviceName))
	err := config.OVHClient.Get(endpoint, &ds)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	configEndpoint := fmt.Sprintf("/hosting/privateDatabase/%s/config", url.PathEscape(serviceName))
	var rawCfg map[string]interface{}
	if err := config.OVHClient.Get(configEndpoint, &rawCfg); err != nil {
		return fmt.Errorf("calling Get %s: %q", configEndpoint, err)
	}

	// The OVH API returns the advanced configuration as an object containing a "details"
	// array (each item includes "key" and "value", plus metadata). The Terraform schema
	// for advanced_configuration is a simple key/value map, so we flatten it here.
	advancedCfg := make(map[string]interface{})

	if detailsRaw, ok := rawCfg["details"]; ok {
		if details, ok := detailsRaw.([]interface{}); ok {
			for _, item := range details {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				key, _ := m["key"].(string)
				if key == "" {
					continue
				}
				if v, exists := m["value"]; exists && v != nil {
					advancedCfg[key] = fmt.Sprint(v)
				}
			}
		}
	} else {
		// Backward/alternate API shapes: keep any top-level keys as strings.
		for k, v := range rawCfg {
			if v != nil {
				advancedCfg[k] = fmt.Sprint(v)
			}
		}
	}

	if err := d.Set("advanced_configuration", advancedCfg); err != nil {
		return fmt.Errorf("setting advanced_configuration from %s: %w", configEndpoint, err)
	}

	return nil
}

func resourceHostingPrivateDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate privateDatabase %s", serviceName)
		endpoint := fmt.Sprintf(
			"/hosting/privateDatabase/%s/terminate",
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
		log.Printf("[DEBUG] Will confirm termination of privateDatabase %s", serviceName)
		endpoint := fmt.Sprintf(
			"/hosting/privateDatabase/%s/confirmTermination",
			url.PathEscape(serviceName),
		)
		if err := config.OVHClient.Post(endpoint, &HostingPrivateDatabaseConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDeleteFromResource(d, meta, terminate, confirmTerminate); err != nil {
		return err
	}

	return nil
}
