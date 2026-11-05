package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesNodepoolDataSource)(nil)

func NewCloudManagedKubernetesNodepoolDataSource() datasource.DataSource {
	return &cloudManagedKubernetesNodepoolDataSource{}
}

type cloudManagedKubernetesNodepoolDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

func (d *cloudManagedKubernetesNodepoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_nodepool"
}

// Schema is the resource schema made read-only, so both share one model;
// service_name, managed_kubernetes_id and id are the lookup keys.
func (d *cloudManagedKubernetesNodepoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{}
	for k, a := range managedKubernetesNodepoolResourceAttributes() {
		attrs[k] = dataSourceAttribute(a)
	}
	attrs["service_name"] = cloudServiceNameAttribute()
	attrs["managed_kubernetes_id"] = dataSourceRequiredString("ID of the Managed Kubernetes cluster")
	attrs["id"] = dataSourceRequiredString("Nodepool ID")

	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve a nodepool of a Managed Kubernetes cluster.",
		Attributes:  attrs,
	}
}

func (d *cloudManagedKubernetesNodepoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudManagedKubernetesNodepoolModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesNodepoolEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString(), data.Id.ValueString())

	var responseData CloudManagedKubernetesNodepoolAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
