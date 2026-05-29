package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSDisksRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceVPSDisksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/disks", url.PathEscape(serviceName))
	var ids []int64
	if err := config.OVHClient.Get(endpoint, &ids); err != nil {
		return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = strconv.FormatInt(id, 10)
	}

	d.SetId(hashcode.Strings(append([]string{serviceName}, stringIDs...)))
	d.Set("result", ids)
	return nil
}
