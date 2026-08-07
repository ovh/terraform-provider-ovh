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

var _ datasource.DataSourceWithConfigure = (*cloudInstanceImagesDataSource)(nil)

func NewCloudInstanceImagesDataSource() datasource.DataSource {
	return &cloudInstanceImagesDataSource{}
}

type cloudInstanceImagesDataSource struct {
	config *Config
}

// CloudInstanceImagesModel is the model for the plural images data source.
type CloudInstanceImagesModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Images      types.List             `tfsdk:"images"`
}

func (d *cloudInstanceImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_images"
}

func (d *cloudInstanceImagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to list the images available in a public cloud project. This is read-only reference data: each entry describes a bootable image the project can create instances from, not an existing resource.",
		MarkdownDescription: "Use this data source to list the images available in a public cloud project. This is read-only reference data: each entry describes a bootable image the project can create instances from, not an existing resource.",
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
				Description:         "Restrict the listing to the images offered in this region. The catalog is per-region: an image returned for one region is not guaranteed to exist in another.",
				MarkdownDescription: "Restrict the listing to the images offered in this region. The catalog is per-region: an image returned for one region is not guaranteed to exist in another.",
			},
			"images": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "Images advertised by the backend for this project, one entry per region the image is offered in.",
				MarkdownDescription: "Images advertised by the backend for this project, one entry per region the image is offered in.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: instanceImageDataSourceAttributes(),
				},
			},
		},
	}
}

func (d *cloudInstanceImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceImagesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/instance/image"
	if !data.Region.IsNull() && data.Region.ValueString() != "" {
		endpoint += "?region=" + url.QueryEscape(data.Region.ValueString())
	}

	var responseData []CloudInstanceImageAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	itemObjType := types.ObjectType{AttrTypes: instanceImageItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		var m CloudInstanceImageModel
		m.MergeWith(ctx, &responseData[i])
		obj, diags := types.ObjectValue(instanceImageItemAttrTypes(), map[string]attr.Value{
			"id":         m.Id,
			"name":       m.Name,
			"status":     m.Status,
			"visibility": m.Visibility,
			"min_disk":   m.MinDisk,
			"min_ram":    m.MinRam,
			"size":       m.Size,
			"created_at": m.CreatedAt,
			"updated_at": m.UpdatedAt,
			"location":   m.Location,
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, obj)
	}

	data.Images = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
