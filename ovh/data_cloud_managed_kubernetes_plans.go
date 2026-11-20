package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedKubernetesPlansDataSource)(nil)

func NewCloudManagedKubernetesPlansDataSource() datasource.DataSource {
	return &cloudManagedKubernetesPlansDataSource{}
}

type cloudManagedKubernetesPlansDataSource struct {
	cloudManagedKubernetesDataSourceBase
}

type cloudManagedKubernetesPlansModel struct {
	ServiceName ovhtypes.TfStringValue           `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue           `tfsdk:"region"`
	Plans       []cloudManagedKubernetesPlanItem `tfsdk:"plans"`
}

type cloudManagedKubernetesPlanItem struct {
	Name         ovhtypes.TfStringValue                  `tfsdk:"name" json:"name"`
	Status       ovhtypes.TfStringValue                  `tfsdk:"status" json:"status"`
	ControlPlane *cloudManagedKubernetesPlanControlPlane `tfsdk:"control_plane" json:"controlPlane"`
	Limits       *cloudManagedKubernetesPlanLimits       `tfsdk:"limits" json:"limits"`
}

type cloudManagedKubernetesPlanControlPlane struct {
	Etcd *cloudManagedKubernetesPlanEtcd `tfsdk:"etcd" json:"etcd"`
}

type cloudManagedKubernetesPlanEtcd struct {
	Dedicated ovhtypes.TfBoolValue                   `tfsdk:"dedicated" json:"dedicated"`
	Storage   *cloudManagedKubernetesPlanEtcdStorage `tfsdk:"storage" json:"storage"`
}

type cloudManagedKubernetesPlanEtcdStorage struct {
	Limit ovhtypes.TfInt64Value  `tfsdk:"limit" json:"limit"`
	Unit  ovhtypes.TfStringValue `tfsdk:"unit" json:"unit"`
}

type cloudManagedKubernetesPlanLimits struct {
	MaxNodes ovhtypes.TfInt64Value `tfsdk:"max_nodes" json:"maxNodes"`
}

func (d *cloudManagedKubernetesPlansDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_plans"
}

func (d *cloudManagedKubernetesPlansDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the Managed Kubernetes plans available for cluster creation in a region, with their control plane characteristics and limits.",
		Attributes: map[string]schema.Attribute{
			"service_name": cloudServiceNameAttribute(),
			"region":       dataSourceRequiredString("Region to list plans for. Plan availability varies by region."),
			"plans": dataSourceComputedList("Plans available in the region.", map[string]schema.Attribute{
				"name":   dataSourceComputedString("Plan name (`FREE`, `STANDARD`)."),
				"status": dataSourceComputedString("Availability of the plan for cluster creation (`AVAILABLE`, `UNAVAILABLE`)."),
				"control_plane": dataSourceComputedObject("Control plane characteristics of the plan.", map[string]schema.Attribute{
					"etcd": dataSourceComputedObject("Etcd characteristics of the control plane.", map[string]schema.Attribute{
						"dedicated": dataSourceComputedBool("Whether etcd runs on storage dedicated to the cluster."),
						"storage": dataSourceComputedObject("Etcd storage quota.", map[string]schema.Attribute{
							"limit": dataSourceComputedInt64("Storage quota, expressed in `unit`."),
							"unit":  dataSourceComputedString("Unit of the storage quota (for example `MiB`, `GiB`)."),
						}),
					}),
				}),
				"limits": dataSourceComputedObject("Limits enforced by the plan.", map[string]schema.Attribute{
					"max_nodes": dataSourceComputedInt64("Maximum number of nodes allowed across the cluster."),
				}),
			}),
		},
	}
}

func (d *cloudManagedKubernetesPlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudManagedKubernetesPlansModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesReferenceEndpoint(data.ServiceName.ValueString(), "plan?region="+url.QueryEscape(data.Region.ValueString()))
	plans, err := helpers.GetAllPagesV2[cloudManagedKubernetesPlanItem](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	// Non-nil so an empty result is stored as [] rather than null.
	if plans == nil {
		plans = []cloudManagedKubernetesPlanItem{}
	}
	data.Plans = plans

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
