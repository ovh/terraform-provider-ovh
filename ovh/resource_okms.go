package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/types"
)

var (
	_ resource.ResourceWithConfigure   = (*okmsResource)(nil)
	_ resource.ResourceWithImportState = (*okmsResource)(nil)
)

func NewOkmsResource() resource.Resource {
	return &okmsResource{}
}

type okmsResource struct {
	config *Config
}

func (r *okmsResource) waitKmsUpdate(ctx context.Context, kmsId string, timeout time.Duration, cb func(*OkmsModel) bool) error {
	var responseData OkmsModel
	// Read updated resource & update corresponding tf resource state
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		// Read updated resource
		if err := r.getOkmsById(kmsId, &responseData); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error reading KMS %s", kmsId))
		}

		if cb(&responseData) {
			return nil
		}

		return retry.RetryableError(errors.New("Waiting okms update"))
	})

	return err
}

func (r *okmsResource) getOkmsById(id string, data *OkmsModel) error {
	endpoint := "/v2/okms/resource/" + url.PathEscape(id) + "?publicCA=true"
	if err := r.config.OVHClient.Get(endpoint, &data); err != nil {
		return err
	}

	data.DisplayName = data.Iam.DisplayName
	return nil
}

func (r *okmsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_okms"
}

func (d *okmsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

func (d *okmsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OkmsResourceSchema(ctx)
}

func (r *okmsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *okmsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OkmsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	order, diags := data.ToCreate(ctx)
	resp.Diagnostics.Append(diags...)
	if order == nil {
		return
	}

	if err := orderCreate(order, r.config, "okms"); err != nil {
		resp.Diagnostics.AddError("failed to create order", err.Error())
	}

	orderID := order.Order.OrderId.ValueInt64()
	plans := []PlanValue{}
	resp.Diagnostics.Append(order.Plan.ElementsAs(ctx, &plans, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Find service name from order
	id, err := serviceNameFromOrder(r.config.OVHClient, orderID, plans[0].PlanCode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to retrieve service name", err.Error())
	}

	data.Id = types.TfStringValue{
		StringValue: basetypes.NewStringValue(id),
	}

	err = r.waitKmsUpdate(ctx, id, 2*time.Minute, func(responseData *OkmsModel) bool {
		// KMS id was updated successfully
		if responseData.Id.ValueString() == id {
			data.MergeWith(responseData)
			return true
		}
		return false
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create KMS, creation timeout", err.Error())
	}

	// Update service displayName
	if serviceUpdateDisplayNameAPIv2(
		r.config,
		data.Id.ValueString(),
		data.DisplayName.ValueString(),
		&resp.Diagnostics,
	) == nil {
		// Update IAM display name as well
		data.Iam.DisplayName = data.DisplayName
	}

	err = r.waitKmsUpdate(ctx, id, 1*time.Minute, func(responseData *OkmsModel) bool {
		// KMS name was updated successfully
		if responseData.DisplayName.ValueString() == data.DisplayName.ValueString() {
			return true
		}
		return false
	})

	if err != nil {
		resp.Diagnostics.AddError("Error getting updated KMS", err.Error())
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OkmsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.getOkmsById(data.Id.ValueString(), &data); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading KMS %s", data.Id.ValueString()),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData OkmsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.OvhSubsidiary.IsUnknown() && (data.OvhSubsidiary.IsUnknown() || data.OvhSubsidiary.IsNull()) {
		// This can happen after an import as the API doesn't return this info
		// This isn't useful once the order is done, so just use what is set in the conf
		data.OvhSubsidiary = planData.OvhSubsidiary
	}

	if planData.DisplayName.ValueString() != data.DisplayName.ValueString() {
		// Update service displayName
		log.Printf("[OKMS] updating display name for %s to %s", data.Id.ValueString(), planData.DisplayName.ValueString())
		if serviceUpdateDisplayNameAPIv2(
			r.config,
			data.Id.ValueString(),
			planData.DisplayName.ValueString(),
			&resp.Diagnostics,
		) == nil {
			log.Printf("[OKMS] update success")
			// Save data into Terraform state
			data.DisplayName = planData.DisplayName
			data.Iam.DisplayName = planData.DisplayName
			err := r.waitKmsUpdate(ctx, data.Id.ValueString(), 1*time.Minute, func(responseData *OkmsModel) bool {
				// KMS name was updated successfully
				if responseData.DisplayName.ValueString() == data.DisplayName.ValueString() {
					return true
				}
				return false
			})

			if err != nil {
				resp.Diagnostics.AddError("Error getting updated KMS", err.Error())
			}

		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OkmsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	serviceId, err := serviceIdFromResourceName(r.config.OVHClient, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error locating service okms %s", id),
			err.Error(),
		)
		return
	}

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate okms %s (service %d)", id, serviceId)
		endpoint := fmt.Sprintf("/services/%d/terminate", serviceId)
		if err := r.config.OVHClient.Post(endpoint, nil, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
				return "", nil
			}
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return id, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of okms %s", id)
		endpoint := fmt.Sprintf("/services/%d/terminate/confirm", serviceId)
		if err := r.config.OVHClient.Post(endpoint, &ConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return nil
	}

	if err := orderDelete(r.config, terminate, confirmTerminate); err != nil {
		resp.Diagnostics.AddError("failed to delete resource", err.Error())
		return
	}
}
