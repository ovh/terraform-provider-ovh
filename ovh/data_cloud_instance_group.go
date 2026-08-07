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
	attrs := instanceGroupDataSourceAttributes()
	attrs["service_name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Service name of the resource representing the id of the cloud project",
		MarkdownDescription: "Service name of the resource representing the id of the cloud project",
	}
	attrs["id"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Unique identifier of the instance group to look up.",
		MarkdownDescription: "Unique identifier of the instance group to look up.",
	}

	resp.Schema = schema.Schema{
		Description:         "Use this data source to retrieve information about an instance group (placement group) in a public cloud project. An instance group is immutable once created: there is no update route and its membership is fixed at instance-creation time.",
		MarkdownDescription: "Use this data source to retrieve information about an instance group (placement group) in a public cloud project. An instance group is immutable once created: there is no update route and its membership is fixed at instance-creation time.",
		Attributes:          attrs,
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

func dsString(description string) schema.StringAttribute {
	return schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Computed:            true,
		Description:         description,
		MarkdownDescription: description,
	}
}

// instanceGroupDataSourceAttributes builds the datasource schema: same shape as
// the resource model but service_name+id are the only inputs and everything else
// is Computed. Shared by the singular data source (which overrides service_name
// and id to Required) and the `instance_groups[]` nested object of the plural one.
func instanceGroupDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":     dsString("Unique identifier of the instance group, assigned at creation and used to reference it on the GET/DELETE routes."),
		"region": dsString("Region where the instance group, and therefore its member instances, are placed. Immutable after creation: a group cannot be moved to another region, and its members must be created in this same location."),
		"name":   dsString("Display name of the instance group. Immutable after creation, as instance groups cannot be updated."),
		"policy": dsString("Placement policy applied to the group's member instances: AFFINITY packs members onto the same hypervisor, ANTI_AFFINITY spreads members across distinct hypervisors. Maps to the underlying OpenStack/Nova server group policy."),
		// The API deliberately omits updatedAt here: an instance group has no
		// update route, so the target spec (and its checksum) never changes.
		"checksum":        dsString("Computed hash of the current target specification, used for optimistic concurrency control. Because an instance group has no update route, this value never changes after creation."),
		"created_at":      dsString("Timestamp at which the instance group was created, in RFC 3339 format."),
		"resource_status": dsString("Instance group readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY)."),
		"current_state": schema.SingleNestedAttribute{
			Computed:            true,
			Description:         "State of the instance group as observed on the backend.",
			MarkdownDescription: "State of the instance group as observed on the backend.",
			Attributes:          instanceGroupCurrentStateDataSourceSchemaAttributes(),
		},
	}
}

// instanceGroupCurrentStateDataSourceSchemaAttributes returns the datasource-schema
// attributes for current_state. Shared by both the singular and plural data sources.
func instanceGroupCurrentStateDataSourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":   dsString("Display name of the instance group as reported by the backend. Fixed for the lifetime of the group, since instance groups cannot be updated."),
		"policy": dsString("Placement policy currently enforced for the group's members (AFFINITY, ANTI_AFFINITY). Mirrors the underlying OpenStack/Nova server group policy and is fixed at creation."),
		"location": schema.SingleNestedAttribute{
			Computed:            true,
			Description:         "Region (and, where applicable, availability zone) where the instance group and its member instances are deployed, as observed on the backend.",
			MarkdownDescription: "Region (and, where applicable, availability zone) where the instance group and its member instances are deployed, as observed on the backend.",
			Attributes: map[string]schema.Attribute{
				"region":            dsString(instanceCatalogLocationRegionDescription),
				"availability_zone": dsString(instanceCatalogLocationAZDescription),
			},
		},
		"members": schema.ListNestedAttribute{
			Computed:            true,
			Description:         "Instances currently belonging to this group. Membership is determined at instance-creation time via the instance's group field and cannot be changed afterwards. Empty when the group has no member instances.",
			MarkdownDescription: "Instances currently belonging to this group. Membership is determined at instance-creation time via the instance's `group` field and cannot be changed afterwards. Empty when the group has no member instances.",
			NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"id": dsString("Identifier of the member instance."),
			}},
		},
	}
}
