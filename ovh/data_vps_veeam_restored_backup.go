package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSVeeamRestoredBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSVeeamRestoredBackupRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS",
			},
			"restore_point_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the currently mounted Veeam restore point",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the restored backup (mounted|restoring|unmounted|unmounting)",
			},
			"access_nfs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NFS access information",
			},
			"access_smb": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMB access information",
			},
		},
	}
}

func dataSourceVPSVeeamRestoredBackupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/veeam/restoredBackup", url.PathEscape(serviceName))
	rb := &VpsVeeamRestoredBackup{}
	if err := config.OVHClient.Get(endpoint, rb); err != nil {
		return fmt.Errorf("error calling GET %s: %s", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s/%d", serviceName, rb.RestorePointId))
	d.Set("restore_point_id", rb.RestorePointId)
	d.Set("state", rb.State)
	d.Set("access_nfs", rb.AccessInfos.Nfs)
	d.Set("access_smb", rb.AccessInfos.Smb)
	return nil
}
