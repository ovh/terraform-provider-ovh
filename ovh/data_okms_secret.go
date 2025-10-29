package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*okmsSecretDataSource)(nil)

func NewOkmsSecretDataSource() datasource.DataSource {
	return &okmsSecretDataSource{}
}

type okmsSecretDataSource struct {
	config *Config
}

func (d *okmsSecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_okms_secret"
}

func (d *okmsSecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *okmsSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OkmsSecretDataSourceSchema(ctx)
}

func (d *okmsSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var configModel OkmsSecretDataSourceModel

	// Read Terraform configuration into the lightweight DS model
	resp.Diagnostics.Append(req.Config.Get(ctx, &configModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	base := "/v2/okms/resource/" + url.PathEscape(configModel.OkmsId.ValueString()) + "/secret/" + url.PathEscape(configModel.Path.ValueString())
	versionProvided := !configModel.Version.IsNull() && !configModel.Version.IsUnknown() && configModel.Version.ValueInt64() > 0
	includeDataRequested := configModel.IncludeData.ValueBool()

	metaEndpoint := base
	if !versionProvided { // only add includeData for latest case; version endpoint handled separately
		if includeDataRequested {
			metaEndpoint += "?includeData=true"
		}
	}
	var apiModel OkmsSecretModel
	if err := d.config.OVHClient.Get(metaEndpoint, &apiModel); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", metaEndpoint),
			err.Error(),
		)
		return
	}

	configModel.Iam = apiModel.Iam
	configModel.Metadata = apiModel.Metadata

	if versionProvided {
		verEndpoint := base + "/version/" + fmt.Sprintf("%d", configModel.Version.ValueInt64())
		if includeDataRequested {
			verEndpoint += "?includeData=true"
		}
		var ver struct {
			Id        int64           `json:"id"`
			CreatedAt string          `json:"createdAt"`
			Data      json.RawMessage `json:"data"`
			State     string          `json:"state"`
			// deactivatedAt may appear
			DeactivatedAt *string `json:"deactivatedAt"`
		}
		if err := d.config.OVHClient.Get(verEndpoint, &ver); err == nil {
			if includeDataRequested && len(ver.Data) > 0 && string(ver.Data) != "null" {
				configModel.Data = ovhtypes.NewTfStringValue(string(ver.Data))
			}
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", verEndpoint),
				err.Error(),
			)
			return
		}
	} else {
		configModel.Version = apiModel.Metadata.CurrentVersion
		if includeDataRequested {
			if !apiModel.Version.Data.IsNull() && !apiModel.Version.Data.IsUnknown() {
				configModel.Data = apiModel.Version.Data
			} else {
				if !apiModel.Metadata.CurrentVersion.IsNull() && !apiModel.Metadata.CurrentVersion.IsUnknown() && apiModel.Metadata.CurrentVersion.ValueInt64() > 0 {
					verEndpoint := base + "/version/" + fmt.Sprintf("%d", apiModel.Metadata.CurrentVersion.ValueInt64()) + "?includeData=true"
					var ver struct {
						Data json.RawMessage `json:"data"`
					}
					if err := d.config.OVHClient.Get(verEndpoint, &ver); err == nil {
						if len(ver.Data) > 0 && string(ver.Data) != "null" {
							configModel.Data = ovhtypes.NewTfStringValue(string(ver.Data))
						}
					}
				}
			}
		}
	}

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &configModel)...)
}
