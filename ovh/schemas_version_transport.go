package ovh

import (
	"net/http"
	"strings"
)

// schemasVersionHeader is the HTTP header to select a specific
// API schema version when processing the request.
const schemasVersionHeader = "X-Schemas-Version"

// schemasVersion is the schema version requested by this client.
// Bump this value when a new schema version is released and the provider
// has been updated to match.
const schemasVersion = "1.0"

// schemasVersionTransport is an http.RoundTripper middleware that injects the
// X-Schemas-Version header on every request whose path starts with "/v2/".
// Requests targeting v1 endpoints are left untouched because v1 instances do
// not support multi-version schemas.
type schemasVersionTransport struct {
	next http.RoundTripper
}

func (t *schemasVersionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.HasPrefix(req.URL.Path, "/v2/") {
		// Clone the request so we don't mutate the caller's original.
		req = req.Clone(req.Context())
		req.Header.Set(schemasVersionHeader, schemasVersion)
	}
	return t.next.RoundTrip(req)
}

// newSchemasVersionTransport wraps the given RoundTripper with schema version
// header injection for v2 API calls.
func newSchemasVersionTransport(next http.RoundTripper) http.RoundTripper {
	return &schemasVersionTransport{next: next}
}
