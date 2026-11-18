package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

var (
	_ resource.Resource                = (*cloudManagedDatabaseServiceResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudManagedDatabaseServiceResource)(nil)
	_ resource.ResourceWithImportState = (*cloudManagedDatabaseServiceResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*cloudManagedDatabaseServiceResource)(nil)
)

const (
	managedDatabaseServiceCreateTimeout = 40 * time.Minute
	managedDatabaseServiceUpdateTimeout = 40 * time.Minute
	managedDatabaseServiceDeleteTimeout = 20 * time.Minute
)

// Variables so that the unit tests don't wait.
var (
	managedDatabaseServicePollDelay      = 10 * time.Second
	managedDatabaseServicePollMinTimeout = 5 * time.Second
)

// managedDatabaseCreationOnly recreates the service when an argument that only the creation takes
// changes from one value to another. Setting it, as after an import, or removing it only changes the
// state.
var managedDatabaseCreationOnly = stringplanmodifier.RequiresReplaceIf(
	func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
		resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.IsNull()
	},
	"Changing it from one value to another recreates the service.",
	"Changing it from one value to another recreates the service.",
)

func NewCloudManagedDatabaseResource() resource.Resource {
	return &cloudManagedDatabaseServiceResource{category: managedDatabaseFamily}
}

func NewCloudManagedAnalyticsResource() resource.Resource {
	return &cloudManagedDatabaseServiceResource{category: managedAnalyticsFamily}
}

// cloudManagedDatabaseServiceResource is a service of one of the families: ovh_cloud_managed_database or
// ovh_cloud_managed_analytics.
type cloudManagedDatabaseServiceResource struct {
	managedDatabaseConfigured
	category managedDatabaseCategory
}

func (r *cloudManagedDatabaseServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.category.typeName("")
}

