package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	sshKey := &MeSshKeyResponse{}

	id := d.Id()
	err := config.OVHClient.Get(
		fmt.Sprintf("/me/sshKey/%s", id),
		sshKey,
	)
	if err != nil {
		return fmt.Errorf("Unable to find SSH key named %s:\n\t %q", id, err)
	}

	d.Set("key_name", sshKey.KeyName)
	d.Set("key", sshKey.Key)
	d.Set("default", sshKey.Default)

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
		nil,
		putParams,
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
		nil,
		params,
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
