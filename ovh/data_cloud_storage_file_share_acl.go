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

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileShareAclDataSource)(nil)

func NewCloudStorageFileShareAclDataSource() datasource.DataSource {
	return &cloudStorageFileShareAclDataSource{}
}

type cloudStorageFileShareAclDataSource struct {
	config *Config
}

func (d *cloudStorageFileShareAclDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share_acl"
}

func (d *cloudStorageFileShareAclDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fileShareAclDataSourceAttributes returns the computed attributes describing a
// file storage share access rule (ACL), shared between the singular and plural data sources.
func fileShareAclDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Access rule ID",
		},
		"access_to": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "IP address or CIDR allowed to access the file storage share",
		},
		"access_level": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Access level granted (READ_WRITE, READ_ONLY)",
		},
		"checksum": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Creation date of the access rule",
		},
		"updated_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last update date of the access rule",
		},
		"resource_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Access rule readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
		},
		"current_state": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current observed state of the access rule from the infrastructure",
			Attributes: map[string]schema.Attribute{
				"access_to": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "IP address or CIDR allowed to access the file storage share",
				},
				"access_level": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Access level granted",
				},
				"state": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Current state of the access rule (ACTIVE, APPLYING, DENYING, ERROR)",
				},
				"created_at": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Creation date of the access rule",
				},
			},
		},
	}
}

func (d *cloudStorageFileShareAclDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name of the resource representing the id of the cloud project",
			MarkdownDescription: "Service name of the resource representing the id of the cloud project",
		},
		"share_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "ID of the file storage share the access rule applies to",
			MarkdownDescription: "ID of the file storage share the access rule applies to",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Access rule ID",
			MarkdownDescription: "Access rule ID",
		},
	}

	for name, attribute := range fileShareAclDataSourceAttributes() {
		if name == "id" {
			continue
		}
		attrs[name] = attribute
	}

	resp.Schema = schema.Schema{
		Description:         "Get an access rule (ACL) of a public cloud file storage share.",
		MarkdownDescription: "Get an access rule (ACL) of a public cloud file storage share.",
		Attributes:          attrs,
	}
}

// cloudStorageFileShareAclDataSourceModel is the Terraform state model for this data source.
type cloudStorageFileShareAclDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	ShareId        ovhtypes.TfStringValue `tfsdk:"share_id"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	AccessTo       ovhtypes.TfStringValue `tfsdk:"access_to"`
	AccessLevel    ovhtypes.TfStringValue `tfsdk:"access_level"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

func (d *cloudStorageFileShareAclDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileShareAclDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fileStorageAclBaseEndpoint(data.ServiceName.ValueString(), data.ShareId.ValueString()) +
		"/" + url.PathEscape(data.Id.ValueString())

	var v CloudStorageFileShareAclAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	mapFileShareAclToDataSourceModel(&v, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapFileShareAclToDataSourceModel populates the data source model from the API response.
// ServiceName and ShareId come from the configuration and are never touched here.
func mapFileShareAclToDataSourceModel(v *CloudStorageFileShareAclAPIResponse, data *cloudStorageFileShareAclDataSourceModel) {
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)}
	data.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Checksum)}
	data.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CreatedAt)}
	data.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.UpdatedAt)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)}

	if v.TargetSpec != nil {
		data.AccessTo = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.AccessTo)}
		data.AccessLevel = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.AccessLevel)}
	}

	if v.CurrentState != nil {
		data.CurrentState = buildStorageFileShareAclCurrentStateObject(v.CurrentState)
	} else {
		data.CurrentState = types.ObjectNull(StorageFileShareAclCurrentStateAttrTypes())
	}
}
