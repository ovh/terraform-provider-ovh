package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceDedicatedInstallationTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedInstallationTemplateRead,
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This template name",
			},

			// Computed properties
			"bit_format": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This template bit format (32 or 64)",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Category of this template (informative only)",
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
			"end_of_install": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "after this date, install of this template will not be possible at OVH",
			},
			"family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "this template family type (bsd,linux,solaris,windows)",
			},
			"filesystems": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Filesystems available",
			},
			"hard_raid_configuration": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports hardware raid configuration through the OVH API",
			},
			"inputs": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mandatory": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enum": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"license": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"usage": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"lvm_ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports Logical Volumes (Linux LVM)",
			},
			"no_partitioning": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The template supports RAID0 and RAID1 on 2 disks.",
			},
			"project": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "This distribution project details",
				Elem: &schema.Resource{
					Description: "A structure describing some interesting facts about an OS template project",
					Schema: map[string]*schema.Schema{
						"usage": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item version",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item url",
									},
									"release_notes": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item release notes",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item name",
									},
									"governance": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "OS template project item governance",
									},
								},
							},
						},
						"os": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "OS template project OS details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item version",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item url",
									},
									"release_notes": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item release notes",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS template project item name",
									},
									"governance": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "OS template project item governance",
									},
								},
							},
						},
					},
				},
			},
			"soft_raid_only_mirroring": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Partitioning customization is available but limited to mirroring for this OS template",
			},
			"subfamily": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "this template subfamily type",
			},
		},
	}
}

func dataSourceDedicatedInstallationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	template, err := getDedicatedInstallationTemplate(d, config.OVHClient)
	if err != nil {
		return err
	}
	if template == nil {
		return fmt.Errorf("Your query returned no results. Please change your search criteria")
	}

	// set attributes
	for k, v := range template.ToMap() {
		d.Set(k, v)
	}

	name := d.Get("template_name").(string)
	d.SetId(name)

	return nil
}

func getDedicatedInstallationTemplate(d *schema.ResourceData, client *ovh.Client) (*InstallationTemplate, error) {
	r := &InstallationTemplate{}

	endpoint := fmt.Sprintf(
		"/dedicated/installationTemplate/%s",
		url.PathEscape(d.Get("template_name").(string)),
	)

	if err := client.Get(endpoint, r); err != nil {
		return nil, helpers.CheckDeleted(d, err, endpoint)
	}

	return r, nil
}
