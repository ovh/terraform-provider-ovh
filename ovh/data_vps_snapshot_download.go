package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceVPSSnapshotDownload exposes the signed download URL and size of
// the snapshot currently held by a VPS. The URL is short-lived and should be
// considered sensitive.
func dataSourceVPSSnapshotDownload() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPSSnapshotDownloadRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Short-lived signed URL to download the snapshot.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the snapshot in bytes.",
			},
		},
	}
}

func dataSourceVPSSnapshotDownloadRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	dl := &VPSSnapshotDownload{}
	endpoint := fmt.Sprintf("/vps/%s/snapshot/download", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, dl); err != nil {
		return diag.Errorf("calling Get %s: %s", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("url", dl.URL)
	d.Set("size", dl.Size)
	return nil
}
