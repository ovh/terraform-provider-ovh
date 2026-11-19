package ovh

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestCloudSshKeyDataSourcesServiceNameFromEnv reads ovh_cloud_ssh_key and
// ovh_cloud_ssh_keys from a fake API with service_name omitted, while another
// argument is unknown at plan time. The read is then deferred to apply, and
// a resource using the service_name of the data source must not see it change
// between plan and apply.
func TestCloudSshKeyDataSourcesServiceNameFromEnv(t *testing.T) {
	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "env-project")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/1.0/auth/time":
			fmt.Fprint(w, time.Now().Unix())
		case r.URL.Path == "/1.0/auth/details":
			fmt.Fprint(w, `{"account": "xx1111-ovh", "method": "provider"}`)
		case r.URL.Path == "/v2/publicCloud/project/env-project/sshKey/my-key":
			fmt.Fprint(w, `{"name": "my-key", "publicKey": "ssh-ed25519 AAAA", "createdAt": "2026-01-01T00:00:00Z", "updatedAt": "2026-01-01T00:00:00Z"}`)
		case r.URL.Path == "/v2/publicCloud/project/env-project/sshKey":
			fmt.Fprint(w, `[{"name": "my-key", "publicKey": "ssh-ed25519 AAAA", "createdAt": "2026-01-01T00:00:00Z", "updatedAt": "2026-01-01T00:00:00Z"}]`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message": "not found"}`)
		}
	}))
	t.Cleanup(server.Close)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
provider "ovh" {
  endpoint           = %q
  application_key    = "key"
  application_secret = "secret"
  consumer_key       = "consumer"
}

resource "terraform_data" "name" {
  input = "my-key"
}

data "ovh_cloud_ssh_key" "key" {
  name = terraform_data.name.output
}

data "ovh_cloud_ssh_keys" "keys" {
  depends_on = [terraform_data.name]
}

resource "terraform_data" "projects" {
  input = [data.ovh_cloud_ssh_key.key.service_name, data.ovh_cloud_ssh_keys.keys.service_name]
}
`, server.URL+"/1.0"),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.ovh_cloud_ssh_key.key", "service_name", "env-project"),
				resource.TestCheckResourceAttr("data.ovh_cloud_ssh_keys.keys", "service_name", "env-project"),
				resource.TestCheckResourceAttr("terraform_data.projects", "output.0", "env-project"),
				resource.TestCheckResourceAttr("terraform_data.projects", "output.1", "env-project"),
			),
		}},
	})
}
