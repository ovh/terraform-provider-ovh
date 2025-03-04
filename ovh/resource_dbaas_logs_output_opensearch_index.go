package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceDbaasLogsOutputOpensearchIndex() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasLogsOutputOpensearchIndexCreate,
		ReadContext:   resourceDbaasLogsOutputOpensearchIndexRead,
		UpdateContext: resourceDbaasLogsOutputOpensearchIndexUpdate,
		DeleteContext: resourceDbaasLogsOutputOpensearchIndexDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDbaasLogsOutputOpensearchIndexImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Index description",
				Required:    true,
			},
			"nb_shard": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of shard",
			},
			"suffix": {
				Type:        schema.TypeString,
				Description: "Index suffix",
				Required:    true,
			},

			// computed
			"alert_notify_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set, notify when size is near 80, 90 or 100 % of its maximum capacity",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation creation",
			},
			"current_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current Index size (in bytes)",
			},
			"index_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Index ID",
			},
			"is_editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if you are allowed to edit entry",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum index size (in bytes)",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Index name",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation last update",
			},
		},
	}
}

func resourceDbaasLogsOutputOpensearchIndexImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	serviceName, id, ok := strings.Cut(givenID, "/")
	if !ok {
		return nil, fmt.Errorf("Import Id is not service_name/id formatted")
	}
	d.SetId(id)
	d.Set("service_name", serviceName)

	return []*schema.ResourceData{d}, nil
}

func resourceDbaasLogsOutputOpensearchIndexCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will create dbaas logs output opensearch Index for: %s", serviceName)

	opts := (&DbaasLogsOutputOpensearchIndexCreateOps{}).FromResource(d)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/index",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Post(endpoint, opts, res); err != nil {
		return diag.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	// Wait for operation status
	op, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId)
	if err != nil {
		return diag.FromErr(err)
	}
	id := op.IndexId
	if id == nil {
		return diag.Errorf("Index Id is nil. This should not happen: operation is %s/%s", serviceName, res.OperationId)
	}

	d.SetId(*id)

	return resourceDbaasLogsOutputOpensearchIndexRead(ctx, d, meta)
}

func resourceDbaasLogsOutputOpensearchIndexUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will update dbaas logs output Opensearch index for: %s", serviceName)

	opts := (&DbaasLogsOutputOpensearchIndexUpdateOps{}).FromResource(d)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/index/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)
	if err := config.OVHClient.Put(endpoint, opts, res); err != nil {
		return diag.Errorf("Error calling Put %s:\n\t %q", endpoint, err)
	}

	// Wait for operation status
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
		return diag.FromErr(err)
	}

	return resourceDbaasLogsOutputOpensearchIndexRead(ctx, d, meta)
}

func resourceDbaasLogsOutputOpensearchIndexRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will read dbaas logs output Opensearch index: %s/%s", serviceName, id)
	res := &DbaasLogsOutputOpensearchIndex{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/index/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		log.Printf("[ERROR] %s: %v", endpoint, err)
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	return nil
}

func resourceDbaasLogsOutputOpensearchIndexDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete dbaas logs output Openserach index: %s/%s", serviceName, id)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/index/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Delete(endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	// Wait for operation status
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
