package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceDataSource)(nil)

func NewCloudInstanceDataSource() datasource.DataSource {
	return &cloudInstanceDataSource{}
}

type cloudInstanceDataSource struct {
	config *Config
}

func (d *cloudInstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance"
}

func (d *cloudInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about an instance in a public cloud project.",
		Attributes:  instanceDataSourceAttributes(ctx),
	}
}

func (d *cloudInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// instanceDataSourceAttributes builds the datasource schema: same shape as the
// resource but service_name+id are the only inputs and everything else is Computed.
func instanceDataSourceAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"service_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: "Service name / cloud project id"},
		"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: "Instance ID"},

		"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Region"},
		"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Availability zone"},
		"ssh_key_name":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "SSH key name"},
		"group_id":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Placement group ID"},
		"name":              schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Instance name"},
		"flavor_id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Flavor ID"},
		"image_id":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Image ID"},
		"power_state":       schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Desired power state"},
		"networks": schema.ListNestedAttribute{Computed: true, Description: "Requested network interfaces", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"public":         schema.BoolAttribute{Computed: true},
			"network_id":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"subnet_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"floating_ip_id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}}},
		"volume_ids":         schema.ListAttribute{CustomType: ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx), Computed: true, Description: "Attached volume IDs"},
		"security_group_ids": schema.ListAttribute{CustomType: ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx), Computed: true, Description: "Security group IDs"},
		"shares": schema.ListNestedAttribute{Computed: true, Description: "Attached shares", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"access_level": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}}},
		"checksum":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Target-spec checksum"},
		"created_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Creation date"},
		"updated_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Last update date"},
		"resource_status": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Readiness status"},
		"current_state":   schema.SingleNestedAttribute{Computed: true, Description: "Observed state", Attributes: instanceCurrentStateDataSourceSchemaAttributes()},
	}
}

// instanceCurrentStateDataSourceSchemaAttributes mirrors
// instanceCurrentStateSchemaAttributes but returns datasource-schema attributes
// (the resource and datasource schema packages define distinct Attribute types).
// Shared by both the singular and plural instance data sources.
func instanceCurrentStateDataSourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Observed instance name"},
		"power_state":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Observed power state"},
		"locked":       schema.BoolAttribute{Computed: true, Description: "Whether the instance is locked"},
		"ssh_key_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "SSH key injected at boot"},
		"host_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Opaque physical host ID"},
		"project_id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Owning project ID"},
		"user_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Owning OpenStack user ID"},
		"flavor": schema.SingleNestedAttribute{Computed: true, Description: "Observed flavor details", Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"name":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"vcpus":     schema.Int64Attribute{Computed: true},
			"ram":       schema.Int64Attribute{Computed: true},
			"disk":      schema.Int64Attribute{Computed: true},
			"swap":      schema.Int64Attribute{Computed: true},
			"ephemeral": schema.Int64Attribute{Computed: true},
		}},
		"image": schema.SingleNestedAttribute{Computed: true, Description: "Observed image details (null for boot-from-volume)", Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"name":       schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"size":       schema.Int64Attribute{Computed: true},
			"status":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"deprecated": schema.BoolAttribute{Computed: true},
		}},
		"location": schema.SingleNestedAttribute{Computed: true, Description: "Observed location", Attributes: map[string]schema.Attribute{
			"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}},
		"networks": schema.ListNestedAttribute{Computed: true, Description: "Observed network interfaces", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"public":         schema.BoolAttribute{Computed: true},
			"subnet_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"gateway_id":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"floating_ip_id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"addresses": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"ip":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
				"mac":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
				"type":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
				"version": schema.Int64Attribute{Computed: true},
			}}},
		}}},
		"volumes": schema.ListNestedAttribute{Computed: true, Description: "Observed attached volumes", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"size": schema.Int64Attribute{Computed: true},
		}}},
		"shares": schema.ListNestedAttribute{Computed: true, Description: "Observed attached shares (only populated on single-instance read)", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"access_level": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"access_to":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"state":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}}},
		"security_groups": schema.ListNestedAttribute{Computed: true, Description: "Observed security groups", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}}},
		"group": schema.SingleNestedAttribute{Computed: true, Description: "Placement group membership", Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}},
	}
}
