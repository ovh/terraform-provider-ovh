package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPSSnapshotRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS.",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	snap := &VPSSnapshot{}
	endpoint := fmt.Sprintf("/vps/%s/snapshot", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, snap); err != nil {
		return diag.Errorf("calling Get %s: %s", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("description", snap.Description)
	if !snap.CreationDate.IsZero() {
		d.Set("creation_date", snap.CreationDate.Format("2006-01-02T15:04:05Z07:00"))
	}
	d.Set("region", snap.Region)
	return nil
}
