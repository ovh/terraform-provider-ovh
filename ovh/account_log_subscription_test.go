package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	acc "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	"go.uber.org/ratelimit"
)

const accountLogSubscriptionTestId = "4c8a4b7e-9d3e-4a0f-8a8e-2f0e6b1d7c11"

// newAccountLogSubscriptionTestResource returns the resource built by newResource, talking to a
// fake API that answers the requests of the given method and path with the given status and body.
// The API is served under /1.0, as the real one, so that the /v2 paths are routed as in production.
func newAccountLogSubscriptionTestResource(t *testing.T, newResource func() resource.Resource, method, apiPath string, status int, body string) *accountLogSubscriptionResource {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	mux.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	r := newResource().(*accountLogSubscriptionResource)
	r.config = &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}
	return r
}

func newAccountLogSubscriptionTestState(t *testing.T, r resource.Resource, model *AccountLogSubscriptionModel) tfsdk.State {
	t.Helper()
	ctx := context.Background()

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", schemaResp.Diagnostics)
	}

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	if model != nil {
		if diags := state.Set(ctx, model); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
	}
	return state
}

func accountLogSubscriptionResources() map[string]struct {
	newResource func() resource.Resource
	typeName    string
	apiPath     string
} {
	return map[string]struct {
		newResource func() resource.Resource
		typeName    string
		apiPath     string
	}{
		"audit": {NewMeLogsAuditLogSubscriptionResource, "ovh_me_logs_audit_log_subscription", "/1.0/me/logs/audit/log/subscription/"},
		"api":   {NewMeApiLogSubscriptionResource, "ovh_me_api_log_subscription", "/1.0/me/api/log/subscription/"},
		"iam":   {NewIamLogSubscriptionResource, "ovh_iam_log_subscription", "/v2/iam/log/subscription/"},
	}
}

func TestAccountLogSubscriptionMetadata(t *testing.T) {
	for name, tc := range accountLogSubscriptionResources() {
		t.Run(name, func(t *testing.T) {
			var resp resource.MetadataResponse
			tc.newResource().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "ovh"}, &resp)
			if resp.TypeName != tc.typeName {
				t.Errorf("unexpected type name %q, want %q", resp.TypeName, tc.typeName)
			}
		})
	}
}

func TestAccountLogSubscriptionRead(t *testing.T) {
	ctx := context.Background()
	body := `{
		"createdAt": "2026-10-01T10:00:00Z",
		"kind": "default",
		"resource": {"name": "xx1111-ovh", "type": "account"},
		"serviceName": "ldp-xx-12345",
		"streamId": "0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10",
		"subscriptionId": "` + accountLogSubscriptionTestId + `",
		"updatedAt": "2026-10-02T10:00:00Z"
	}`

	for name, tc := range accountLogSubscriptionResources() {
		t.Run(name, func(t *testing.T) {
			r := newAccountLogSubscriptionTestResource(t, tc.newResource, http.MethodGet, tc.apiPath+accountLogSubscriptionTestId, http.StatusOK, body)

			// The state of an imported subscription, which only has its id.
			prior := AccountLogSubscriptionModel{ID: ovhtypes.NewTfStringValue(accountLogSubscriptionTestId)}.CreatedState(accountLogSubscriptionTestId, "")
			prior.Kind = ovhtypes.NewTfStringNull()
			prior.StreamId = ovhtypes.NewTfStringNull()
			prior.LdpServiceName = ovhtypes.NewTfStringNull()

			resp := resource.ReadResponse{State: newAccountLogSubscriptionTestState(t, r, &prior)}
			r.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			var got AccountLogSubscriptionModel
			resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			for attribute, pair := range map[string][2]string{
				"id":               {got.ID.ValueString(), accountLogSubscriptionTestId},
				"kind":             {got.Kind.ValueString(), "default"},
				"stream_id":        {got.StreamId.ValueString(), "0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"},
				"ldp_service_name": {got.LdpServiceName.ValueString(), "ldp-xx-12345"},
				"resource_name":    {got.ResourceName.ValueString(), "xx1111-ovh"},
				"resource_type":    {got.ResourceType.ValueString(), "account"},
				"created_at":       {got.CreatedAt.ValueString(), "2026-10-01T10:00:00Z"},
				"updated_at":       {got.UpdatedAt.ValueString(), "2026-10-02T10:00:00Z"},
			} {
				if pair[0] != pair[1] {
					t.Errorf("unexpected %s %q, want %q", attribute, pair[0], pair[1])
				}
			}
		})
	}
}

