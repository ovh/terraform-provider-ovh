package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketObjectsDataSource)(nil)

func NewCloudStorageObjectBucketObjectsDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketObjectsDataSource{}
}

type cloudStorageObjectBucketObjectsDataSource struct {
	config *Config
}

func (d *cloudStorageObjectBucketObjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
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

func (d *cloudStorageObjectBucketObjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List the objects of an S3 compatible object storage bucket in a public cloud project.",
		MarkdownDescription: "List the objects of an S3 compatible object storage bucket in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
			},
			"bucket_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Bucket identifier, in the `<REGION>_<name>` form",
			},
			"prefix": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "If set, only objects whose key starts with this prefix are returned",
			},
			"delimiter": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "If set, keys sharing a prefix up to the delimiter are grouped into one common prefix entry",
			},
			"with_versions": schema.BoolAttribute{
				Optional:    true,
				Description: "Include every object version and delete marker",
			},
			"objects": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of objects",
				NestedObject: schema.NestedAttributeObject{
					Attributes: storageObjectBucketObjectDataSourceAttributes(),
				},
			},
		},
	}
}

type cloudStorageObjectBucketObjectsDataSourceModel struct {
	ServiceName  ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketId     ovhtypes.TfStringValue `tfsdk:"bucket_id"`
	Prefix       ovhtypes.TfStringValue `tfsdk:"prefix"`
	Delimiter    ovhtypes.TfStringValue `tfsdk:"delimiter"`
	WithVersions types.Bool             `tfsdk:"with_versions"`
	Objects      types.List             `tfsdk:"objects"`
}

func (d *cloudStorageObjectBucketObjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketObjectsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() || data.ServiceName.ValueString() == "" {
		envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
		if envServiceName == "" {
			resp.Diagnostics.AddError(
				"Missing service_name",
				"The service_name attribute is required. Provide it in the data source "+
					"configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
			)
			return
		}
		data.ServiceName = ovhtypes.NewTfStringValue(envServiceName)
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketId.ValueString()) + "/object"

	queryParams := url.Values{}
	if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
		queryParams.Set("prefix", data.Prefix.ValueString())
	}
	if !data.Delimiter.IsNull() && !data.Delimiter.IsUnknown() {
		queryParams.Set("delimiter", data.Delimiter.ValueString())
	}
	if !data.WithVersions.IsNull() && !data.WithVersions.IsUnknown() {
		queryParams.Set("withVersions", strconv.FormatBool(data.WithVersions.ValueBool()))
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	apiObjects, err := helpers.GetAllPagesV2[CloudStorageObjectBucketObjectAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	objectValues := make([]attr.Value, 0, len(apiObjects))
	for i := range apiObjects {
		obj, diags := buildStorageObjectBucketObjectValue(&apiObjects[i])
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		objectValues = append(objectValues, obj)
	}

	objectsList, diags := types.ListValue(types.ObjectType{AttrTypes: storageObjectBucketObjectAttrTypes()}, objectValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Objects = objectsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
