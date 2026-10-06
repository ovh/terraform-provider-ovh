package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// dedicatedCephUserTimeout bounds each wait of a change: for the tasks in progress on the
// cluster, and for the task of the change. Tasks may last long, such as the deletion of a pool
// which took more than 10 minutes in testing, and a change may first wait for the tasks of others.
const dedicatedCephUserTimeout = 30 * time.Minute

var (
	_ resource.Resource                   = (*dedicatedCephUserResource)(nil)
	_ resource.ResourceWithConfigure      = (*dedicatedCephUserResource)(nil)
	_ resource.ResourceWithImportState    = (*dedicatedCephUserResource)(nil)
	_ resource.ResourceWithValidateConfig = (*dedicatedCephUserResource)(nil)
)

func NewDedicatedCephUserResource() resource.Resource {
	return &dedicatedCephUserResource{}
}

type dedicatedCephUserResource struct {
	config *Config
}

type DedicatedCephUserModel struct {
	ID              ovhtypes.TfStringValue `tfsdk:"id"`
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name            ovhtypes.TfStringValue `tfsdk:"name"`
	PoolPermissions types.List             `tfsdk:"pool_permissions"`
	Key             ovhtypes.TfStringValue `tfsdk:"key"`
	MonCaps         ovhtypes.TfStringValue `tfsdk:"mon_caps"`
	OsdCaps         ovhtypes.TfStringValue `tfsdk:"osd_caps"`
	MdsCaps         ovhtypes.TfStringValue `tfsdk:"mds_caps"`
}

type DedicatedCephUserPoolPermissionModel struct {
	PoolName   types.String `tfsdk:"pool_name"`
	Read       types.Bool   `tfsdk:"read"`
	Write      types.Bool   `tfsdk:"write"`
	Execute    types.Bool   `tfsdk:"execute"`
	ClassRead  types.Bool   `tfsdk:"class_read"`
	ClassWrite types.Bool   `tfsdk:"class_write"`
}

// GrantsNothing tells whether the configured permission grants nothing on its pool, which the API
// refuses. An unknown flag may grant something, and an omitted one defaults to false.
func (m DedicatedCephUserPoolPermissionModel) GrantsNothing() bool {
	for _, flag := range []types.Bool{m.Read, m.Write, m.Execute, m.ClassRead, m.ClassWrite} {
		if flag.IsUnknown() || flag.ValueBool() {
			return false
		}
	}
	return true
}

func DedicatedCephUserPoolPermissionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"pool_name":   types.StringType,
		"read":        types.BoolType,
		"write":       types.BoolType,
		"execute":     types.BoolType,
		"class_read":  types.BoolType,
		"class_write": types.BoolType,
	}
}

func (m *DedicatedCephUserModel) MergeWith(user *DedicatedCephUser) {
	m.ID = ovhtypes.NewTfStringValue(m.ServiceName.ValueString() + "/" + user.Name)
	m.Name = ovhtypes.NewTfStringValue(user.Name)
	m.Key = ovhtypes.NewTfStringValue(user.Key)
	m.MonCaps = ovhtypes.NewTfStringValue(user.MonCaps)
	m.OsdCaps = ovhtypes.NewTfStringValue(user.OsdCaps)
	if user.MdsCaps != nil {
		m.MdsCaps = ovhtypes.NewTfStringValue(*user.MdsCaps)
	} else {
		m.MdsCaps = ovhtypes.NewTfStringNull()
	}
}

// CreatedState returns the state of a user whose creation the API accepted, before its attributes
// can be read and its pool permissions set: they are null rather than unknown, which a state cannot
// hold.
func (m DedicatedCephUserModel) CreatedState() DedicatedCephUserModel {
	m.ID = ovhtypes.NewTfStringValue(m.ServiceName.ValueString() + "/" + m.Name.ValueString())
	m.PoolPermissions = types.ListNull(types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()})
	m.Key = ovhtypes.NewTfStringNull()
	m.MonCaps = ovhtypes.NewTfStringNull()
	m.OsdCaps = ovhtypes.NewTfStringNull()
	m.MdsCaps = ovhtypes.NewTfStringNull()
	return m
}

