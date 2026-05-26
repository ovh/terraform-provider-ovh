package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceVPSTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSTasksRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: helpers.ValidateEnum(vpsTaskStates),
			},
			"type_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: helpers.ValidateEnum(vpsTaskTypes),
			},
			"task_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceVPSTasksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	queryParams := url.Values{}
	if v, ok := d.GetOk("state_filter"); ok {
		queryParams.Add("state", v.(string))
	}
	if v, ok := d.GetOk("type_filter"); ok {
		queryParams.Add("type", v.(string))
	}

	endpoint := fmt.Sprintf("/vps/%s/tasks", url.PathEscape(serviceName))
	if encoded := queryParams.Encode(); encoded != "" {
		endpoint = endpoint + "?" + encoded
	}

	ids := []int64{}
	if err := config.OVHClient.Get(endpoint, &ids); err != nil {
		return fmt.Errorf("Error calling GET %s: %q", endpoint, err)
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	stringIds := make([]string, len(ids))
	for i, id := range ids {
		stringIds[i] = strconv.FormatInt(id, 10)
	}
	d.SetId(hashcode.Strings(stringIds))
	d.Set("task_ids", ids)
	return nil
}
