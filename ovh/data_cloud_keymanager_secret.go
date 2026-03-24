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

var _ datasource.DataSourceWithConfigure = (*cloudKeymanagerSecretDataSource)(nil)

func NewCloudKeymanagerSecretDataSource() datasource.DataSource {
	return &cloudKeymanagerSecretDataSource{}
}

type cloudKeymanagerSecretDataSource struct {
	config *Config
}

func (d *cloudKeymanagerSecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_keymanager_secret"
}

func (d *cloudKeymanagerSecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudKeymanagerSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get information about a single secret in the Barbican Key Manager service.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"secret_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the secret",
			},

			// Computed
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Secret ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Computed hash representing the current resource state",
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation date of the secret",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date of the secret",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Secret readiness status",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Region of the secret",
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Name of the secret",
			},
			"secret_type": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Type of the secret",
			},
			"metadata": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Key-value metadata for the secret",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the secret",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"secret_type": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"algorithm": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"bit_length": schema.Int64Attribute{
						Computed: true,
					},
					"mode": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"payload_content_type": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"expiration": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"secret_ref": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"status": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"region": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"metadata": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *cloudKeymanagerSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudKeymanagerSecretDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/secret/" + url.PathEscape(data.SecretId.ValueString())

	var responseData CloudKeymanagerSecretAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