func (r *cloudManagedDatabaseServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	computedString := func(description string) schema.StringAttribute {
		return schema.StringAttribute{Computed: true, Description: description}
	}

	resp.Schema = schema.Schema{
		Description: "Creates a managed " + r.category.family + " service (" + strings.Join(r.category.engines, ", ") + ") in a public cloud project, using the Cloud Databases API v2.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the public cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used. Changing it recreates the service.",
				PlanModifiers: []planmodifier.String{
					EnvDefaultString("OVH_CLOUD_PROJECT_SERVICE", true),
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"engine": managedDatabaseEngineArgument(r.category.engines, "service"),
			"region": schema.StringAttribute{
				Required:      true,
				Description:   "Region of the nodes, such as eu-west-gra. Changing it recreates the service.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "Engine version. Available versions are returned by the capabilities endpoint.",
			},
			"plan": schema.StringAttribute{
				Required:    true,
				Description: "Plan of the service. Available plans are returned by the capabilities endpoint.",
			},
			"flavor": schema.StringAttribute{
				Required:    true,
				Description: "Flavor of the nodes. Available flavors are returned by the capabilities endpoint.",
			},
			"nodes_count": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Description:   "Number of nodes. Defaults to the minimum of the plan.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the service.",
			},
			"deletion_protected": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether the service is protected from deletion. It must be false to destroy the service.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"prometheus_enabled": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether the Prometheus endpoint is enabled.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"advanced_configuration": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Advanced configuration properties requested for the service, with the defaults the API adds.",
			},
			"backup": schema.SingleNestedAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Backup configuration.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"regions": schema.ListAttribute{
						ElementType:   types.StringType,
						Optional:      true,
						Computed:      true,
						Description:   "Regions where the backups are stored.",
						PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
					},
					"time": schema.StringAttribute{
						Optional:      true,
						Computed:      true,
						Description:   "Time of the daily backup, as HH:MM UTC.",
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
				},
			},
			"maintenance_window": schema.SingleNestedAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Maintenance window.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"start_time_utc": schema.StringAttribute{
						Required:    true,
						Description: "Start time of the maintenance window, as HH:MM UTC.",
					},
				},
			},
			"network": schema.SingleNestedAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Network configuration.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"network_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "ID of the private network the service is attached to. Changing it recreates the service.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"subnet_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "ID of the subnet the service is attached to. Changing it recreates the service.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"authorized_ips": schema.SingleNestedAttribute{
						Optional:      true,
						Computed:      true,
						Description:   "IP addresses allowed to connect to the service.",
						PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Attributes: map[string]schema.Attribute{
							"rules": schema.SetNestedAttribute{
								Required:    true,
								Description: "CIDR blocks allowed to connect, from the private network too. An empty set denies every client.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"cidr": schema.StringAttribute{
											Required:    true,
											Description: "CIDR block, such as 203.0.113.0/24.",
										},
										"description": schema.StringAttribute{
											Optional:    true,
											Description: "Description of the rule.",
										},
									},
								},
							},
						},
					},
					"public": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						Description:   "Whether the service is published on a public address, with a public domain name and TLS. The authorized IPs apply to it too. Changing it sends a PATCH request, which a consumer key must allow.",
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
				},
			},
			"storage": schema.SingleNestedAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Storage configuration.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"size": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Storage size.",
						Attributes: map[string]schema.Attribute{
							"unit": schema.StringAttribute{
								Required:    true,
								Description: "Unit of the size: GB or MB.",
								Validators:  []validator.String{stringvalidator.OneOf("GB", "MB")},
							},
							"value": schema.Int64Attribute{
								Required:    true,
								Description: "Size, in the unit.",
							},
						},
					},
				},
			},
			"fork_source_service_id": schema.StringAttribute{
				Optional:      true,
				Description:   "ID of the service to create this one from, using one of its backups. Changing it to another value recreates the service.",
				PlanModifiers: []planmodifier.String{managedDatabaseCreationOnly},
			},
			"fork_source_backup_id": schema.StringAttribute{
				Optional:      true,
				Description:   "ID of the backup to create the service from. When it and fork_source_point_in_time are null, the latest backup of fork_source_service_id is used. Changing it to another value recreates the service.",
				PlanModifiers: []planmodifier.String{managedDatabaseCreationOnly},
				Validators:    []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("fork_source_point_in_time"))},
			},
			"fork_source_point_in_time": schema.StringAttribute{
				Optional:      true,
				Description:   "Point in time to create the service from, in RFC 3339 format. The source service must have point-in-time recovery enabled. Changing it to another value recreates the service.",
				PlanModifiers: []planmodifier.String{managedDatabaseCreationOnly},
			},
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "ID of the service.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"checksum": computedString("Hash of the service, computed by the API."),
			"created_at": schema.StringAttribute{
				Computed:      true,
				Description:   "Creation date of the service.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": computedString("Last update date of the service."),
			"resource_status": schema.StringAttribute{
				Computed:      true,
				Description:   "Readiness of the service: CREATING, UPDATING, READY, ERROR, OUT_OF_SYNC, DELETING… The next apply updates an OUT_OF_SYNC service, sending its configuration again.",
				PlanModifiers: []planmodifier.String{OutOfSyncPlanModifier()},
			},
			"current_tasks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Asynchronous operations running on the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"errors": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Errors of the task.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"message": computedString("Description of the error."),
								},
							},
						},
						"id":     computedString("ID of the task."),
						"link":   computedString("Link to the details of the task."),
						"status": computedString("Status of the task: ERROR, PENDING, RUNNING, SCHEDULED or WAITING_USER_INPUT."),
						"type":   computedString("Type of the task."),
					},
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Configuration running on the service.",
				Attributes: map[string]schema.Attribute{
					"advanced_configuration": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "Advanced configuration properties.",
					},
					"backup": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Backup configuration.",
						Attributes: map[string]schema.Attribute{
							"pitr": computedString("Date until which point-in-time recovery is available."),
							"regions": schema.ListAttribute{
								ElementType: types.StringType,
								Computed:    true,
								Description: "Regions where the backups are stored.",
							},
							"retention_days": schema.Int64Attribute{
								Computed:    true,
								Description: "Number of days the backups are kept.",
							},
							"time": computedString("Time of the daily backup, as HH:MM UTC."),
						},
					},
					"capabilities": schema.MapNestedAttribute{
						Computed:    true,
						Description: "Capabilities of the service, such as USERS or CONNECTION_POOLS, with the state of each action: ENABLED or DISABLED.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"create": computedString("State of the create action."),
								"delete": computedString("State of the delete action."),
								"read":   computedString("State of the read action."),
								"update": computedString("State of the update action."),
							},
						},
					},
					"deletion_protected": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the service is protected from deletion.",
					},
					"engine_capabilities": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Capabilities of the engine of the service.",
						Attributes: map[string]schema.Attribute{
							"advanced_configuration_properties": schema.ListNestedAttribute{
								Computed:    true,
								Description: "Advanced configuration properties that can be set, with their constraints.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"default":     computedString("Default value of the property."),
										"description": computedString("Description of the property."),
										"maximum": schema.Float64Attribute{
											Computed:    true,
											Description: "Maximum value of a numeric property.",
										},
										"minimum": schema.Float64Attribute{
											Computed:    true,
											Description: "Minimum value of a numeric property.",
										},
										"name": computedString("Name of the property."),
										"requires_restart": schema.BoolAttribute{
											Computed:    true,
											Description: "Whether changing the property restarts the service.",
										},
										"type": computedString("Type of the property: BOOLEAN, DOUBLE, LONG or STRING."),
										"values": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
											Description: "Possible values of a string property.",
										},
									},
								},
							},
							"backup_regions": schema.ListAttribute{
								ElementType: types.StringType,
								Computed:    true,
								Description: "Regions where the backups can be stored.",
							},
							"categories": schema.ListAttribute{
								ElementType: types.StringType,
								Computed:    true,
								Description: "Command categories a user can be given. Valkey only.",
							},
							"commands": schema.ListAttribute{
								ElementType: types.StringType,
								Computed:    true,
								Description: "Commands a user can be given. Valkey only.",
							},
							"integrations": schema.ListNestedAttribute{
								Computed:    true,
								Description: "Integrations the service can take part in.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"destination_engine": computedString("Engine of the destination service."),
										"parameters": schema.ListNestedAttribute{
											Computed:    true,
											Description: "Parameters of the integration.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"name": computedString("Name of the parameter."),
													"type": computedString("Type of the parameter: INTEGER or STRING."),
												},
											},
										},
										"source_engine": computedString("Engine of the source service."),
										"type":          computedString("Type of the integration."),
									},
								},
							},
						},
					},
					"description": computedString("Description of the service."),
					"endpoints": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Connection endpoints of the service.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"component": computedString("Component reached through the endpoint."),
								"domain":    computedString("Domain name of the endpoint."),
								"path":      computedString("Path of the endpoint."),
								"port": schema.Int64Attribute{
									Computed:    true,
									Description: "Port of the endpoint.",
								},
								"public": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether the endpoint uses the public domain name of the service.",
								},
								"scheme": computedString("Scheme of the endpoint."),
								"ssl": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether TLS is enabled on the endpoint.",
								},
								"ssl_mode": computedString("TLS mode of the endpoint."),
								"uri":      computedString("Connection URI of the endpoint."),
							},
						},
					},
					"flavor": computedString("Flavor of the nodes."),
					"maintenance_window": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Maintenance window.",
						Attributes: map[string]schema.Attribute{
							"start_time_utc": computedString("Start time of the maintenance window, as HH:MM UTC."),
						},
					},
					"network": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Network configuration.",
						Attributes: map[string]schema.Attribute{
							"network_id": computedString("ID of the private network."),
							"subnet_id":  computedString("ID of the subnet."),
							"authorized_ips": schema.SingleNestedAttribute{
								Computed:    true,
								Description: "IP addresses allowed to connect to the service.",
								Attributes: map[string]schema.Attribute{
									"rules": schema.SetNestedAttribute{
										Computed:    true,
										Description: "Allowed CIDR blocks.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"cidr":        computedString("CIDR block."),
												"created_at":  computedString("Creation date of the rule."),
												"description": computedString("Description of the rule."),
											},
										},
									},
								},
							},
							"public": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether the service is published on a public address.",
							},
						},
					},
					"nodes": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Nodes of the service.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"created_at": computedString("Creation date of the node."),
								"flavor":     computedString("Flavor of the node."),
								"id":         computedString("ID of the node."),
								"location": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Location of the node.",
									Attributes: map[string]schema.Attribute{
										"availability_zone": computedString("Availability zone of the node."),
										"region":            computedString("Region of the node."),
									},
								},
								"name": computedString("Name of the node."),
								"port": schema.Int64Attribute{
									Computed:    true,
									Description: "Port of the node.",
								},
								"role":   computedString("Role of the node: PRIMARY, SECONDARY, STANDBY…"),
								"status": computedString("Readiness of the node."),
							},
						},
					},
					"nodes_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of nodes.",
					},
					"plan": computedString("Plan of the service."),
					"prometheus_enabled": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the Prometheus endpoint is enabled.",
					},
					"status": computedString("Status of the service reported by the engine."),
					"storage": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Storage configuration.",
						Attributes: map[string]schema.Attribute{
							"size": schema.SingleNestedAttribute{
								Computed:    true,
								Description: "Storage size.",
								Attributes: map[string]schema.Attribute{
									"unit": computedString("Unit of the size."),
									"value": schema.Int64Attribute{
										Computed:    true,
										Description: "Size, in the unit.",
									},
								},
							},
							"type": computedString("Type of storage."),
						},
					},
					"version": computedString("Engine version."),
				},
			},
		},
	}
}

