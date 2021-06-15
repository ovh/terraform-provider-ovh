package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMeSshKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMeSshKeyRead,
		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this public Ssh key",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ASCII encoded public Ssh key",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True when this public Ssh key is used for rescue mode and reinstallations",
			},
		},
	}
}

func dataSourceMeSshKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sshKey := &MeSshKeyResponse{}

	keyName := d.Get("key_name").(string)
	err := config.OVHClient.Get(
		fmt.Sprintf("/me/sshKey/%s", keyName),
		sshKey,
	)
	if err != nil {
		return fmt.Errorf("Unable to find SSH key named %s:\n\t %q", keyName, err)
	}

	d.SetId(sshKey.KeyName)
	d.Set("key_name", sshKey.KeyName)
	d.Set("key", sshKey.Key)
	d.Set("default", sshKey.Default)

	return nil
}
