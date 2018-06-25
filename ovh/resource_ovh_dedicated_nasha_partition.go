package ovh

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDedicatedNASHAPartition() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedNASHAPartitionCreate,
		Read:   resourceDedicatedNASHAPartitionRead,
		Update: resourceDedicatedNASHAPartitionUpdate,
		Delete: resourceDedicatedNASHAPartitionDelete,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"capacity": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_by_snapshots": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDedicatedNASHAPartitionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	name := d.Get("name").(string)

	partition := &DedicatedNASHAPartition{
		Name:     name,
		Protocol: d.Get("protocol").(string),
		Size:     d.Get("size").(int),
	}

	resp := &DedicatedNASHATask{}

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition", serviceName)

	err := config.OVHClient.Post(endpoint, partition, resp)
	if err != nil {
		return fmt.Errorf("calling %s with params %q:\n\t %s", endpoint, partition, err.Error())
	}

	stateConf := resp.StateChangeConf(d, meta)

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for NASHA partition (%q): %s", partition, err.Error())
	}
	log.Printf("[DEBUG] Created NASHA partition")

	d.SetId(fmt.Sprintf("dedicated_nasha_%s_partition_%s", serviceName, name))

	// resourceDedicatedNASHAPartitionRead(d, meta)

	return nil
}

func resourceDedicatedNASHAPartitionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	name := d.Get("name").(string)

	resp := &DedicatedNASHAPartition{}

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s", serviceName, name)

	err := config.OVHClient.Get(endpoint, resp)
	if err != nil {
		if err.Error() == fmt.Sprintf("Error 404: \"The requested object (partitionName = %s) does not exist\"", name) {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("Error calling %s:\n\t '%q'", endpoint, err)
		}
	}

	d.Set("size", resp.Size)
	d.Set("protocol", resp.Protocol)
	d.Set("capacity", resp.Capacity)
	d.Set("used_by_snapshots", resp.UsedBySnapshots)

	return nil
}

func resourceDedicatedNASHAPartitionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	name := d.Get("name").(string)

	partition := &DedicatedNASHAPartition{
		Size: d.Get("size").(int),
	}
	resp := &DedicatedNASHAPartition{}

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s", serviceName, name)

	err := config.OVHClient.Put(endpoint, partition, resp)
	if err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	stateConf := &resource.StateChangeConf{
		Target: []string{string(d.Get("size").(int))},
		Refresh: func() (interface{}, string, error) {
			resp := &DedicatedNASHAPartition{}
			err := config.OVHClient.Get(endpoint, resp)
			if err != nil {
				return nil, "", err
			}
			return d, string(resp.Size), nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for NASHA partition update: %s", err)
	}

	return nil
}

func resourceDedicatedNASHAPartitionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	name := d.Get("name").(string)

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s", serviceName, name)

	resp := &DedicatedNASHATask{}

	err := config.OVHClient.Delete(endpoint, resp)
	if err != nil {
		return fmt.Errorf("calling %s :\n\t %q", endpoint, err)
	}

	stateConf := resp.StateChangeConf(d, meta)

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for NASHA partition delete: %s", err)
	}

	log.Printf("[DEBUG] Deleted NASHA partition")

	return nil
}
