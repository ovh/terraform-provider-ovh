package ovh

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSAutomatedBackupRestorePoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSAutomatedBackupRestorePointsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS",
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter restore points by state (available, restored, restoring)",
				ValidateFunc: helpers.ValidateEnum([]string{"available", "restored", "restoring"}),
			},
			"restore_points": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of restore points (RFC3339 datetimes)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVPSAutomatedBackupRestorePointsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/automatedBackup/restorePoints", url.PathEscape(serviceName))
	if v, ok := d.GetOk("state"); ok {
		endpoint = fmt.Sprintf("%s?state=%s", endpoint, url.QueryEscape(v.(string)))
	}

	points := []string{}
	if err := config.OVHClient.Get(endpoint, &points); err != nil {
		return fmt.Errorf("error calling GET %s: %w", endpoint, err)
	}

	sort.Strings(points)

	d.SetId(hashcode.Strings(append([]string{serviceName, d.Get("state").(string)}, points...)))
	d.Set("restore_points", points)
	return nil
}
