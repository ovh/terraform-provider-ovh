package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func resourceDedicatedServerNetworking() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedServerNetworkingCreate,
		Read:   resourceDedicatedServerNetworkingRead,
		Delete: resourceDedicatedServerNetworkingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your dedicated server.",
			},
			"interfaces": {
				// we want to have the interfaces in a determinist order
				Type: schema.TypeSet,
				// we can bond all 4 interfaces
				// we can bond 2x2 interfaces
				// we cannot bond 3 interfaces and leave one alone
				// we cannot separate all 4 interfaces
				MinItems:    1,
				MaxItems:    2,
				Required:    true,
				ForceNew:    true,
				Description: "Interface or interfaces aggregation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"macs": {
							// we want to have the MACs in a determinist order
							Type:        schema.TypeSet,
							Required:    true,
							ForceNew:    true,
							Description: "Interface Mac address",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Interface type",
							// The ValidateFunc expect a signature func(val any, key string) (warns []string, errs []error)
							// we are not yet using go 1.18+ as such we cannot use any
							// Once in 1.18 we can add a validation to enforce type is either public or vrack
						},
					},
				},
			},

			//Computed
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation description",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation status",
			},
		},
	}
}

func resourceDedicatedServerNetworkingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	opts := (&DedicatedServerNetworkingCreateOpts{}).FromResource(d)
	log.Printf("%v\n", opts)

	// Before trying to manipulate networking details on the given server, let's make sure no operations are in progress
	if err := waitForDedicatedServerNetworking(serviceName, config.OVHClient); err != nil {
		return err
	}

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/networking",
		url.PathEscape(serviceName),
	)

	dedicatedServerNetworking := DedicatedServerNetworking{}
	if err := config.OVHClient.Post(endpoint, opts, dedicatedServerNetworking); err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	// Once new networking details have been sent let's wait until the changes are active
	if err := waitForDedicatedServerNetworking(serviceName, config.OVHClient); err != nil {
		return err
	}

	d.SetId(serviceName)
	return resourceDedicatedServerNetworkingRead(d, meta)
}

func resourceDedicatedServerNetworkingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	// To make it possible to use terraform import let's use d.Id() instead of d.Get("service_name").(string)
	serviceName := d.Id()

	serverNetworkingDetails, err := getDedicatedServerNetworkingDetails(serviceName, config.OVHClient)
	if err != nil {
		return fmt.Errorf("Error retrieving networking details for %s:\n\t %q", serviceName, err)
	}

	d.Set("service_name", serviceName)
	d.Set("description", serverNetworkingDetails.Status)
	d.Set("status", serverNetworkingDetails.Status)

	networkInterfaces := make([]map[string]interface{}, len(serverNetworkingDetails.Interfaces))
	for i, networkInterfaceDetails := range serverNetworkingDetails.Interfaces {

		networkInterface := make(map[string]interface{})
		networkInterface["type"] = networkInterfaceDetails.Type
		macs := networkInterfaceDetails.Macs

		// we want the MACs associated to an interface to be in a determist order to avoid false positive diff
		sort.Strings(macs)
		networkInterface["macs"] = macs

		networkInterfaces[i] = networkInterface
	}

	// we want interfaces to be in a determist order to avoid false positive diff
	sort.SliceStable(networkInterfaces, func(i, j int) bool {
		return networkInterfaces[i]["type"].(string) < networkInterfaces[j]["type"].(string)
	})

	err = d.Set("interfaces", networkInterfaces)
	if err != nil {
		return fmt.Errorf("Error persisting interfaces in state for %s:\n\t %q", serviceName, err)
	}
	return nil
}

func resourceDedicatedServerNetworkingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	// Before trying to manipulate networking details on the given server, let's make sure no operations are in progress
	if err := waitForDedicatedServerNetworking(serviceName, config.OVHClient); err != nil {
		return err
	}

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/networking",
		url.PathEscape(serviceName),
	)

	dedicatedServerNetworking := DedicatedServerNetworking{}
	if err := config.OVHClient.Delete(endpoint, dedicatedServerNetworking); err != nil {
		return fmt.Errorf("Error calling DELETE %s:\n\t %q", endpoint, err)
	}

	// Once new networking details have been sent let's wait until the changes are active
	if err := waitForDedicatedServerNetworking(serviceName, config.OVHClient); err != nil {
		return err
	}

	// we cant delete the task through the API, just forget about its Id
	d.SetId("")
	return nil
}

func waitForDedicatedServerNetworking(serviceName string, c *ovh.Client) error {

	refreshFunc := func() (interface{}, string, error) {
		var taskErr error
		var serverNetworkingDetails *DedicatedServerNetworking

		// The Dedicated Server API often returns 500/404 errors
		// in such case we retry to retrieve task status
		// 404 may happen because of some inconsistency between the
		// api endpoint call and the target region executing the task
		retryErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			serverNetworkingDetails, err = getDedicatedServerNetworkingDetails(serviceName, c)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 500) {
					return resource.RetryableError(err)
				}
				// other error dont retry and fail
				taskErr = err
			}
			return nil
		})

		if retryErr != nil {
			return serverNetworkingDetails, "", retryErr
		}

		if taskErr != nil {
			return serverNetworkingDetails, "", taskErr
		}

		log.Printf("[INFO] Networking parameter for %s: %s", serviceName, serverNetworkingDetails.Status)
		return serverNetworkingDetails, serverNetworkingDetails.Status, nil
	}

	log.Printf("[INFO] Waiting for networking details to be applied for %s", serviceName)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deploying"},
		Target:     []string{"active"},
		Refresh:    refreshFunc,
		Timeout:    45 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Dedicated Server networking task for %s to complete: %s", serviceName, err)
	}

	return nil
}

func getDedicatedServerNetworkingDetails(serviceName string, c *ovh.Client) (*DedicatedServerNetworking, error) {
	serverNetworkingDetails := &DedicatedServerNetworking{}
	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/networking",
		url.PathEscape(serviceName),
	)

	if err := c.Get(endpoint, serverNetworkingDetails); err != nil {
		return nil, err
	}

	return serverNetworkingDetails, nil
}
