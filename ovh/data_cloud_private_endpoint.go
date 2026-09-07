package ovh

import (
	"context"
	"fmt"
	"maps"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudPrivateEndpointDataSource)(nil)

func NewCloudPrivateEndpointDataSource() datasource.DataSource {
	return &cloudPrivateEndpointDataSource{}
}

type cloudPrivateEndpointDataSource struct {
	config *Config
}

func (d *cloudPrivateEndpointDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_private_endpoint"
}

func (d *cloudPrivateEndpointDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// privateEndpointComputedAttributes returns every attribute a private endpoint
// exports, excluding the service_name and id used to address it. Shared by the
// singular data source and by the elements of the plural one, so each
// description is written once.
func privateEndpointComputedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Private endpoint name",
		},
		"description": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Private endpoint description",
		},
		"region": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Region where the private endpoint is created",
		},
		"service": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Service exposed through the endpoint",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Name of the service",
				},
				"tier": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Tier of the service",
				},
			},
		},
		"network_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Network the endpoint port is allocated in",
		},
		"subnet_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Subnet the endpoint port is allocated in",
		},
		"vip_address": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Private IP address pinned on the endpoint port",
		},
		"public_dns_enabled": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether the endpoint hostname is registered in the public DNS zone",
		},
		"checksum": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Creation date of the private endpoint",
		},
		"updated_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last update date of the private endpoint",
		},
		"resource_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Private endpoint readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
		},
		"current_state": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current state of the private endpoint",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Private endpoint name",
				},
				"description": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Private endpoint description",
				},
				"endpoint_id": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Identifier of the endpoint on the service side, as pe-<id>",
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
					},
				},
				"service": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Service exposed through the endpoint",
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Name of the service",
						},
						"tier": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Tier of the service",
						},
						"port": schema.Int64Attribute{
							Computed:    true,
							Description: "Port the service listens on through the endpoint",
						},
						"transport": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Transport protocol of the service",
						},
						"bandwidth": schema.Int64Attribute{
							Computed:    true,
							Description: "Bandwidth of the tier, in Mbit/s. Null when the catalog does not declare it",
						},
					},
				},
				"network": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Network the endpoint port is attached to",
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Network identifier",
						},
					},
				},
				"subnet": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Subnet the endpoint port is attached to",
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Subnet identifier",
						},
					},
				},
				"vip_address": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Private IP address of the endpoint in the customer network",
				},
				"hostname": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Fully qualified domain name of the endpoint. Null until the service publishes it. See public_dns_enabled for public resolvability",
				},
				"public_dns_enabled": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether the endpoint hostname is registered in the public DNS zone",
				},
			},
		},
	}
}

func (d *cloudPrivateEndpointDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Required:    true,
			Description: "Service name of the resource representing the id of the cloud project",
		},
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Required:    true,
			Description: "Private endpoint ID",
		},
	}
	maps.Copy(attributes, privateEndpointComputedAttributes())

	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about a private endpoint in a public cloud project.",
		Attributes:  attributes,
	}
}

func (d *cloudPrivateEndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// The resource model carries exactly the attributes of this schema, so it
	// doubles as the data source model and MergeWith is reused as-is.
	var data CloudPrivateEndpointModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudPrivateEndpointAPIResponse
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
