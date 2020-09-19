package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"
)

const EnterpriseCloudDBBaseUrl = "/cloudDB/enterprise/cluster/%s"

func dataSourceEnterpriseCloudDb() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEnterpriseCloudDB,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: false,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateEnterpriseCloudDBStatus(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func dataSourceEnterpriseCloudDB(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	clusterId := d.Get("cluster_id").(string)
	log.Printf("[DEBUG] Will retrieve enterprise cloud db %s", clusterId)

	db := &EnterpriseCloudDB{}
	if err := config.OVHClient.Get(fmt.Sprintf(EnterpriseCloudDBBaseUrl, clusterId), &db); err != nil {
		return fmt.Errorf("Error calling %s/%s:\n\t %q", EnterpriseCloudDBBaseUrl, clusterId, err)
	}
	d.SetId(db.Id)
	d.Set("region", db.RegionName)
	d.Set("status", db.Status)
	d.Set("cluster_id", db.Id)

	return nil
}
