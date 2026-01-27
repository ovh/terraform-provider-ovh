package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
)

var _ resource.ResourceWithConfigure = (*MeIdentityUserTokenResource)(nil)
var _ resource.ResourceWithImportState = (*MeIdentityUserTokenResource)(nil)

func NewMeIdentityUserTokenResource() resource.Resource {
	return &MeIdentityUserTokenResource{}
}

type MeIdentityUserTokenResource struct {
	config *Config
}

func (r *MeIdentityUserTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_me_identity_user_token"
}

func (r *MeIdentityUserTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

func (r *MeIdentityUserTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a token for a specific identity user.",
		Attributes: map[string]schema.Attribute{
			"user_login": schema.StringAttribute{
				Required:    true,
				Description: "User's login suffix",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Token name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "Token description",
			},
			"expires_at": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Token expiration date",
			},
			"expires_in": schema.Int64Attribute{
				Optional:    true,
				Description: "Token validity duration in seconds",
				// Using Int64 for integer values in framework
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The token value",
			},
			"creation": schema.StringAttribute{
				Computed:    true,
				Description: "Creation date of this token",
			},
			"last_used": schema.StringAttribute{
				Computed:    true,
				Description: "Last use of this token",
			},
		},
	}
}

func (r *MeIdentityUserTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'user_login/name', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_login"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func (r *MeIdentityUserTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MeIdentityUserTokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := &MeIdentityUserTokenCreateOpts{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		opts.ExpiresAt = plan.ExpiresAt.ValueString()
	}
	if !plan.ExpiresIn.IsNull() && !plan.ExpiresIn.IsUnknown() {
		opts.ExpiresIn = int(plan.ExpiresIn.ValueInt64())
	}

	endpoint := fmt.Sprintf("/me/identity/user/%s/token", url.PathEscape(plan.UserLogin.ValueString()))
	var apiResp MeIdentityUserTokenResponse
	err := r.config.OVHClient.PostWithContext(ctx, endpoint, opts, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating identity user token %s", plan.Name.ValueString()),
			err.Error(),
		)
		return
	}

	plan.Token = types.StringValue(apiResp.Token)
	plan.Creation = types.StringValue(apiResp.Creation)

	// Handle ExpiresAt format preservation
	expiresAtVal := normalizeTime(apiResp.ExpiresAt)
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		planTime, err1 := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		apiTime, err2 := time.Parse(time.RFC3339, apiResp.ExpiresAt)
		if err1 == nil && err2 == nil && planTime.Equal(apiTime) {
			expiresAtVal = plan.ExpiresAt.ValueString()
		}
	}
	plan.ExpiresAt = types.StringValue(expiresAtVal)

	if apiResp.LastUsed != nil {
		plan.LastUsed = types.StringValue(*apiResp.LastUsed)
	} else {
		plan.LastUsed = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MeIdentityUserTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MeIdentityUserTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/me/identity/user/%s/token/%s",
		url.PathEscape(state.UserLogin.ValueString()),
		url.PathEscape(state.Name.ValueString()))

	var apiResp MeIdentityUserTokenResponse
	err := r.config.OVHClient.GetWithContext(ctx, endpoint, &apiResp)
	if err != nil {
		if ovhErr, ok := err.(*ovh.APIError); ok && ovhErr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading identity user token %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}

	state.Name = types.StringValue(apiResp.Name)
	state.Description = types.StringValue(apiResp.Description)
	state.Creation = types.StringValue(apiResp.Creation)

	// Handle ExpiresAt format preservation
	expiresAtVal := normalizeTime(apiResp.ExpiresAt)
	if !state.ExpiresAt.IsNull() && !state.ExpiresAt.IsUnknown() {
		stateTime, err1 := time.Parse(time.RFC3339, state.ExpiresAt.ValueString())
		apiTime, err2 := time.Parse(time.RFC3339, apiResp.ExpiresAt)
		if err1 == nil && err2 == nil && stateTime.Equal(apiTime) {
			expiresAtVal = state.ExpiresAt.ValueString()
		}
	}
	state.ExpiresAt = types.StringValue(expiresAtVal)

	if apiResp.LastUsed != nil {
		state.LastUsed = types.StringValue(*apiResp.LastUsed)
	} else {
		state.LastUsed = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MeIdentityUserTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MeIdentityUserTokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state MeIdentityUserTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := &MeIdentityUserTokenUpdateOpts{}

	if !plan.Description.IsNull() {
		opts.Description = plan.Description.ValueString()
	}
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		opts.ExpiresAt = plan.ExpiresAt.ValueString()
	}
	if !plan.ExpiresIn.IsNull() && !plan.ExpiresIn.IsUnknown() {
		opts.ExpiresIn = int(plan.ExpiresIn.ValueInt64())
	}

	endpoint := fmt.Sprintf("/me/identity/user/%s/token/%s",
		url.PathEscape(plan.UserLogin.ValueString()),
		url.PathEscape(plan.Name.ValueString()))

	var apiResp MeIdentityUserTokenResponse
	err := r.config.OVHClient.PutWithContext(ctx, endpoint, opts, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating identity user token %s", plan.Name.ValueString()),
			err.Error(),
		)
		return
	}

	// Update state. Keep existing token secret because API doesn't return it.
	if apiResp.Token != "" {
		plan.Token = types.StringValue(apiResp.Token)
	} else {
		plan.Token = state.Token
	}

	plan.Creation = types.StringValue(apiResp.Creation)

	// Handle ExpiresAt format preservation
	expiresAtVal := normalizeTime(apiResp.ExpiresAt)
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		planTime, err1 := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		apiTime, err2 := time.Parse(time.RFC3339, apiResp.ExpiresAt)
		if err1 == nil && err2 == nil && planTime.Equal(apiTime) {
			expiresAtVal = plan.ExpiresAt.ValueString()
		}
	}
	plan.ExpiresAt = types.StringValue(expiresAtVal)

	if apiResp.LastUsed != nil {
		plan.LastUsed = types.StringValue(*apiResp.LastUsed)
	} else {
		plan.LastUsed = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MeIdentityUserTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MeIdentityUserTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/me/identity/user/%s/token/%s",
		url.PathEscape(state.UserLogin.ValueString()),
		url.PathEscape(state.Name.ValueString()))

	err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	if err != nil {
		// Ignore 404
		if ovhErr, ok := err.(*ovh.APIError); ok && ovhErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting identity user token %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}
}

type MeIdentityUserTokenModel struct {
	UserLogin   types.String `tfsdk:"user_login"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ExpiresAt   types.String `tfsdk:"expires_at"`
	ExpiresIn   types.Int64  `tfsdk:"expires_in"`
	Token       types.String `tfsdk:"token"`
	Creation    types.String `tfsdk:"creation"`
	LastUsed    types.String `tfsdk:"last_used"`
}

func normalizeTime(t string) string {
	if t == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return t
	}
	return parsed.UTC().Format(time.RFC3339)
}
