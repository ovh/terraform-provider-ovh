package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceMeInstallationTemplatePartitionSchemeHardwareRaid() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeInstallationTemplatePartitionSchemeHardwareRaidCreate,
		Read:   resourceMeInstallationTemplatePartitionSchemeHardwareRaidRead,
		Update: resourceMeInstallationTemplatePartitionSchemeHardwareRaidUpdate,
		Delete: resourceMeInstallationTemplatePartitionSchemeHardwareRaidDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMeInstallationTemplatePartitionSchemeHardwareRaidImportState,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hardware RAID name",
			},
			"disks": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Disk List. Syntax is cX:dY for disks and [cX:dY,cX:dY] for groups. With X and Y resp. the controller id and the disk id",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "RAID mode (raid0, raid1, raid10, raid5, raid50, raid6, raid60)",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateRAIDMode(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"step": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the creation order of the hardware RAID",
			},
		},
	}
}

func resourceMeInstallationTemplatePartitionSchemeHardwareRaidImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("Import Id is not template_name/scheme_name/name formatted")
	}
	templateName := splitId[0]
	schemeName := splitId[1]
	name := splitId[2]
	d.Set("template_name", templateName)
	d.Set("scheme_name", schemeName)
	d.Set("name", name)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceMeInstallationTemplatePartitionSchemeHardwareRaidCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)

	opts := (&HardwareRaidCreateOrUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/hardwareRaid",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
	)

	if err := config.OVHClient.Post(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Calling %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", templateName, schemeName, opts.Name))

	return resourceMeInstallationTemplatePartitionSchemeHardwareRaidRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemeHardwareRaidUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)
	name := d.Get("name").(string)

	opts := (&HardwareRaidCreateOrUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/hardwareRaid/%s",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
		url.PathEscape(name),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Calling %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	return resourceMeInstallationTemplatePartitionSchemeHardwareRaidRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemeHardwareRaidDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)
	name := d.Get("name").(string)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/hardwareRaid/%s",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
		url.PathEscape(name),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	return nil
}

func resourceMeInstallationTemplatePartitionSchemeHardwareRaidRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	schemeName := d.Get("scheme_name").(string)
	name := d.Get("name").(string)

	r := &HardwareRaid{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/hardwareRaid/%s",
		url.PathEscape(templateName),
		url.PathEscape(schemeName),
		url.PathEscape(name),
	)

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", templateName, schemeName, name))

	return nil
}
