package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VPSDiskMonitoring mirrors complexType.UnitAndValues_vps.VpsTimestampValue
// returned by GET /vps/{serviceName}/disks/{id}/monitoring.
type VPSDiskMonitoring struct {
	Unit   string `json:"unit"`
	Values []struct {
		Timestamp string  `json:"timestamp"`
		Value     float64 `json:"value"`
	} `json:"values"`
}

func dataSourceVPSDiskMonitoring() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSDiskMonitoringRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"period": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Period of the monitoring window: lastday, lastweek, lastmonth, lastyear, today.",
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					switch v.(string) {
					case "lastday", "lastweek", "lastmonth", "lastyear", "today":
					default:
						errs = append(errs, fmt.Errorf("%s must be one of lastday,lastweek,lastmonth,lastyear,today", k))
					}
					return
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Statistic type to query. The API accepts vps.VpsStatisticTypeEnum values; only the disk-related ones are meaningful here.",
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					switch v.(string) {
					case "cpu:iowait", "cpu:max", "cpu:nice", "cpu:sys", "cpu:used", "cpu:user",
						"mem:max", "mem:used",
						"net:rx", "net:tx":
					default:
						errs = append(errs, fmt.Errorf("%s must be one of cpu:iowait,cpu:max,cpu:nice,cpu:sys,cpu:used,cpu:user,mem:max,mem:used,net:rx,net:tx", k))
					}
					return
				},
			},
			// Computed
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVPSDiskMonitoringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	diskID := int64(d.Get("disk_id").(int))
	period := d.Get("period").(string)
	statType := d.Get("type").(string)

	qs := url.Values{}
	qs.Set("period", period)
	qs.Set("type", statType)

	endpoint := fmt.Sprintf("/vps/%s/disks/%d/monitoring?%s",
		url.PathEscape(serviceName), diskID, qs.Encode())

	resp := &VPSDiskMonitoring{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%s|%d|%s|%s", serviceName, diskID, period, statType))
	d.Set("unit", resp.Unit)

	values := make([]map[string]interface{}, 0, len(resp.Values))
	for _, v := range resp.Values {
		values = append(values, map[string]interface{}{
			"timestamp": v.Timestamp,
			"value":     v.Value,
		})
	}
	d.Set("values", values)
	return nil
}
