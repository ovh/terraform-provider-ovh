package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesVersionsDataSource)(nil)

func NewCloudManagedKubernetesVersionsDataSource() datasource.DataSource {
	return &cloudManagedKubernetesVersionsDataSource{}
}

type cloudManagedKubernetesVersionsDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

type cloudManagedKubernetesVersionsModel struct {
	ServiceName       ovhtypes.TfStringValue              `tfsdk:"service_name" json:"-"`
	DefaultVersion    ovhtypes.TfStringValue              `tfsdk:"default_version" json:"defaultVersion"`
	AvailableVersions []cloudManagedKubernetesVersionItem `tfsdk:"available_versions" json:"availableVersions"`
}

type cloudManagedKubernetesVersionItem struct {
	Version ovhtypes.TfStringValue `tfsdk:"version" json:"version"`
	Status  ovhtypes.TfStringValue `tfsdk:"status" json:"status"`
}

func (d *cloudManagedKubernetesVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_versions"
}

func (d *cloudManagedKubernetesVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the Kubernetes minor versions available for Managed Kubernetes cluster creation, and the default version.",
		Attributes: map[string]schema.Attribute{
			"service_name":    cloudServiceNameAttribute(),
			"default_version": dataSourceComputedString("Minor version applied when a creation request does not specify one (for example `1.34`)."),
			"available_versions": dataSourceComputedList("Minor versions available for cluster creation, newest first. Creation resolves a minor version to its latest patch.", map[string]schema.Attribute{
				"version": dataSourceComputedString("Kubernetes minor version (for example `1.34`)."),
				"status":  dataSourceComputedString("Support status of the version (`SUPPORTED`, `END_OF_SUPPORT`)."),
			}),
		},
	}
}

func (d *cloudManagedKubernetesVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudManagedKubernetesVersionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesReferenceEndpoint(data.ServiceName.ValueString(), "version")
	if err := d.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
