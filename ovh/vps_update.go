package ovh

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// vpsPrivateInstallationPending is the key of the resource private state that
// holds []byte("true") while the installation requested when the VPS was
// ordered is not applied: the next plan reinstalls the VPS.
const vpsPrivateInstallationPending = "installation_pending"

// vpsPrivateState is the resource private state of the framework, whose type
// is internal to it.
type vpsPrivateState interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// vpsChanges are the changes of a VPS a plan applies through the API v1
// actions; nil fields are left as they are.
type vpsChanges struct {
	DisplayName *string
	Rebuild     *vpsRebuild
	Schedule    *string // automated backup schedule, HH:MM:SS
	PowerState  *string
}

func (c vpsChanges) empty() bool {
	return c.DisplayName == nil && c.Rebuild == nil && c.Schedule == nil && c.PowerState == nil
}

// vpsRebuild is the installation to reinstall a VPS with. ImageId is nil when
// the image comes from ImageName: it is resolved when the reinstallation is
// applied.
type vpsRebuild struct {
	ImageId, ImageName, SshKey, PostInstallScript *string
	SendPasswordByEmail                           bool
}

// vpsV1RebuildBody is the body of POST /vps/{serviceName}/rebuild.
// doNotSendPassword only applies with an SSH key.
type vpsV1RebuildBody struct {
	ImageId           string  `json:"imageId"`
	PublicSshKey      *string `json:"publicSshKey,omitempty"`
	PostInstallScript *string `json:"postInstallScript,omitempty"`
	DoNotSendPassword *bool   `json:"doNotSendPassword,omitempty"`
}

func (b *vpsRebuild) body(imageId string) vpsV1RebuildBody {
	body := vpsV1RebuildBody{ImageId: imageId, PublicSshKey: b.SshKey, PostInstallScript: b.PostInstallScript}
	if b.SshKey != nil {
		doNotSend := !b.SendPasswordByEmail
		body.DoNotSendPassword = &doNotSend
	}
	return body
}

// vpsPlannedChanges returns the changes from state to plan that call the API:
// a field is set only when the known value of the plan differs from the state
// (vpsValueChanged), changes without effect (vpsChangeWithoutEffect) left out.
// The VPS is reinstalled, with the whole installation of the plan, when
// installation.from, ssh_key or post_installation_script change, or when the
// installation requested at the order is pending; its power state is then
// applied too.
func vpsPlannedChanges(plan, state VPSResourceModel, installationPending bool) vpsChanges {
	ctx := context.Background()
	changed := func(path string, plan, state attr.Value) bool {
		p, err := plan.ToTerraformValue(ctx)
		if err != nil {
			return false
		}
		s, err := state.ToTerraformValue(ctx)
		if err != nil {
			return false
		}
		// A plan value unknown with a null configuration is not a change.
		return vpsValueChanged(path, p, s, tftypes.NewValue(p.Type(), nil))
	}

	var changes vpsChanges
	if changed("display_name", plan.DisplayName, state.DisplayName) {
		changes.DisplayName = vpsKnownString(plan.DisplayName)
	}
	if installationPending || changed("installation", plan.Installation, state.Installation) {
		changes.Rebuild = vpsRebuildOf(plan.Installation)
	}
	if changed("options", plan.Options, state.Options) {
		if cron := vpsKnownString(vpsAttributeAt(plan.Options, "auto_backup", "cron")); cron != nil {
			// The schema checks the format.
			if schedule, err := vpsScheduleFromCron(*cron); err == nil {
				changes.Schedule = &schedule
			}
		}
	}
	// A reinstallation starts the VPS: the power state of the plan is applied
	// after it even when it does not change. vpsSetPowerState calls nothing
	// for a VPS already in this state.
	if changes.Rebuild != nil || changed("power_state", plan.PowerState, state.PowerState) {
		changes.PowerState = vpsKnownString(plan.PowerState)
	}
	return changes
}

// vpsRebuildOf returns the reinstallation with installation.
func vpsRebuildOf(installation types.Object) *vpsRebuild {
	send := true
	if value, ok := vpsAttributeAt(installation, "send_password_by_email").(types.Bool); ok && !value.IsNull() && !value.IsUnknown() {
		send = value.ValueBool()
	}
	return &vpsRebuild{
		ImageId:             vpsKnownString(vpsAttributeAt(installation, "from", "image_id")),
		ImageName:           vpsKnownString(vpsAttributeAt(installation, "from", "image_name")),
		SshKey:              vpsKnownString(vpsAttributeAt(installation, "ssh_key")),
		PostInstallScript:   vpsKnownString(vpsAttributeAt(installation, "post_installation_script")),
		SendPasswordByEmail: send,
	}
}

// vpsKnownString returns the string of a known string value, nil otherwise.
func vpsKnownString(value attr.Value) *string {
	s, ok := value.(ovhtypes.TfStringValue)
	if !ok || s.IsNull() || s.IsUnknown() {
		return nil
	}
	v := s.ValueString()
	return &v
}

