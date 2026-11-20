package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesCapabilitiesVersionsDataSource)(nil)

func NewCloudManagedKubernetesCapabilitiesVersionsDataSource() datasource.DataSource {
	return &cloudManagedKubernetesCapabilitiesVersionsDataSource{}
}

type cloudManagedKubernetesCapabilitiesVersionsDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

type cloudManagedKubernetesCapabilitiesVersionsModel struct {
	ServiceName         ovhtypes.TfStringValue              `tfsdk:"service_name" json:"-"`
	ManagedKubernetesId ovhtypes.TfStringValue              `tfsdk:"managed_kubernetes_id" json:"-"`
	CurrentVersion      ovhtypes.TfStringValue              `tfsdk:"current_version" json:"currentVersion"`
	AvailableUpdates    []cloudManagedKubernetesVersionItem `tfsdk:"available_updates" json:"availableUpdates"`
	AvailableUpgrades   []cloudManagedKubernetesVersionItem `tfsdk:"available_upgrades" json:"availableUpgrades"`
}

func (d *cloudManagedKubernetesCapabilitiesVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_capabilities_versions"
}

func (d *cloudManagedKubernetesCapabilitiesVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	versionAttributes := map[string]schema.Attribute{
		"version": dataSourceComputedString("Full Kubernetes version (for example `1.33.1-2`)."),
		"status":  dataSourceComputedString("Support status of the version (`SUPPORTED`, `END_OF_SUPPORT`)."),
	}
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the Kubernetes version deployed on a Managed Kubernetes cluster and the versions it can move to.",
		Attributes: map[string]schema.Attribute{
			"service_name":          cloudServiceNameAttribute(),
			"managed_kubernetes_id": dataSourceRequiredString("Cluster ID"),
			"current_version":       dataSourceComputedString("Full Kubernetes version currently deployed on the cluster (for example `1.32.8-1`)."),
			"available_updates":     dataSourceComputedList("Patch versions of the current minor the cluster can update to (for example `1.32.8-1` to `1.32.8-2`).", versionAttributes),
			"available_upgrades":    dataSourceComputedList("Versions of the next minor the cluster can upgrade to (for example `1.32.8-1` to `1.33.1-2`). Multi-minor jumps are never listed.", versionAttributes),
		},
	}
}

func (d *cloudManagedKubernetesCapabilitiesVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudManagedKubernetesCapabilitiesVersionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString()) + "/capabilities/version"
	if err := d.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
