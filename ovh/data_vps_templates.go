package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

// vpsTemplate is the JSON shape returned by /vps/{sn}/templates/{id}.
// bitFormat comes back as a string ("32"/"64") in the API; we convert downstream.
type vpsTemplate struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Distribution      string   `json:"distribution"`
	BitFormat         string   `json:"bitFormat"`
	AvailableLanguage []string `json:"availableLanguage"`
	Locale            string   `json:"locale"`
}

func dataSourceVPSTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSTemplatesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"distribution_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Substring match against template.distribution",
			},
			"bit_format_filter": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Filter on template.bitFormat (32 or 64)",
				ValidateFunc: validateVPSBitFormat,
			},

			// Computed
			"template_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":           {Type: schema.TypeInt, Computed: true},
						"name":         {Type: schema.TypeString, Computed: true},
						"distribution": {Type: schema.TypeString, Computed: true},
						"bit_format":   {Type: schema.TypeInt, Computed: true},
						"available_language": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"locale": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func validateVPSBitFormat(v interface{}, k string) (ws []string, errors []error) {
	i, ok := v.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("%s must be an int", k))
		return
	}
	if i != 32 && i != 64 {
		errors = append(errors, fmt.Errorf("%s must be 32 or 64, got %d", k, i))
	}
	return
}

func dataSourceVPSTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	distFilter := strings.ToLower(d.Get("distribution_filter").(string))
	bitFilterRaw, bitFilterSet := d.GetOk("bit_format_filter")
	bitFilter := 0
	if bitFilterSet {
		bitFilter = bitFilterRaw.(int)
	}

	ids := []int{}
	listEndpoint := fmt.Sprintf("/vps/%s/templates", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(listEndpoint, &ids); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
			msg := apiErr.Message
			switch {
			case strings.Contains(msg, "Got an invalid (or empty) URL"):
				return fmt.Errorf(
					"the OVHcloud API endpoint %s is not available on this VPS lineup. "+
						"This data source may only work on legacy VPS plans, or the endpoint "+
						"may have been deprecated. See the data source's documentation for "+
						"supported VPS generations.",
					listEndpoint)
			case strings.Contains(msg, "does not exist"):
				return fmt.Errorf(
					"the requested resource at %s does not exist (the VPS may not have "+
						"the required option subscribed, or the resource ID is wrong)",
					listEndpoint)
			}
		}
		return fmt.Errorf("calling GET %s: %w", listEndpoint, err)
	}

	sort.Ints(ids)

	// Fan-out GET per template id; cap at 4 workers.
	type job struct {
		idx int
		id  int
	}
	type result struct {
		idx  int
		tmpl *vpsTemplate
		err  error
	}

	jobs := make(chan job)
	results := make(chan result)

	var wg sync.WaitGroup
	const workers = 4
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				t := &vpsTemplate{}
				ep := fmt.Sprintf(
					"/vps/%s/templates/%d",
					url.PathEscape(serviceName),
					j.id,
				)
				err := config.OVHClient.Get(ep, t)
				results <- result{idx: j.idx, tmpl: t, err: err}
			}
		}()
	}

	go func() {
		for i, id := range ids {
			jobs <- job{idx: i, id: id}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	collected := make([]*vpsTemplate, len(ids))
	for r := range results {
		if r.err != nil {
			if apiErr, ok := r.err.(*ovh.APIError); ok && apiErr.Code == 404 {
				msg := apiErr.Message
				switch {
				case strings.Contains(msg, "Got an invalid (or empty) URL"):
					return fmt.Errorf(
						"the OVHcloud API endpoint for VPS templates is not available on this VPS lineup. " +
							"This data source may only work on legacy VPS plans, or the endpoint " +
							"may have been deprecated. See the data source's documentation for " +
							"supported VPS generations.")
				case strings.Contains(msg, "does not exist"):
					return fmt.Errorf(
						"a requested VPS template does not exist (the VPS may not have " +
							"the required option subscribed, or the resource ID is wrong)")
				}
			}
			return fmt.Errorf("Error fetching VPS template: %w", r.err)
		}
		collected[r.idx] = r.tmpl
	}

	filteredIDs := []int{}
	filteredTemplates := []map[string]interface{}{}
	for _, t := range collected {
		if t == nil {
			continue
		}
		if distFilter != "" && !strings.Contains(strings.ToLower(t.Distribution), distFilter) {
			continue
		}
		bf, _ := strconv.Atoi(t.BitFormat)
		if bitFilterSet && bf != bitFilter {
			continue
		}
		filteredIDs = append(filteredIDs, t.ID)
		filteredTemplates = append(filteredTemplates, map[string]interface{}{
			"id":                 t.ID,
			"name":               t.Name,
			"distribution":       t.Distribution,
			"bit_format":         bf,
			"available_language": t.AvailableLanguage,
			"locale":             t.Locale,
		})
	}

	d.SetId(vpsTemplatesHashID(serviceName, filteredIDs))
	d.Set("template_ids", filteredIDs)
	d.Set("templates", filteredTemplates)
	return nil
}

func vpsTemplatesHashID(serviceName string, ids []int) string {
	parts := []string{serviceName}
	for _, i := range ids {
		parts = append(parts, strconv.Itoa(i))
	}
	return strings.Join(parts, "_")
}
