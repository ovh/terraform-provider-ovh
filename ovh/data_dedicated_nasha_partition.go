package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedNashaPartition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedNashaPartitionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_by_snapshots": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDedicatedNashaPartitionRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	partitionName := d.Get("name").(string)

	ds := &DedicatedNASHAPartition{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/nasha/%s/partition/%s",
			url.PathEscape(serviceName),
			url.PathEscape(partitionName),
		),
		&ds,
	)

	if err != nil {
		return diag.Errorf(
			"Error calling /dedicated/nasha/%s/partition/%s:\n\t %q",
			serviceName,
			partitionName,
			err,
		)
	}

	d.SetId(ds.Name)
	d.Set("name", ds.Name)
	d.Set("description", ds.Description)
	d.Set("protocol", ds.Protocol)
	d.Set("size", ds.Size)
	d.Set("capacity", ds.Capacity)
	d.Set("used_by_snapshots", ds.UsedBySnapshots)

	return nil
}
