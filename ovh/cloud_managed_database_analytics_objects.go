package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Shared by the resources of the managed database and managed analytics services.

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

// managedDatabaseEngineArgument returns the engine argument of a resource of several engines.
func managedDatabaseEngineArgument(engines []string, object string) schema.StringAttribute {
	return schema.StringAttribute{
		Required:      true,
		Description:   "Engine of the service: " + strings.Join(engines, ", ") + ". Changing it recreates the " + object + ".",
		Validators:    []validator.String{stringvalidator.OneOf(engines...)},
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
}
