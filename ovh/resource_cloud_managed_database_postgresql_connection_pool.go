package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

var (
	_ resource.Resource                = (*cloudManagedDatabasePostgresqlConnectionPoolResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudManagedDatabasePostgresqlConnectionPoolResource)(nil)
	_ resource.ResourceWithImportState = (*cloudManagedDatabasePostgresqlConnectionPoolResource)(nil)
)

func NewCloudManagedDatabasePostgresqlConnectionPoolResource() resource.Resource {
	return &cloudManagedDatabasePostgresqlConnectionPoolResource{managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"postgresql"}, suffix: "_postgresql_connection_pool"}}
}

type cloudManagedDatabasePostgresqlConnectionPoolResource struct {
	managedDatabaseObjectResource
}

func (r *cloudManagedDatabasePostgresqlConnectionPoolResource) endpoint(data *CloudManagedDatabaseConnectionPoolModel) string {
	return r.serviceEndpoint(data.ServiceName, data.Engine, data.ServiceId) + "/connectionPool"
}

func (r *cloudManagedDatabasePostgresqlConnectionPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := r.objectAttributes("connection pool")
	attributes["name"] = schema.StringAttribute{
		Required:      true,
		Description:   "Name of the connection pool. Changing it recreates the connection pool.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
	attributes["database_id"] = schema.StringAttribute{
		Required:    true,
		Description: "ID of the database the pool connects to.",
	}
	attributes["mode"] = schema.StringAttribute{
		Required:    true,
		Description: "PgBouncer pooling mode: SESSION, STATEMENT or TRANSACTION.",
		Validators:  []validator.String{stringvalidator.OneOf("SESSION", "STATEMENT", "TRANSACTION")},
	}
	attributes["size"] = schema.Int64Attribute{
		Required:    true,
		Description: "Maximum number of server connections in the pool.",
		Validators:  []validator.Int64{int64validator.AtLeast(1)},
	}
	attributes["user_id"] = schema.StringAttribute{
		Optional:    true,
		Description: "ID of the user the pool connects with, which can't be the administrator of the service. If omitted, the pool uses a system user.",
	}
	attributes["port"] = schema.Int64Attribute{
		Computed:      true,
		Description:   "Port the connection pool listens on.",
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	}
	attributes["ssl_mode"] = schema.StringAttribute{
		Computed:      true,
		Description:   "TLS mode of the connection pool.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["uri"] = schema.StringAttribute{
		Computed:    true,
		Description: "Connection URI of the connection pool.",
	}
	attributes["checksum"] = schema.StringAttribute{
		Computed:    true,
		Description: "Hash of the connection pool, computed by the API.",
	}
	attributes["created_at"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Creation date of the connection pool.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["updated_at"] = schema.StringAttribute{
		Computed:    true,
		Description: "Last update date of the connection pool.",
	}

	resp.Schema = schema.Schema{
		Description: "Creates a PgBouncer connection pool in a PostgreSQL service, using the Cloud Databases API v2.",
		Attributes:  attributes,
	}
}

func (r *cloudManagedDatabasePostgresqlConnectionPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.importState(ctx, req, resp, "id")
}

func (r *cloudManagedDatabasePostgresqlConnectionPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudManagedDatabaseConnectionPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data)
	var pool managedDatabaseConnectionPool
	if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseObjectTimeout, func() error {
		return retry.RetryContext(ctx, managedDatabaseObjectTimeout, func() *retry.RetryError {
			err := r.config.OVHClient.PostWithContext(ctx, endpoint, data.ToCreate(), &pool)
			// The forward path refuses a pool whose database Aiven doesn't know yet, right after its creation.
			var apiErr *ovh.APIError
			if errors.As(err, &apiErr) && apiErr.Code == http.StatusForbidden && strings.Contains(apiErr.Details["message"], "does not exist") {
				return retry.RetryableError(err)
			}
			if err != nil {
				return retry.NonRetryableError(err)
			}
			return nil
		})
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}
	data.MergeWith(&pool)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if pool.Uri != "" || resp.Diagnostics.HasError() {
		return
	}

	// The native path gives a pool its URI and its port once PgBouncer serves it.
	endpoint += "/" + url.PathEscape(pool.Id)
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING", managedDatabaseAPIUnavailable},
		Target:  []string{"READY"},
		Refresh: func() (any, string, error) {
			var res managedDatabaseConnectionPool
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, &res)
			switch {
			case isManagedDatabaseAPIUnavailable(err):
				return res, managedDatabaseAPIUnavailable, nil
			case err != nil:
				return nil, "", err
			case res.Uri == "":
				return res, "PENDING", nil
			}
			return res, "READY", nil
		},
		Timeout:    managedDatabaseObjectTimeout,
		MinTimeout: managedDatabaseServicePollMinTimeout,
	}
	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the connection pool to be served", err.Error())
		return
	}
	pool = res.(managedDatabaseConnectionPool)
	data.MergeWith(&pool)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedDatabasePostgresqlConnectionPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudManagedDatabaseConnectionPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	var pool managedDatabaseConnectionPool
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &pool); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	data.MergeWith(&pool)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudManagedDatabasePostgresqlConnectionPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CloudManagedDatabaseConnectionPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&state) + "/" + url.PathEscape(state.Id.ValueString())
	var pool managedDatabaseConnectionPool
	if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseObjectTimeout, func() error {
		return r.config.OVHClient.PutWithContext(ctx, endpoint, plan.ToPut(), &pool)
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return
	}
	plan.MergeWith(&pool)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cloudManagedDatabasePostgresqlConnectionPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudManagedDatabaseConnectionPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(&data) + "/" + url.PathEscape(data.Id.ValueString())
	if err := deleteManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}
	// On the native path, the pool stays until its workflow removes it, and its database and its user
	// can't be deleted before.
	if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, true, managedDatabaseObjectTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for the connection pool to be deleted", err.Error())
	}
}
