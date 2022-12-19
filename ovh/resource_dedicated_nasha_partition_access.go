package ovh

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDedicatedNASHAPartitionAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedNASHAPartitionAccessCreate,
		ReadContext:   resourceDedicatedNASHAPartitionAccessRead,
		DeleteContext: resourceDedicatedNASHAPartitionAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(c context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"partition_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDedicatedNASHAPartitionAccessCreate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	partitionName := d.Get("partition_name").(string)
	ipsubnet, _ := d.Get("ip").(string)

	access := &DedicatedNASHAPartitionAccess{
		IP:   d.Get("ip").(string),
		Type: d.Get("type").(string),
	}

	resp := &DedicatedNASHATask{}

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/access", serviceName, partitionName)

	err := config.OVHClient.Post(endpoint, access, resp)
	if err != nil {
		return diag.Errorf("calling %s with params %s:\n\t %q", endpoint, access, err)
	}

	stateConf := resp.StateChangeConf(d, meta)

	_, err = stateConf.WaitForStateContext(c)
	if err != nil {
		return diag.Errorf("waiting for NASHA partition access (%s): %s", access, err)
	}
	log.Printf("[DEBUG] Created NASHA partition access")

	d.SetId(fmt.Sprintf("dedicated_nasha_%s_partition_%s_access_%s", serviceName, partitionName, ipsubnet))

	return nil
}

func resourceDedicatedNASHAPartitionAccessRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	partitionName := d.Get("partition_name").(string)
	ipsubnet, _ := d.Get("ip").(string)

	resp := &DedicatedNASHAPartitionAccess{}

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/access/%s", serviceName, partitionName, url.PathEscape(ipsubnet))

	err := config.OVHClient.Get(endpoint, resp)
	if err != nil {
		if err.Error() == fmt.Sprintf("Error 404: \"The requested object (ip = %s) does not exist\"", ipsubnet) ||
			err.Error() == fmt.Sprintf("Error 404: \"The requested object (partitionName = %s) does not exist\"", partitionName) {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("Error calling %s:\n\t '%q'", endpoint, err)
		}
	}
	d.Set("type", resp.Type)

	return nil
}

func resourceDedicatedNASHAPartitionAccessDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	partitionName := d.Get("partition_name").(string)
	ipsubnet, _ := d.Get("ip").(string)

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/access/%s", serviceName, partitionName, url.PathEscape(ipsubnet))

	resp := &DedicatedNASHATask{}

	err := config.OVHClient.Delete(endpoint, resp)
	if err != nil {
		return diag.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	stateConf := resp.StateChangeConf(d, meta)

	_, err = stateConf.WaitForStateContext(c)
	if err != nil {
		return diag.Errorf("waiting for NASHA partition access: %s", err)
	}

	return nil
}
