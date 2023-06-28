package ovh

import (
	"fmt"
	"log"
	"sync"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/ovh/go-ovh/ovh"
)

var providerVersion, providerCommit string

type Config struct {
	Account           string
	Plate             string
	Endpoint          string
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
	OVHClient         *ovh.Client
	authenticated     bool
	authFailed        error
	lockAuth          *sync.Mutex
}

func clientDefault(c *Config) (*ovh.Client, error) {
	client, err := ovh.NewClient(
		c.Endpoint,
		c.ApplicationKey,
		c.ApplicationSecret,
		c.ConsumerKey,
	)
	if err != nil {
		return nil, err
	}

	client.UserAgent = "Terraform/" + providerVersion + "/" + providerCommit
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
			c.authFailed = fmt.Errorf("OVH client seems to be misconfigured: %q\n", err)
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
	validEndpoint := false

	ovhEndpoints := [7]string{ovh.OvhEU, ovh.OvhCA, ovh.OvhUS, ovh.KimsufiEU, ovh.KimsufiCA, ovh.SoyoustartEU, ovh.SoyoustartCA}

	for _, e := range ovhEndpoints {
		if ovh.Endpoints[c.Endpoint] == e {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("%s must be one of %#v endpoints\n", c.Endpoint, ovh.Endpoints)
	}

	targetClient, err := clientDefault(c)
	if err != nil {
		return fmt.Errorf("Error getting ovh client: %q\n", err)
	}

	// decorating the OVH http client with logs
	httpClient := targetClient.Client
	if targetClient.Client.Transport == nil {
		targetClient.Client.Transport = cleanhttp.DefaultTransport()
	}

	httpClient.Transport = logging.NewTransport("OVH", httpClient.Transport)
	c.OVHClient = targetClient

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
