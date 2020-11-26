package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Here come all the computed items
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"displayname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"netbootmode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slamonitoring": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"keymap": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcore": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"offertype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"datacenter": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceVPSRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	vps := &VPS{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/vps/%s",
			url.PathEscape(serviceName),
		),
		&vps,
	)

	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(vps.Name)
	d.Set("name", vps.Name)
	d.Set("zone", vps.Zone)
	d.Set("state", vps.State)
	d.Set("vcore", vps.Vcore)
	d.Set("cluster", vps.Cluster)
	d.Set("memory", vps.Memory)
	d.Set("netbootmode", vps.NetbootMode)
	d.Set("keymap", vps.Keymap)
	d.Set("offertype", vps.OfferType)
	d.Set("slamonitoring", vps.SlaMonitorting)
	d.Set("displayname", vps.DisplayName)
	var model = make(map[string]string)
	model["name"] = vps.Model.Name
	model["offer"] = vps.Model.Offer
	model["version"] = vps.Model.Version
	d.Set("model", model)
	d.Set("type", ovhvps_getType(vps.OfferType, vps.Model.Name, vps.Model.Version))

	ips := []string{}
	err = config.OVHClient.Get(
		fmt.Sprintf("/vps/%s/ips", d.Id()),
		&ips,
	)

	d.Set("ips", ips)

	vpsDatacenter := VPSDatacenter{}
	err = config.OVHClient.Get(
		fmt.Sprintf("/vps/%s/datacenter", d.Id()),
		&vpsDatacenter,
	)
	datacenter := make(map[string]string)
	datacenter["name"] = vpsDatacenter.Name
	datacenter["longname"] = vpsDatacenter.Longname
	d.Set("datacenter", datacenter)
	return nil
}
