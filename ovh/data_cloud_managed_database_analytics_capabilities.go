package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedDatabaseCapabilitiesDataSource)(nil)

func NewCloudManagedDatabaseCapabilitiesDataSource() datasource.DataSource {
	return &cloudManagedDatabaseCapabilitiesDataSource{category: managedDatabaseFamily}
}

func NewCloudManagedAnalyticsCapabilitiesDataSource() datasource.DataSource {
	return &cloudManagedDatabaseCapabilitiesDataSource{category: managedAnalyticsFamily}
}

type cloudManagedDatabaseCapabilitiesDataSource struct {
	managedDatabaseDataSourceConfigured
	category managedDatabaseCategory
}

// The publicCloud.database capabilities models. Their tfsdk tags map them onto the attributes.

type managedDatabaseCapabilityLifecycle struct {
	EndOfLifeAt *string `json:"endOfLifeAt" tfsdk:"end_of_life_at"`
	EndOfSaleAt *string `json:"endOfSaleAt" tfsdk:"end_of_sale_at"`
	StartedAt   string  `json:"startedAt" tfsdk:"started_at"`
	Status      string  `json:"status" tfsdk:"status"`
}

type managedDatabaseCapabilityEngineVersion struct {
	Default   bool                               `json:"default" tfsdk:"default"`
	Lifecycle managedDatabaseCapabilityLifecycle `json:"lifecycle" tfsdk:"lifecycle"`
	Name      string                             `json:"name" tfsdk:"name"`
	Tags      []string                           `json:"tags" tfsdk:"tags"`
}

type managedDatabaseCapabilityEngine struct {
	Category       string                                   `json:"category" tfsdk:"category"`
	DefaultVersion string                                   `json:"defaultVersion" tfsdk:"default_version"`
	Description    string                                   `json:"description" tfsdk:"description"`
	Lifecycle      managedDatabaseCapabilityLifecycle       `json:"lifecycle" tfsdk:"lifecycle"`
	Name           string                                   `json:"name" tfsdk:"name"`
	Order          int64                                    `json:"order" tfsdk:"order"`
	SslModes       []string                                 `json:"sslModes" tfsdk:"ssl_modes"`
	Storage        string                                   `json:"storage" tfsdk:"storage"`
	Tags           []string                                 `json:"tags" tfsdk:"tags"`
	Versions       []managedDatabaseCapabilityEngineVersion `json:"versions" tfsdk:"versions"`
}

type managedDatabaseCapabilityFlavorSpecifications struct {
	Memory managedDatabaseStorageSize `json:"memory" tfsdk:"memory"`
	Vcpus  int64                      `json:"vcpus" tfsdk:"vcpus"`
}

type managedDatabaseCapabilityFlavor struct {
	Generation     *string                                       `json:"generation" tfsdk:"generation"`
	Lifecycle      managedDatabaseCapabilityLifecycle            `json:"lifecycle" tfsdk:"lifecycle"`
	Name           string                                        `json:"name" tfsdk:"name"`
	Order          int64                                         `json:"order" tfsdk:"order"`
	Specifications managedDatabaseCapabilityFlavorSpecifications `json:"specifications" tfsdk:"specifications"`
	Tags           []string                                      `json:"tags" tfsdk:"tags"`
}

type managedDatabaseCapabilityPlan struct {
	BackupRetention *string                            `json:"backupRetention" tfsdk:"backup_retention"`
	Description     string                             `json:"description" tfsdk:"description"`
	Lifecycle       managedDatabaseCapabilityLifecycle `json:"lifecycle" tfsdk:"lifecycle"`
	Name            string                             `json:"name" tfsdk:"name"`
	Order           int64                              `json:"order" tfsdk:"order"`
	Tags            []string                           `json:"tags" tfsdk:"tags"`
}

