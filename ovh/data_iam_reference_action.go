package ovh

import (
	"context"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIamReferenceActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIamReferenceActionRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"categories": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceIamReferenceActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	resType := d.Get("type").(string)

	e, err := url.Parse("/v2/iam/reference/action")
	if err != nil {
		return diag.Errorf("Unable to get actions:\n\t %q", err)
	}
	q := e.Query()
	q.Add("resourceType", resType)
	e.RawQuery = q.Encode()
	log.Printf("[DEBUG] actions url: %+s", e.String())

	var actions []IamReferenceAction
	err = config.OVHClient.GetWithContext(
		ctx,
		e.String(),
		&actions,
	)
	if err != nil {
		return diag.Errorf("Unable to get actions:\n\t %q", err)
	}
	log.Printf("[DEBUG] actions: %+v", actions)

	saveAct := make([]map[string]any, 0, len(actions))
	for _, a := range actions {
		saveAct = append(saveAct, a.ToMap())
	}
	d.SetId(resType)
	d.Set("actions", saveAct)

	return nil
}
