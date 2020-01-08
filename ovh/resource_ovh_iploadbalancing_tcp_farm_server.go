package ovh

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type IpLoadbalancingTcpFarmServer struct {
	BackendId            int     `json:"backendId,omitempty"`
	ServerId             int     `json:"serverId,omitempty"`
	FarmId               int     `json:"farmId,omitempty"`
	DisplayName          *string `json:"displayName,omitempty"`
	Address              string  `json:"address"`
	Cookie               *string `json:"cookie,omitempty"`
	Port                 *int    `json:"port"`
	ProxyProtocolVersion *string `json:"proxyProtocolVersion"`
	Chain                *string `json:"chain"`
	Weight               *int    `json:"weight"`
	Probe                *bool   `json:"probe"`
	Ssl                  *bool   `json:"ssl"`
	Backup               *bool   `json:"backup"`
	Status               string  `json:"status"`
}

func resourceIpLoadbalancingTcpFarmServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpLoadbalancingTcpFarmServerCreate,
		Read:   resourceIpLoadbalancingTcpFarmServerRead,
		Update: resourceIpLoadbalancingTcpFarmServerUpdate,
		Delete: resourceIpLoadbalancingTcpFarmServerDelete,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"farm_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					ip := v.(string)
					if net.ParseIP(ip).To4() == nil {
						errors = append(errors, fmt.Errorf("Address %s is not an IPv4", ip))
					}
					return
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cookie": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"proxy_protocol_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := validateStringEnum(v.(string), []string{"v1", "v2", "v2-ssl", "v2-ssl-cn"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"chain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"probe": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"backup": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := validateStringEnum(v.(string), []string{"active", "inactive"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func resourceIpLoadbalancingTcpFarmServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	newBackendServer := &IpLoadbalancingTcpFarmServer{
		DisplayName:          getNilStringPointerFromData(d, "display_name"),
		Address:              d.Get("address").(string),
		Port:                 getNilIntPointerFromData(d, "port"),
		ProxyProtocolVersion: getNilStringPointerFromData(d, "proxy_protocol_version"),
		Chain:                getNilStringPointerFromData(d, "chain"),
		Weight:               getNilIntPointerFromData(d, "weight"),
		Probe:                getNilBoolPointerFromData(d, "probe"),
		Ssl:                  getNilBoolPointerFromData(d, "ssl"),
		Backup:               getNilBoolPointerFromData(d, "backup"),
		Status:               d.Get("status").(string),
	}

	service := d.Get("service_name").(string)
	farmid := d.Get("farm_id").(int)
	r := &IpLoadbalancingTcpFarmServer{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d/server", service, farmid)

	err := config.OVHClient.Post(endpoint, newBackendServer, r)
	if err != nil {
		return fmt.Errorf("calling %s with %d:\n\t %s", endpoint, farmid, err.Error())
	}

	//set id
	d.SetId(fmt.Sprintf("%d", r.ServerId))

	return nil
}

func resourceIpLoadbalancingTcpFarmServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	farmid := d.Get("farm_id").(int)
	r := &IpLoadbalancingTcpFarmServer{}

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d/server/%s", service, farmid, d.Id())

	err := config.OVHClient.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling %s :\n\t %q", endpoint, err)
	}
	log.Printf("[DEBUG] Response object from OVH : %v", r)

	d.Set("probe", *r.Probe)
	d.Set("ssl", *r.Ssl)
	d.Set("backup", *r.Backup)
	d.Set("address", r.Address)
	if r.DisplayName != nil {
		d.Set("display_name", *r.DisplayName)
	}
	if r.Cookie != nil {
		d.Set("cookie", *r.Cookie)
	}
	d.Set("port", *r.Port)
	if r.ProxyProtocolVersion != nil {
		d.Set("proxy_protocol_version", *r.ProxyProtocolVersion)
	}
	if r.Chain != nil {
		d.Set("chain", *r.Chain)
	}
	d.Set("weight", *r.Weight)
	d.Set("status", r.Status)

	return nil
}

func resourceIpLoadbalancingTcpFarmServerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	update := &IpLoadbalancingTcpFarmServer{
		DisplayName:          getNilStringPointerFromData(d, "display_name"),
		Address:              d.Get("address").(string),
		Port:                 getNilIntPointerFromData(d, "port"),
		ProxyProtocolVersion: getNilStringPointerFromData(d, "proxy_protocol_version"),
		Chain:                getNilStringPointerFromData(d, "chain"),
		Weight:               getNilIntPointerFromData(d, "weight"),
		Probe:                getNilBoolPointerFromData(d, "probe"),
		Ssl:                  getNilBoolPointerFromData(d, "ssl"),
		Backup:               getNilBoolPointerFromData(d, "backup"),
		Status:               d.Get("status").(string),
	}

	service := d.Get("service_name").(string)
	farmid := d.Get("farm_id").(int)
	r := &IpLoadbalancingTcpFarmServer{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d/server/%s", service, farmid, d.Id())
	js, _ := json.Marshal(update)
	log.Printf("[DEBUG] PUT %s : %v", endpoint, string(js))
	err := config.OVHClient.Put(endpoint, update, r)
	if err != nil {
		return fmt.Errorf("calling %s with %d:\n\t %s", endpoint, farmid, err.Error())
	}
	return nil
}

func resourceIpLoadbalancingTcpFarmServerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	farmid := d.Get("farm_id").(int)

	r := &IpLoadbalancingTcpFarmServer{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d/server/%s", service, farmid, d.Id())

	err := config.OVHClient.Delete(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling %s :\n\t %s", endpoint, err.Error())
	}

	return nil
}
