package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func waitForDbaasLogsOperation(ctx context.Context, c *ovhwrap.Client, serviceName, id string) (*DbaasLogsOperation, error) {
	// Wait for operation status
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING", "RECEIVED", "STARTED", "RETRY", "RUNNING"},
		Target:     []string{"SUCCESS"},
		Refresh:    waitForDbaasLogsOperationCheck(c, serviceName, id),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for dbaas logs operation %s/%s: %s", serviceName, id, err)
	}

	op, ok := res.(*DbaasLogsOperation)
	if !ok {
		return nil, fmt.Errorf(
			"Error waiting for operation %s/%s: got %v instead of DbaasLogsOperation",
			serviceName,
			id,
			reflect.TypeOf(res),
		)
	}

	return op, nil
}

func waitForDbaasLogsOperationCheck(c *ovhwrap.Client, serviceName, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res := &DbaasLogsOperation{}

		endpoint := fmt.Sprintf("/dbaas/logs/%s/operation/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		if err := c.Get(endpoint, res); err != nil {
			log.Printf("[WARNING] error while waiting for dbaas logs operation id %s: %v", id, err)
			return nil, "", err
		}

		log.Printf("[DEBUG] Pending dbaas logs operation: %s", id)
		return res, res.State, nil
	}
}
