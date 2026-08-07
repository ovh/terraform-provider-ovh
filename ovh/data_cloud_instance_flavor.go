package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceFlavorDataSource)(nil)

func NewCloudInstanceFlavorDataSource() datasource.DataSource {
	return &cloudInstanceFlavorDataSource{}
}

type cloudInstanceFlavorDataSource struct {
	config *Config
}

func (d *cloudInstanceFlavorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_flavor"
}

func (d *cloudInstanceFlavorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

const (
	instanceCatalogLocationRegionDescription = "Region code"
	instanceCatalogLocationAZDescription     = "Availability zone within the region"
)

func instanceCatalogLocationDataSourceAttribute(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:            true,
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         instanceCatalogLocationRegionDescription,
				MarkdownDescription: instanceCatalogLocationRegionDescription,
			},
			"availability_zone": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         instanceCatalogLocationAZDescription,
				MarkdownDescription: instanceCatalogLocationAZDescription,
			},
		},
	}
}

// instanceFlavorDataSourceAttributes returns the flavor attributes as fully
// Computed. Shared by the singular data source (which overrides id to Required)
// and the `flavors[]` nested object of the plural one.
func instanceFlavorDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "The OpenStack/Nova flavor ID. Stable within a region and used to reference the flavor when creating an instance.",
			MarkdownDescription: "The OpenStack/Nova flavor ID. Stable within a region and used to reference the flavor when creating an instance.",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "The backend flavor name (for example b2-7, c2-15). This is the commercial/technical name used to identify the sizing in the catalog.",
			MarkdownDescription: "The backend flavor name (for example `b2-7`, `c2-15`). This is the commercial/technical name used to identify the sizing in the catalog.",
		},
		"vcpus": schema.Int64Attribute{
			Computed:            true,
			Description:         "Number of virtual CPUs provided by the flavor.",
			MarkdownDescription: "Number of virtual CPUs provided by the flavor.",
		},
		"ram": schema.Int64Attribute{
			Computed:            true,
			Description:         "Amount of memory provided by the flavor, expressed in MB.",
			MarkdownDescription: "Amount of memory provided by the flavor, expressed in MB.",
		},
		"disk": schema.Int64Attribute{
			Computed:            true,
			Description:         "Size of the flavor's root disk in GB. This is the primary system disk provisioned for instances created from this flavor.",
			MarkdownDescription: "Size of the flavor's root disk in GB. This is the primary system disk provisioned for instances created from this flavor.",
		},
		"swap": schema.Int64Attribute{
			Computed:            true,
			Description:         "Size of the flavor's swap space in MB. Zero when the flavor provides no swap.",
			MarkdownDescription: "Size of the flavor's swap space in MB. Zero when the flavor provides no swap.",
		},
		"ephemeral": schema.Int64Attribute{
			Computed:            true,
			Description:         "Size of the flavor's ephemeral disk in GB. Ephemeral storage is transient: its contents do not survive a rebuild or deletion of the instance. Zero when the flavor provides no ephemeral disk.",
			MarkdownDescription: "Size of the flavor's ephemeral disk in GB. Ephemeral storage is transient: its contents do not survive a rebuild or deletion of the instance. Zero when the flavor provides no ephemeral disk.",
		},
		"is_public": schema.BoolAttribute{
			Computed:            true,
			Description:         "Whether the flavor is publicly available to the project. Private flavors are only visible to the projects they have been explicitly shared with.",
			MarkdownDescription: "Whether the flavor is publicly available to the project. Private flavors are only visible to the projects they have been explicitly shared with.",
		},
		"description": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Free-form description of the flavor as reported by the backend. May be empty when no description is advertised for this flavor.",
			MarkdownDescription: "Free-form description of the flavor as reported by the backend. May be empty when no description is advertised for this flavor.",
		},
		"location": instanceCatalogLocationDataSourceAttribute(
			"Region (and, where applicable, availability zone) where this flavor is offered. The flavor catalog is per-region: a flavor returned for one region is not guaranteed to exist, or to carry the same characteristics, in another.",
		),
	}
}

func (d *cloudInstanceFlavorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := instanceFlavorDataSourceAttributes()
	attrs["service_name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Service name of the resource representing the id of the cloud project",
		MarkdownDescription: "Service name of the resource representing the id of the cloud project",
	}
	attrs["id"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "The OpenStack/Nova flavor ID to look up. Stable within a region and used to reference the flavor when creating an instance.",
		MarkdownDescription: "The OpenStack/Nova flavor ID to look up. Stable within a region and used to reference the flavor when creating an instance.",
	}

	resp.Schema = schema.Schema{
		Description:         "Use this data source to retrieve information about a flavor available in a public cloud project. This is read-only reference data: it describes a hardware sizing (vCPUs, RAM, disks) the project can create instances from, not an existing resource.",
		MarkdownDescription: "Use this data source to retrieve information about a flavor available in a public cloud project. This is read-only reference data: it describes a hardware sizing (vCPUs, RAM, disks) the project can create instances from, not an existing resource.",
		Attributes:          attrs,
	}
}

func (d *cloudInstanceFlavorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceFlavorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/instance/flavor/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceFlavorAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
