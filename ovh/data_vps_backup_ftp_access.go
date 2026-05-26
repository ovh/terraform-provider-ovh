package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSBackupFtpAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSBackupFtpAccessRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS.",
			},
			"ip_block": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CIDR-formatted IP block identifying the ACL entry.",
			},

			// Computed
			"ftp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether FTP access is granted to this IP block.",
			},
			"cifs": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether CIFS / SMB access is granted to this IP block.",
			},
			"nfs": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether NFS access is granted to this IP block.",
			},
			"is_applied": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the ACL entry has been applied on the backup FTP server.",
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the last ACL update.",
			},
		},
	}
}

func dataSourceVPSBackupFtpAccessRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipBlock := d.Get("ip_block").(string)

	endpoint := fmt.Sprintf(
		"/vps/%s/backupftp/access/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipBlock),
	)

	acl := &VPSBackupFtpAcl{}
	if err := config.OVHClient.Get(endpoint, acl); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s|%s", serviceName, ipBlock))
	d.Set("ftp", acl.Ftp)
	d.Set("cifs", acl.Cifs)
	d.Set("nfs", acl.Nfs)
	d.Set("is_applied", acl.IsApplied)
	d.Set("last_update", acl.LastUpdate)

	return nil
}
