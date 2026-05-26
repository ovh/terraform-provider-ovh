package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSVeeamRestorePoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSVeeamRestorePointsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter restore points by creation time (RFC3339 format)",
			},
			"result": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Veeam restore point IDs",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceVPSVeeamRestorePointsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/veeam/restorePoints", url.PathEscape(serviceName))
	if v, ok := d.GetOk("creation_time"); ok {
		qs := url.Values{}
		qs.Add("creationTime", v.(string))
		endpoint = endpoint + "?" + qs.Encode()
	}

	ids := []int64{}
	if err := config.OVHClient.Get(endpoint, &ids); err != nil {
		return fmt.Errorf("error calling GET %s: %s", endpoint, err)
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	strs := make([]string, len(ids))
	for i, v := range ids {
		strs[i] = strconv.FormatInt(v, 10)
	}
	d.SetId(hashcode.Strings(strs))
	d.Set("result", ids)
	return nil
}
