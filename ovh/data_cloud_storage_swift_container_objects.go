package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageSwiftContainerObjectsDataSource)(nil)

func NewCloudStorageSwiftContainerObjectsDataSource() datasource.DataSource {
	return &cloudStorageSwiftContainerObjectsDataSource{}
}

type cloudStorageSwiftContainerObjectsDataSource struct {
	config *Config
}

func (d *cloudStorageSwiftContainerObjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_swift_container_objects"
}

func (d *cloudStorageSwiftContainerObjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageSwiftContainerObjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List the objects of a Swift object storage container in a public cloud project.",
		MarkdownDescription: "List the objects of a Swift object storage container in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"container_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Container identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Container identifier, in the `<REGION>_<name>` form",
			},
			"prefix": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "If set, only objects whose key starts with this prefix are returned",
				MarkdownDescription: "If set, only objects whose key starts with this prefix are returned",
			},
			"limit": schema.Int64Attribute{
				Optional:            true,
				Description:         "Maximum number of objects returned. When unset, every object of the container is returned",
				MarkdownDescription: "Maximum number of objects returned. When unset, every object of the container is returned",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"objects": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of objects",
				MarkdownDescription: "List of objects",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Object key, which is its full path within the container",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "Object size in bytes",
						},
						"content_type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "MIME type of the object content",
						},
						"last_modified_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last modification timestamp of the object",
						},
						"etag": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Entity tag of the object: the MD5 hash of its content, or a hash of the segment hashes for a large object",
						},
						"retrieval_state": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Retrieval state of the object in an archive container (SEALED, UNSEALING, UNSEALED), absent in a standard container",
						},
						"retrieval_delay": schema.Int64Attribute{
							Computed:    true,
							Description: "Seconds to wait before the object becomes readable, only set while retrieval_state is UNSEALING",
						},
					},
				},
			},
		},
	}
}

type cloudStorageSwiftContainerObjectsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ContainerId ovhtypes.TfStringValue `tfsdk:"container_id"`
	Prefix      ovhtypes.TfStringValue `tfsdk:"prefix"`
	Limit       types.Int64            `tfsdk:"limit"`
	Objects     types.List             `tfsdk:"objects"`
}

func (d *cloudStorageSwiftContainerObjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageSwiftContainerObjectsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseEndpoint := storageSwiftContainerEndpoint(data.ServiceName.ValueString()) +
		"/" + url.PathEscape(data.ContainerId.ValueString()) + "/object"

	limit := int64(0)
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		limit = data.Limit.ValueInt64()
	}

	// The route pages Swift-style: marker is the last key of the previous page.
	var objects []CloudStorageSwiftContainerAPIObject
	marker := ""
	for {
		queryParams := url.Values{}
		if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
			queryParams.Set("prefix", data.Prefix.ValueString())
		}
		if marker != "" {
			queryParams.Set("marker", marker)
		}
		if limit > 0 {
			queryParams.Set("limit", strconv.FormatInt(limit-int64(len(objects)), 10))
		}
		endpoint := baseEndpoint
		if len(queryParams) > 0 {
			endpoint += "?" + queryParams.Encode()
		}

		var page []CloudStorageSwiftContainerAPIObject
		if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &page); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}
		if len(page) == 0 {
			break
		}
		objects = append(objects, page...)

		if limit > 0 && int64(len(objects)) >= limit {
			objects = objects[:limit]
			break
		}
		// Guards against a route that ignores the marker and keeps returning the same page.
		if page[len(page)-1].Key == marker {
			break
		}
		marker = page[len(page)-1].Key
	}

	objectValues := make([]attr.Value, 0, len(objects))
	for i := range objects {
		objectValues = append(objectValues, buildStorageSwiftContainerObjectValue(&objects[i]))
	}

	objectsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: StorageSwiftContainerObjectAttrTypes()},
		objectValues,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Objects = objectsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
