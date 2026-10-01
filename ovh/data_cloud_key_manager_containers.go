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

var _ datasource.DataSourceWithConfigure = (*cloudKeyManagerContainersDataSource)(nil)

func NewCloudKeyManagerContainersDataSource() datasource.DataSource {
	return &cloudKeyManagerContainersDataSource{}
}

type cloudKeyManagerContainersDataSource struct {
	config *Config
}

func (d *cloudKeyManagerContainersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_key_manager_containers"
}

func (d *cloudKeyManagerContainersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudKeyManagerContainersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all containers in the Barbican Key Manager service for a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"containers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of containers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Container ID",
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current resource state",
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
							Description: "Container readiness status",
						},
						"location": keyManagerLocationDataSourceSchema(),
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Name of the container",
						},
						"type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Type of the container",
						},
						"current_state": keyManagerContainerDataSourceCurrentStateSchema(),
					},
				},
			},
		},
	}
}

func (d *cloudKeyManagerContainersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudKeyManagerContainersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/container"

	var apiContainers []CloudKeyManagerContainerAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiContainers); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	var containerValues []attr.Value
	for _, c := range apiContainers {
		name := ""
		containerType := ""
		locationObj := types.ObjectNull(keyManagerLocationAttrTypes())
		if c.TargetSpec != nil {
			name = c.TargetSpec.Name
			containerType = c.TargetSpec.Type
			locationObj = buildKeyManagerLocationObject(c.TargetSpec.Location)
		}

		var currentStateObj types.Object
		if c.CurrentState != nil {
			currentStateObj = buildKeyManagerContainerCurrentStateObject(ctx, c.CurrentState)
		} else {
			currentStateObj = types.ObjectNull(KeyManagerContainerCurrentStateAttrTypes())
		}

		itemObj, _ := types.ObjectValue(
			KeyManagerContainerListItemAttrTypes(),
			map[string]attr.Value{
				"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(c.Id)},
				"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(c.Checksum)},
				"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(c.CreatedAt)},
				"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(c.UpdatedAt)},
				"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(c.ResourceStatus)},
				"location":        locationObj,
				"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(name)},
				"type":            ovhtypes.TfStringValue{StringValue: types.StringValue(containerType)},
				"current_state":   currentStateObj,
			},
		)
		containerValues = append(containerValues, itemObj)
	}

	containersList, _ := types.ListValue(
		types.ObjectType{AttrTypes: KeyManagerContainerListItemAttrTypes()},
		containerValues,
	)
	data.Containers = containersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