func (r *cloudManagedDatabaseServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<engine>/<id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("engine"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

// ModifyPlan keeps the state when the plan only differs from it by unknown values and the
// configuration has none: the update would send the same service. Terraform plans one when the
// configuration leaves out maintenance_window, network or storage: the API fills in their required
// attributes, so Terraform proposes null for them, which makes every computed attribute unknown. An
// unknown value of the configuration, such as the ID of a subnet being replaced, can change the
// service.
func (r *cloudManagedDatabaseServiceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// OutOfSyncPlanModifier plans an update of an OUT_OF_SYNC service, which Terraform doesn't see as a
	// change: what the API computes is unknown until the PUT, as in any update.
	var status types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("resource_status"), &status)...)
	if status.ValueString() == "OUT_OF_SYNC" {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("checksum"), types.StringUnknown())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("updated_at"), types.StringUnknown())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("advanced_configuration"), types.MapUnknown(types.StringType))...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("current_tasks"), types.ListUnknown(types.ObjectType{AttrTypes: managedDatabaseTaskAttrTypes()}))...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("current_state"), types.ObjectUnknown(managedDatabaseCurrentStateAttrTypes()))...)
		return
	}

	if !req.Config.Raw.IsFullyKnown() {
		return
	}
	known, err := tftypes.Transform(req.Plan.Raw, func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
		if v.IsKnown() {
			return v, nil
		}
		prior, _, err := tftypes.WalkAttributePath(req.State.Raw, p)
		if err != nil {
			return v, nil
		}
		return prior.(tftypes.Value), nil
	})
	if err == nil && known.Equal(req.State.Raw) {
		resp.Plan.Raw = req.State.Raw
	}
}

