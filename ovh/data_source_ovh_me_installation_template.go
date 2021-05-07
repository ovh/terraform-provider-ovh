package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func dataSourceMeInstallationTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMeInstallationTemplateRead,
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This template name",
			},

			// computed
			"default_language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default language of this template",
			},
			"customization": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"change_log": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "field is not used anymore",
							Description: "Template change log details",
						},
						"custom_hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Set up the server using the provided hostname instead of the default hostname",
						},
						"post_installation_script_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate the URL where your postinstall customisation script is located",
						},
						"post_installation_script_return": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'",
						},
						"rating": {
							Type:       schema.TypeInt,
							Deprecated: "field is not used anymore",
							Computed:   true,
						},
						"ssh_key_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the ssh key that should be installed. Password login will be disabled",
						},
						"use_distribution_kernel": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Use the distribution's native kernel instead of the recommended OVH Kernel",
						},
					},
				},
			},

			"partition_scheme": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of this partitioning scheme",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "on a reinstall, if a partitioning scheme is not specified, the one with the higher priority will be used by default, among all the compatible partitioning schemes (given the underlying hardware specifications)",
						},
						"hardware_raid": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hardware RAID name",
									},
									"disks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Disk List. Syntax is cX:dY for disks and [cX:dY,cX:dY] for groups. With X and Y resp. the controller id and the disk id",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "RAID mode (raid0, raid1, raid10, raid5, raid50, raid6, raid60)",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the creation order of the hardware RAID",
									},
								},
							},
						},
						"partition": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filesystem": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Partition filesystem",
									},
									"mountpoint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "partition mount point",
									},
									"raid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "raid partition type",
									},
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "size of partition in MB, 0 => rest of the space",
									},
									"order": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "step or order. specifies the creation order of the partition on the disk",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "partition type",
									},
									"volume_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The volume name needed for proxmox distribution",
									},
								},
							},
						},
					},
				},
			},

			//Computed
			"available_languages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all language available for this template",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"beta": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution is new and, although tested and functional, may still display odd behaviour",
			},
			"bit_format": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This template bit format (32 or 64)",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Category of this template (informative only). (basic, customer, hosting, other, readyToUse, virtualisation)",
			},
			"deprecated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "is this distribution deprecated",
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
				Description: "this template family type (bsd,linux,solaris,windows)",
			},
			"hard_raid_configuration": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports hardware raid configuration through the OVH API",
			},
			"filesystems": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Filesystems available (btrfs,ext3,ext4,ntfs,reiserfs,swap,ufs,xfs,zfs)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_modification": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of last modification of the base image",
			},
			"lvm_ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports Logical Volumes (Linux LVM)",
			},
			"supports_distribution_kernel": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports installation using the distribution's native kernel instead of the recommended OVH kernel",
			},
			"supports_gpt_label": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports the GUID Partition Table (GPT), providing up to 128 partitions that can have more than 2TB",
			},
			"supports_rtm": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports RTM software",
			},
			"supports_sql_server": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This distribution supports the microsoft SQL server",
			},
			"supports_uefi": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This distribution supports UEFI setup (no,only,yes)",
			},
		},
	}
}

func dataSourceMeInstallationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	template, err := getInstallationTemplate(d, config.OVHClient)
	if err != nil {
		return err
	}

	// set attributes
	for k, v := range template.ToMap() {
		d.Set(k, v)
	}

	// set partitionSchemes
	err = partialMeInstallationTemplatePartitionSchemesRead(d, meta)
	if err != nil {
		return err
	}

	name := d.Get("template_name").(string)
	d.SetId(name)

	return nil
}

func partialMeInstallationTemplatePartitionSchemesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("template_name").(string)

	schemes, err := getPartitionSchemes(name, config.OVHClient)
	if err != nil {
		return err
	}

	partitionSchemes := make([]interface{}, len(schemes))

	for i, scheme := range schemes {
		partitionScheme := scheme.ToMap()

		// set partitionScheme Partitions
		partitions, err := getPartitionSchemePartitions(name, scheme.Name, config.OVHClient)
		if err != nil {
			return err
		}

		partitionList := make([]interface{}, len(partitions))
		for ii, partition := range partitions {
			partitionList[ii] = partition.ToMap()
		}
		partitionScheme["partition"] = partitionList

		// set partitionScheme HardwareRaids
		hardwareRaids, err := getPartitionSchemeHardwareRaids(name, scheme.Name, config.OVHClient)
		if err != nil {
			return err
		}

		hardwareRaidList := make([]interface{}, len(hardwareRaids))
		for ii, hardwareRaid := range hardwareRaids {
			hardwareRaidList[ii] = hardwareRaid.ToMap()
		}
		partitionScheme["hardware_raid"] = hardwareRaidList

		partitionSchemes[i] = partitionScheme
	}

	d.Set("partition_scheme", partitionSchemes)

	return nil
}

func getPartitionSchemes(template string, client *ovh.Client) ([]*PartitionScheme, error) {
	schemes, err := getPartitionSchemeIds(template, client)
	if err != nil {
		return nil, err
	}

	partitionSchemes := []*PartitionScheme{}
	for _, scheme := range schemes {
		partitionScheme, err := getPartitionScheme(template, scheme, client)
		if err != nil {
			return nil, err
		}

		partitionSchemes = append(partitionSchemes, partitionScheme)
	}

	return partitionSchemes, nil
}

func getPartitionScheme(template, scheme string, client *ovh.Client) (*PartitionScheme, error) {
	r := &PartitionScheme{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s",
		url.PathEscape(template),
		url.PathEscape(scheme),
	)

	if err := client.Get(endpoint, r); err != nil {
		return nil, fmt.Errorf("Error calling GET %s: %s \n", endpoint, err.Error())
	}

	return r, nil
}

func getPartitionSchemeIds(template string, client *ovh.Client) ([]string, error) {
	schemes := []string{}
	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme",
		url.PathEscape(template),
	)
	err := client.Get(endpoint, &schemes)

	if err != nil {
		return nil, fmt.Errorf("Error calling GET %s: %s \n", endpoint, err.Error())
	}
	return schemes, nil
}

func getPartitionSchemePartitions(template, scheme string, client *ovh.Client) ([]*Partition, error) {
	mountPoints := []string{}
	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/partition",
		url.PathEscape(template),
		url.PathEscape(scheme),
	)
	err := client.Get(endpoint, &mountPoints)

	if err != nil {
		return nil, fmt.Errorf("Error calling GET %s: %s \n", endpoint, err.Error())
	}

	partitions := []*Partition{}
	for _, mountPoint := range mountPoints {
		partition, err := getPartitionSchemePartition(template, scheme, mountPoint, client)
		if err != nil {
			return nil, err
		}

		partitions = append(partitions, partition)
	}

	return partitions, nil
}

func getPartitionSchemePartition(template, scheme, mountPoint string, client *ovh.Client) (*Partition, error) {
	r := &Partition{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/partition/%s",
		url.PathEscape(template),
		url.PathEscape(scheme),
		url.PathEscape(mountPoint),
	)

	if err := client.Get(endpoint, r); err != nil {
		return nil, fmt.Errorf("Calling GET %s: %s \n", endpoint, err.Error())
	}

	return r, nil
}

func getPartitionSchemeHardwareRaids(template, scheme string, client *ovh.Client) ([]*HardwareRaid, error) {
	names := []string{}
	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/hardwareRaid",
		url.PathEscape(template),
		url.PathEscape(scheme),
	)
	err := client.Get(endpoint, &names)

	if err != nil {
		return nil, fmt.Errorf("Error calling GET %s: %s \n", endpoint, err.Error())
	}

	hardwareRaids := []*HardwareRaid{}
	for _, name := range names {
		hardwareRaid, err := getPartitionSchemeHardwareRaid(template, scheme, name, client)
		if err != nil {
			return nil, err
		}

		hardwareRaids = append(hardwareRaids, hardwareRaid)
	}

	return hardwareRaids, nil
}

func getPartitionSchemeHardwareRaid(template, scheme, name string, client *ovh.Client) (*HardwareRaid, error) {
	r := &HardwareRaid{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s/hardwareRaid/%s",
		url.PathEscape(template),
		url.PathEscape(scheme),
		url.PathEscape(name),
	)

	if err := client.Get(endpoint, r); err != nil {
		return nil, fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	return r, nil
}
