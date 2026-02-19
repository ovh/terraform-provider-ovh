package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceCloudProjectFailoverIpAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectFailoverIpAttachCreate,
		Read:   resourceCloudProjectFailoverIpAttachRead,
		Delete: resourceCloudProjectFailoverIpAttachDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: resourceCloudProjectFailoverIpAttachSchema(),
	}
}

func resourceCloudProjectFailoverIpAttachSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_name": {
			Type:        schema.TypeString,
			Description: "The service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
			ForceNew:    true,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
		},

		"block": {
			Type:        schema.TypeString,
			Description: "IP block",
			Optional:    true,
			Computed:    true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				err := helpers.ValidateIp(v.(string))
				if err != nil {
					errors = append(errors, err)
				}
				return
			},
		},
		"continent_code": {
			Type:        schema.TypeString,
			Description: "Ip continent",
			Optional:    true,
			Computed:    true,
		},
		"geo_loc": {
			Type:        schema.TypeString,
			Description: "Ip location",
			Optional:    true,
			Computed:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Ip id",
			Computed:    true,
		},
		"ip": {
			Type:        schema.TypeString,
			Description: "Ip",
			Optional:    true,
			Computed:    true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				err := helpers.ValidateIp(v.(string))
				if err != nil {
					errors = append(errors, err)
				}
				return
			},
		},
		"progress": {
			Type:        schema.TypeInt,
			Description: "Current operation progress in percent",
			Computed:    true,
		},
		"routed_to": {
			Type:        schema.TypeString,
			Description: "Instance where ip is routed to",
			Computed:    true,
			ForceNew:    true,
			Optional:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "Ip status",
			Computed:    true,
		},
		"sub_type": {
			Type:        schema.TypeString,
			Description: "IP sub type",
			Computed:    true,
		},
	}
}

func resourceCloudProjectFailoverIpAttachRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read cloud project ip addresses %s", serviceName)
	endpoint := fmt.Sprintf("/cloud/project/%s/ip/failover",
		url.PathEscape(serviceName),
	)

	ips := []FailoverIp{}
	if err := config.OVHClient.Get(endpoint, &ips); err != nil {
		return fmt.Errorf("error calling GET %s:\n\t %q", endpoint, err)
	}

	match := false
	for _, ip := range ips {
		if ip.Ip == d.Get("ip").(string) {
			for k, v := range ip.ToMap() {
				match = true
				if k == "id" {
					d.SetId(v.(string))
				} else {
					err := d.Set(k, v)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	if !match {
		return fmt.Errorf("failover IP %s cannot be found in cloud project %s", d.Get("ip").(string), serviceName)
	}

	return nil
}

func resourceCloudProjectFailoverIpAttachCreate(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("service_name").(string)
	config := meta.(*Config)

	//Fetch Failover IP address to populate ID field
	log.Printf("[DEBUG] Will read cloud project ip addresses %s", serviceName)
	endpoint := fmt.Sprintf("/cloud/project/%s/ip/failover",
		url.PathEscape(serviceName),
	)

	ips := []FailoverIp{}
	if err := config.OVHClient.Get(endpoint, &ips); err != nil {
		return fmt.Errorf("error calling GET %s:\n\t %q", endpoint, err)
	}

	match := false
	for _, ip := range ips {
		if ip.Ip == d.Get("ip").(string) {
			for k, v := range ip.ToMap() {
				match = true
				if k == "id" {
					d.SetId(v.(string))
				}
			}
		}
	}

	if !match {
		return fmt.Errorf("failover IP %s cannot be found in cloud project %s", d.Get("ip").(string), serviceName)
	}

	id := d.Get("id").(string)

	log.Printf("[DEBUG] Will attach failover ip to an instance: %s", serviceName)
	opts := (&ProjectIpFailoverAttachCreation{}).FromResource(d)
	endpoint = fmt.Sprintf("/cloud/project/%s/ip/failover/%s/attach",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)

	err := retry.RetryContext(context.Background(), 5*time.Minute, func() *retry.RetryError {
		ip := &FailoverIp{}
		if err := config.OVHClient.Post(endpoint, opts, ip); err != nil {
			// Retry 400 errors because it can mean that the instance IP
			// is not allocated yet.
			ovhError, isOvhApiError := err.(*ovh.APIError)
			if isOvhApiError && ovhError.Code == 400 {
				log.Printf("[INFO] IP with id=%s not attached yet, retrying…", id)
				return retry.RetryableError(fmt.Errorf("error calling POST %s: %q", endpoint, err))
			} else {
				return retry.NonRetryableError(fmt.Errorf("failed to attach failover IP: %s", err))
			}
		}

		for k, v := range ip.ToMap() {
			if k != "id" {
				err := d.Set(k, v)
				if err != nil {
					return retry.NonRetryableError(err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	for d.Get("status").(string) == "operationPending" {
		if err := resourceCloudProjectFailoverIpAttachRead(d, meta); err != nil {
			return err
		}
	}

	return nil
}

func resourceCloudProjectFailoverIpAttachDelete(d *schema.ResourceData, meta interface{}) error {
	// Failover IPs cannot be detached from an instance, so nothing done here.
	d.SetId("")

	return nil
}
