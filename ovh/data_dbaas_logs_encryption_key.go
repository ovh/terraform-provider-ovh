package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*dbaasLogsEncryptionKeyDataSource)(nil)

func NewDbaasLogsEncryptionKeyDataSource() datasource.DataSource {
	return &dbaasLogsEncryptionKeyDataSource{}
}

type dbaasLogsEncryptionKeyDataSource struct {
	config *Config
}

func (d *dbaasLogsEncryptionKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_logs_encryption_key"
}

func (d *dbaasLogsEncryptionKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dbaasLogsEncryptionKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DbaasLogsEncryptionKeyDataSourceSchema(ctx)
}

func (d *dbaasLogsEncryptionKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DbaasLogsEncryptionKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.EncryptionKeyId.IsNull() && data.Title.IsNull() {
		resp.Diagnostics.AddError("missing encryption_key_id or title",
			"You need to provide encryption_key_id or title")
		return
	}

	serviceName := data.ServiceName.ValueString()

	// Direct lookup by encryption_key_id
	if !data.EncryptionKeyId.IsNull() {
		endpoint := "/dbaas/logs/" + url.PathEscape(serviceName) + "/encryptionKey/" + url.PathEscape(data.EncryptionKeyId.ValueString())
		if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &data); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	// Lookup by title: list all encryption keys and filter
	var encryptionKeyIDs []string
	endpoint := "/dbaas/logs/" + url.PathEscape(serviceName) + "/encryptionKey"

	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &encryptionKeyIDs); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	titleFilter := data.Title.ValueString()
	var matches []DbaasLogsEncryptionKeyDataSourceModel

	for _, id := range encryptionKeyIDs {
		var keyData DbaasLogsEncryptionKeyDataSourceModel
		keyEndpoint := "/dbaas/logs/" + url.PathEscape(serviceName) + "/encryptionKey/" + url.PathEscape(id)

		if err := d.config.OVHClient.GetWithContext(ctx, keyEndpoint, &keyData); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", keyEndpoint), err.Error())
			return
		}

		if strings.EqualFold(keyData.Title.ValueString(), titleFilter) {
			keyData.ServiceName = ovhtypes.NewTfStringValue(serviceName)
			matches = append(matches, keyData)
		} else {
			tflog.Info(ctx, fmt.Sprintf("skipping encryption key %s with title %s", id, keyData.Title.ValueString()))
		}
	}

	if len(matches) == 0 {
		resp.Diagnostics.AddError("encryption key not found",
			fmt.Sprintf("no encryption key was found with title %q", titleFilter))
		return
	}

	if len(matches) > 1 {
		resp.Diagnostics.AddError("multiple encryption keys found",
			fmt.Sprintf("multiple encryption keys were found with title %q, please use encryption_key_id instead", titleFilter))
		return
	}

	data.MergeWith(&matches[0])

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
