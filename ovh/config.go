package ovh

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"sync"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
	"gopkg.in/ini.v1"
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

	// Profile name to load from .ovh.conf
	Profile string

	OVHClient     *ovhwrap.Client
	authenticated bool
	authFailed    error
	lockAuth      *sync.Mutex

	ApiRateLimit ratelimit.Limiter
}

// ovhConfPaths holds the default paths to look for OVH configuration files.
var ovhConfPaths = []string{
	"/etc/ovh.conf",
	"~/.ovh.conf",
	"./ovh.conf",
}

// loadOvhConf loads the OVH INI configuration from the standard configuration file paths.
func loadOvhConf() (*ini.File, error) {
	paths := []interface{}{}

	var home string

	for _, path := range ovhConfPaths {
		if strings.HasPrefix(path, "~/") {
			if home == "" {
				if usr, err := user.Current(); err == nil {
					home = usr.HomeDir
				} else if h := os.Getenv("HOME"); h != "" {
					home = h
				} else {
					continue
				}
			}
			path = home + path[1:]
		}
		paths = append(paths, path)
	}

	if len(paths) == 0 {
		return ini.Empty(), nil
	}

	return ini.LooseLoad(paths[0], paths[1:]...)
}

// applyProfile loads configuration values from a named profile in the OVH configuration
// files and applies them to the Config struct. Profile sections use the format [profile:<name>].
//
// Values are only applied if they are not already set in the Config struct (i.e., not
// provided via Terraform HCL). Environment variables take precedence over profile values;
// if an env var is set for a field, the profile value is skipped so that go-ovh's own
// config loading can pick up the env var.
func (c *Config) applyProfile() error {
	if c.Profile == "" {
		return nil
	}

	cfg, err := loadOvhConf()
	if err != nil {
		return fmt.Errorf("cannot load OVH configuration file for profile %q: %w", c.Profile, err)
	}

	sectionName := "profile:" + c.Profile
	if !cfg.HasSection(sectionName) {
		return fmt.Errorf("profile %q not found in OVH configuration files", c.Profile)
	}

	section := cfg.Section(sectionName)

	// getVal returns the value of the given key from the profile section,
	// but returns an empty string if the corresponding env variable is already
	// set (so that go-ovh's own loadConfig can handle the env var with higher
	// priority than the profile).
	getVal := func(key, envVar string) string {
		if os.Getenv(envVar) != "" {
			return ""
		}
		if section.HasKey(key) {
			return section.Key(key).String()
		}
		return ""
	}

	if c.Endpoint == "" {
		c.Endpoint = getVal("endpoint", "OVH_ENDPOINT")
	}
	if c.AccessToken == "" {
		c.AccessToken = getVal("access_token", "OVH_ACCESS_TOKEN")
	}
	if c.ApplicationKey == "" {
		c.ApplicationKey = getVal("application_key", "OVH_APPLICATION_KEY")
	}
	if c.ApplicationSecret == "" {
		c.ApplicationSecret = getVal("application_secret", "OVH_APPLICATION_SECRET")
	}
	if c.ConsumerKey == "" {
		c.ConsumerKey = getVal("consumer_key", "OVH_CONSUMER_KEY")
	}
	if c.ClientID == "" {
		c.ClientID = getVal("client_id", "OVH_CLIENT_ID")
	}
	if c.ClientSecret == "" {
		c.ClientSecret = getVal("client_secret", "OVH_CLIENT_SECRET")
	}

	return nil
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
	// Apply profile configuration before creating the client, so that profile
	// values are available to clientDefault.
	if err := c.applyProfile(); err != nil {
		return err
	}

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

	// Chain transports: schemasVersion (adds X-Schemas-Version for /v2/ paths)
	// → logging (logs request/response) → original transport (sends over the wire).
	httpClient.Transport = newSchemasVersionTransport(logging.NewTransport("OVH", httpClient.Transport))
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
