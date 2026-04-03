package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*emailDomainAccountsDataSource)(nil)

func NewEmailDomainAccountsDataSource() datasource.DataSource {
	return &emailDomainAccountsDataSource{}
}

type emailDomainAccountsDataSource struct {
	config *Config
}

func (d *emailDomainAccountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_domain_accounts"
}

func (d *emailDomainAccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *emailDomainAccountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = EmailDomainAccountsDataSourceSchema(ctx)
}

func (d *emailDomainAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EmailDomainAccountsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build query string for optional filters
	query := url.Values{}
	if !data.AccountName.IsNull() && !data.AccountName.IsUnknown() {
		query.Set("accountName", data.AccountName.ValueString())
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		query.Set("description", data.Description.ValueString())
	}

	endpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account"
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	if err := d.config.OVHClient.Get(endpoint, &data.Accounts); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
