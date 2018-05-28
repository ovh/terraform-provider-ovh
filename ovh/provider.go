package ovh

import (
	"fmt"
	"log"
	"os"
	"os/user"

	ini "gopkg.in/ini.v1"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a schema.Provider for OVH.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_ENDPOINT", nil),
				Description: descriptions["endpoint"],
			},
			"application_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_KEY", ""),
				Description: descriptions["application_key"],
			},
			"application_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_SECRET", ""),
				Description: descriptions["application_secret"],
			},
			"consumer_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CONSUMER_KEY", ""),
				Description: descriptions["consumer_key"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			// New naming schema (issue #23)
			"ovh_cloud_region":  dataSourcePublicCloudRegion(),
			"ovh_cloud_regions": dataSourcePublicCloudRegions(),
			// Legacy naming schema (new datasources should not be added here)
			"ovh_publiccloud_region":  dataSourcePublicCloudRegion(),
			"ovh_publiccloud_regions": dataSourcePublicCloudRegions(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ovh_iploadbalancing_tcp_farm":        resourceIpLoadbalancingTcpFarm(),
			"ovh_iploadbalancing_tcp_farm_server": resourceIpLoadbalancingTcpFarmServer(),
			"ovh_domain_zone_record":              resourceOvhDomainZoneRecord(),
			"ovh_domain_zone_redirection":         resourceOvhDomainZoneRedirection(),
			// New naming schema (issue #23)
			"ovh_cloud_network_private":        resourcePublicCloudPrivateNetwork(),
			"ovh_cloud_network_private_subnet": resourcePublicCloudPrivateNetworkSubnet(),
			"ovh_cloud_user":                   resourcePublicCloudUser(),
			"ovh_vrack_cloudproject":           resourceVRackPublicCloudAttachment(),
			// Legacy naming schema (new resources should not be added here)
			"ovh_publiccloud_private_network":        resourcePublicCloudPrivateNetwork(),
			"ovh_publiccloud_private_network_subnet": resourcePublicCloudPrivateNetworkSubnet(),
			"ovh_publiccloud_user":                   resourcePublicCloudUser(),
			"ovh_vrack_publiccloud_attachment":       resourceVRackPublicCloudAttachment(),
		},

		ConfigureFunc: configureProvider,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"endpoint": "The OVH API endpoint to target (ex: \"ovh-eu\").",

		"application_key": "The OVH API Application Key.",

		"application_secret": "The OVH API Application Secret.",
		"consumer_key":       "The OVH API Consumer key.",
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	userHome, err := currentUserHome()
	if err != nil {
		log.Fatal(err)
	}
	config := Config{
		Endpoint: d.Get("endpoint").(string),
	}
	configFile := fmt.Sprintf("%s/.ovh.conf", userHome)
	if _, err := os.Stat(configFile); err == nil {
		c, err := ini.Load(configFile)
		if err != nil {
			return nil, err
		}

		section, err := c.GetSection(d.Get("endpoint").(string))
		if err != nil {
			return nil, err
		}
		config.ApplicationKey = section.Key("application_key").String()
		config.ApplicationSecret = section.Key("application_secret").String()
		config.ConsumerKey = section.Key("consumer_key").String()
	}
	if v, ok := d.GetOk("application_key"); ok {
		config.ApplicationKey = v.(string)
	}
	if v, ok := d.GetOk("application_secret"); ok {
		config.ApplicationSecret = v.(string)
	}
	if v, ok := d.GetOk("consumer_key"); ok {
		config.ConsumerKey = v.(string)
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// currentUserHome attempts to get current user's home directory
func currentUserHome() (string, error) {
	userHome := ""
	usr, err := user.Current()
	if err != nil {
		// Fallback by trying to read $HOME
		userHome = os.Getenv("HOME")
		if userHome != "" {
			err = nil
		}
	} else {
		userHome = usr.HomeDir
	}
	return userHome, nil
}
