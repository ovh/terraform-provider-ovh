package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VPSReinstallCreateOpts is the request body for POST /vps/{serviceName}/reinstall.
//
// This wraps the LEGACY template-based reinstall path (numeric templateId from
// /vps/templates), which is distinct from the newer image-based /vps/{sn}/rebuild
// flow exposed by ovh_vps.image_id.
type VPSReinstallCreateOpts struct {
	TemplateId        int64    `json:"templateId"`
	Language          string   `json:"language,omitempty"`
	SoftwareId        []int64  `json:"softwareId,omitempty"`
	SshKey            []string `json:"sshKey,omitempty"`
	PublicSshKey      string   `json:"publicSshKey,omitempty"`
	DoNotSendPassword bool     `json:"doNotSendPassword,omitempty"`
}

func resourceVPSReinstall() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSReinstallCreate,
		Read:   resourceVPSReinstallRead,
		Delete: resourceVPSReinstallDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS.",
			},
			"template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Legacy numeric template id to install. Discoverable via the ovh_vps_templates data sources.",
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "en",
				Description: "Display language for the installation (default: \"en\").",
			},
			"software_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional list of legacy software ids to install on top of the template.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of SSH key names (from /me/sshKey) to install on the VPS.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ssh_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Raw public SSH key to install on the VPS (alternative to ssh_keys).",
			},
			"do_not_send_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "If true, do not generate or email a root password. Ensure at least one SSH key is provided or you will not be able to log in.",
			},
			"triggers": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Arbitrary map of values that, when changed, will trigger a re-creation of the reinstall task.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Computed
			"task_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Identifier of the reinstall task.",
			},
			"task_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Final state of the reinstall task (should be \"done\").",
			},
		},
	}
}

func resourceVPSReinstallCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	opts := vpsReinstallOptsFromResource(d)

	// API-allowed but unloggable-into: warn (don't error) when no key is provided
	// and password delivery is suppressed.
	if opts.DoNotSendPassword && len(opts.SshKey) == 0 && opts.PublicSshKey == "" {
		log.Printf("[WARN] do_not_send_password=true with no ssh_keys and no public_ssh_key: VPS %s will be reinstalled without a way to log in", serviceName)
	}

	endpoint := fmt.Sprintf("/vps/%s/reinstall", url.PathEscape(serviceName))
	task := &VPSTask{}

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		if task.Id == 0 {
			return fmt.Errorf("failed to create reinstall task on %s: %w", serviceName, err)
		}
		log.Printf("[WARN] Ignored error when calling POST %s: %v", endpoint, err)
	}

	if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(task.Id, 10))
	d.Set("task_id", task.Id)
	d.Set("task_state", task.State)

	return nil
}

func resourceVPSReinstallRead(d *schema.ResourceData, meta interface{}) error {
	// Nothing to do on READ.
	//
	// IMPORTANT: This resource represents a one-shot reinstall task, not a real
	// long-lived resource. OVH purges old tasks from its database after a while,
	// which would surface as a 404 on the task endpoint. If we returned that 404
	// up to Terraform it would treat the resource as deleted and trigger a fresh
	// reinstall on next apply — which is exactly what we must avoid.
	return nil
}

func resourceVPSReinstallDelete(d *schema.ResourceData, meta interface{}) error {
	// A reinstall task cannot be undone or removed via the API; clear the ID so
	// the state forgets about it.
	d.SetId("")
	return nil
}

func vpsReinstallOptsFromResource(d *schema.ResourceData) *VPSReinstallCreateOpts {
	opts := &VPSReinstallCreateOpts{
		TemplateId:        int64(d.Get("template_id").(int)),
		Language:          d.Get("language").(string),
		PublicSshKey:      d.Get("public_ssh_key").(string),
		DoNotSendPassword: d.Get("do_not_send_password").(bool),
	}

	if v, ok := d.GetOk("software_ids"); ok {
		raw := v.([]interface{})
		ids := make([]int64, 0, len(raw))
		for _, r := range raw {
			ids = append(ids, int64(r.(int)))
		}
		opts.SoftwareId = ids
	}

	if v, ok := d.GetOk("ssh_keys"); ok {
		raw := v.([]interface{})
		keys := make([]string, 0, len(raw))
		for _, r := range raw {
			keys = append(keys, r.(string))
		}
		opts.SshKey = keys
	}

	return opts
}
