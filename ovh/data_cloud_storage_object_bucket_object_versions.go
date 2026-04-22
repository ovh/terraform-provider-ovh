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

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketObjectVersionsDataSource)(nil)

func NewCloudStorageObjectBucketObjectVersionsDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketObjectVersionsDataSource{}
}

type cloudStorageObjectBucketObjectVersionsDataSource struct {
	config *Config
}

// cloudStorageObjectBucketObjectVersionsDataSourceModel is the Terraform state
// model for the list-object-versions data source.
type cloudStorageObjectBucketObjectVersionsDataSourceModel struct {
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketName      ovhtypes.TfStringValue `tfsdk:"bucket_name"`
	ObjectKey       ovhtypes.TfStringValue `tfsdk:"object_key"`
	Limit           types.Int64            `tfsdk:"limit"`
	KeyMarker       ovhtypes.TfStringValue `tfsdk:"key_marker"`
	VersionIdMarker ovhtypes.TfStringValue `tfsdk:"version_id_marker"`

	Objects             types.List             `tfsdk:"objects"`
	NextKeyMarker       ovhtypes.TfStringValue `tfsdk:"next_key_marker"`
	NextVersionIdMarker ovhtypes.TfStringValue `tfsdk:"next_version_id_marker"`
	IsTruncated         types.Bool             `tfsdk:"is_truncated"`
}

func (d *cloudStorageObjectBucketObjectVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_object_versions"
}

func (d *cloudStorageObjectBucketObjectVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketObjectVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists versions (and delete markers) of a single S3 object using the v2 API. " +
			"Pagination is a single page controlled by the `limit`, `key_marker` and `version_id_marker` inputs.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"bucket_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Name of the bucket",
			},
			"object_key": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Object key whose versions to list",
			},
			"limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of entries to return (API hard cap: 1000)",
			},
			"key_marker": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Key to resume listing from (pagination)",
			},
			"version_id_marker": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Version ID to resume listing from (must be paired with `key_marker`)",
			},

			"objects": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "Object versions and delete markers returned for the current page",
				NestedObject: schema.NestedAttributeObject{Attributes: objectDataSourceAttributesForList()},
			},
			"next_key_marker": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Key to pass as `key_marker` for the next page; empty when `is_truncated` is false",
			},
			"next_version_id_marker": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Version ID to pass as `version_id_marker` for the next page; empty when `is_truncated` is false",
			},
			"is_truncated": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether more versions exist beyond the current page",
			},
		},
	}
}

func (d *cloudStorageObjectBucketObjectVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketObjectVersionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketName.ValueString()) +
		"/object/" + url.PathEscape(data.ObjectKey.ValueString()) + "/version"

	query := url.Values{}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		query.Set("limit", fmt.Sprintf("%d", data.Limit.ValueInt64()))
	}
	if !data.KeyMarker.IsNull() && !data.KeyMarker.IsUnknown() {
		query.Set("keyMarker", data.KeyMarker.ValueString())
	}
	if !data.VersionIdMarker.IsNull() && !data.VersionIdMarker.IsUnknown() {
		query.Set("versionIdMarker", data.VersionIdMarker.ValueString())
	}

	if encoded := query.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	var response CloudBucketObjectVersionListResponse
	if err := d.config.OVHClient.Get(endpoint, &response); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Objects = apiObjectListToTFList(response.Objects)
	data.IsTruncated = types.BoolValue(response.IsTruncated)
	if response.NextKeyMarker != nil {
		data.NextKeyMarker = ovhtypes.TfStringValue{StringValue: types.StringValue(*response.NextKeyMarker)}
	} else {
		data.NextKeyMarker = ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	if response.NextVersionIDMarker != nil {
		data.NextVersionIdMarker = ovhtypes.TfStringValue{StringValue: types.StringValue(*response.NextVersionIDMarker)}
	} else {
		data.NextVersionIdMarker = ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
