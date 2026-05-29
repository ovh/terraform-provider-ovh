package ovh

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// VPSSnapshot represents a VPS snapshot as returned by /vps/{serviceName}/snapshot.
type VPSSnapshot struct {
	ID           string    `json:"id,omitempty"`
	Description  string    `json:"description,omitempty"`
	CreationDate time.Time `json:"creationDate,omitempty"`
	Region       string    `json:"region,omitempty"`
}

// VPSSnapshotDownload represents the download payload returned by
// /vps/{serviceName}/snapshot/download.
type VPSSnapshotDownload struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

// VPSCreateSnapshotOpts is the body passed to /vps/{serviceName}/createSnapshot.
type VPSCreateSnapshotOpts struct {
	Description string `json:"description,omitempty"`
}

// vpsSnapshotTaskStateChangeConf polls a VPS task until it reaches a terminal state.
func vpsSnapshotTaskStateChangeConf(serviceName string, taskID int64, meta interface{}) *retry.StateChangeConf {
	config := meta.(*Config)
	return &retry.StateChangeConf{
		Pending: []string{"waitingAck", "todo", "doing", "paused", "init", "waiting"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {
			resp := &VPSTask{}
			endpoint := fmt.Sprintf("/vps/%s/tasks/%d", url.PathEscape(serviceName), taskID)
			if err := config.OVHClient.Get(endpoint, resp); err != nil {
				return nil, "", err
			}
			switch resp.State {
			case "error", "cancelled", "blocked":
				return resp, resp.State, fmt.Errorf("VPS task %d ended in state %q", resp.Id, resp.State)
			}
			return resp, resp.State, nil
		},
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
