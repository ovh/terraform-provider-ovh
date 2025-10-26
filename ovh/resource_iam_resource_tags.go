package ovh

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceIamResourceTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIamResourceTagsCreate,
		ReadContext:   resourceIamResourceTagsRead,
		UpdateContext: resourceIamResourceTagsUpdate,
		DeleteContext: resourceIamResourceTagsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIamResourceTagsImport,
		},
		Schema: map[string]*schema.Schema{
			"urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "URN of the resource",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Map of tags to apply to the resource. Keys must match ^[a-zA-Z0-9_.:/=+@-]{1,128}$ and values must match ^[a-zA-Z0-9_.:/=+@-]{0,256}$. Tags prefixed with 'ovh:' are managed by OVH and cannot be set",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateDiagFunc: validateTags,
				DiffSuppressFunc: suppressOvhTagsDiff,
			},
		},
	}
}

// validateTags validates tag keys and values according to the requirements
func validateTags(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	tags := v.(map[string]interface{})

	keyRegex := regexp.MustCompile(`^[a-zA-Z0-9_.:/=+@-]{1,128}$`)
	valueRegex := regexp.MustCompile(`^[a-zA-Z0-9_.:/=+@-]{0,256}$`)

	for key, val := range tags {
		// Check for ovh: prefix
		if strings.HasPrefix(key, "ovh:") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid tag key",
				Detail:   fmt.Sprintf("Tag key '%s' cannot start with 'ovh:' prefix. Tags with this prefix are managed by OVH and cannot be set via the API", key),
			})
		}

		if !keyRegex.MatchString(key) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid tag key",
				Detail:   fmt.Sprintf("Tag key '%s' must match pattern ^[a-zA-Z0-9_.:/=+@-]{1,128}$", key),
			})
		}

		value := val.(string)
		if !valueRegex.MatchString(value) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid tag value",
				Detail:   fmt.Sprintf("Tag value '%s' for key '%s' must match pattern ^[a-zA-Z0-9_.:/=+@-]{0,256}$", value, key),
			})
		}
	}

	return diags
}

// suppressOvhTagsDiff suppresses diffs for ovh: prefixed tags
// These tags are managed by OVH and should not cause drift detection
func suppressOvhTagsDiff(k, old, new string, d *schema.ResourceData) bool {
	// k will be in format "tags.key_name" or just "tags" for the whole map
	// We only want to suppress individual tag keys that start with "ovh:"
	if strings.HasPrefix(k, "tags.") {
		tagKey := strings.TrimPrefix(k, "tags.")
		// Suppress diff if this is an ovh: prefixed tag
		if strings.HasPrefix(tagKey, "ovh:") {
			return true
		}
	}
	return false
}

func resourceIamResourceTagsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resourceURN := d.Get("urn").(string)
	endpoint := fmt.Sprintf("/v2/iam/resource/%s", url.PathEscape(resourceURN))

	// Get tags from configuration and convert to map[string]string
	configTags := d.Get("tags").(map[string]interface{})
	tagsMap := make(map[string]string)
	for key, val := range configTags {
		tagsMap[key] = val.(string)
	}

	// Update resource with tags using PUT endpoint
	// The API will preserve ovh: prefixed system tags automatically
	updateBody := map[string]interface{}{
		"tags": tagsMap,
	}

	err := config.OVHClient.PutWithContext(ctx, endpoint, updateBody, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create tags: %w", err))
	}

	// Set the resource ID to the URN
	d.SetId(resourceURN)

	// Read back the resource to sync state
	return resourceIamResourceTagsRead(ctx, d, meta)
}

func resourceIamResourceTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resourceURN := d.Id()

	// Fetch the resource details
	var resourceDetails IamResourceDetails
	endpoint := fmt.Sprintf("/v2/iam/resource/%s", url.PathEscape(resourceURN))
	err := config.OVHClient.GetWithContext(ctx, endpoint, &resourceDetails)
	if err != nil {
		return diag.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	// Sync all tags to state (including ovh: prefixed system tags)
	// This allows users to see all tags, but ovh: prefixed tags cannot be modified
	allTags := make(map[string]string)
	for key, value := range resourceDetails.Tags {
		allTags[key] = value
	}

	if err := d.Set("tags", allTags); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIamResourceTagsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// The ID passed to import is the resource URN
	resourceURN := d.Id()

	// Set the urn field from the import ID
	if err := d.Set("urn", resourceURN); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceIamResourceTagsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resourceURN := d.Id()

	if d.HasChange("tags") {
		endpoint := fmt.Sprintf("/v2/iam/resource/%s", url.PathEscape(resourceURN))

		// // Get tags from configuration
		configTags := d.Get("tags").(map[string]interface{})

		// Update resource with tags using PUT endpoint
		// The API will preserve ovh: prefixed system tags automatically
		updateBody := map[string]interface{}{
			"tags": configTags,
		}

		err := config.OVHClient.PutWithContext(ctx, endpoint, updateBody, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update tags: %w", err))
		}
	}

	// Read back the resource to sync state
	return resourceIamResourceTagsRead(ctx, d, meta)
}

func resourceIamResourceTagsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resourceURN := d.Id()
	endpoint := fmt.Sprintf("/v2/iam/resource/%s", url.PathEscape(resourceURN))

	// Send empty tags to remove all managed tags
	// The API will preserve ovh: prefixed system tags automatically
	updateBody := map[string]interface{}{
		"tags": map[string]string{},
	}

	err := config.OVHClient.PutWithContext(ctx, endpoint, updateBody, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.SetId("")
	return nil
}
