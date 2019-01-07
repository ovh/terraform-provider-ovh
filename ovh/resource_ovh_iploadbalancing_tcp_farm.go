package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

type IpLoadbalancingTcpFarmBackendProbe struct {
	Match    string `json:"match,omitempty"`
	Port     int    `json:"port,omitempty"`
	Interval int    `json:"interval,omitempty"`
	Negate   bool   `json:"negate,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
	ForceSsl bool   `json:"forceSsl,omitempty"`
	URL      string `json:"url,omitempty"`
	Method   string `json:"method,omitempty"`
	Type     string `json:"type,omitempty"`
}

type IpLoadbalancingTcpFarm struct {
	FarmId         int                                 `json:"farmId,omitempty"`
	Zone           string                              `json:"zone,omitempty"`
	VrackNetworkId int                                 `json:"vrackNetworkId,omitempty"`
	Port           int                                 `json:"port,omitempty"`
	Stickiness     string                              `json:"stickiness,omitempty"`
	Balance        string                              `json:"balance,omitempty"`
	Probe          *IpLoadbalancingTcpFarmBackendProbe `json:"probe,omitempty"`
	DisplayName    string                              `json:"displayName,omitempty"`
}

func resourceIpLoadbalancingTcpFarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpLoadbalancingTcpFarmCreate,
		Read:   resourceIpLoadbalancingTcpFarmRead,
		Update: resourceIpLoadbalancingTcpFarmUpdate,
		Delete: resourceIpLoadbalancingTcpFarmDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"balance": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := validateStringEnum(v.(string), []string{"first", "leastconn", "roundrobin", "source"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"stickiness": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := validateStringEnum(v.(string), []string{"sourceIp"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"vrack_network_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"probe": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								err := validateStringEnum(v.(string), []string{"contains", "default", "internal", "matches", "status"})
								if err != nil {
									errors = append(errors, err)
								}
								return
							},
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value < 30 || value > 3600 {
									errors = append(errors, fmt.Errorf("Probe interval not in 30..3600 range"))
								}
								return
							},
						},
						"negate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"force_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								err := validateStringEnum(v.(string), []string{"GET", "HEAD", "OPTIONS", "internal"})
								if err != nil {
									errors = append(errors, err)
								}
								return
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								err := validateStringEnum(v.(string), []string{"http", "internal", "mysql", "oko", "pgsql", "smtp", "tcp"})
								if err != nil {
									errors = append(errors, err)
								}
								return
							},
						},
					},
				},
			},
		},
	}
}

func resourceIpLoadbalancingTcpFarmCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	probe := &IpLoadbalancingTcpFarmBackendProbe{}
	probeSet := d.Get("probe").(*schema.Set)
	if probeSet.Len() > 0 {
		probeData := probeSet.List()[0].(map[string]interface{})
		probe.Match = probeData["match"].(string)
		probe.Port = probeData["port"].(int)
		probe.Interval = probeData["interval"].(int)
		probe.Negate = probeData["negate"].(bool)
		probe.Pattern = probeData["pattern"].(string)
		probe.ForceSsl = probeData["force_ssl"].(bool)
		probe.URL = probeData["url"].(string)
		probe.Method = probeData["method"].(string)
		probe.Type = probeData["type"].(string)
	}

	farm := &IpLoadbalancingTcpFarm{
		Zone:           d.Get("zone").(string),
		VrackNetworkId: d.Get("vrack_network_id").(int),
		Port:           d.Get("port").(int),
		Stickiness:     d.Get("stickiness").(string),
		Balance:        d.Get("balance").(string),
		Probe:          probe,
		DisplayName:    d.Get("display_name").(string),
	}

	service := d.Get("service_name").(string)
	resp := &IpLoadbalancingTcpFarm{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm", service)

	err := config.OVHClient.Post(endpoint, farm, resp)
	if err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", resp.FarmId))

	return nil
}

func resourceIpLoadbalancingTcpFarmRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	r := &IpLoadbalancingTcpFarm{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%s", service, d.Id())

	err := config.OVHClient.Get(endpoint, &r)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}

	d.Set("display_name", r.DisplayName)

	return nil
}

func resourceIpLoadbalancingTcpFarmUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%s", service, d.Id())

	probe := &IpLoadbalancingTcpFarmBackendProbe{}
	probeSet := d.Get("probe").(*schema.Set)
	if probeSet.Len() > 0 {
		probeData := probeSet.List()[0].(map[string]interface{})
		probe.Match = probeData["match"].(string)
		probe.Port = probeData["port"].(int)
		probe.Interval = probeData["interval"].(int)
		probe.Negate = probeData["negate"].(bool)
		probe.Pattern = probeData["pattern"].(string)
		probe.ForceSsl = probeData["force_ssl"].(bool)
		probe.URL = probeData["url"].(string)
		probe.Method = probeData["method"].(string)
		probe.Type = probeData["type"].(string)
	}

	farm := &IpLoadbalancingTcpFarm{
		VrackNetworkId: d.Get("vrack_network_id").(int),
		Port:           d.Get("port").(int),
		Stickiness:     d.Get("stickiness").(string),
		Balance:        d.Get("balance").(string),
		Probe:          probe,
		DisplayName:    d.Get("display_name").(string),
	}

	err := config.OVHClient.Put(endpoint, farm, nil)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}

	return nil
}

func resourceIpLoadbalancingTcpFarmDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	r := &IpLoadbalancingTcpFarm{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%s", service, d.Id())

	err := config.OVHClient.Delete(endpoint, &r)
	if err != nil {
		return fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	return nil
}
