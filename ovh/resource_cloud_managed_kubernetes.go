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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudManagedKubernetesResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudManagedKubernetesResource)(nil)
	_ resource.ResourceWithImportState = (*cloudManagedKubernetesResource)(nil)
)

func NewCloudManagedKubernetesResource() resource.Resource {
	return &cloudManagedKubernetesResource{}
}

type cloudManagedKubernetesResource struct {
	config *Config
}

func (r *cloudManagedKubernetesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_kubernetes"
}

func (r *cloudManagedKubernetesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var managedKubernetesMutableAttrs = MutableAttrs{
	Strings: []string{"name", "version"},
	Objects: []string{"api_server", "authentication", "auto_upgrade", "cni", "network"},
}

func (r *cloudManagedKubernetesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Managed Kubernetes cluster in a public cloud project.",
		Attributes:  managedKubernetesResourceAttributes(),
	}
}

func ciliumResourceRequirements(desc string) schema.SingleNestedAttribute {
	list := func(d string) schema.SingleNestedAttribute {
		return optionalObject(d, map[string]schema.Attribute{
			"cpu":    optionalString("CPU resource quantity, for example 250m or 1"),
			"memory": optionalString("Memory resource quantity, for example 128Mi or 1Gi"),
		})
	}
	return optionalObject(desc, map[string]schema.Attribute{
		"limits":   list("Maximum compute resources allowed"),
		"requests": list("Minimum compute resources guaranteed"),
	})
}

func kubeProxyDuration(desc string) schema.StringAttribute {
	a := optionalString(desc + ", as an ISO 8601 duration (for example PT60S)")
	a.CustomType = ovhtypes.TfISO8601DurationType{}
	return a
}