// ToPoolPermissions converts the pool_permissions attribute to its API representation.
func (m *DedicatedCephUserModel) ToPoolPermissions(ctx context.Context) ([]DedicatedCephUserPoolPermission, diag.Diagnostics) {
	var permissions []DedicatedCephUserPoolPermission
	if m.PoolPermissions.IsNull() || m.PoolPermissions.IsUnknown() {
		return permissions, nil
	}

	var models []DedicatedCephUserPoolPermissionModel
	diags := m.PoolPermissions.ElementsAs(ctx, &models, false)
	for _, model := range models {
		permissions = append(permissions, DedicatedCephUserPoolPermission{
			PoolName:   model.PoolName.ValueString(),
			Read:       model.Read.ValueBool(),
			Write:      model.Write.ValueBool(),
			Execute:    model.Execute.ValueBool(),
			ClassRead:  model.ClassRead.ValueBool(),
			ClassWrite: model.ClassWrite.ValueBool(),
		})
	}

	return permissions, diags
}

// MergeWithPoolPermissions sets the pool_permissions attribute from the API permissions. The
// permissions granting nothing are left out, and the attribute stays null when it was not set and
// the user has no permission, so that omitting it in the configuration shows no drift. The
// permissions the model already holds keep their order, the others coming after them in the order
// of the API.
func (m *DedicatedCephUserModel) MergeWithPoolPermissions(ctx context.Context, permissions []DedicatedCephUserPoolPermission) diag.Diagnostics {
	objectType := types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()}

	known, diags := m.ToPoolPermissions(ctx)
	if diags.HasError() {
		return diags
	}
	knownPools := make(map[string]int, len(known))
	for i, permission := range known {
		knownPools[permission.PoolName] = i
	}

	sorted := make([]DedicatedCephUserPoolPermission, len(permissions))
	copy(sorted, permissions)
	sort.SliceStable(sorted, func(i, j int) bool {
		posI, knownI := knownPools[sorted[i].PoolName]
		posJ, knownJ := knownPools[sorted[j].PoolName]
		if knownI && knownJ {
			return posI < posJ
		}
		return knownI && !knownJ
	})

	// Not nil, as a nil slice would make a null list instead of the configured empty one.
	models := []DedicatedCephUserPoolPermissionModel{}
	for _, permission := range sorted {
		if permission.IsEmpty() {
			continue
		}
		models = append(models, DedicatedCephUserPoolPermissionModel{
			PoolName:   types.StringValue(permission.PoolName),
			Read:       types.BoolValue(permission.Read),
			Write:      types.BoolValue(permission.Write),
			Execute:    types.BoolValue(permission.Execute),
			ClassRead:  types.BoolValue(permission.ClassRead),
			ClassWrite: types.BoolValue(permission.ClassWrite),
		})
	}

	if len(models) == 0 && (m.PoolPermissions.IsNull() || m.PoolPermissions.IsUnknown()) {
		m.PoolPermissions = types.ListNull(objectType)
		return diags
	}

	var listDiags diag.Diagnostics
	m.PoolPermissions, listDiags = types.ListValueFrom(ctx, objectType, models)
	diags.Append(listDiags...)
	return diags
}

func (r *dedicatedCephUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_ceph_user"
}

func (r *dedicatedCephUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// dedicatedCephUserPermissionDefault defaults an omitted permission flag to false.
type dedicatedCephUserPermissionDefault struct{}

var _ defaults.Bool = dedicatedCephUserPermissionDefault{}

func (d dedicatedCephUserPermissionDefault) Description(ctx context.Context) string {
	return "value defaults to false"
}

func (d dedicatedCephUserPermissionDefault) MarkdownDescription(ctx context.Context) string {
	return "value defaults to `false`"
}

func (d dedicatedCephUserPermissionDefault) DefaultBool(ctx context.Context, req defaults.BoolRequest, resp *defaults.BoolResponse) {
	resp.PlanValue = types.BoolValue(false)
}

func dedicatedCephUserPermissionFlag(description string) schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             dedicatedCephUserPermissionDefault{},
		Description:         description,
		MarkdownDescription: description,
	}
}

