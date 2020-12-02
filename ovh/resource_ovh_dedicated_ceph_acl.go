package ovh

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDedicatedCephACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedCephACLCreate,
		Read:   resourceDedicatedCephACLRead,
		Delete: resourceDedicatedCephACLDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDedicatedCephACLImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Computed: false,
				ForceNew: true,
				Required: true,
			},
			"family": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIp(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"netmask": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIp(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func resourceDedicatedCephACLImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import id is not service_name/acl_id formatted")
	}
	id := splitId[0]
	aclId := splitId[1]
	d.SetId(aclId)
	d.Set("service_name", id)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceDedicatedCephACLList(d *schema.ResourceData, meta interface{}) ([]DedicatedCephACL, error) {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	url := fmt.Sprintf("/dedicated/ceph/%s/acl", serviceName)
	var aclResp []DedicatedCephACL
	err := config.OVHClient.Get(url, &aclResp)
	if err != nil {
		return nil, fmt.Errorf("Error calling GET %s:\n\t%q", url, err)
	}
	return aclResp, nil
}

func resourceDedicatedCephACLCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	acl := (&DedicatedCephACLCreateOpts{}).FromResource(d)
	serviceName := d.Get("service_name").(string)
	url := fmt.Sprintf("/dedicated/ceph/%s/acl", serviceName)

	// create the ACL
	var taskId string
	err := config.OVHClient.Post(url, acl, &taskId)
	if err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t%q", url, err)
	}

	// monitor task execution
	stateConf := &resource.StateChangeConf{
		Target: []string{"DONE"},
		Refresh: func() (interface{}, string, error) {
			url = fmt.Sprintf("/dedicated/ceph/%s/task/%s", serviceName, taskId)
			var stateResp []DedicatedCephTask
			err := config.OVHClient.Get(url, &stateResp)
			if err != nil {
				return nil, "", err
			}
			return d, stateResp[0].State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for CEPH ACL creation:\n\t %q", err)
	}

	// grab the id of the ACL
	acls, err := resourceDedicatedCephACLList(d, meta)
	if err != nil {
		return err
	}
	found := false
	for _, item := range acls {
		if item.Netmask == d.Get("netmask") && item.Network == d.Get("network") {
			d.SetId(fmt.Sprintf("%d", item.Id))
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Error listing CEPH ACL, :\n\t cannot find created ACL")
	}

	return resourceDedicatedCephACLRead(d, meta)
}

func resourceDedicatedCephACLRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	id := d.Get("service_name").(string)
	url := fmt.Sprintf("/dedicated/ceph/%s/acl/%s", id, d.Id())
	resp := &DedicatedCephACL{}

	if err := config.OVHClient.Get(url, resp); err != nil {
		return helpers.CheckDeleted(d, err, url)
	}

	d.Set("netmask", resp.Netmask)
	d.Set("family", resp.Family)
	d.Set("network", resp.Network)
	return nil
}

func resourceDedicatedCephACLDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	url := fmt.Sprintf("/dedicated/ceph/%s/acl/%s", serviceName, d.Id())
	var taskId string
	err := config.OVHClient.Delete(url, &taskId)
	if err != nil {
		return fmt.Errorf("Error calling DELETE %s:\n\t%q", url, err)
	}

	// monitor task execution
	stateConf := &resource.StateChangeConf{
		Target: []string{"DONE"},
		Refresh: func() (interface{}, string, error) {
			url = fmt.Sprintf("/dedicated/ceph/%s/task/%s", serviceName, taskId)
			var stateResp []DedicatedCephTask
			err := config.OVHClient.Get(url, &stateResp)
			if err != nil {
				return nil, "", err
			}
			return d, stateResp[0].State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for CEPH ACL deletion:\n\t %q", err)
	}
	d.SetId("")
	return nil
}
