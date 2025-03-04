package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func datasourceCloudProjectUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCloudProjectUserRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user ID",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceCloudProjectUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userId := d.Get("user_id").(string)

	user := &CloudProjectUser{}

	log.Printf("[DEBUG] Will read public cloud user %s from project: %s", d.Id(), serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(userId),
	)

	if err := config.OVHClient.Get(endpoint, user); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.SetId(strconv.Itoa(user.Id))
	// set resource attributes
	for k, v := range user.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read Public Cloud User %s", user)
	return nil
}
