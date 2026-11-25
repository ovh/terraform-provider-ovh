package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

var (
	_ resource.Resource                = (*cloudManagedDatabaseIntegrationResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudManagedDatabaseIntegrationResource)(nil)
	_ resource.ResourceWithImportState = (*cloudManagedDatabaseIntegrationResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*cloudManagedDatabaseIntegrationResource)(nil)

	_ resource.ResourceWithValidateConfig = (*cloudManagedDatabaseIntegrationResource)(nil)
)

// managedDatabaseIntegrationTypes is the publicCloud.database.IntegrationTypeEnum. The pairs of
// engines each type links are in current_state.engine_capabilities.integrations of the service.
var managedDatabaseIntegrationTypes = []string{
	"GRAFANA_DASHBOARD", "GRAFANA_DATASOURCE", "KAFKA_CONNECT", "KAFKA_LOGS", "KAFKA_MIRROR_MAKER",
	"OPENSEARCH_LOGS", "POSTGRESQL_METRICS", "PROMETHEUS", "RSYSLOG",
}

func NewCloudManagedDatabaseIntegrationResource() resource.Resource {
	return &cloudManagedDatabaseIntegrationResource{managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"mysql", "postgresql", "valkey"}, suffix: "_integration"}}
}

func NewCloudManagedAnalyticsIntegrationResource() resource.Resource {
	return &cloudManagedDatabaseIntegrationResource{managedDatabaseObjectResource{category: managedAnalyticsFamily, engines: []string{"clickhouse", "grafana", "kafka", "kafkaConnect", "kafkaMirrorMaker", "opensearch"}, suffix: "_integration"}}
}

type cloudManagedDatabaseIntegrationResource struct {
	managedDatabaseObjectResource
}

func (r *cloudManagedDatabaseIntegrationResource) endpoint(data *CloudManagedDatabaseIntegrationModel) string {
	return r.serviceEndpoint(data.ServiceName, data.Engine, data.ServiceId) + "/integration"
}

func (r *cloudManagedDatabaseIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replace := []planmodifier.String{stringplanmodifier.RequiresReplace()}
	attributes := r.objectAttributes("integration")
	attributes["source_service_id"] = schema.StringAttribute{
		Required:      true,
		Description:   "ID of the source service. Changing it recreates the integration.",
		PlanModifiers: replace,
	}
	attributes["destination_service_id"] = schema.StringAttribute{
		Required:      true,
		Description:   "ID of the destination service. Changing it recreates the integration.",
		PlanModifiers: replace,
	}
	attributes["type"] = schema.StringAttribute{
		Required:      true,
		Description:   "Type of the integration. Changing it recreates the integration.",
		PlanModifiers: replace,
		Validators:    []validator.String{stringvalidator.OneOf(managedDatabaseIntegrationTypes...)},
	}
	attributes["parameters"] = schema.MapAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Description: "Parameters of the integration, which depend on its type. The API doesn't return them: a change made outside of Terraform isn't detected, and an import leaves them unset until the configuration sets them, which only records them. Changing them recreates the integration.",
	}
	attributes["status"] = schema.StringAttribute{
		Computed:    true,
		Description: "Readiness of the integration: CREATING, UPDATING, READY, ERROR or DELETING.",
	}
	attributes["checksum"] = schema.StringAttribute{
		Computed:    true,
		Description: "Hash of the integration, computed by the API.",
	}

	resp.Schema = schema.Schema{
		Description: r.describe("an integration with another service"),
		Attributes:  attributes,
	}
}

// ModifyPlan recreates the integration when its parameters change: the API can't update it. An
// imported integration has none: setting them only records them.
func (r *cloudManagedDatabaseIntegrationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var planned, current types.Map
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("parameters"), &planned)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("parameters"), &current)...)
	if !current.IsNull() && !planned.Equal(current) {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("parameters"))
	}
}

// ValidateConfig checks that the service is the source or the destination of the integration, as v1
// requires.
func (r *cloudManagedDatabaseIntegrationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var service, source, destination types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service_id"), &service)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("source_service_id"), &source)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("destination_service_id"), &destination)...)
	for _, id := range []types.String{service, source, destination} {
		if id.IsNull() || id.IsUnknown() {
			return
		}
	}
	if !service.Equal(source) && !service.Equal(destination) {
		resp.Diagnostics.AddAttributeError(path.Root("service_id"), "Service outside of the integration",
			"service_id must be source_service_id or destination_service_id.")
	}
}

func (r *cloudManagedDatabaseIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.importState(ctx, req, resp, "id")
}

func (r *cloudManagedDatabaseIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudManagedDatabaseIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body, diags := data.ToCreate(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data)
	var integration managedDatabaseIntegration
	if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseObjectTimeout, func() error {
		return r.config.OVHClient.PostWithContext(ctx, endpoint, body, &integration)
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Saved right away, so that the integration stays tracked if the wait fails.
	created := data
	created.MergeWith(&integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &created)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint += "/" + url.PathEscape(integration.Id)
	if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, false, managedDatabaseObjectTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for the integration to be ready", err.Error())
		return
	}
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &integration); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	data.MergeWith(&integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedDatabaseIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudManagedDatabaseIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	var integration managedDatabaseIntegration
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &integration); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	data.MergeWith(&integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update only records the parameters of an imported integration: any other change recreates it.
func (r *cloudManagedDatabaseIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudManagedDatabaseIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	var integration managedDatabaseIntegration
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &integration); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	data.MergeWith(&integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedDatabaseIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudManagedDatabaseIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	if err := deleteManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}
	if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, true, managedDatabaseObjectTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for the integration to be deleted", err.Error())
	}
}
