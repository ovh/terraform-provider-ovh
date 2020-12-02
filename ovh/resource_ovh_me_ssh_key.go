package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceMeSshKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeSshKeyCreate,
		Read:   resourceMeSshKeyRead,
		Update: resourceMeSshKeyUpdate,
		Delete: resourceMeSshKeyDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of this public Ssh key",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ASCII encoded public Ssh key",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "True when this public Ssh key is used for rescue mode and reinstallations",
			},
		},
	}
}

// Common function with the datasource
func resourceMeSshKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	r := &MeSshKeyResponse{}
	endpoint := fmt.Sprintf("/me/sshKey/%s", d.Id())

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("key_name", r.KeyName)
	d.Set("key", r.Key)
	d.Set("default", r.Default)

	return nil
}

func resourceMeSshKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyName := d.Get("key_name").(string)
	key := d.Get("key").(string)
	params := &MeSshKeyCreateOpts{
		KeyName: keyName,
		Key:     key,
	}

	log.Printf("[DEBUG] Will create Ssh key: %s", params)

	err := config.OVHClient.Post("/me/sshKey", params, nil)
	if err != nil {
		return fmt.Errorf("Error creating SSH Key with params %s:\n\t %q", params, err)
	}

	d.SetId(keyName)

	// Update the resource in all cases in order to set Default if it is
	// different from default value (false)
	putParams := &MeSshKeyUpdateOpts{
		Default: d.Get("default").(bool),
	}
	err = config.OVHClient.Put(
		fmt.Sprintf("/me/sshKey/%s", keyName),
		putParams,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Unable to update SSH key named %s:\n\t %q", keyName, err)
	}

	return resourceMeSshKeyRead(d, meta)
}

func resourceMeSshKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyName := d.Get("key_name").(string)
	params := &MeSshKeyUpdateOpts{
		Default: d.Get("default").(bool),
	}
	err := config.OVHClient.Put(
		fmt.Sprintf("/me/sshKey/%s", keyName),
		params,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Unable to update SSH key named %s:\n\t %q", keyName, err)
	}

	log.Printf("[DEBUG] Updated SSH Key %s", keyName)
	return resourceMeSshKeyRead(d, meta)
}

func resourceMeSshKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyName := d.Get("key_name").(string)
	err := config.OVHClient.Delete(
		fmt.Sprintf("/me/sshKey/%s", keyName),
		nil,
	)
	if err != nil {
		return fmt.Errorf("Unable to delete SSH key named %s:\n\t %q", keyName, err)
	}

	log.Printf("[DEBUG] Deleted SSH Key %s", keyName)
	d.SetId("")
	return nil
}
