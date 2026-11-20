package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudManagedKubernetesNodepoolResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudManagedKubernetesNodepoolResource)(nil)
	_ resource.ResourceWithImportState = (*cloudManagedKubernetesNodepoolResource)(nil)
)

func NewCloudManagedKubernetesNodepoolResource() resource.Resource {
	return &cloudManagedKubernetesNodepoolResource{}
}

type cloudManagedKubernetesNodepoolResource struct {
	config *Config
}

func (r *cloudManagedKubernetesNodepoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes_nodepool"
}

func (r *cloudManagedKubernetesNodepoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.config = config
}

var managedKubernetesNodepoolMutableAttrs = MutableAttrs{
	Int64s:  []string{"desired_nodes"},
	Objects: []string{"autoscaler", "template"},
}

func (r *cloudManagedKubernetesNodepoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a nodepool in a Managed Kubernetes cluster.",
		Attributes:  managedKubernetesNodepoolResourceAttributes(),
	}
}

func managedKubernetesNodepoolResourceAttributes() map[string]schema.Attribute {
	replace := stringplanmodifier.RequiresReplace()
	replaceIfSet := stringplanmodifier.RequiresReplaceIfConfigured()
	boolReplaceIfSet := boolplanmodifier.RequiresReplaceIfConfigured()

	attrs := map[string]schema.Attribute{
		// Immutable
		"service_name":          requiredString("Service name of the resource representing the id of the cloud project", replace),
		"managed_kubernetes_id": requiredString("ID of the Managed Kubernetes cluster", replace),
		"name":                  optionalString("Nodepool name, generated when omitted", replaceIfSet),
		"availability_zone":     optionalString("Availability zone the nodes are placed in", replaceIfSet),
		"billing_type":          withOneOf(optionalString("Billing period of the nodes: HOURLY or MONTHLY", replaceIfSet), "HOURLY", "MONTHLY"),
		"node_config": schema.SingleNestedAttribute{
			Required:    true,
			Description: "Hardware configuration of the nodes (immutable)",
			Attributes: map[string]schema.Attribute{
				"flavor": requiredString("Node flavor, for example b3-8", replace),
				"attach_floating_ip": optionalObject("Floating IP attachment of the nodes", map[string]schema.Attribute{
					"enabled": optionalBool("Whether a floating IP is attached to each node", boolReplaceIfSet),
				}),
			},
		},
		"anti_affinity": optionalObject("Anti-affinity placement of the nodes (immutable)", map[string]schema.Attribute{
			"enabled": optionalBool("Whether anti-affinity placement is enabled", boolReplaceIfSet),
			"policy":  withOneOf(optionalString("Placement policy: HARD spreads nodes across distinct hypervisors", replaceIfSet), "HARD"),
		}),

		// Mutable
		"desired_nodes": optionalInt64("Desired number of nodes, `0` when omitted. With the autoscaler enabled, must be at least `autoscaler.min_nodes`, and the autoscaler updates it: add it to `lifecycle.ignore_changes`"),
		"autoscaler": optionalObject("Automatic node scaling between min_nodes and max_nodes", map[string]schema.Attribute{
			"enabled":   optionalBool("Whether automatic scaling is enabled"),
			"min_nodes": optionalInt64("Minimum number of nodes"),
			"max_nodes": optionalInt64("Maximum number of nodes"),
			"scale_down": optionalObject("Thresholds governing node removal", map[string]schema.Attribute{
				"unneeded_time_seconds": optionalInt64("Time in seconds a node stays below the utilization threshold before removal"),
				"unready_time_seconds":  optionalInt64("Time in seconds an unready node is kept before removal"),
				"utilization_threshold": optionalNumber("Node utilization ratio, between 0 and 1, below which a node can be removed"),
			}),
		}),
		"template": optionalObject("Kubernetes metadata and scheduling settings applied to every node", map[string]schema.Attribute{
			"metadata": optionalObject("Kubernetes metadata of each node", map[string]schema.Attribute{
				"annotations": optionalStringMap("Annotations of each node"),
				"finalizers":  optionalStringList("Finalizers of each node"),
				"labels":      optionalStringMap("Labels of each node"),
			}),
			"spec": optionalObject("Kubernetes scheduling settings of each node", map[string]schema.Attribute{
				"taints": optionalNestedList("Taints of each node", map[string]schema.Attribute{
					"effect": withOneOf(optionalString("Taint effect: NoExecute, NoSchedule or PreferNoSchedule"), "NoExecute", "NoSchedule", "PreferNoSchedule"),
					"key":    optionalString("Taint key"),
					"value":  optionalString("Taint value"),
				}),
				"unschedulable": optionalBool("Whether the nodes are unschedulable"),
			}),
		}),

		// Computed envelope
		"id":              schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Nodepool ID", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"created_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Creation date", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"checksum":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Computed hash of the target spec", PlanModifiers: []planmodifier.String{UnknownDuringUpdateStringModifier(managedKubernetesNodepoolMutableAttrs)}},
		"updated_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Last update date", PlanModifiers: []planmodifier.String{UnknownDuringUpdateStringModifier(managedKubernetesNodepoolMutableAttrs)}},
		"resource_status": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, MarkdownDescription: "Nodepool readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`…)", PlanModifiers: []planmodifier.String{OutOfSyncPlanModifier()}},
	}

	cs := map[string]schema.Attribute{
		"desired_nodes": computedAttribute(optionalInt64("Observed desired number of nodes")),
		"rollout": computedAttribute(optionalObject("Strategy used to apply node changes", map[string]schema.Attribute{
			"strategy": optionalString("IN_PLACE or ROLLING_UPDATE"),
		})),
		"status": computedAttribute(optionalObject("Size and readiness of the nodes", map[string]schema.Attribute{
			"current_nodes":    optionalInt64("Number of nodes provisioned"),
			"ready_nodes":      optionalInt64("Number of nodes reporting Ready"),
			"up_to_date_nodes": optionalInt64("Number of nodes running the current nodepool spec"),
		})),
	}
	for _, k := range []string{"name", "availability_zone", "billing_type", "node_config", "anti_affinity", "autoscaler", "template"} {
		cs[k] = computedAttribute(attrs[k])
	}
	attrs["current_state"] = schema.SingleNestedAttribute{
		Computed:      true,
		Description:   "Observed state of the nodepool",
		PlanModifiers: []planmodifier.Object{UnknownDuringUpdateObjectModifier(managedKubernetesNodepoolMutableAttrs)},
		Attributes:    cs,
	}
	return attrs
}

