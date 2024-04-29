package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"

	"github.com/ovh/go-ovh/ovh"
)

func resourceMeInstallationTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeInstallationTemplateCreate,
		Read:   resourceMeInstallationTemplateRead,
		Update: resourceMeInstallationTemplateUpdate,
		Delete: resourceMeInstallationTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMeInstallationTemplateImportState,
		},

		Schema: map[string]*schema.Schema{
			"base_template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH template name yours will be based on, choose one among the list given by compatibleTemplates function",
			},
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This template name",
			},
			"remove_default_partition_schemes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Remove default partition schemes at creation",
			},
			"customization": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Set up the server using the provided hostname instead of the default hostname",
						},
						"post_installation_script_link": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicate the URL where your postinstall customisation script is located",
						},
						"post_installation_script_return": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'",
						},
					},
				},
			},

			"bit_format": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This template bit format (32 or 64)",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "information about this template",
			},
			"distribution": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the distribution this template is based on",
			},
			"family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "this template family type",
			},
			"hard_raid_configuration": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports hardware raid configuration through the OVH API",
				Deprecated:  "This will be deprecated in the next release",
			},
			"filesystems": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Filesystems available",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lvm_ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports Logical Volumes (Linux LVM)",
			},
		},
	}
}

func resourceMeInstallationTemplateImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not base_template_name/template_name formatted")
	}
	baseTemplateName := splitId[0]
	templateName := splitId[1]
	d.SetId(templateName)
	d.Set("base_template_name", baseTemplateName)
	d.Set("template_name", templateName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceMeInstallationTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	opts := (&InstallationTemplateCreateOpts{}).FromResource(d)

	endpoint := "/me/installationTemplate"

	// the resource is created via the POST endpoint, then updated
	// via the PUT endpoint to apply customizations.
	if err := config.OVHClient.Post(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	d.SetId(d.Get("template_name").(string))

	// We call the update method to put customization opts
	updateOpts := (&InstallationTemplateUpdateOpts{}).FromResource(d)
	endpoint = fmt.Sprintf(
		"/me/installationTemplate/%s",
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Put(endpoint, updateOpts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	// handle remove_default_partitions option
	removeDefaultPartitions := false

	if v, ok := d.GetOk("remove_default_partition_schemes"); ok {
		removeDefaultPartitions = v.(bool)
	}

	d.Set("remove_default_partition_schemes", removeDefaultPartitions)

	if removeDefaultPartitions {
		templateName := d.Get("template_name").(string)
		defaultSchemes, err := getPartitionSchemeIds(templateName, config.OVHClient)
		if err != nil {
			return err
		}

		for _, scheme := range defaultSchemes {
			endpoint := fmt.Sprintf(
				"/me/installationTemplate/%s/partitionScheme/%s",
				url.PathEscape(templateName),
				url.PathEscape(scheme),
			)

			if err := config.OVHClient.Delete(endpoint, nil); err != nil {
				return fmt.Errorf("Error calling DELETE %s: %s \n", endpoint, err.Error())
			}
		}
	}

	return resourceMeInstallationTemplateRead(d, meta)
}

func resourceMeInstallationTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	opts := (&InstallationTemplateUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s",
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	return resourceMeInstallationTemplateRead(d, meta)
}

func resourceMeInstallationTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("template_name").(string)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s",
		url.PathEscape(name),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Error calling DELETE %s: %s \n", endpoint, err.Error())
	}

	return nil
}

func resourceMeInstallationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	template, err := getInstallationTemplate(d, config.OVHClient)
	if err != nil {
		return err
	}

	// set attributes
	for k, v := range template.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func getInstallationTemplate(d *schema.ResourceData, client *ovh.Client) (*InstallationTemplate, error) {
	r := &InstallationTemplate{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s",
		url.PathEscape(d.Get("template_name").(string)),
	)

	if err := client.Get(endpoint, r); err != nil {
		return nil, helpers.CheckDeleted(d, err, endpoint)
	}

	return r, nil
}