// vpsRebuildImageId returns the id of the image to reinstall VPS id with:
// rebuild.ImageId, or else the image named rebuild.ImageName. An image VPS id
// cannot be installed from is an error on its attribute, which lists the
// images it can be installed from.
func vpsRebuildImageId(ctx context.Context, client *ovhwrap.Client, id string, rebuild *vpsRebuild) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	from := path.Root("installation").AtName("from")
	notAvailable := func(attribute, image string, names []string) {
		diags.AddAttributeError(from.AtName(attribute),
			"Image not available for this VPS",
			fmt.Sprintf("VPS %s cannot be installed from image %s. Possible images: %s", id, image, strings.Join(names, ", ")),
		)
	}

	switch {
	case rebuild.ImageId != nil:
		images, err := vpsAvailableImages(ctx, client, id)
		if err != nil {
			diags.AddError(fmt.Sprintf("Error reading the images of VPS %s", id), err.Error())
			return "", diags
		}
		names := make([]string, 0, len(images))
		for _, image := range images {
			if image.Id == *rebuild.ImageId {
				return image.Id, diags
			}
			names = append(names, image.Name)
		}
		slices.Sort(names)
		notAvailable("image_id", *rebuild.ImageId, names)
	case rebuild.ImageName != nil:
		imageId, names, err := vpsImageIdForName(ctx, client, id, *rebuild.ImageName)
		if err != nil {
			diags.AddError(fmt.Sprintf("Error reading the images of VPS %s", id), err.Error())
			return "", diags
		}
		if imageId != "" {
			return imageId, diags
		}
		notAvailable("image_name", fmt.Sprintf("%q", *rebuild.ImageName), names)
	default:
		diags.AddAttributeError(from, "No image to install the VPS from", "Set installation.from.image_id or installation.from.image_name.")
	}
	return "", diags
}

// applyChanges applies changes to the VPS of data, which starts from the
// state, in this order: display name, reinstallation, automated backup
// schedule, power state, each waited for before the next one. Each step done
// sets its attributes of plan in data; on an error, the steps after it are
// not done, so that the next plan shows them again. Once all are done, data
// takes every known value of plan, the changes without effect included. A
// reinstallation done removes installation_pending from private. data is then
// read from the VPS.
func (r *vpsResource) applyChanges(ctx context.Context, data *VPSResourceModel, plan VPSResourceModel, changes vpsChanges, deadline time.Time, private vpsPrivateState) diag.Diagnostics {
	diags := r.applySteps(ctx, data, plan, changes, deadline, private)
	if !diags.HasError() {
		data.PlanCode = vpsPlanned(plan.PlanCode, data.PlanCode)
		data.DisplayName = vpsPlanned(plan.DisplayName, data.DisplayName)
		data.PowerState = vpsPlanned(plan.PowerState, data.PowerState)
		data.Location = vpsPlanned(plan.Location, data.Location)
		data.Installation = vpsPlanned(plan.Installation, data.Installation)
		data.Options = vpsPlanned(plan.Options, data.Options)
		data.Order = vpsPlanned(plan.Order, data.Order)
		data.Timeouts = vpsPlanned(plan.Timeouts, data.Timeouts)
	}

	found, refreshDiags := r.refresh(ctx, data)
	diags.Append(refreshDiags...)
	if !refreshDiags.HasError() && !found {
		diags.AddError(fmt.Sprintf("VPS %s not found", data.Id.ValueString()), "The VPS is not found after its update.")
	}
	return diags
}

// vpsPlanned returns the plan value, or current when it is unknown, as the
// computed values a creation plan leaves unknown.
func vpsPlanned[T attr.Value](plan, current T) T {
	if plan.IsUnknown() {
		return current
	}
	return plan
}

// applySteps applies the steps of applyChanges, and stops at the first error.
func (r *vpsResource) applySteps(ctx context.Context, data *VPSResourceModel, plan VPSResourceModel, changes vpsChanges, deadline time.Time, private vpsPrivateState) diag.Diagnostics {
	var diags diag.Diagnostics
	client := r.config.OVHClient
	id := data.Id.ValueString()
	base := "/vps/" + url.PathEscape(id)

	if changes.DisplayName != nil {
		body := struct {
			DisplayName string `json:"displayName"`
		}{*changes.DisplayName}
		if err := client.PutWithContext(ctx, base, body, nil); err != nil {
			diags.AddError(fmt.Sprintf("Error renaming VPS %s", id), fmt.Sprintf("calling PUT %s: %s", base, err))
			return diags
		}
		data.DisplayName = plan.DisplayName
	}

	if changes.Rebuild != nil {
		imageId, imageDiags := vpsRebuildImageId(ctx, client, id, changes.Rebuild)
		diags.Append(imageDiags...)
		if diags.HasError() {
			return diags
		}
		task, err := vpsPostAction(ctx, client, base+"/rebuild", changes.Rebuild.body(imageId), deadline)
		if err == nil {
			err = vpsWaitTask(ctx, client, id, task, deadline)
		}
		if err != nil {
			diags.AddError(fmt.Sprintf("Error reinstalling VPS %s", id), err.Error())
			return diags
		}
		data.Installation = plan.Installation
		diags.Append(private.SetKey(ctx, vpsPrivateInstallationPending, nil)...)
	}

	if changes.Schedule != nil {
		body := struct {
			Schedule string `json:"schedule"`
		}{*changes.Schedule}
		task, err := vpsPostAction(ctx, client, base+"/automatedBackup/reschedule", body, deadline)
		if err == nil {
			err = vpsWaitTask(ctx, client, id, task, deadline)
		}
		if err != nil {
			diags.AddError(fmt.Sprintf("Error scheduling the automated backup of VPS %s", id), err.Error())
			return diags
		}
		data.Options = plan.Options
	}

	if changes.PowerState != nil {
		if err := vpsSetPowerState(ctx, client, id, *changes.PowerState, deadline); err != nil {
			diags.AddError(fmt.Sprintf("Error changing the power state of VPS %s to %s", id, *changes.PowerState), err.Error())
			return diags
		}
		data.PowerState = plan.PowerState
	}
	return diags
}
