package ovh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

// newCloudProjectKubeReadTestConfig returns a provider config talking to a fake API serving a
// cluster last updated at updatedAt, and the number of kubeconfigs it generated.
func newCloudProjectKubeReadTestConfig(t *testing.T, updatedAt string) (*Config, *int) {
	t.Helper()

	generated := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	mux.HandleFunc("/1.0/cloud/project/xxx/kube/my-kube", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"id": "my-kube", "status": "READY", "version": "1.35.2-4", "updatedAt": %q}`, updatedAt)
	})
	mux.HandleFunc("/1.0/cloud/project/xxx/kube/my-kube/kubeconfig", func(w http.ResponseWriter, r *http.Request) {
		generated++
		content := fmt.Sprintf(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ca-%[1]d
    server: https://my-kube.example
  name: my-kube
users:
- name: admin
  user:
    client-certificate-data: cert-%[1]d
    client-key-data: key-%[1]d
`, generated)
		body, _ := json.Marshal(map[string]string{"content": content})
		w.Write(body)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	return &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}, &generated
}

func TestCloudProjectKubeReadKubeconfig(t *testing.T) {
	const updatedAt = "2026-10-27T16:47:58Z"

	for name, tc := range map[string]struct {
		stateUpdatedAt string
		wantGenerated  int
	}{
		// Each generation creates new credentials: keep the kubeconfig of the state
		"cluster unchanged": {stateUpdatedAt: updatedAt, wantGenerated: 0},
		// A reset of the kubeconfig renews the certificate authority
		"cluster changed": {stateUpdatedAt: "2026-10-27T16:40:41Z", wantGenerated: 1},
		// State written before updated_at existed
		"no updated_at in state": {stateUpdatedAt: "", wantGenerated: 1},
	} {
		t.Run(name, func(t *testing.T) {
			config, generated := newCloudProjectKubeReadTestConfig(t, updatedAt)

			d := schema.TestResourceDataRaw(t, resourceCloudProjectKube().Schema, map[string]interface{}{})
			d.SetId("my-kube")
			d.Set(kubeServiceNameKey, "xxx")
			d.Set(kubeClusterUpdatedAtKey, tc.stateUpdatedAt)
			d.Set(kubeClusterKubeconfigKey, "previous kubeconfig")
			d.Set(kubeClusterKubeconfigAttributesKey, []map[string]interface{}{{
				kubeClusterKubeconfigHostKey:                 "https://my-kube.example",
				kubeClusterKubeconfigClusterCaCertificateKey: "ca-0",
				kubeClusterKubeconfigClientCertificateKey:    "cert-0",
				kubeClusterKubeconfigClientKeyKey:            "key-0",
			}})

			if err := resourceCloudProjectKubeRead(d, config); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if *generated != tc.wantGenerated {
				t.Errorf("generated %d kubeconfigs, want %d", *generated, tc.wantGenerated)
			}
			wantCa := "ca-0"
			if tc.wantGenerated > 0 {
				wantCa = "ca-1"
			}
			if ca := d.Get(kubeClusterKubeconfigAttributesKey + ".0." + kubeClusterKubeconfigClusterCaCertificateKey); ca != wantCa {
				t.Errorf("cluster_ca_certificate is %q, want %q", ca, wantCa)
			}
			if got := d.Get(kubeClusterUpdatedAtKey); got != updatedAt {
				t.Errorf("updated_at is %q, want %q", got, updatedAt)
			}
		})
	}
}
