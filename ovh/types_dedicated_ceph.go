package ovh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

type DedicatedCephCrushTunable string
type DedicatedCephState string
type DedicatedCephACLType string
type DedicatedCephStatus string

const (
	CrushTunableOptimal  DedicatedCephCrushTunable = "OPTIMAL"
	CrushTunableDefault  DedicatedCephCrushTunable = "DEFAULT"
	CrushTunableLegacy   DedicatedCephCrushTunable = "LEGACY"
	CrushTunableBobtail  DedicatedCephCrushTunable = "BOBTAIL"
	CrushTunableArgonaut DedicatedCephCrushTunable = "ARGONAUT"
	CrushTunableFirefly  DedicatedCephCrushTunable = "FIREFLY"
	CrushTunableHammer   DedicatedCephCrushTunable = "HAMMER"
	CrushTunablEjewel    DedicatedCephCrushTunable = "JEWEL"
	StateActive          DedicatedCephState        = "ACTIVE"
	StateSuspended       DedicatedCephState        = "SUSPENDED"
	ACLTypeIPv4          DedicatedCephACLType      = "IPV4"
	ACLTypeIPv6          DedicatedCephACLType      = "IPV6"
	StatusCreating       DedicatedCephStatus       = "CREATING"
	StatusInstalled      DedicatedCephStatus       = "INSTALLED"
	StatusDeleting       DedicatedCephStatus       = "DELETING"
	StatusDeleted        DedicatedCephStatus       = "DELETED"
	StatusTaskInProgress DedicatedCephStatus       = "TASK_IN_PROGRESS"
)

type DedicatedCeph struct {
	ServiceName        string                    `json:"serviceName"`
	CephMonitors       []string                  `json:"cephMons"`
	CephVersion        string                    `json:"cephVersion"`
	CrushTunables      DedicatedCephCrushTunable `json:"crushTunables"`
	Label              string                    `json:"label"`
	Region             string                    `json:"region"`
	Size               float32                   `json:"size"`
	State              DedicatedCephState        `json:"state"`
	Status             DedicatedCephStatus       `json:"status"`
	IamResourceDetails `json:"iam"`
}

type DedicatedCephACL struct {
	Id      int                  `json:"id"`
	Family  DedicatedCephACLType `json:"family"`
	Netmask string               `json:"netmask"`
	Network string               `json:"network"`
}

type DedicatedCephACLCreateOpts struct {
	AclList []string `json:"aclList"`
}

// DedicatedCephTaskInProgress is a task listed by /dedicated/ceph/{serviceName}/task, which lists
// the tasks in progress on a cluster.
type DedicatedCephTaskInProgress struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DedicatedCephTask struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	FinishDate string `json:"finishDate"`
	Type       string `json:"type"`
	CreateDate string `json:"createDate"`
}

func (opts *DedicatedCephACLCreateOpts) FromResource(d *schema.ResourceData) *DedicatedCephACLCreateOpts {
	network := helpers.GetNilStringPointerFromData(d, "network")
	netmask := helpers.GetNilStringPointerFromData(d, "netmask")
	opts.AclList = []string{fmt.Sprintf("%s/%s", *network, *netmask)}
	return opts
}

// DedicatedCephBool is a boolean of the dedicated CEPH API. Although documented as a JSON boolean,
// the API may return it as a string, such as "False" for the backup of a pool, so both are
// accepted.
type DedicatedCephBool bool

func (b *DedicatedCephBool) UnmarshalJSON(data []byte) error {
	var value bool
	if err := json.Unmarshal(data, &value); err == nil {
		*b = DedicatedCephBool(value)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("expected a boolean from the dedicated CEPH API, got %s", data)
	}

	value, err := strconv.ParseBool(str)
	if err != nil {
		return fmt.Errorf("expected a boolean from the dedicated CEPH API, got %q", str)
	}

	*b = DedicatedCephBool(value)
	return nil
}

// DedicatedCephPool is the API representation of /dedicated/ceph/{serviceName}/pool/{poolName}.
type DedicatedCephPool struct {
	ServiceName       string            `json:"serviceName"`
	Name              string            `json:"name"`
	PoolType          string            `json:"poolType"`
	Backup            DedicatedCephBool `json:"backup"`
	MinActiveReplicas int64             `json:"minActiveReplicas"`
	ReplicaCount      int64             `json:"replicaCount"`
}

type DedicatedCephPoolCreateOpts struct {
	PoolName string `json:"poolName"`
}

// DedicatedCephUser is the API representation of /dedicated/ceph/{serviceName}/user/{userName}.
type DedicatedCephUser struct {
	ServiceName string  `json:"serviceName"`
	Name        string  `json:"name"`
	Key         string  `json:"key"`
	MonCaps     string  `json:"monCaps"`
	OsdCaps     string  `json:"osdCaps"`
	MdsCaps     *string `json:"mdsCaps"`
}

type DedicatedCephUserCreateOpts struct {
	UserName string `json:"userName"`
}

// DedicatedCephUserPoolPermission is a user-pool permission, as listed by and sent to
// /dedicated/ceph/{serviceName}/user/{userName}/pool.
type DedicatedCephUserPoolPermission struct {
	PoolName   string `json:"poolName"`
	Read       bool   `json:"read"`
	Write      bool   `json:"write"`
	Execute    bool   `json:"execute"`
	ClassRead  bool   `json:"classRead"`
	ClassWrite bool   `json:"classWrite"`
}

