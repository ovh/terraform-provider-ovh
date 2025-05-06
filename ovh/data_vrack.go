package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVrack() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVrackRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Here come all the computed items
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceVrackRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	vrack := &Vrack{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/vrack/%s",
			serviceName,
		),
		&vrack,
	)

	if err != nil {
		return fmt.Errorf("Error calling /vrack/%s:\n\t %q", serviceName, err)
	}

	d.SetId(serviceName)
	d.Set("name", vrack.Name)
	d.Set("description", vrack.Description)
	iam := make(map[string]string)
	iam["id"] = vrack.IamResourceDetails.Id
	iam["urn"] = vrack.IamResourceDetails.URN
	d.Set("iam", iam)
	return nil
}
