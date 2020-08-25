package ovh

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"

	"github.com/ovh/go-ovh/ovh"
)

func resourceOvhCloudNetworkPrivateImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not OVH_PROJECT_ID/network_id formatted")
	}
	d.SetId(splitId[1])
	d.Set("project_id", splitId[0])
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudNetworkPrivate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudNetworkPrivateCreate,
		Read:   resourceCloudNetworkPrivateRead,
		Update: resourceCloudNetworkPrivateUpdate,
		Delete: resourceCloudNetworkPrivateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOvhCloudNetworkPrivateImportState,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_PROJECT_ID", nil),
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
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"regions_status": {
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
					},
				},
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

func resourceCloudNetworkPrivateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	regions, _ := helpers.StringsFromSchema(d, "regions")

	params := &CloudNetworkPrivateCreateOpts{
		ProjectId: projectId,
		VlanId:    d.Get("vlan_id").(int),
		Name:      d.Get("name").(string),
		Regions:   regions,
	}

	r := &CloudNetworkPrivateResponse{}

	log.Printf("[DEBUG] Will create public cloud private network: %s", params)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private", params.ProjectId)

	err := config.OVHClient.Post(endpoint, params, r)
	if err != nil {
		return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for Private Network %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForCloudNetworkPrivateActive(config.OVHClient, projectId, r.Id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for private network (%s): %s", params, err)
	}
	log.Printf("[DEBUG] Created Private Network %s", r)

	//set id
	d.SetId(r.Id)

	return nil
}

func resourceCloudNetworkPrivateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)

	r := &CloudNetworkPrivateResponse{}

	log.Printf("[DEBUG] Will read public cloud private network for project: %s, id: %s", projectId, d.Id())

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, d.Id())

	d.Partial(true)
	err := config.OVHClient.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	err = readCloudNetworkPrivate(config, d, r)
	if err != nil {
		return err
	}
	d.Partial(false)

	log.Printf("[DEBUG] Read Public Cloud Private Network %s", r)
	return nil
}

func resourceCloudNetworkPrivateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	params := &CloudNetworkPrivateUpdateOpts{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Will update public cloud private network: %s", params)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, d.Id())

	err := config.OVHClient.Put(endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Updated Public cloud %s Private Network %s:", projectId, d.Id())

	return resourceCloudNetworkPrivateRead(d, meta)
}

func resourceCloudNetworkPrivateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud private network for project: %s, id: %s", projectId, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, id)

	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %q", endpoint, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForCloudNetworkPrivateDelete(config.OVHClient, projectId, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("deleting for private network (%s): %s", id, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Private Network %s", projectId, id)
	return nil
}

func readCloudNetworkPrivate(config *Config, d *schema.ResourceData, r *CloudNetworkPrivateResponse) error {
	d.Set("name", r.Name)
	d.Set("status", r.Status)
	d.Set("type", r.Type)
	d.Set("vlan_id", r.Vlanid)

	regions_status := make([]map[string]interface{}, 0)
	regions := make([]string, 0)

	for i := range r.Regions {
		region := make(map[string]interface{})
		region["region"] = r.Regions[i].Region
		region["status"] = r.Regions[i].Status
		regions_status = append(regions_status, region)
		regions = append(regions, fmt.Sprintf(r.Regions[i].Region))
	}
	d.Set("regions_status", regions_status)
	d.Set("regions", regions)

	d.SetId(r.Id)
	return nil
}

func cloudNetworkPrivateExists(projectId, id string, c *ovh.Client) error {
	r := &CloudNetworkPrivateResponse{}

	log.Printf("[DEBUG] Will read public cloud private network for project: %s, id: %s", projectId, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, id)

	err := c.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling %s:\n\t %q", endpoint, err)
	}
	log.Printf("[DEBUG] Read public cloud private network: %s", r)

	return nil
}

// AttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an Attachment Task.
func waitForCloudNetworkPrivateActive(c *ovh.Client, projectId, CloudNetworkPrivateId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudNetworkPrivateResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, CloudNetworkPrivateId)
		err := c.Get(endpoint, r)
		if err != nil {
			return r, "", err
		}

		log.Printf("[DEBUG] Pending Private Network: %s", r)
		return r, r.Status, nil
	}
}

// AttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an Attachment Task.
func waitForCloudNetworkPrivateDelete(c *ovh.Client, projectId, CloudNetworkPrivateId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudNetworkPrivateResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, CloudNetworkPrivateId)
		err := c.Get(endpoint, r)
		if err != nil {
			if err.(*ovh.APIError).Code == 404 {
				log.Printf("[DEBUG] private network id %s on project %s deleted", CloudNetworkPrivateId, projectId)
				return r, "DELETED", nil
			} else {
				return r, "", err
			}
		}
		log.Printf("[DEBUG] Pending Private Network: %s", r)
		return r, r.Status, nil
	}
}
