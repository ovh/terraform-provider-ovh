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

	RawOVHClient  *ovh.Client
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

	c.lockAuth.Lock()
	defer c.lockAuth.Unlock()

	if c.authFailed != nil {
		return c.authFailed
	}

	if !c.authenticated {
		var details OvhAuthDetails
		if err := c.OVHClient.Get("/auth/details", &details); err != nil {
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
		return fmt.Errorf("error getting ovh client: %q", err)
	}

	// decorating the OVH http client with logs
	httpClient := targetClient.Client
	if targetClient.Client.Transport == nil {
		targetClient.Client.Transport = cleanhttp.DefaultTransport()
	}

	httpClient.Transport = logging.NewTransport("OVH", httpClient.Transport)
	c.RawOVHClient = targetClient
	c.OVHClient = &ovhwrap.Client{
		RawClient:   targetClient,
		RateLimiter: c.ApiRateLimit,
	}

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
