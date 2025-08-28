package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*okmsSecretResource)(nil)

func NewOkmsSecretResource() resource.Resource {
	return &okmsSecretResource{}
}

type okmsSecretResource struct {
	config *Config
}

func (r *okmsSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_okms_secret"
}

func (d *okmsSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.config = config
}

func (d *okmsSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OkmsSecretResourceSchema(ctx)
}

func (r *okmsSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData OkmsSecretModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Cas.IsNull() && !data.Cas.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"CAS Ignored On Create",
			"The 'cas' attribute is only used on update operations and ignored during creation.",
		)
	}

	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/secret"
	createPayload := buildSecretPayload(&data, true)
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)
	if responseData.Version.Data.IsNull() || responseData.Version.Data.IsUnknown() {
		responseData.Version = data.Version
	}

	populateVersionComputedFields(r, &responseData, data.OkmsId.ValueString(), data.Path.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *okmsSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData OkmsSecretModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/secret/" + url.PathEscape(data.Path.ValueString()) + ""

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData OkmsSecretModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update resource
	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/secret/" + url.PathEscape(data.Path.ValueString()) + ""

	// Avoid creating a new secret version when only metadata (or other fields) changed and
	// the version data itself is unchanged. The API creates a new version whenever a
	// version payload is sent, even if the content is identical. We therefore don't include the
	// version field in the update payload when the user-specified data matches the
	// prior state value.
	planForPayload := planData // shallow copy
	if !planData.Version.Data.IsNull() && !planData.Version.Data.IsUnknown() &&
		!data.Version.Data.IsNull() && !data.Version.Data.IsUnknown() &&
		planData.Version.Data.ValueString() == data.Version.Data.ValueString() {
		// Mark version data null in the payload model so buildSecretPayload skips it.
		planForPayload.Version.Data = ovhtypes.NewTfStringNull()
	}
	updatePayload := buildSecretPayload(&planForPayload, false)
	// cas (check-and-set) must be passed as query parameter
	casQuery := ""
	if !planData.Cas.IsNull() && !planData.Cas.IsUnknown() {
		casQuery = "?cas=" + url.QueryEscape(fmt.Sprintf("%d", planData.Cas.ValueInt64()))
	}
	if err := r.config.OVHClient.Put(endpoint+casQuery, updatePayload, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint+casQuery),
			err.Error(),
		)
		return
	}

	// Read updated resource
	endpoint = "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/secret/" + url.PathEscape(data.Path.ValueString()) + ""
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)

	populateVersionComputedFields(r, &responseData, data.OkmsId.ValueString(), data.Path.ValueString())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *okmsSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OkmsSecretModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/secret/" + url.PathEscape(data.Path.ValueString()) + ""
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}

// buildSecretPayload constructs the payload for create/update.
// On create include path; on update path is immutable so omitted.
func buildSecretPayload(m *OkmsSecretModel, isCreate bool) map[string]any {
	payload := map[string]any{}
	if isCreate && !m.Path.IsNull() && !m.Path.IsUnknown() {
		payload["path"] = m.Path.ValueString()
	}
	if !m.Version.Data.IsNull() && !m.Version.Data.IsUnknown() {
		if vp := buildVersionData(m.Version.Data.ValueString()); vp != nil {
			payload["version"] = vp
		}
	}
	if meta := buildMetadataPayload(&m.Metadata); meta != nil {
		payload["metadata"] = meta
	}
	return payload
}

// buildVersionData attempts to JSON decode the user provided string; if structured returns structured form.
func buildVersionData(raw string) map[string]any {
	versionPayload := map[string]any{}
	var decoded any
	if err := json.Unmarshal([]byte(raw), &decoded); err == nil {
		switch decoded.(type) {
		case map[string]any, []any:
			versionPayload["data"] = decoded
		default:
			versionPayload["data"] = raw
		}
	} else {
		versionPayload["data"] = raw
	}
	return versionPayload
}

// buildMetadataPayload extracts settable metadata fields.
func buildMetadataPayload(meta *MetadataValue) map[string]any {
	mp := map[string]any{}
	if !meta.CustomMetadata.IsNull() && !meta.CustomMetadata.IsUnknown() {
		mp["customMetadata"] = meta.CustomMetadata
	}
	if !meta.MaxVersions.IsNull() && !meta.MaxVersions.IsUnknown() {
		mp["maxVersions"] = meta.MaxVersions
	}
	if !meta.DeactivateVersionAfter.IsNull() && !meta.DeactivateVersionAfter.IsUnknown() {
		mp["deactivateVersionAfter"] = meta.DeactivateVersionAfter
	}
	if !meta.CasRequired.IsNull() && !meta.CasRequired.IsUnknown() {
		mp["casRequired"] = meta.CasRequired
	}
	if len(mp) == 0 {
		return nil
	}
	return mp
}

// populateVersionComputedFields fills secret version attributes
func populateVersionComputedFields(r *okmsSecretResource, model *OkmsSecretModel, okmsId, path string) {
	// If currentVersion unknown or zero, nothing to enrich
	if model.Metadata.CurrentVersion.IsNull() || model.Metadata.CurrentVersion.IsUnknown() || model.Metadata.CurrentVersion.ValueInt64() == 0 {
		return
	}

	current := model.Metadata.CurrentVersion.ValueInt64()

	// First try the efficient direct version endpoint
	versionEndpoint := "/v2/okms/resource/" + url.PathEscape(okmsId) + "/secret/" + url.PathEscape(path) + "/version/" + fmt.Sprintf("%d", current)
	var ver struct {
		Id            *int64  `json:"id"`
		CreatedAt     *string `json:"createdAt"`
		State         *string `json:"state"`
		DeactivatedAt *string `json:"deactivatedAt"`
	}
	if err := r.config.OVHClient.Get(versionEndpoint, &ver); err != nil || ver.Id == nil {
		// Best-effort enrichment; silently skip on error
		return
	}
	// Populate from direct call
	model.Version.Id = ovhtypes.NewTfInt64Value(*ver.Id)
	if ver.CreatedAt != nil {
		model.Version.CreatedAt = ovhtypes.NewTfStringValue(*ver.CreatedAt)
	}
	if ver.State != nil {
		model.Version.State = ovhtypes.NewTfStringValue(*ver.State)
	}
	if ver.DeactivatedAt != nil {
		model.Version.DeactivatedAt = ovhtypes.NewTfStringValue(*ver.DeactivatedAt)
	} else {
		model.Version.DeactivatedAt = ovhtypes.NewTfStringNull()
	}
}
