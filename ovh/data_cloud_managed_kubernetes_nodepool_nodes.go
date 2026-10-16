package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesNodepoolNodesDataSource)(nil)

func NewCloudManagedKubernetesNodepoolNodesDataSource() datasource.DataSource {
	return &cloudManagedKubernetesNodepoolNodesDataSource{}
}

type cloudManagedKubernetesNodepoolNodesDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

type cloudManagedKubernetesNodepoolNodesModel struct {
	ServiceName         ovhtypes.TfStringValue                   `tfsdk:"service_name"`
	ManagedKubernetesId ovhtypes.TfStringValue                   `tfsdk:"managed_kubernetes_id"`
	NodepoolId          ovhtypes.TfStringValue                   `tfsdk:"nodepool_id"`
	Nodes               []cloudManagedKubernetesNodepoolNodeItem `tfsdk:"nodes"`
}

type cloudManagedKubernetesNodepoolNodeItem struct {
	Id             ovhtypes.TfStringValue                          `tfsdk:"id" json:"id"`
	CreatedAt      ovhtypes.TfStringValue                          `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt      ovhtypes.TfStringValue                          `tfsdk:"updated_at" json:"updatedAt"`
	ResourceStatus ovhtypes.TfStringValue                          `tfsdk:"resource_status" json:"resourceStatus"`
	CurrentState   *cloudManagedKubernetesNodepoolNodeCurrentState `tfsdk:"current_state" json:"currentState"`
}

type cloudManagedKubernetesNodepoolNodeCurrentState struct {
	Name       ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	NodepoolId ovhtypes.TfStringValue `tfsdk:"nodepool_id" json:"nodePoolId"`
	Flavor     ovhtypes.TfStringValue `tfsdk:"flavor" json:"flavor"`
	InstanceId ovhtypes.TfStringValue `tfsdk:"instance_id" json:"instanceId"`
	Version    ovhtypes.TfStringValue `tfsdk:"version" json:"version"`
	UpToDate   ovhtypes.TfBoolValue   `tfsdk:"up_to_date" json:"upToDate"`
	DeployedAt ovhtypes.TfStringValue `tfsdk:"deployed_at" json:"deployedAt"`
}

func (d *cloudManagedKubernetesNodepoolNodesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_nodepool_nodes"
}

func (d *cloudManagedKubernetesNodepoolNodesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the nodes of a Managed Kubernetes cluster nodepool.",
		Attributes: map[string]schema.Attribute{
			"service_name":          cloudServiceNameAttribute(),
			"managed_kubernetes_id": dataSourceRequiredString("ID of the Managed Kubernetes cluster."),
			"nodepool_id":           dataSourceRequiredString("ID of the nodepool."),
			"nodes": dataSourceComputedList("Nodes of the nodepool.", map[string]schema.Attribute{
				"id":              dataSourceComputedString("Node ID."),
				"created_at":      dataSourceComputedString("Node creation date."),
				"updated_at":      dataSourceComputedString("Last node modification date."),
				"resource_status": dataSourceComputedString("Lifecycle status of the node (`CREATING`, `READY`, `UPDATING`, `OUT_OF_SYNC`, `SUSPENDED`, `ERROR`, `DELETING`, `UNKNOWN`)."),
				"current_state": dataSourceComputedObject("Observed state of the node.", map[string]schema.Attribute{
					"name":        dataSourceComputedString("Node name."),
					"nodepool_id": dataSourceComputedString("ID of the nodepool the node belongs to."),
					"flavor":      dataSourceComputedString("Flavor of the Public Cloud instance backing the node."),
					"instance_id": dataSourceComputedString("ID of the Public Cloud instance backing the node. Absent until the instance is provisioned."),
					"version":     dataSourceComputedString("Kubernetes version installed on the node."),
					"up_to_date":  dataSourceComputedBool("Whether the node runs the current nodepool specification."),
					"deployed_at": dataSourceComputedString("Date the node joined the cluster. Absent until the node has joined."),
				}),
			}),
		},
	}
}

func (d *cloudManagedKubernetesNodepoolNodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudManagedKubernetesNodepoolNodesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesNodepoolEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString(), data.NodepoolId.ValueString()) + "/node"
	nodes, err := helpers.GetAllPagesV2[cloudManagedKubernetesNodepoolNodeItem](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	// Non-nil so an empty result is stored as [] rather than null.
	if nodes == nil {
		nodes = []cloudManagedKubernetesNodepoolNodeItem{}
	}
	data.Nodes = nodes

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
