package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesKubeconfigDataSource)(nil)

func NewCloudManagedKubernetesKubeconfigDataSource() datasource.DataSource {
	return &cloudManagedKubernetesKubeconfigDataSource{}
}

type cloudManagedKubernetesKubeconfigDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

type cloudManagedKubernetesKubeconfigModel struct {
	ServiceName         ovhtypes.TfStringValue `tfsdk:"service_name" json:"-"`
	ManagedKubernetesId ovhtypes.TfStringValue `tfsdk:"managed_kubernetes_id" json:"-"`
	Content             ovhtypes.TfStringValue `tfsdk:"content" json:"content"`
}

func (d *cloudManagedKubernetesKubeconfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_kubeconfig"
}

func (d *cloudManagedKubernetesKubeconfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve the admin kubeconfig of a Managed Kubernetes cluster.",
		Attributes: map[string]schema.Attribute{
			"service_name":          cloudServiceNameAttribute(),
			"managed_kubernetes_id": dataSourceRequiredString("Cluster ID"),
			"content":               schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Sensitive: true, Description: "Admin kubeconfig of the cluster, in YAML"},
		},
	}
}

func (d *cloudManagedKubernetesKubeconfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudManagedKubernetesKubeconfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString()) + "/kubeconfig"
	if err := d.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
