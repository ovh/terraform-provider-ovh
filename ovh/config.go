package ovh

import (
	"fmt"
	"log"
	"sync"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

var providerVersion, providerCommit string

type Config struct {
	Account  string
	Plate    string
	Endpoint string

	// Access token
	AccessToken string

	// AK / AS / CK authentication information
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string

	// oAuth2 authentication information
	ClientID     string
	ClientSecret string

	// Extra user-agent information
	UserAgentExtra string

	// Ignore initialization errors
	IgnoreInitError bool

	OVHClient     *ovhwrap.Client
	authenticated bool
	authFailed    error
	lockAuth      *sync.Mutex

	ApiRateLimit ratelimit.Limiter
}

func clientDefault(c *Config) (*ovh.Client, error) {
	var (
		client *ovh.Client
		err    error
	)

	switch {
	case c.AccessToken != "":
		client, err = ovh.NewAccessTokenClient(
			c.Endpoint,
			c.AccessToken,
		)
	case c.ClientID != "":
		client, err = ovh.NewOAuth2Client(
			c.Endpoint,
			c.ClientID,
			c.ClientSecret,
		)
	default:
		client, err = ovh.NewClient(
			c.Endpoint,
			c.ApplicationKey,
			c.ApplicationSecret,
			c.ConsumerKey,
		)
	}

	if err != nil {
		return nil, err
	}

	// Retrieve endpoint that is used
	for k, v := range ovh.Endpoints {
		if v == client.Endpoint() {
			c.Endpoint = k
		}
	}

	client.UserAgent = "Terraform/" + providerVersion + "/" + providerCommit
	if c.UserAgentExtra != "" {
		client.UserAgent += " - " + c.UserAgentExtra
	}

	return client, nil
}

func (c *Config) loadAndValidate() error {
	if err := c.load(); err != nil {
		return err
	}

	// If OVHClient is nil (because we ignored client creation error),
	// skip the authentication validation step entirely.
	if c.IgnoreInitError && c.OVHClient == nil {
		log.Printf("[WARN] OVH client is nil, skipping authentication validation")
		return nil
	}

	c.lockAuth.Lock()
	defer c.lockAuth.Unlock()

	if c.authFailed != nil {
		return c.authFailed
	}

	if !c.authenticated {
		// If OVHClient is nil (OAuth error was ignored during load), skip validation
		if c.OVHClient == nil {
			if c.IgnoreInitError {
				log.Printf("[WARN] OVH client not initialized, skipping authentication validation")
				return nil
			}
			return fmt.Errorf("OVH client not initialized")
		}

		var details OvhAuthDetails
		if err := c.OVHClient.Get("/auth/details", &details); err != nil {
			// Allow ignoring the /auth/details verification step when OVH_IGNORE_INIT_ERROR=true
			// and using OAuth2 (ClientID provided).
			if c.IgnoreInitError {
				log.Printf("[WARN] Ignoring /auth/details verification error: %v", err)
				// Leave c.authenticated=false so runtime calls may retry auth
				return nil
			}
			c.authFailed = fmt.Errorf("OVH client seems to be misconfigured: %q", err)
			return c.authFailed
		}

		log.Printf("[DEBUG] Logged in on OVH API")
		c.Account = details.Account
		c.authenticated = true
	}

	if c.Plate == "" {
		c.Plate = plateFromEndpoint(c.Endpoint)
	}

	return nil
}

func (c *Config) load() error {
	targetClient, err := clientDefault(c)
	if err != nil {
		// Allow ignoring client creation errors when OVH_IGNORE_INIT_ERROR=true
		if c.IgnoreInitError {
			log.Printf("[WARN] Ignoring client creation error: %v", err)
			return nil
		}
		return fmt.Errorf("error getting ovh client: %q", err)
	}

	// decorating the OVH http client with logs
	httpClient := targetClient.Client
	if targetClient.Client.Transport == nil {
		targetClient.Client.Transport = cleanhttp.DefaultTransport()
	}

	httpClient.Transport = logging.NewTransport("OVH", httpClient.Transport)
	c.OVHClient = ovhwrap.NewClient(targetClient, c.ApiRateLimit)
	return nil
}

var plateMapping map[string]string = map[string]string{
	"ovh-eu":        "eu",
	"ovh-ca":        "ca",
	"ovh-us":        "us",
	"kimsufi-eu":    "eu",
	"kimsufi-ca":    "ca",
	"soyoustart-eu": "eu",
	"soyoustart-ca": "ca",
}

func plateFromEndpoint(endpoint string) string {
	return plateMapping[endpoint]
}
