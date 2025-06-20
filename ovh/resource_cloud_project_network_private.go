package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"golang.org/x/exp/slices"
)

func resourceOvhCloudProjectNetworkPrivateImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not OVH_CLOUD_PROJECT/network_id formatted")
	}
	d.SetId(splitId[1])
	d.Set("service_name", splitId[0])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectNetworkPrivate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectNetworkPrivateCreate,
		Read:   resourceCloudProjectNetworkPrivateRead,
		Update: resourceCloudProjectNetworkPrivateUpdate,
		Delete: resourceCloudProjectNetworkPrivateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOvhCloudProjectNetworkPrivateImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"regions_status": {
				Type:       schema.TypeSet,
				Computed:   true,
				Deprecated: "use the regions_attributes field instead",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Required: true,
						},

						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"regions_attributes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Required: true,
						},

						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"openstackid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: RegionAttributesHash,
			},
			"regions_openstack_ids": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudProjectNetworkPrivateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	regions, _ := helpers.StringsFromSchema(d, "regions")

	params := &CloudProjectNetworkPrivateCreateOpts{
		ServiceName: serviceName,
		VlanId:      d.Get("vlan_id").(int),
		Name:        d.Get("name").(string),
		Regions:     regions,
	}

	r := &CloudProjectNetworkPrivateResponse{}

	log.Printf("[DEBUG] Will create public cloud private network: %s", params)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private", params.ServiceName)

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for Private Network %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForCloudProjectNetworkPrivateActive(config.OVHClient, serviceName, r.Id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for private network (%s): %s", params, err)
	}
	log.Printf("[DEBUG] Created Private Network %s", r)

	//set id
	d.SetId(r.Id)

	return resourceCloudProjectNetworkPrivateRead(d, meta)
}

func resourceCloudProjectNetworkPrivateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	r := &CloudProjectNetworkPrivateResponse{}

	log.Printf("[DEBUG] Will read public cloud private network for project: %s, id: %s", serviceName, d.Id())

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", serviceName, d.Id())

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("name", r.Name)
	d.Set("status", r.Status)
	d.Set("type", r.Type)
	d.Set("vlan_id", r.Vlanid)

	regions_status := make([]map[string]interface{}, 0)
	regions_attributes := make([]map[string]interface{}, 0)
	regions_openstack_ids := map[string]string{}
	regions := make([]string, 0)

	for i := range r.Regions {
		region_attributes := make(map[string]interface{})
		region_attributes["region"] = r.Regions[i].Region
		region_attributes["status"] = r.Regions[i].Status
		region_attributes["openstackid"] = r.Regions[i].OpenStackId
		regions_attributes = append(regions_attributes, region_attributes)

		regions_openstack_ids[r.Regions[i].Region] = r.Regions[i].OpenStackId

		region_status := make(map[string]interface{})
		region_status["region"] = r.Regions[i].Region
		region_status["status"] = r.Regions[i].Status
		regions_status = append(regions_status, region_status)

		regions = append(regions, fmt.Sprintf(r.Regions[i].Region))
	}
	d.Set("regions_attributes", regions_attributes)
	d.Set("regions_openstack_ids", regions_openstack_ids)
	d.Set("regions_status", regions_status)
	d.Set("regions", regions)

	d.SetId(r.Id)
	d.Set("service_name", serviceName)

	log.Printf("[DEBUG] Read Public Cloud Private Network %s", r)
	return nil
}

func resourceCloudProjectNetworkPrivateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regions, _ := helpers.StringsFromSchema(d, "regions")

	params := &CloudProjectNetworkPrivateUpdateOpts{
		Name:    d.Get("name").(string),
		Regions: regions,
	}

	log.Printf("[DEBUG] params %s", params)
	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s/region",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)
	for _, reg := range params.Regions {
		param := CloudProjectNetworkPrivateUpdateOptsAlone{
			Region: reg,
		}
		log.Printf("[DEBUG] Will update public cloud private network: %s", param)
		err := config.OVHClient.Post(endpoint, param, nil)
		if err != nil {
			if strings.Contains(err.Error(), "already activated") {
				log.Printf("[DEBUG] Region %s already activated", reg)
				continue
			} else {
				return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, param, err)
			}
		}

		log.Printf("[DEBUG] Waiting for Private Network %s:", reg)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"BUILDING"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForCloudProjectNetworkPrivateActive(config.OVHClient, serviceName, d.Id()),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for private network (%s): %s", params, err)
		}
		log.Printf("[DEBUG] Created Private Network %s", reg)
	}

	log.Printf("[DEBUG] Will read public cloud private network for project: %s, id: %s", serviceName, d.Id())

	endpoint = fmt.Sprintf("/cloud/project/%s/network/private/%s",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)
	r := &CloudProjectNetworkPrivateResponse{}
	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	currentRegions := make([]string, 0)
	for _, r := range r.Regions {
		currentRegions = append(currentRegions, r.Region)
	}

	regionsToRemove := make([]string, 0)
	for _, apiinput := range currentRegions {
		if !slices.Contains(params.Regions, apiinput) {
			regionsToRemove = append(regionsToRemove, apiinput)
		}
	}

	for _, reg := range regionsToRemove {
		endpoint = fmt.Sprintf("/cloud/project/%s/network/private/%s/region/%s",
			url.PathEscape(serviceName),
			url.PathEscape(d.Id()),
			url.PathEscape(reg),
		)

		if err := config.OVHClient.Delete(endpoint, params); err != nil {
			return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, params.Name, err)
		}
	}
	return resourceCloudProjectNetworkPrivateRead(d, meta)
}

func resourceCloudProjectNetworkPrivateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud private network for project: %s, id: %s", serviceName, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("calling %s:\n\t %q", endpoint, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForCloudProjectNetworkPrivateDelete(config.OVHClient, serviceName, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("deleting for private network (%s): %s", id, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Private Network %s", serviceName, id)
	return nil
}

// AttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an Attachment Task.
func waitForCloudProjectNetworkPrivateActive(c *ovhwrap.Client, serviceName, CloudProjectNetworkPrivateId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudProjectNetworkPrivateResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", serviceName, CloudProjectNetworkPrivateId)
		if err := c.Get(endpoint, r); err != nil {
			return r, "", err
		}

		log.Printf("[DEBUG] Pending Private Network: %s", r)
		return r, r.Status, nil
	}
}

// AttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an Attachment Task.
func waitForCloudProjectNetworkPrivateDelete(c *ovhwrap.Client, serviceName, CloudProjectNetworkPrivateId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudProjectNetworkPrivateResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", serviceName, CloudProjectNetworkPrivateId)
		if err := c.Get(endpoint, r); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				log.Printf("[DEBUG] private network id %s on project %s deleted", CloudProjectNetworkPrivateId, serviceName)
				return r, "DELETED", nil
			} else {
				return r, "", err
			}
		}
		log.Printf("[DEBUG] Pending Private Network: %s", r)
		return r, r.Status, nil
	}
}
