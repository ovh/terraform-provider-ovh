package ovh

import (
	"context"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceIamReferenceResourceType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIamReferenceResourceTypeRead,
		Schema: map[string]*schema.Schema{
			"types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceIamReferenceResourceTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	e, err := url.Parse("/v2/iam/reference/resource/type")
	if err != nil {
		return diag.Errorf("Unable to get actions:\n\t %q", err)
	}

	var types []string
	err = config.OVHClient.GetWithContext(
		ctx,
		e.String(),
		&types,
	)
	if err != nil {
		return diag.Errorf("Unable to get actions:\n\t %q", err)
	}
	log.Printf("[DEBUG] resource types: %+v", types)

	d.SetId(hashcode.Strings(types))
	d.Set("types", types)

	return nil
}
