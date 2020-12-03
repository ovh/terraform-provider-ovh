package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIpLoadbalancingTcpFarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpLoadbalancingTcpFarmCreate,
		Read:   resourceIpLoadbalancingTcpFarmRead,
		Update: resourceIpLoadbalancingTcpFarmUpdate,
		Delete: resourceIpLoadbalancingTcpFarmDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpLoadbalancingTcpFarmImportState,
		},

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
					err := helpers.ValidateStringEnum(v.(string), []string{"first", "leastconn", "roundrobin", "source"})
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
					err := helpers.ValidateStringEnum(v.(string), []string{"sourceIp"})
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
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								err := helpers.ValidateStringEnum(v.(string), []string{"contains", "default", "internal", "matches", "status"})
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
							Type:         schema.TypeBool,
							Default:      false,
							RequiredWith: []string{"probe.0.match"},
							Optional:     true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"force_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								err := helpers.ValidateStringEnum(v.(string), []string{"GET", "HEAD", "OPTIONS", "internal"})
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
								err := helpers.ValidateStringEnum(v.(string), []string{"http", "internal", "mysql", "oco", "pgsql", "smtp", "tcp"})
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

func resourceIpLoadbalancingTcpFarmImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/farm id formatted")
	}
	serviceName := splitId[0]
	farmId := splitId[1]
	d.SetId(farmId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceIpLoadbalancingTcpFarmCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	farm := (&IpLoadbalancingFarmCreateOrUpdateOpts{}).FromResource(d)
	service := d.Get("service_name").(string)
	resp := &IpLoadbalancingFarm{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm", service)

	err := config.OVHClient.Post(endpoint, farm, resp)
	if err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", resp.FarmId))

	return resourceIpLoadbalancingTcpFarmRead(d, meta)
}

func resourceIpLoadbalancingTcpFarmRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%s", service, d.Id())
	r := &IpLoadbalancingFarm{}

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	probes := make([]map[string]interface{}, 0)
	if r.Probe != nil && r.Probe.ToMap() != nil {
		probes = append(probes, r.Probe.ToMap())
	}

	d.Set("display_name", r.DisplayName)
	d.Set("zone", r.Zone)
	d.Set("port", r.Port)
	d.Set("balance", r.Balance)
	d.Set("probe", probes)
	d.Set("vrack_network_id", r.VrackNetworkId)
	d.Set("stickiness", r.Stickiness)
	return nil
}

func resourceIpLoadbalancingTcpFarmUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%s", service, d.Id())

	farm := (&IpLoadbalancingFarmCreateOrUpdateOpts{}).FromResource(d)

	err := config.OVHClient.Put(endpoint, farm, nil)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}

	return resourceIpLoadbalancingTcpFarmRead(d, meta)
}

func resourceIpLoadbalancingTcpFarmDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	r := &IpLoadbalancingFarm{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%s", service, d.Id())

	err := config.OVHClient.Delete(endpoint, &r)
	if err != nil {
		return fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	d.SetId("")
	return nil
}