func managedKubernetesResourceAttributes() map[string]schema.Attribute {
	replace := stringplanmodifier.RequiresReplace()
	replaceIfSet := stringplanmodifier.RequiresReplaceIfConfigured()

	authorizedIps := optionalObject("IP allowlist for the API server", map[string]schema.Attribute{
		"rules": optionalNestedList("Authorized CIDR blocks", map[string]schema.Attribute{
			"cidr": optionalString("Authorized CIDR block, for example 203.0.113.0/24. Allow-all is 0.0.0.0/0"),
		}),
	})

	attrs := map[string]schema.Attribute{
		// Immutable
		"service_name": requiredString("Service name of the resource representing the id of the cloud project", replace),
		"region":       requiredString("Region of the cluster", replace),
		"plan":         withOneOf(requiredString("Platform tier: FREE or STANDARD", replace), "FREE", "STANDARD"),
		"ip_allocation_policy": optionalObject("Internal pod and service addressing (immutable)", map[string]schema.Attribute{
			"pods":     optionalString("IPv4 CIDR block allocated for pods, for example 10.0.0.0/16", replaceIfSet),
			"services": optionalString("IPv4 CIDR block allocated for services, for example 10.96.0.0/12", replaceIfSet),
		}),

		// Mutable
		"name":    optionalString("Cluster name, generated when omitted"),
		"version": optionalString("Kubernetes version: a minor (1.34) resolves to its latest patch, or a full version (1.34.1-2). Changing it upgrades the cluster"),
		"api_server": optionalObject("Control-plane exposure, IP filtering and admission plugins", map[string]schema.Attribute{
			"access": optionalObject("Exposure and IP filtering of the API server", map[string]schema.Attribute{
				"authorized_ips": authorizedIps,
			}),
			"admission_plugins": optionalNestedMap("Admission plugin configurations, keyed by plugin name", map[string]schema.Attribute{
				"enabled": optionalBool("Whether the admission plugin is enabled"),
			}),
		}),
		"authentication": withVariants(optionalObject("Control-plane authentication", map[string]schema.Attribute{
			"type": withOneOf(optionalString("Authentication variant: OIDC"), "OIDC"),
			"oidc": optionalObject("OpenID Connect provider, set when type is OIDC", map[string]schema.Attribute{
				"client_id":       optionalString("OIDC client identifier"),
				"issuer_url":      optionalString("OIDC issuer URL"),
				"ca_content":      optionalString("CA certificate of the OIDC provider, PEM-encoded"),
				"groups_claim":    optionalStringList("JWT claims used as Kubernetes groups"),
				"groups_prefix":   optionalString("Prefix prepended to group claims"),
				"required_claim":  optionalStringList("Claims that must be present in the JWT token"),
				"signing_algs":    optionalStringList("Accepted signing algorithms (ES256, ES384, ES512, PS256, PS384, PS512, RS256, RS384, RS512)"),
				"username_claim":  optionalString("JWT claim used as the Kubernetes username"),
				"username_prefix": optionalString("Prefix prepended to username claims"),
			}),
		}), "type", "oidc"),
		"auto_upgrade": optionalObject("Automatic Kubernetes version upgrade policy", map[string]schema.Attribute{
			"patch": optionalObject("Automatic upgrade of patch versions (enabled by default)", map[string]schema.Attribute{
				"enabled": optionalBool("Whether automatic patch upgrade is enabled"),
			}),
		}),

		// Required, partially mutable: the variant (type) forces a replacement
		"cni": withVariants(schema.SingleNestedAttribute{
			Required:    true,
			Description: "Container Network Interface plugin",
			Attributes: map[string]schema.Attribute{
				"type": withOneOf(requiredString("CNI variant: CANAL (FREE plan) or CILIUM (STANDARD plan)", replace), "CANAL", "CILIUM"),
				"canal": optionalObject("Canal settings, set when type is CANAL", map[string]schema.Attribute{
					"kube_proxy": withVariants(optionalObject("kube-proxy configuration", map[string]schema.Attribute{
						"mode": withOneOf(optionalString("kube-proxy mode: IPTABLES or IPVS (immutable)", replaceIfSet), "IPTABLES", "IPVS"),
						"iptables": optionalObject("iptables settings, set when mode is IPTABLES", map[string]schema.Attribute{
							"min_sync_period": kubeProxyDuration("Minimum period between iptables rules refreshes"),
							"sync_period":     kubeProxyDuration("Period between iptables rules refreshes"),
						}),
						"ipvs": optionalObject("IPVS settings, set when mode is IPVS", map[string]schema.Attribute{
							"min_sync_period": kubeProxyDuration("Minimum period between IPVS rules refreshes"),
							"sync_period":     kubeProxyDuration("Period between IPVS rules refreshes"),
							"scheduler": withOneOf(optionalString("IPVS scheduling algorithm"),
								"DESTINATION_HASHING", "LEAST_CONNECTION", "NEVER_QUEUE", "ROUND_ROBIN", "SHORTEST_EXPECTED_DELAY", "SOURCE_HASHING"),
							"tcp_fin_timeout": kubeProxyDuration("Timeout of IPVS TCP sessions after a FIN"),
							"tcp_timeout":     kubeProxyDuration("Timeout of idle IPVS TCP sessions"),
							"udp_timeout":     kubeProxyDuration("Timeout of IPVS UDP packets"),
						}),
					}), "mode", "iptables", "ipvs"),
				}),
				"cilium": optionalObject("Cilium settings, set when type is CILIUM", map[string]schema.Attribute{
					"cluster_id": optionalInt64("ClusterMesh cluster identifier, between 1 and 255"),
					"cluster_mesh": optionalObject("ClusterMesh multi-cluster connectivity", map[string]schema.Attribute{
						"enabled": optionalBool("Whether ClusterMesh is enabled"),
						"api_server": optionalObject("ClusterMesh API server exposure", map[string]schema.Attribute{
							"node_port":    optionalInt64("Node port exposing the ClusterMesh API server"),
							"service_type": withOneOf(optionalString("Service type exposing the ClusterMesh API server"), "LoadBalancer", "NodePort"),
						}),
					}),
					"hubble": optionalObject("Hubble observability", map[string]schema.Attribute{
						"enabled": optionalBool("Whether Hubble is enabled"),
						"relay": optionalObject("Hubble relay", map[string]schema.Attribute{
							"enabled": optionalBool("Whether the Hubble relay is enabled"),
						}),
						"ui": optionalObject("Hubble UI", map[string]schema.Attribute{
							"enabled":            optionalBool("Whether the Hubble UI is enabled"),
							"backend_resources":  ciliumResourceRequirements("Compute resources of the Hubble UI backend"),
							"frontend_resources": ciliumResourceRequirements("Compute resources of the Hubble UI frontend"),
						}),
					}),
				}),
			},
		}, "type", "canal", "cilium"),
		"network": withVariants(schema.SingleNestedAttribute{
			Required:            true,
			MarkdownDescription: "Network topology of the nodes. The variant and its identifiers are immutable; `default_egress` is mutable",
			Attributes: map[string]schema.Attribute{
				"type": withOneOf(requiredString("Network variant: PUBLIC or VRACK", replace), "PUBLIC", "VRACK"),
				"vrack": optionalObject("Private network (vRack) settings, required when type is VRACK", map[string]schema.Attribute{
					"private_network_id":      optionalString("Private network ID", replaceIfSet),
					"subnet_id":               optionalString("Subnet ID of the nodes", replaceIfSet),
					"load_balancer_subnet_id": optionalString("Subnet ID of the load balancers", replaceIfSet),
					"default_egress": withVariants(optionalObject("Default egress route of the nodes. Defaults to MANAGED_GATEWAY on STANDARD, NODE_PUBLIC_INTERFACE on FREE", map[string]schema.Attribute{
						"type": withOneOf(optionalString("Egress route type"), "MANAGED_GATEWAY", "NODE_PUBLIC_INTERFACE", "USER_GATEWAY"),
						"user_gateway": optionalObject("User-provided egress gateway, required when type is USER_GATEWAY", map[string]schema.Attribute{
							"ip": optionalString("IP address of the gateway"),
						}),
					}), "type", "user_gateway"),
				}),
			},
		}, "type", "vrack"),

		// Computed envelope
		"id":              schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Cluster ID", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"created_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Creation date", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"checksum":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Computed hash of the target spec", PlanModifiers: []planmodifier.String{UnknownDuringUpdateStringModifier(managedKubernetesMutableAttrs)}},
		"updated_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Last update date", PlanModifiers: []planmodifier.String{UnknownDuringUpdateStringModifier(managedKubernetesMutableAttrs)}},
		"resource_status": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, MarkdownDescription: "Cluster readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`…)", PlanModifiers: []planmodifier.String{OutOfSyncPlanModifier()}},
	}

	attrs["current_state"] = schema.SingleNestedAttribute{
		Computed:      true,
		Description:   "Observed state of the cluster",
		PlanModifiers: []planmodifier.Object{UnknownDuringUpdateObjectModifier(managedKubernetesMutableAttrs)},
		Attributes:    managedKubernetesCurrentStateAttributes(attrs),
	}
	return attrs
}

