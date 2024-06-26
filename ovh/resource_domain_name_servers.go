package ovh

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"log"
	"strings"
)

type OvhDomainNameServerUpdate struct {
	NameServers []OvhDomainNameServer `json:"nameServers,omitempty"`
}

type OvhDomainNameServerUpdateResult struct {
	CanAccelerate bool   `json:"canAccelerate,omitempty"`
	CanCancel     bool   `json:"canCancel,omitempty"`
	CanRelaunch   bool   `json:"canRelaunch,omitempty"`
	Comment       string `json:"comment,omitempty"`
	CreationDate  string `json:"creationDate,omitempty"`
	Domain        string `json:"domain,omitempty"`
	DoneDate      string `json:"doneDate,omitempty"`
	Function      string `json:"function,omitempty"`
	Id            int64  `json:"id,omitempty"`
	LastUpdate    string `json:"lastUpdate,omitempty"`
	Status        string `json:"status,omitempty"`
	TodoDate      string `json:"todoDate,omitempty"`
}

type OvhDomainNameServer struct {
	Id       int64  `json:"id,omitempty"`
	Host     string `json:"host,omitempty"`
	Ip       string `json:"ip,omitempty"`
	IsUsed   bool   `json:"isUsed,omitempty"`
	ToDelete bool   `json:"toDelete,omitempty"`
}

type OvhDomainNameServers struct {
	Id      string                `json:"id,omitempty"`
	Domain  string                `json:"domain,omitempty"`
	Servers []OvhDomainNameServer `json:"servers,omitempty"`
}

func (r *OvhDomainNameServers) String() string {
	domains := make([]string, 0)
	for _, server := range r.Servers {
		domains = append(domains, fmt.Sprintf(
			"server[id: %v, host: %s, ip: %s, isUsed: %v, toDelete: %v]",
			server.Id,
			server.Host,
			server.Ip,
			server.IsUsed,
			server.ToDelete,
		))
	}
	return fmt.Sprintf(
		"nameservers[id: %v, domain: %s, servers: [%s]]",
		r.Id,
		r.Domain,
		strings.Join(domains, ", "),
	)
}

func resourceOvhDomainNameServersImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	d.SetId(givenId)
	d.Set("domain", d.Id())
	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceOvhDomainNameServers() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvhDomainNameServersCreate,
		Read:   resourceOvhDomainNameServersRead,
		Update: resourceOvhDomainNameServersCreate, // Update is the same as create
		Delete: resourceOvhDomainNameServersDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOvhDomainNameServersImportState,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
				MinItems: 1,
			},
		},
	}
}

func resourceOvhDomainNameServersCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("domain").(string)

	// Servers to put
	servers := &OvhDomainNameServerUpdate{
		NameServers: make([]OvhDomainNameServer, 0),
	}
	// Loop due to rename of the field in the API
	for _, server := range d.Get("servers").([]interface{}) {
		s := server.(map[string]interface{})
		servers.NameServers = append(servers.NameServers, OvhDomainNameServer{
			Host: s["host"].(string),
			Ip:   s["ip"].(string),
		})
	}
	log.Printf("[DEBUG] OVH Record create configuration: %#v", servers)
	err := config.OVHClient.Post(fmt.Sprintf("/domain/%s/nameServers/update", serviceName), servers, nil)

	if err != nil {
		return fmt.Errorf("failed to register OVH Nameservers: %s", err)
	}

	d.SetId(serviceName)

	return resourceOvhDomainNameServersRead(d, meta)
}

func resourceOvhDomainNameServersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	record, err := ovhDomainNameServers(config.OVHClient, d)

	if err != nil {
		return err
	}

	if record == nil {
		return fmt.Errorf("domain %v has been deleted", d.Id())
	}

	d.SetId(record.Id)
	d.Set("domain", record.Domain)
	d.Set("servers", flattenOvhDomainNameServers(record.Servers))

	return nil
}

func flattenOvhDomainNameServers(servers []OvhDomainNameServer) []interface{} {
	result := make([]interface{}, 0)
	for _, server := range servers {
		result = append(result, map[string]interface{}{
			"host": server.Host,
			"ip":   server.Ip,
		})
	}
	return result
}

func resourceOvhDomainNameServersDelete(d *schema.ResourceData, meta interface{}) error {
	// Note that NameServers can not be deleted, only updated
	d.SetId("")
	return nil
}

func ovhDomainNameServers(client *ovh.Client, d *schema.ResourceData) (*OvhDomainNameServers, error) {
	domain := d.Get("domain").(string)
	nameServers := &OvhDomainNameServers{
		Servers: make([]OvhDomainNameServer, 0),
		Domain:  domain,
		Id:      domain,
	}
	rec := &([]int{})

	endpoint := fmt.Sprintf("/domain/%s/nameServer", domain)

	err := client.Get(endpoint, rec)

	if err != nil {
		if err.(*ovh.APIError).Code == 404 {
			return nil, nil
		}
		return nil, err
	}

	// Read each name server
	for _, id := range *rec {
		server, err := ovhDomainNameServer(client, domain, id)
		if err != nil {
			return nil, err
		}
		nameServers.Servers = append(nameServers.Servers, *server)
	}

	return nameServers, nil
}

func ovhDomainNameServer(client *ovh.Client, domain string, id int) (*OvhDomainNameServer, error) {
	rec := &OvhDomainNameServer{}

	endpoint := fmt.Sprintf("/domain/%s/nameServer/%d", domain, id)
	err := client.Get(endpoint, rec)
	if err != nil {
		return nil, err
	}

	return rec, nil
}
