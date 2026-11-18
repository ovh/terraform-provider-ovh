package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// managedDatabaseDataSourceConfigured gives a data source the provider configuration.
type managedDatabaseDataSourceConfigured struct {
	config *Config
}

func (d *managedDatabaseDataSourceConfigured) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func managedDatabaseComputedString(description string) schema.StringAttribute {
	return schema.StringAttribute{Computed: true, Description: description}
}
