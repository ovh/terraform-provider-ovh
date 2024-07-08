package ovh

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type DedicatedNASHA struct {
	ServiceName        string `json:"serviceName,omitempty"`
	Monitored          bool   `json:"monitored,omitempty"`
	ZpoolSize          int    `json:"zpoolSize,omitempty"`
	CustomName         string `json:"customName,omitempty"`
	Datacenter         string `json:"datacenter,omitempty"`
	DiskType           string `json:"diskType,omitempty"`
	CanCreatePartition bool   `json:"canCreatePartition,omitempty"`
	Ip                 string `json:"ip,omitempty"`
	ZpoolCapacity      int    `json:"zpoolCapacity,omitempty"`
	IamResourceDetails `json:"iam"`
}

type DedicatedNASHAPartition struct {
	Name            string `json:"partitionName,omitempty"`
	Description     string `json:"partitionDescription,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	Size            int    `json:"size,omitempty"`
	Capacity        int    `json:"partitionCapacity,omitempty"`
	UsedBySnapshots int    `json:"usedBySnapshots,omitempty"`
}

type DedicatedNASHAPartitionAccess struct {
	IP             string `json:"ip"`
	Type           string `json:"type,omitempty"` // readonly, readwrite
	AclDescription string `json:"aclDescription,omitempty"`
}

type DedicatedNASHAPartitionSnapshot struct {
	Type string `json:"snapshotType"` // day-1, day-2, day-3, day-7, hour-1, hour-6
}

type DedicatedNASHATask struct {
	ID          int    `json:"taskId"`
	StorageName string `json:"storageName"`
	Status      string `json:"status"`
	// "cancelled"
	// "customerError"
	// "doing"
	// "done"
	// "init"
	// "ovhError"
	// "todo"
	Details       string    `json:"details"`
	LastUpdate    time.Time `json:"lastUpdate"`
	TodoDate      time.Time `json:"todoDate"`
	PartitionName string    `json:"partitionName"`
	Operation     string    `json:"operation"`
	DoneDate      time.Time `json:"doneDate"`
}

func (t *DedicatedNASHATask) StateChangeConf(d *schema.ResourceData, meta interface{}) *resource.StateChangeConf {
	config := meta.(*Config)
	return &resource.StateChangeConf{
		Pending: []string{"todo", "init", "doing"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {
			resp := &DedicatedNASHATask{}
			endpoint := fmt.Sprintf("/dedicated/nasha/%s/task/%d", d.Get("service_name").(string), t.ID)
			err := config.OVHClient.Get(endpoint, resp)
			if err != nil {
				return nil, "", err
			}
			return d, resp.Status, nil

		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
}
