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

var _ datasource.DataSourceWithConfigure = (*cloudKeyManagerSecretPayloadDataSource)(nil)

func NewCloudKeyManagerSecretPayloadDataSource() datasource.DataSource {
	return &cloudKeyManagerSecretPayloadDataSource{}
}

type cloudKeyManagerSecretPayloadDataSource struct {
	config *Config
}

func (d *cloudKeyManagerSecretPayloadDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_key_manager_secret_payload"
}

func (d *cloudKeyManagerSecretPayloadDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudKeyManagerSecretPayloadDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the payload (secret material) of a Barbican Key Manager secret.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"secret_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "UUID of the secret",
			},

			// Computed
			"payload": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Sensitive:   true,
				Description: "The payload (secret material) of the secret. This value is sensitive.",
			},
		},
	}
}

func (d *cloudKeyManagerSecretPayloadDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudKeyManagerSecretPayloadDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/keyManager/secret/" + url.PathEscape(data.SecretId.ValueString()) + "/payload"

	var responseData CloudKeyManagerSecretPayloadAPIResponse
	if err := d.config.OVHClient.Post(endpoint, nil, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Payload = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.Payload)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
