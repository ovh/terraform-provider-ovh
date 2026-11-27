package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

// vpsV1Task is the part of vps.Task, the answer of the API v1 actions of a
// VPS and of GET /vps/{serviceName}/tasks/{id}, read by the resource.
type vpsV1Task struct {
	Id    int64  `json:"id"`
	State string `json:"state"`
	Type  string `json:"type"`
}

// vpsV1Power is the part of GET /vps/{serviceName} read to change the power
// state of a VPS.
type vpsV1Power struct {
	State       string `json:"state"`
	NetbootMode string `json:"netbootMode"` // local or rescue
}

// vpsPostAction calls POST endpoint with body (none when nil) and returns the
// task it starts. API v1 answers 409 while another task runs on the VPS: the
// call is tried again every vpsPollInterval until deadline.
func vpsPostAction(ctx context.Context, client *ovhwrap.Client, endpoint string, body any, deadline time.Time) (*vpsV1Task, error) {
	for {
		var task vpsV1Task
		err := client.PostWithContext(ctx, endpoint, body, &task)
		if err == nil {
			return &task, nil
		}
		if vpsAPIStatus(err) != http.StatusConflict {
			return nil, fmt.Errorf("calling POST %s: %w", endpoint, err)
		}
		if time.Now().Add(vpsPollInterval).After(deadline) {
			return nil, fmt.Errorf("calling POST %s: another task still runs on the VPS: %w", endpoint, err)
		}
		log.Printf("[INFO] calling POST %s: %s, trying again", endpoint, err)
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("calling POST %s: %w", endpoint, ctx.Err())
		case <-time.After(vpsPollInterval):
		}
	}
}

// vpsWaitTask waits until deadline for task, started on VPS id, to end, then
// for the VPS to leave its transitional states. A task API v1 no longer knows
// (404) has ended. A task that ends in error or cancelled is an error; any
// other state, or a failure to read the task, keeps the wait going.
func vpsWaitTask(ctx context.Context, client *ovhwrap.Client, id string, task *vpsV1Task, deadline time.Time) error {
	endpoint := "/vps/" + url.PathEscape(id) + "/tasks/" + url.PathEscape(strconv.FormatInt(task.Id, 10))
	conf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"done"},
		Refresh: func() (any, string, error) {
			var current vpsV1Task
			err := client.GetWithContext(ctx, endpoint, &current)
			switch {
			case vpsAPIStatus(err) == http.StatusNotFound:
				return task, "done", nil
			case err != nil:
				log.Printf("[WARN] VPS %s: calling GET %s: %s", id, endpoint, err)
				return task, "pending", nil
			}
			switch current.State {
			case "done":
				return &current, "done", nil
			case "error", "cancelled":
				return nil, "", fmt.Errorf("VPS task %d (%s) ended in state %s", task.Id, task.Type, current.State)
			}
			return &current, "pending", nil
		},
		Timeout:      time.Until(deadline),
		Delay:        vpsPollDelay,
		PollInterval: vpsPollInterval,
	}
	if _, err := conf.WaitForStateContext(ctx); err != nil {
		var timeout *retry.TimeoutError
		if errors.As(err, &timeout) {
			return fmt.Errorf("waiting for VPS task %d (%s): %w", task.Id, task.Type, err)
		}
		return err
	}

	if err := vpsWaitReady(ctx, client, id, time.Until(deadline)); err != nil {
		return fmt.Errorf("waiting for VPS %s after its task %d (%s): %w", id, task.Id, task.Type, err)
	}
	return nil
}

// vpsCronPattern matches the cron expressions API v1 can apply: daily, at a
// minute ("M H * * *", UTC), written as GET reads them back, without leading
// zeros.
var vpsCronPattern = regexp.MustCompile(`^([0-9]|[1-5][0-9]) ([0-9]|1[0-9]|2[0-3]) \* \* \*$`)

