package ovh

// Request body parameters for a oauth2 client POST request
type ApiOauth2ClientCreateOpts struct {
	CallbackUrls []string `json:"callbackUrls"`
	Description  string   `json:"description"`
	Flow         string   `json:"flow"`
	Name         string   `json:"name"`
}

// Response body parameters from a successful oauth2 client POST request
type ApiOauth2ClientCreateResponse struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Response body parameters from a successful oauth2 client GET request
type ApiOauth2ClientReadResponse struct {
	CallbackUrls []string `json:"callbackUrls"`
	ClientId     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Description  string   `json:"description"`
	Flow         string   `json:"flow"`
	Name         string   `json:"name"`
	Identity     string   `json:"identity"`
}

// Request body parameters for a oauth2 client PUT request
type ApiOauth2ClientUpdateOpts struct {
	CallbackUrls []string `json:"callbackUrls"`
	Description  string   `json:"description"`
	Name         string   `json:"name"`
}
