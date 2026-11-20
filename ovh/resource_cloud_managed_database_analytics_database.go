package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

var (
	_ resource.Resource                = (*cloudManagedDatabaseDatabaseResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudManagedDatabaseDatabaseResource)(nil)
	_ resource.ResourceWithImportState = (*cloudManagedDatabaseDatabaseResource)(nil)
)

func NewCloudManagedDatabaseDatabaseResource() resource.Resource {
	return &cloudManagedDatabaseDatabaseResource{managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"mysql", "postgresql"}, suffix: "_database"}}
}

func NewCloudManagedAnalyticsClickhouseDatabaseResource() resource.Resource {
	return &cloudManagedDatabaseDatabaseResource{managedDatabaseObjectResource{category: managedAnalyticsFamily, engines: []string{"clickhouse"}, suffix: "_clickhouse_database"}}
}

type cloudManagedDatabaseDatabaseResource struct {
	managedDatabaseObjectResource
}

func (r *cloudManagedDatabaseDatabaseResource) endpoint(data *CloudManagedDatabaseDatabaseModel) string {
	return r.serviceEndpoint(data.ServiceName, data.Engine, data.ServiceId) + "/database"
}

func (r *cloudManagedDatabaseDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := r.objectAttributes("database")
	attributes["name"] = schema.StringAttribute{
		Required:      true,
		Description:   "Name of the database. Changing it recreates the database.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
	attributes["default"] = schema.BoolAttribute{
		Computed:      true,
		Description:   "Whether the database is the default one, created with the service.",
		PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
	}
	attributes["checksum"] = schema.StringAttribute{
		Computed:    true,
		Description: "Hash of the database, computed by the API.",
	}
	attributes["created_at"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Creation date of the database.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["updated_at"] = schema.StringAttribute{
		Computed:    true,
		Description: "Last update date of the database.",
	}

	resp.Schema = schema.Schema{
		Description: r.describe("a logical database"),
		Attributes:  attributes,
	}
}

func (r *cloudManagedDatabaseDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.importState(ctx, req, resp, "id")
}

func (r *cloudManagedDatabaseDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudManagedDatabaseDatabaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data)
	var database managedDatabaseDatabase
	if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseObjectTimeout, func() error {
		return r.config.OVHClient.PostWithContext(ctx, endpoint, managedDatabaseDatabaseCreation{Name: data.Name.ValueString()}, &database)
	}); err != nil {
		// A POST can fail after making the database: its answer was lost, or the forward path didn't see
		// the database in Aiven at once, and lists it once it resyncs the service, within a minute.
		existing := r.find(ctx, endpoint, data.Name.ValueString())
		var apiErr *ovh.APIError
		switch {
		case existing == nil:
			resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
			return
		case errors.As(err, &apiErr) && apiErr.Class == "Client::Conflict::DatabaseAlready":
			id := []string{data.ServiceName.ValueString()}
			if !r.named() {
				id = append(id, data.Engine.ValueString())
			}
			id = append(id, data.ServiceId.ValueString(), existing.Id)
			resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint),
				fmt.Sprintf("%s\n\nTo manage the existing database, import it with the ID %s.", err, strings.Join(id, "/")))
			return
		}
		database = *existing
	}
	data.MergeWith(&database)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// find returns the database of the service with the name, nil when there is none or the list fails.
func (r *cloudManagedDatabaseDatabaseResource) find(ctx context.Context, endpoint, name string) *managedDatabaseDatabase {
	databases, err := helpers.GetAllPagesV2[managedDatabaseDatabase](ctx, r.config.OVHClient, endpoint)
	if err != nil {
		return nil
	}
	for _, database := range databases {
		if database.Name == name {
			return &database
		}
	}
	return nil
}

func (r *cloudManagedDatabaseDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudManagedDatabaseDatabaseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	var database managedDatabaseDatabase
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &database); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	data.MergeWith(&database)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is never called: every argument recreates the database.
func (r *cloudManagedDatabaseDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unexpected update", "Every argument of a database recreates it. Please report this issue to the provider developers.")
}

func (r *cloudManagedDatabaseDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudManagedDatabaseDatabaseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Default.ValueBool() {
		resp.Diagnostics.AddWarning("Default database not deleted",
			"The API can't delete the default database of the service, created with it: "+data.Name.ValueString()+" is only removed from the state.")
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	if err := deleteManagedDatabaseObject(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}
	// On the native path, the database stays until its workflow drops it.
	if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, true, managedDatabaseObjectTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for the database to be deleted", err.Error())
	}
}
