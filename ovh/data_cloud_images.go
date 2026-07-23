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

var _ datasource.DataSourceWithConfigure = (*cloudImagesDataSource)(nil)

func NewCloudImagesDataSource() datasource.DataSource {
	return &cloudImagesDataSource{}
}

type cloudImagesDataSource struct {
	config *Config
}

// CloudImagesModel is the model for the plural images data source.
type CloudImagesModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Images      types.List             `tfsdk:"images"`
}

func (d *cloudImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_images"
}

func (d *cloudImagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the images available in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Filter the images by region",
			},
			"images": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of images",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Image ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Image name",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Availability status of the image",
						},
						"visibility": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Visibility scope of the image",
						},
						"min_disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum root disk size (GB) required to boot from this image",
						},
						"min_ram": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum RAM (MB) required to boot from this image",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the image (bytes)",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the image",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the image",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region where the image is offered",
						},
					},
				},
			},
		},
	}
}

func (d *cloudImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudImagesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/instance/image"
	if !data.Region.IsNull() && data.Region.ValueString() != "" {
		endpoint += "?region=" + url.QueryEscape(data.Region.ValueString())
	}

	var responseData []CloudImageAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	itemObjType := types.ObjectType{AttrTypes: imageAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		m := buildImageModel(&responseData[i])
		obj, diags := types.ObjectValue(imageAttrTypes(), map[string]attr.Value{
			"id":         m.Id,
			"name":       m.Name,
			"status":     m.Status,
			"visibility": m.Visibility,
			"min_disk":   m.MinDisk,
			"min_ram":    m.MinRam,
			"size":       m.Size,
			"created_at": m.CreatedAt,
			"updated_at": m.UpdatedAt,
			"region":     m.Region,
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
