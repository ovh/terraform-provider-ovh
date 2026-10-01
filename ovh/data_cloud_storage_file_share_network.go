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

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileShareNetworkDataSource)(nil)

func NewCloudStorageFileShareNetworkDataSource() datasource.DataSource {
	return &cloudStorageFileShareNetworkDataSource{}
}

type cloudStorageFileShareNetworkDataSource struct {
	config *Config
}

func (d *cloudStorageFileShareNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share_network"
}

func (d *cloudStorageFileShareNetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fileShareNetworkDataSourceAttributes returns the computed attributes describing a
// file storage share network, shared between the singular and plural data sources.
func fileShareNetworkDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Share network ID",
		},
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Share network name",
		},
		"description": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Share network description",
		},
		"network_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "ID of the private network",
		},
		"subnet_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "ID of the subnet",
		},
		"location": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Location of the share network",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Region where the share network resides",
				},
				"availability_zone": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Availability zone where the share network resides",
				},
			},
		},
		"checksum": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Creation date of the share network",
		},
		"updated_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last update date of the share network",
		},
		"resource_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Share network readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
		},
		"current_state": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current state of the file storage share network",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Share network name",
				},
				"description": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Share network description",
				},
				"network_id": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "ID of the private network",
				},
				"subnet_id": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "ID of the subnet",
				},
				"location": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Current location",
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
	}
}

func (d *cloudStorageFileShareNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name of the resource representing the id of the cloud project",
			MarkdownDescription: "Service name of the resource representing the id of the cloud project",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Share network ID",
			MarkdownDescription: "Share network ID",
		},
	}

	for name, attribute := range fileShareNetworkDataSourceAttributes() {
		if name == "id" {
			continue
		}
		attrs[name] = attribute
	}

	resp.Schema = schema.Schema{
		Description:         "Get a file storage share network in a public cloud project.",
		MarkdownDescription: "Get a file storage share network in a public cloud project.",
		Attributes:          attrs,
	}
}

// cloudStorageFileShareNetworkDataSourceModel is the Terraform state model for this data source.
type cloudStorageFileShareNetworkDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Description    ovhtypes.TfStringValue `tfsdk:"description"`
	NetworkId      ovhtypes.TfStringValue `tfsdk:"network_id"`
	SubnetId       ovhtypes.TfStringValue `tfsdk:"subnet_id"`
	Location       types.Object           `tfsdk:"location"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

func (d *cloudStorageFileShareNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileShareNetworkDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/file/network/" + url.PathEscape(data.Id.ValueString())

	var v CloudStorageFileShareNetworkAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	mapFileShareNetworkToDataSourceModel(&v, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapFileShareNetworkToDataSourceModel populates the data source model from the API response.
func mapFileShareNetworkToDataSourceModel(v *CloudStorageFileShareNetworkAPIResponse, data *cloudStorageFileShareNetworkDataSourceModel) {
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)}
	data.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Checksum)}
	data.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CreatedAt)}
	data.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.UpdatedAt)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)}

	if v.TargetSpec != nil {
		data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Name)}
		data.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Description)}
		if v.TargetSpec.Network != nil {
			data.NetworkId = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Network.Id)}
		}
		if v.TargetSpec.Subnet != nil {
			data.SubnetId = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Subnet.Id)}
		}
		if v.TargetSpec.Location != nil {
			data.Location, _ = types.ObjectValue(
				fileShareNetworkLocationAttrTypes(),
				map[string]attr.Value{
					"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Location.Region)},
					"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Location.AvailabilityZone)},
				},
			)
		} else {
			data.Location = types.ObjectNull(fileShareNetworkLocationAttrTypes())
		}
	} else {
		data.Location = types.ObjectNull(fileShareNetworkLocationAttrTypes())
	}

	if v.CurrentState != nil {
		data.CurrentState = buildFileShareNetworkCurrentStateObject(v.CurrentState)
	} else {
		data.CurrentState = types.ObjectNull(FileShareNetworkCurrentStateAttrTypes())
	}
}
