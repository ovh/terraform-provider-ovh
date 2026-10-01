package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceDataSource)(nil)

func NewCloudInstanceDataSource() datasource.DataSource {
	return &cloudInstanceDataSource{}
}

type cloudInstanceDataSource struct {
	config *Config
}

func (d *cloudInstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance"
}

func (d *cloudInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about an instance in a public cloud project.",
		Attributes:  instanceDataSourceAttributes(ctx),
	}
}

func (d *cloudInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData, data.priorSpec())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// instanceDataSourceAttributes builds the datasource schema: same shape as the
// resource but service_name+id are the only inputs and everything else is Computed.
func instanceDataSourceAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"service_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: instanceDescServiceName, MarkdownDescription: instanceDescServiceName},
		"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: instanceDescId, MarkdownDescription: instanceDescId},

		"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescLocationRegion, MarkdownDescription: instanceDescLocationRegion},
		"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescLocationAZ, MarkdownDescription: instanceDescLocationAZ},
		"ssh_key_name":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsSSHKeyName, MarkdownDescription: instanceDescCsSSHKeyName},
		"group_id":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescGroupId, MarkdownDescription: instanceDescGroupId},
		"name":              schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescDsName, MarkdownDescription: instanceDescDsName},
		"flavor_id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescFlavorId, MarkdownDescription: instanceDescFlavorId},
		"image_id":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescDsImageId, MarkdownDescription: instanceDescDsImageId},
		"power_state":       schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescDsPowerState, MarkdownDescription: instanceDescDsPowerStateMd},
		"networks": schema.ListNestedAttribute{Computed: true, Description: instanceDescDsNetworks, MarkdownDescription: instanceDescDsNetworks, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"network_id":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescNetworkRefNetworkId, MarkdownDescription: instanceDescNetworkRefNetworkId},
			"subnet_id":             schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescNetworkRefSubnetId, MarkdownDescription: instanceDescNetworkRefSubnetIdMd},
			"ip":                    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescNetworkRefIp, MarkdownDescription: instanceDescNetworkRefIpMd},
			"auto_assign_public_ip": schema.BoolAttribute{Computed: true, Description: instanceDescNetworkRefAutoAssign, MarkdownDescription: instanceDescNetworkRefAutoAssignMd},
		}}},
		"volume_ids":         schema.ListAttribute{CustomType: ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx), Computed: true, Description: instanceDescVolumeIds, MarkdownDescription: instanceDescVolumeIds},
		"security_group_ids": schema.ListAttribute{CustomType: ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx), Computed: true, Description: instanceDescDsSecurityGroupIds, MarkdownDescription: instanceDescDsSecurityGroupIds},
		"shares": schema.ListNestedAttribute{Computed: true, Description: instanceDescDsShares, MarkdownDescription: instanceDescDsShares, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescShareRefId, MarkdownDescription: instanceDescShareRefId},
			"access_level": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescShareRefAccessLevel, MarkdownDescription: instanceDescShareRefAccessLevelMd},
		}}},
		"checksum":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescChecksum, MarkdownDescription: instanceDescChecksum},
		"created_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCreatedAt, MarkdownDescription: instanceDescCreatedAt},
		"updated_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescUpdatedAt, MarkdownDescription: instanceDescUpdatedAt},
		"resource_status": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescResourceStatus, MarkdownDescription: instanceDescResourceStatusMd},
		"current_state":   schema.SingleNestedAttribute{Computed: true, Description: instanceDescCurrentState, MarkdownDescription: instanceDescCurrentState, Attributes: instanceCurrentStateDataSourceSchemaAttributes()},
	}
}

