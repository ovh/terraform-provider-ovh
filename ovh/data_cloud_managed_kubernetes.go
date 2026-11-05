package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesDataSource)(nil)

func NewCloudManagedKubernetesDataSource() datasource.DataSource {
	return &cloudManagedKubernetesDataSource{}
}

type cloudManagedKubernetesDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

func (d *cloudManagedKubernetesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes"
}

// Schema is the resource schema made read-only, so both share one model;
// service_name and id are the lookup keys.
func (d *cloudManagedKubernetesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{}
	for k, a := range managedKubernetesResourceAttributes() {
		attrs[k] = dataSourceAttribute(a)
	}
	attrs["service_name"] = cloudServiceNameAttribute()
	attrs["id"] = dataSourceRequiredString("Cluster ID")

	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve a Managed Kubernetes cluster.",
		Attributes:  attrs,
	}
}

func (d *cloudManagedKubernetesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudManagedKubernetesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesEndpoint(data.ServiceName.ValueString(), data.Id.ValueString())

	var responseData CloudManagedKubernetesAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