// managedKubernetesNodepoolAttrTypes returns the attribute types of every object root attribute.
var managedKubernetesNodepoolAttrTypes = sync.OnceValue(func() map[string]map[string]attr.Type {
	return objectAttrTypes(managedKubernetesNodepoolResourceAttributes())
})

func (r *cloudManagedKubernetesNodepoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<managed_kubernetes_id>/<nodepool_id>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("managed_kubernetes_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func managedKubernetesNodepoolEndpoint(serviceName, clusterId, id string) string {
	endpoint := managedKubernetesEndpoint(serviceName, clusterId) + "/nodepool"
	if id != "" {
		endpoint += "/" + url.PathEscape(id)
	}
	return endpoint
}

func (r *cloudManagedKubernetesNodepoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudManagedKubernetesNodepoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	endpoint := managedKubernetesNodepoolEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString(), "")

	var responseData CloudManagedKubernetesNodepoolAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Save state immediately so the ID is tracked even if the workflow fails.
	data.MergeWith(ctx, &responseData, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	endpoint = managedKubernetesNodepoolEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString(), responseData.Id)
	if err := waitForCloudResourceStatus(ctx, r.config.OVHClient, endpoint, "managed Kubernetes nodepool", responseData.Id, cloudResourceBusy, []string{"READY"}); err != nil {
		resp.Diagnostics.AddError("Error waiting for managed Kubernetes nodepool to be ready", err.Error())
		return
	}

	var fresh CloudManagedKubernetesNodepoolAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &fresh); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &fresh, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedKubernetesNodepoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudManagedKubernetesNodepoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesNodepoolEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString(), data.Id.ValueString())

	var responseData CloudManagedKubernetesNodepoolAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedKubernetesNodepoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudManagedKubernetesNodepoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesNodepoolEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString(), data.Id.ValueString())

	// Refresh the checksum right before PUT to avoid a 409 if server-side drift bumped it.
	var currentData CloudManagedKubernetesNodepoolAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &currentData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	updatePayload := data.ToUpdate(ctx, currentData.Checksum, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseData CloudManagedKubernetesNodepoolAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return
	}

	if err := waitForCloudResourceUpdate(ctx, r.config.OVHClient, endpoint, "managed Kubernetes nodepool", data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error waiting for managed Kubernetes nodepool to be ready after update", err.Error())
		return
	}

	var fresh CloudManagedKubernetesNodepoolAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &fresh); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &fresh, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete treats a 404 as gone, which also covers a cluster deleted first.
func (r *cloudManagedKubernetesNodepoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudManagedKubernetesNodepoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesNodepoolEndpoint(data.ServiceName.ValueString(), data.ManagedKubernetesId.ValueString(), data.Id.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}

	// The nodepool can still report its previous status right after the DELETE call.
	if err := waitForCloudResourceStatus(ctx, r.config.OVHClient, endpoint, "managed Kubernetes nodepool", data.Id.ValueString(), append([]string{"DELETING", "READY"}, cloudResourceBusy...), []string{"DELETED"}); err != nil {
		resp.Diagnostics.AddError("Error waiting for managed Kubernetes nodepool to be deleted", err.Error())
	}
}
