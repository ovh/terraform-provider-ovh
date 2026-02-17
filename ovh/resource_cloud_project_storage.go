package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	"github.com/peterhellberg/duration"
)

var _ resource.ResourceWithConfigure = (*cloudProjectStorageResource)(nil)

func NewCloudProjectStorageResource() resource.Resource {
	return &cloudProjectStorageResource{}
}

type cloudProjectStorageResource struct {
	config *Config
}

func (r *cloudProjectStorageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage"
}

func (d *cloudProjectStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *cloudProjectStorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectRegionStorageResourceSchema(ctx)
}

func (r *cloudProjectStorageResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// 1. Skip if the entire resource is being created or destroyed
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// 2. Extract the Configuration and State for the object_lock attribute
	var config, state, plan ObjectLockValue

	// We check the Config to see if the user removed the block
	req.Config.GetAttribute(ctx, path.Root("object_lock"), &config)
	req.State.GetAttribute(ctx, path.Root("object_lock"), &state)
	req.Plan.GetAttribute(ctx, path.Root("object_lock"), &plan)

	// 3. Logic: If it's in State (Enabled) but NOT in Config (Removed)
	if config.IsNull() && !state.IsNull() {
		if state.Status.ValueString() == "enabled" {
			// A. Force the attribute in the Plan to be Null (this breaks the "No Changes" loop)
			resp.Plan.SetAttribute(ctx, path.Root("object_lock"), types.ObjectNull(plan.AttributeTypes(ctx)))

			// B. Explicitly trigger the resource replacement
			resp.RequiresReplace = append(resp.RequiresReplace, path.Root("object_lock"))

			// C. Add a helpful message to the CLI
			resp.Diagnostics.AddWarning(
				"Removing Object Lock - Bucket Replacement Required",
				"The object_lock configuration has been removed from the Terraform configuration. "+
					"Object Lock cannot be disabled once enabled on a bucket. "+
					"Terraform will DESTROY the existing bucket and CREATE a new one without Object Lock. "+
					"WARNING: All objects in the bucket will be permanently deleted.",
			)
		}
	}

	// 4. Logic: If object_lock exists in both config and state, check for status changes
	if !config.IsNull() && !state.IsNull() {
		if !config.Status.IsNull() && !state.Status.IsNull() {
			configStatus := config.Status.ValueString()
			stateStatus := state.Status.ValueString()

			// Status change from enabled to disabled requires replacement
			if stateStatus == "enabled" && configStatus == "disabled" {
				resp.RequiresReplace = append(resp.RequiresReplace, path.Root("object_lock"))

				resp.Diagnostics.AddWarning(
					"Disabling Object Lock - Bucket Replacement Required",
					"Object Lock status is being changed from 'enabled' to 'disabled'. "+
						"Object Lock cannot be disabled once enabled on a bucket. "+
						"Terraform will DESTROY the existing bucket and CREATE a new one. "+
						"WARNING: All objects in the bucket will be permanently deleted.",
				)
			}

			// Status change from disabled to enabled requires replacement
			if stateStatus == "disabled" && configStatus == "enabled" {
				resp.RequiresReplace = append(resp.RequiresReplace, path.Root("object_lock"))

				resp.Diagnostics.AddWarning(
					"Enabling Object Lock - Bucket Replacement Required",
					"Object Lock status is being changed from 'disabled' to 'enabled'. "+
						"Object Lock must be enabled at bucket creation. "+
						"Terraform will DESTROY the existing bucket and CREATE a new one. "+
						"WARNING: All objects in the bucket will be permanently deleted.",
				)
			}
		}
	}
}

func (r *cloudProjectStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<region_name>/<storage_name>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), splits[2])...)
}

