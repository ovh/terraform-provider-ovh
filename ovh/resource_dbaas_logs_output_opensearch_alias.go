package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"golang.org/x/exp/slices"
)

func resourceDbaasLogsOutputOpensearchAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasLogsOutputOpensearchAliasCreate,
		ReadContext:   resourceDbaasLogsOutputOpensearchAliasRead,
		UpdateContext: resourceDbaasLogsOutputOpensearchAliasUpdate,
		DeleteContext: resourceDbaasLogsOutputOpensearchAliasDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDbaasLogsOutputOpensearchAliasImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Alias description",
				Required:    true,
			},
			"suffix": {
				Type:        schema.TypeString,
				Description: "Alias suffix",
				Required:    true,
			},

			// computed
			"alias_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alias used",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation creation",
			},
			"current_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current alias size (in bytes)",
			},
			"is_editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if you are allowed to edit entry",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alias name",
			},
			"nb_index": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of index",
			},
			"nb_stream": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of shard",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation last update",
			},

			"indexes": {
				Type:        schema.TypeSet,
				Description: "Indexes attached to alias",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"streams": {
				Type:        schema.TypeSet,
				Description: "Streams attached to alias",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDbaasLogsOutputOpensearchAliasImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	serviceName, id, ok := strings.Cut(givenID, "/")
	if !ok {
		return nil, fmt.Errorf("Import Id is not service_name/id formatted")
	}
	d.SetId(id)
	d.Set("service_name", serviceName)
	return []*schema.ResourceData{d}, nil
}

func resourceDbaasLogsOutputOpensearchAliasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will create dbaas logs output opensearch alias for: %s", serviceName)

	opts := (&DbaasLogsOutputOpensearchAliasCreateOps{}).FromResource(d)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/alias",
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

	id := op.AliasId
	if id == nil {
		return diag.Errorf("Alias Id is nil. This should not happen: operation is %s/%s", serviceName, res.OperationId)
	}

	d.SetId(*id)

	indexes := d.Get("indexes").(*schema.Set)
	for _, index := range indexes.List() {
		if err = resourceDbaasLogsOutputOpensearchAliasAttachIndex(ctx, config, serviceName, *id, index.(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	streams := d.Get("streams").(*schema.Set)
	for _, stream := range streams.List() {
		if err = resourceDbaasLogsOutputOpensearchAliasAttachStream(ctx, config, serviceName, *id, stream.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDbaasLogsOutputOpensearchAliasRead(ctx, d, meta)
}

func resourceDbaasLogsOutputOpensearchAliasUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will update dbaas logs output Opensearch alias for: %s", serviceName)

	if d.HasChange("description") {
		opts := (&DbaasLogsOutputOpensearchAliasUpdateOps{}).FromResource(d)
		res := &DbaasLogsOperation{}
		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/output/opensearch/alias/%s",
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
	}

	if d.HasChange("indexes") {
		old, new := d.GetChange("indexes")
		oldIndexesSet := old.(*schema.Set)
		newIndexesSet := new.(*schema.Set)
		oldIndexes := oldIndexesSet.List()
		newIndexes := newIndexesSet.List()
		for _, idx := range oldIndexes {
			if !slices.Contains(newIndexes, idx) {
				if err := resourceDbaasLogsOutputOpensearchAliasDetachIndex(ctx, config, serviceName, id, idx.(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		for _, idx := range newIndexes {
			if !slices.Contains(oldIndexes, idx) {
				if err := resourceDbaasLogsOutputOpensearchAliasAttachIndex(ctx, config, serviceName, id, idx.(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("streams") {
		old, new := d.GetChange("streams")
		oldStreamsSet := old.(*schema.Set)
		newStreamsSet := new.(*schema.Set)
		oldStreams := oldStreamsSet.List()
		newStreams := newStreamsSet.List()
		for _, idx := range oldStreams {
			if !slices.Contains(newStreams, idx) {
				if err := resourceDbaasLogsOutputOpensearchAliasDetachStream(ctx, config, serviceName, id, idx.(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		for _, idx := range newStreams {
			if !slices.Contains(oldStreams, idx) {
				if err := resourceDbaasLogsOutputOpensearchAliasAttachStream(ctx, config, serviceName, id, idx.(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return resourceDbaasLogsOutputOpensearchAliasRead(ctx, d, meta)
}

func resourceDbaasLogsOutputOpensearchAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will read dbaas logs output Opensearch alias: %s/%s", serviceName, id)
	res := &DbaasLogsOutputOpensearchAlias{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/alias/%s",
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

	idx, err := resourceDbaasLogsOutputOpensearchAliasIndexRead(ctx, config, serviceName, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("indexes", idx)

	streams, err := resourceDbaasLogsOutputOpensearchAliasStreamRead(ctx, config, serviceName, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("streams", streams)

	return nil
}

func resourceDbaasLogsOutputOpensearchAliasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	id := d.Id()

	indexes := d.Get("indexes").(*schema.Set)
	for _, index := range indexes.List() {
		if err := resourceDbaasLogsOutputOpensearchAliasDetachIndex(ctx, config, serviceName, id, index.(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	streams := d.Get("streams").(*schema.Set)
	for _, stream := range streams.List() {
		if err := resourceDbaasLogsOutputOpensearchAliasDetachStream(ctx, config, serviceName, id, stream.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] Will delete dbaas logs output Opensearch alias: %s/%s", serviceName, id)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/alias/%s",
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

func resourceDbaasLogsOutputOpensearchAliasAttachIndex(ctx context.Context, config *Config, serviceName, aliasID, indexID string) error {
	endpoint := fmt.Sprintf("/dbaas/logs/%s/output/opensearch/alias/%s/index", url.PathEscape(serviceName), url.PathEscape(aliasID))
	res := &DbaasLogsOperation{}

	if err := config.OVHClient.Post(endpoint, &DbaasLogsOutputOpensearchAliasIndexCreate{IndexID: indexID}, &res); err != nil {
		return fmt.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	_, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId)
	if err != nil {
		return err
	}

	return nil
}

func resourceDbaasLogsOutputOpensearchAliasDetachIndex(ctx context.Context, config *Config, serviceName, aliasID, indexID string) error {
	endpoint := fmt.Sprintf("/dbaas/logs/%s/output/opensearch/alias/%s/index/%s", url.PathEscape(serviceName), url.PathEscape(aliasID), url.PathEscape(indexID))
	res := &DbaasLogsOperation{}

	if err := config.OVHClient.DeleteWithContext(ctx, endpoint, &res); err != nil {
		return fmt.Errorf("Error calling delete %s:\n\t %q", endpoint, err)
	}

	_, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId)
	if err != nil {
		return err
	}

	return nil
}

func resourceDbaasLogsOutputOpensearchAliasAttachStream(ctx context.Context, config *Config, serviceName, aliasID, streamId string) error {
	endpoint := fmt.Sprintf("/dbaas/logs/%s/output/opensearch/alias/%s/stream", url.PathEscape(serviceName), url.PathEscape(aliasID))
	res := &DbaasLogsOperation{}

	if err := config.OVHClient.Post(endpoint, &DbaasLogsOutputOpensearchAliasStreamCreate{StreamID: streamId}, &res); err != nil {
		return fmt.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	_, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId)
	if err != nil {
		return err
	}

	return nil
}

func resourceDbaasLogsOutputOpensearchAliasDetachStream(ctx context.Context, config *Config, serviceName, aliasID, streamId string) error {
	endpoint := fmt.Sprintf("/dbaas/logs/%s/output/opensearch/alias/%s/stream/%s", url.PathEscape(serviceName), url.PathEscape(aliasID), url.PathEscape(streamId))
	res := &DbaasLogsOperation{}

	if err := config.OVHClient.DeleteWithContext(ctx, endpoint, &res); err != nil {
		return fmt.Errorf("Error calling delete %s:\n\t %q", endpoint, err)
	}

	_, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId)
	if err != nil {
		return err
	}

	return nil
}

func resourceDbaasLogsOutputOpensearchAliasIndexRead(ctx context.Context, config *Config, serviceName, aliasId string) ([]string, error) {
	var (
		endpoint = fmt.Sprintf("/dbaas/logs/%s/output/opensearch/alias/%s/index", url.PathEscape(serviceName), url.PathEscape(aliasId))
		indexes  []string
	)

	if err := config.OVHClient.GetWithContext(ctx, endpoint, &indexes); err != nil {
		return nil, fmt.Errorf("failed to list attached indexes: %w", err)
	}

	return indexes, nil
}

func resourceDbaasLogsOutputOpensearchAliasStreamRead(ctx context.Context, config *Config, serviceName, aliasId string) ([]string, error) {
	var (
		endpoint = fmt.Sprintf("/dbaas/logs/%s/output/opensearch/alias/%s/stream", url.PathEscape(serviceName), url.PathEscape(aliasId))
		streams  []string
	)

	if err := config.OVHClient.GetWithContext(ctx, endpoint, &streams); err != nil {
		return nil, fmt.Errorf("failed to list attached indexes: %w", err)
	}

	return streams, nil
}
