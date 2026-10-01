package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudFloatingIPsDataSource)(nil)

func NewCloudFloatingIPsDataSource() datasource.DataSource {
	return &cloudFloatingIPsDataSource{}
}

type cloudFloatingIPsDataSource struct {
	config *Config
}

// cloudFloatingIPsDataSourceModel is the model for the plural floating IPs data source.
type cloudFloatingIPsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	FloatingIPs types.List             `tfsdk:"floating_ips"`
}

func (d *cloudFloatingIPsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_floating_ips"
}

func (d *cloudFloatingIPsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

func (d *cloudFloatingIPsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the floating IPs of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Computed:    true,
				Description: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
			},
			"floating_ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of floating IPs",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "IP address of the floating IP",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Description of the floating IP",
						},
						"location": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Location of the floating IP",
							Attributes: map[string]schema.Attribute{
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Region",
								},
								"availability_zone": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Availability zone",
								},
							},
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current target specification value",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the floating IP",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the floating IP",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Floating IP readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Current state of the floating IP",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "OpenStack identifier of the floating IP",
								},
								"ip": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "IP address of the floating IP",
								},
								"status": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "OpenStack status of the floating IP (ACTIVE, DOWN, ERROR)",
								},
								"description": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Description of the floating IP",
								},
								"network": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "External network the floating IP belongs to",
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											CustomType:  ovhtypes.TfStringType{},
											Computed:    true,
											Description: "Network ID",
										},
									},
								},
								"associated_resource": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Resource the floating IP is currently attached to",
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											CustomType:  ovhtypes.TfStringType{},
											Computed:    true,
											Description: "ID of the associated resource",
										},
										"type": schema.StringAttribute{
											CustomType:  ovhtypes.TfStringType{},
											Computed:    true,
											Description: "Type of the associated resource",
										},
									},
								},
								"location": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Location details",
									Attributes: map[string]schema.Attribute{
										"region": schema.StringAttribute{
											CustomType:  ovhtypes.TfStringType{},
											Computed:    true,
											Description: "Region",
										},
										"availability_zone": schema.StringAttribute{
											CustomType:  ovhtypes.TfStringType{},
											Computed:    true,
											Description: "Availability zone",
										},
									},
								},
							},
						},
						"current_tasks": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Ongoing asynchronous tasks related to the floating IP",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"errors": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Errors that occured on the task",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"message": schema.StringAttribute{
													CustomType:  ovhtypes.TfStringType{},
													Computed:    true,
													Description: "Error description",
												},
											},
										},
									},
									"id": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Identifier of the current task",
									},
									"link": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Link to the task details",
									},
									"status": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Current global status of the current task",
									},
									"type": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Type of the current task",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudFloatingIPsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudFloatingIPsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() || data.ServiceName.ValueString() == "" {
		envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
		if envServiceName == "" {
			resp.Diagnostics.AddError(
				"Missing service_name",
				"The service_name attribute is required. Provide it in the data source "+
					"configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
			)
			return
		}
		data.ServiceName = ovhtypes.NewTfStringValue(envServiceName)
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicIp/floating"

	responseData, err := helpers.GetAllPagesV2[CloudFloatingIPAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: CloudFloatingIPListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildCloudFloatingIPListItemObject(ctx, &responseData[i]))
	}

	data.FloatingIPs = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
