package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSAutomatedBackupAttached() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSAutomatedBackupAttachedRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS",
			},
			"attached_backups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Currently attached automated backup restore points",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restore_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nfs": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"smb": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"additional_disk": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVPSAutomatedBackupAttachedRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/automatedBackup/attachedBackup", url.PathEscape(serviceName))
	attached := []VPSAutomatedBackupAttached{}
	if err := config.OVHClient.Get(endpoint, &attached); err != nil {
		return fmt.Errorf("error calling GET %s: %w", endpoint, err)
	}

	out := make([]map[string]interface{}, 0, len(attached))
	keys := []string{serviceName}
	for _, a := range attached {
		out = append(out, map[string]interface{}{
			"restore_point":   a.RestorePoint,
			"nfs":             a.Access.NFS,
			"smb":             a.Access.SMB,
			"additional_disk": a.Access.AdditionalDisk,
		})
		keys = append(keys, a.RestorePoint)
	}

	d.SetId(hashcode.Strings(keys))
	d.Set("attached_backups", out)
	return nil
}
