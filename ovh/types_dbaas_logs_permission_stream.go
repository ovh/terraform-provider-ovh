package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DbaasLogsRolePermissionStreamOpts holds the payload for appending a stream permission.
type DbaasLogsRolePermissionStreamOpts struct {
	StreamId string `json:"streamId"`
}

func (opts *DbaasLogsRolePermissionStreamOpts) FromResource(d *schema.ResourceData) *DbaasLogsRolePermissionStreamOpts {
	opts.StreamId = d.Get("stream_id").(string)
	return opts
}

// DbaasLogsRolePermission represents the permission as returned by the GET endpoint.
type DbaasLogsRolePermissionStream struct {
	PermissionId   string `json:"permissionId"`
	PermissionType string `json:"permissionType"`
	StreamId       string `json:"streamId"`
}
