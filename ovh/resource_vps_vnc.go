package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

// VPSVnc maps the response of POST /vps/{serviceName}/openConsoleAccess.
type VPSVnc struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type VPSVncOpenConsoleAccessOpts struct {
	Protocol string `json:"protocol"`
}

func resourceVPSVnc() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSVncCreate,
		Read:   resourceVPSVncRead,
		Delete: resourceVPSVncDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service_name of your VPS.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"VNC", "VNCOverWebSocket"}),
				Description:  "The VNC protocol to open: VNC or VNCOverWebSocket.",
			},

			// Computed
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VNC host to connect to.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VNC port to connect to.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The one-shot VNC password. Issued by OVHcloud on each create and intended for immediate use; the session auto-expires server-side.",
			},
		},
	}
}

func resourceVPSVncCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	protocol := d.Get("protocol").(string)

	endpoint := fmt.Sprintf("/vps/%s/openConsoleAccess", url.PathEscape(serviceName))
	opts := VPSVncOpenConsoleAccessOpts{Protocol: protocol}

	log.Printf("[DEBUG] Opening VPS VNC console for %s (protocol=%s)", serviceName, protocol)

	vnc := &VPSVnc{}
	if err := config.OVHClient.Post(endpoint, opts, vnc); err != nil {
		return fmt.Errorf("calling POST %s: %w", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, protocol))
	d.Set("host", vnc.Host)
	d.Set("port", vnc.Port)
	d.Set("password", vnc.Password)

	return nil
}

func resourceVPSVncRead(d *schema.ResourceData, meta interface{}) error {
	// The OVH API exposes no GET for a previously opened VNC session: the POST
	// returns one-shot credentials that the operator is expected to use
	// immediately. We therefore keep whatever the create call stored in state
	// and treat the credentials as opaque. Any change to protocol or
	// service_name forces a new resource which re-issues fresh credentials.
	return nil
}

func resourceVPSVncDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no close endpoint; the VNC session auto-expires server-side.
	// Removing the resource from state is sufficient.
	d.SetId("")
	return nil
}
