package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"golang.org/x/net/publicsuffix"

	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.ResourceWithConfigure = (*domainNameResource)(nil)

func NewDomainNameResource() resource.Resource {
	return &domainNameResource{}
}

type domainNameResource struct {
	config *Config
}

func (r *domainNameResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_name"
}

func (d *domainNameResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *domainNameResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DomainNameResourceSchema(ctx)
}

func (d *domainNameResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_name"), req.ID)...)
}

// setDefaultDomainOrderValues fills the plan configuration with a default configuration
// required to order the given domain name
func setDefaultDomainOrderValues(ctx context.Context, order *OrderModel, domain string) {
	var plan PlanValue
	if len(order.Plan.Elements()) > 0 {
		plan = order.Plan.Elements()[0].(PlanValue)
	}

	if plan.PricingMode.IsNull() || plan.PricingMode.IsUnknown() {
		plan.PricingMode = types.NewTfStringValue("create-default")
	}

	if plan.Duration.IsNull() || plan.Duration.IsUnknown() {
		plan.Duration = types.NewTfStringValue("P1Y")
	}

	if plan.PlanCode.IsNull() || plan.PlanCode.IsUnknown() {
		tld, _ := publicsuffix.PublicSuffix(domain)
		plan.PlanCode = types.NewTfStringValue(tld)
	}

	plan.Domain = types.NewTfStringValue(domain)

	order.Plan = types.TfListNestedValue[PlanValue]{
		ListValue: basetypes.NewListValueMust(plan.Type(ctx), []attr.Value{plan}),
	}
}

func (r *domainNameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData DomainNameModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create order and wait for service to be delivered
	order := data.ToOrder()
	setDefaultDomainOrderValues(ctx, order, data.DomainName.ValueString())
	if err := orderCreate(order, r.config, "domain", true); err != nil {
		resp.Diagnostics.AddError("failed to create order", err.Error())
		return
	}

	endpoint := "/v2/domain/name/" + url.PathEscape(data.DomainName.ValueString())

	// Only trigger an update if the target spec is defined
	if !data.TargetSpec.IsNull() && !data.TargetSpec.IsUnknown() {
		// Read resource to get the checksum required to update it
		var initialData DomainNameModel
		if err := r.config.OVHClient.Get(endpoint, &initialData); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		// Update resource
		updateData := data.ToUpdate()
		updateData.Checksum = &initialData.Checksum
		if err := r.config.OVHClient.Put(endpoint, updateData, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Put %s", endpoint),
				err.Error(),
			)
			return
		}
	}

	// Read updated resource
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData, true)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainNameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData DomainNameModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/domain/name/" + url.PathEscape(data.DomainName.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData, true)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainNameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData DomainNameModel

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

	endpoint := "/v2/domain/name/" + url.PathEscape(data.DomainName.ValueString())

	if !planData.TargetSpec.IsNull() && !planData.TargetSpec.IsUnknown() {
		// Set the checksum to its current value
		updateData := planData.ToUpdate()
		updateData.Checksum = &data.Checksum

		// Update resource
		if err := r.config.OVHClient.Put(endpoint, updateData, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Put %s", endpoint),
				err.Error(),
			)
			return
		}
	}

	// Read updated resource
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData, false)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *domainNameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainNameModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.DomainName.ValueString()

	serviceID, err := serviceIdFromRouteAndResourceName(r.config.OVHClient, "/domain/{serviceName}", serviceName)
	if err != nil {
		resp.Diagnostics.AddError("failed to retrieve service ID", err.Error())
	}

	terminate := func() (string, error) {
		log.Printf("[DEBUG] Will terminate service %s", serviceName)
		endpoint := fmt.Sprintf("/services/%d/terminate", serviceID)
		if err := r.config.OVHClient.Post(endpoint, nil, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
				return "", nil
			}
			return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
		return serviceName, nil
	}

	confirmTerminate := func(token string) error {
		log.Printf("[DEBUG] Will confirm termination of service %s", serviceName)
		endpoint := fmt.Sprintf("/services/%d/terminate/confirm", serviceID)
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
