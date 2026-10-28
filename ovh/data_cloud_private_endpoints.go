package ovh

import (
	"context"
	"fmt"
	"maps"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudPrivateEndpointsDataSource)(nil)

func NewCloudPrivateEndpointsDataSource() datasource.DataSource {
	return &cloudPrivateEndpointsDataSource{}
}

type cloudPrivateEndpointsDataSource struct {
	config *Config
}

// cloudPrivateEndpointsDataSourceModel is the model for the plural private
// endpoints data source.
type cloudPrivateEndpointsDataSourceModel struct {
	ServiceName      ovhtypes.TfStringValue `tfsdk:"service_name"`
	PrivateEndpoints types.List             `tfsdk:"private_endpoints"`
}

func (d *cloudPrivateEndpointsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_private_endpoints"
}

func (d *cloudPrivateEndpointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// privateEndpointListItemAttrTypes returns the attribute types for a single
// private endpoint element of the plural data source list.
func privateEndpointListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 ovhtypes.TfStringType{},
		"name":               ovhtypes.TfStringType{},
		"description":        ovhtypes.TfStringType{},
		"region":             ovhtypes.TfStringType{},
		"service":            types.ObjectType{AttrTypes: privateEndpointServiceAttrTypes()},
		"network_id":         ovhtypes.TfStringType{},
		"subnet_id":          ovhtypes.TfStringType{},
		"vip_address":        ovhtypes.TfStringType{},
		"public_dns_enabled": types.BoolType,
		"checksum":           ovhtypes.TfStringType{},
		"created_at":         ovhtypes.TfStringType{},
		"updated_at":         ovhtypes.TfStringType{},
		"resource_status":    ovhtypes.TfStringType{},
		"current_state":      types.ObjectType{AttrTypes: privateEndpointCurrentStateAttrTypes()},
	}
}

func (d *cloudPrivateEndpointsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	itemAttributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Private endpoint ID",
		},
	}
	maps.Copy(itemAttributes, privateEndpointComputedAttributes())

	resp.Schema = schema.Schema{
		Description: "Use this data source to list the private endpoints of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"private_endpoints": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of private endpoints",
				NestedObject: schema.NestedAttributeObject{Attributes: itemAttributes},
			},
		},
	}
}

func (d *cloudPrivateEndpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudPrivateEndpointsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint"

	responseData, err := helpers.GetAllPagesV2[CloudPrivateEndpointAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		item, diags := buildPrivateEndpointListItemObject(ctx, &responseData[i])
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, item)
	}

	endpoints, diags := types.ListValue(types.ObjectType{AttrTypes: privateEndpointListItemAttrTypes()}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.PrivateEndpoints = endpoints

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildPrivateEndpointListItemObject maps one API response onto a list element.
// It reuses the resource model's mapping so the plural data source cannot drift
// from the resource.
func buildPrivateEndpointListItemObject(ctx context.Context, response *CloudPrivateEndpointAPIResponse) (basetypes.ObjectValue, diag.Diagnostics) {
	// MergeWith leaves an attribute untouched when the response omits it, so
	// seed the nested objects with typed nulls: a zero value carries no
	// attribute types and would fail the object type check below.
	m := CloudPrivateEndpointModel{
		Service:      types.ObjectNull(privateEndpointServiceAttrTypes()),
		CurrentState: types.ObjectNull(privateEndpointCurrentStateAttrTypes()),
	}
	m.MergeWith(ctx, response)

	return types.ObjectValue(privateEndpointListItemAttrTypes(), map[string]attr.Value{
		"id":                 m.Id,
		"name":               m.Name,
		"description":        m.Description,
		"region":             m.Region,
		"service":            m.Service,
		"network_id":         m.NetworkId,
		"subnet_id":          m.SubnetId,
		"vip_address":        m.VipAddress,
		"public_dns_enabled": m.PublicDnsEnabled,
		"checksum":           m.Checksum,
		"created_at":         m.CreatedAt,
		"updated_at":         m.UpdatedAt,
		"resource_status":    m.ResourceStatus,
		"current_state":      m.CurrentState,
	})
}