// A subscription deleted outside Terraform is removed from the state, so that the next apply
// creates it again.
func TestAccountLogSubscriptionReadNotFound(t *testing.T) {
	ctx := context.Background()
	body := `{"class":"Client::NotFound","message":"The requested object (subscriptionId = ` + accountLogSubscriptionTestId + `) does not exist"}`

	for name, tc := range accountLogSubscriptionResources() {
		t.Run(name, func(t *testing.T) {
			r := newAccountLogSubscriptionTestResource(t, tc.newResource, http.MethodGet, tc.apiPath+accountLogSubscriptionTestId, http.StatusNotFound, body)

			prior := AccountLogSubscriptionModel{
				Kind:     ovhtypes.NewTfStringValue("default"),
				StreamId: ovhtypes.NewTfStringValue("0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"),
			}.CreatedState(accountLogSubscriptionTestId, "ldp-xx-12345")

			resp := resource.ReadResponse{State: newAccountLogSubscriptionTestState(t, r, &prior)}
			r.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if !resp.State.Raw.IsNull() {
				t.Errorf("expected the resource to be removed from the state, got %v", resp.State.Raw)
			}
		})
	}
}

func TestAccountLogSubscriptionReadError(t *testing.T) {
	ctx := context.Background()
	body := `{"class":"Server::InternalServerError","message":"Internal server error"}`

	r := newAccountLogSubscriptionTestResource(t, NewMeLogsAuditLogSubscriptionResource, http.MethodGet,
		"/1.0/me/logs/audit/log/subscription/"+accountLogSubscriptionTestId, http.StatusInternalServerError, body)

	prior := AccountLogSubscriptionModel{
		Kind:     ovhtypes.NewTfStringValue("default"),
		StreamId: ovhtypes.NewTfStringValue("0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"),
	}.CreatedState(accountLogSubscriptionTestId, "ldp-xx-12345")

	resp := resource.ReadResponse{State: newAccountLogSubscriptionTestState(t, r, &prior)}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error")
	}
	if resp.State.Raw.IsNull() {
		t.Error("expected the resource to stay in the state")
	}
}

