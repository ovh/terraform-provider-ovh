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

var _ datasource.DataSourceWithConfigure = (*cloudKeyManagerContainerConsumerDataSource)(nil)

func NewCloudKeyManagerContainerConsumerDataSource() datasource.DataSource {
	return &cloudKeyManagerContainerConsumerDataSource{}
}

type cloudKeyManagerContainerConsumerDataSource struct {
	config *Config
}

func (d *cloudKeyManagerContainerConsumerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_key_manager_container_consumer"
}

func (d *cloudKeyManagerContainerConsumerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudKeyManagerContainerConsumerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get information about a single consumer registered on a Barbican Key Manager container.",
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
			"consumer_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Consumer identifier, as returned in the `id` attribute of the `ovh_cloud_key_manager_container_consumer` resource and the `id` field of the `ovh_cloud_key_manager_container_consumers` data source (the base64 identifier expected by the get-one endpoint).",
			},

			// Computed
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
	}
}

func (d *cloudKeyManagerContainerConsumerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudKeyManagerContainerConsumerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/keyManager/container/" + url.PathEscape(data.ContainerId.ValueString()) +
		"/consumer/" + url.PathEscape(data.ConsumerId.ValueString())

	var responseData CloudKeyManagerContainerConsumerAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.Id)}
	data.Service = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.Service)}
	data.ResourceType = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.ResourceType)}
	data.ResourceId = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.ResourceId)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
