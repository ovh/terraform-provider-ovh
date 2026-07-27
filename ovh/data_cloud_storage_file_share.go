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

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileShareDataSource)(nil)

func NewCloudStorageFileShareDataSource() datasource.DataSource {
	return &cloudStorageFileShareDataSource{}
}

type cloudStorageFileShareDataSource struct {
	config *Config
}

func (d *cloudStorageFileShareDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share"
}

func (d *cloudStorageFileShareDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fileShareDataSourceAttributes returns the computed attributes describing a file
// storage share, shared between the singular data source (root attributes) and the
// plural data source (list element attributes).
func fileShareDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "File share ID",
		},
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "File share name",
		},
		"description": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "File share description",
		},
		"size": schema.Int64Attribute{
			Computed:    true,
			Description: "Size of the file share in GB",
		},
		"protocol": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "File share protocol",
		},
		"share_type": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "File share type",
		},
		"location": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Location of the file share",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Region where the file share resides",
				},
				"availability_zone": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Availability zone where the file share resides",
				},
			},
		},
		"share_network_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "ID of the share network the file share is attached to",
		},
		"checksum": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Creation date of the file share",
		},
		"updated_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last update date of the file share",
		},
		"resource_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "File share readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
		},
		"current_state": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current state of the file storage share",
			Attributes: map[string]schema.Attribute{
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
				"name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "File share name",
				},
				"description": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "File share description",
				},
				"size": schema.Int64Attribute{
					Computed:    true,
					Description: "Size of the file share in GB",
				},
				"protocol": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "File share protocol",
				},
				"share_type": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "File share type",
				},
				"share_network_id": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "ID of the share network the file share is attached to",
				},
				"export_locations": schema.ListNestedAttribute{
					Computed:    true,
					Description: "Export locations for the file share",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Export path",
							},
							"preferred": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether this is the preferred export location",
							},
						},
					},
				},
				"capabilities": schema.ListNestedAttribute{
					Computed:    true,
					Description: "Action-availability flags derived from the file share status",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Capability name",
							},
							"enabled": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether the capability is enabled",
							},
							"reason": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Reason why the capability is disabled",
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudStorageFileShareDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			Description:         "File share ID",
			MarkdownDescription: "File share ID",
		},
	}

	// Merge in the shared computed attributes (id is overridden by the Required one above).
	for name, attribute := range fileShareDataSourceAttributes() {
		if name == "id" {
			continue
		}
		attrs[name] = attribute
	}

	resp.Schema = schema.Schema{
		Description:         "Get a file storage share (NFS) in a public cloud project.",
		MarkdownDescription: "Get a file storage share (NFS) in a public cloud project.",
		Attributes:          attrs,
	}
}

// cloudStorageFileShareDataSourceModel is the Terraform state model for this data source.
// It mirrors the resource read shape (flattened target spec + envelope + current_state).
type cloudStorageFileShareDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Description    ovhtypes.TfStringValue `tfsdk:"description"`
	Size           types.Int64            `tfsdk:"size"`
	Protocol       ovhtypes.TfStringValue `tfsdk:"protocol"`
	ShareType      ovhtypes.TfStringValue `tfsdk:"share_type"`
	Location       types.Object           `tfsdk:"location"`
	ShareNetworkId ovhtypes.TfStringValue `tfsdk:"share_network_id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

func (d *cloudStorageFileShareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileShareDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/file/share/" + url.PathEscape(data.Id.ValueString())

	var v CloudStorageFileShareAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	mapFileShareToDataSourceModel(ctx, &v, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapFileShareToDataSourceModel populates the data source model from the API response.
func mapFileShareToDataSourceModel(ctx context.Context, v *CloudStorageFileShareAPIResponse, data *cloudStorageFileShareDataSourceModel) {
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)}
	data.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Checksum)}
	data.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CreatedAt)}
	data.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.UpdatedAt)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)}

	if v.TargetSpec != nil {
		data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Name)}
		data.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Description)}
		data.Size = types.Int64Value(v.TargetSpec.Size)
		data.Protocol = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Protocol)}
		data.ShareType = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.ShareType)}
		if v.TargetSpec.Location != nil {
			data.Location, _ = types.ObjectValue(
				fileShareLocationAttrTypes(),
				map[string]attr.Value{
					"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Location.Region)},
					"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Location.AvailabilityZone)},
				},
			)
		} else {
			data.Location = types.ObjectNull(fileShareLocationAttrTypes())
		}
		if v.TargetSpec.ShareNetwork != nil {
			data.ShareNetworkId = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.ShareNetwork.Id)}
		}
	} else {
		data.Location = types.ObjectNull(fileShareLocationAttrTypes())
	}

	if v.CurrentState != nil {
		data.CurrentState = buildFileShareCurrentStateObject(ctx, v.CurrentState)
	} else {
		data.CurrentState = types.ObjectNull(FileShareCurrentStateAttrTypes())
	}
}
