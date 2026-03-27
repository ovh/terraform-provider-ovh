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

var _ datasource.DataSourceWithConfigure = (*cloudKeymanagerContainersDataSource)(nil)

func NewCloudKeymanagerContainersDataSource() datasource.DataSource {
	return &cloudKeymanagerContainersDataSource{}
}

type cloudKeymanagerContainersDataSource struct {
	config *Config
}

func (d *cloudKeymanagerContainersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_keymanager_containers"
}

func (d *cloudKeymanagerContainersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudKeymanagerContainersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region of the container",
						},
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
						"current_state": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Current state of the container",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"type": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Type of the container. Possible values: CERTIFICATE, GENERIC, RSA",
								},
								"container_ref": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"status": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Status of the container. Possible values: ACTIVE, ERROR",
								},
								"region": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"secret_refs": schema.ListNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												CustomType: ovhtypes.TfStringType{},
												Computed:   true,
											},
											"secret_id": schema.StringAttribute{
												CustomType: ovhtypes.TfStringType{},
												Computed:   true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudKeymanagerContainersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudKeymanagerContainersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/container"

	var apiContainers []CloudKeymanagerContainerAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiContainers); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	var containerValues []attr.Value
	for _, c := range apiContainers {
		region := ""
		name := ""
		containerType := ""
		if c.TargetSpec != nil {
			name = c.TargetSpec.Name
			containerType = c.TargetSpec.Type
			if c.TargetSpec.Location != nil {
				region = c.TargetSpec.Location.Region
			}
		}

		var currentStateObj types.Object
		if c.CurrentState != nil {
			currentStateObj = buildKeymanagerContainerCurrentStateObject(ctx, c.CurrentState)
		} else {
			currentStateObj = types.ObjectNull(KeymanagerContainerCurrentStateAttrTypes())
		}

		itemObj, _ := types.ObjectValue(
			KeymanagerContainerListItemAttrTypes(),
			map[string]attr.Value{
				"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(c.Id)},
				"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(c.Checksum)},
				"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(c.CreatedAt)},
				"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(c.UpdatedAt)},
				"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(c.ResourceStatus)},
				"region":          ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
				"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(name)},
				"type":            ovhtypes.TfStringValue{StringValue: types.StringValue(containerType)},
				"current_state":   currentStateObj,
			},
		)
		containerValues = append(containerValues, itemObj)
	}

	containersList, _ := types.ListValue(
		types.ObjectType{AttrTypes: KeymanagerContainerListItemAttrTypes()},
		containerValues,
	)
	data.Containers = containersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
