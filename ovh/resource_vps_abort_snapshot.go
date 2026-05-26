package ovh

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

// resourceVpsAbortSnapshot is a one-shot Terraform resource that POSTs
// /vps/{service_name}/abortSnapshot to cancel an in-flight snapshot or
// automated-backup operation.
//
// One-shot semantics:
//   - Create: POST the endpoint; the abort is fire-and-forget server-side
//     (no task is returned). The current RFC3339 timestamp is recorded in
//     `aborted_at`.
//   - Read:   no-op (the abort is an event, not a queryable resource).
//   - Update: never invoked — every field is ForceNew, so mutating the
//     `triggers` map triggers destroy+create, which re-runs the abort.
//   - Delete: state-only — there is nothing to undo.
func resourceVpsAbortSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsAbortSnapshotCreate,
		Read:   resourceVpsAbortSnapshotRead,
		Delete: resourceVpsAbortSnapshotDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS service.",
			},
			"triggers": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Map of arbitrary string values; changing any value re-runs the abort.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"aborted_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RFC3339 timestamp of when the abort was issued.",
			},
		},
	}
}

func resourceVpsAbortSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/vps/%s/abortSnapshot",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Post(endpoint, nil, nil); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok {
			switch {
			case apiErr.Code == 404:
				return fmt.Errorf(
					"VPS %q has no snapshot or automated-backup operation to abort (POST %s returned 404).",
					serviceName, endpoint,
				)
			case strings.Contains(strings.ToLower(apiErr.Message), "no operation"),
				strings.Contains(strings.ToLower(apiErr.Message), "not in progress"),
				strings.Contains(strings.ToLower(apiErr.Message), "no task"):
				return fmt.Errorf(
					"VPS %q has no in-flight snapshot or automated-backup operation to abort: %s",
					serviceName, apiErr.Message,
				)
			}
		}
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	d.SetId(fmt.Sprintf("%s/abortSnapshot/%s", serviceName, now))
	_ = d.Set("aborted_at", now)

	return resourceVpsAbortSnapshotRead(d, meta)
}

// resourceVpsAbortSnapshotRead is intentionally a no-op: the abort is an
// event, not a queryable resource. We deliberately do NOT propagate any
// remote 404 here, since there is nothing remote to read.
func resourceVpsAbortSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourceVpsAbortSnapshotDelete drops the resource id. The abort itself
// cannot be undone — the resource only existed to record that the call ran.
func resourceVpsAbortSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
