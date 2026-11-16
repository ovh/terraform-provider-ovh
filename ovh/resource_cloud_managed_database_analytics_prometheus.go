package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

// managedDatabasePrometheusKind is the model of a kind of Prometheus endpoint: the attributes every
// endpoint has, and the fields of its engines.
type managedDatabasePrometheusKind[M any] interface {
	*M
	prometheus() *managedDatabasePrometheusModel
	mergeEndpoint(ctx context.Context, endpoint *managedDatabasePrometheusEndpoint) diag.Diagnostics
}

var (
	_ resource.ResourceWithConfigure   = (*managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel])(nil)
	_ resource.ResourceWithImportState = (*managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel])(nil)
	_ resource.ResourceWithModifyPlan  = (*managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel])(nil)
)

func NewCloudManagedDatabasePrometheusResource() resource.Resource {
	return &managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"mysql", "postgresql", "valkey"}, suffix: "_prometheus"},
		fields:                        managedDatabasePrometheusTargetsAttributes(),
	}
}

func NewCloudManagedDatabaseMongodbPrometheusResource() resource.Resource {
	return &managedDatabasePrometheusResource[cloudManagedDatabaseMongodbPrometheusModel, *cloudManagedDatabaseMongodbPrometheusModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"mongodb"}, suffix: "_mongodb_prometheus"},
		fields: map[string]schema.Attribute{
			"srv_domain": schema.StringAttribute{
				Computed:      true,
				Description:   "SRV domain that resolves to the scraping targets.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func NewCloudManagedAnalyticsPrometheusResource() resource.Resource {
	return &managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedAnalyticsFamily, engines: []string{"clickhouse", "kafka", "kafkaConnect", "kafkaMirrorMaker", "opensearch"}, suffix: "_prometheus"},
		fields:                        managedDatabasePrometheusTargetsAttributes(),
	}
}

func managedDatabasePrometheusTargetsAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"targets": schema.ListNestedAttribute{
			Computed:      true,
			Description:   "Scraping targets.",
			PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{Computed: true, Description: "Hostname of the target."},
					"port": schema.Int64Attribute{Computed: true, Description: "Port of the target."},
				},
			},
		},
	}
}

// managedDatabasePrometheusResource has no object of its own in the API: it turns the prometheusEnabled
// flag of the service on and off, and manages the credentials of the endpoint.
type managedDatabasePrometheusResource[M any, P managedDatabasePrometheusKind[M]] struct {
	managedDatabaseObjectResource
	fields map[string]schema.Attribute
}

func (r *managedDatabasePrometheusResource[M, P]) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := r.objectAttributes("Prometheus endpoint")
	attributes["id"] = schema.StringAttribute{
		Computed:      true,
		Description:   "ID of the service.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["password_reset"] = schema.StringAttribute{
		Optional:    true,
		Description: "Arbitrary string to change to reset the password. Removing it keeps the password.",
	}
	attributes["username"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Username for Prometheus authentication.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["password"] = schema.StringAttribute{
		Computed:      true,
		Sensitive:     true,
		Description:   "Password for Prometheus authentication.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["timeouts"] = schema.SingleNestedAttribute{
		Optional:    true,
		Description: "How long to wait for the service, as durations such as 1h30m.",
		Attributes: map[string]schema.Attribute{
			"create": schema.StringAttribute{Optional: true, Description: "Wait for the service to enable the endpoint. Defaults to 40m.", Validators: []validator.String{managedDatabaseDuration}},
			"delete": schema.StringAttribute{Optional: true, Description: "Wait for the service to disable the endpoint. Defaults to 40m.", Validators: []validator.String{managedDatabaseDuration}},
		},
	}
	for name, attribute := range r.fields {
		attributes[name] = attribute
	}
	resp.Schema = schema.Schema{
		Description: r.describe("the Prometheus endpoint, which it enables and whose credentials it manages,"),
		Attributes:  attributes,
	}
}

// ModifyPlan announces a new password only when password_reset resets it, not for a change of
// timeouts.
func (r *managedDatabasePrometheusResource[M, P]) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var planned, current types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("password_reset"), &planned)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("password_reset"), &current)...)
	if managedDatabasePasswordResets(planned, current) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("password"), types.StringUnknown())...)
	}
}

func (r *managedDatabasePrometheusResource[M, P]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.importState(ctx, req, resp, "")
}

