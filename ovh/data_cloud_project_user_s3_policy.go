package ovh

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataCloudProjectUserS3Policy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCloudProjectUserS3PolicyRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the ID of the cloud project.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user ID",
			},
			"policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy document. This is a JSON formatted string.",
			},
		},
	}
}

func dataCloudProjectUserS3PolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userId := d.Get("user_id").(string)

	s3Policy := &CloudProjectUserS3Policy{}

	log.Printf("[DEBUG] Will read the S3 Policy for user: %s from project: %s", userId, serviceName)

	endpoint := buildUserS3PolicyEndpoint(serviceName, userId)
	if err := config.OVHClient.Get(endpoint, s3Policy); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.SetId(buildUserS3PolicyId(serviceName, userId))
	d.Set("policy", s3Policy.Policy)

	return nil
}
