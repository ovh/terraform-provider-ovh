package ovh

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sort"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

// vpsAvailableImagesFetchWorkers caps concurrent per-image GETs when
// populating the bulk `images` attribute. Kept small to be gentle on
// the OVH API while still being meaningfully faster than sequential.
const vpsAvailableImagesFetchWorkers = 4

func dataSourceVPSAvailableImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSAvailableImagesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_pattern": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   {Type: schema.TypeString, Computed: true},
						"name": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceVPSAvailableImagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	var pattern *regexp.Regexp
	if raw, ok := d.GetOk("name_pattern"); ok {
		p, err := regexp.Compile(raw.(string))
		if err != nil {
			return fmt.Errorf("invalid name_pattern %q: %w", raw, err)
		}
		pattern = p
	}

	ids := []string{}
	listEndpoint := fmt.Sprintf("/vps/%s/images/available", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(listEndpoint, &ids); err != nil {
		return fmt.Errorf("Error calling GET %s: %w", listEndpoint, err)
	}
	sort.Strings(ids)

	// Fan out per-id GETs with a small worker pool so we can both
	// filter on name and expose the bulk `images` field. Failures
	// on individual images are logged and skipped rather than
	// failing the whole data source.
	type result struct {
		idx int
		img VPSImage
		ok  bool
	}

	results := make([]result, len(ids))
	jobs := make(chan int)
	var wg sync.WaitGroup

	workers := vpsAvailableImagesFetchWorkers
	if len(ids) < workers {
		workers = len(ids)
	}
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				img := VPSImage{}
				endpoint := fmt.Sprintf(
					"/vps/%s/images/available/%s",
					url.PathEscape(serviceName),
					url.PathEscape(ids[i]),
				)
				if err := config.OVHClient.Get(endpoint, &img); err != nil {
					log.Printf("[WARN] skipping VPS image %q: %s", ids[i], err)
					results[i] = result{idx: i, ok: false}
					continue
				}
				// API may not echo id back; fall back to listing value.
				if img.ID == "" {
					img.ID = ids[i]
				}
				results[i] = result{idx: i, img: img, ok: true}
			}
		}()
	}
	for i := range ids {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	filteredIDs := make([]string, 0, len(ids))
	images := make([]map[string]interface{}, 0, len(ids))
	for _, r := range results {
		if !r.ok {
			continue
		}
		if pattern != nil && !pattern.MatchString(r.img.Name) {
			continue
		}
		filteredIDs = append(filteredIDs, r.img.ID)
		images = append(images, map[string]interface{}{
			"id":   r.img.ID,
			"name": r.img.Name,
		})
	}

	d.SetId(hashcode.Strings(filteredIDs))
	d.Set("image_ids", filteredIDs)
	d.Set("images", images)
	return nil
}
