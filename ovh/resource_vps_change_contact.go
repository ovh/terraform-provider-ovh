package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPSChangeContact() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPSChangeContactCreate,
		ReadContext:   resourceVPSChangeContactRead,
		DeleteContext: resourceVPSChangeContactDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service_name of the VPS",
				Required:    true,
				ForceNew:    true,
			},
			"contact_admin": {
				Type:        schema.TypeString,
				Description: "The OVH nichandle of the new admin contact",
				Optional:    true,
				ForceNew:    true,
			},
			"contact_billing": {
				Type:        schema.TypeString,
				Description: "The OVH nichandle of the new billing contact",
				Optional:    true,
				ForceNew:    true,
			},
			"contact_tech": {
				Type:        schema.TypeString,
				Description: "The OVH nichandle of the new tech contact",
				Optional:    true,
				ForceNew:    true,
			},
			"task_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceVPSChangeContactCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	opts := &VPSChangeContactOpts{}
	count := 0
	if v, ok := d.GetOk("contact_admin"); ok {
		s := v.(string)
		opts.ContactAdmin = &s
		count++
	}
	if v, ok := d.GetOk("contact_billing"); ok {
		s := v.(string)
		opts.ContactBilling = &s
		count++
	}
	if v, ok := d.GetOk("contact_tech"); ok {
		s := v.(string)
		opts.ContactTech = &s
		count++
	}

	if count == 0 {
		return diag.Errorf("at least one of contact_admin, contact_billing or contact_tech must be set")
	}

	var taskIDs []int64
	endpoint := fmt.Sprintf("/vps/%s/changeContact", url.PathEscape(serviceName))
	if err := config.OVHClient.Post(endpoint, opts, &taskIDs); err != nil {
		return diag.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	ids := make([]int, 0, len(taskIDs))
	idStrings := make([]string, 0, len(taskIDs))
	for _, id := range taskIDs {
		ids = append(ids, int(id))
		idStrings = append(idStrings, strconv.FormatInt(id, 10))
	}
	d.Set("task_ids", ids)

	// Build a stable id, unique per change request, since this resource is a
	// one-shot action and not state-bearing on the OVH side.
	id := fmt.Sprintf("%s-%d", serviceName, time.Now().UnixNano())
	if len(idStrings) > 0 {
		id = fmt.Sprintf("%s-%s", serviceName, strings.Join(idStrings, "-"))
	}
	d.SetId(id)

	return nil
}

func resourceVPSChangeContactRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// The OVH API does not expose a way to read back a previously submitted
	// change-contact request, so there is nothing to refresh here. The new
	// contacts can be observed through ovh_vps_service_info.
	return nil
}

func resourceVPSChangeContactDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Contact changes are not reversible from this resource. Removing the
	// resource only drops it from state.
	d.SetId("")
	return nil
}
