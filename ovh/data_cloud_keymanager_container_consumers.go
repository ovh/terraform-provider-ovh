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

var _ datasource.DataSourceWithConfigure = (*cloudKeymanagerContainerConsumersDataSource)(nil)

func NewCloudKeymanagerContainerConsumersDataSource() datasource.DataSource {
	return &cloudKeymanagerContainerConsumersDataSource{}
}

type cloudKeymanagerContainerConsumersDataSource struct {
	config *Config
}

func (d *cloudKeymanagerContainerConsumersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_keymanager_container_consumers"
}

func (d *cloudKeymanagerContainerConsumersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudKeymanagerContainerConsumersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all consumers registered on a Barbican Key Manager container.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"container_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "UUID of the container",
			},
			"consumers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of consumers registered on the container",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed consumer identifier",
						},
						"service": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "OpenStack service type of the consumer",
						},
						"resource_type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Type of the resource consuming the container",
						},
						"resource_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "UUID of the resource consuming the container",
						},
					},
				},
			},
		},
	}
}

func (d *cloudKeymanagerContainerConsumersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudKeymanagerContainerConsumersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/keyManager/container/" + url.PathEscape(data.ContainerId.ValueString()) + "/consumer"

	var apiConsumers []CloudKeymanagerContainerConsumerAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiConsumers); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	var consumerValues []attr.Value
	for _, c := range apiConsumers {
		obj, _ := types.ObjectValue(
			KeymanagerContainerConsumerAttrTypes(),
			map[string]attr.Value{
				"id":            ovhtypes.TfStringValue{StringValue: types.StringValue(c.Id)},
				"service":       ovhtypes.TfStringValue{StringValue: types.StringValue(c.Service)},
				"resource_type": ovhtypes.TfStringValue{StringValue: types.StringValue(c.ResourceType)},
				"resource_id":   ovhtypes.TfStringValue{StringValue: types.StringValue(c.ResourceId)},
			},
		)
		consumerValues = append(consumerValues, obj)
	}

	consumersList, _ := types.ListValue(
		types.ObjectType{AttrTypes: KeymanagerContainerConsumerAttrTypes()},
		consumerValues,
	)
	data.Consumers = consumersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
