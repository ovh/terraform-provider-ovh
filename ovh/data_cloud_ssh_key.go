package ovh

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudSshKeyDataSource)(nil)

func NewCloudSshKeyDataSource() datasource.DataSource {
	return &cloudSshKeyDataSource{}
}

type cloudSshKeyDataSource struct {
	config *Config
}

func (d *cloudSshKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_ssh_key"
}

func (d *cloudSshKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudSshKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudSshKeyDataSourceSchema(ctx)
}

// resolveCloudSshKeyServiceName returns the service name from the given value, falling
// back to the OVH_CLOUD_PROJECT_SERVICE environment variable when it is not set. When the
// environment variable is used, the provided value is updated so it ends up in the state.
func resolveCloudSshKeyServiceName(serviceName *ovhtypes.TfStringValue, diags interface{ AddError(string, string) }) string {
	if !serviceName.IsNull() && !serviceName.IsUnknown() && serviceName.ValueString() != "" {
		return serviceName.ValueString()
	}

	envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
	if envServiceName == "" {
		diags.AddError(
			"Missing service_name",
			"The service_name attribute is required. Please provide it in the configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
		)
		return ""
	}

	*serviceName = ovhtypes.NewTfStringValue(envServiceName)
	return envServiceName
}

func (d *cloudSshKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudSshKeyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := resolveCloudSshKeyServiceName(&data.ServiceName, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := cloudSshKeyEndpoint(serviceName, data.Name.ValueString())

	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &data); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
