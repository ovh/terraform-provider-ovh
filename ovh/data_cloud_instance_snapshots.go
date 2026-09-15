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

var _ datasource.DataSourceWithConfigure = (*cloudInstanceSnapshotsDataSource)(nil)

func NewCloudInstanceSnapshotsDataSource() datasource.DataSource {
	return &cloudInstanceSnapshotsDataSource{}
}

type cloudInstanceSnapshotsDataSource struct {
	config *Config
}

func (d *cloudInstanceSnapshotsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_snapshots"
}

func (d *cloudInstanceSnapshotsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceSnapshotsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List instance snapshots for a given instance in a public cloud project.",
		MarkdownDescription: "List instance snapshots for a given instance in a public cloud project.",
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
				Description:         "Region where the instance snapshots reside",
				MarkdownDescription: "Region where the instance snapshots reside",
			},
			"instance_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the instance whose snapshots to list",
				MarkdownDescription: "ID of the instance whose snapshots to list",
			},
			"snapshots": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of snapshots for the instance",
				MarkdownDescription: "List of snapshots for the instance",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Snapshot ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Snapshot name",
						},
						"location": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Location of the snapshot",
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
							Description: "ID of the snapshotted instance",
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
							Description: "Snapshot readiness status",
						},
					},
				},
			},
		},
	}
}

type cloudInstanceSnapshotsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	InstanceId  ovhtypes.TfStringValue `tfsdk:"instance_id"`
	Snapshots   types.List             `tfsdk:"snapshots"`
}

func instanceSnapshotListItemAttrTypes() map[string]attr.Type {
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

func (d *cloudInstanceSnapshotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudInstanceSnapshotsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/compute/snapshot?instanceId=" + url.QueryEscape(data.InstanceId.ValueString())

	var apiSnapshots []CloudInstanceSnapshotAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiSnapshots); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	snapshotObjs := make([]attr.Value, 0, len(apiSnapshots))
	for _, b := range apiSnapshots {
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
			instanceSnapshotListItemAttrTypes(),
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

		snapshotObjs = append(snapshotObjs, itemObj)
	}

	snapshotsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: instanceSnapshotListItemAttrTypes()},
		snapshotObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Snapshots = snapshotsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
