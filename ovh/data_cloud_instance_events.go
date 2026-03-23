package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceEventsDataSource)(nil)

func NewCloudInstanceEventsDataSource() datasource.DataSource {
	return &cloudInstanceEventsDataSource{}
}

type cloudInstanceEventsDataSource struct {
	config *Config
}

func (d *cloudInstanceEventsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_events"
}

func (d *cloudInstanceEventsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceEventsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List events (action history) for a compute instance in a public cloud project.",
		MarkdownDescription: "List events (action history) for a compute instance in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Required
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"instance_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the compute instance",
				MarkdownDescription: "ID of the compute instance",
			},

			// Computed
			"events": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of events for the instance, sorted by most recent first",
				MarkdownDescription: "List of events for the instance, sorted by most recent first",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Computed:            true,
							Description:         "Creation date of the event",
							MarkdownDescription: "Creation date of the event",
						},
						"kind": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Computed:            true,
							Description:         "Nature of the event (e.g. LOCK, UNLOCK, REBOOT, CREATE, STOP, START)",
							MarkdownDescription: "Nature of the event (e.g. `LOCK`, `UNLOCK`, `REBOOT`, `CREATE`, `STOP`, `START`)",
						},
						"link": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Computed:            true,
							Description:         "Link to the event related resource",
							MarkdownDescription: "Link to the event related resource",
						},
						"message": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Computed:            true,
							Description:         "Description of what happened on the event",
							MarkdownDescription: "Description of what happened on the event",
						},
						"type": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Computed:            true,
							Description:         "Type of the event (TASK_START, TASK_SUCCESS, TASK_ERROR, TARGET_SPEC_UPDATE)",
							MarkdownDescription: "Type of the event (`TASK_START`, `TASK_SUCCESS`, `TASK_ERROR`, `TARGET_SPEC_UPDATE`)",
						},
					},
				},
			},
		},
	}
}

// cloudInstanceEventAPIResponse maps a single event from the API response.
type cloudInstanceEventAPIResponse struct {
	CreatedAt string  `json:"createdAt"`
	Kind      string  `json:"kind"`
	Link      *string `json:"link"`
	Message   string  `json:"message"`
	Type      string  `json:"type"`
}

// cloudInstanceEventsDataSourceModel is the Terraform state model.
type cloudInstanceEventsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	InstanceId  ovhtypes.TfStringValue `tfsdk:"instance_id"`
	Events      types.List             `tfsdk:"events"`
}

// eventAttrTypes returns the attribute types for a single event object.
func eventAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at": ovhtypes.TfStringType{},
		"kind":       ovhtypes.TfStringType{},
		"link":       ovhtypes.TfStringType{},
		"message":    ovhtypes.TfStringType{},
		"type":       ovhtypes.TfStringType{},
	}
}

func (d *cloudInstanceEventsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudInstanceEventsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/compute/instance/" + url.PathEscape(data.InstanceId.ValueString()) + "/events"

	var apiEvents []cloudInstanceEventAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiEvents); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Convert API events to Terraform list of objects
	eventObjs := make([]attr.Value, len(apiEvents))
	for i, evt := range apiEvents {
		linkVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
		if evt.Link != nil {
			linkVal = ovhtypes.TfStringValue{StringValue: types.StringValue(*evt.Link)}
		}

		obj, diags := types.ObjectValue(
			eventAttrTypes(),
			map[string]attr.Value{
				"created_at": ovhtypes.TfStringValue{StringValue: types.StringValue(evt.CreatedAt)},
				"kind":       ovhtypes.TfStringValue{StringValue: types.StringValue(evt.Kind)},
				"link":       linkVal,
				"message":    ovhtypes.TfStringValue{StringValue: types.StringValue(evt.Message)},
				"type":       ovhtypes.TfStringValue{StringValue: types.StringValue(evt.Type)},
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		eventObjs[i] = obj
	}

	eventsList, diags := types.ListValue(types.ObjectType{AttrTypes: eventAttrTypes()}, eventObjs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Events = eventsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
