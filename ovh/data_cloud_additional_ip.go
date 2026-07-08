package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudAdditionalIPDataSource)(nil)

func NewCloudAdditionalIPDataSource() datasource.DataSource {
	return &cloudAdditionalIPDataSource{}
}

type cloudAdditionalIPDataSource struct {
	config *Config
}

func (d *cloudAdditionalIPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_additional_ip"
}

func (d *cloudAdditionalIPDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudAdditionalIPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about an additional IP in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "IP address of the additional IP",
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Checksum field of the API envelope. Always empty for this read-only IP type.",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Additional IP readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the additional IP",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Identifier of the additional IP",
					},
					"ip": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "IP address of the additional IP",
					},
					"ip_block": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "IP block the additional IP belongs to",
					},
					"associated_resource": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Resource the additional IP is currently attached to",
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
				Description: "Ongoing asynchronous tasks related to the additional IP",
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
	}
}

// cloudAdditionalIPDataSourceModel is the Terraform state model for this data source.
type cloudAdditionalIPDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
	CurrentTasks   types.List             `tfsdk:"current_tasks"`
}

func (d *cloudAdditionalIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudAdditionalIPDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicIp/additional/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudAdditionalIPAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.Checksum)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.ResourceStatus)}

	if responseData.CurrentState != nil {
		data.CurrentState = buildCloudAdditionalIPCurrentStateObject(ctx, responseData.CurrentState)
	} else {
		data.CurrentState = types.ObjectNull(CloudAdditionalIPCurrentStateAttrTypes())
	}

	data.CurrentTasks = buildCloudPublicIPCurrentTasksList(responseData.CurrentTasks)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
