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

var _ datasource.DataSourceWithConfigure = (*cloudInstanceBackupDataSource)(nil)

func NewCloudInstanceBackupDataSource() datasource.DataSource {
	return &cloudInstanceBackupDataSource{}
}

type cloudInstanceBackupDataSource struct {
	config *Config
}

func (d *cloudInstanceBackupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_backup"
}

func (d *cloudInstanceBackupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceBackupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get an instance backup in a public cloud project.",
		MarkdownDescription: "Get an instance backup in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Backup ID",
				MarkdownDescription: "Backup ID",
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Backup name",
			},
			"location": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Location of the backup",
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Region",
					},
				},
			},
			"instance_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "ID of the backed-up instance",
			},
			"min_disk": schema.Int64Attribute{
				Computed:    true,
				Description: "Minimum disk size in GB required to boot",
			},
			"min_ram": schema.Int64Attribute{
				Computed:    true,
				Description: "Minimum RAM in MB required to boot",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "Image size in bytes",
			},
			"status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Image status in the backend",
			},
			"visibility": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Image visibility",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Backup readiness status",
			},
		},
	}
}

type cloudInstanceBackupDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Location       types.Object           `tfsdk:"location"`
	InstanceId     ovhtypes.TfStringValue `tfsdk:"instance_id"`
	MinDisk        types.Int64            `tfsdk:"min_disk"`
	MinRam         types.Int64            `tfsdk:"min_ram"`
	Size           types.Int64            `tfsdk:"size"`
	Status         ovhtypes.TfStringValue `tfsdk:"status"`
	Visibility     ovhtypes.TfStringValue `tfsdk:"visibility"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
}

func (d *cloudInstanceBackupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudInstanceBackupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/compute/backup/" + url.PathEscape(data.Id.ValueString())

	var b CloudInstanceBackupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &b); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	name := ""
	instanceId := ""
	region := ""
	if b.TargetSpec != nil {
		name = b.TargetSpec.Name
		if b.TargetSpec.Instance != nil {
			instanceId = b.TargetSpec.Instance.Id
		}
		if b.TargetSpec.Location != nil {
			region = b.TargetSpec.Location.Region
		}
	}

	minDisk := int64(0)
	minRam := int64(0)
	size := int64(0)
	status := ""
	visibility := ""
	if b.CurrentState != nil {
		minDisk = b.CurrentState.MinDisk
		minRam = b.CurrentState.MinRam
		size = b.CurrentState.Size
		status = b.CurrentState.Status
		visibility = b.CurrentState.Visibility
		if b.CurrentState.Instance != nil && b.CurrentState.Instance.Id != "" {
			instanceId = b.CurrentState.Instance.Id
		}
		if b.CurrentState.Name != "" {
			name = b.CurrentState.Name
		}
		if b.CurrentState.Location != nil {
			region = b.CurrentState.Location.Region
		}
	}

	locObj, diags := types.ObjectValue(
		map[string]attr.Type{"region": ovhtypes.TfStringType{}},
		map[string]attr.Value{
			"region": ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(b.Id)}
	data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(name)}
	data.Location = locObj
	data.InstanceId = ovhtypes.TfStringValue{StringValue: types.StringValue(instanceId)}
	data.MinDisk = types.Int64Value(minDisk)
	data.MinRam = types.Int64Value(minRam)
	data.Size = types.Int64Value(size)
	data.Status = ovhtypes.TfStringValue{StringValue: types.StringValue(status)}
	data.Visibility = ovhtypes.TfStringValue{StringValue: types.StringValue(visibility)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(b.ResourceStatus)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
