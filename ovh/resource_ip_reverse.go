package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func resourceIpReverse() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpReverseCreate,
		Read:   resourceIpReverseRead,
		Delete: resourceIpReverseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpReverseImportState,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			"ip_reverse": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIp(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			"reverse": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"readiness_timeout_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				ForceNew: true,
			},
		},
	}
}

func resourceIpReverseImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "|", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not ip|ip_reverse formatted")
	}
	ip := splitId[0]
	ipReverse := splitId[1]
	d.SetId(ipReverse)
	d.Set("ip", ip)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceIpReverseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Create the new reverse
	ip := d.Get("ip").(string)
	opts := (&IpReverseCreateOpts{}).FromResource(d)
	res := &IpReverse{}

	readinessTimeoutDurationAttr, isReadinessTimeoutDefined := d.GetOk("readiness_timeout_duration")

	retryDuration := 1 * time.Minute
	if isReadinessTimeoutDefined {
		var err error
		retryDuration, err = time.ParseDuration(readinessTimeoutDurationAttr.(string))
		if err != nil {
			return fmt.Errorf("failed to create OVH IP Reverse: cannot parse readiness_timeout_seconds attribute: %s", err)
		}
	}

	err := postIpReverseWithRetry(context.TODO(), *config.OVHClient, fmt.Sprintf("/ip/%s/reverse", url.PathEscape(ip)), opts, res, retryDuration)
	if err != nil {
		return fmt.Errorf("failed to create OVH IP Reverse: %s", err)
	}

	d.SetId(res.IpReverse)

	return resourceIpReverseRead(d, meta)
}

func postIpReverseWithRetry(parentCtx context.Context, client ovhwrap.Client, endpoint string, opts *IpReverseCreateOpts, result *IpReverse, retryDuration time.Duration) error {
	ctx, cancel := context.WithTimeout(parentCtx, retryDuration)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		err := client.Post(
			endpoint,
			opts,
			&result,
		)

		if err == nil {
			return err
		}

		errOvh, ok := err.(*ovh.APIError)
		if ok && errOvh.Code != 400 {
			// don't retry non-400 errors
			return err
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry duration exhausted for ovh_ip_reverse resource creation: %w", err)
		case <-ticker.C:
			continue
		}
	}
}

func resourceIpReverseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ip := d.Get("ip").(string)

	res := &IpReverse{}
	endpoint := fmt.Sprintf(
		"/ip/%s/reverse/%s",
		url.PathEscape(ip),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range res.ToMap() {
		d.Set(k, v)
	}

	return nil
}

func resourceIpReverseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] Deleting OVH IP Reverse: %s->%s", d.Get("reverse").(string), d.Get("ip_reverse").(string))
	ip := d.Get("ip").(string)
	endpoint := fmt.Sprintf(
		"/ip/%s/reverse/%s",
		url.PathEscape(ip),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId("")
	return nil
}