func (r *dedicatedCephUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a user in a dedicated CEPH cluster (Cloud Disk Array) and manages its pool permissions.",
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
				Description:         "Name of the user",
				MarkdownDescription: "Name of the user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pool_permissions": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "Permissions of the user on the pools of the cluster. The permissions on the pools not listed are cleared",
				MarkdownDescription: "Permissions of the user on the pools of the cluster. The permissions on the pools not listed are cleared",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"pool_name": schema.StringAttribute{
							Required:            true,
							Description:         "Name of the pool",
							MarkdownDescription: "Name of the pool",
						},
						"read":        dedicatedCephUserPermissionFlag("Read permission"),
						"write":       dedicatedCephUserPermissionFlag("Write permission"),
						"execute":     dedicatedCephUserPermissionFlag("Execute permission"),
						"class_read":  dedicatedCephUserPermissionFlag("Class read permission"),
						"class_write": dedicatedCephUserPermissionFlag("Class write permission"),
					},
				},
			},

			// Computed
			"key": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Sensitive:           true,
				Description:         "Key of the user to connect to the cluster",
				MarkdownDescription: "Key of the user to connect to the cluster",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mon_caps": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Capabilities of the user on the MON daemons",
				MarkdownDescription: "Capabilities of the user on the MON daemons",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"osd_caps": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Capabilities of the user on the OSD daemons, derived from its pool permissions",
				MarkdownDescription: "Capabilities of the user on the OSD daemons, derived from its pool permissions",
			},
			"mds_caps": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Capabilities of the user on the MDS daemons",
				MarkdownDescription: "Capabilities of the user on the MDS daemons",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

// validateDedicatedCephUserPoolPermissions refuses the permissions granting nothing, which the API
// refuses with "Empty association not allowed", the permissions on a pool being cleared by removing
// it from the list. It also refuses the pools listed more than once, as the permissions are set per
// pool. A permission with an unknown flag is accepted, as it may grant something, and a pool with
// an unknown name is checked once its name is known.
func validateDedicatedCephUserPoolPermissions(ctx context.Context, permissions types.List) diag.Diagnostics {
	var diags diag.Diagnostics
	if permissions.IsNull() || permissions.IsUnknown() {
		return diags
	}

	var models []DedicatedCephUserPoolPermissionModel
	diags.Append(permissions.ElementsAs(ctx, &models, false)...)

	firstIndex := make(map[string]int, len(models))
	for i, model := range models {
		if model.GrantsNothing() {
			diags.AddAttributeError(
				path.Root("pool_permissions").AtListIndex(i),
				"Pool permission granting nothing",
				"A pool permission must grant at least one of read, write, execute, class_read or class_write. "+
					"To clear the permissions of the user on a pool, remove the pool from pool_permissions.",
			)
		}

		if model.PoolName.IsNull() || model.PoolName.IsUnknown() {
			continue
		}
		poolName := model.PoolName.ValueString()
		if first, found := firstIndex[poolName]; found {
			diags.AddAttributeError(
				path.Root("pool_permissions").AtListIndex(i).AtName("pool_name"),
				"Duplicate pool permission",
				fmt.Sprintf("The pool %q is listed more than once in pool_permissions, at index %d and %d. "+
					"Merge its permissions into a single entry.", poolName, first, i),
			)
			continue
		}
		firstIndex[poolName] = i
	}
	return diags
}

