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

var _ datasource.DataSourceWithConfigure = (*cloudInstanceBackupsDataSource)(nil)

func NewCloudInstanceBackupsDataSource() datasource.DataSource {
	return &cloudInstanceBackupsDataSource{}
}

type cloudInstanceBackupsDataSource struct {
	config *Config
}

func (d *cloudInstanceBackupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_backups"
}

func (d *cloudInstanceBackupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceBackupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List instance backups for a given instance in a public cloud project.",
		MarkdownDescription: "List instance backups for a given instance in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the instance backups reside",
				MarkdownDescription: "Region where the instance backups reside",
			},
			"instance_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the instance whose backups to list",
				MarkdownDescription: "ID of the instance whose backups to list",
			},
			"backups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of backups for the instance",
				MarkdownDescription: "List of backups for the instance",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Backup ID",
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
				},
			},
		},
	}
}

type cloudInstanceBackupsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	InstanceId  ovhtypes.TfStringValue `tfsdk:"instance_id"`
	Backups     types.List             `tfsdk:"backups"`
}

func instanceBackupListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   ovhtypes.TfStringType{},
		"name": ovhtypes.TfStringType{},
		"location": types.ObjectType{AttrTypes: map[string]attr.Type{
			"region": ovhtypes.TfStringType{},
		}},
		"instance_id":     ovhtypes.TfStringType{},
		"size":            types.Int64Type,
		"status":          ovhtypes.TfStringType{},
		"visibility":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
	}
}

func (d *cloudInstanceBackupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudInstanceBackupsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/compute/backup?instanceId=" + url.QueryEscape(data.InstanceId.ValueString())

	var apiBackups []CloudInstanceBackupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiBackups); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	backupObjs := make([]attr.Value, 0, len(apiBackups))
	for _, b := range apiBackups {
		if b.CurrentState != nil && b.CurrentState.Location != nil &&
			b.CurrentState.Location.Region != data.Region.ValueString() {
			continue
		}

		name := ""
		instanceId := ""
		if b.TargetSpec != nil {
			name = b.TargetSpec.Name
			if b.TargetSpec.Instance != nil {
				instanceId = b.TargetSpec.Instance.Id
			}
		}

		size := int64(0)
		status := ""
		visibility := ""
		region := data.Region.ValueString()
		if b.CurrentState != nil {
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

		itemObj, diags := types.ObjectValue(
			instanceBackupListItemAttrTypes(),
			map[string]attr.Value{
				"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(b.Id)},
				"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(name)},
				"location":        locObj,
				"instance_id":     ovhtypes.TfStringValue{StringValue: types.StringValue(instanceId)},
				"size":            types.Int64Value(size),
				"status":          ovhtypes.TfStringValue{StringValue: types.StringValue(status)},
				"visibility":      ovhtypes.TfStringValue{StringValue: types.StringValue(visibility)},
				"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(b.ResourceStatus)},
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		backupObjs = append(backupObjs, itemObj)
	}

	backupsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: instanceBackupListItemAttrTypes()},
		backupObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Backups = backupsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
