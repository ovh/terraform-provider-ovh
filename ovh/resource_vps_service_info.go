package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPSServiceInfo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPSServiceInfoCreateOrUpdate,
		UpdateContext: resourceVPSServiceInfoCreateOrUpdate,
		ReadContext:   resourceVPSServiceInfoRead,
		DeleteContext: resourceVPSServiceInfoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service_name of the VPS",
				Required:    true,
				ForceNew:    true,
			},
			"renew_automatic": {
				Type:        schema.TypeBool,
				Description: "Whether automatic renewal is enabled",
				Required:    true,
			},
			"renew_delete_at_expiration": {
				Type:        schema.TypeBool,
				Description: "Whether the service should be deleted at expiration",
				Optional:    true,
				Default:     false,
			},
			"renew_forced": {
				Type:        schema.TypeBool,
				Description: "Whether renewal is forced",
				Optional:    true,
				Default:     false,
			},
			"renew_manual_payment": {
				Type:        schema.TypeBool,
				Description: "Whether renewal requires manual payment",
				Optional:    true,
			},
			"renew_period": {
				Type:        schema.TypeInt,
				Description: "Renewal period in months",
				Optional:    true,
			},
			// Computed
			"service_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engaged_up_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renewal_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"possible_renew_period": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func vpsServiceInfoFetch(config *Config, serviceName string) (*VPSServiceInfo, error) {
	info := &VPSServiceInfo{}
	endpoint := fmt.Sprintf("/vps/%s/serviceInfos", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, info); err != nil {
		return nil, fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}
	return info, nil
}

func resourceVPSServiceInfoCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	// Fetch the current service infos to preserve non-renew fields on PUT.
	info, err := vpsServiceInfoFetch(config, serviceName)
	if err != nil {
		return diag.FromErr(err)
	}

	info.Renew.Automatic = d.Get("renew_automatic").(bool)
	info.Renew.DeleteAtExpiration = d.Get("renew_delete_at_expiration").(bool)
	info.Renew.Forced = d.Get("renew_forced").(bool)

	if v, ok := d.GetOkExists("renew_manual_payment"); ok {
		b := v.(bool)
		info.Renew.ManualPayment = &b
	} else {
		info.Renew.ManualPayment = nil
	}

	if v, ok := d.GetOk("renew_period"); ok {
		p := v.(int)
		info.Renew.Period = &p
	} else {
		info.Renew.Period = nil
	}

	endpoint := fmt.Sprintf("/vps/%s/serviceInfos", url.PathEscape(serviceName))
	if err := config.OVHClient.Put(endpoint, info, nil); err != nil {
		return diag.Errorf("Error calling PUT %s:\n\t %q", endpoint, err)
	}

	d.SetId(serviceName)
	return resourceVPSServiceInfoRead(ctx, d, meta)
}

func resourceVPSServiceInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	if serviceName == "" {
		serviceName = d.Id()
		d.Set("service_name", serviceName)
	}

	info, err := vpsServiceInfoFetch(config, serviceName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("service_id", info.ServiceID)
	d.Set("status", info.Status)
	d.Set("creation", info.Creation)
	d.Set("expiration", info.Expiration)
	if info.EngagedUpTo != nil {
		d.Set("engaged_up_to", *info.EngagedUpTo)
	} else {
		d.Set("engaged_up_to", "")
	}
	d.Set("renewal_type", info.RenewalType)
	d.Set("possible_renew_period", info.PossibleRenewPeriod)

	d.Set("renew_automatic", info.Renew.Automatic)
	d.Set("renew_delete_at_expiration", info.Renew.DeleteAtExpiration)
	d.Set("renew_forced", info.Renew.Forced)
	if info.Renew.ManualPayment != nil {
		d.Set("renew_manual_payment", *info.Renew.ManualPayment)
	}
	if info.Renew.Period != nil {
		d.Set("renew_period", *info.Renew.Period)
	}

	d.SetId(serviceName)
	return nil
}

func resourceVPSServiceInfoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not own the underlying VPS service. Removing it from
	// state restores the default renew configuration (automatic renewal off,
	// not forced, not scheduled for deletion at expiration) on the OVH side.
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	info, err := vpsServiceInfoFetch(config, serviceName)
	if err != nil {
		return diag.FromErr(err)
	}

	info.Renew.Automatic = false
	info.Renew.DeleteAtExpiration = false
	info.Renew.Forced = false
	info.Renew.ManualPayment = nil
	info.Renew.Period = nil

	endpoint := fmt.Sprintf("/vps/%s/serviceInfos", url.PathEscape(serviceName))
	if err := config.OVHClient.Put(endpoint, info, nil); err != nil {
		return diag.Errorf("Error calling PUT %s:\n\t %q", endpoint, err)
	}

	d.SetId("")
	return nil
}
