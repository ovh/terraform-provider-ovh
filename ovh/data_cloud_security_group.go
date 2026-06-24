package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudSecurityGroupDataSource)(nil)

func NewCloudSecurityGroupDataSource() datasource.DataSource {
	return &cloudSecurityGroupDataSource{}
}

type cloudSecurityGroupDataSource struct {
	config *Config
}

func (d *cloudSecurityGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_security_group"
}

func (d *cloudSecurityGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// securityGroupDataSourceRuleSchemaAttributes returns the schema attributes for
// a computed rule (used in both the root rule list and the current_state rule lists).
func securityGroupDataSourceRuleSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Rule ID",
		},
		"direction": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Direction of the rule",
		},
		"ethernet_type": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Ethernet type",
		},
		"protocol": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Protocol",
		},
		"port_range_min": schema.Int64Attribute{
			Computed:    true,
			Description: "Minimum port number",
		},
		"port_range_max": schema.Int64Attribute{
			Computed:    true,
			Description: "Maximum port number",
		},
		"remote_group_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Remote security group ID",
		},
		"remote_ip_prefix": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Remote IP prefix",
		},
		"description": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Description of the rule",
		},
	}
}

// securityGroupDataSourceCurrentStateSchema returns the schema for the computed
// current_state nested object, reused by the singular and plural data sources.
func securityGroupDataSourceCurrentStateSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Current state of the security group",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Name of the security group",
			},
			"description": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Description of the security group",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Region of the security group",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "User-specified security group rules with their IDs",
				NestedObject: schema.NestedAttributeObject{
					Attributes: securityGroupDataSourceRuleSchemaAttributes(),
				},
			},
			"default_rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Default egress rules auto-created by OpenStack",
				NestedObject: schema.NestedAttributeObject{
					Attributes: securityGroupDataSourceRuleSchemaAttributes(),
				},
			},
		},
	}
}

func (d *cloudSecurityGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about a security group in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Security group ID",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Region of the security group",
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Name of the security group",
			},
			"description": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Description of the security group",
			},
			"rule": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of security group rules",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Direction of the rule",
						},
						"ethernet_type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Ethernet type",
						},
						"protocol": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Protocol",
						},
						"port_range_min": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum port number",
						},
						"port_range_max": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum port number",
						},
						"remote_group_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Remote security group ID",
						},
						"remote_ip_prefix": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Remote IP prefix",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Description of the rule",
						},
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
				Description: "Creation date of the security group",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date of the security group",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Security group readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
			},
			"current_state": securityGroupDataSourceCurrentStateSchema(),
		},
	}
}

func (d *cloudSecurityGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudSecurityGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudSecurityGroupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
