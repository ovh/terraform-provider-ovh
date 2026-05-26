package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

// vpsOption mirrors the OVH vps.Option payload returned by
// GET /vps/{serviceName}/option/{option}. We only expose `state` for now —
// the upstream schema is small and `state` is the only stable, useful field.
type vpsOption struct {
	State string `json:"state"`
}

func dataSourceVPSOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSOptionsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS.",
			},
			"options": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of option names currently subscribed on the VPS.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"options_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Per-option detail (state) fetched from /vps/{serviceName}/option/{option}.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVPSOptionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	names := []string{}
	endpoint := fmt.Sprintf("/vps/%s/option", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, &names); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	sort.Strings(names)

	// Fan-out option detail fetch capped at 4 concurrent workers.
	type detail struct {
		Name  string
		State string
		Err   error
	}
	details := make([]detail, len(names))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for i, name := range names {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, name string) {
			defer wg.Done()
			defer func() { <-sem }()
			opt := &vpsOption{}
			ep := fmt.Sprintf("/vps/%s/option/%s",
				url.PathEscape(serviceName),
				url.PathEscape(name),
			)
			if err := config.OVHClient.Get(ep, opt); err != nil {
				details[i] = detail{Name: name, Err: fmt.Errorf("Error calling GET %s:\n\t %q", ep, err)}
				return
			}
			details[i] = detail{Name: name, State: opt.State}
		}(i, name)
	}
	wg.Wait()

	optionsDetail := make([]map[string]interface{}, 0, len(details))
	for _, dt := range details {
		if dt.Err != nil {
			return dt.Err
		}
		optionsDetail = append(optionsDetail, map[string]interface{}{
			"name":  dt.Name,
			"state": dt.State,
		})
	}

	d.SetId(hashcode.Strings(names))
	d.Set("options", names)
	d.Set("options_detail", optionsDetail)
	return nil
}
