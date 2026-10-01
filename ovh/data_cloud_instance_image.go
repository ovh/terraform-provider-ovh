package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceImageDataSource)(nil)

func NewCloudInstanceImageDataSource() datasource.DataSource {
	return &cloudInstanceImageDataSource{}
}

type cloudInstanceImageDataSource struct {
	config *Config
}

func (d *cloudInstanceImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_image"
}

func (d *cloudInstanceImageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// instanceImageDataSourceAttributes returns the image attributes as fully
// Computed. Shared by the singular data source (which overrides id to Required)
// and the `images[]` nested object of the plural one.
func instanceImageDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "The OpenStack/Glance image ID. Stable within a region and used to reference the image when creating an instance.",
			MarkdownDescription: "The OpenStack/Glance image ID. Stable within a region and used to reference the image when creating an instance.",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Display name of the image as reported by the backend (for example the distribution and version, such as Debian 12).",
			MarkdownDescription: "Display name of the image as reported by the backend (for example the distribution and version, such as `Debian 12`).",
		},
		"status": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Availability status of the image as reported by the backend. Only images in an active status can be used to create an instance.",
			MarkdownDescription: "Availability status of the image as reported by the backend. Only images in an active status can be used to create an instance.",
		},
		"visibility": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Visibility scope of the image, for example whether it is a public OVHcloud-provided image or private to the project.",
			MarkdownDescription: "Visibility scope of the image, for example whether it is a public OVHcloud-provided image or private to the project.",
		},
		"min_disk": schema.Int64Attribute{
			Computed:            true,
			Description:         "Minimum root disk size, in GB, that an instance must provide to boot from this image. A flavor whose disk is smaller than this value cannot be used with the image.",
			MarkdownDescription: "Minimum root disk size, in GB, that an instance must provide to boot from this image. A flavor whose disk is smaller than this value cannot be used with the image.",
		},
		"min_ram": schema.Int64Attribute{
			Computed:            true,
			Description:         "Minimum amount of memory, in MB, that an instance must provide to boot from this image. A flavor whose RAM is below this value cannot be used with the image.",
			MarkdownDescription: "Minimum amount of memory, in MB, that an instance must provide to boot from this image. A flavor whose RAM is below this value cannot be used with the image.",
		},
		"size": schema.Int64Attribute{
			Computed:            true,
			Description:         "Size of the image on the backend, expressed in bytes.",
			MarkdownDescription: "Size of the image on the backend, expressed in bytes.",
		},
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Timestamp at which the image was created on the backend, in RFC 3339 format.",
			MarkdownDescription: "Timestamp at which the image was created on the backend, in RFC 3339 format.",
		},
		"updated_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Timestamp of the last modification of the image on the backend, in RFC 3339 format.",
			MarkdownDescription: "Timestamp of the last modification of the image on the backend, in RFC 3339 format.",
		},
		"location": instanceCatalogLocationDataSourceAttribute(
			"Region (and, where applicable, availability zone) where this image is offered. The image catalog is per-region: an image returned for one region is not guaranteed to exist in another.",
		),
	}
}

func (d *cloudInstanceImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := instanceImageDataSourceAttributes()
	attrs["service_name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Service name of the resource representing the id of the cloud project",
		MarkdownDescription: "Service name of the resource representing the id of the cloud project",
	}
	attrs["id"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "The OpenStack/Glance image ID to look up. Stable within a region and used to reference the image when creating an instance.",
		MarkdownDescription: "The OpenStack/Glance image ID to look up. Stable within a region and used to reference the image when creating an instance.",
	}

	resp.Schema = schema.Schema{
		Description:         "Use this data source to retrieve information about an image available in a public cloud project. This is read-only reference data: it describes a bootable image the project can create instances from, not an existing resource.",
		MarkdownDescription: "Use this data source to retrieve information about an image available in a public cloud project. This is read-only reference data: it describes a bootable image the project can create instances from, not an existing resource.",
		Attributes:          attrs,
	}
}

func (d *cloudInstanceImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceImageModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/instance/image/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceImageAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
