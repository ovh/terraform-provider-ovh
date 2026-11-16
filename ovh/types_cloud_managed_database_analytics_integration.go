package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudManagedDatabaseIntegrationModel is the Terraform model of an integration between two managed
// database services.
type CloudManagedDatabaseIntegrationModel struct {
	ServiceName          types.String `tfsdk:"service_name"`
	ServiceId            types.String `tfsdk:"service_id"`
	Engine               types.String `tfsdk:"engine"`
	SourceServiceId      types.String `tfsdk:"source_service_id"`
	DestinationServiceId types.String `tfsdk:"destination_service_id"`
	Type                 types.String `tfsdk:"type"`
	Parameters           types.Map    `tfsdk:"parameters"`

	Id       types.String `tfsdk:"id"`
	Status   types.String `tfsdk:"status"`
	Checksum types.String `tfsdk:"checksum"`
}

// managedDatabaseIntegrationCreation is the publicCloud.database.IntegrationCreation model.
type managedDatabaseIntegrationCreation struct {
	DestinationServiceId string            `json:"destinationServiceId"`
	Parameters           map[string]string `json:"parameters,omitempty"`
	SourceServiceId      string            `json:"sourceServiceId"`
	Type                 string            `json:"type"`
}

// managedDatabaseIntegration is the publicCloud.database.Integration model, without what the API
// doesn't fill in: v1 gives neither the parameters nor the dates, which apiv2 sets to 1970.
type managedDatabaseIntegration struct {
	Checksum             string  `json:"checksum"`
	DestinationServiceId *string `json:"destinationServiceId"`
	Id                   string  `json:"id"`
	SourceServiceId      string  `json:"sourceServiceId"`
	Status               string  `json:"status"`
	Type                 string  `json:"type"`
}

func (m *CloudManagedDatabaseIntegrationModel) ToCreate(ctx context.Context) (*managedDatabaseIntegrationCreation, diag.Diagnostics) {
	var diags diag.Diagnostics
	body := &managedDatabaseIntegrationCreation{
		DestinationServiceId: m.DestinationServiceId.ValueString(),
		SourceServiceId:      m.SourceServiceId.ValueString(),
		Type:                 m.Type.ValueString(),
	}
	if !m.Parameters.IsNull() && !m.Parameters.IsUnknown() {
		diags.Append(m.Parameters.ElementsAs(ctx, &body.Parameters, false)...)
	}
	return body, diags
}

// MergeWith copies the integration into the model, whose parameters it leaves as they are.
func (m *CloudManagedDatabaseIntegrationModel) MergeWith(integration *managedDatabaseIntegration) {
	m.Id = types.StringValue(integration.Id)
	m.SourceServiceId = types.StringValue(integration.SourceServiceId)
	if integration.DestinationServiceId != nil {
		m.DestinationServiceId = types.StringValue(*integration.DestinationServiceId)
	}
	m.Type = types.StringValue(integration.Type)
	m.Status = types.StringValue(integration.Status)
	m.Checksum = types.StringValue(integration.Checksum)
}
