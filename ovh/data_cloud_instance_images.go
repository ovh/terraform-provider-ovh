package ovh

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceImagesDataSource)(nil)

func NewCloudInstanceImagesDataSource() datasource.DataSource {
	return &cloudInstanceImagesDataSource{}
}

type cloudInstanceImagesDataSource struct {
	config *Config
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
	resp.Schema = CloudInstanceImagesDataSourceSchema(ctx)
}

func (d *cloudInstanceImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceImagesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/reference/instance/image?region=" + url.QueryEscape(data.RegionName.ValueString())

	var arr []CloudInstanceImagesValue
	if err := d.config.OVHClient.Get(endpoint, &arr); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Filter by name regexp if provided
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		namePattern := data.Name.ValueString()
		re, err := regexp.Compile(namePattern)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid name filter",
				fmt.Sprintf("The name filter %q is not a valid regular expression: %s", namePattern, err.Error()),
			)
			return
		}

		filtered := make([]CloudInstanceImagesValue, 0, len(arr))
		for _, img := range arr {
			if re.MatchString(img.Name.ValueString()) {
				filtered = append(filtered, img)
			}
		}
		arr = filtered
	}

	var b []attr.Value
	for _, a := range arr {
		b = append(b, a)
	}

	data.Images = ovhtypes.TfListNestedValue[CloudInstanceImagesValue]{
		ListValue: basetypes.NewListValueMust(CloudInstanceImagesValue{}.Type(ctx), b),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
