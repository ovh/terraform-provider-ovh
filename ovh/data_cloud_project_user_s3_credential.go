package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataCloudProjectUserS3Credential() *schema.Resource {
	return &schema.Resource{
		Read: dataCloudProjectUserS3CredentialRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
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
			"access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The access key",
			},

			//Computed
			"secret_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataCloudProjectUserS3CredentialRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userID := d.Get("user_id").(string)
	accessKey := d.Get("access_key_id").(string)

	log.Printf("[DEBUG] Will read public cloud secret access key for access key %s user %s on project: %s", accessKey, userID, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s/s3Credentials/%s",
		url.PathEscape(serviceName),
		url.PathEscape(userID),
		url.PathEscape(accessKey),
	)

	s3Credential := &CloudProjectUserS3Credential{}
	if err := config.OVHClient.Get(endpoint, &s3Credential); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(serviceName)
	d.Set("secret_access_key", s3Credential.Secret)

	return nil
}
