package ovh

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDedicatedNASHAPartitionSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedNASHAPartitionSnapshotCreate,
		ReadContext:   resourceDedicatedNASHAPartitionSnapshotRead,
		DeleteContext: resourceDedicatedNASHAPartitionSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"partition_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDedicatedNASHAPartitionSnapshotCreate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	partitionName := d.Get("partition_name").(string)
	snapshotType := d.Get("type").(string)

	snapshot := &DedicatedNASHAPartitionSnapshot{
		Type: snapshotType,
	}

	resp := &DedicatedNASHATask{}

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/snapshot", serviceName, partitionName)

	err := config.OVHClient.Post(endpoint, snapshot, resp)
	if err != nil {
		return diag.Errorf("calling %s with params %q:\n\t %s", endpoint, snapshot, err.Error())
	}

	stateConf := resp.StateChangeConf(d, meta)

	_, err = stateConf.WaitForStateContext(c)
	if err != nil {
		return diag.Errorf("waiting for HA-NAS snapshot (%q): %s", snapshot, err.Error())
	}
	log.Printf("[DEBUG] Created HA-NAS snapshot")

	d.SetId(fmt.Sprintf("%s/%s/%s", serviceName, partitionName, snapshotType))

	return nil
}

func resourceDedicatedNASHAPartitionSnapshotRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	if strings.Contains(d.Id(), "/") {
		tab := strings.Split(d.Id(), "/")
		if len(tab) != 3 {
			return diag.Errorf("can't parse snapshot ID: %s", d.Id())
		}

		d.Set("service_name", tab[0])
		d.Set("partition_name", tab[1])
		d.Set("type", tab[2])
	}

	serviceName := d.Get("service_name").(string)
	partitionName := d.Get("partition_name").(string)
	snapshotType := d.Get("type").(string)

	resp := &DedicatedNASHAPartition{}

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/snapshot/%s", serviceName, partitionName, snapshotType)

	err := config.OVHClient.Get(endpoint, resp)
	if err != nil {
		if err.Error() == fmt.Sprintf("Error 404: \"The requested object (type = %s) does not exist\"", snapshotType) {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("Error calling %s:\n\t '%q'", endpoint, err)
		}
	}

	return nil
}

func resourceDedicatedNASHAPartitionSnapshotDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	partitionName := d.Get("partition_name").(string)
	snapshotType := d.Get("type").(string)

	endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/snapshot/%s", serviceName, partitionName, snapshotType)

	resp := &DedicatedNASHATask{}

	err := config.OVHClient.Delete(endpoint, resp)
	if err != nil {
		return diag.Errorf("calling %s :\n\t %q", endpoint, err)
	}

	stateConf := resp.StateChangeConf(d, meta)

	_, err = stateConf.WaitForStateContext(c)
	if err != nil {
		return diag.Errorf("waiting for HA-NAS shapshot delete: %s", err)
	}

	log.Printf("[DEBUG] Deleted HA-NAS shapshot")

	return nil
}
