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

var _ datasource.DataSourceWithConfigure = (*cloudKeymanagerSecretsDataSource)(nil)

func NewCloudKeymanagerSecretsDataSource() datasource.DataSource {
	return &cloudKeymanagerSecretsDataSource{}
}

type cloudKeymanagerSecretsDataSource struct {
	config *Config
}

func (d *cloudKeymanagerSecretsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_keymanager_secrets"
}

func (d *cloudKeymanagerSecretsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudKeymanagerSecretsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all secrets in the Barbican Key Manager service for a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"secrets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of secrets",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Secret ID",
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current resource state",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the secret",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the secret",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Secret readiness status",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region of the secret",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Name of the secret",
						},
						"secret_type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Type of the secret",
						},
						"current_state": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Current state of the secret",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"secret_type": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"algorithm": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"bit_length": schema.Int64Attribute{
									Computed: true,
								},
								"mode": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"payload_content_type": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"expiration": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"secret_ref": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"status": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"region": schema.StringAttribute{
									CustomType: ovhtypes.TfStringType{},
									Computed:   true,
								},
								"metadata": schema.MapAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudKeymanagerSecretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudKeymanagerSecretsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/secret"

	var apiSecrets []CloudKeymanagerSecretAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiSecrets); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	var secretValues []attr.Value
	for _, s := range apiSecrets {
		region := ""
		name := ""
		secretType := ""
		if s.TargetSpec != nil {
			name = s.TargetSpec.Name
			secretType = s.TargetSpec.SecretType
			if s.TargetSpec.Location != nil {
				region = s.TargetSpec.Location.Region
			}
		}

		var currentStateObj types.Object
		if s.CurrentState != nil {
			currentStateObj = buildKeymanagerSecretCurrentStateObject(ctx, s.CurrentState)
		} else {
			currentStateObj = types.ObjectNull(KeymanagerSecretCurrentStateAttrTypes())
		}

		itemObj, _ := types.ObjectValue(
			KeymanagerSecretListItemAttrTypes(),
			map[string]attr.Value{
				"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(s.Id)},
				"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(s.Checksum)},
				"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(s.CreatedAt)},
				"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(s.UpdatedAt)},
				"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(s.ResourceStatus)},
				"region":          ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
				"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(name)},
				"secret_type":     ovhtypes.TfStringValue{StringValue: types.StringValue(secretType)},
				"current_state":   currentStateObj,
			},
		)
		secretValues = append(secretValues, itemObj)
	}

	secretsList, _ := types.ListValue(
		types.ObjectType{AttrTypes: KeymanagerSecretListItemAttrTypes()},
		secretValues,
	)
	data.Secrets = secretsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
