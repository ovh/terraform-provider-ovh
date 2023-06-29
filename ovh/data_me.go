package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMe() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMeRead,
		Schema: map[string]*schema.Schema{
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity URN of the account",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address of nichandle",
			},
			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Area of nichandle",
			},
			"birth_city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City of birth",
			},
			"birth_day": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Birth date",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City of nichandle",
			},
			"company_national_identification_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Company National Identification Number",
			},
			"corporation_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Corporation type",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Customer country",
			},
			"currency": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Currency code",
						},
						"symbol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Currency symbol",
						},
					},
				},
				Description: "Customer currency",
			},
			"customer_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Your customer code (a numerical value used for identification when contacting support via phone call)",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address",
			},
			"fax": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fax number",
			},
			"firstname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "First name",
			},
			"italian_sdi": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Italian SDI",
			},
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Language",
			},
			"legalform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Customer legal form",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Customer name",
			},
			"national_identification_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "National Identification Number",
			},
			"nichandle": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Customer identifier",
			},
			"organisation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of organisation",
			},
			"ovh_company": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OVH subsidiary",
			},
			"ovh_subsidiary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OVH subsidiary",
			},
			"phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Phone number",
			},
			"phone_country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Phone number's country code",
			},
			"sex": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gender",
			},
			"spare_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Spare email",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Nichandle state",
			},
			"vat": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VAT number",
			},
			"zip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zipcode",
			},
		},
	}
}

func dataSourceMeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	me := &MeResponse{}

	err := config.OVHClient.Get("/me", me)
	if err != nil {
		return fmt.Errorf("Unable to retrieve /me information:\n\t %q", err)
	}
	log.Printf("[DEBUG] /me information: %+v", me)

	d.SetId(me.Nichandle)

	for k, v := range me.ToMap() {
		d.Set(k, v)
	}

	d.Set("urn", fmt.Sprintf("urn:v1:%s:resource:account:%s", config.Plate, me.Nichandle))

	return nil
}