func (r *cloudManagedDatabaseServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cloudManagedDatabaseServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := data.ToCreate(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.category.endpoint(data.ServiceName.ValueString(), data.Engine.ValueString())
	var res CloudManagedDatabaseServiceAPIResponse
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, payload, &res); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Saved right away so that the service stays tracked if what follows fails. It is a copy, so
	// that what follows still reads the plan.
	created := data
	resp.Diagnostics.Append(created.MergeWith(ctx, &res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &created)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint += "/" + url.PathEscape(res.Id)
	if err := waitForManagedDatabaseServiceReady(ctx, r.config.OVHClient, endpoint, managedDatabaseServiceCreateTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for the service to be ready", err.Error())
		return
	}

	// The creation ignores the deletion protection: a PUT applies it, when the service doesn't have it yet.
	protect := func(service *CloudManagedDatabaseServiceModel) {
		if !data.DeletionProtected.IsUnknown() && !data.DeletionProtected.IsNull() {
			service.DeletionProtected = data.DeletionProtected
		}
	}
	if err := putManagedDatabaseServiceWhenReady(ctx, r.config.OVHClient, endpoint, protect, managedDatabaseServiceCreateTimeout); err != nil {
		resp.Diagnostics.AddError("Error applying the deletion protection of the service", err.Error())
		return
	}

	res = CloudManagedDatabaseServiceAPIResponse{}
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &res); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	resp.Diagnostics.Append(data.MergeWith(ctx, &res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedDatabaseServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cloudManagedDatabaseServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.category.serviceEndpoint(data.ServiceName.ValueString(), data.Engine.ValueString(), data.Id.ValueString())
	var res CloudManagedDatabaseServiceAPIResponse
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &res); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	resp.Diagnostics.Append(data.MergeWith(ctx, &res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedDatabaseServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state cloudManagedDatabaseServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.category.serviceEndpoint(state.ServiceName.ValueString(), state.Engine.ValueString(), state.Id.ValueString())

	// The PUT keeps the public exposure, which a PATCH changes: after the PUT when it publishes the
	// service and before it when it unpublishes it, so that the public address only ever filters with
	// the authorized IPs configured along with it.
	planned := managedDatabaseObjectAs[managedDatabaseNetwork](ctx, plan.Network, &resp.Diagnostics)
	prior := managedDatabaseObjectAs[managedDatabaseNetwork](ctx, state.Network, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	var public *bool
	if planned != nil && planned.Public != nil && (prior == nil || !reflect.DeepEqual(prior.Public, planned.Public)) {
		public = planned.Public
	}
	if public != nil && !*public {
		if err := publishManagedDatabaseServiceWhenReady(ctx, r.config.OVHClient, endpoint, false, managedDatabaseServiceUpdateTimeout); err != nil {
			resp.Diagnostics.AddError("Error unpublishing the service", err.Error())
			return
		}
	}

	// advanced_configuration is read-only: the PUT sends the current one back.
	replace := func(service *CloudManagedDatabaseServiceModel) {
		current := service.AdvancedConfiguration
		*service = plan.CloudManagedDatabaseServiceModel
		service.AdvancedConfiguration = current
	}
	if err := putManagedDatabaseServiceWhenReady(ctx, r.config.OVHClient, endpoint, replace, managedDatabaseServiceUpdateTimeout); err != nil {
		resp.Diagnostics.AddError("Error updating the service", err.Error())
		return
	}

	if public != nil && *public {
		if err := publishManagedDatabaseServiceWhenReady(ctx, r.config.OVHClient, endpoint, true, managedDatabaseServiceUpdateTimeout); err != nil {
			resp.Diagnostics.AddError("Error publishing the service", err.Error())
			return
		}
	}

	var res CloudManagedDatabaseServiceAPIResponse
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &res); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	resp.Diagnostics.Append(plan.MergeWith(ctx, &res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cloudManagedDatabaseServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data cloudManagedDatabaseServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.category.serviceEndpoint(data.ServiceName.ValueString(), data.Engine.ValueString(), data.Id.ValueString())
	err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseServiceDeleteTimeout, func() error {
		return r.config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	})
	if helpers.IsNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}

	// The service can still report its previous status right after the DELETE.
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING", "READY", "UPDATING"},
		Target:  []string{"DELETED"},
		Refresh: func() (any, string, error) {
			res := &CloudManagedDatabaseServiceAPIResponse{}
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if helpers.IsNotFound(err) {
				return res, "DELETED", nil
			}
			if err != nil {
				return nil, "", err
			}
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("service", res.Id, res.CurrentTasks)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    managedDatabaseServiceDeleteTimeout,
		Delay:      managedDatabaseServicePollDelay,
		MinTimeout: managedDatabaseServicePollMinTimeout,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError("Error waiting for the service to be deleted", err.Error())
	}
}

// putManagedDatabaseServiceWhenReady waits for the service to be READY or OUT_OF_SYNC, the statuses in
// which it takes a change, then replaces its whole targetSpec with the current one after change, and
// waits for the service to apply it. Nothing is sent when change changes nothing, unless the service
// is OUT_OF_SYNC: the same targetSpec brings it back in sync.
func putManagedDatabaseServiceWhenReady(ctx context.Context, client *ovhwrap.Client, endpoint string, change func(*CloudManagedDatabaseServiceModel), timeout time.Duration) error {
	if err := waitForManagedDatabaseServiceReady(ctx, client, endpoint, timeout, "OUT_OF_SYNC"); err != nil {
		return err
	}
	if err := retryWhileManagedDatabaseServiceBusy(ctx, timeout, func() error {
		var res CloudManagedDatabaseServiceAPIResponse
		if err := client.GetWithContext(ctx, endpoint, &res); err != nil {
			return err
		}
		var service CloudManagedDatabaseServiceModel
		diags := service.MergeWith(ctx, &res)
		current, d := service.ToPut(ctx)
		diags.Append(d...)
		change(&service)
		body, d := service.ToPut(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return fmt.Errorf("%s: %s", diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
		}
		if reflect.DeepEqual(body, current) && res.ResourceStatus != "OUT_OF_SYNC" {
			return nil
		}
		return client.PutWithContext(ctx, endpoint, body, nil)
	}); err != nil {
		return fmt.Errorf("calling Put %s: %w", endpoint, err)
	}
	return waitForManagedDatabaseServiceReady(ctx, client, endpoint, timeout)
}

// publishManagedDatabaseServiceWhenReady waits for the service to be READY or OUT_OF_SYNC, publishes it
// on a public address or takes it off with a PATCH, and waits for the service to apply it.
func publishManagedDatabaseServiceWhenReady(ctx context.Context, client *ovhwrap.Client, endpoint string, public bool, timeout time.Duration) error {
	if err := waitForManagedDatabaseServiceReady(ctx, client, endpoint, timeout, "OUT_OF_SYNC"); err != nil {
		return err
	}
	body := map[string]any{"targetSpec": map[string]any{"network": map[string]any{"public": public}}}
	if err := retryWhileManagedDatabaseServiceBusy(ctx, timeout, func() error {
		return client.PatchWithContext(ctx, endpoint, body, nil)
	}); err != nil {
		return fmt.Errorf("calling Patch %s: %w", endpoint, err)
	}
	return waitForManagedDatabaseServiceReady(ctx, client, endpoint, timeout)
}

// waitForManagedDatabaseServiceReady polls the service until it is READY and so is its engine, the
// status v1 waits for: resourceStatus can be READY before the engine has provisioned every node. A
// status of alsoReady ends the wait too. ERROR is terminal and reports the errors of the failed tasks.
func waitForManagedDatabaseServiceReady(ctx context.Context, client *ovhwrap.Client, endpoint string, timeout time.Duration, alsoReady ...string) error {
	pending := []string{"CREATING", "CREATION_PLANNED", "PENDING", "UPDATING", "OUT_OF_SYNC"}
	stateConf := &retry.StateChangeConf{
		Pending: slices.DeleteFunc(pending, func(status string) bool { return slices.Contains(alsoReady, status) }),
		Target:  append([]string{"READY"}, alsoReady...),
		Refresh: func() (any, string, error) {
			res := &CloudManagedDatabaseServiceAPIResponse{}
			if err := client.GetWithContext(ctx, endpoint, res); err != nil {
				return nil, "", err
			}
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("service", res.Id, res.CurrentTasks)
			}
			if res.ResourceStatus == "READY" && res.CurrentState != nil && res.CurrentState.Status != nil {
				return res, *res.CurrentState.Status, nil
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    timeout,
		Delay:      managedDatabaseServicePollDelay,
		MinTimeout: managedDatabaseServicePollMinTimeout,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// retryWhileManagedDatabaseServiceBusy retries call while the API answers 409, which it does while
// another operation runs on the service.
func retryWhileManagedDatabaseServiceBusy(ctx context.Context, timeout time.Duration, call func() error) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err := call()
		var apiErr *ovh.APIError
		switch {
		case err == nil:
			return nil
		case errors.As(err, &apiErr) && apiErr.Code == http.StatusConflict:
			return retry.RetryableError(err)
		default:
			return retry.NonRetryableError(err)
		}
	})
}
