package ovh

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// headersTransport is an http.RoundTripper middleware that injects a fixed
// set of extra HTTP headers on every outgoing request.
type headersTransport struct {
	headers map[string]string
	next    http.RoundTripper
}

func (t *headersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request so we don't mutate the caller's original.
	req = req.Clone(req.Context())
	for k, v := range t.headers {
		// Set the header key exactly as configured, bypassing
		// http.Header.Set's canonicalization (e.g. "x-ovh-nic" would
		// otherwise be rewritten to "X-Ovh-Nic" on the wire), since some
		// upstream gateways match header names literally.
		req.Header[k] = []string{v}
	}
	return t.next.RoundTrip(req)
}

// newHeadersTransport wraps the given RoundTripper with extra header
// injection. If headers is empty, next is returned unchanged.
func newHeadersTransport(headers map[string]string, next http.RoundTripper) http.RoundTripper {
	if len(headers) == 0 {
		return next
	}
	return &headersTransport{headers: headers, next: next}
}

// httpHeadersFromEnv reads extra HTTP headers from numbered
// OVH_HTTP_HEADERS_N environment variables, each formatted like a raw HTTP
// header line (e.g. OVH_HTTP_HEADERS_0="x-ovh-nic: ab123456-ovh"). It exists
// so that headers can be injected without editing provider configuration,
// e.g. from `make testacc`. Reading stops at the first unset index.
func httpHeadersFromEnv() map[string]string {
	headers := make(map[string]string)
	for i := 0; ; i++ {
		v, ok := os.LookupEnv(fmt.Sprintf("OVH_HTTP_HEADERS_%d", i))
		if !ok {
			break
		}
		k, val, found := strings.Cut(v, ":")
		if !found {
			continue
		}
		headers[strings.TrimSpace(k)] = strings.TrimSpace(val)
	}
	if len(headers) == 0 {
		return nil
	}
	return headers
}
