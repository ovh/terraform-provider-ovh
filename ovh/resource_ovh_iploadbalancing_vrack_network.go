package ovh

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceIPLoadbalancingVrackNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingVrackNetworkCreate,
		Read:   resourceIPLoadbalancingVrackNetworkRead,
		Update: resourceIPLoadbalancingVrackNetworkUpdate,
		Delete: resourceIPLoadbalancingVrackNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIPLoadbalancingVrackNetworkImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your IPloadbalancer",
				Required:    true,
			},

			"farm_id": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "This attribute is there for documentation purpose only and isnt passed to the OVH API as it may conflicts with http/tcp farms `vrack_network_id` attribute",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Human readable name for your vrack network",
				Optional:    true,
			},
			"nat_ip": {
				Type:        schema.TypeString,
				Description: "An IP block used as a pool of IPs by this Load Balancer to connect to the servers in this private network. The blck must be in the private network and reserved for the Load Balancer",
				Required:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "IP block of the private network in the vRack",
				Required:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"vlan": {
				Type:        schema.TypeInt,
				Description: "VLAN of the private network in the vRack. 0 if the private network is not in a VLAN",
				Optional:    true,
				Computed:    true,
			},

			//Computed
			"vrack_network_id": {
				Type:        schema.TypeInt,
				Description: "Internal Load Balancer identifier of the vRack private network",
				Computed:    true,
			},
		},
	}
}

func resourceIPLoadbalancingVrackNetworkImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/vrack_network_id formatted")
	}
	serviceName := splitId[0]
	vrackNetworkId := splitId[1]
	d.SetId(fmt.Sprintf("%s_%s", serviceName, vrackNetworkId))
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceIPLoadbalancingVrackNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	opts := (&IpLoadbalancingVrackNetworkCreateOpts{}).FromResource(d)
	vrackNetwork := &IpLoadbalancingVrackNetwork{}

	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Post(endpoint, opts, vrackNetwork); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}
	d.SetId(fmt.Sprintf("%s_%d", serviceName, vrackNetwork.VrackNetworkId))

	// set resource attributes
	for k, v := range vrackNetwork.ToMap() {
		d.Set(k, v)
	}

	return resourceIPLoadbalancingVrackNetworkRead(d, meta)
}

func resourceIPLoadbalancingVrackNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	opts := (&IpLoadbalancingVrackNetworkUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network/%d",
		url.PathEscape(serviceName),
		d.Get("vrack_network_id").(int),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	return resourceIPLoadbalancingVrackNetworkRead(d, meta)
}

func resourceIPLoadbalancingVrackNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	// delete network
	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network/%d",
		url.PathEscape(serviceName),
		d.Get("vrack_network_id").(int),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Error calling DELETE %s: %s \n", endpoint, err.Error())
	}

	d.SetId("")
	return nil
}

func resourceIPLoadbalancingVrackNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	networkId, err := strconv.ParseInt(strings.TrimPrefix(d.Id(), fmt.Sprintf("%s_", serviceName)), 10, 64)
	if err != nil {
		return fmt.Errorf(
			"Could not parse iploadbalancing vrack network id %s,%s:\n\t %q",
			serviceName,
			d.Id(),
			err,
		)
	}

	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network/%d",
		url.PathEscape(d.Get("service_name").(string)),
		networkId,
	)

	vn := &IpLoadbalancingVrackNetwork{}
	if err := config.OVHClient.Get(endpoint, vn); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	if networkId != vn.VrackNetworkId {
		return fmt.Errorf(
			"Network Id inconsistency for iploadbalancing %s. asked %d, got %d",
			serviceName,
			networkId,
			vn.VrackNetworkId,
		)
	}

	// set resource attributes
	for k, v := range vn.ToMap() {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprintf("%s_%d", d.Get("service_name").(string), d.Get("vrack_network_id").(int)))
	return nil
}
