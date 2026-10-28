package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageSwiftContainerDataSource)(nil)

func NewCloudStorageSwiftContainerDataSource() datasource.DataSource {
	return &cloudStorageSwiftContainerDataSource{}
}

type cloudStorageSwiftContainerDataSource struct {
	config *Config
}

func (d *cloudStorageSwiftContainerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_swift_container"
}

func (d *cloudStorageSwiftContainerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func storageSwiftContainerDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Container identifier, in the <REGION>_<name> form",
		},
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Container name",
		},
		"location": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Region hosting the container",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Region identifier",
				},
			},
		},
		"archive": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether the container uses the Public Cloud Archive storage policy",
		},
		"type": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Access type of the container (PRIVATE, PUBLIC, STATIC)",
		},
		"cors": schema.SetAttribute{
			ElementType: ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Origins allowed to access the container from a browser",
		},
		"checksum": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Creation date of the container",
		},
		"updated_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last update date of the container",
		},
		"resource_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Container readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UNKNOWN, UPDATING)",
		},
		"current_state": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current observed state of the container, read from the Swift backend",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Container name",
				},
				"location": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Region hosting the container",
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region identifier",
						},
					},
				},
				"archive": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether the container uses the Public Cloud Archive storage policy",
				},
				"type": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Access type of the container (PRIVATE, PUBLIC, STATIC)",
				},
				"cors": schema.SetAttribute{
					ElementType: ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Origins currently allowed to access the container from a browser",
				},
				"objects_count": schema.Int64Attribute{
					Computed:    true,
					Description: "Number of objects stored in the container",
				},
				"objects_size": schema.Int64Attribute{
					Computed:    true,
					Description: "Total size of the objects stored in the container, in bytes",
				},
				"static_url": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Public URL of the container, only set when type is PUBLIC or STATIC",
				},
			},
		},
	}
}

func (d *cloudStorageSwiftContainerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			Description:         "Container identifier, in the <REGION>_<name> form",
			MarkdownDescription: "Container identifier, in the `<REGION>_<name>` form",
		},
	}

	for name, attribute := range storageSwiftContainerDataSourceAttributes() {
		if name == "id" {
			continue
		}
		attrs[name] = attribute
	}

	resp.Schema = schema.Schema{
		Description:         "Get a Swift object storage container in a public cloud project.",
		MarkdownDescription: "Get a Swift object storage container in a public cloud project.",
		Attributes:          attrs,
	}
}

type cloudStorageSwiftContainerDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Location       types.Object           `tfsdk:"location"`
	Archive        types.Bool             `tfsdk:"archive"`
	Type           ovhtypes.TfStringValue `tfsdk:"type"`
	Cors           types.Set              `tfsdk:"cors"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

func (d *cloudStorageSwiftContainerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageSwiftContainerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageSwiftContainerEndpoint(data.ServiceName.ValueString()) + "/" + url.PathEscape(data.Id.ValueString())

	var v CloudStorageSwiftContainerAPIResponse
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	mapStorageSwiftContainerToDataSourceModel(ctx, &v, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapStorageSwiftContainerToDataSourceModel(ctx context.Context, v *CloudStorageSwiftContainerAPIResponse, data *cloudStorageSwiftContainerDataSourceModel) {
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)}
	data.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Checksum)}
	data.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CreatedAt)}
	data.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.UpdatedAt)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)}

	if v.TargetSpec != nil {
		data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Name)}
		data.Location = buildStorageSwiftContainerLocationObject(v.TargetSpec.Location)
		archive := false
		if v.TargetSpec.Archive != nil {
			archive = *v.TargetSpec.Archive
		} else if v.CurrentState != nil {
			archive = v.CurrentState.Archive
		}
		data.Archive = types.BoolValue(archive)
		containerType := v.TargetSpec.Type
		if containerType == "" && v.CurrentState != nil {
			containerType = v.CurrentState.Type
		}
		data.Type = nullableTfString(containerType)
		data.Cors = buildStorageSwiftContainerCorsSet(v.TargetSpec.Cors)
	} else {
		data.Location = types.ObjectNull(StorageSwiftContainerLocationAttrTypes())
		data.Archive = types.BoolNull()
		data.Type = nullableTfString("")
		data.Cors = types.SetNull(ovhtypes.TfStringType{})
	}

	if v.CurrentState != nil {
		if v.CurrentState.Name != "" {
			data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CurrentState.Name)}
		}
		data.CurrentState = buildStorageSwiftContainerCurrentStateObject(ctx, v.CurrentState)
	} else {
		data.CurrentState = types.ObjectNull(StorageSwiftContainerCurrentStateAttrTypes())
	}
}
