package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func resourceDomainNameServers() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a domain name servers",
		Schema:      resourceDomainNameServersSchema(),
		Create:      resourceDomainNameServersUpdate,
		Read:        resourceDomainNameServersRead,
		Update:      resourceDomainNameServersUpdate,
		Delete:      resourceDomainNameServersDelete,
		Importer: &schema.ResourceImporter{
			State: func(resourceData *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				resourceData.Set("domain", resourceData.Id())
				return []*schema.ResourceData{resourceData}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(10 * time.Minute),
			Read:    schema.DefaultTimeout(30 * time.Second),
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Update:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceDomainNameServersSchema() map[string]*schema.Schema {
	resourceSchema := map[string]*schema.Schema{
		"domain": {
			Type:        schema.TypeString,
			Description: "Domain name",
			Required:    true,
			ForceNew:    true,
		},
		"servers": {
			Type:        schema.TypeSet,
			Description: "Name servers for the domain",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host": {
						Type:        schema.TypeString,
						Description: "DNS name server hostname",
						ValidateFunc: func(v any, k string) (ws []string, err []error) {
							if strings.HasSuffix(v.(string), ".") {
								return nil, []error{
									fmt.Errorf(`field "host" must not end by a dot`),
								}
							}

							return nil, nil
						},
						Required: true,
					},
					"ip": {
						Type:        schema.TypeString,
						Description: "DNS name server IP address",
						ValidateFunc: func(v any, k string) (ws []string, errors []error) {
							if v.(string) == "" {
								return
							}

							err := helpers.ValidateIp(v.(string))

							if err != nil {
								errors = append(errors, err)
							}

							return
						},
						Optional: true,
						Default:  "",
					},
				},
			},
			Required: true,
			MinItems: 2,
			MaxItems: 8,
		},
	}

	return resourceSchema
}

func resourceDomainNameServersRead(resourceData *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	domainName := resourceData.Id()

	domainNameServers := &DomainNameServers{
		Domain: domainName,
	}

	log.Printf("[DEBUG] Will read domain name servers: %s\n", domainName)

	var nameservers []int
	endpoint := fmt.Sprintf("/domain/%s/nameServer", url.PathEscape(domainName))

	if err := config.OVHClient.Get(endpoint, &nameservers); err != nil {
		return fmt.Errorf("calling GET %s:\n\t%s", endpoint, err.Error())
	}

	for _, nameServerId := range nameservers {
		responseData := DomainNameServer{}
		endpoint := fmt.Sprintf("/domain/%s/nameServer/%d", url.PathEscape(domainName), nameServerId)

		if err := config.OVHClient.Get(endpoint, &responseData); err != nil {
			return helpers.CheckDeleted(resourceData, err, endpoint)
		}

		domainNameServers.Servers = append(domainNameServers.Servers, responseData)
	}

	resourceData.SetId(domainName)
	for k, v := range domainNameServers.ToMap() {
		resourceData.Set(k, v)
	}

	return nil
}

func resourceDomainNameServersUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	domainName := resourceData.Get("domain").(string)
	task := DomainTask{}

	nameServersUpdate := &DomainNameServerUpdateOpts{
		NameServers: make([]DomainNameServer, 0),
	}

	for _, nameServer := range resourceData.Get("servers").(*schema.Set).List() {
		ns := nameServer.(map[string]interface{})

		nameServersUpdate.NameServers = append(nameServersUpdate.NameServers, DomainNameServer{
			Host: ns["host"].(string),
			Ip:   ns["ip"].(string),
		})
	}

	log.Printf("[DEBUG] Will update domain name nameServersUpdate: %s\n", domainName)

	endpoint := fmt.Sprintf("/domain/%s/nameServers/update", url.PathEscape(domainName))

	if err := config.OVHClient.Post(endpoint, nameServersUpdate, &task); err != nil {
		return fmt.Errorf("calling POST %s:\n\t%s", endpoint, err.Error())
	}

	if err := waitDomainTask(config.OVHClient, domainName, task.TaskID); err != nil {
		return fmt.Errorf("waiting for %s name servers to be updated:\n\t%s", domainName, err.Error())
	}

	resourceData.SetId(domainName)

	return resourceDomainNameServersRead(resourceData, meta)
}

func resourceDomainNameServersDelete(resourceData *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	domainName := resourceData.Get("domain").(string)
	domainNameServerTypeOpts := &DomainNameServerTypeOpts{
		NameServerType: "hosted",
	}

	log.Printf("[DEBUG] Will reset domain name servers to default: %s\n", domainName)

	endpoint := fmt.Sprintf("/domain/%s", url.PathEscape(domainName))

	if err := config.OVHClient.Put(endpoint, domainNameServerTypeOpts, nil); err != nil {
		return fmt.Errorf("calling PUT %s:\n\t%s", endpoint, err.Error())
	}

	if err := waitDomainNameServersHosted(config.OVHClient, domainName); err != nil {
		return fmt.Errorf("waiting for %s name servers to be updated:\n\t%s", domainName, err.Error())
	}

	resourceData.SetId("")

	return nil
}

func waitDomainNameServersHosted(client *ovhwrap.Client, domainName string) error {
	endpoint := fmt.Sprintf("/domain/%s", url.PathEscape(domainName))

	stateConf := &retry.StateChangeConf{
		Pending: []string{"UPDATING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: func() (result interface{}, state string, err error) {
			var status = "UPDATING"
			var domainNameServerType DomainNameServerTypeOpts

			if err := client.Get(endpoint, &domainNameServerType); err != nil {
				log.Printf("[ERROR] couldn't fetch name server type for domain %s:\n\t%s\n", domainName, err.Error())
				return nil, "ERROR", err
			}

			if domainNameServerType.NameServerType != "external" {
				status = "SUCCEEDED"
			}

			return domainNameServerType, status, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err := stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error waiting for domain %s name server type to reset:\n\t%s", domainName, err.Error())
	}

	return err
}