func (r *managedDatabasePrometheusResource[M, P]) service(model P) string {
	p := model.prometheus()
	return r.serviceEndpoint(p.ServiceName, p.Engine, p.ServiceId)
}

func (r *managedDatabasePrometheusResource[M, P]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data M
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	model := P(&data)

	service := r.service(model)
	enable := func(service *CloudManagedDatabaseServiceModel) {
		service.PrometheusEnabled = types.BoolValue(true)
	}
	timeout := managedDatabaseTimeout(model.prometheus().Timeouts, "create", managedDatabaseServiceUpdateTimeout)
	if err := putManagedDatabaseServiceWhenReady(ctx, r.config.OVHClient, service, enable, timeout); err != nil {
		resp.Diagnostics.AddError("Error enabling the Prometheus endpoint of the service", err.Error())
		return
	}

	resp.Diagnostics.Append(r.resetPassword(ctx, service, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.readEndpoint(ctx, service, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *managedDatabasePrometheusResource[M, P]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data M
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	model := P(&data)

	service := r.service(model)
	var res CloudManagedDatabaseServiceAPIResponse
	if err := r.config.OVHClient.GetWithContext(ctx, service, &res); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", service), err.Error())
		return
	}
	// Disabled outside of this resource: the next apply enables it again.
	if ts := res.TargetSpec; ts != nil && ts.PrometheusEnabled != nil && !*ts.PrometheusEnabled {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(r.readEndpoint(ctx, service, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update resets the password when password_reset resets it, and otherwise only records the change,
// of timeouts or the removal of password_reset: the other arguments recreate the resource.
func (r *managedDatabasePrometheusResource[M, P]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, prior M
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	model, current := P(&data), P(&prior)

	service := r.service(model)
	if managedDatabasePasswordResets(model.prometheus().PasswordReset, current.prometheus().PasswordReset) {
		resp.Diagnostics.Append(r.resetPassword(ctx, service, model)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if model.prometheus().Password.IsUnknown() {
		// A password_reset unknown in the plan, which could have reset the password, turned out null.
		model.prometheus().Password = current.prometheus().Password
	}
	resp.Diagnostics.Append(r.readEndpoint(ctx, service, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *managedDatabasePrometheusResource[M, P]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data M
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disable := func(service *CloudManagedDatabaseServiceModel) {
		service.PrometheusEnabled = types.BoolValue(false)
	}
	model := P(&data)
	timeout := managedDatabaseTimeout(model.prometheus().Timeouts, "delete", managedDatabaseServiceUpdateTimeout)
	err := putManagedDatabaseServiceWhenReady(ctx, r.config.OVHClient, r.service(model), disable, timeout)
	var apiErr *ovh.APIError
	switch {
	case err == nil || helpers.IsNotFound(err):
	// The API refuses any change to a locked service, which would then block its own destruction.
	case errors.As(err, &apiErr) && apiErr.Code == http.StatusConflict && apiErr.Class == "Client::Conflict::ServiceLocked":
		resp.Diagnostics.AddWarning("Prometheus endpoint left enabled",
			"The service is locked and refuses any change, so its Prometheus endpoint stays enabled. Once the service is unlocked, set prometheus_enabled to false on the service to disable it. "+err.Error())
	default:
		resp.Diagnostics.AddError("Error disabling the Prometheus endpoint of the service", err.Error())
	}
}

func (r *managedDatabasePrometheusResource[M, P]) resetPassword(ctx context.Context, service string, model P) diag.Diagnostics {
	var diags diag.Diagnostics
	endpoint := service + "/prometheus/credentials/reset"
	var access managedDatabasePrometheusAccess
	if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseObjectTimeout, func() error {
		return r.config.OVHClient.PostWithContext(ctx, endpoint, nil, &access)
	}); err != nil {
		diags.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return diags
	}
	model.prometheus().Password = types.StringValue(access.Password)
	return diags
}

func (r *managedDatabasePrometheusResource[M, P]) readEndpoint(ctx context.Context, service string, model P) diag.Diagnostics {
	var diags diag.Diagnostics
	var endpoint managedDatabasePrometheusEndpoint
	if err := r.config.OVHClient.GetWithContext(ctx, service+"/prometheus", &endpoint); err != nil {
		diags.AddError(fmt.Sprintf("Error calling Get %s/prometheus", service), err.Error())
		return diags
	}
	model.prometheus().mergeWith(&endpoint)
	diags.Append(model.mergeEndpoint(ctx, &endpoint)...)
	return diags
}
