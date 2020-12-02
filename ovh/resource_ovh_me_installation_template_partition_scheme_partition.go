package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceMeInstallationTemplatePartitionSchemePartition() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeInstallationTemplatePartitionSchemePartitionCreate,
		Read:   resourceMeInstallationTemplatePartitionSchemePartitionRead,
		Update: resourceMeInstallationTemplatePartitionSchemePartitionUpdate,
		Delete: resourceMeInstallationTemplatePartitionSchemePartitionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMeInstallationTemplatePartitionSchemePartitionImportState,
		},

		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Template name",
			},
			"scheme_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of this partitioning scheme",
			},
			"mountpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition mount point",
			},
			"filesystem": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Partition filesystem",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateFilesystem(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "size of partition in MB, 0 => rest of the space",
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "step or order. specifies the creation order of the partition on the disk",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition type",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidatePartitionType(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			// optional
			"raid": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "raid partition type",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidatePartitionRAIDMode(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The volume name needed for proxmox distribution",
			},
		},
	}
}

func resourceMeInstallationTemplatePartitionSchemePartitionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("Import Id is not template_name/scheme_name/mountpoint formatted")
	}
	templateName := splitId[0]
	schemeName := splitId[1]
	mountpoint := splitId[2]
	d.Set("template_name", templateName)
	d.Set("scheme_name", schemeName)
	d.Set("mountpoint", mountpoint)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceMeInstallationTemplatePartitionSchemePartitionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)

	opts := (&PartitionCreateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/partition",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
	)

	if err := config.OVHClient.Post(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", templateName, schemeName, opts.Mountpoint))

	return resourceMeInstallationTemplatePartitionSchemePartitionRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemePartitionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)

	opts := (&PartitionUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/partition/%s",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
		url.PathEscape(opts.Mountpoint),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Calling PUT %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	return resourceMeInstallationTemplatePartitionSchemePartitionRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemePartitionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)
	mountpoint := d.Get("mountpoint").(string)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/partition/%s",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
		url.PathEscape(mountpoint),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Calling DELETE %s: %s \n", endpoint, err.Error())
	}

	return nil
}

func resourceMeInstallationTemplatePartitionSchemePartitionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)
	mountpoint := d.Get("mountpoint").(string)

	r := &Partition{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/partition/%s",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
		url.PathEscape(mountpoint),
	)

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", templateName, schemeName, mountpoint))

	return nil
}
