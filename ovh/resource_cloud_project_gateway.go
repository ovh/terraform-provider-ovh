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
)

func resourceOvhCloudProjectGatewayImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()

	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("given ID is not ovh_cloud_project/region/gateway_id formatted")
	}

	d.SetId(splitId[2])
	d.Set("service_name", splitId[0])
	d.Set("region", splitId[1])

	return []*schema.ResourceData{d}, nil
}

func resourceCloudProjectGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectGatewayCreate,
		Read:   resourceCloudProjectGatewayRead,
		Update: resourceCloudProjectGatewayUpdate,
		Delete: resourceCloudProjectGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOvhCloudProjectGatewayImportState,
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
			},
			"model": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_information": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "External information of the gateway",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:        schema.TypeString,
							Description: "External network ID of the gateway",
							Computed:    true,
						},
						"ips": {
							Type:        schema.TypeList,
							Description: "List of external ips of the gateway",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Description: "External IP of the gateway",
										Computed:    true,
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Description: "Subnet ID of the ip",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interfaces list of the gateway",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the interface",
							Computed:    true,
						},
						"ip": {
							Type:        schema.TypeString,
							Description: "IP of the interface",
							Computed:    true,
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "Network ID of the interface",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet ID of the interface",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudProjectGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)
	network := d.Get("network_id").(string)
	subnet := d.Get("subnet_id").(string)

	params := &CloudProjectGatewayCreateOpts{
		Name:  d.Get("name").(string),
		Model: d.Get("model").(string),
	}

	r := &CloudProjectGatewayResponse{}

	log.Printf("[DEBUG] Will create public cloud gateway: %s", params)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/network/%s/subnet/%s/gateway",
		url.PathEscape(serviceName),
		url.PathEscape(region),
		url.PathEscape(network),
		url.PathEscape(subnet))

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for Gateway %+v:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in-progress"},
		Target:     []string{"active"},
		Refresh:    waitForCloudProjectGatewayActive(config.OVHClient, serviceName, r.Id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for gateway (%s): %s", params, err)
	}

	ro := &CloudProjectOperationResponse{}
	endpointo := fmt.Sprintf("/cloud/project/%s/operation/%s",
		url.PathEscape(serviceName),
		url.PathEscape(r.Id))
	if err := config.OVHClient.Get(endpointo, ro); err != nil {
		return nil
	}
	log.Printf("[DEBUG] Created Gateway %+v", ro)

	gatewayId := ro.ResourceId

	// set id
	d.SetId(*gatewayId)

	return resourceCloudProjectGatewayRead(d, meta)
}

func resourceCloudProjectGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)

	r := &CloudProjectGatewayResponse{}

	log.Printf("[DEBUG] Will read public cloud gateway for project: %s, region: %s, id: %s", serviceName, region, d.Id())

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/gateway/%s",
		url.PathEscape(serviceName),
		url.PathEscape(region),
		url.PathEscape(d.Id()))

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("name", r.Name)
	d.Set("model", r.Model)
	d.Set("status", r.Status)
	d.Set("region", region)
	d.SetId(r.Id)
	d.Set("service_name", serviceName)

	externalInfos := make([]map[string]interface{}, 0)
	if r.ExternalInformation != nil {
		externalInfo := make(map[string]interface{})
		ips := make([]map[string]interface{}, 0)
		for _, externalIp := range r.ExternalInformation.Ips {
			ip := make(map[string]interface{})
			ip["ip"] = externalIp.Ip
			ip["subnet_id"] = externalIp.SubnetId
			ips = append(ips, ip)
		}
		externalInfo["ips"] = ips
		externalInfo["network_id"] = r.ExternalInformation.NetworkId
		externalInfos = append(externalInfos, externalInfo)
	}
	d.Set("external_information", externalInfos)

	interfaces := make([]map[string]string, 0)
	for _, responseInterface := range r.Interfaces {
		itf := make(map[string]string)
		itf["id"] = responseInterface.Id
		itf["ip"] = responseInterface.Ip
		itf["subnet_id"] = responseInterface.SubnetId
		itf["network_id"] = responseInterface.NetworkId
		interfaces = append(interfaces, itf)
	}
	d.Set("interfaces", interfaces)

	log.Printf("[DEBUG] Read Public Cloud Gateway %+v", r)
	return nil
}

func resourceCloudProjectGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	params := &CloudProjectGatewayUpdateOpts{
		Name:  d.Get("name").(string),
		Model: d.Get("model").(string),
	}
	region := d.Get("region").(string)

	log.Printf("[DEBUG] Will update public cloud gateway: %s", params)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/gateway/%s",
		url.PathEscape(serviceName),
		url.PathEscape(region),
		url.PathEscape(d.Id()))

	if err := config.OVHClient.Put(endpoint, params, nil); err != nil {
		return fmt.Errorf("calling %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Updated Public cloud %s Gateway %s:", serviceName, d.Id())

	return resourceCloudProjectGatewayRead(d, meta)
}

func resourceCloudProjectGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	region := d.Get("region").(string)

	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud gateway for project: %s, region: %s, id: %s", serviceName, region, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/gateway/%s",
		url.PathEscape(serviceName),
		url.PathEscape(region),
		url.PathEscape(id))

	r := &CloudProjectOperationResponse{}
	if err := config.OVHClient.Delete(endpoint, r); err != nil {
		return fmt.Errorf("calling %s:\n\t %q", endpoint, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in-progress"},
		Target:     []string{"completed"},
		Refresh:    waitForCloudProjectGatewayDelete(config.OVHClient, serviceName, r.Id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("deleting for gateway (%s): %s", id, err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Deleted Public Cloud %s Gateway %s", serviceName, id)
	return nil
}

// AttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an Attachment Task.
func waitForCloudProjectGatewayActive(c *ovhwrap.Client, serviceName, operationId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ro := &CloudProjectOperationResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/operation/%s",
			url.PathEscape(serviceName),
			url.PathEscape(operationId))
		if err := c.Get(endpoint, ro); err != nil {
			return ro, "", err
		}

		log.Printf("[DEBUG] Pending Operation: %+v", ro)

		if ro.ResourceId != nil {
			rg := &CloudProjectGatewayResponse{}
			gatewayId := ro.ResourceId
			endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/gateway/%s",
				url.PathEscape(serviceName),
				url.PathEscape(ro.Regions[0]),
				url.PathEscape(*gatewayId))
			if err := c.Get(endpoint, rg); err != nil {
				return rg, "", err
			}
			return rg, rg.Status, nil
		}

		return ro, ro.Status, nil
	}
}

// AttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an Attachment Task.
func waitForCloudProjectGatewayDelete(c *ovhwrap.Client, serviceName, OperationId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudProjectOperationResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/operation/%s",
			url.PathEscape(serviceName),
			url.PathEscape(OperationId))
		if err := c.Get(endpoint, r); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				log.Printf("[DEBUG] gateway id %s on project %s deleted", OperationId, serviceName)
				return r, "DELETED", nil
			} else {
				return r, "", err
			}
		}
		log.Printf("[DEBUG] Pending Gateway: %+v", r)
		return r, r.Status, nil
	}
}
