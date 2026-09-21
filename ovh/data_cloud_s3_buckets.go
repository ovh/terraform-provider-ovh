package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudS3BucketsDataSource)(nil)

func NewCloudS3BucketsDataSource() datasource.DataSource {
	return &cloudS3BucketsDataSource{}
}

type cloudS3BucketsDataSource struct {
	config *Config
}

func (d *cloudS3BucketsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_s3_buckets"
}

func (d *cloudS3BucketsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudS3BucketsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List the S3 compatible object storage buckets in a public cloud project.",
		MarkdownDescription: "List the S3 compatible object storage buckets in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "If set, only buckets located in this region are returned",
				MarkdownDescription: "If set, only buckets located in this region are returned",
			},
			"buckets": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of buckets",
				MarkdownDescription: "List of buckets",
				NestedObject: schema.NestedAttributeObject{
					Attributes: s3BucketDataSourceAttributes(),
				},
			},
		},
	}
}

// cloudS3BucketsDataSourceModel is the Terraform state model for this data source.
type cloudS3BucketsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Buckets     types.List             `tfsdk:"buckets"`
}

// s3BucketListItemAttrTypes returns the attribute types for a single bucket item in the list.
func s3BucketListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"location":        types.ObjectType{AttrTypes: S3BucketLocationAttrTypes()},
		"owner_user_id":   ovhtypes.TfStringType{},
		"encryption":      types.ObjectType{AttrTypes: S3BucketEncryptionAttrTypes()},
		"versioning":      types.ObjectType{AttrTypes: S3BucketVersioningAttrTypes()},
		"object_lock":     types.ObjectType{AttrTypes: S3BucketObjectLockAttrTypes()},
		"tags":            types.MapType{ElemType: ovhtypes.TfStringType{}},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: S3BucketCurrentStateAttrTypes()},
	}
}

func (d *cloudS3BucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudS3BucketsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket"

	var apiBuckets []CloudS3BucketAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiBuckets); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	regionFilter := ""
	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		regionFilter = data.Region.ValueString()
	}

	bucketObjs := make([]attr.Value, 0, len(apiBuckets))
	for i := range apiBuckets {
		v := apiBuckets[i]

		region := ""
		if v.TargetSpec != nil && v.TargetSpec.Location != nil {
			region = v.TargetSpec.Location.Region
		}
		if v.CurrentState != nil && v.CurrentState.Location != nil && v.CurrentState.Location.Region != "" {
			region = v.CurrentState.Location.Region
		}
		if regionFilter != "" && region != regionFilter {
			continue
		}

		// Reuse the singular mapping to keep the shape consistent.
		var item cloudS3BucketDataSourceModel
		mapS3BucketToDataSourceModel(ctx, &v, &item)

		obj, diags := types.ObjectValue(
			s3BucketListItemAttrTypes(),
			map[string]attr.Value{
				"id":              item.Id,
				"name":            item.Name,
				"location":        item.Location,
				"owner_user_id":   item.OwnerUserId,
				"encryption":      item.Encryption,
				"versioning":      item.Versioning,
				"object_lock":     item.ObjectLock,
				"tags":            item.Tags,
				"checksum":        item.Checksum,
				"created_at":      item.CreatedAt,
				"updated_at":      item.UpdatedAt,
				"resource_status": item.ResourceStatus,
				"current_state":   item.CurrentState,
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		bucketObjs = append(bucketObjs, obj)
	}

	bucketsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: s3BucketListItemAttrTypes()},
		bucketObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Buckets = bucketsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
