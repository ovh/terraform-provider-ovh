package ovh

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedDatabaseServiceDataSource)(nil)

func NewCloudManagedDatabaseDataSource() datasource.DataSource {
	return &cloudManagedDatabaseServiceDataSource{category: managedDatabaseFamily}
}

func NewCloudManagedAnalyticsDataSource() datasource.DataSource {
	return &cloudManagedDatabaseServiceDataSource{category: managedAnalyticsFamily}
}

type cloudManagedDatabaseServiceDataSource struct {
	managedDatabaseDataSourceConfigured
	category managedDatabaseCategory
}

func (d *cloudManagedDatabaseServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.category.typeName("")
}

func (d *cloudManagedDatabaseServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	computedInt64 := func(description string) schema.Int64Attribute {
		return schema.Int64Attribute{Computed: true, Description: description}
	}
	computedBool := func(description string) schema.BoolAttribute {
		return schema.BoolAttribute{Computed: true, Description: description}
	}
	computedStringList := func(description string) schema.ListAttribute {
		return schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: description}
	}
	computedStringMap := func(description string) schema.MapAttribute {
		return schema.MapAttribute{ElementType: types.StringType, Computed: true, Description: description}
	}
	maintenanceWindow := schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Maintenance window.",
		Attributes: map[string]schema.Attribute{
			"start_time_utc": managedDatabaseComputedString("Start time of the maintenance window, as HH:MM UTC."),
		},
	}
	// The rules of the current state also carry their creation date.
	networkWithRules := func(rule map[string]schema.Attribute) schema.SingleNestedAttribute {
		return schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Network configuration.",
			Attributes: map[string]schema.Attribute{
				"network_id": managedDatabaseComputedString("ID of the private network the service is attached to."),
				"subnet_id":  managedDatabaseComputedString("ID of the subnet the service is attached to."),
				"authorized_ips": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "IP addresses allowed to connect to the service.",
					Attributes: map[string]schema.Attribute{
						"rules": schema.SetNestedAttribute{
							Computed:     true,
							Description:  "Allowed CIDR blocks.",
							NestedObject: schema.NestedAttributeObject{Attributes: rule},
						},
					},
				},
				"public": computedBool("Whether the service is published on a public address."),
			},
		}
	}
	network := networkWithRules(map[string]schema.Attribute{
		"cidr":        managedDatabaseComputedString("CIDR block."),
		"description": managedDatabaseComputedString("Description of the rule."),
	})
	networkState := networkWithRules(map[string]schema.Attribute{
		"cidr":        managedDatabaseComputedString("CIDR block."),
		"created_at":  managedDatabaseComputedString("Creation date of the rule."),
		"description": managedDatabaseComputedString("Description of the rule."),
	})
	storageSize := schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Storage size.",
		Attributes: map[string]schema.Attribute{
			"unit":  managedDatabaseComputedString("Unit of the size: GB or MB."),
			"value": computedInt64("Size, in the unit."),
		},
	}

	resp.Schema = schema.Schema{
		Description: "Use this data source to get a managed " + d.category.family + " service (" + strings.Join(d.category.engines, ", ") + ") of a public cloud project, using the Cloud Databases API v2.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required:    true,
				Description: "ID of the public cloud project.",
			},
			"engine": schema.StringAttribute{
				Required:    true,
				Description: "Engine of the service: " + strings.Join(d.category.engines, ", ") + ".",
				Validators:  []validator.String{stringvalidator.OneOf(d.category.engines...)},
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the service.",
			},
			"region":                 managedDatabaseComputedString("Region of the nodes."),
			"version":                managedDatabaseComputedString("Engine version."),
			"plan":                   managedDatabaseComputedString("Plan of the service."),
			"flavor":                 managedDatabaseComputedString("Flavor of the nodes."),
			"nodes_count":            computedInt64("Number of nodes."),
			"description":            managedDatabaseComputedString("Description of the service."),
			"deletion_protected":     computedBool("Whether the service is protected from deletion."),
			"prometheus_enabled":     computedBool("Whether the Prometheus endpoint is enabled."),
			"advanced_configuration": computedStringMap("Advanced configuration properties."),
			"backup": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Backup configuration.",
				Attributes: map[string]schema.Attribute{
					"regions": computedStringList("Regions where the backups are stored."),
					"time":    managedDatabaseComputedString("Time of the daily backup, as HH:MM UTC."),
				},
			},
			"maintenance_window": maintenanceWindow,
			"network":            network,
			"storage": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Storage configuration.",
				Attributes: map[string]schema.Attribute{
					"size": storageSize,
				},
			},
			"checksum":        managedDatabaseComputedString("Hash of the service, computed by the API."),
			"created_at":      managedDatabaseComputedString("Creation date of the service."),
			"updated_at":      managedDatabaseComputedString("Last update date of the service."),
			"resource_status": managedDatabaseComputedString("Readiness of the service: CREATING, UPDATING, READY, ERROR, OUT_OF_SYNC, DELETING…"),
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
									"message": managedDatabaseComputedString("Description of the error."),
								},
							},
						},
						"id":     managedDatabaseComputedString("ID of the task."),
						"link":   managedDatabaseComputedString("Link to the details of the task."),
						"status": managedDatabaseComputedString("Status of the task: ERROR, PENDING, RUNNING, SCHEDULED or WAITING_USER_INPUT."),
						"type":   managedDatabaseComputedString("Type of the task."),
					},
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Configuration running on the service.",
				Attributes: map[string]schema.Attribute{
					"advanced_configuration": computedStringMap("Advanced configuration properties."),
					"backup": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Backup configuration.",
						Attributes: map[string]schema.Attribute{
							"pitr":           managedDatabaseComputedString("Date until which point-in-time recovery is available."),
							"regions":        computedStringList("Regions where the backups are stored."),
							"retention_days": computedInt64("Number of days the backups are kept."),
							"time":           managedDatabaseComputedString("Time of the daily backup, as HH:MM UTC."),
						},
					},
					"capabilities": schema.MapNestedAttribute{
						Computed:    true,
						Description: "Capabilities of the service, such as USERS or CONNECTION_POOLS, with the state of each action: ENABLED or DISABLED.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"create": managedDatabaseComputedString("State of the create action."),
								"delete": managedDatabaseComputedString("State of the delete action."),
								"read":   managedDatabaseComputedString("State of the read action."),
								"update": managedDatabaseComputedString("State of the update action."),
							},
						},
					},
					"deletion_protected": computedBool("Whether the service is protected from deletion."),
					"description":        managedDatabaseComputedString("Description of the service."),
					"engine_capabilities": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Capabilities of the engine of the service.",
						Attributes: map[string]schema.Attribute{
							"advanced_configuration_properties": schema.ListNestedAttribute{
								Computed:    true,
								Description: "Advanced configuration properties that can be set, with their constraints.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"default":          managedDatabaseComputedString("Default value of the property."),
										"description":      managedDatabaseComputedString("Description of the property."),
										"maximum":          schema.Float64Attribute{Computed: true, Description: "Maximum value of a numeric property."},
										"minimum":          schema.Float64Attribute{Computed: true, Description: "Minimum value of a numeric property."},
										"name":             managedDatabaseComputedString("Name of the property."),
										"requires_restart": computedBool("Whether changing the property restarts the service."),
										"type":             managedDatabaseComputedString("Type of the property: BOOLEAN, DOUBLE, LONG or STRING."),
										"values":           computedStringList("Possible values of a string property."),
									},
								},
							},
							"backup_regions": computedStringList("Regions where the backups can be stored."),
							"categories":     computedStringList("Command categories a user can be given. Valkey only."),
							"commands":       computedStringList("Commands a user can be given. Valkey only."),
							"integrations": schema.ListNestedAttribute{
								Computed:    true,
								Description: "Integrations the service can take part in.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"destination_engine": managedDatabaseComputedString("Engine of the destination service."),
										"parameters": schema.ListNestedAttribute{
											Computed:    true,
											Description: "Parameters of the integration.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"name": managedDatabaseComputedString("Name of the parameter."),
													"type": managedDatabaseComputedString("Type of the parameter: INTEGER or STRING."),
												},
											},
										},
										"source_engine": managedDatabaseComputedString("Engine of the source service."),
										"type":          managedDatabaseComputedString("Type of the integration."),
									},
								},
							},
						},
					},
					"endpoints": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Connection endpoints of the service.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"component": managedDatabaseComputedString("Component reached through the endpoint."),
								"domain":    managedDatabaseComputedString("Domain name of the endpoint."),
								"path":      managedDatabaseComputedString("Path of the endpoint."),
								"port":      computedInt64("Port of the endpoint."),
								"public":    computedBool("Whether the endpoint uses the public domain name of the service."),
								"scheme":    managedDatabaseComputedString("Scheme of the endpoint."),
								"ssl":       computedBool("Whether TLS is enabled on the endpoint."),
								"ssl_mode":  managedDatabaseComputedString("TLS mode of the endpoint."),
								"uri":       managedDatabaseComputedString("Connection URI of the endpoint."),
							},
						},
					},
					"flavor":             managedDatabaseComputedString("Flavor of the nodes."),
					"maintenance_window": maintenanceWindow,
					"network":            networkState,
					"nodes": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Nodes of the service.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"created_at": managedDatabaseComputedString("Creation date of the node."),
								"flavor":     managedDatabaseComputedString("Flavor of the node."),
								"id":         managedDatabaseComputedString("ID of the node."),
								"location": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Location of the node.",
									Attributes: map[string]schema.Attribute{
										"availability_zone": managedDatabaseComputedString("Availability zone of the node."),
										"region":            managedDatabaseComputedString("Region of the node."),
									},
								},
								"name":   managedDatabaseComputedString("Name of the node."),
								"port":   computedInt64("Port of the node."),
								"role":   managedDatabaseComputedString("Role of the node: PRIMARY, SECONDARY, STANDBY…"),
								"status": managedDatabaseComputedString("Readiness of the node."),
							},
						},
					},
					"nodes_count":        computedInt64("Number of nodes."),
					"plan":               managedDatabaseComputedString("Plan of the service."),
					"prometheus_enabled": computedBool("Whether the Prometheus endpoint is enabled."),
					"status":             managedDatabaseComputedString("Status of the service reported by the engine."),
					"storage": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Storage configuration.",
						Attributes: map[string]schema.Attribute{
							"size": storageSize,
							"type": managedDatabaseComputedString("Type of storage."),
						},
					},
					"version": managedDatabaseComputedString("Engine version."),
				},
			},
		},
	}
}

func (d *cloudManagedDatabaseServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudManagedDatabaseServiceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := d.category.serviceEndpoint(data.ServiceName.ValueString(), data.Engine.ValueString(), data.Id.ValueString())
	var res CloudManagedDatabaseServiceAPIResponse
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &res); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	// Unknown makes MergeWith return every key of the advanced configuration.
	data.AdvancedConfiguration = types.MapUnknown(types.StringType)
	resp.Diagnostics.Append(data.MergeWith(ctx, &res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
