package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// cloudManagedKubernetesDataSourceBase holds the Configure shared by the MKS data sources.
type cloudManagedKubernetesDataSourceBase struct {
	config *Config
}

func (d *cloudManagedKubernetesDataSourceBase) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.config = config
}

func managedKubernetesReferenceEndpoint(serviceName, path string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/reference/managedKubernetes/" + path
}
