package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudGatewayDataSource)(nil)

func NewCloudGatewayDataSource() datasource.DataSource {
	return &cloudGatewayDataSource{}
}

type cloudGatewayDataSource struct {
	config *Config
}

func (d *cloudGatewayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_gateway"
}

func (d *cloudGatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

func (d *cloudGatewayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about a gateway in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Gateway ID",
			},
			"location": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Location of the gateway",
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Region",
					},
					"availability_zone": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Availability zone",
					},
				},
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Gateway name",
			},
			"description": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Gateway description",
			},
			"external_gateway": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "External gateway configuration",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the external gateway is enabled",
					},
					"model": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "External gateway sizing model",
					},
				},
			},
			"subnet_ids": schema.ListAttribute{
				CustomType:  ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Computed:    true,
				Description: "List of subnet IDs attached as router interfaces",
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Computed hash representing the current target specification value",
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation date of the gateway",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date of the gateway",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Gateway readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the gateway",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Gateway name",
					},
					"description": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Gateway description",
					},
					"status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "OpenStack router status (ACTIVE, BUILD, DOWN, ERROR)",
					},
					"external_gateway": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "External gateway configuration",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether the external gateway is enabled",
							},
							"model": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "External gateway sizing model",
							},
						},
					},
					"external_ip": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "External IP address assigned to the gateway",
					},
					"subnets": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Currently attached subnets",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Subnet ID",
								},
							},
						},
					},
					"location": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Location details",
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Region",
							},
							"availability_zone": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Availability zone",
							},
						},
					},
				},
			},
		},
	}
}

// cloudGatewayDataSourceModel is the Terraform state model for this data source.
// It mirrors the resource model but exposes the location as a nested object
// (the gateway is fetched by ID, so region/AZ are computed, not user input).
type cloudGatewayDataSourceModel struct {
	ServiceName     ovhtypes.TfStringValue                             `tfsdk:"service_name"`
	Id              ovhtypes.TfStringValue                             `tfsdk:"id"`
	Location        types.Object                                       `tfsdk:"location"`
	Name            ovhtypes.TfStringValue                             `tfsdk:"name"`
	Description     ovhtypes.TfStringValue                             `tfsdk:"description"`
	ExternalGateway types.Object                                       `tfsdk:"external_gateway"`
	SubnetIds       ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"subnet_ids"`
	Checksum        ovhtypes.TfStringValue                             `tfsdk:"checksum"`
	CreatedAt       ovhtypes.TfStringValue                             `tfsdk:"created_at"`
	UpdatedAt       ovhtypes.TfStringValue                             `tfsdk:"updated_at"`
	ResourceStatus  ovhtypes.TfStringValue                             `tfsdk:"resource_status"`
	CurrentState    types.Object                                       `tfsdk:"current_state"`
}

func (d *cloudGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudGatewayDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudGatewayAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Reuse the resource model's mapping, then project it onto the data source
	// model which exposes the location as a nested computed object.
	var m CloudGatewayModel
	m.ServiceName = data.ServiceName
	m.Id = data.Id
	m.MergeWith(ctx, &responseData)

	locationVal, diags := types.ObjectValue(
		gatewayLocationAttrTypes(),
		map[string]attr.Value{
			"region":            m.Region,
			"availability_zone": m.AvailabilityZone,
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Location = locationVal
	data.Name = m.Name
	data.Description = m.Description
	data.ExternalGateway = m.ExternalGateway
	data.SubnetIds = m.SubnetIds
	data.Checksum = m.Checksum
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
	data.ResourceStatus = m.ResourceStatus
	data.CurrentState = m.CurrentState

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
