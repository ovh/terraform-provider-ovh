package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudManagedDatabaseDatabaseModel is the Terraform model of a logical database of a managed
// database service.
type CloudManagedDatabaseDatabaseModel struct {
	ServiceName types.String `tfsdk:"service_name"`
	ServiceId   types.String `tfsdk:"service_id"`
	Engine      types.String `tfsdk:"engine"`
	Name        types.String `tfsdk:"name"`

	Id        types.String `tfsdk:"id"`
	Default   types.Bool   `tfsdk:"default"`
	Checksum  types.String `tfsdk:"checksum"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// managedDatabaseDatabaseCreation is the publicCloud.database.DatabaseCreation model.
type managedDatabaseDatabaseCreation struct {
	Name string `json:"name"`
}

// managedDatabaseDatabase is the publicCloud.database.Database model. Its tfsdk tags map it onto the
// attributes of the database data sources.
type managedDatabaseDatabase struct {
	Checksum  string `json:"checksum" tfsdk:"checksum"`
	CreatedAt string `json:"createdAt" tfsdk:"created_at"`
	Default   bool   `json:"default" tfsdk:"default"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updatedAt" tfsdk:"updated_at"`
}

func (m *CloudManagedDatabaseDatabaseModel) MergeWith(database *managedDatabaseDatabase) {
	m.Id = types.StringValue(database.Id)
	m.Name = types.StringValue(database.Name)
	m.Default = types.BoolValue(database.Default)
	m.Checksum = types.StringValue(database.Checksum)
	m.CreatedAt = types.StringValue(database.CreatedAt)
	m.UpdatedAt = types.StringValue(database.UpdatedAt)
}
