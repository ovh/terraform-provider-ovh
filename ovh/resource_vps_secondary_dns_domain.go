package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceVPSSecondaryDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSSecondaryDNSDomainCreate,
		Read:   resourceVPSSecondaryDNSDomainRead,
		Delete: resourceVPSSecondaryDNSDomainDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVPSSecondaryDNSDomainImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// PUT deprecated 2025-10-15, ip changes require recreate.
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if err := helpers.ValidateIpV4(v.(string)); err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_master": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPSSecondaryDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	domain := d.Get("domain").(string)

	opts := &VPSSecondaryDNSDomainCreateOpts{Domain: domain}
	if v, ok := d.GetOk("ip"); ok {
		s := v.(string)
		opts.Ip = &s
	}

	endpoint := fmt.Sprintf("/vps/%s/secondaryDnsDomains", url.PathEscape(serviceName))
	if err := config.OVHClient.Post(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t%q", endpoint, err)
	}

	d.SetId(serviceName + "|" + domain)
	return resourceVPSSecondaryDNSDomainRead(d, meta)
}

func resourceVPSSecondaryDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	domain := d.Get("domain").(string)

	endpoint := fmt.Sprintf("/vps/%s/secondaryDnsDomains/%s",
		url.PathEscape(serviceName), url.PathEscape(domain))
	resp := &VPSSecondaryDNSDomain{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", serviceName)
	d.Set("domain", resp.Domain)
	d.Set("dns", resp.Dns)
	d.Set("ip_master", resp.IpMaster)
	d.Set("creation_date", resp.CreationDate)
	return nil
}

func resourceVPSSecondaryDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	domain := d.Get("domain").(string)

	endpoint := fmt.Sprintf("/vps/%s/secondaryDnsDomains/%s",
		url.PathEscape(serviceName), url.PathEscape(domain))
	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}
	d.SetId("")
	return nil
}

func resourceVPSSecondaryDNSDomainImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "|", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("import id must be service_name|domain")
	}
	d.Set("service_name", parts[0])
	d.Set("domain", parts[1])
	return []*schema.ResourceData{d}, nil
}
