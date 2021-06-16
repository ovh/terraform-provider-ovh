package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceDedicatedServerBoots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedServerBootsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your dedicated server.",
				Required:    true,
			},

			"boot_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the value of bootType property",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateBootType(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			"kernel": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the value of kernel property",
			},

			// Computed
			"result": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Server compatibles netboots ids",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceDedicatedServerBootsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	ids := []int64{}

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/boot",
		url.PathEscape(serviceName),
	)

	if bootType, ok := d.GetOk("boot_type"); ok {
		endpoint = fmt.Sprintf(
			"%s?bootType=%s",
			endpoint,
			url.PathEscape(bootType.(string)),
		)
	}

	if err := config.OVHClient.Get(endpoint, &ids); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	if kernel, ok := d.GetOk("kernel"); ok {
		// refine boots on kernel name
		new_ids := []int64{}
		for _, id := range ids {
			boot := &DedicatedServerBoot{}
			endpoint := fmt.Sprintf(
				"/dedicated/server/%s/boot/%d",
				url.PathEscape(serviceName),
				id,
			)
			if err := config.OVHClient.Get(endpoint, boot); err != nil {
				return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
			}

			if boot.Kernel == kernel {
				new_ids = append(new_ids, id)
			}
		}
		ids = new_ids
	}

	// setting id by computing a hashcode of all the ids
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.FormatInt(id, 10)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(idsStr)

	d.SetId(hashcode.Strings(idsStr))
	d.Set("result", ids)
	return nil
}
