package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func datasourceCloudProjectUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCloudProjectUsersRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			// Computed
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
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
					},
				},
			},
		},
	}
}

func datasourceCloudProjectUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	users := make([]CloudProjectUser, 0)

	log.Printf("[DEBUG] Will read public cloud user %s from project: %s", d.Id(), serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Get(endpoint, &users); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	mapUsers := make([]map[string]interface{}, len(users))
	ids := make([]string, len(users))

	for i, user := range users {
		mapUsers[i] = user.ToMap()
		mapUsers[i]["user_id"] = strconv.Itoa(user.Id)
		ids = append(ids, strconv.Itoa(user.Id))
	}

	d.SetId(hashcode.Strings(ids))
	d.Set("users", mapUsers)

	log.Printf("[DEBUG] Read Public Cloud User %s", mapUsers)
	return nil
}
