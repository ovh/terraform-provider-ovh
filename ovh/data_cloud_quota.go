package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudQuotaDataSource)(nil)

// NewCloudQuotaDataSource returns a new quota data source.
func NewCloudQuotaDataSource() datasource.DataSource {
	return &cloudQuotaDataSource{}
}

type cloudQuotaDataSource struct {
	config *Config
}

func (d *cloudQuotaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_quota"
}

func (d *cloudQuotaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudQuotaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	usageAttrs := map[string]schema.Attribute{
		"limit": schema.Int64Attribute{
			Computed:    true,
			Description: "Maximum authorized value for this quota.",
		},
		"used": schema.Int64Attribute{
			Computed:    true,
			Description: "Current usage reported by OpenStack. Null when not exposed by the service.",
		},
		"unit": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Unit of the limit and used values (e.g. count, GB, MB).",
		},
	}
	limitAttrs := map[string]schema.Attribute{
		"limit": schema.Int64Attribute{
			Computed:    true,
			Description: "Maximum authorized value for this limit.",
		},
		"unit": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Unit of the limit value.",
		},
	}

	resp.Schema = schema.Schema{
		Description: "Fetch read-only quota information (applied profile, available profiles and per-region usage) for a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Input
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "If set, restricts the per-region quota usage to this single region.",
			},

			// Envelope
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource identifier (the project id).",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Readiness of the resource.",
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Computed hash representing the current target specification value.",
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation date (RFC3339).",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date (RFC3339).",
			},

			// Flattened targetSpec
			"prevent_automatic_quota_upgrade": schema.BoolAttribute{
				Computed:    true,
				Description: "When true, automatic quota upgrades are disabled for this project.",
			},
			"regions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Target quota profile per region.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region where the profile applies.",
						},
						"profile": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Quota profile to apply in this region.",
						},
					},
				},
			},

			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current quota state of the project.",
				Attributes: map[string]schema.Attribute{
					"prevent_automatic_quota_upgrade": schema.BoolAttribute{
						Computed:    true,
						Description: "When true, automatic quota upgrades are disabled for this project.",
					},

					// --- Available profiles ---
					"available_profiles": schema.ListNestedAttribute{
						Computed:    true,
						Description: "List of available quota profiles with their caps.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Profile name.",
								},
								"compute": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"cores":     schema.Int64Attribute{Computed: true},
										"instances": schema.Int64Attribute{Computed: true},
										"memory":    schema.Int64Attribute{Computed: true},
									},
								},
								"volume": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"backup_size_total": schema.Int64Attribute{Computed: true},
										"backups":           schema.Int64Attribute{Computed: true},
										"size_total":        schema.Int64Attribute{Computed: true},
										"snapshots":         schema.Int64Attribute{Computed: true},
										"volumes":           schema.Int64Attribute{Computed: true},
									},
								},
								"network": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"floating_ips":         schema.Int64Attribute{Computed: true},
										"gateways":             schema.Int64Attribute{Computed: true},
										"networks":             schema.Int64Attribute{Computed: true},
										"security_group_rules": schema.Int64Attribute{Computed: true},
										"security_groups":      schema.Int64Attribute{Computed: true},
										"subnets":              schema.Int64Attribute{Computed: true},
									},
								},
								"loadbalancer": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"health_monitors": schema.Int64Attribute{Computed: true},
										"l7_policies":     schema.Int64Attribute{Computed: true},
										"l7_rules":        schema.Int64Attribute{Computed: true},
										"listeners":       schema.Int64Attribute{Computed: true},
										"loadbalancers":   schema.Int64Attribute{Computed: true},
										"members":         schema.Int64Attribute{Computed: true},
										"pools":           schema.Int64Attribute{Computed: true},
									},
								},
								"key_manager": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"containers": schema.Int64Attribute{Computed: true},
										"secrets":    schema.Int64Attribute{Computed: true},
									},
								},
								"share": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"backup_size_total": schema.Int64Attribute{Computed: true},
										"backups":           schema.Int64Attribute{Computed: true},
										"shares":            schema.Int64Attribute{Computed: true},
										"size_total":        schema.Int64Attribute{Computed: true},
										"snapshots":         schema.Int64Attribute{Computed: true},
									},
								},
								"keypair": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"keypairs": schema.Int64Attribute{Computed: true},
									},
								},
							},
						},
					},

					// --- Per-region usage ---
					"regions": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Per-region quota state reported by OpenStack.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Region name (e.g. GRA11).",
								},
								"profile": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Currently applied quota profile name in this region.",
								},
								"compute": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"cores":     schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"instances": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"memory":    schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"volume": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"backup_size_total": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"backups":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"per_volume_size":   schema.SingleNestedAttribute{Computed: true, Attributes: limitAttrs},
										"size_total":        schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"snapshots":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"volumes":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"network": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"floating_ips":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"gateways":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"networks":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"security_group_rules": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"security_groups":      schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"subnets":              schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"loadbalancer": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"health_monitors": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"l7_policies":     schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"l7_rules":        schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"listeners":       schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"loadbalancers":   schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"members":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"pools":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"key_manager": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"containers": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"secrets":    schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"share": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"backup_size_total":   schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"backups":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"per_share_size":      schema.SingleNestedAttribute{Computed: true, Attributes: limitAttrs},
										"share_networks":      schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"shares":              schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"size_total":          schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"snapshot_size_total": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"snapshots":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"keypair": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"keypairs": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudQuotaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudQuotaModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/quota"
	if data.Region.ValueString() != "" {
		endpoint += "?region=" + url.QueryEscape(data.Region.ValueString())
	}

	var responseData CloudQuotaAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
