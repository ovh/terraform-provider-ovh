package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudManagedDatabaseLogSubscriptionModel is the Terraform model of the forwarding of one kind of
// logs of a managed database service to a Logs Data Platform stream.
type CloudManagedDatabaseLogSubscriptionModel struct {
	ServiceName types.String `tfsdk:"service_name"`
	ServiceId   types.String `tfsdk:"service_id"`
	Engine      types.String `tfsdk:"engine"`
	Kind        types.String `tfsdk:"kind"`
	StreamId    types.String `tfsdk:"stream_id"`

	Id             types.String `tfsdk:"id"`
	LdpServiceName types.String `tfsdk:"ldp_service_name"`
	Resource       types.Object `tfsdk:"resource"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

// managedDatabaseLogSubscriptionCreation is the dbaas.logs.LogSubscriptionCreation model.
type managedDatabaseLogSubscriptionCreation struct {
	Kind     string `json:"kind"`
	StreamId string `json:"streamId"`
}

// managedDatabaseLogSubscriptionOperation is the dbaas.logs.LogSubscriptionResponse model: the Logs
// Data Platform operation that creates or deletes the subscription.
type managedDatabaseLogSubscriptionOperation struct {
	OperationId string `json:"operationId"`
	ServiceName string `json:"serviceName"`
}

type managedDatabaseLogSubscriptionResource struct {
	Name string `json:"name" tfsdk:"name"`
	Type string `json:"type" tfsdk:"type"`
}

// managedDatabaseLogSubscription is the dbaas.logs.LogSubscription model. Its tfsdk tags map it onto
// the attributes of the log subscription data source.
type managedDatabaseLogSubscription struct {
	CreatedAt      string                                  `json:"createdAt" tfsdk:"created_at"`
	Kind           string                                  `json:"kind" tfsdk:"kind"`
	Resource       *managedDatabaseLogSubscriptionResource `json:"resource" tfsdk:"resource"`
	ServiceName    string                                  `json:"serviceName" tfsdk:"ldp_service_name"`
	StreamId       string                                  `json:"streamId" tfsdk:"stream_id"`
	SubscriptionId string                                  `json:"subscriptionId" tfsdk:"id"`
	UpdatedAt      string                                  `json:"updatedAt" tfsdk:"updated_at"`
}

func managedDatabaseLogSubscriptionResourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"type": types.StringType,
	}
}

func (m *CloudManagedDatabaseLogSubscriptionModel) MergeWith(ctx context.Context, subscription *managedDatabaseLogSubscription) diag.Diagnostics {
	var diags diag.Diagnostics
	m.Id = types.StringValue(subscription.SubscriptionId)
	m.Kind = types.StringValue(subscription.Kind)
	m.StreamId = types.StringValue(subscription.StreamId)
	m.LdpServiceName = types.StringValue(subscription.ServiceName)
	m.CreatedAt = types.StringValue(subscription.CreatedAt)
	m.UpdatedAt = types.StringValue(subscription.UpdatedAt)
	m.Resource = managedDatabaseObjectFrom(ctx, managedDatabaseLogSubscriptionResourceAttrTypes(), subscription.Resource, &diags)
	return diags
}