// ValidateConfig refuses the permissions granting nothing as early as possible. Their flags may be
// unknown until applying, so Create and Update check them again once known.
func (r *dedicatedCephUserResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var permissions types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("pool_permissions"), &permissions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateDedicatedCephUserPoolPermissions(ctx, permissions)...)
}

func dedicatedCephUserEndpoint(serviceName, userName string) string {
	return "/dedicated/ceph/" + url.PathEscape(serviceName) + "/user/" + url.PathEscape(userName)
}

// setPoolPermissions brings the permissions of the user to the wanted ones. The POST replaces all
// the permissions of the user at once; when no permission is wanted, the current ones are cleared
// pool by pool instead.
func (r *dedicatedCephUserResource) setPoolPermissions(ctx context.Context, serviceName, userName string, current, wanted []DedicatedCephUserPoolPermission) error {
	endpoint := dedicatedCephUserEndpoint(serviceName, userName) + "/pool"

	if len(wanted) > 0 {
		var taskId string
		opts := &DedicatedCephUserPoolPermissionsCreateOpts{Permissions: wanted}
		err := retryDedicatedCephChange(ctx, dedicatedCephUserTimeout, func() error {
			return r.config.OVHClient.PostWithContext(ctx, endpoint, opts, &taskId)
		})
		if err != nil {
			return fmt.Errorf("error calling Post %s: %w", endpoint, err)
		}
		return waitDedicatedCephTask(ctx, r.config.OVHClient, serviceName, taskId, dedicatedCephUserTimeout)
	}

	for _, permission := range current {
		var taskId string
		poolEndpoint := endpoint + "/" + url.PathEscape(permission.PoolName)
		err := retryDedicatedCephChange(ctx, dedicatedCephUserTimeout, func() error {
			return r.config.OVHClient.DeleteWithContext(ctx, poolEndpoint, &taskId)
		})
		if err != nil {
			return fmt.Errorf("error calling Delete %s: %w", poolEndpoint, err)
		}
		if err := waitDedicatedCephTask(ctx, r.config.OVHClient, serviceName, taskId, dedicatedCephUserTimeout); err != nil {
			return err
		}
	}

	return nil
}

// fetch gets the user and its pool permissions.
func (r *dedicatedCephUserResource) fetch(ctx context.Context, serviceName, userName string) (*DedicatedCephUser, []DedicatedCephUserPoolPermission, error) {
	endpoint := dedicatedCephUserEndpoint(serviceName, userName)

	var user DedicatedCephUser
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &user); err != nil {
		return nil, nil, fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	var permissions []DedicatedCephUserPoolPermission
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint+"/pool", &permissions); err != nil {
		return nil, nil, fmt.Errorf("error calling Get %s/pool: %w", endpoint, err)
	}

	return &user, permissions, nil
}

// refresh fetches the user and its pool permissions into the model.
func (r *dedicatedCephUserResource) refresh(ctx context.Context, data *DedicatedCephUserModel) diag.Diagnostics {
	var diags diag.Diagnostics

	userName := data.Name.ValueString()
	user, permissions, err := r.fetch(ctx, data.ServiceName.ValueString(), userName)
	if err != nil {
		diags.AddError(fmt.Sprintf("Error reading user %s", userName), err.Error())
		return diags
	}

	data.MergeWith(user)
	diags.Append(data.MergeWithPoolPermissions(ctx, permissions)...)
	return diags
}

