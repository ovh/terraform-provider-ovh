package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceDbaasLogsInputEngine() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return dataSourceDbaasLogsInputEngineRead(d, meta)
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_deprecated": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDbaasLogsInputEngineRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[DEBUG] Will read dbaas logs input engines")
	res := []string{}
	endpoint := fmt.Sprintf("/dbaas/logs/input/engine")
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	nameFilter := helpers.GetNilStringPointerFromData(d, "name")
	isDeprecatedFilter := helpers.GetNilBoolPointerFromData(d, "is_deprecated")
	versionFilter := helpers.GetNilStringPointerFromData(d, "version")

	engines := []*DbaasLogsInputEngine{}

	for _, id := range res {
		log.Printf("[DEBUG] Will read dbaas logs input engine %s", id)
		endpoint := fmt.Sprintf("/dbaas/logs/input/engine/%s",
			url.PathEscape(id),
		)
		engine := &DbaasLogsInputEngine{}
		if err := config.OVHClient.Get(endpoint, &engine); err != nil {
			return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
		}

		match := true
		if nameFilter != nil && strings.ToLower(engine.Name) != strings.ToLower(*nameFilter) {
			match = false
		}
		if versionFilter != nil && strings.ToLower(engine.Version) != strings.ToLower(*versionFilter) {
			match = false
		}
		if isDeprecatedFilter != nil && engine.IsDeprecated != *isDeprecatedFilter {
			match = false
		}
		if match {
			engines = append(engines, engine)
		}
	}

	if len(engines) == 0 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(engines) > 1 {
		return fmt.Errorf("Your query returned more than one result. " +
			"Please change your search criteria and try again.")
	}

	for k, v := range engines[0].ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	return nil
}
