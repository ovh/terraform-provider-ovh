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

func resourceDbaasLogsRolePermissionStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasLogsRolePermissionStreamCreate,
		ReadContext:   resourceDbaasLogsRolePermissionStreamRead,
		DeleteContext: resourceDbaasLogsRolePermissionStreamDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDbaasLogsRolePermissionStreamImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Description: "Role ID to which the permission will be appended",
				Required:    true,
				ForceNew:    true,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Description: "Graylog stream ID to be associated as a permission",
				Required:    true,
				ForceNew:    true,
			},
			// Computed fields from the API GET response.
			"permission_id": {
				Type:        schema.TypeString,
				Description: "Permission ID",
				Computed:    true,
			},
			"permission_type": {
				Type:        schema.TypeString,
				Description: "Permission type (e.g., READ_ONLY)",
				Computed:    true,
			},
		},
	}
}

// resourceDbaasLogsRolePermissionStreamCreate adds a stream permission to a role.
func resourceDbaasLogsRolePermissionStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	roleId := d.Get("role_id").(string)
	streamId := d.Get("stream_id").(string)

	log.Printf("[DEBUG] Adding Graylog stream permission: service=%s, role=%s", serviceName, roleId)

	opts := (&DbaasLogsRolePermissionStreamOpts{}).FromResource(d)
	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s/permission/stream", url.PathEscape(serviceName), url.PathEscape(roleId))
	if err := config.OVHClient.Post(endpoint, opts, res); err != nil {
		return diag.Errorf("Error calling POST %s:\n\t%q", endpoint, err)
	}

	// Wait for the asynchronous operation to complete.
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
		return diag.FromErr(err)
	}

	// Retrieve the full list of permission IDs for the role.
	listEndpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s/permission", url.PathEscape(serviceName), url.PathEscape(roleId))
	var permissionIDs []string
	if err := config.OVHClient.Get(listEndpoint, &permissionIDs); err != nil {
		return diag.Errorf("Error retrieving permissions list (%s):\n\t%q", listEndpoint, err)
	}

	// Iterate over each permission ID to retrieve full details and filter by streamId.
	var foundPermission *DbaasLogsRolePermissionStream
	for _, pid := range permissionIDs {
		permEndpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s/permission/%s", url.PathEscape(serviceName), url.PathEscape(roleId), url.PathEscape(pid))
		var permission DbaasLogsRolePermissionStream
		if err := config.OVHClient.Get(permEndpoint, &permission); err != nil {
			// Optionally, log error and continue with next permission.
			log.Printf("[WARN] Failed to get permission details for ID %s: %v", pid, err)
			continue
		}
		if permission.StreamId == streamId {
			foundPermission = &permission
			break
		}
	}

	if foundPermission == nil {
		return diag.Errorf("No permission found for streamId %s", streamId)
	}

	d.SetId(foundPermission.PermissionId)
	return resourceDbaasLogsRolePermissionStreamRead(ctx, d, meta)
}

// resourceDbaasLogsRolePermissionStreamRead reads the details of a permission.
func resourceDbaasLogsRolePermissionStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	roleId := d.Get("role_id").(string)
	permissionId := d.Id()

	log.Printf("[DEBUG] Reading permission: service=%s, role=%s, permission=%s", serviceName, roleId, permissionId)

	permission := &DbaasLogsRolePermissionStream{}
	endpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s/permission/%s", url.PathEscape(serviceName), url.PathEscape(roleId), url.PathEscape(permissionId))
	if err := config.OVHClient.Get(endpoint, permission); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	d.Set("permission_id", permission.PermissionId)
	d.Set("permission_type", permission.PermissionType)
	d.Set("stream_id", permission.StreamId)

	return nil
}

// resourceDbaasLogsRolePermissionStreamDelete deletes the specified permission.
func resourceDbaasLogsRolePermissionStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	roleId := d.Get("role_id").(string)
	permissionId := d.Id()

	log.Printf("[DEBUG] Deleting permission: service=%s, role=%s, permission=%s", serviceName, roleId, permissionId)

	res := &DbaasLogsOperation{}
	endpoint := fmt.Sprintf("/dbaas/logs/%s/role/%s/permission/%s", url.PathEscape(serviceName), url.PathEscape(roleId), url.PathEscape(permissionId))
	if err := config.OVHClient.Delete(endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	// Wait for the asynchronous deletion operation to complete.
	if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// resourceDbaasLogsRolePermissionStreamImportState imports a permission given an ID in the format "service_name/role_id/permission_id".
func resourceDbaasLogsRolePermissionStreamImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenID := d.Id()
	splitID := strings.SplitN(givenID, "/", 3)
	if len(splitID) != 3 {
		return nil, fmt.Errorf("Import ID must be in service_name/role_id/permission_id format")
	}
	serviceName := splitID[0]
	roleId := splitID[1]
	permissionId := splitID[2]

	d.Set("service_name", serviceName)
	d.Set("role_id", roleId)
	d.SetId(permissionId)

	return []*schema.ResourceData{d}, nil
}
