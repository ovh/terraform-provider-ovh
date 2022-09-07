package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type CloudProjectUserS3Policy struct {
	Policy string `json:"policy"`
}

const defaultS3Policy = `{"Statement":[]}`

func resourceCloudProjectUserS3Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectUserS3PolicyUpdate,
		UpdateContext: resourceCloudProjectUserS3PolicyUpdate,
		ReadContext:   resourceCloudProjectUserS3PolicyRead,
		DeleteContext: resourceCloudProjectUserS3PolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudProjectUserS3PolicyImportState,
		},

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
				Required:    true,
				Description: "The policy document. This is a JSON formatted string.",
			},
		},
	}
}

func resourceCloudProjectUserS3PolicyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	serviceName := ""
	userId := d.Id()

	splitId := strings.SplitN(userId, "/", 2)
	if len(splitId) == 2 {
		serviceName = splitId[0]
		userId = splitId[1]
	}

	if serviceName == "" || userId == "" {
		return nil, fmt.Errorf("Import ID is not service_name/user_id formatted. Got %s", d.Id())
	}

	d.SetId(buildUserS3PolicyId(serviceName, userId))
	d.Set("service_name", serviceName)
	d.Set("user_id", userId)

	return []*schema.ResourceData{d}, nil
}

func resourceCloudProjectUserS3PolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userId := d.Get("user_id").(string)
	policy := d.Get("policy").(string)

	log.Printf("[DEBUG] Will set the S3 Policy for user: %s from project: %s", userId, serviceName)
	if err := postPolicy(serviceName, userId, policy, config); err != nil {
		return err
	}

	d.SetId(buildUserS3PolicyId(serviceName, userId))

	log.Printf("[DEBUG] Policy %s set for user %s", policy, userId)

	return resourceCloudProjectUserS3PolicyRead(ctx, d, meta)
}

func resourceCloudProjectUserS3PolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	log.Printf("[DEBUG] Policy %s set for user %s", s3Policy.Policy, userId)
	return nil
}

func resourceCloudProjectUserS3PolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userId := d.Get("user_id").(string)

	log.Printf("[DEBUG] Will delete the S3 Policy for user: %s from project: %s", userId, serviceName)

	if err := postPolicy(serviceName, userId, defaultS3Policy, config); err != nil {
		return err
	}

	log.Printf("[DEBUG] Policy deleted for user %s", userId)

	d.SetId("")

	return nil
}

func postPolicy(serviceName string, userID string, policy string, config *Config) diag.Diagnostics {
	s3PolicyReq := &CloudProjectUserS3Policy{Policy: policy}
	endpoint := buildUserS3PolicyEndpoint(serviceName, userID)

	if err := config.OVHClient.Post(endpoint, s3PolicyReq, nil); err != nil {
		return diag.Errorf("calling Post %s :\n\t %q", endpoint, err)
	}
	return nil
}

func buildUserS3PolicyId(serviceName string, userId string) string {
	return fmt.Sprintf(
		"%s/%s",
		url.PathEscape(serviceName),
		url.PathEscape(userId),
	)
}

func buildUserS3PolicyEndpoint(serviceName string, userId string) string {
	return fmt.Sprintf(
		"/cloud/project/%s/user/%s/policy",
		url.PathEscape(serviceName),
		url.PathEscape(userId),
	)
}
