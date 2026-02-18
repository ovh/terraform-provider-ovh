package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectStorageLifecycleConfigurationDataSource)(nil)

func NewCloudProjectStorageLifecycleConfigurationDataSource() datasource.DataSource {
	return &cloudProjectStorageLifecycleConfigurationDataSource{}
}

type cloudProjectStorageLifecycleConfigurationDataSource struct {
	config *Config
}

func (d *cloudProjectStorageLifecycleConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage_object_bucket_lifecycle_configuration"
}

func (d *cloudProjectStorageLifecycleConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectStorageLifecycleConfigurationDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectStorageLifecycleConfigurationDataSourceSchema(ctx)
}

func (d *cloudProjectStorageLifecycleConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectStorageLifecycleConfigurationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() {
		data.ServiceName.StringValue = basetypes.NewStringValue(os.Getenv("OVH_CLOUD_PROJECT_SERVICE"))
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.ContainerName.ValueString()) +
		"/lifecycle"

	var responseData CloudProjectStorageLifecycleConfigurationModel
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Lifecycle configuration not found",
				fmt.Sprintf("No lifecycle configuration found at %s: %s", endpoint, err.Error()),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.ServiceName = data.ServiceName
	responseData.RegionName = data.RegionName
	responseData.ContainerName = data.ContainerName
	responseData.ID = ovhtypes.NewTfStringValue(
		fmt.Sprintf("%s/%s/%s",
			data.ServiceName.ValueString(),
			data.RegionName.ValueString(),
			data.ContainerName.ValueString(),
		),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}
