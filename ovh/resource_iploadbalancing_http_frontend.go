package ovh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIpLoadbalancingHttpFrontend() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpLoadbalancingHttpFrontendCreate,
		Read:   resourceIpLoadbalancingHttpFrontendRead,
		Update: resourceIpLoadbalancingHttpFrontendUpdate,
		Delete: resourceIpLoadbalancingHttpFrontendDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpLoadbalancingHttpFrontendImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"allowed_source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dedicated_ipfo": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_farm_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"default_ssl_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"redirect_location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceIpLoadbalancingHttpFrontendImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/frontend id formatted")
	}
	serviceName := splitId[0]
	frontendId := splitId[1]
	d.SetId(frontendId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceIpLoadbalancingHttpFrontendCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	allowedSources, _ := helpers.StringsFromSchema(d, "allowed_source")
	dedicatedIpFo, _ := helpers.StringsFromSchema(d, "dedicated_ipfo")

	for _, s := range allowedSources {
		if err := helpers.ValidateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `allowed_source` value: %s", err)
		}
	}

	for _, s := range dedicatedIpFo {
		if err := helpers.ValidateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `dedicated_ipfo` value: %s", err)
		}
	}

	frontend := &IpLoadbalancingHttpFrontend{
		Port:             d.Get("port").(string),
		Zone:             d.Get("zone").(string),
		AllowedSource:    allowedSources,
		DedicatedIpFo:    dedicatedIpFo,
		Disabled:         d.Get("disabled").(bool),
		Ssl:              d.Get("ssl").(bool),
		RedirectLocation: d.Get("redirect_location").(string),
		DisplayName:      d.Get("display_name").(string),
	}

	frontend.DefaultFarmId = helpers.GetNilIntPointerFromData(d, "default_farm_id")
	frontend.DefaultSslId = helpers.GetNilIntPointerFromData(d, "default_ssl_id")

	service := d.Get("service_name").(string)
	resp := &IpLoadbalancingHttpFrontend{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend", service)

	err := config.OVHClient.Post(endpoint, frontend, resp)
	if err != nil {
		return fmt.Errorf("calling POST %s:\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", resp.FrontendId))

	return resourceIpLoadbalancingHttpFrontendRead(d, meta)
}

func resourceIpLoadbalancingHttpFrontendRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	r := &IpLoadbalancingHttpFrontend{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%s", service, d.Id())

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(fmt.Sprintf("%d", r.FrontendId))

	allowedSources := make([]string, 0)
	allowedSources = append(allowedSources, r.AllowedSource...)

	dedicatedIpFos := make([]string, 0)
	dedicatedIpFos = append(dedicatedIpFos, r.DedicatedIpFo...)

	d.Set("allowed_source", allowedSources)
	d.Set("dedicated_ipfo", dedicatedIpFos)
	d.Set("default_farm_id", r.DefaultFarmId)
	d.Set("default_ssl_id", r.DefaultSslId)
	d.Set("disabled", r.Disabled)
	d.Set("display_name", r.DisplayName)
	d.Set("port", r.Port)
	d.Set("ssl", r.Ssl)
	d.Set("zone", r.Zone)
	d.Set("redirect_location", r.RedirectLocation)

	return nil
}

func resourceIpLoadbalancingHttpFrontendUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%s", service, d.Id())

	allowedSources, _ := helpers.StringsFromSchema(d, "allowed_source")
	dedicatedIpFo, _ := helpers.StringsFromSchema(d, "dedicated_ipfo")

	for _, s := range allowedSources {
		if err := helpers.ValidateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `allowed_source` value: %s", err)
		}
	}

	for _, s := range dedicatedIpFo {
		if err := helpers.ValidateIpBlock(s); err != nil {
			return fmt.Errorf("Error validating `dedicated_ipfo` value: %s", err)
		}
	}

	frontend := &IpLoadbalancingHttpFrontend{
		Port:             d.Get("port").(string),
		Zone:             d.Get("zone").(string),
		AllowedSource:    allowedSources,
		DedicatedIpFo:    dedicatedIpFo,
		Disabled:         d.Get("disabled").(bool),
		Ssl:              d.Get("ssl").(bool),
		RedirectLocation: d.Get("redirect_location").(string),
		DisplayName:      d.Get("display_name").(string),
	}

	frontend.DefaultFarmId = helpers.GetNilIntPointerFromData(d, "default_farm_id")
	frontend.DefaultSslId = helpers.GetNilIntPointerFromData(d, "default_ssl_id")

	err := config.OVHClient.Put(endpoint, frontend, nil)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
	}

	return resourceIpLoadbalancingHttpFrontendRead(d, meta)
}

func resourceIpLoadbalancingHttpFrontendDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := d.Get("service_name").(string)
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%s", service, d.Id())

	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	d.SetId("")
	return nil
}
