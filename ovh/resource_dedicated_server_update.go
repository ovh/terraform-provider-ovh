package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDedicatedServerUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedServerUpdateCreateOrUpdate,
		UpdateContext: resourceDedicatedServerUpdateCreateOrUpdate,
		ReadContext:   resourceDedicatedServerUpdateRead,
		DeleteContext: resourceDedicatedServerUpdateDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your dedicated server.",
				Required:    true,
			},
			"boot_id": {
				Type:        schema.TypeInt,
				Description: "The boot id of your dedicated server.",
				Computed:    true,
				Optional:    true,
			},
			"boot_script": {
				Type:        schema.TypeString,
				Description: "The boot script of your dedicated server.",
				Optional:    true,
			},
			"efi_bootloader_path": {
				Type:        schema.TypeString,
				Description: "The path of the EFI bootloader.",
				Computed:    true,
				Optional:    true,
			},
			"monitoring": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Icmp monitoring state",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "error, hacked, hackedBlocked, ok",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(v.(string), []string{"error", "hacked", "hackedBlocked", "ok"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Display name of the dedicated server",
			},
		},
	}
}

func resourceDedicatedServerUpdateCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	opts := (&DedicatedServerUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return diag.Errorf("Error calling PUT %s:\n\t %q", endpoint, err)
	}

	if d.HasChange("display_name") {
		newDisplayName := d.Get("display_name").(string)
		if err := serviceUpdateDisplayName(ctx, config, "dedicated/server", serviceName, newDisplayName); err != nil {
			return diag.Errorf("failed to update display name: %s", err)
		}
	}

	//set fake id
	d.SetId(serviceName)

	return resourceDedicatedServerUpdateRead(ctx, d, meta)
}

func resourceDedicatedServerUpdateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	ds := &DedicatedServer{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s",
			url.PathEscape(serviceName),
		),
		&ds,
	)

	if err != nil {
		return diag.Errorf(
			"Error calling GET /dedicated/server/%s:\n\t %q",
			serviceName,
			err,
		)
	}

	d.Set("boot_id", ds.BootId)
	d.Set("boot_script", ds.BootScript)
	d.Set("efi_bootloader_path", ds.EfiBootloaderPath)
	d.Set("monitoring", ds.Monitoring)
	d.Set("state", ds.State)
	d.Set("display_name", ds.DisplayName)

	//set fake id
	d.SetId(serviceName)
	return nil
}

func resourceDedicatedServerUpdateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
