package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudManagedDatabaseConnectionPoolModel is the Terraform model of a PgBouncer connection pool of a
// managed PostgreSQL service.
type CloudManagedDatabaseConnectionPoolModel struct {
	ServiceName types.String `tfsdk:"service_name"`
	ServiceId   types.String `tfsdk:"service_id"`
	Engine      types.String `tfsdk:"engine"`
	Name        types.String `tfsdk:"name"`
	DatabaseId  types.String `tfsdk:"database_id"`
	Mode        types.String `tfsdk:"mode"`
	Size        types.Int64  `tfsdk:"size"`
	UserId      types.String `tfsdk:"user_id"`

	Id        types.String `tfsdk:"id"`
	Port      types.Int64  `tfsdk:"port"`
	SslMode   types.String `tfsdk:"ssl_mode"`
	Uri       types.String `tfsdk:"uri"`
	Checksum  types.String `tfsdk:"checksum"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// managedDatabaseConnectionPoolCreation is the publicCloud.database.postgresql.ConnectionPoolCreation
// model.
type managedDatabaseConnectionPoolCreation struct {
	DatabaseId string  `json:"databaseId"`
	Mode       string  `json:"mode"`
	Name       string  `json:"name"`
	Size       int64   `json:"size"`
	UserId     *string `json:"userId,omitempty"`
}

// managedDatabaseConnectionPool is the publicCloud.database.postgresql.ConnectionPool model. Its tfsdk
// tags map it onto the attributes of the connection pool data sources.
type managedDatabaseConnectionPool struct {
	Checksum   string  `json:"checksum" tfsdk:"checksum"`
	CreatedAt  string  `json:"createdAt" tfsdk:"created_at"`
	DatabaseId string  `json:"databaseId" tfsdk:"database_id"`
	Id         string  `json:"id" tfsdk:"id"`
	Mode       string  `json:"mode" tfsdk:"mode"`
	Name       string  `json:"name" tfsdk:"name"`
	Port       int64   `json:"port" tfsdk:"port"`
	Size       int64   `json:"size" tfsdk:"size"`
	SslMode    *string `json:"sslMode" tfsdk:"ssl_mode"`
	UpdatedAt  string  `json:"updatedAt" tfsdk:"updated_at"`
	Uri        string  `json:"uri" tfsdk:"uri"`
	UserId     *string `json:"userId" tfsdk:"user_id"`
}

func (m *CloudManagedDatabaseConnectionPoolModel) ToCreate() *managedDatabaseConnectionPoolCreation {
	return &managedDatabaseConnectionPoolCreation{
		DatabaseId: m.DatabaseId.ValueString(),
		Mode:       m.Mode.ValueString(),
		Name:       m.Name.ValueString(),
		Size:       m.Size.ValueInt64(),
		UserId:     knownStringPointer(m.UserId),
	}
}

// ToPut builds the body of the publicCloud.database.postgresql.ConnectionPoolReplace model, which takes
// every argument. A pool without user sends it as null.
func (m *CloudManagedDatabaseConnectionPoolModel) ToPut() map[string]any {
	return map[string]any{
		"databaseId": m.DatabaseId.ValueString(),
		"mode":       m.Mode.ValueString(),
		"size":       m.Size.ValueInt64(),
		"userId":     knownStringPointer(m.UserId),
	}
}

func (m *CloudManagedDatabaseConnectionPoolModel) MergeWith(pool *managedDatabaseConnectionPool) {
	m.Id = types.StringValue(pool.Id)
	m.Name = types.StringValue(pool.Name)
	m.DatabaseId = types.StringValue(pool.DatabaseId)
	m.Mode = types.StringValue(pool.Mode)
	m.Size = types.Int64Value(pool.Size)
	m.UserId = types.StringPointerValue(pool.UserId)
	m.Port = types.Int64Value(pool.Port)
	m.SslMode = types.StringPointerValue(pool.SslMode)
	m.Uri = types.StringValue(pool.Uri)
	m.Checksum = types.StringValue(pool.Checksum)
	m.CreatedAt = types.StringValue(pool.CreatedAt)
	m.UpdatedAt = types.StringValue(pool.UpdatedAt)
}
