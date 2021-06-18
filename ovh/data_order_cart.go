package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

const (
	orderCartExpireFormat = "2006-01-02T15:04:05+00:00"
)

func dataSourceOrderCart() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrderCartRead,

		Schema: map[string]*schema.Schema{
			"ovh_subsidiary": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Ovh Subsidiary",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateSubsidiary(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of your cart",
				Optional:    true,
				ForceNew:    true,
			},
			"expire": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Expiration time (format: %s)", orderCartExpireFormat),
				DefaultFunc: func() (interface{}, error) {
					// by default cart expires in 30 minute as it should
					// be sufficient to retrieve info from api and proceed to order
					now := time.Now().Add(time.Minute * 30)
					return now.Format(orderCartExpireFormat), nil
				},
				Optional: true,
				Computed: true,
			},

			// Computed
			"cart_id": {
				Type:        schema.TypeString,
				Description: "Cart identifier",
				Computed:    true,
			},

			"read_only": {
				Description: "Indicates if the cart has already been validated",
				Computed:    true,
				Type:        schema.TypeBool,
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Items of your cart",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceOrderCartRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var r *OrderCart

	// to be created
	if d.Id() == "" {
		params := (&OrderCartCreateOpts{}).FromResource(d)
		newCart, err := orderCartCreate(meta, params, false)
		if err != nil {
			return err
		}
		r = newCart
		d.SetId(r.CartId)

	} else {
		r = &OrderCart{}
		log.Printf("[DEBUG] Will get order cart: %v", d.Id())
		endpoint := fmt.Sprintf(
			"/order/cart/%s",
			url.PathEscape(d.Id()),
		)

		err := config.OVHClient.Get(endpoint, r)
		if err != nil {
			return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
		}

	}
	// set attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	return nil
}
