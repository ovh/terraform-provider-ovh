package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSBackupFtp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSBackupFtpRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS.",
			},
			"ftp_backup_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the backup FTP server.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backup FTP offer type.",
			},
			"read_only_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when backup FTP will be set in read-only mode.",
			},
			"quota": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Backup FTP storage quota (unit/value).",
			},
			"usage": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Backup FTP storage usage (unit/value).",
			},
		},
	}
}

func dataSourceVPSBackupFtpRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	backup := &VPSBackupFtp{}
	endpoint := fmt.Sprintf("/vps/%s/backupftp", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, backup); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("ftp_backup_name", backup.FtpBackupName)
	d.Set("type", backup.Type)
	if backup.ReadOnlyDate != nil {
		d.Set("read_only_date", *backup.ReadOnlyDate)
	} else {
		d.Set("read_only_date", "")
	}

	quota := map[string]string{
		"unit":  backup.Quota.Unit,
		"value": fmt.Sprintf("%v", backup.Quota.Value),
	}
	d.Set("quota", quota)

	usage := map[string]string{
		"unit":  backup.Usage.Unit,
		"value": fmt.Sprintf("%v", backup.Usage.Value),
	}
	d.Set("usage", usage)

	return nil
}
