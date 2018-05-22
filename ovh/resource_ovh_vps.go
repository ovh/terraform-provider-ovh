package ovh

import (
	"fmt"
	//	"log"
	//	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

type VPSModel struct {
	Name    string `json:"name"`
	Offer   string `json:"offer"`
	Memory  int    `json:"memory"`
	Vcore   int    `json:"vcore"`
	Version string `json:"version"`
	Disk    int    `json:"disk"`
}

type VPS struct {
	Name           string   `json:"name"`
	Cluster        string   `json:"cluster"`
	Memory         int      `json:"memoryLimit"`
	NetbootMode    string   `json:"netbootMode"`
	Keymap         string   `json:"keymap"`
	Zone           string   `json:"zone"`
	State          string   `json:"state"`
	Vcore          int      `json:"vcore"`
	OfferType      string   `json:"offerType"`
	SlaMonitorting bool     `json:"slaMonitoring"`
	DisplayName    string   `json:"displayName"`
	Model          VPSModel `json:"model"`
}

type VPSDatacenter struct {
	Name     string `json:"name"`
	Longname string `json:"longname"`
}

type VPSProperties struct {
	NetbootMode    *string `json:"netbootMode"`
	Keymap         *string `json:"keymap"`
	SlaMonitorting bool    `json:"slaMonitoring"`
	DisplayName    *string `json:"displayName"`
}

func resourceOVHVPS() *schema.Resource {
	return &schema.Resource{
		Create: resourceOVHVPSCreate,
		Read:   resourceOVHVPSRead,
		Update: resourceOVHVPSUpdate,
		Delete: resourceOVHVPSDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "running",
			},
			"displayname": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"netbootmode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "local",
			},
			"slamonitoring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"keymap": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			// Here come all the computed items
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
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"offer": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"datacenter": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"longname": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceOVHVPSCreate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Creation is not supported at the moment, please use import")
}

func resourceOVHVPSRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)
	vps := VPS{}
	err := provider.OVHClient.Get(
		fmt.Sprintf("/vps/%s", d.Id()),
		&vps,
	)

	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(vps.Name)
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
	err = provider.OVHClient.Get(
		fmt.Sprintf("/vps/%s/ips", d.Id()),
		&ips,
	)

	d.Set("ips", ips)

	vpsDatacenter := VPSDatacenter{}
	err = provider.OVHClient.Get(
		fmt.Sprintf("/vps/%s/datacenter", d.Id()),
		&vpsDatacenter,
	)
	datacenter := make(map[string]string)
	datacenter["name"] = vpsDatacenter.Name
	datacenter["longname"] = vpsDatacenter.Longname
	d.Set("datacenter", datacenter)
	return nil
}

func resourceOVHVPSUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)
	const Null = "\xff"
	if d.HasChange("type") {
		return fmt.Errorf("Type modification is not supported at the moment")
	}
	if d.HasChange("keymap") || d.HasChange("displayname") || d.HasChange("netbootmode") || d.HasChange("slamonitoring") {
		newProperties := &VPSProperties{
			Keymap:         strPtr(d.Get("keymap").(string)),
			DisplayName:    strPtr(d.Get("displayname").(string)),
			NetbootMode:    strPtr(d.Get("netbootmode").(string)),
			SlaMonitorting: d.Get("slamonitoring").(bool),
		}
		if *newProperties.Keymap == "" {
			newProperties.Keymap = nil
		}
		err := provider.OVHClient.Put(
			fmt.Sprintf("/vps/%s", d.Id()),
			newProperties,
			nil,
		)
		if err != nil {
			return fmt.Errorf("Failed to update VPS: %s", err)
		}
	}
	if d.HasChange("state") {
		command := ""
		switch d.Get("state") {
		case "running":
			command = "start"
		case "stopped":
			command = "stop"
		default:
			return fmt.Errorf("Unknown wanted state")
		}
		err := provider.OVHClient.Post(
			fmt.Sprintf("/vps/%s/%s", d.Id(), command),
			nil,
			nil,
		)
		if err != nil {
			return fmt.Errorf("Failed to update VPS: %s", err)
		}
	}
	return nil
}

func resourceOVHVPSDelete(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Deletion is not supported at the moment, please terminate manually and terraform state rm")
}

func ovhvps_getType(offertype string, model_name string, model_version string) string {
	var offertypeToOfferPrefix = make(map[string]string)
	offertypeToOfferPrefix["cloud"] = "ceph-nvme"
	offertypeToOfferPrefix["cloud-ram"] = "ceph-nvme-ram"
	offertypeToOfferPrefix["ssd"] = "ssd"
	offertypeToOfferPrefix["classic"] = "classic"
	return (fmt.Sprintf("vps_%s_%s_%s", offertypeToOfferPrefix[offertype],
		model_name,
		model_version))
}

func strPtr(s string) *string {
	return &s
}
