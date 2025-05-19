package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

type DomainZone struct {
	DnssecSupported    bool     `json:"dnssecSupported"`
	HasDnsAnycast      bool     `json:"hasDnsAnycast"`
	LastUpdate         string   `json:"lastUpdate"`
	Name               string   `json:"name"`
	NameServers        []string `json:"nameServers"`
	IamResourceDetails `json:"iam"`
}

func (v DomainZone) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["dnssec_supported"] = v.DnssecSupported
	obj["has_dns_anycast"] = v.HasDnsAnycast
	obj["last_update"] = v.LastUpdate
	obj["name"] = v.Name
	obj["urn"] = v.URN

	if v.NameServers != nil {
		obj["name_servers"] = v.NameServers
	}

	return obj
}

type DomainZoneConfirmTerminationOpts struct {
	Token string `json:"token"`
}

type DomainTask struct {
	CanAccelerate bool   `json:"canAccelerate"`
	CanCancel     bool   `json:"canCancel"`
	CanRelaunch   bool   `json:"canRelaunch"`
	Comment       string `json:"comment"`
	CreationDate  string `json:"creationDate"`
	DoneDate      string `json:"doneDate"`
	Function      string `json:"function"`
	TaskID        int    `json:"id"`
	LastUpdate    string `json:"lastUpdate"`
	Status        string `json:"status"`
	TodoDate      string `json:"todoDate"`
}

func waitDomainTask(client *ovhwrap.Client, domainName string, taskId int) error {
	endpoint := fmt.Sprintf("/domain/%s/task/%d", url.PathEscape(domainName), taskId)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"todo", "doing"},
		Target:  []string{"done"},
		Refresh: func() (result interface{}, state string, err error) {
			var task DomainTask

			if err := client.Get(endpoint, &task); err != nil {
				log.Printf("[ERROR] couldn't fetch task %d for domain %s:\n\t%s\n", taskId, domainName, err.Error())
				return nil, "error", err
			}

			return task, task.Status, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err := stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error waiting for domain: %s task: %d to complete:\n\t%s", domainName, taskId, err.Error())
	}

	return err
}
