package ovh

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCloudManagedAnalyticsKafkaUserAccessRead(t *testing.T) {
	ctx := context.Background()
	state := managedDatabaseDataSourceTestRead(t, NewCloudManagedAnalyticsKafkaUserAccessDataSource(),
		map[string]string{"service_name": "pj", "service_id": managedDatabasePostgresqlTestId, "user_id": "u1"},
		map[string][]managedDatabaseTestAnswer{
			"GET /v2/publicCloud/project/pj/managedAnalytics/kafka/" + managedDatabasePostgresqlTestId + "/user/u1/access": {{http.StatusOK, `{"cert": "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n", "key": "-----BEGIN PRIVATE KEY-----\nMIIE\n-----END PRIVATE KEY-----\n"}`}},
		})
	var engine, cert, key types.String
	state.GetAttribute(ctx, path.Root("engine"), &engine)
	state.GetAttribute(ctx, path.Root("cert"), &cert)
	state.GetAttribute(ctx, path.Root("key"), &key)
	if engine.ValueString() != "kafka" || cert.ValueString() == "" || key.ValueString() == "" {
		t.Errorf("unexpected state: engine %s, cert %s, key %s", engine, cert, key)
	}
}
