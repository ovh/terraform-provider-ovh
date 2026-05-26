package ovh

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSTaskRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSTaskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	taskId := int64(d.Get("id").(int))

	endpoint := fmt.Sprintf(
		"/vps/%s/tasks/%s",
		url.PathEscape(serviceName),
		url.PathEscape(strconv.FormatInt(taskId, 10)),
	)

	task := &VPSTask{}
	if err := config.OVHClient.Get(endpoint, task); err != nil {
		return fmt.Errorf("Error calling GET %s: %q", endpoint, err)
	}

	d.SetId(strconv.FormatInt(task.Id, 10))
	d.Set("date", task.Date)
	d.Set("type", task.Type)
	d.Set("state", task.State)
	d.Set("progress", task.Progress)
	return nil
}