func (r *dedicatedCephUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DedicatedCephUserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The flags are all known now, which they may not be when validating the configuration.
	resp.Diagnostics.Append(validateDedicatedCephUserPoolPermissions(ctx, data.PoolPermissions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	userName := data.Name.ValueString()

	// Changes are made one at a time on a cluster, see lockDedicatedCephCluster.
	unlock, err := lockDedicatedCephCluster(ctx, serviceName)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for the other changes of %s", serviceName), err.Error())
		return
	}
	defer unlock()

	if err := waitDedicatedCephIdle(ctx, r.config.OVHClient, serviceName, dedicatedCephUserTimeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for the tasks in progress on %s", serviceName), err.Error())
		return
	}

	permissions, diags := data.ToPoolPermissions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var taskId string
	endpoint := "/dedicated/ceph/" + url.PathEscape(serviceName) + "/user"
	opts := &DedicatedCephUserCreateOpts{UserName: userName}
	err = retryDedicatedCephChange(ctx, dedicatedCephUserTimeout, func() error {
		return r.config.OVHClient.PostWithContext(ctx, endpoint, opts, &taskId)
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Save the user as soon as its creation is accepted, so that Terraform tracks it, and replaces it
	// on the next apply, when a later step fails, such as setting its pool permissions.
	created := data.CreatedState()
	resp.Diagnostics.Append(resp.State.Set(ctx, &created)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := waitDedicatedCephTask(ctx, r.config.OVHClient, serviceName, taskId, dedicatedCephUserTimeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for user %s to be created", userName), err.Error())
		return
	}

	if err := r.setPoolPermissions(ctx, serviceName, userName, nil, permissions); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error setting the pool permissions of user %s", userName), err.Error())
		return
	}

	resp.Diagnostics.Append(r.refresh(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dedicatedCephUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DedicatedCephUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userName := data.Name.ValueString()
	user, permissions, err := r.fetch(ctx, data.ServiceName.ValueString(), userName)
	if err != nil {
		// fetch wraps the client error, so unwrap it to reach the API error.
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error reading user %s", userName), err.Error())
		return
	}

	data.MergeWith(user)
	resp.Diagnostics.Append(data.MergeWithPoolPermissions(ctx, permissions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dedicatedCephUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state DedicatedCephUserModel

	// Read Terraform plan and prior state data into the models
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The flags are all known now, which they may not be when validating the configuration.
	resp.Diagnostics.Append(validateDedicatedCephUserPoolPermissions(ctx, data.PoolPermissions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	userName := data.Name.ValueString()

	// Changes are made one at a time on a cluster, see lockDedicatedCephCluster.
	unlock, err := lockDedicatedCephCluster(ctx, serviceName)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for the other changes of %s", serviceName), err.Error())
		return
	}
	defer unlock()

	if err := waitDedicatedCephIdle(ctx, r.config.OVHClient, serviceName, dedicatedCephUserTimeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for the tasks in progress on %s", serviceName), err.Error())
		return
	}

	// pool_permissions is the only attribute that can be updated in place.
	wanted, diags := data.ToPoolPermissions(ctx)
	resp.Diagnostics.Append(diags...)
	current, diags := state.ToPoolPermissions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.setPoolPermissions(ctx, serviceName, userName, current, wanted); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error setting the pool permissions of user %s", userName), err.Error())
		return
	}

	resp.Diagnostics.Append(r.refresh(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dedicatedCephUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DedicatedCephUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	userName := data.Name.ValueString()

	// Changes are made one at a time on a cluster, see lockDedicatedCephCluster.
	unlock, err := lockDedicatedCephCluster(ctx, serviceName)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for the other changes of %s", serviceName), err.Error())
		return
	}
	defer unlock()

	if err := waitDedicatedCephIdle(ctx, r.config.OVHClient, serviceName, dedicatedCephUserTimeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for the tasks in progress on %s", serviceName), err.Error())
		return
	}

	endpoint := dedicatedCephUserEndpoint(serviceName, userName)

	var taskId string
	err = retryDedicatedCephChange(ctx, dedicatedCephUserTimeout, func() error {
		return r.config.OVHClient.DeleteWithContext(ctx, endpoint, &taskId)
	})
	if err != nil {
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}

	if err := waitDedicatedCephTask(ctx, r.config.OVHClient, serviceName, taskId, dedicatedCephUserTimeout); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error waiting for user %s to be deleted", userName), err.Error())
	}
}

func (r *dedicatedCephUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.SplitN(req.ID, "/", 2)
	if len(splits) != 2 || splits[0] == "" || splits[1] == "" {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<user_name>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
