package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedConfigForRegion returns a common Config setup needed for the sweeper
// functions for a given region
func sharedConfigForRegion(region string) (*Config, error) {
	v := os.Getenv("OVH_ENDPOINT")
	if v == "" {
		return nil, fmt.Errorf("OVH_ENDPOINT must be set")
	}

	v = os.Getenv("OVH_APPLICATION_KEY")
	if v == "" {
		return nil, fmt.Errorf("OVH_APPLICATION_KEY must be set")
	}

	v = os.Getenv("OVH_APPLICATION_SECRET")
	if v == "" {
		return nil, fmt.Errorf("OVH_APPLICATION_SECRET must be set")
	}

	v = os.Getenv("OVH_CONSUMER_KEY")
	if v == "" {
		return nil, fmt.Errorf("OVH_CONSUMER_KEY must be set")
	}

	config := Config{
		Endpoint:          os.Getenv("OVH_ENDPOINT"),
		ApplicationKey:    os.Getenv("OVH_APPLICATION_KEY"),
		ApplicationSecret: os.Getenv("OVH_APPLICATION_SECRET"),
		ConsumerKey:       os.Getenv("OVH_CONSUMER_KEY"),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, fmt.Errorf("couln't load OVH Client: %s", err)

	}

	return &config, nil
}

// sharedClientForRegion returns a common OVHClient setup needed for the sweeper
// functions for a given region
func sharedClientForRegion(region string) (*ovh.Client, error) {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return nil, err
	}
	return config.OVHClient, nil
}