// vpsScheduleFromCron converts the cron expression of the API v2 contract
// ("M H * * *") to the API v1 automated backup schedule ("HH:MM:SS", UTC).
// It is the inverse of vpsCronFromV1Schedule.
func vpsScheduleFromCron(cron string) (string, error) {
	match := vpsCronPattern.FindStringSubmatch(cron)
	if match == nil {
		return "", fmt.Errorf("invalid automated backup cron %q, expected \"M H * * *\" with M from 0 to 59 and H from 0 to 23", cron)
	}
	// The pattern only matches numbers.
	minute, _ := strconv.Atoi(match[1])
	hour, _ := strconv.Atoi(match[2])
	return fmt.Sprintf("%02d:%02d:00", hour, minute), nil
}

// vpsImageIdForName returns the id of the image named name among the images
// VPS id can be installed from. When none has this name, imageId is empty and
// names lists the names of the available images, sorted.
func vpsImageIdForName(ctx context.Context, client *ovhwrap.Client, id, name string) (imageId string, names []string, err error) {
	images, err := vpsAvailableImages(ctx, client, id)
	if err != nil {
		return "", nil, err
	}
	for _, image := range images {
		if image.Name == name {
			return image.Id, nil, nil
		}
		names = append(names, image.Name)
	}
	slices.Sort(names)
	return "", names, nil
}

// vpsAvailableImages reads the images VPS id can be installed from,
// vpsV1ListParallelism at once, in the order of GET images/available.
func vpsAvailableImages(ctx context.Context, client *ovhwrap.Client, id string) ([]vpsV1Image, error) {
	var ids []string
	endpoint := "/vps/" + url.PathEscape(id) + "/images/available"
	if err := client.GetWithContext(ctx, endpoint, &ids); err != nil {
		return nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	images := make([]vpsV1Image, len(ids))
	errs := make([]error, len(ids))
	slots := make(chan struct{}, vpsV1ListParallelism)
	var wg sync.WaitGroup
	for i, imageId := range ids {
		slots <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			endpoint := endpoint + "/" + url.PathEscape(imageId)
			if err := client.GetWithContext(ctx, endpoint, &images[i]); err != nil {
				errs[i] = fmt.Errorf("calling GET %s: %w", endpoint, err)
			}
		}()
	}
	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return images, nil
}

// vpsSetPowerState brings VPS id to the power state target (ACTIVE, SHUTOFF
// or RESCUE) and waits until deadline for each task it starts.
//
// API v1 boots a VPS on the rescue system through its netboot mode, which only
// applies at the next boot: a VPS whose netboot mode changes, or that runs on
// another system than the target one, is rebooted, or started when stopped. A
// VPS stopped for SHUTOFF gets the local netboot mode back, so that it next
// starts on its disk.
func vpsSetPowerState(ctx context.Context, client *ovhwrap.Client, id, target string, deadline time.Time) error {
	netboot := "local"
	switch target {
	case "ACTIVE", "SHUTOFF":
	case "RESCUE":
		netboot = "rescue"
	default:
		return fmt.Errorf("cannot bring VPS %s to the power state %q", id, target)
	}

	base := "/vps/" + url.PathEscape(id)
	var vps vpsV1Power
	if err := client.GetWithContext(ctx, base, &vps); err != nil {
		return fmt.Errorf("calling GET %s: %w", base, err)
	}

	changed := false
	if vps.NetbootMode != netboot {
		body := struct {
			NetbootMode string `json:"netbootMode"`
		}{netboot}
		if err := client.PutWithContext(ctx, base, body, nil); err != nil {
			return fmt.Errorf("calling PUT %s: %w", base, err)
		}
		changed = true
	}

	var action string
	switch {
	case target == "SHUTOFF":
		if vps.State != "stopped" {
			action = "stop"
		}
	case vps.State == "stopped":
		action = "start"
	case changed || (vps.State == "rescued") != (target == "RESCUE"):
		// Also when the netboot mode is already the right one but the VPS has
		// not booted on it yet.
		action = "reboot"
	}
	if action == "" {
		return nil
	}

	task, err := vpsPostAction(ctx, client, base+"/"+action, nil, deadline)
	if err != nil {
		return err
	}
	return vpsWaitTask(ctx, client, id, task, deadline)
}
