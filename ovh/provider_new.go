package ovh

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &OvhProvider{}
)

// OvhProvider is the provider implementation.
type OvhProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *OvhProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ovh"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *OvhProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: descriptions["endpoint"],
			},
			"application_key": schema.StringAttribute{
				Optional:    true,
				Description: descriptions["application_key"],
			},
			"application_secret": schema.StringAttribute{
				Optional:    true,
				Description: descriptions["application_secret"],
			},
			"consumer_key": schema.StringAttribute{
				Optional:    true,
				Description: descriptions["consumer_key"],
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *OvhProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ovhProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown OVH API endpoint",
			"The provider cannot create the OVH API client as there is a missing or empty value for the API endpoint. "+
				"Set the endpoint value in the configuration and ensure the value is not empty.",
		)
	}

	if config.ApplicationKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("application_key"),
			"Unknown OVH API application_key",
			"The provider cannot create the OVH API client as there is a missing or empty value for the API application key. "+
				"Set the application key value in the configuration or use the OVH_APPLICATION_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if config.ApplicationSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("application_secret"),
			"Unknown OVH API application_secret",
			"The provider cannot create the OVH API client as there is a missing or empty value for the API application secret. "+
				"Set the application secret value in the configuration or use the OVH_APPLICATION_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if config.ConsumerKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("consumer_key"),
			"Unknown OVH API consumer_key",
			"The provider cannot create the OVH API client as there is a missing or empty value for the API consumer key. "+
				"Set the consumer key value in the configuration or use the OVH_CONSUMER_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	clientConfig := Config{
		lockAuth: &sync.Mutex{},
	}

	// Check if API variables has been set directly in the configuration
	if !config.Endpoint.IsNull() {
		clientConfig.Endpoint = config.Endpoint.ValueString()
	}
	if !config.ApplicationKey.IsNull() {
		clientConfig.ApplicationKey = config.ApplicationKey.ValueString()
	}
	if !config.ApplicationSecret.IsNull() {
		clientConfig.ApplicationSecret = config.ApplicationSecret.ValueString()
	}
	if !config.ConsumerKey.IsNull() {
		clientConfig.ConsumerKey = config.ConsumerKey.ValueString()
	}

	if err := clientConfig.loadAndValidate(); err != nil {
		resp.Diagnostics.AddError(err.Error(), "failed to init OVH API client")
		return
	}

	resp.DataSourceData = &clientConfig
	resp.ResourceData = &clientConfig
}

// DataSources defines the data sources implemented in the provider.
func (p *OvhProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCloudProjectDatabaseIPRestrictionsDataSource,
		NewCloudProjectDataSource,
		NewCloudProjectsDataSource,
		NewDedicatedServerSpecificationsHardwareDataSource,
		NewDedicatedServerSpecificationsNetworkDataSource,
		NewDomainZoneDnssecDataSource,
		NewIpFirewallDataSource,
		NewIpFirewallRuleDataSource,
		NewIpMitigationDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *OvhProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCloudProjectAlertingResource,
		NewDomainZoneDnssecResource,
		NewIpFirewallResource,
		NewIpFirewallRuleResource,
		NewIploadbalancingUdpFrontendResource,
		NewIpMitigationResource,
		NewVpsResource,
	}
}

type ovhProviderModel struct {
	Endpoint          types.String `tfsdk:"endpoint"`
	ApplicationKey    types.String `tfsdk:"application_key"`
	ApplicationSecret types.String `tfsdk:"application_secret"`
	ConsumerKey       types.String `tfsdk:"consumer_key"`
}
