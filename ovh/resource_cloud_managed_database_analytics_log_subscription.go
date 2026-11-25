package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

var (
	_ resource.Resource                = (*cloudManagedDatabaseLogSubscriptionResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudManagedDatabaseLogSubscriptionResource)(nil)
	_ resource.ResourceWithImportState = (*cloudManagedDatabaseLogSubscriptionResource)(nil)
)

func NewCloudManagedDatabaseLogSubscriptionResource() resource.Resource {
	return &cloudManagedDatabaseLogSubscriptionResource{managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"mongodb", "mysql", "postgresql", "valkey"}, suffix: "_log_subscription"}}
}

func NewCloudManagedAnalyticsLogSubscriptionResource() resource.Resource {
	return &cloudManagedDatabaseLogSubscriptionResource{managedDatabaseObjectResource{category: managedAnalyticsFamily, engines: []string{"clickhouse", "grafana", "kafka", "kafkaConnect", "kafkaMirrorMaker", "opensearch"}, suffix: "_log_subscription"}}
}

type cloudManagedDatabaseLogSubscriptionResource struct {
	managedDatabaseObjectResource
}

func (r *cloudManagedDatabaseLogSubscriptionResource) endpoint(data *CloudManagedDatabaseLogSubscriptionModel) string {
	return r.serviceEndpoint(data.ServiceName, data.Engine, data.ServiceId) + "/log/subscription"
}

func (r *cloudManagedDatabaseLogSubscriptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	kept := []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
	attributes := r.objectAttributes("log subscription")
	attributes["kind"] = schema.StringAttribute{
		Required:      true,
		Description:   "Kind of logs to forward. Available kinds are returned by the log/kind endpoint of the service. Changing it recreates the subscription.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
	attributes["stream_id"] = schema.StringAttribute{
		Required:      true,
		Description:   "ID of the Logs Data Platform stream receiving the logs. Changing it recreates the subscription.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
	attributes["ldp_service_name"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Name of the Logs Data Platform service of the stream.",
		PlanModifiers: kept,
	}
	attributes["resource"] = schema.SingleNestedAttribute{
		Computed:      true,
		Description:   "Resource the logs come from.",
		PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Computed: true, Description: "Name of the resource."},
			"type": schema.StringAttribute{Computed: true, Description: "Type of the resource."},
		},
	}
	attributes["created_at"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Creation date of the subscription.",
		PlanModifiers: kept,
	}
	attributes["updated_at"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Last update date of the subscription.",
		PlanModifiers: kept,
	}

	resp.Schema = schema.Schema{
		Description: r.describe("the forwarding of one kind of logs to a Logs Data Platform stream"),
		Attributes:  attributes,
	}
}

func (r *cloudManagedDatabaseLogSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.importState(ctx, req, resp, "id")
}

