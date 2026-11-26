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

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	acc "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	"go.uber.org/ratelimit"
)

const (
	okmsLogSubscriptionTestOkmsId = "7f3a82ac-a8d8-4c2a-ab0c-f6e86ddf6a7c"
	okmsLogSubscriptionTestId     = "4c8a4b7e-9d3e-4a0f-8a8e-2f0e6b1d7c11"
	okmsLogSubscriptionTestPath   = "/v2/okms/resource/" + okmsLogSubscriptionTestOkmsId + "/log/subscription/" + okmsLogSubscriptionTestId
)

func init() {
	acc.AddTestSweepers("ovh_okms_log_subscription", &acc.Sweeper{
		Name: "ovh_okms_log_subscription",
		F:    testSweepOkmsLogSubscription,
	})
}

// testSweepOkmsLogSubscription deletes the subscriptions of the test KMS that forward to a stream
// of the test Logs Data Platform service whose title has the test prefix. The subscriptions have
// no name of their own, so the stream is what tells the ones left by the tests apart.
func testSweepOkmsLogSubscription(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	okmsId := os.Getenv("OVH_OKMS")
	ldpServiceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	if okmsId == "" || ldpServiceName == "" {
		log.Print("[DEBUG] OVH_OKMS or OVH_DBAAS_LOGS_SERVICE_TEST is not set. No KMS log subscription to sweep")
		return nil
	}

	endpoint := okmsLogSubscriptionEndpoint(okmsId)
	var subs []AccountLogSubscription
	if err := client.Get(endpoint, &subs); err != nil {
		return fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	for _, sub := range subs {
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

		subscriptionEndpoint := endpoint + "/" + url.PathEscape(sub.SubscriptionId)
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

// newOkmsLogSubscriptionTestResource returns the resource talking to a fake API that answers the
// requests of the given method on the test subscription with the given status and body.
func newOkmsLogSubscriptionTestResource(t *testing.T, method string, status int, body string) *okmsLogSubscriptionResource {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	mux.HandleFunc(okmsLogSubscriptionTestPath, func(w http.ResponseWriter, r *http.Request) {
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

	r := NewOkmsLogSubscriptionResource().(*okmsLogSubscriptionResource)
	r.config = &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}
	return r
}

func newOkmsLogSubscriptionTestState(t *testing.T, model *OkmsLogSubscriptionModel) tfsdk.State {
	t.Helper()
	ctx := context.Background()

	var schemaResp resource.SchemaResponse
	NewOkmsLogSubscriptionResource().Schema(ctx, resource.SchemaRequest{}, &schemaResp)
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

func okmsLogSubscriptionTestModel() OkmsLogSubscriptionModel {
	return OkmsLogSubscriptionModel{
		OkmsId:   ovhtypes.NewTfStringValue(okmsLogSubscriptionTestOkmsId),
		Kind:     ovhtypes.NewTfStringValue("audit-secret"),
		StreamId: ovhtypes.NewTfStringValue("0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"),
	}.CreatedState(okmsLogSubscriptionTestId, "ldp-xx-12345")
}

func TestOkmsLogSubscriptionMetadata(t *testing.T) {
	var resp resource.MetadataResponse
	NewOkmsLogSubscriptionResource().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "ovh"}, &resp)
	if resp.TypeName != "ovh_okms_log_subscription" {
		t.Errorf("unexpected type name %q", resp.TypeName)
	}
}

func TestOkmsLogSubscriptionRead(t *testing.T) {
	ctx := context.Background()
	body := `{
		"createdAt": "2026-10-01T10:00:00Z",
		"kind": "audit-secret",
		"resource": {"name": "` + okmsLogSubscriptionTestOkmsId + `", "type": "okms"},
		"serviceName": "ldp-xx-12345",
		"streamId": "0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10",
		"subscriptionId": "` + okmsLogSubscriptionTestId + `",
		"updatedAt": "2026-10-02T10:00:00Z"
	}`
	r := newOkmsLogSubscriptionTestResource(t, http.MethodGet, http.StatusOK, body)

	// The state of an imported subscription, which only has the id of its KMS and its own.
	prior := OkmsLogSubscriptionModel{
		OkmsId:   ovhtypes.NewTfStringValue(okmsLogSubscriptionTestOkmsId),
		Kind:     ovhtypes.NewTfStringNull(),
		StreamId: ovhtypes.NewTfStringNull(),
	}.CreatedState(okmsLogSubscriptionTestId, "")
	prior.LdpServiceName = ovhtypes.NewTfStringNull()

	resp := resource.ReadResponse{State: newOkmsLogSubscriptionTestState(t, &prior)}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	var got OkmsLogSubscriptionModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	for attribute, pair := range map[string][2]string{
		"id":               {got.ID.ValueString(), okmsLogSubscriptionTestId},
		"okms_id":          {got.OkmsId.ValueString(), okmsLogSubscriptionTestOkmsId},
		"kind":             {got.Kind.ValueString(), "audit-secret"},
		"stream_id":        {got.StreamId.ValueString(), "0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"},
		"ldp_service_name": {got.LdpServiceName.ValueString(), "ldp-xx-12345"},
		"resource_name":    {got.ResourceName.ValueString(), okmsLogSubscriptionTestOkmsId},
		"resource_type":    {got.ResourceType.ValueString(), "okms"},
		"created_at":       {got.CreatedAt.ValueString(), "2026-10-01T10:00:00Z"},
		"updated_at":       {got.UpdatedAt.ValueString(), "2026-10-02T10:00:00Z"},
	} {
		if pair[0] != pair[1] {
			t.Errorf("unexpected %s %q, want %q", attribute, pair[0], pair[1])
		}
	}
}

// A subscription deleted outside Terraform is removed from the state, so that the next apply
// creates it again.
func TestOkmsLogSubscriptionReadNotFound(t *testing.T) {
	ctx := context.Background()
	body := `{"class":"Client::NotFound","message":"The requested object (subscriptionId = ` + okmsLogSubscriptionTestId + `) does not exist"}`
	r := newOkmsLogSubscriptionTestResource(t, http.MethodGet, http.StatusNotFound, body)

	prior := okmsLogSubscriptionTestModel()
	resp := resource.ReadResponse{State: newOkmsLogSubscriptionTestState(t, &prior)}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("expected the resource to be removed from the state, got %v", resp.State.Raw)
	}
}

func TestOkmsLogSubscriptionReadError(t *testing.T) {
	ctx := context.Background()
	body := `{"class":"Server::InternalServerError","message":"Internal server error"}`
	r := newOkmsLogSubscriptionTestResource(t, http.MethodGet, http.StatusInternalServerError, body)

	prior := okmsLogSubscriptionTestModel()
	resp := resource.ReadResponse{State: newOkmsLogSubscriptionTestState(t, &prior)}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error")
	}
	if resp.State.Raw.IsNull() {
		t.Error("expected the resource to stay in the state")
	}
}

// Deleting a subscription that is already gone succeeds, without waiting for an operation.
func TestOkmsLogSubscriptionDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	body := `{"class":"Client::NotFound","message":"not found"}`
	r := newOkmsLogSubscriptionTestResource(t, http.MethodDelete, http.StatusNotFound, body)

	prior := okmsLogSubscriptionTestModel()
	state := newOkmsLogSubscriptionTestState(t, &prior)

	resp := resource.DeleteResponse{State: state}
	r.Delete(ctx, resource.DeleteRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestOkmsLogSubscriptionImportState(t *testing.T) {
	ctx := context.Background()
	r := NewOkmsLogSubscriptionResource().(resource.ResourceWithImportState)

	resp := resource.ImportStateResponse{State: newOkmsLogSubscriptionTestState(t, nil)}
	r.ImportState(ctx, resource.ImportStateRequest{ID: okmsLogSubscriptionTestOkmsId + "/" + okmsLogSubscriptionTestId}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	var got OkmsLogSubscriptionModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if got.OkmsId.ValueString() != okmsLogSubscriptionTestOkmsId {
		t.Errorf("unexpected okms_id %q, want %q", got.OkmsId.ValueString(), okmsLogSubscriptionTestOkmsId)
	}
	if got.ID.ValueString() != okmsLogSubscriptionTestId {
		t.Errorf("unexpected id %q, want %q", got.ID.ValueString(), okmsLogSubscriptionTestId)
	}
}

func TestOkmsLogSubscriptionImportStateMalformed(t *testing.T) {
	ctx := context.Background()
	r := NewOkmsLogSubscriptionResource().(resource.ResourceWithImportState)

	for _, id := range []string{okmsLogSubscriptionTestId, "/" + okmsLogSubscriptionTestId, okmsLogSubscriptionTestOkmsId + "/", "a/b/c"} {
		resp := resource.ImportStateResponse{State: newOkmsLogSubscriptionTestState(t, nil)}
		r.ImportState(ctx, resource.ImportStateRequest{ID: id}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Errorf("expected an error for the id %q", id)
		}
	}
}

// The state saved as soon as a subscription is created must hold no unknown value, which a state
// cannot hold, while the plan it comes from has the attributes read from the API unknown.
func TestOkmsLogSubscriptionCreatedState(t *testing.T) {
	unknown := ovhtypes.TfStringValue{StringValue: basetypes.NewStringUnknown()}
	created := OkmsLogSubscriptionModel{
		ID:             unknown,
		OkmsId:         ovhtypes.NewTfStringValue(okmsLogSubscriptionTestOkmsId),
		Kind:           ovhtypes.NewTfStringValue("audit-secret"),
		StreamId:       ovhtypes.NewTfStringValue("0a6e2b8f-7c34-4d4c-9b55-1f4e2b3c9d10"),
		LdpServiceName: unknown,
		ResourceName:   unknown,
		ResourceType:   unknown,
		CreatedAt:      unknown,
		UpdatedAt:      unknown,
	}.CreatedState(okmsLogSubscriptionTestId, "ldp-xx-12345")

	state := newOkmsLogSubscriptionTestState(t, &created)
	if !state.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", state.Raw)
	}
	if got := created.ID.ValueString(); got != okmsLogSubscriptionTestId {
		t.Errorf("unexpected id %q", got)
	}
	if got := created.LdpServiceName.ValueString(); got != "ldp-xx-12345" {
		t.Errorf("unexpected ldp_service_name %q", got)
	}
}

const testAccOkmsLogSubscriptionConfig = `
resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "%s"
  title        = "%s"
  description  = "%s"
}

resource "ovh_okms_log_subscription" "sub" {
  okms_id   = "%s"
  kind      = "audit-kmip"
  stream_id = ovh_dbaas_logs_output_graylog_stream.stream.stream_id
}
`

// TestAccOkmsLogSubscription_basic subscribes the test KMS to a new stream, checks that a second
// plan is empty, then imports the subscription.
func TestAccOkmsLogSubscription_basic(t *testing.T) {
	okmsId := os.Getenv("OVH_OKMS")
	ldpServiceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	streamTitle := acctest.RandomWithPrefix(test_prefix)
	resourceName := "ovh_okms_log_subscription.sub"

	config := fmt.Sprintf(testAccOkmsLogSubscriptionConfig, ldpServiceName, streamTitle, streamTitle, okmsId)

	acc.Test(t, acc.TestCase{
		PreCheck: func() {
			testAccPreCheckOkms(t)
			testAccPreCheckDbaasLogs(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []acc.TestStep{
			{
				Config: config,
				Check: acc.ComposeTestCheckFunc(
					acc.TestCheckResourceAttrSet(resourceName, "id"),
					acc.TestCheckResourceAttr(resourceName, "okms_id", okmsId),
					acc.TestCheckResourceAttr(resourceName, "kind", "audit-kmip"),
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
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource %s not found in the state", resourceName)
					}
					return rs.Primary.Attributes["okms_id"] + "/" + rs.Primary.ID, nil
				},
			},
		},
	})
}
