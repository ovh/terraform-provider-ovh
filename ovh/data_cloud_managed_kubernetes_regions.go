package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesRegionsDataSource)(nil)

func NewCloudManagedKubernetesRegionsDataSource() datasource.DataSource {
	return &cloudManagedKubernetesRegionsDataSource{}
}

type cloudManagedKubernetesRegionsDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

type cloudManagedKubernetesRegionsModel struct {
	ServiceName ovhtypes.TfStringValue             `tfsdk:"service_name"`
	Regions     []cloudManagedKubernetesRegionItem `tfsdk:"regions"`
}

type cloudManagedKubernetesRegionItem struct {
	Name              ovhtypes.TfStringValue   `tfsdk:"name" json:"name"`
	AvailabilityZones []ovhtypes.TfStringValue `tfsdk:"availability_zones" json:"availabilityZones"`
	AvailablePlans    []ovhtypes.TfStringValue `tfsdk:"available_plans" json:"availablePlans"`
}

func (d *cloudManagedKubernetesRegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_regions"
}

func (d *cloudManagedKubernetesRegionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the regions where a Managed Kubernetes cluster can be created, with their availability zones and the plans offered there.",
		Attributes: map[string]schema.Attribute{
			"service_name": cloudServiceNameAttribute(),
			"regions": dataSourceComputedList("Regions available for cluster creation.", map[string]schema.Attribute{
				"name":               dataSourceComputedString("Region identifier (for example `EU-WEST-PAR`)."),
				"availability_zones": dataSourceComputedStringList("Availability zones in the region. Empty for a single-zone region."),
				"available_plans":    dataSourceComputedStringList("Plans that can be selected when creating a cluster in this region."),
			}),
		},
	}
}

func (d *cloudManagedKubernetesRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudManagedKubernetesRegionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesReferenceEndpoint(data.ServiceName.ValueString(), "region")
	regions, err := helpers.GetAllPagesV2[cloudManagedKubernetesRegionItem](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	// Non-nil so an empty result is stored as [] rather than null.
	if regions == nil {
		regions = []cloudManagedKubernetesRegionItem{}
	}
	data.Regions = regions

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
