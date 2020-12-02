package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceMeInstallationTemplatePartitionScheme() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeInstallationTemplatePartitionSchemeCreate,
		Read:   resourceMeInstallationTemplatePartitionSchemeRead,
		Update: resourceMeInstallationTemplatePartitionSchemeUpdate,
		Delete: resourceMeInstallationTemplatePartitionSchemeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMeInstallationTemplatePartitionSchemeImportState,
		},

		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This template name",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of this partitioning scheme",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "on a reinstall, if a partitioning scheme is not specified, the one with the higher priority will be used by default, among all the compatible partitioning schemes (given the underlying hardware specifications)",
			},
		},
	}
}

func resourceMeInstallationTemplatePartitionSchemeImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not template_name/scheme_name formatted")
	}
	templateName := splitId[0]
	schemeName := splitId[1]
	d.Set("template_name", templateName)
	d.Set("name", schemeName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceMeInstallationTemplatePartitionSchemeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)

	opts := (&PartitionSchemeCreateOrUpdateOpts{}).FromResource(d)
	endpoint := fmt.Sprintf("/me/installationTemplate/%s/partitionScheme", templateName)

	if err := config.OVHClient.Post(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	d.SetId(fmt.Sprintf(
		"%s/%s",
		url.PathEscape(templateName),
		url.PathEscape(opts.Name),
	))

	return resourceMeInstallationTemplatePartitionSchemeRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)

	opts := (&PartitionSchemeCreateOrUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s",
		url.PathEscape(templateName),
		url.PathEscape(opts.Name),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	return resourceMeInstallationTemplatePartitionSchemeRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	name := d.Get("name").(string)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s",
		url.PathEscape(templateName),
		url.PathEscape(name),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Error calling DELETE %s: %s \n", endpoint, err.Error())
	}

	return nil
}

func resourceMeInstallationTemplatePartitionSchemeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	templateName := d.Get("template_name").(string)
	name := d.Get("name").(string)

	r := &PartitionScheme{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s",
		url.PathEscape(templateName),
		url.PathEscape(name),
	)

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprintf("%s-%s", templateName, name))
	return nil
}
