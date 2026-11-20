package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

// Shared by the resources of the managed database and managed analytics services and of their
// objects: users, databases and connection pools.

const managedDatabaseObjectTimeout = 20 * time.Minute

// managedDatabaseCategory is a family of services of the Cloud Databases API v2: the ovh_cloud_managed_<family>
// resources, under the URL segment of the family, for its engines.
type managedDatabaseCategory struct {
	family  string
	segment string
	engines []string
}

var (
	managedDatabaseFamily  = managedDatabaseCategory{family: "database", segment: "managedDatabase", engines: []string{"mongodb", "mysql", "postgresql", "valkey"}}
	managedAnalyticsFamily = managedDatabaseCategory{family: "analytics", segment: "managedAnalytics", engines: []string{"clickhouse", "grafana", "kafka", "kafkaConnect", "kafkaMirrorMaker", "opensearch"}}
)

// typeName returns the type name of a resource or data source of the family, after the provider name.
func (c managedDatabaseCategory) typeName(suffix string) string {
	return "_cloud_managed_" + c.family + suffix
}

// endpoint returns the path of the services of an engine.
func (c managedDatabaseCategory) endpoint(serviceName, engine string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/" + c.segment + "/" + url.PathEscape(engine)
}

func (c managedDatabaseCategory) serviceEndpoint(serviceName, engine, serviceId string) string {
	return c.endpoint(serviceName, engine) + "/" + url.PathEscape(serviceId)
}

// managedDatabaseConfigured gives a resource the provider configuration.
type managedDatabaseConfigured struct {
	config *Config
}

func (r *managedDatabaseConfigured) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// managedDatabaseStringFormat validates a string with parse, whose error says what is wrong with it.
type managedDatabaseStringFormat struct {
	name  string
	parse func(string) error
}

func (v managedDatabaseStringFormat) Description(context.Context) string {
	return "value must be a valid " + v.name
}

func (v managedDatabaseStringFormat) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v managedDatabaseStringFormat) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if err := v.parse(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid "+v.name, err.Error())
	}
}

var managedDatabaseDuration = managedDatabaseStringFormat{"duration", func(s string) error {
	_, err := time.ParseDuration(s)
	return err
}}

// managedDatabaseCIDR is the strict CIDR block that v1 takes: no bit set after the prefix.
var managedDatabaseCIDR = managedDatabaseStringFormat{"CIDR block", func(s string) error {
	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return err
	}
	if masked := prefix.Masked().String(); masked != s {
		return fmt.Errorf("%s has bits set after its prefix length: use %s", s, masked)
	}
	return nil
}}

// managedDatabaseEngineArgument returns the engine argument of a resource of several engines.
func managedDatabaseEngineArgument(engines []string, object string) schema.StringAttribute {
	return schema.StringAttribute{
		Required:      true,
		Description:   "Engine of the service: " + strings.Join(engines, ", ") + ". Changing it recreates the " + object + ".",
		Validators:    []validator.String{stringvalidator.OneOf(engines...)},
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
}

// managedDatabaseObjectResource holds what the resources of the objects of a service share: the
// provider configuration, the family, the engines that have the object and the end of the type name.
type managedDatabaseObjectResource struct {
	managedDatabaseConfigured
	category managedDatabaseCategory
	engines  []string
	suffix   string
}

// managedDatabaseEngineInName tells whether the type name ending with suffix holds the engine, which
// then isn't an argument but a read-only attribute, for the resources and data sources alike.
func managedDatabaseEngineInName(engines []string, suffix string) bool {
	return len(engines) == 1 && strings.HasPrefix(suffix, "_"+strings.ToLower(engines[0])+"_")
}

func (r *managedDatabaseObjectResource) named() bool {
	return managedDatabaseEngineInName(r.engines, r.suffix)
}

func (r *managedDatabaseObjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.category.typeName(r.suffix)
}

// serviceEndpoint returns the path of the service of an object.
func (r *managedDatabaseObjectResource) serviceEndpoint(serviceName, engine, serviceId types.String) string {
	return r.category.serviceEndpoint(serviceName.ValueString(), engine.ValueString(), serviceId.ValueString())
}

// describe returns the description of a resource that creates an object in a service of the engines.
func (r *managedDatabaseObjectResource) describe(what string) string {
	return "Creates " + what + " in a managed " + r.category.family + " service (" + strings.Join(r.engines, ", ") + "), using the Cloud Databases API v2."
}

// objectAttributes returns the attributes that every object of a service has: the project, the engine
// and the service, and the ID of the object, here named.
func (r *managedDatabaseObjectResource) objectAttributes(object string) map[string]schema.Attribute {
	engine := managedDatabaseEngineArgument(r.engines, object)
	if r.named() {
		engine = schema.StringAttribute{
			Computed:    true,
			Description: "Engine of the service: " + r.engines[0] + ".",
			Default:     stringdefault.StaticString(r.engines[0]),
		}
	}
	return map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "ID of the public cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used. Changing it recreates the " + object + ".",
			PlanModifiers: []planmodifier.String{
				EnvDefaultString("OVH_CLOUD_PROJECT_SERVICE", true),
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"engine": engine,
		"service_id": schema.StringAttribute{
			Required:      true,
			Description:   "ID of the service. Changing it recreates the " + object + ".",
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"id": schema.StringAttribute{
			Computed:      true,
			Description:   "ID of the " + object + ".",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
	}
}

// importState imports an object from <service_name>/<engine>/<service_id>/<key>, without the engine
// when the type name holds it.
func (r *managedDatabaseObjectResource) importState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, key string) {
	format := []string{"service_name"}
	if !r.named() {
		format = append(format, "engine")
	}
	format = append(format, "service_id", key)

	splits := strings.Split(req.ID, "/")
	if len(splits) != len(format) {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <"+strings.Join(format, ">/<")+">")
		return
	}
	for i, name := range format {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(name), splits[i])...)
	}
	if r.named() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("engine"), r.engines[0])...)
	}
}

