package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudExtNetIPsDataSource)(nil)

func NewCloudExtNetIPsDataSource() datasource.DataSource {
	return &cloudExtNetIPsDataSource{}
}

type cloudExtNetIPsDataSource struct {
	config *Config
}

// cloudExtNetIPsDataSourceModel is the model for the plural external network IPs data source.
type cloudExtNetIPsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ExtNetIPs   types.List             `tfsdk:"ext_net_ips"`
}

func (d *cloudExtNetIPsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_ext_net_ips"
}

func (d *cloudExtNetIPsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudExtNetIPsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the external network IPs of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"ext_net_ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of external network IPs",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "IP address of the external network IP",
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Checksum field of the API envelope. Always empty for this read-only IP type.",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the external network IP",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the external network IP",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "External network IP readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Current state of the external network IP",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Identifier of the external network IP",
								},
								"ip": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "IP address of the external network IP",
								},
								"associated_resource": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Resource the external network IP is currently attached to",
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
							Description: "Ongoing asynchronous tasks related to the external network IP",
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

func (d *cloudExtNetIPsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudExtNetIPsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicIp/extNet"

	responseData, err := helpers.GetAllPagesV2[CloudExtNetIPAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: CloudExtNetIPListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildCloudExtNetIPListItemObject(ctx, &responseData[i]))
	}

	data.ExtNetIPs = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
