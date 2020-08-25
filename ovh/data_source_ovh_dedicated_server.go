package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedServerRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"boot_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "",
			},
			"commercial_range": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "dedicater server commercial range",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "dedicated datacenter localisation (bhs1,bhs2,...)",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "dedicated server ip (IPv4)",
			},
			"ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "dedicated server ip blocks",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"link_speed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "",
			},
			"monitoring": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Icmp monitoring state",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "dedicated server name",
			},
			"os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating system",
			},
			"professional_use": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Does this server have professional use option",
			},
			"rack": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"rescue_mail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"reverse": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "dedicated server reverse",
			},
			"root_device": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"server_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "your server id",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "error, hacked, hackedBlocked, ok",
			},
			"support_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dedicated server support level (critical, fastpath, gs, pro)",
			},
			"vnis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "VirtualNetworkInterface activation state",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VirtualNetworkInterface mode (public,vrack,vrack_aggregation)",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User defined VirtualNetworkInterface name",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VirtualNetworkInterface unique id",
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "server name",
						},
						"vrack": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vRack name",
						},
						"nics": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "NetworkInterfaceControllers bound to this VirtualNetworkInterface",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"enabled_vrack_vnis": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of enabled vrack VNI uuids",
			},
			"enabled_vrack_aggregation_vnis": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of enabled vrack_aggregation VNI uuids",
			},
			"enabled_public_vnis": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of enabled public VNI uuids",
			},
		},
	}
}

func dataSourceDedicatedServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	ds := &DedicatedServer{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s",
			url.PathEscape(serviceName),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"Error calling /dedicated/server/%s:\n\t %q",
			serviceName,
			err,
		)
	}

	d.SetId(ds.Name)
	d.Set("boot_id", ds.BootId)
	d.Set("commercial_range", ds.CommercialRange)
	d.Set("datacenter", ds.Datacenter)
	d.Set("ip", ds.Ip)
	d.Set("link_speed", ds.LinkSpeed)
	d.Set("monitoring", ds.Monitoring)
	d.Set("os", ds.Os)
	d.Set("name", ds.Name)
	d.Set("professional_use", ds.ProfessionalUse)
	d.Set("rack", ds.Rack)
	d.Set("rescue_mail", ds.RescueMail)
	d.Set("reverse", ds.Reverse)
	d.Set("root_device", ds.RootDevice)
	d.Set("server_id", ds.ServerId)
	d.Set("state", ds.State)
	d.Set("support_level", ds.SupportLevel)

	dsIps := &[]string{}
	err = config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s/ips",
			url.PathEscape(serviceName),
		),
		&dsIps,
	)

	if err != nil {
		return fmt.Errorf(
			"Error reading Dedicated Server IPs for %s: %q",
			serviceName,
			err,
		)
	}

	d.Set("ips", dsIps)

	// Set VNIs attributes
	vnis, err := getDedicatedServerVNIs(d, meta)

	if err != nil {
		return fmt.Errorf("Error reading Dedicated Server VNIs: %s", err)
	}

	mapvnis := make([]map[string]interface{}, len(vnis))
	enabledVrackVnis := []string{}
	enabledVrackAggregationVnis := []string{}
	enabledPublicVnis := []string{}

	for i, vni := range vnis {
		mapvnis[i] = vni.ToMap()

		if vni.Enabled {
			switch vni.Mode {
			case "vrack":
				enabledVrackVnis = append(enabledVrackVnis, vni.Uuid)
			case "vrack_aggregation":
				enabledVrackAggregationVnis = append(enabledVrackAggregationVnis, vni.Uuid)
			case "public":
				enabledPublicVnis = append(enabledPublicVnis, vni.Uuid)
			default:
				log.Printf("[WARN] unknown VNI mode. DS {%v} VNI {%v}", ds, vni)
			}
		}
	}

	d.Set("vnis", mapvnis)
	d.Set("enabled_vrack_vnis", enabledVrackVnis)
	d.Set("enabled_vrack_aggregation_vnis", enabledVrackAggregationVnis)
	d.Set("enabled_public_vnis", enabledPublicVnis)

	return nil
}

func getDedicatedServerVNIs(d *schema.ResourceData, meta interface{}) ([]*DedicatedServerVNI, error) {
	config := meta.(*Config)

	log.Printf("[INFO] Getting VNIs for dedicated server: %s", d.Id())

	serviceName := d.Get("service_name").(string)

	// First get ids unfiltered
	ids := []string{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s/virtualNetworkInterface",
			url.PathEscape(serviceName),
		),
		&ids,
	)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving VNIs for dedicated server %s: %s", d.Id(), err)
	}

	if len(ids) < 1 {
		log.Printf("[WARN] Dedicated server %s returned no VNI. Your server might be on legacy network infrastructure.", d.Id())
		return nil, nil
	}

	vnis := []*DedicatedServerVNI{}

	for i := 0; i < len(ids); i++ {
		vni := &DedicatedServerVNI{}
		err := config.OVHClient.Get(
			fmt.Sprintf("/dedicated/server/%s/virtualNetworkInterface/%s", serviceName, ids[i]),
			vni,
		)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving VNI info for dedicated server %s: %s", d.Id(), err)
		}
		vnis = append(vnis, vni)
	}

	return vnis, nil
}
