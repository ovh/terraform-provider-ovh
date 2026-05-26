package ovh

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSServiceInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSServiceInfoRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
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
			"contact_admin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_billing": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_tech": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"can_delete_at_expiration": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"possible_renew_period": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"renew_automatic": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"renew_delete_at_expiration": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"renew_forced": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"renew_manual_payment": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSServiceInfoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	info := &VPSServiceInfo{}
	endpoint := fmt.Sprintf("/vps/%s/serviceInfos", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, info); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.SetId(strconv.FormatInt(info.ServiceID, 10))
	return vpsServiceInfoApplyToState(d, info)
}

func vpsServiceInfoApplyToState(d *schema.ResourceData, info *VPSServiceInfo) error {
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
	d.Set("contact_admin", info.ContactAdmin)
	d.Set("contact_billing", info.ContactBilling)
	d.Set("contact_tech", info.ContactTech)
	d.Set("domain", info.Domain)
	d.Set("can_delete_at_expiration", info.CanDeleteAtExpiration)
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
	return nil
}
