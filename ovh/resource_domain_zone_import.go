package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*domainZoneImportResource)(nil)

func NewDomainZoneImportResource() resource.Resource {
	return &domainZoneImportResource{}
}

type domainZoneImportResource struct {
	config *Config
}

func (r *domainZoneImportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_zone_import"
}

func (d *domainZoneImportResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *domainZoneImportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DomainZoneImportResourceSchema(ctx)
}

func (r *domainZoneImportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		data   DomainZoneImportModel
		task   DomainTask
		export string
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Import zone file
	endpoint := "/domain/zone/" + url.PathEscape(data.ZoneName.ValueString()) + "/import"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Wait for import task completion
	if err := waitDNSTask(ctx, r.config, data.ZoneName.ValueString(), task.TaskID); err != nil {
		resp.Diagnostics.AddError("Error waiting for task completion", err.Error())
		return
	}

	// Export zone file
	endpoint = "/domain/zone/" + url.PathEscape(data.ZoneName.ValueString()) + "/export"
	if err := r.config.OVHClient.Get(endpoint, &export); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	data.ExportedContent = types.NewTfStringValue(export)

	data.ID = data.ZoneName

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainZoneImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		data         DomainZoneImportModel
		responseData string
	)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/domain/zone/" + url.PathEscape(data.ZoneName.ValueString()) + "/export"
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Force diff if records have been modified
	if data.ExportedContent.ValueString() != responseData {
		data.ZoneFile = types.NewTfStringValue(responseData)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainZoneImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("update should never happen", "this code should be unreachable")
}

func (r *domainZoneImportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainZoneImportModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/domain/zone/" + url.PathEscape(data.ZoneName.ValueString()) + "/reset"
	if err := r.config.OVHClient.Post(endpoint, json.RawMessage(`{"minimized":false}`), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
	}
}

func waitDNSTask(ctx context.Context, config *Config, zoneName string, taskID int) error {
	endpoint := fmt.Sprintf("/domain/zone/%s/task/%d", url.PathEscape(zoneName), taskID)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"doing", "todo"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {
			var task DomainTask
			if err := config.OVHClient.Get(endpoint, &task); err != nil {
				log.Printf("[ERROR] couldn't fetch task %d: error: %v", taskID, err)
				return nil, "error", err
			}
			return task.Status, task.Status, nil
		},
		Timeout:    time.Hour,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
