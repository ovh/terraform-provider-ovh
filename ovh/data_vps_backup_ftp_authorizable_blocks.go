package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSBackupFtpAuthorizableBlocks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSBackupFtpAuthorizableBlocksRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS.",
			},
			"blocks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of CIDR blocks authorized to be granted backup FTP access.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVPSBackupFtpAuthorizableBlocksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	blocks := []string{}
	endpoint := fmt.Sprintf("/vps/%s/backupftp/authorizableBlocks", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, &blocks); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("blocks", blocks)

	return nil
}
