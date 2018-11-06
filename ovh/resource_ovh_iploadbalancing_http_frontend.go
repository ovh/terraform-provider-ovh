package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

type IpLoadbalancingHttpFrontend struct {
	FrontendId    int      `json:"frontendId,omitempty"`
	Port          string   `json:"port"`
	Zone          string   `json:"zone"`
	AllowedSource []string `json:"allowedSource,omitempty"`
	DedicatedIpFo []string `json:"dedicatedIpfo,omitempty"`
	DefaultFarmId *int     `json:"defaultFarmId,omitempty"`
	DefaultSslId  *int     `json:"defaultSslId,omitempty"`
	Disabled      *bool    `json:"disabled"`
	Ssl           *bool    `json:"ssl"`
	DisplayName   string   `json:"displayName,omitempty"`
}

func resourceIpLoadbalancinghttpFrontend() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpLoadbalancinghttpFrontendCreate,
		Read:   resourceIpLoadbalancinghttpFrontendRead,
		Update: resourceIpLoadbalancinghttpFrontendUpdate,
		Delete: resourceIpLoadbalancinghttpFrontendDelete,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"allowed_source": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"dedicated_ipfo": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"default_farm_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"default_ssl_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: false,
			},
			"ssl": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: false,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceIpLoadbalancinghttpFrontendCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	allowedSources := stringsFromSchema(d, "allowed_source")
	dedicatedIpFo := stringsFromSchema(d, "dedicated_ipfo")

	for _, s := range allowedSources {
		if err := validateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `allowed_source` value: %s", err)
		}
	}

	for _, s := range dedicatedIpFo {
		if err := validateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `dedicated_ipfo` value: %s", err)
		}
	}

	frontend := &IpLoadbalancinghttpFrontend{
		Port:          d.Get("port").(string),
		Zone:          d.Get("zone").(string),
		AllowedSource: allowedSources,
		DedicatedIpFo: dedicatedIpFo,
		Disabled:      getNilBoolPointer(d.Get("disabled").(bool)),
		Ssl:           getNilBoolPointer(d.Get("ssl").(bool)),
		DisplayName:   d.Get("display_name").(string),
	}

	if farmId, ok := d.GetOk("default_farm_id"); ok {
		frontend.DefaultFarmId = getNilIntPointer(farmId.(int))
	}
	if sslId, ok := d.GetOk("default_ssl_id"); ok {
		frontend.DefaultSslId = getNilIntPointer(sslId.(int))
	}

	service := d.Get("service_name").(string)
	resp := &IpLoadbalancinghttpFrontend{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend", service)

	err := config.OVHClient.Post(endpoint, frontend, resp)
	if err != nil {
		return fmt.Errorf("calling POST %s:\n\t %s", endpoint, err.Error())
	}
	return readIpLoadbalancinghttpFrontend(resp, d)

}

func resourceIpLoadbalancinghttpFrontendRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	r := &IpLoadbalancinghttpFrontend{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%s", service, d.Id())

	err := config.OVHClient.Get(endpoint, &r)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}
	return readIpLoadbalancinghttpFrontend(r, d)
}

func resourceIpLoadbalancinghttpFrontendUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%s", service, d.Id())

	allowedSources := stringsFromSchema(d, "allowed_source")
	dedicatedIpFo := stringsFromSchema(d, "dedicated_ipfo")

	for _, s := range allowedSources {
		if err := validateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `allowed_source` value: %s", err)
		}
	}

	for _, s := range dedicatedIpFo {
		if err := validateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `dedicated_ipfo` value: %s", err)
		}
	}

	frontend := &IpLoadbalancinghttpFrontend{
		Port:          d.Get("port").(string),
		Zone:          d.Get("zone").(string),
		AllowedSource: allowedSources,
		DedicatedIpFo: dedicatedIpFo,
		Disabled:      getNilBoolPointer(d.Get("disabled").(bool)),
		Ssl:           getNilBoolPointer(d.Get("ssl").(bool)),
		DisplayName:   d.Get("display_name").(string),
	}

	if farmId, ok := d.GetOk("default_farm_id"); ok {
		frontend.DefaultFarmId = getNilIntPointer(farmId.(int))
	}
	if sslId, ok := d.GetOk("default_ssl_id"); ok {
		frontend.DefaultSslId = getNilIntPointer(sslId.(int))
	}

	err := config.OVHClient.Put(endpoint, frontend, nil)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}
	return nil
}

func readIpLoadbalancinghttpFrontend(r *IpLoadbalancinghttpFrontend, d *schema.ResourceData) error {
	d.Set("display_name", r.DisplayName)
	d.Set("port", r.Port)
	d.Set("zone", r.Zone)

	allowedSources := make([]string, 0)
	for _, s := range r.AllowedSource {
		allowedSources = append(allowedSources, s)
	}
	d.Set("allowed_source", allowedSources)

	dedicatedIpFos := make([]string, 0)
	for _, s := range r.DedicatedIpFo {
		dedicatedIpFos = append(dedicatedIpFos, s)
	}
	d.Set("dedicated_ipfo", dedicatedIpFos)

	if r.DefaultFarmId != nil {
		d.Set("default_farm_id", r.DefaultFarmId)
	}

	if r.DefaultSslId != nil {
		d.Set("default_ssl_id", r.DefaultSslId)
	}
	if r.Disabled != nil {
		d.Set("disabled", r.Disabled)
	}
	if r.Ssl != nil {
		d.Set("ssl", r.Ssl)
	}

	d.SetId(fmt.Sprintf("%d", r.FrontendId))

	return nil
}

func resourceIpLoadbalancinghttpFrontendDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	r := &IpLoadbalancinghttpFrontend{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%s", service, d.Id())

	err := config.OVHClient.Delete(endpoint, &r)
	if err != nil {
		return fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	return nil
}
