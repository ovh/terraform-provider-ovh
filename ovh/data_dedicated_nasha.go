package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedNasha() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedNashaRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The storage service name",
			},

			// Computed
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"can_create_partition": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True, if partition creation is allowed on this HA-NAS",
			},
			"custom_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name you give to the HA-NAS",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "area of HA-NAS",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the disk type of the HA-NAS",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access IP of HA-NAS",
			},
			"monitored": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Send an email to customer if any issue is detected",
			},
			"zpool_capacity": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "percentage of HA-NAS space used in %",
			},
			"zpool_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the size of the HA-NAS",
			},

			"partitions_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of partition names for this HA-NAS",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceDedicatedNashaRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	ds := &DedicatedNASHA{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/nasha/%s",
			url.PathEscape(serviceName),
		),
		&ds,
	)

	if err != nil {
		return diag.Errorf(
			"Error calling /dedicated/nasha/%s:\n\t %q",
			serviceName,
			err,
		)
	}

	var partitionsList []string
	err = config.OVHClient.Get(
		fmt.Sprintf("/dedicated/nasha/%s/partition", url.PathEscape(serviceName)),
		&partitionsList,
	)
	if err != nil {
		return diag.Errorf(
			"Error calling /dedicated/nasha/%s/partition: %s",
			serviceName,
			err,
		)
	}

	d.SetId(ds.ServiceName)
	d.Set("urn", ds.URN)
	d.Set("service_name", ds.ServiceName)
	d.Set("monitored", ds.Monitored)
	d.Set("zpool_size", ds.ZpoolSize)
	d.Set("zpool_capacity", ds.ZpoolCapacity)
	d.Set("partitions_list", partitionsList)
	d.Set("datacenter", ds.Datacenter)
	d.Set("disk_type", ds.DiskType)
	d.Set("can_create_partition", ds.CanCreatePartition)
	d.Set("ip", ds.Ip)

	return nil
}