// instanceCurrentStateDataSourceSchemaAttributes mirrors
// instanceCurrentStateSchemaAttributes but returns datasource-schema attributes
// (the resource and datasource schema packages define distinct Attribute types).
// Shared by both the singular and plural instance data sources.
func instanceCurrentStateDataSourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsName, MarkdownDescription: instanceDescCsName},
		"power_state":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsPowerState, MarkdownDescription: instanceDescCsPowerState},
		"locked":       schema.BoolAttribute{Computed: true, Description: instanceDescCsLocked, MarkdownDescription: instanceDescCsLocked},
		"ssh_key_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsSSHKeyName, MarkdownDescription: instanceDescCsSSHKeyName},
		"host_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsHostId, MarkdownDescription: instanceDescCsHostId},
		"project_id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsProjectId, MarkdownDescription: instanceDescCsProjectId},
		"user_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsUserId, MarkdownDescription: instanceDescCsUserId},
		"flavor": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsFlavor, MarkdownDescription: instanceDescCsFlavor, Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescFlavorId, MarkdownDescription: instanceDescFlavorId},
			"name":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescFlavorName, MarkdownDescription: instanceDescFlavorName},
			"vcpus":     schema.Int64Attribute{Computed: true, Description: instanceDescFlavorVcpus, MarkdownDescription: instanceDescFlavorVcpus},
			"ram":       schema.Int64Attribute{Computed: true, Description: instanceDescFlavorRam, MarkdownDescription: instanceDescFlavorRam},
			"disk":      schema.Int64Attribute{Computed: true, Description: instanceDescFlavorDisk, MarkdownDescription: instanceDescFlavorDisk},
			"swap":      schema.Int64Attribute{Computed: true, Description: instanceDescFlavorSwap, MarkdownDescription: instanceDescFlavorSwap},
			"ephemeral": schema.Int64Attribute{Computed: true, Description: instanceDescFlavorEphemeral, MarkdownDescription: instanceDescFlavorEphemeral},
		}},
		"image": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsImage, MarkdownDescription: instanceDescCsImage, Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescImageId, MarkdownDescription: instanceDescImageId},
			"name":       schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescImageName, MarkdownDescription: instanceDescImageName},
			"size":       schema.Int64Attribute{Computed: true, Description: instanceDescImageSize, MarkdownDescription: instanceDescImageSize},
			"status":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescImageStatus, MarkdownDescription: instanceDescImageStatus},
			"deprecated": schema.BoolAttribute{Computed: true, Description: instanceDescImageDeprecated, MarkdownDescription: instanceDescImageDeprecated},
		}},
		"location": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsLocation, MarkdownDescription: instanceDescCsLocation, Attributes: map[string]schema.Attribute{
			"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescLocationRegion, MarkdownDescription: instanceDescLocationRegion},
			"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescLocationAZ, MarkdownDescription: instanceDescLocationAZ},
		}},
		"networks": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsNetworks, MarkdownDescription: instanceDescCsNetworks, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsNetworkId, MarkdownDescription: instanceDescCsNetworkId},
			"subnet_id":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsNetworkSubnetId, MarkdownDescription: instanceDescCsNetworkSubnetId},
			"gateway_id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsNetworkGatewayId, MarkdownDescription: instanceDescCsNetworkGatewayId},
			"addresses": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsAddresses, MarkdownDescription: instanceDescCsAddressesMd, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"ip":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescAddressIp, MarkdownDescription: instanceDescAddressIp},
				"mac":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescAddressMac, MarkdownDescription: instanceDescAddressMac},
				"type":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescAddressType, MarkdownDescription: instanceDescAddressTypeMd},
				"version": schema.Int64Attribute{Computed: true, Description: instanceDescAddressVersion, MarkdownDescription: instanceDescAddressVersion},
			}}},
		}}},
		"volumes": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsVolumes, MarkdownDescription: instanceDescCsVolumes, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescVolumeId, MarkdownDescription: instanceDescVolumeId},
			"name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescVolumeName, MarkdownDescription: instanceDescVolumeName},
			"size": schema.Int64Attribute{Computed: true, Description: instanceDescVolumeSize, MarkdownDescription: instanceDescVolumeSize},
		}}},
		"shares": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsShares, MarkdownDescription: instanceDescCsShares, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareId, MarkdownDescription: instanceDescCsShareId},
			"access_level": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareAccessLevel, MarkdownDescription: instanceDescCsShareAccessLevelMd},
			"access_to":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareAccessTo, MarkdownDescription: instanceDescCsShareAccessTo},
			"state":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareState, MarkdownDescription: instanceDescCsShareStateMd},
		}}},
		"security_groups": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsSecurityGroups, MarkdownDescription: instanceDescCsSecurityGroups, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescSecurityGroupId, MarkdownDescription: instanceDescSecurityGroupId},
		}}},
		"group": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsGroup, MarkdownDescription: instanceDescCsGroup, Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescGroupId, MarkdownDescription: instanceDescGroupId},
		}},
	}
}
