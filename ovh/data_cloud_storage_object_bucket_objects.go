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

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketObjectsDataSource)(nil)

func NewCloudStorageObjectBucketObjectsDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketObjectsDataSource{}
}

type cloudStorageObjectBucketObjectsDataSource struct {
	config *Config
}

// cloudStorageObjectBucketObjectsDataSourceModel is the Terraform state model
// for the list-objects data source.
type cloudStorageObjectBucketObjectsDataSourceModel struct {
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketName      ovhtypes.TfStringValue `tfsdk:"bucket_name"`
	Prefix          ovhtypes.TfStringValue `tfsdk:"prefix"`
	Delimiter       ovhtypes.TfStringValue `tfsdk:"delimiter"`
	Limit           types.Int64            `tfsdk:"limit"`
	KeyMarker       ovhtypes.TfStringValue `tfsdk:"key_marker"`
	VersionIdMarker ovhtypes.TfStringValue `tfsdk:"version_id_marker"`
	WithVersions    types.Bool             `tfsdk:"with_versions"`

	Objects       types.List             `tfsdk:"objects"`
	NextKeyMarker ovhtypes.TfStringValue `tfsdk:"next_key_marker"`
	IsTruncated   types.Bool             `tfsdk:"is_truncated"`
}

func (d *cloudStorageObjectBucketObjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_objects"
}

func (d *cloudStorageObjectBucketObjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketObjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists objects in an S3 bucket using the v2 API. Pagination is a single page controlled by the " +
			"`limit`, `key_marker` and `version_id_marker` inputs; use `next_key_marker` (and `is_truncated`) from the " +
			"response to fetch subsequent pages.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"bucket_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Name of the bucket to list objects from",
			},
			"prefix": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Return only objects whose key begins with this prefix",
			},
			"delimiter": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Character to group keys into common prefixes (directory-like listing)",
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
				Description: "Version ID to resume listing from (only used when `with_versions` is true and `key_marker` is set)",
			},
			"with_versions": schema.BoolAttribute{
				Optional:    true,
				Description: "When true, list all versions (and delete markers) rather than only current objects",
			},

			"objects": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "Objects, common prefixes and (when `with_versions` is set) delete markers returned for the current page",
				NestedObject: schema.NestedAttributeObject{Attributes: objectDataSourceAttributesForList()},
			},
			"next_key_marker": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Key to pass as `key_marker` for the next page; empty when `is_truncated` is false",
			},
			"is_truncated": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether more objects exist beyond the current page",
			},
		},
	}
}

// objectDataSourceAttributesForList returns the schema for one element of the
// list produced by the objects / object_versions data sources. It mirrors the
// attribute types declared in cloudBucketObjectAttrTypes.
func objectDataSourceAttributesForList() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"e_tag": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Entity tag (ETag) of the object",
		},
		"is_common_prefix": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this entry is a common-prefix grouping rather than an object",
		},
		"is_delete_marker": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this version is a delete marker",
		},
		"is_latest": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this version is the latest version of the object",
		},
		"key": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Object key (path within the bucket)",
		},
		"last_modified": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Date and time at which the object was last modified",
		},
		"legal_hold": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Current legal hold status (ON, OFF)",
		},
		"replication_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Replication status (COMPLETED, FAILED, NONE, PENDING, REPLICA)",
		},
		"restore_status": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Restore status for objects stored in archive tiers",
			Attributes: map[string]schema.Attribute{
				"expire_date": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Date and time at which the restored copy expires",
				},
				"in_progress": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether a restore operation is currently in progress",
				},
			},
		},
		"retention": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current retention (object lock) configuration",
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Retention mode (COMPLIANCE, GOVERNANCE, NONE)",
				},
				"retain_until_date": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Date and time until which the object is retained",
				},
			},
		},
		"size": schema.Int64Attribute{
			Computed:    true,
			Description: "Size of the object in bytes",
		},
		"storage_class": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Storage class (DEEP_ARCHIVE, GLACIER, GLACIER_IR, HIGH_PERF, INTELLIGENT_TIERING, ONEZONE_IA, STANDARD, STANDARD_IA)",
		},
		"version_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Version identifier of the object (when versioning is enabled)",
		},
	}
}

func (d *cloudStorageObjectBucketObjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketObjectsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketName.ValueString()) + "/object"

	query := url.Values{}
	if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
		query.Set("prefix", data.Prefix.ValueString())
	}
	if !data.Delimiter.IsNull() && !data.Delimiter.IsUnknown() {
		query.Set("delimiter", data.Delimiter.ValueString())
	}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		query.Set("limit", fmt.Sprintf("%d", data.Limit.ValueInt64()))
	}
	if !data.KeyMarker.IsNull() && !data.KeyMarker.IsUnknown() {
		query.Set("keyMarker", data.KeyMarker.ValueString())
	}
	if !data.VersionIdMarker.IsNull() && !data.VersionIdMarker.IsUnknown() {
		query.Set("versionIdMarker", data.VersionIdMarker.ValueString())
	}
	if !data.WithVersions.IsNull() && !data.WithVersions.IsUnknown() && data.WithVersions.ValueBool() {
		query.Set("withVersions", "true")
	}

	if encoded := query.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	var response CloudBucketObjectListResponse
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