// managedKubernetesCurrentStateAttributes mirrors the target spec as Computed,
// plus the observed-only fields.
func managedKubernetesCurrentStateAttributes(target map[string]schema.Attribute) map[string]schema.Attribute {
	cs := map[string]schema.Attribute{
		"deployed_version": computedAttribute(optionalString("Full Kubernetes version deployed, for example 1.34.1-2")),
	}
	for _, k := range []string{"name", "region", "plan", "version", "api_server", "authentication", "auto_upgrade", "cni", "ip_allocation_policy", "network"} {
		cs[k] = computedAttribute(target[k])
	}

	apiServer := cs["api_server"].(schema.SingleNestedAttribute)
	access := apiServer.Attributes["access"].(schema.SingleNestedAttribute)
	access.Attributes["public_endpoint"] = computedAttribute(optionalObject("Public endpoint of the API server", map[string]schema.Attribute{
		"url": optionalString("URL of the control plane, available once the cluster is ready"),
	}))
	return cs
}

// managedKubernetesAttrTypes returns the attribute types of every object root attribute.
var managedKubernetesAttrTypes = sync.OnceValue(func() map[string]map[string]attr.Type {
	return objectAttrTypes(managedKubernetesResourceAttributes())
})

func (r *cloudManagedKubernetesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<cluster_id>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func managedKubernetesEndpoint(serviceName, id string) string {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/managedKubernetes"
	if id != "" {
		endpoint += "/" + url.PathEscape(id)
	}
	return endpoint
}

func (r *cloudManagedKubernetesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudManagedKubernetesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	endpoint := managedKubernetesEndpoint(data.ServiceName.ValueString(), "")

	var responseData CloudManagedKubernetesAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Save state immediately so the ID is tracked even if the workflow fails.
	data.MergeWith(ctx, &responseData, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	endpoint = managedKubernetesEndpoint(data.ServiceName.ValueString(), responseData.Id)
	if err := waitForCloudResourceStatus(ctx, r.config.OVHClient, endpoint, "managed Kubernetes cluster", responseData.Id, cloudResourceBusy, []string{"READY"}); err != nil {
		resp.Diagnostics.AddError("Error waiting for managed Kubernetes cluster to be ready", err.Error())
		return
	}

	var fresh CloudManagedKubernetesAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &fresh); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &fresh, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedKubernetesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudManagedKubernetesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesEndpoint(data.ServiceName.ValueString(), data.Id.ValueString())

	var responseData CloudManagedKubernetesAPIResponse
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

func (r *cloudManagedKubernetesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudManagedKubernetesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesEndpoint(data.ServiceName.ValueString(), data.Id.ValueString())

	// Refresh the checksum right before PUT to avoid a 409 if server-side drift bumped it.
	var currentData CloudManagedKubernetesAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &currentData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	updatePayload := data.ToUpdate(ctx, currentData.Checksum, &resp.Diagnostics)
	var priorVersion ovhtypes.TfStringValue
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("version"), &priorVersion)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Unchanged version: send its minor, a no-op for the API. Resending a full version fails
	// once auto-upgrade has moved past it (declared 1.35.2-9, deployed 1.35.3-1 → 400).
	if !data.Version.IsNull() && data.Version.Equal(priorVersion) {
		updatePayload.TargetSpec.Version = ovhtypes.NewTfStringValue(minorVersion(data.Version.ValueString()))
	}

	var responseData CloudManagedKubernetesAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return
	}

	if err := waitForCloudResourceUpdate(ctx, r.config.OVHClient, endpoint, "managed Kubernetes cluster", data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error waiting for managed Kubernetes cluster to be ready after update", err.Error())
		return
	}

	var fresh CloudManagedKubernetesAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &fresh); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &fresh, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedKubernetesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudManagedKubernetesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := managedKubernetesEndpoint(data.ServiceName.ValueString(), data.Id.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}

	if err := waitForCloudResourceStatus(ctx, r.config.OVHClient, endpoint, "managed Kubernetes cluster", data.Id.ValueString(), append([]string{"DELETING", "READY"}, cloudResourceBusy...), []string{"DELETED"}); err != nil {
		resp.Diagnostics.AddError("Error waiting for managed Kubernetes cluster to be deleted", err.Error())
	}
}
