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

var _ datasource.DataSourceWithConfigure = (*cloudNetworkPrivateVrackDataSource)(nil)

func NewCloudNetworkPrivateVrackDataSource() datasource.DataSource {
	return &cloudNetworkPrivateVrackDataSource{}
}

type cloudNetworkPrivateVrackDataSource struct {
	config *Config
}

func (d *cloudNetworkPrivateVrackDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_network_private_vrack"
}

func (d *cloudNetworkPrivateVrackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudNetworkPrivateVrackDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about a private network (vRack) in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Network ID",
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Network name",
			},
			"location": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Location of the network",
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Region of the network",
					},
				},
			},
			"description": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Network description",
			},
			"vlan_id": schema.Int64Attribute{
				CustomType:  ovhtypes.TfInt64Type{},
				Computed:    true,
				Description: "VLAN ID of the network",
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Computed hash representing the current target specification value",
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation date of the network",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date of the network",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Network readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the network",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Network name",
					},
					"description": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Network description",
					},
					"vlan_id": schema.Int64Attribute{
						CustomType:  ovhtypes.TfInt64Type{},
						Computed:    true,
						Description: "VLAN ID of the network",
					},
					"location": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Location details",
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Region code",
							},
						},
					},
				},
			},
		},
	}
}

// cloudNetworkPrivateVrackDataSourceModel is the Terraform state model for this data source.
// It mirrors the resource model but exposes the location as a nested object
// (the network is fetched by ID, so region is computed, not user input).
type cloudNetworkPrivateVrackDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Location       types.Object           `tfsdk:"location"`
	Description    ovhtypes.TfStringValue `tfsdk:"description"`
	VlanId         ovhtypes.TfInt64Value  `tfsdk:"vlan_id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

func (d *cloudNetworkPrivateVrackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudNetworkPrivateVrackDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudNetworkAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Reuse the resource model's mapping, then project it onto the data source
	// model which exposes the location as a nested computed object.
	var m CloudNetworkPrivateVrackModel
	m.ServiceName = data.ServiceName
	m.Id = data.Id
	m.MergeWith(ctx, &responseData)

	locationVal, diags := types.ObjectValue(
		vrackLocationAttrTypes(),
		map[string]attr.Value{
			"region": m.Region,
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Location = locationVal
	data.Name = m.Name
	data.Description = m.Description
	data.VlanId = m.VlanId
	data.Checksum = m.Checksum
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
	data.ResourceStatus = m.ResourceStatus
	data.CurrentState = m.CurrentState

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
