package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectUserS3Credential() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectUserS3CredentialCreate,
		Read:   resourceCloudProjectUserS3CredentialRead,
		Delete: resourceCloudProjectUserS3CredentialDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudProjectUserS3CredentialImportState,
		},

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
			"internal_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCloudProjectUserS3CredentialImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	access := d.Id()
	serviceName := ""
	userId := ""

	splitId := strings.SplitN(access, "/", 3)
	if len(splitId) == 3 {
		serviceName = splitId[0]
		userId = splitId[1]
		access = splitId[2]
	}

	if serviceName == "" || userId == "" || access == "" {
		return nil, fmt.Errorf("Import Id is not service_name/user_id/access_key formatted")
	}

	d.SetId(access)
	d.Set("service_name", serviceName)
	d.Set("user_id", userId)

	return []*schema.ResourceData{d}, nil
}

func resourceCloudProjectUserS3CredentialCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userID := d.Get("user_id").(string)

	s3Credential := &CloudProjectUserS3Credential{}

	log.Printf("[DEBUG] Will create Public Cloud S3 AccessKey for user: %s from project: %s", userID, serviceName)
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s/s3Credentials",
		url.PathEscape(serviceName),
		url.PathEscape(userID),
	)

	if err := config.OVHClient.Post(endpoint, nil, s3Credential); err != nil {
		return fmt.Errorf("calling Post %s :\n\t %q", endpoint, err)
	}

	d.SetId(s3Credential.Access)
	d.Set("secret_access_key", s3Credential.Secret)
	for k, v := range s3Credential.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Created the Public Cloud S3 AccessKey %s", s3Credential)

	return resourceCloudProjectUserS3CredentialRead(d, meta)
}

func resourceCloudProjectUserS3CredentialRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userID := d.Get("user_id").(string)

	s3Credential := &CloudProjectUserS3Credential{}

	log.Printf("[DEBUG] Will read the Public Cloud S3 AccessKey %s for user %s from project: %s", d.Id(), userID, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s/s3Credentials/%s",
		url.PathEscape(serviceName),
		userID,
		d.Id(),
	)

	if err := config.OVHClient.Get(endpoint, s3Credential); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(s3Credential.Access)
	// set resource attributes
	for k, v := range s3Credential.ToMap() {
		d.Set(k, v)
	}

	log.Printf("[DEBUG] Read the Public Cloud S3 AccessKey %s", s3Credential)
	return nil
}

func resourceCloudProjectUserS3CredentialDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userID := d.Get("user_id").(string)

	id := d.Id()

	log.Printf("[DEBUG] Will delete the Public Cloud S3 AccessKey %s for user %s from project: %s", id, userID, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/user/%s/s3Credentials/%s",
		url.PathEscape(serviceName),
		userID,
		id,
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("calling Delete %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Deleted the Public Cloud S3 AccessKey %s for user %s from project %s", id, userID, serviceName)

	d.SetId("")

	return nil
}
