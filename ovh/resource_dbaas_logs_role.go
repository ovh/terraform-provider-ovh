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
)

func resourceDbaasLogsRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasLogsRoleCreate,
		ReadContext:   resourceDbaasLogsRoleRead,
		UpdateContext: resourceDbaasLogsRoleUpdate,
		DeleteContext: resourceDbaasLogsRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDbaasLogsRoleImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The role name",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The role description",
				Required:    true,
			},
			// Computed fields from the API response
			"created_at": {
				Type:        schema.TypeString,
				Description: "Role creation date",
				Computed:    true,
			},
			"nb_member": {
				Type:        schema.TypeInt,
				Description: "Number of members in the role",
				Computed:    true,
			},
			"nb_permission": {
				Type:        schema.TypeInt,
				Description: "Number of permissions assigned to the role",
				Computed:    true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Description: "Role identifier",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Role last update date",
				Computed:    true,
			},
		},
	}
}

// resourceDbaasLogsRoleCreate creates a new role using an asynchronous operation.
func resourceDbaasLogsRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Creating dbaas logs role for service: %s", serviceName)

	opts := (&DbaasLogsRoleCreateOpts{}).FromResource(d)
	opRes := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf("/dbaas/logs/%s/role", url.PathEscape(serviceName))
	if err := config.OVHClient.Post(endpoint, opts, opRes); err != nil {
		return diag.Errorf("Error calling POST %s:\n\t%q", endpoint, err)
	}

	// Wait for asynchronous operation to complete.
	op, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, opRes.OperationId)
	if err != nil {
		return diag.FromErr(err)
	}

	if op.RoleId == nil {
		return diag.Errorf("Role Id is nil. This should not happen: operation %s/%s", serviceName, opRes.OperationId)
	}
	d.SetId(*op.RoleId)

	return resourceDbaasLogsRoleRead(ctx, d, meta)
}

// resourceDbaasLogsRoleRead reads the role resource.
func resourceDbaasLogsRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	roleId := d.Id()

	log.Printf("[DEBUG] Reading dbaas logs role for service: %s, role: %s", serviceName, roleId)

	role := &DbaasLogsRole{}
	endpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s", url.PathEscape(serviceName), url.PathEscape(roleId))
	if err := config.OVHClient.Get(endpoint, role); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range role.ToMap() {
		if k != "role_id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	return nil
}

// resourceDbaasLogsRoleUpdate updates the role using an asynchronous operation.
func resourceDbaasLogsRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	roleId := d.Id()

	log.Printf("[DEBUG] Updating dbaas logs role for service: %s, role: %s", serviceName, roleId)

	opts := (&DbaasLogsRoleUpdateOpts{}).FromResource(d)
	opRes := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s", url.PathEscape(serviceName), url.PathEscape(roleId))
	if err := config.OVHClient.Post(endpoint, opts, opRes); err != nil {
		return diag.Errorf("Error calling POST %s:\n\t%q", endpoint, err)
	}

	// Wait for asynchronous operation to complete.
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, opRes.OperationId); err != nil {
		return diag.FromErr(err)
	}

	return resourceDbaasLogsRoleRead(ctx, d, meta)
}

// resourceDbaasLogsRoleDelete deletes the role using an asynchronous operation.
func resourceDbaasLogsRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	roleId := d.Id()

	log.Printf("[DEBUG] Deleting dbaas logs role for service: %s, role: %s", serviceName, roleId)
	opRes := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s", url.PathEscape(serviceName), url.PathEscape(roleId))
	if err := config.OVHClient.Delete(endpoint, opRes); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	// Wait for asynchronous operation to complete.
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, opRes.OperationId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// resourceDbaasLogsRoleImportState allows importing a role using service_name/roleId format.
func resourceDbaasLogsRoleImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	splitID := strings.SplitN(givenID, "/", 2)
	if len(splitID) != 2 {
		return nil, fmt.Errorf("Import ID must be in service_name/roleId format")
	}
	serviceName := splitID[0]
	roleId := splitID[1]
	d.SetId(roleId)
	d.Set("service_name", serviceName)

	return []*schema.ResourceData{d}, nil
}
