package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// cloudManagedDatabaseIPRestrictionsDataSource wraps cloudProjectDatabaseIPRestrictionsDataSource
// with the new resource name ovh_cloud_managed_database_ip_restrictions.
type cloudManagedDatabaseIPRestrictionsDataSource struct {
	cloudProjectDatabaseIPRestrictionsDataSource
}

func NewCloudManagedDatabaseIPRestrictionsDataSource() datasource.DataSource {
	return &cloudManagedDatabaseIPRestrictionsDataSource{}
}

func (d *cloudManagedDatabaseIPRestrictionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_database_ip_restrictions"
}

// cloudManagedAnalyticsIPRestrictionsDataSource wraps cloudProjectDatabaseIPRestrictionsDataSource
// with the analytics alias name ovh_cloud_managed_analytics_ip_restrictions.
type cloudManagedAnalyticsIPRestrictionsDataSource struct {
	cloudProjectDatabaseIPRestrictionsDataSource
}

func NewCloudManagedAnalyticsIPRestrictionsDataSource() datasource.DataSource {
	return &cloudManagedAnalyticsIPRestrictionsDataSource{}
}

func (d *cloudManagedAnalyticsIPRestrictionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_managed_analytics_ip_restrictions"
}
