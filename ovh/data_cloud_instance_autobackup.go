package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceAutobackupDataSource)(nil)

func NewCloudInstanceAutobackupDataSource() datasource.DataSource {
	return &cloudInstanceAutobackupDataSource{}
}

type cloudInstanceAutobackupDataSource struct {
	config *Config
}

func (d *cloudInstanceAutobackupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_autobackup"
}

func (d *cloudInstanceAutobackupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceAutobackupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads an autobackup (Mistral crontrigger) for a Public Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Public Cloud project ID",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Autobackup (crontrigger) ID",
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Crontrigger name",
			},
			"image_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Backup image name prefix",
			},
			"cron": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Cron schedule pattern",
			},
			"rotation": schema.Int64Attribute{
				CustomType:  ovhtypes.TfInt64Type{},
				Computed:    true,
				Description: "Number of backup versions to keep",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "OpenStack region",
			},
			"instance_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Source instance ID",
			},
			"distant": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Cross-region backup configuration",
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Target region for cross-region backup",
					},
					"image_name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Image name for the cross-region backup copy",
					},
				},
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource checksum for concurrency control",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource status (READY, CREATING, DELETING, ERROR)",
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation timestamp (RFC 3339)",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update timestamp (RFC 3339)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state from the infrastructure",
				Attributes:  cloudInstanceAutobackupDataSourceCurrentStateAttributes(),
			},
		},
	}
}

func cloudInstanceAutobackupDataSourceCurrentStateAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Crontrigger name",
		},
		"image_name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Backup image name prefix",
		},
		"cron": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Cron schedule pattern",
		},
		"rotation": schema.Int64Attribute{
			CustomType:  ovhtypes.TfInt64Type{},
			Computed:    true,
			Description: "Number of backup versions to keep",
		},
		"region": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "OpenStack region",
		},
		"instance_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Source instance ID",
		},
		"workflow_name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "OpenStack Mistral workflow name",
		},
		"next_execution_time": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Next scheduled execution time",
		},
		"distant": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Cross-region backup configuration",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Target region for cross-region backup",
				},
				"image_name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Image name for the cross-region backup copy",
				},
			},
		},
	}
}

func (d *cloudInstanceAutobackupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceAutobackupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/compute/autobackup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceAutobackupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
