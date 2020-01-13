package ovh

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description: "Farm id your vRack network is attached to and their type",
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
					err := validateIpBlock(v.(string))
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
					err := validateIpBlock(v.(string))
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
		"/ipLoadbalancing/%s/vrack/network",
		url.PathEscape(serviceName),
	)

	// start of update procedure (put + post)
	d.Partial(true)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	d.SetPartial("nat_ip")
	d.SetPartial("display_name")
	d.SetPartial("vlan")

	// update farm id
	farmIdOpts := (&IpLoadbalancingVrackNetworkFarmIdUpdateOpts{}).FromResource(d)
	vrackNetwork := &IpLoadbalancingVrackNetwork{}

	endpoint = fmt.Sprintf("%s/updateFarmId", endpoint)
	if err := config.OVHClient.Post(endpoint, farmIdOpts, vrackNetwork); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, farmIdOpts, err)
	}
	d.SetPartial("farm_id")

	// end of update procedure
	d.Partial(false)

	return resourceIPLoadbalancingVrackNetworkRead(d, meta)
}

func resourceIPLoadbalancingVrackNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	// start of delete procedure (put + post)
	d.Partial(true)

	// update farm id to remove all farm ids from vrack network
	farmIdOpts := &IpLoadbalancingVrackNetworkFarmIdUpdateOpts{
		FarmId: []int64{},
	}

	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network/%d/updateFarmId",
		url.PathEscape(serviceName),
		d.Get("vrack_network_id").(int),
	)
	if err := config.OVHClient.Post(endpoint, farmIdOpts, nil); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, farmIdOpts, err)
	}

	d.Set("farm_id", nil)
	d.SetPartial("farm_id")

	// delete network
	endpoint = fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network/%d",
		url.PathEscape(serviceName),
		d.Get("vrack_network_id").(int),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Error calling DELETE %s: %s \n", endpoint, err.Error())
	}
	// end of update procedure
	d.Partial(false)

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
	if err := config.OVHClient.Get(endpoint, &vn); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
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
