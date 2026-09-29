package ovh

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

// dedicatedCephACLTimeout bounds each wait of a change: for the tasks in progress on the cluster,
// the retries of the change while the cluster is locked, and the task of the change.
const dedicatedCephACLTimeout = 30 * time.Minute

func resourceDedicatedCephACL() *schema.Resource {
	return &schema.Resource{
		// Without the timeout of the SDK, which would cut the waits for the other changes of the cluster
		// short: each wait is bounded by dedicatedCephACLTimeout, and cancelled along with Terraform.
		CreateWithoutTimeout: resourceDedicatedCephACLCreate,
		Read:                 resourceDedicatedCephACLRead,
		DeleteWithoutTimeout: resourceDedicatedCephACLDelete,
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

func resourceDedicatedCephACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	acl := (&DedicatedCephACLCreateOpts{}).FromResource(d)
	serviceName := d.Get("service_name").(string)
	url := fmt.Sprintf("/dedicated/ceph/%s/acl", serviceName)

	// Changes are made one at a time on a cluster, see lockDedicatedCephCluster.
	unlock, err := lockDedicatedCephCluster(ctx, serviceName)
	if err != nil {
		return diag.Errorf("Error waiting for the other changes of %s:\n\t%q", serviceName, err)
	}
	defer unlock()

	if err := waitDedicatedCephIdle(ctx, config.OVHClient, serviceName, dedicatedCephACLTimeout); err != nil {
		return diag.Errorf("Error waiting for the tasks in progress on %s:\n\t%q", serviceName, err)
	}

	// create the ACL
	var taskId string
	err = retryDedicatedCephChange(ctx, dedicatedCephACLTimeout, func() error {
		return config.OVHClient.PostWithContext(ctx, url, acl, &taskId)
	})
	if err != nil {
		return diag.Errorf("Error calling POST %s:\n\t%q", url, err)
	}

	// monitor task execution
	if err := waitDedicatedCephTask(ctx, config.OVHClient, serviceName, taskId, dedicatedCephACLTimeout); err != nil {
		return diag.Errorf("Error waiting for CEPH ACL creation:\n\t %q", err)
	}

	// grab the id of the ACL
	acls, err := resourceDedicatedCephACLList(d, meta)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.Errorf("Error listing CEPH ACL, :\n\t cannot find created ACL")
	}

	return diag.FromErr(resourceDedicatedCephACLRead(d, meta))
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

func resourceDedicatedCephACLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	url := fmt.Sprintf("/dedicated/ceph/%s/acl/%s", serviceName, d.Id())

	// Changes are made one at a time on a cluster, see lockDedicatedCephCluster.
	unlock, err := lockDedicatedCephCluster(ctx, serviceName)
	if err != nil {
		return diag.Errorf("Error waiting for the other changes of %s:\n\t%q", serviceName, err)
	}
	defer unlock()

	if err := waitDedicatedCephIdle(ctx, config.OVHClient, serviceName, dedicatedCephACLTimeout); err != nil {
		return diag.Errorf("Error waiting for the tasks in progress on %s:\n\t%q", serviceName, err)
	}

	var taskId string
	err = retryDedicatedCephChange(ctx, dedicatedCephACLTimeout, func() error {
		return config.OVHClient.DeleteWithContext(ctx, url, &taskId)
	})
	if err != nil {
		return diag.Errorf("Error calling DELETE %s:\n\t%q", url, err)
	}

	// monitor task execution
	if err := waitDedicatedCephTask(ctx, config.OVHClient, serviceName, taskId, dedicatedCephACLTimeout); err != nil {
		return diag.Errorf("Error waiting for CEPH ACL deletion:\n\t %q", err)
	}
	d.SetId("")
	return nil
}
