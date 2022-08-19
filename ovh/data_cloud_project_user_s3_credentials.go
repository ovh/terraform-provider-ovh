package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataCloudProjectUserS3Credentials() *schema.Resource {
	return &schema.Resource{
		Read: dataCloudProjectUserS3CredentialsRead,
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
			"access_key_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataCloudProjectUserS3CredentialsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userID := d.Get("user_id").(string)

	log.Printf("[DEBUG] Will read public cloud access key ids for user %s on project: %s", userID, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s/s3Credentials",
		url.PathEscape(serviceName),
		url.PathEscape(userID),
	)

	credentials := make([]CloudProjectUserS3Credential, 0)
	if err := config.OVHClient.Get(endpoint, &credentials); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	accessKeys := make([]string, 0, len(credentials))

	for _, key := range credentials {
		accessKeys = append(accessKeys, key.Access)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, userID))
	d.Set("access_key_ids", accessKeys)

	return nil
}