// Deleting a subscription that is already gone succeeds, without waiting for an operation.
func TestAccountLogSubscriptionDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	body := `{"class":"Client::NotFound","message":"not found"}`

	for name, tc := range accountLogSubscriptionResources() {
		t.Run(name, func(t *testing.T) {
			r := newAccountLogSubscriptionTestResource(t, tc.newResource, http.MethodDelete, tc.apiPath+accountLogSubscriptionTestId, http.StatusNotFound, body)

			prior := AccountLogSubscriptionModel{
				Kind:     ovhtypes.NewTfStringValue("default"),
				StreamId: ovhtypes.NewTfStringValue("0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"),
			}.CreatedState(accountLogSubscriptionTestId, "ldp-xx-12345")
			state := newAccountLogSubscriptionTestState(t, r, &prior)

			resp := resource.DeleteResponse{State: state}
			r.Delete(ctx, resource.DeleteRequest{State: state}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}

func TestAccountLogSubscriptionImportState(t *testing.T) {
	ctx := context.Background()
	r := NewIamLogSubscriptionResource()

	resp := resource.ImportStateResponse{State: newAccountLogSubscriptionTestState(t, r, nil)}
	r.(resource.ResourceWithImportState).ImportState(ctx, resource.ImportStateRequest{ID: accountLogSubscriptionTestId}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	var got ovhtypes.TfStringValue
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("id"), &got)...)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if got.ValueString() != accountLogSubscriptionTestId {
		t.Errorf("unexpected id %q, want %q", got.ValueString(), accountLogSubscriptionTestId)
	}
}

// The state saved as soon as a subscription is created must hold no unknown value, which a state
// cannot hold, while the plan it comes from has the attributes read from the API unknown.
func TestAccountLogSubscriptionCreatedState(t *testing.T) {
	unknown := ovhtypes.TfStringValue{StringValue: basetypes.NewStringUnknown()}
	created := AccountLogSubscriptionModel{
		ID:             unknown,
		Kind:           ovhtypes.NewTfStringValue("default"),
		StreamId:       ovhtypes.NewTfStringValue("0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"),
		LdpServiceName: unknown,
		ResourceName:   unknown,
		ResourceType:   unknown,
		CreatedAt:      unknown,
		UpdatedAt:      unknown,
	}.CreatedState(accountLogSubscriptionTestId, "ldp-xx-12345")

	state := newAccountLogSubscriptionTestState(t, NewMeApiLogSubscriptionResource(), &created)
	if !state.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", state.Raw)
	}
	if got := created.ID.ValueString(); got != accountLogSubscriptionTestId {
		t.Errorf("unexpected id %q", got)
	}
	if got := created.LdpServiceName.ValueString(); got != "ldp-xx-12345" {
		t.Errorf("unexpected ldp_service_name %q", got)
	}
}

// testSweepAccountLogSubscription deletes the subscriptions listed by endpoint that forward to a
// stream of the test Logs Data Platform service whose title has the test prefix. The subscriptions
// have no name of their own, so the stream is what tells the ones left by the tests apart.
func testSweepAccountLogSubscription(region, endpoint string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	ldpServiceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	if ldpServiceName == "" {
		log.Printf("[DEBUG] OVH_DBAAS_LOGS_SERVICE_TEST is not set. No subscription of %s to sweep", endpoint)
		return nil
	}

	var ids []string
	if err := client.Get(endpoint, &ids); err != nil {
		return fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	for _, id := range ids {
		subscriptionEndpoint := endpoint + "/" + url.PathEscape(id)

		var sub AccountLogSubscription
		if err := client.Get(subscriptionEndpoint, &sub); err != nil {
			return fmt.Errorf("error calling Get %s: %w", subscriptionEndpoint, err)
		}
		if sub.ServiceName != ldpServiceName {
			continue
		}

		streamEndpoint := fmt.Sprintf("/dbaas/logs/%s/output/graylog/stream/%s", url.PathEscape(sub.ServiceName), url.PathEscape(sub.StreamId))
		var stream DbaasLogsOutputGraylogStream
		if err := client.Get(streamEndpoint, &stream); err != nil {
			var errOvh *ovh.APIError
			if errors.As(err, &errOvh) && errOvh.Code == http.StatusNotFound {
				continue
			}
			return fmt.Errorf("error calling Get %s: %w", streamEndpoint, err)
		}
		if !strings.HasPrefix(stream.Title, test_prefix) {
			continue
		}

		log.Printf("[INFO] Deleting subscription %s forwarding to stream %s", subscriptionEndpoint, stream.Title)
		var op AccountLogSubscriptionOperation
		if err := client.Delete(subscriptionEndpoint, &op); err != nil {
			return fmt.Errorf("error calling Delete %s: %w", subscriptionEndpoint, err)
		}
		if _, err := waitForDbaasLogsOperation(context.Background(), client, op.ServiceName, op.OperationId); err != nil {
			return err
		}
	}

	return nil
}

const testAccAccountLogSubscriptionConfig = `
resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "%s"
  title        = "%s"
  description  = "%s"
}

resource "%s" "sub" {
  kind      = "default"
  stream_id = ovh_dbaas_logs_output_graylog_stream.stream.stream_id
}
`

// testAccAccountLogSubscription creates a subscription of the given resource type to a new stream,
// checks that a second plan is empty, then imports it.
func testAccAccountLogSubscription(t *testing.T, resourceType string) {
	ldpServiceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	streamTitle := acctest.RandomWithPrefix(test_prefix)
	resourceName := resourceType + ".sub"

	config := fmt.Sprintf(testAccAccountLogSubscriptionConfig, ldpServiceName, streamTitle, streamTitle, resourceType)

	acc.Test(t, acc.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogs(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []acc.TestStep{
			{
				Config: config,
				Check: acc.ComposeTestCheckFunc(
					acc.TestCheckResourceAttrSet(resourceName, "id"),
					acc.TestCheckResourceAttr(resourceName, "kind", "default"),
					acc.TestCheckResourceAttrPair(resourceName, "stream_id", "ovh_dbaas_logs_output_graylog_stream.stream", "stream_id"),
					acc.TestCheckResourceAttr(resourceName, "ldp_service_name", ldpServiceName),
					acc.TestCheckResourceAttrSet(resourceName, "resource_name"),
					acc.TestCheckResourceAttrSet(resourceName, "resource_type"),
					acc.TestCheckResourceAttrSet(resourceName, "created_at"),
					acc.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