func (r *cloudProjectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData CloudProjectRegionStorageModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle service_name: use provided value or fall back to environment variable
	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() || data.ServiceName.ValueString() == "" {
		envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
		if envServiceName == "" {
			resp.Diagnostics.AddError(
				"Missing service_name",
				"The service_name attribute is required. Please provide it in the resource configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
			)
			return
		}
		data.ServiceName = ovhtypes.NewTfStringValue(envServiceName)
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)
	r.fixISO8601Diff(&data, &responseData)

	// Set the ID as composite key: service_name/region_name/name
	compositeID := fmt.Sprintf("%s/%s/%s",
		data.ServiceName.ValueString(),
		data.RegionName.ValueString(),
		data.Name.ValueString())
	responseData.ID = ovhtypes.NewTfStringValue(compositeID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectRegionStorageModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Add filters if defined
	queryParams := url.Values{}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		queryParams.Add("limit", strconv.FormatInt(data.Limit.ValueInt64(), 10))
	}
	if !data.Marker.IsNull() && !data.Marker.IsUnknown() {
		queryParams.Add("marker", data.Marker.ValueString())
	}
	if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
		queryParams.Add("prefix", data.Prefix.ValueString())
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage/" + url.PathEscape(data.Name.ValueString()) + "?" + queryParams.Encode()
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		// Handle 404 errors by removing the resource from state
		// This allows Terraform to recreate resources deleted outside of Terraform
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Resource not found, removing from state: %s", endpoint))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)
	r.fixISO8601Diff(&data, &responseData)

	if data.HideObjects.ValueBool() {
		responseData.Objects = ovhtypes.NewListNestedObjectValueOfNull[ObjectsValue](ctx)
	}

	// Set the ID as composite key: service_name/region_name/name
	compositeID := fmt.Sprintf("%s/%s/%s",
		data.ServiceName.ValueString(),
		data.RegionName.ValueString(),
		data.Name.ValueString())
	responseData.ID = ovhtypes.NewTfStringValue(compositeID)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData CloudProjectRegionStorageModel

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
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage/" + url.PathEscape(data.Name.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Add filters if defined
	queryParams := url.Values{}
	if !planData.Limit.IsNull() && !planData.Limit.IsUnknown() {
		queryParams.Add("limit", strconv.FormatInt(planData.Limit.ValueInt64(), 10))
	}
	if !planData.Marker.IsNull() && !planData.Marker.IsUnknown() {
		queryParams.Add("marker", planData.Marker.ValueString())
	}
	if !planData.Prefix.IsNull() && !planData.Prefix.IsUnknown() {
		queryParams.Add("prefix", planData.Prefix.ValueString())
	}

	// Read updated resource
	endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage/" + url.PathEscape(data.Name.ValueString()) + "?" + queryParams.Encode()
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)
	r.fixISO8601Diff(&planData, &responseData)

	if planData.HideObjects.ValueBool() {
		responseData.Objects = ovhtypes.NewListNestedObjectValueOfNull[ObjectsValue](ctx)
	}

	// Set the ID as composite key: service_name/region_name/name
	compositeID := fmt.Sprintf("%s/%s/%s",
		data.ServiceName.ValueString(),
		data.RegionName.ValueString(),
		data.Name.ValueString())
	responseData.ID = ovhtypes.NewTfStringValue(compositeID)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectRegionStorageModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.deleteBucket(ctx, data.ServiceName.ValueString(), data.RegionName.ValueString(), data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting bucket: "+data.Name.ValueString(),
			err.Error(),
		)
		return
	}

	// If replicas exist, delete them manually after the main bucket deletion if option is set
	for _, rule := range data.Replication.Rules.Elements() {
		destination := rule.(ReplicationRulesValue).Destination

		// Empty region means that replica is already deleted
		if destination.Region.ValueString() == "" {
			continue
		}

		if destination.RemoveOnMainBucketDeletion.ValueBool() {
			tflog.Info(ctx, fmt.Sprintf("removing replica bucket %s", destination.Name.ValueString()))
			if err := r.deleteBucket(ctx, data.ServiceName.ValueString(), destination.Region.ValueString(), destination.Name.ValueString()); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error removing replica %s", destination.Name.ValueString()),
					err.Error(),
				)
			}
		}
	}
}

func (r *cloudProjectStorageResource) deleteBucket(ctx context.Context, serviceName, regionName, storageName string) error {
	endpoint := "/cloud/project/" + url.PathEscape(serviceName) +
		"/region/" + url.PathEscape(regionName) +
		"/storage/" + url.PathEscape(storageName) +
		"/object?withVersions=true"
	bulkDeleteEndpoint := "/cloud/project/" + url.PathEscape(serviceName) +
		"/region/" + url.PathEscape(regionName) +
		"/storage/" + url.PathEscape(storageName) +
		"/bulkDeleteObjects"

	// Try to empty bucket to be able to destroy it completely.
	// This operation can fail if objects are locked. In this case, they must be unlocked manually
	// before retrying the operation.
	// A maximum of 1000 objects are returned at each GET call, so we have to iterate to empty bucket.
	for {
		var (
			objects     []CloudProjectStorageObjectsValue
			idsToDelete []map[string]string
		)

		if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &objects); err != nil {
			return fmt.Errorf("error calling GET %s: %w", endpoint, err)
		}

		if len(objects) == 0 {
			break
		}

		for _, obj := range objects {
			idsToDelete = append(idsToDelete, map[string]string{
				"key":       obj.Key.ValueString(),
				"versionId": obj.VersionId.ValueString(),
			})
		}

		tflog.Info(ctx, fmt.Sprintf("removing objects %s", idsToDelete))
		if err := r.config.OVHClient.PostWithContext(ctx, bulkDeleteEndpoint, map[string]any{
			"objects": idsToDelete,
		}, nil); err != nil {
			return fmt.Errorf("error calling POST %s: %w", bulkDeleteEndpoint, err)
		}
	}

	endpoint = "/cloud/project/" + url.PathEscape(serviceName) +
		"/region/" + url.PathEscape(regionName) +
		"/storage/" + url.PathEscape(storageName)

	// Delete bucket itself
	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, nil); err != nil {
		if ovhErr, ok := err.(*ovh.APIError); ok && ovhErr.Code == http.StatusNotFound {
			// If bucket was already deleted, ignore the error
			return nil
		}

		return fmt.Errorf("error calling DELETE %s: %w", endpoint, err)
	}

	return nil
}

func (r *cloudProjectStorageResource) fixISO8601Diff(expected, actual *CloudProjectRegionStorageModel) {
	if expected.ObjectLock.IsNull() || expected.ObjectLock.IsUnknown() ||
		actual.ObjectLock.IsNull() || actual.ObjectLock.IsUnknown() {
		return
	}

	if expected.ObjectLock.Rule.IsNull() || expected.ObjectLock.Rule.IsUnknown() ||
		actual.ObjectLock.Rule.IsNull() || actual.ObjectLock.Rule.IsUnknown() {
		return
	}

	expectedPeriod := expected.ObjectLock.Rule.Period.ValueString()
	actualPeriod := actual.ObjectLock.Rule.Period.ValueString()

	if expectedPeriod == actualPeriod {
		return
	}

	// Parse both periods using the duration library
	expectedDur, err1 := duration.Parse(expectedPeriod)
	actualDur, err2 := duration.Parse(actualPeriod)

	if err1 != nil || err2 != nil {
		// If parsing fails, skip normalization
		return
	}

	// Compare the underlying time.Duration values
	if expectedDur == actualDur {
		// If semantically equal, update actual to match expected to avoid Terraform diff
		actual.ObjectLock.Rule.Period = expected.ObjectLock.Rule.Period
	}
}
