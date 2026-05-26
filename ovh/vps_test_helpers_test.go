package ovh

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
)

// checkVPSOptionSubscribed skips the test if OVH_VPS doesn't have the given option subscribed.
// Requires testAccPreCheckVPS to have run first (i.e. OVH_VPS env var set + credentials loaded).
func checkVPSOptionSubscribed(t *testing.T, optionName string) {
	t.Helper()
	serviceName := os.Getenv("OVH_VPS")
	if serviceName == "" {
		return // testAccPreCheckVPS handles this
	}
	checkVPSServiceOptionSubscribed(t, serviceName, optionName)
}

// checkVPSServiceOptionSubscribed skips the test if the given VPS service doesn't have
// the given option subscribed.
func checkVPSServiceOptionSubscribed(t *testing.T, serviceName, optionName string) {
	t.Helper()
	if serviceName == "" {
		return
	}
	var options []string
	endpoint := fmt.Sprintf("/vps/%s/option", url.PathEscape(serviceName))
	if err := testAccOVHClient.Get(endpoint, &options); err != nil {
		t.Skipf("could not list VPS options on %s: %s — skipping", serviceName, err)
		return
	}
	for _, o := range options {
		if o == optionName {
			return // option present, run the test
		}
	}
	t.Skipf("VPS %s does not have %q option subscribed — skipping (subscribe the option to enable this test)", serviceName, optionName)
}

// skipIfEndpointMissing performs a HEAD-like GET probe against the given path and
// skips the test if the API returns 404. Used for endpoints that are region-conditional
// (e.g., /vps/{sn}/models exists on EU/CA but not on the US schema).
func skipIfEndpointMissing(t *testing.T, path string) {
	t.Helper()
	if testAccOVHClient == nil {
		return
	}
	var raw any
	err := testAccOVHClient.Get(path, &raw)
	if err == nil {
		return
	}
	if apiErr, ok := err.(interface{ Error() string }); ok {
		msg := apiErr.Error()
		if strings.Contains(msg, "404") {
			t.Skipf("endpoint %s not available on this region/lineup (got: %s) — skipping", path, msg)
		}
	}
}
