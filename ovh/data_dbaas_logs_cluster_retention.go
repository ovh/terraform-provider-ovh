package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*dbaasLogsClusterRetentionDataSource)(nil)

func NewDbaasLogsClusterRetentionDataSource() datasource.DataSource {
	return &dbaasLogsClusterRetentionDataSource{}
}

type dbaasLogsClusterRetentionDataSource struct {
	config *Config
}

func (d *dbaasLogsClusterRetentionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_logs_cluster_retention"
}

func (d *dbaasLogsClusterRetentionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dbaasLogsClusterRetentionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DbaasLogsClusterRetentionDataSourceSchema(ctx)
}

func (d *dbaasLogsClusterRetentionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DbaasLogsClusterRetentionModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.RetentionId.IsNull() && data.Duration.IsNull() {
		resp.Diagnostics.AddError("missing retention_id or duration",
			"You need to provide retention_id or duration")
		return
	}

	// Retention ID given, fetch it directly
	if !data.RetentionId.IsNull() {
		endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/cluster/" + url.PathEscape(data.ClusterId.ValueString()) + "/retention/" + url.PathEscape(data.RetentionId.ValueString())
		if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &data); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	// No retention ID given, try to fetch a retention with given type and duration
	var (
		retentionIDs       []string
		availableDurations []string
		endpoint           = "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/cluster/" + url.PathEscape(data.ClusterId.ValueString()) + "/retention"
	)

	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &retentionIDs); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error calling get %s", endpoint), err.Error())
		return
	}

	// If no retention_type given, default on LOGS_INDEXING value
	if data.RetentionType.IsNull() {
		tflog.Info(ctx, "no retention type defined, defaulting to LOGS_INDEXING")
		data.RetentionType = ovhtypes.NewTfStringValue("LOGS_INDEXING")
	}

	for _, id := range retentionIDs {
		var (
			retentionData DbaasLogsClusterRetentionModel
			endpoint      = "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/cluster/" + url.PathEscape(data.ClusterId.ValueString()) + "/retention/" + url.PathEscape(id)
		)

		if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &retentionData); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("error calling get %s", endpoint), err.Error())
			return
		}

		if !data.RetentionType.Equal(retentionData.RetentionType) {
			tflog.Info(ctx, fmt.Sprintf("skipping retention %s with wrong type %s", retentionData.RetentionId, retentionData.RetentionType))
			continue
		}

		if !retentionData.IsSupported.ValueBool() {
			tflog.Info(ctx, fmt.Sprintf("skipping retention %s as it is not supported", retentionData.RetentionId))
			continue
		}

		availableDurations = append(availableDurations, retentionData.Duration.ValueString())
		if data.Duration.Equal(retentionData.Duration) {
			data.MergeWith(&retentionData)

			// Save data into Terraform state
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}

	// No retention found with given duration and type, error
	errorDetails := ""
	if len(availableDurations) > 0 {
		errorDetails = ", available values are: " + strings.Join(availableDurations, ",")
	}
	resp.Diagnostics.AddError("retention not found",
		fmt.Sprintf("no retention was found with duration %s and type %s%s", data.Duration.ValueString(), data.RetentionType.ValueString(), errorDetails))
}
