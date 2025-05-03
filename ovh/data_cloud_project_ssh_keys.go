package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectSshKeysDataSource)(nil)

func NewCloudProjectSshKeysDataSource() datasource.DataSource {
	return &cloudProjectSshKeysDataSource{}
}

type cloudProjectSshKeysDataSource struct {
	config *Config
}

func (d *cloudProjectSshKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_ssh_keys"
}

func (d *cloudProjectSshKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectSshKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectSshKeysDataSourceSchema(ctx)
}

func (d *cloudProjectSshKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectSshKeysModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	query := ""
	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		query = "?region=" + url.QueryEscape(data.Region.ValueString())
	}
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/sshkey" + query

	if err := d.config.OVHClient.Get(endpoint, &data.SshKeys); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