// Create adopts the subscription of the kind to the stream if there is one, which a creation stopped
// after its POST leaves behind, rather than forwarding the logs twice.
func (r *cloudManagedDatabaseLogSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudManagedDatabaseLogSubscriptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscription, err := r.find(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error looking for an existing log subscription", err.Error())
		return
	}
	if subscription == nil {
		if subscription, err = r.create(ctx, &data); err != nil {
			resp.Diagnostics.AddError("Error creating the log subscription", err.Error())
			// The subscription can exist all the same, for instance when the credentials can't read the
			// operation: saved, it stays tracked.
			if subscription, _ = r.find(ctx, &data); subscription == nil {
				return
			}
		}
	}
	resp.Diagnostics.Append(data.MergeWith(ctx, subscription)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// find returns the subscription of the kind to the stream, or nil.
func (r *cloudManagedDatabaseLogSubscriptionResource) find(ctx context.Context, data *CloudManagedDatabaseLogSubscriptionModel) (*managedDatabaseLogSubscription, error) {
	endpoint := r.endpoint(data)
	list := endpoint + managedDatabaseLogSubscriptionsOfKind(data.Kind.ValueString())
	ids, err := helpers.GetAllPagesV2[string](ctx, r.config.OVHClient, list)
	if err != nil {
		return nil, fmt.Errorf("calling Get %s: %w", list, err)
	}
	for _, id := range ids {
		item := endpoint + "/" + url.PathEscape(id)
		var subscription managedDatabaseLogSubscription
		if err := r.config.OVHClient.GetWithContext(ctx, item, &subscription); err != nil {
			return nil, fmt.Errorf("calling Get %s: %w", item, err)
		}
		if subscription.StreamId == data.StreamId.ValueString() {
			return &subscription, nil
		}
	}
	return nil, nil
}

// create posts the subscription and waits for the Logs Data Platform operation that creates it.
func (r *cloudManagedDatabaseLogSubscriptionResource) create(ctx context.Context, data *CloudManagedDatabaseLogSubscriptionModel) (*managedDatabaseLogSubscription, error) {
	endpoint := r.endpoint(data)
	body := managedDatabaseLogSubscriptionCreation{Kind: data.Kind.ValueString(), StreamId: data.StreamId.ValueString()}
	var operation managedDatabaseLogSubscriptionOperation
	if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseObjectTimeout, func() error {
		return r.config.OVHClient.PostWithContext(ctx, endpoint, body, &operation)
	}); err != nil {
		return nil, fmt.Errorf("calling Post %s: %w", endpoint, err)
	}

	done, err := waitForManagedDatabaseLogsOperation(ctx, r.config.OVHClient, &operation)
	if err != nil {
		return nil, err
	}
	if done.SubscriptionID == nil {
		return nil, fmt.Errorf("operation %s of %s returned no subscription", operation.OperationId, operation.ServiceName)
	}

	endpoint += "/" + url.PathEscape(*done.SubscriptionID)
	var subscription managedDatabaseLogSubscription
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &subscription); err != nil {
		return nil, fmt.Errorf("calling Get %s: %w", endpoint, err)
	}
	return &subscription, nil
}

func (r *cloudManagedDatabaseLogSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudManagedDatabaseLogSubscriptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	var subscription managedDatabaseLogSubscription
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &subscription); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	resp.Diagnostics.Append(data.MergeWith(ctx, &subscription)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is never called: every argument recreates the subscription.
func (r *cloudManagedDatabaseLogSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unexpected update", "Every argument of a log subscription recreates it. Please report this issue to the provider developers.")
}

func (r *cloudManagedDatabaseLogSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudManagedDatabaseLogSubscriptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	var operation managedDatabaseLogSubscriptionOperation
	if err := deleteManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, &operation); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}
	if operation.OperationId == "" {
		return
	}
	if _, err := waitForManagedDatabaseLogsOperation(ctx, r.config.OVHClient, &operation); err != nil {
		resp.Diagnostics.AddError("Error waiting for the log subscription to be deleted", err.Error())
	}
}

// waitForManagedDatabaseLogsOperation polls the Logs Data Platform operation that creates or deletes a
// subscription until it succeeds. A 5xx or a 429 is one more read.
func waitForManagedDatabaseLogsOperation(ctx context.Context, client *ovhwrap.Client, operation *managedDatabaseLogSubscriptionOperation) (*DbaasLogsOperation, error) {
	endpoint := fmt.Sprintf("/dbaas/logs/%s/operation/%s", url.PathEscape(operation.ServiceName), url.PathEscape(operation.OperationId))
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING", "RECEIVED", "STARTED", "RETRY", "RUNNING", managedDatabaseAPIUnavailable},
		Target:  []string{"SUCCESS"},
		Refresh: func() (any, string, error) {
			var res DbaasLogsOperation
			err := client.GetWithContext(ctx, endpoint, &res)
			if isManagedDatabaseAPIUnavailable(err) {
				return &res, managedDatabaseAPIUnavailable, nil
			}
			if err != nil {
				return nil, "", err
			}
			return &res, res.State, nil
		},
		Timeout:    managedDatabaseObjectTimeout,
		MinTimeout: managedDatabaseServicePollMinTimeout,
	}
	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for the operation %s of %s: %w", operation.OperationId, operation.ServiceName, err)
	}
	return res.(*DbaasLogsOperation), nil
}