// IsEmpty tells whether the permission grants nothing on its pool.
func (p DedicatedCephUserPoolPermission) IsEmpty() bool {
	return !p.Read && !p.Write && !p.Execute && !p.ClassRead && !p.ClassWrite
}

type DedicatedCephUserPoolPermissionsCreateOpts struct {
	Permissions []DedicatedCephUserPoolPermission `json:"permissions"`
}

const (
	dedicatedCephTaskDone   = "DONE"
	dedicatedCephTaskFailed = "FAILED"
)

// dedicatedCephClusterLocks holds the lock of each dedicated CEPH changed by the provider, by
// service name.
var dedicatedCephClusterLocks sync.Map

// lockDedicatedCephCluster waits until no other change of the provider is in progress on the given
// dedicated CEPH, and returns the function releasing it. The cluster is locked by the API while a
// task is in progress, and a change made meanwhile fails: the first one with a 403, the next ones
// with a 500. So the changes Terraform makes in parallel are made one at a time, the lock being
// held until the task of the change is done.
func lockDedicatedCephCluster(ctx context.Context, serviceName string) (func(), error) {
	lock, _ := dedicatedCephClusterLocks.LoadOrStore(serviceName, make(chan struct{}, 1))
	ch := lock.(chan struct{})

	select {
	case ch <- struct{}{}:
		return func() { <-ch }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// dedicatedCephServerErrorAttempts is the number of times a change failing on a server error is
// made before giving up.
const dedicatedCephServerErrorAttempts = 3

// retryDedicatedCephChange makes the given change of a dedicated CEPH, and makes it again for as
// long as the API refuses it because the cluster is locked by a task in progress, as when Terraform
// changes several pools or users of a cluster at once. It is also made again a few times when the
// API fails on a server error, which it may do right after the cluster is unlocked without making
// the change.
func retryDedicatedCephChange(ctx context.Context, timeout time.Duration, change func() error) error {
	serverErrors := 0

	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err := change()
		if err == nil {
			return nil
		}

		var errOvh *ovh.APIError
		if !errors.As(err, &errOvh) {
			return retry.NonRetryableError(err)
		}

		switch {
		case errOvh.Code == 403 && strings.Contains(errOvh.Message, "is locked"):
			return retry.RetryableError(err)
		case errOvh.Code >= 500:
			serverErrors++
			if serverErrors < dedicatedCephServerErrorAttempts {
				return retry.RetryableError(err)
			}
		}
		return retry.NonRetryableError(err)
	})
}

// waitDedicatedCephIdle waits until no task is in progress on the given dedicated CEPH. A change
// made while a task is in progress fails, and lockDedicatedCephCluster only keeps the changes of
// the provider from overlapping: this also waits for the tasks started elsewhere, such as from
// another Terraform run or the control panel. A few server errors in a row are tolerated.
func waitDedicatedCephIdle(ctx context.Context, client *ovhwrap.Client, serviceName string, timeout time.Duration) error {
	endpoint := "/dedicated/ceph/" + url.PathEscape(serviceName) + "/task"
	serverErrors := 0

	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		var tasks []DedicatedCephTaskInProgress
		if err := client.GetWithContext(ctx, endpoint, &tasks); err != nil {
			err = fmt.Errorf("error calling Get %s: %w", endpoint, err)

			var errOvh *ovh.APIError
			if errors.As(err, &errOvh) && errOvh.Code >= 500 {
				serverErrors++
				if serverErrors < dedicatedCephServerErrorAttempts {
					return retry.RetryableError(err)
				}
			}
			return retry.NonRetryableError(err)
		}
		serverErrors = 0

		if len(tasks) > 0 {
			return retry.RetryableError(fmt.Errorf("%d task(s) in progress on %s, such as %s (%s)", len(tasks), serviceName, tasks[0].Name, tasks[0].Id))
		}
		return nil
	})
}

// waitDedicatedCephTask waits for the given task of a dedicated CEPH to be done. The API returns
// the task as a list, which may be empty while the task is being registered. A few server errors in
// a row are tolerated while polling, as the API may fail on them without the task failing.
func waitDedicatedCephTask(ctx context.Context, client *ovhwrap.Client, serviceName, taskId string, timeout time.Duration) error {
	endpoint := "/dedicated/ceph/" + url.PathEscape(serviceName) + "/task/" + url.PathEscape(taskId)
	serverErrors := 0

	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		var tasks []DedicatedCephTask
		if err := client.GetWithContext(ctx, endpoint, &tasks); err != nil {
			err = fmt.Errorf("error calling Get %s: %w", endpoint, err)

			var errOvh *ovh.APIError
			if errors.As(err, &errOvh) {
				switch {
				case errOvh.Code == 404:
					return retry.RetryableError(fmt.Errorf("task %s of %s not found yet", taskId, serviceName))
				case errOvh.Code >= 500:
					serverErrors++
					if serverErrors < dedicatedCephServerErrorAttempts {
						return retry.RetryableError(err)
					}
				}
			}
			return retry.NonRetryableError(err)
		}
		serverErrors = 0

		if len(tasks) == 0 {
			return retry.RetryableError(fmt.Errorf("task %s of %s not listed yet", taskId, serviceName))
		}

		switch tasks[0].State {
		case dedicatedCephTaskDone:
			return nil
		case dedicatedCephTaskFailed:
			return retry.NonRetryableError(fmt.Errorf("task %s (%s) of %s failed", taskId, tasks[0].Name, serviceName))
		default:
			return retry.RetryableError(fmt.Errorf("task %s of %s is %s", taskId, serviceName, tasks[0].State))
		}
	})
}
