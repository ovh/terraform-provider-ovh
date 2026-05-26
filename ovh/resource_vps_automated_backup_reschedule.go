package ovh

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

// vpsAutomatedBackupRescheduleOpts is the body for POST /vps/{sn}/automatedBackup/reschedule.
type vpsAutomatedBackupRescheduleOpts struct {
	Schedule string `json:"schedule"`
}

var hhmmssRegexp = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`)

func resourceVPSAutomatedBackupReschedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPSAutomatedBackupRescheduleCreate,
		ReadContext:   resourceVPSAutomatedBackupRescheduleRead,
		UpdateContext: resourceVPSAutomatedBackupRescheduleUpdate,
		DeleteContext: resourceVPSAutomatedBackupRescheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS",
			},
			"schedule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Backup schedule as an HH:MM:SS time (e.g. 02:00:00)",
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					s, ok := v.(string)
					if !ok {
						errs = append(errs, fmt.Errorf("%q must be a string", k))
						return
					}
					if !hhmmssRegexp.MatchString(s) {
						errs = append(errs, fmt.Errorf("%q must match HH:MM:SS (got %q)", k, s))
					}
					return
				},
			},
			"rotation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of backups retained by the automated backup policy",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the automated backup option (enabled / disabled)",
			},
		},
	}
}

func resourceVPSAutomatedBackupRescheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	schedule := d.Get("schedule").(string)

	if err := postRescheduleAndWait(ctx, config, serviceName, schedule); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceName)

	return resourceVPSAutomatedBackupRescheduleRead(ctx, d, meta)
}

func resourceVPSAutomatedBackupRescheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	if serviceName == "" {
		// Importer path: SetId(serviceName) was called by the import passthrough.
		serviceName = d.Id()
		if err := d.Set("service_name", serviceName); err != nil {
			return diag.FromErr(err)
		}
	}

	endpoint := fmt.Sprintf("/vps/%s/automatedBackup", url.PathEscape(serviceName))
	resp := &VPSAutomatedBackup{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("calling GET %s:\n\t %q", endpoint, err)
	}

	// If the automated backup option is disabled, drop the resource from state -
	// the reschedule endpoint can only retune an already-enabled schedule.
	if resp.State == "disabled" {
		d.SetId("")
		return nil
	}

	if err := d.Set("schedule", resp.Schedule); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rotation", resp.Rotation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", resp.State); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVPSAutomatedBackupRescheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	schedule := d.Get("schedule").(string)

	if err := postRescheduleAndWait(ctx, config, serviceName, schedule); err != nil {
		return diag.FromErr(err)
	}

	return resourceVPSAutomatedBackupRescheduleRead(ctx, d, meta)
}

// Delete is a no-op: the /vps/{sn}/automatedBackup/reschedule endpoint can only
// retune an existing schedule; it cannot disable the automated backup option.
// We simply clear the resource from Terraform state and leave the VPS schedule
// untouched. Users who want to fully disable automated backups must do so via
// the VPS option management endpoints / control panel.
func resourceVPSAutomatedBackupRescheduleDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func postRescheduleAndWait(ctx context.Context, config *Config, serviceName, schedule string) error {
	endpoint := fmt.Sprintf("/vps/%s/automatedBackup/reschedule", url.PathEscape(serviceName))
	body := &vpsAutomatedBackupRescheduleOpts{Schedule: schedule}
	task := &VPSTask{}
	if err := config.OVHClient.Post(endpoint, body, task); err != nil {
		return fmt.Errorf("calling POST %s with %v:\n\t %w", endpoint, body, err)
	}

	if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
		return fmt.Errorf("waiting for VPS task %d on %s: %w", task.Id, serviceName, err)
	}
	return nil
}