// waitForManagedDatabaseObject polls an object until it is READY, or until it is gone when deleted is
// true. ERROR is terminal. An object the API gives no status, a database or a connection pool, is READY
// while it exists. Unlike a service, an object is often ready at once: the first read has no delay.
func waitForManagedDatabaseObject(ctx context.Context, client *ovhwrap.Client, endpoint string, deleted bool, timeout time.Duration) error {
	pending, target := []string{"CREATING", "UPDATING", managedDatabaseAPIUnavailable}, []string{"READY"}
	if deleted {
		// The object can still report its previous status right after the DELETE.
		pending, target = []string{"CREATING", "UPDATING", "DELETING", "READY", managedDatabaseAPIUnavailable}, []string{"DELETED"}
	}
	stateConf := &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (any, string, error) {
			var res struct {
				Status string `json:"status"`
			}
			err := client.GetWithContext(ctx, endpoint, &res)
			if deleted && helpers.IsNotFound(err) {
				return res, "DELETED", nil
			}
			if isManagedDatabaseAPIUnavailable(err) {
				return res, managedDatabaseAPIUnavailable, nil
			}
			if err != nil {
				return nil, "", err
			}
			if res.Status == "" {
				res.Status = "READY"
			}
			if res.Status == "ERROR" {
				return nil, "", fmt.Errorf("%s is in ERROR status", endpoint)
			}
			return res, res.Status, nil
		},
		Timeout:    timeout,
		MinTimeout: managedDatabaseServicePollMinTimeout,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// managedDatabaseObjectNotReadyClass is the 409 of a change to a user that isn't READY: MongoDB applies
// a change asynchronously, and the PostgreSQL engine of apiv2 refuses a change meanwhile. A user in
// ERROR gets it too, for good.
const managedDatabaseObjectNotReadyClass = "Client::Conflict::UserStatusNotReady"

// retryWhileManagedDatabaseObjectBusy calls call, retrying while the service is busy, and while the
// object at endpoint isn't READY: it then waits for the object, which fails on ERROR, and calls again.
func retryWhileManagedDatabaseObjectBusy(ctx context.Context, client *ovhwrap.Client, endpoint string, timeout time.Duration, call func() error) error {
	return retryWhileManagedDatabaseServiceBusy(ctx, timeout, func() error {
		err := call()
		var apiErr *ovh.APIError
		if !errors.As(err, &apiErr) || apiErr.Code != http.StatusConflict || apiErr.Class != managedDatabaseObjectNotReadyClass {
			return err
		}
		if err := waitForManagedDatabaseObject(ctx, client, endpoint, false, timeout); err != nil {
			return err
		}
		return call()
	})
}

// deleteManagedDatabaseObject deletes an object, retrying while the service or the object is busy. An
// object already gone isn't an error.
func deleteManagedDatabaseObject(ctx context.Context, client *ovhwrap.Client, endpoint string) error {
	err := retryWhileManagedDatabaseObjectBusy(ctx, client, endpoint, managedDatabaseObjectTimeout, func() error {
		return client.DeleteWithContext(ctx, endpoint, nil)
	})
	if helpers.IsNotFound(err) {
		return nil
	}
	return err
}
