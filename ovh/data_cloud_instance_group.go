package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceGroupDataSource)(nil)

func NewCloudInstanceGroupDataSource() datasource.DataSource {
	return &cloudInstanceGroupDataSource{}
}

type cloudInstanceGroupDataSource struct {
	config *Config
}

func (d *cloudInstanceGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_group"
}

func (d *cloudInstanceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about an instance group (placement group) in a public cloud project.",
		Attributes:  instanceGroupDataSourceAttributes(),
	}
}

func (d *cloudInstanceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instanceGroup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceGroupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// instanceGroupDataSourceAttributes builds the datasource schema: same shape as
// the resource model but service_name+id are the only inputs and everything else
// is Computed.
func instanceGroupDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"service_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: "Service name / cloud project id"},
		"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: "Instance group ID"},

		"region":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Region"},
		"name":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Instance group name"},
		"policy":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Placement policy (AFFINITY, ANTI_AFFINITY)"},
		"checksum":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Target-spec checksum"},
		"created_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Creation date"},
		"updated_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Last update date"},
		"resource_status": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Readiness status"},
		"current_state":   schema.SingleNestedAttribute{Computed: true, Description: "Observed state", Attributes: instanceGroupCurrentStateDataSourceSchemaAttributes()},
	}
}

// instanceGroupCurrentStateDataSourceSchemaAttributes returns the datasource-schema
// attributes for current_state. Shared by both the singular and plural data sources.
func instanceGroupCurrentStateDataSourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Observed instance group name"},
		"policy": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Observed placement policy"},
		"location": schema.SingleNestedAttribute{Computed: true, Description: "Observed location", Attributes: map[string]schema.Attribute{
			"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}},
		"members": schema.ListNestedAttribute{Computed: true, Description: "Instances currently belonging to the group", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}}},
	}
}
