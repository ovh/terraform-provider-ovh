package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectQuotaDataSource)(nil)

// NewCloudProjectQuotaDataSource returns a new quota data source.
func NewCloudProjectQuotaDataSource() datasource.DataSource {
	return &cloudProjectQuotaDataSource{}
}

type cloudProjectQuotaDataSource struct {
	config *Config
}

func (d *cloudProjectQuotaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_quota"
}

func (d *cloudProjectQuotaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectQuotaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"target_spec": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Desired quota specification for the project.",
				Attributes: map[string]schema.Attribute{
					"profile": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Name of the quota profile applied to the project.",
					},
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current quota state of the project.",
				Attributes: map[string]schema.Attribute{
					"profile": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Name of the currently applied quota profile.",
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
										"instances":            schema.Int64Attribute{Computed: true},
										"cores":                schema.Int64Attribute{Computed: true},
										"ram":                  schema.Int64Attribute{Computed: true},
										"security_groups":      schema.Int64Attribute{Computed: true},
										"security_group_rules": schema.Int64Attribute{Computed: true},
										"server_groups":        schema.Int64Attribute{Computed: true},
										"server_group_members": schema.Int64Attribute{Computed: true},
									},
								},
								"block_storage": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"volumes":          schema.Int64Attribute{Computed: true},
										"gigabytes":        schema.Int64Attribute{Computed: true},
										"snapshots":        schema.Int64Attribute{Computed: true},
										"backups":          schema.Int64Attribute{Computed: true},
										"backup_gigabytes": schema.Int64Attribute{Computed: true},
									},
								},
								"network": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"networks":             schema.Int64Attribute{Computed: true},
										"subnets":              schema.Int64Attribute{Computed: true},
										"floating_ips":         schema.Int64Attribute{Computed: true},
										"gateways":             schema.Int64Attribute{Computed: true},
										"security_groups":      schema.Int64Attribute{Computed: true},
										"security_group_rules": schema.Int64Attribute{Computed: true},
									},
								},
								"loadbalancer": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"loadbalancers":  schema.Int64Attribute{Computed: true},
										"listeners":      schema.Int64Attribute{Computed: true},
										"pools":          schema.Int64Attribute{Computed: true},
										"members":        schema.Int64Attribute{Computed: true},
										"healthmonitors": schema.Int64Attribute{Computed: true},
										"l7_policies":    schema.Int64Attribute{Computed: true},
										"l7_rules":       schema.Int64Attribute{Computed: true},
									},
								},
								"key_manager": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"secrets":    schema.Int64Attribute{Computed: true},
										"containers": schema.Int64Attribute{Computed: true},
									},
								},
								"share": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"shares":           schema.Int64Attribute{Computed: true},
										"gigabytes":        schema.Int64Attribute{Computed: true},
										"snapshots":        schema.Int64Attribute{Computed: true},
										"backups":          schema.Int64Attribute{Computed: true},
										"backup_gigabytes": schema.Int64Attribute{Computed: true},
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
						Description: "Per-region quota usage reported by OpenStack.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Region name (e.g. GRA7).",
								},
								"compute": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"instances": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"cores":     schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"memory":    schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"volume": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"volumes":          schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"gigabytes":        schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"snapshots":        schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"backups":          schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"backup_gigabytes": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"per_volume_size":  schema.SingleNestedAttribute{Computed: true, Attributes: limitAttrs},
									},
								},
								"network": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"networks":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"subnets":              schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"floating_ips":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"gateways":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"security_groups":      schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"security_group_rules": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"loadbalancer": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"loadbalancers":  schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"listeners":      schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"pools":          schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"members":        schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"healthmonitors": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"l7_policies":    schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"l7_rules":       schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"key_manager": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"secrets":    schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"containers": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
									},
								},
								"share": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"shares":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"size_total":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"snapshots":          schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"snapshot_gigabytes": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"backups":            schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"backup_gigabytes":   schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"share_networks":     schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
										"per_share_size":     schema.SingleNestedAttribute{Computed: true, Attributes: limitAttrs},
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

func (d *cloudProjectQuotaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectQuotaModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/quota"
	if !data.Region.IsNull() && !data.Region.IsUnknown() && data.Region.ValueString() != "" {
		endpoint += "?region=" + url.QueryEscape(data.Region.ValueString())
	}

	var responseData CloudProjectQuotaAPIResponse
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