type managedDatabaseCapabilityRegion struct {
	Lifecycle managedDatabaseCapabilityLifecycle `json:"lifecycle" tfsdk:"lifecycle"`
	Name      string                             `json:"name" tfsdk:"name"`
	Order     int64                              `json:"order" tfsdk:"order"`
	Tags      []string                           `json:"tags" tfsdk:"tags"`
	Type      string                             `json:"type" tfsdk:"type"`
}

type managedDatabaseCapabilityDisk struct {
	IopsPerGigabyte       *int64   `json:"iopsPerGigabyte" tfsdk:"iops_per_gigabyte"`
	MaximumIops           int64    `json:"maximumIops" tfsdk:"maximum_iops"`
	MaximumThroughput     *float64 `json:"maximumThroughput" tfsdk:"maximum_throughput"`
	Name                  string   `json:"name" tfsdk:"name"`
	ThroughputPerGigabyte *float64 `json:"throughputPerGigabyte" tfsdk:"throughput_per_gigabyte"`
}

type managedDatabaseCapabilityOption struct {
	Name string `json:"name" tfsdk:"name"`
	Type string `json:"type" tfsdk:"type"`
}

type cloudManagedDatabaseCapabilitiesModel struct {
	ServiceName types.String                      `tfsdk:"service_name"`
	Disks       []managedDatabaseCapabilityDisk   `tfsdk:"disks"`
	Engines     []managedDatabaseCapabilityEngine `tfsdk:"engines"`
	Flavors     []managedDatabaseCapabilityFlavor `tfsdk:"flavors"`
	Options     []managedDatabaseCapabilityOption `tfsdk:"options"`
	Plans       []managedDatabaseCapabilityPlan   `tfsdk:"plans"`
	Regions     []managedDatabaseCapabilityRegion `tfsdk:"regions"`
}

func (d *cloudManagedDatabaseCapabilitiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.category.typeName("_capabilities")
}

func (d *cloudManagedDatabaseCapabilitiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	computedInt64 := func(description string) schema.Int64Attribute {
		return schema.Int64Attribute{Computed: true, Description: description}
	}
	computedFloat64 := func(description string) schema.Float64Attribute {
		return schema.Float64Attribute{Computed: true, Description: description}
	}
	tags := schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "Display tags."}
	lifecycle := schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Lifecycle: started_at, end_of_sale_at, end_of_life_at and status (BETA, STABLE, DEPRECATED, END_OF_SALE, END_OF_LIFE or UNAVAILABLE).",
		Attributes: map[string]schema.Attribute{
			"started_at":     managedDatabaseComputedString("Release date."),
			"end_of_sale_at": managedDatabaseComputedString("End of sale date."),
			"end_of_life_at": managedDatabaseComputedString("End of life date."),
			"status":         managedDatabaseComputedString("Current status."),
		},
	}
	list := func(description string, attributes map[string]schema.Attribute) schema.ListNestedAttribute {
		return schema.ListNestedAttribute{Computed: true, Description: description, NestedObject: schema.NestedAttributeObject{Attributes: attributes}}
	}

	resp.Schema = schema.Schema{
		Description: "Use this data source to get what managed " + d.category.family + " services a public cloud project can create: engines and their versions, flavors, plans, regions, disk types and options, using the Cloud Databases API v2.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{Required: true, Description: "ID of the public cloud project."},
			"engines": list("Engines.", map[string]schema.Attribute{
				"name":            managedDatabaseComputedString("Name of the engine."),
				"description":     managedDatabaseComputedString("Description of the engine."),
				"category":        managedDatabaseComputedString("Category of the engine: DATABASE or ANALYTICS."),
				"default_version": managedDatabaseComputedString("Version used when the creation sets none."),
				"ssl_modes":       schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "TLS modes of the connections."},
				"storage":         managedDatabaseComputedString("Storage strategy: DISTRIBUTED, NOT_APPLICABLE or REPLICATED."),
				"order":           computedInt64("Display order."),
				"tags":            tags,
				"lifecycle":       lifecycle,
				"versions": list("Versions of the engine.", map[string]schema.Attribute{
					"name":      managedDatabaseComputedString("Name of the version."),
					"default":   schema.BoolAttribute{Computed: true, Description: "Whether it is the default version."},
					"tags":      tags,
					"lifecycle": lifecycle,
				}),
			}),
			"flavors": list("Flavors.", map[string]schema.Attribute{
				"name":       managedDatabaseComputedString("Name of the flavor."),
				"generation": managedDatabaseComputedString("Generation of the flavor: GEN1, GEN2 or GEN3."),
				"specifications": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Specifications of the flavor.",
					Attributes: map[string]schema.Attribute{
						"vcpus": computedInt64("Number of vCPU cores."),
						"memory": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "RAM size.",
							Attributes: map[string]schema.Attribute{
								"unit":  managedDatabaseComputedString("Unit of the size: GB or MB."),
								"value": computedInt64("Size, in the unit."),
							},
						},
					},
				},
				"order":     computedInt64("Display order."),
				"tags":      tags,
				"lifecycle": lifecycle,
			}),
			"plans": list("Plans.", map[string]schema.Attribute{
				"name":             managedDatabaseComputedString("Name of the plan."),
				"description":      managedDatabaseComputedString("Description of the plan."),
				"backup_retention": managedDatabaseComputedString("Retention of the automatic backups, as a duration."),
				"order":            computedInt64("Display order."),
				"tags":             tags,
				"lifecycle":        lifecycle,
			}),
			"regions": list("Regions.", map[string]schema.Attribute{
				"name":      managedDatabaseComputedString("Name of the region."),
				"type":      managedDatabaseComputedString("Type of the region: LOCAL-ZONE, REGION-1-AZ or REGION-3-AZ."),
				"order":     computedInt64("Display order."),
				"tags":      tags,
				"lifecycle": lifecycle,
			}),
			"disks": list("Disk types.", map[string]schema.Attribute{
				"name":                    managedDatabaseComputedString("Name of the disk type."),
				"maximum_iops":            computedInt64("Maximum IOPS per volume."),
				"iops_per_gigabyte":       computedInt64("IOPS per GB of disk. Null when the IOPS are fixed."),
				"maximum_throughput":      computedFloat64("Maximum throughput in MB/s per volume. Null without limit."),
				"throughput_per_gigabyte": computedFloat64("Throughput in MB/s per GB of disk. Null when it doesn't scale."),
			}),
			"options": list("Options.", map[string]schema.Attribute{
				"name": managedDatabaseComputedString("Name of the option."),
				"type": managedDatabaseComputedString("Type of the option: BOOLEAN, DOUBLE, DURATION, LONG or STRING."),
			}),
		},
	}
}

func (d *cloudManagedDatabaseCapabilitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudManagedDatabaseCapabilitiesModel
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service_name"), &data.ServiceName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	base := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/" + d.category.segment + "/capabilities/"
	data.Disks = getManagedDatabaseCapabilities[managedDatabaseCapabilityDisk](ctx, d.config, base+"disks", &resp.Diagnostics)
	data.Engines = getManagedDatabaseCapabilities[managedDatabaseCapabilityEngine](ctx, d.config, base+"engines", &resp.Diagnostics)
	data.Flavors = getManagedDatabaseCapabilities[managedDatabaseCapabilityFlavor](ctx, d.config, base+"flavors", &resp.Diagnostics)
	data.Options = getManagedDatabaseCapabilities[managedDatabaseCapabilityOption](ctx, d.config, base+"options", &resp.Diagnostics)
	data.Plans = getManagedDatabaseCapabilities[managedDatabaseCapabilityPlan](ctx, d.config, base+"plans", &resp.Diagnostics)
	data.Regions = getManagedDatabaseCapabilities[managedDatabaseCapabilityRegion](ctx, d.config, base+"regions", &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getManagedDatabaseCapabilities[T any](ctx context.Context, config *Config, endpoint string, diags *diag.Diagnostics) []T {
	items, err := helpers.GetAllPagesV2[T](ctx, config.OVHClient, endpoint)
	if err != nil {
		diags.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
	}
	return items
}
