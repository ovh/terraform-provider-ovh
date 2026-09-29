package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

const dedicatedCephPoolTimeout = 10 * time.Minute

var (
	_ resource.Resource                = (*dedicatedCephPoolResource)(nil)
	_ resource.ResourceWithConfigure   = (*dedicatedCephPoolResource)(nil)
	_ resource.ResourceWithImportState = (*dedicatedCephPoolResource)(nil)
)

func NewDedicatedCephPoolResource() resource.Resource {
	return &dedicatedCephPoolResource{}
}

type dedicatedCephPoolResource struct {
	config *Config
}

type DedicatedCephPoolModel struct {
	ID                ovhtypes.TfStringValue `tfsdk:"id"`
	ServiceName       ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name              ovhtypes.TfStringValue `tfsdk:"name"`
	PoolType          ovhtypes.TfStringValue `tfsdk:"pool_type"`
	Backup            ovhtypes.TfBoolValue   `tfsdk:"backup"`
	MinActiveReplicas ovhtypes.TfInt64Value  `tfsdk:"min_active_replicas"`
	ReplicaCount      ovhtypes.TfInt64Value  `tfsdk:"replica_count"`
}

func (m *DedicatedCephPoolModel) MergeWith(pool *DedicatedCephPool) {
	m.ID = ovhtypes.NewTfStringValue(m.ServiceName.ValueString() + "/" + pool.Name)
	m.Name = ovhtypes.NewTfStringValue(pool.Name)
	m.PoolType = ovhtypes.NewTfStringValue(pool.PoolType)
	m.Backup = ovhtypes.NewTfBoolValue(bool(pool.Backup))
	m.MinActiveReplicas = ovhtypes.NewTfInt64Value(pool.MinActiveReplicas)
	m.ReplicaCount = ovhtypes.NewTfInt64Value(pool.ReplicaCount)
}

func (r *dedicatedCephPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_ceph_pool"
}

func (r *dedicatedCephPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dedicatedCephPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a pool in a dedicated CEPH cluster (Cloud Disk Array).",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "The internal name of your dedicated CEPH",
				MarkdownDescription: "The internal name of your dedicated CEPH",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Name of the pool",
				MarkdownDescription: "Name of the pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Computed
			"pool_type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Type of the pool (ERASURE_CODED or REPLICATED)",
				MarkdownDescription: "Type of the pool (`ERASURE_CODED` or `REPLICATED`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backup": schema.BoolAttribute{
				CustomType:          ovhtypes.TfBoolType{},
				Computed:            true,
				Description:         "Whether the pool is backed up",
				MarkdownDescription: "Whether the pool is backed up",
			},
			"min_active_replicas": schema.Int64Attribute{
				CustomType:          ovhtypes.TfInt64Type{},
				Computed:            true,
				Description:         "Minimum number of active replicas",
				MarkdownDescription: "Minimum number of active replicas",
			},
			"replica_count": schema.Int64Attribute{
				CustomType:          ovhtypes.TfInt64Type{},
				Computed:            true,
				Description:         "Number of replicas",
				MarkdownDescription: "Number of replicas",
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Unique identifier for the resource, formatted as service_name/name",
				MarkdownDescription: "Unique identifier for the resource, formatted as `service_name/name`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func dedicatedCephPoolEndpoint(serviceName, poolName string) string {
	return "/dedicated/ceph/" + url.PathEscape(serviceName) + "/pool/" + url.PathEscape(poolName)
}

func (r *dedicatedCephPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DedicatedCephPoolModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	poolName := data.Name.ValueString()

	var taskId string
	endpoint := "/dedicated/ceph/" + url.PathEscape(serviceName) + "/pool"
	opts := &DedicatedCephPoolCreateOpts{PoolName: poolName}
	err := retryDedicatedCephChange(ctx, dedicatedCephPoolTimeout, func() error {
		return r.config.OVHClient.PostWithContext(ctx, endpoint, opts, &taskId)
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	if err := waitDedicatedCephTask(ctx, r.config.OVHClient, serviceName, taskId, dedicatedCephPoolTimeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for pool %s to be created", poolName), err.Error())
		return
	}

	var pool DedicatedCephPool
	endpoint = dedicatedCephPoolEndpoint(serviceName, poolName)
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &pool); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&pool)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dedicatedCephPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DedicatedCephPoolModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pool DedicatedCephPool
	endpoint := dedicatedCephPoolEndpoint(data.ServiceName.ValueString(), data.Name.ValueString())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &pool); err != nil {
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&pool)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dedicatedCephPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every configurable attribute requires a replacement, so there is nothing to update.
	resp.Diagnostics.AddError("Update not supported", "A dedicated CEPH pool cannot be updated, it must be replaced.")
}

func (r *dedicatedCephPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DedicatedCephPoolModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	poolName := data.Name.ValueString()
	endpoint := dedicatedCephPoolEndpoint(serviceName, poolName)

	// A pool can only be deleted in the 5 minutes window opened by allowDeletion, which opens it
	// right away and answers "Success" rather than a task. The window is opened again along with
	// the deletion when the cluster is locked, as it may close while waiting for the lock.
	var taskId string
	err := retryDedicatedCephChange(ctx, dedicatedCephPoolTimeout, func() error {
		if err := r.config.OVHClient.PutWithContext(ctx, endpoint+"/allowDeletion", nil, nil); err != nil {
			return fmt.Errorf("error calling Put %s/allowDeletion: %w", endpoint, err)
		}
		if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, &taskId); err != nil {
			return fmt.Errorf("error calling Delete %s: %w", endpoint, err)
		}
		return nil
	})
	if err != nil {
		// The client error is wrapped, so unwrap it to reach the API error.
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error deleting pool %s", poolName), err.Error())
		return
	}

	if err := waitDedicatedCephTask(ctx, r.config.OVHClient, serviceName, taskId, dedicatedCephPoolTimeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for pool %s to be deleted", poolName), err.Error())
	}
}

func (r *dedicatedCephPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.SplitN(req.ID, "/", 2)
	if len(splits) != 2 || splits[0] == "" || splits[1] == "" {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<pool_name>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
