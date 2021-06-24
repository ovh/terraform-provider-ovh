package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceDedicatedCeph() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedCephRead,
		Schema: map[string]*schema.Schema{
			"ceph_mons": {
				Type:     schema.TypeList,
				Optional: false,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ceph_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"crush_tunables": {
				Type:     schema.TypeString,
				Optional: false,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: false,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: false,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: false,
				Computed: false,
				Required: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Optional: false,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: false,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateDedicatedCephStatus(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}
func dataSourceDedicatedCephRead(d *schema.ResourceData, meta interface{}) error {
	url := "/dedicated/ceph"
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	log.Printf("[DEBUG] Will retrieve dedicated CEPH %s", serviceName)

	ceph := &DedicatedCeph{}
	err := config.OVHClient.Get(fmt.Sprintf("%s/%s", url, serviceName), &ceph)
	if err != nil {
		return fmt.Errorf("Error calling %s/%s:\n\t %q", url, serviceName, err)
	}
	log.Printf("[DEBUG] CEPH is %v", ceph.CephMonitors)
	d.SetId(ceph.ServiceName)
	d.Set("service_name", ceph.ServiceName)
	d.Set("ceph_mons", ceph.CephMonitors)
	d.Set("ceph_version", ceph.CephVersion)
	d.Set("crush_tunables", ceph.CrushTunables)
	d.Set("label", ceph.Label)
	d.Set("region", ceph.Region)
	d.Set("size", ceph.Size)
	d.Set("state", ceph.State)
	d.Set("status", ceph.Status)

	return nil
}
