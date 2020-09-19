package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"
)

type EnterpriseCloudDBStatus string
type EnterpriseCloudDBSecurityGroupStatus string
type EnterpriseCloudDBSecurityGroupRuleStatus string

const (
	EnterpriseCloudDBStatusCreated    EnterpriseCloudDBStatus = "created"
	EnterpriseCloudDBStatusCreating   EnterpriseCloudDBStatus = "creating"
	EnterpriseCloudDBStatusDeleting   EnterpriseCloudDBStatus = "deleting"
	EnterpriseCloudDBStatusReopening  EnterpriseCloudDBStatus = "reopening"
	EnterpriseCloudDBStatusRestarting EnterpriseCloudDBStatus = "restarting"
	EnterpriseCloudDBStatusScaling    EnterpriseCloudDBStatus = "scaling"
	EnterpriseCloudDBStatusSuspended  EnterpriseCloudDBStatus = "suspended"
	EnterpriseCloudDBStatusSuspending EnterpriseCloudDBStatus = "suspending"
	EnterpriseCloudDBStatusUpdating   EnterpriseCloudDBStatus = "updating"

	EnterpriseCloudDBSecurityGroupStatusCreated  EnterpriseCloudDBSecurityGroupStatus = "created"
	EnterpriseCloudDBSecurityGroupStatusCreating EnterpriseCloudDBSecurityGroupStatus = "creating"
	EnterpriseCloudDBSecurityGroupStatusDeleting EnterpriseCloudDBSecurityGroupStatus = "deleting"
	EnterpriseCloudDBSecurityGroupStatusUpdated  EnterpriseCloudDBSecurityGroupStatus = "updated"
	EnterpriseCloudDBSecurityGroupStatusUpdating EnterpriseCloudDBSecurityGroupStatus = "updating"

	EnterpriseCloudDBSecurityGroupRuleStatusCreated  EnterpriseCloudDBSecurityGroupRuleStatus = "created"
	EnterpriseCloudDBSecurityGroupRuleStatusCreating EnterpriseCloudDBSecurityGroupRuleStatus = "creating"
	EnterpriseCloudDBSecurityGroupRuleStatusDeleting EnterpriseCloudDBSecurityGroupRuleStatus = "deleting"
	EnterpriseCloudDBSecurityGroupRuleStatusUpdated  EnterpriseCloudDBSecurityGroupRuleStatus = "updated"
	EnterpriseCloudDBSecurityGroupRuleStatusUpdating EnterpriseCloudDBSecurityGroupRuleStatus = "updating"
)

type EnterpriseCloudDB struct {
	Id         string                  `json:"id"`
	Status     EnterpriseCloudDBStatus `json:"status"`
	RegionName string                  `json:"regionName"`
}

type EnterpriseCloudDBSecurityGroupCreateUpdateOpts struct {
	Name      string `json:"name"`
	ClusterId string `json:"clusterId"`
}

type EnterpriseCloudDBSecurityGroup struct {
	Id     string                  `json:"id"`
	Name   string                  `json:"name"`
	Status EnterpriseCloudDBStatus `json:"status"`
	TaskId string                  `json:"taskId"`
}

func (opts *EnterpriseCloudDBSecurityGroupCreateUpdateOpts) FromResource(d *schema.ResourceData) *EnterpriseCloudDBSecurityGroupCreateUpdateOpts {
	name := helpers.GetNilStringPointerFromData(d, "name")
	opts.Name = *name
	clusterId := helpers.GetNilStringPointerFromData(d, "cluster_id")
	opts.ClusterId = *clusterId
	return opts
}

type EnterpriseCloudDBSecurityGroupRuleCreateUpdateOpts struct {
	Source string `json:"source"`
}

type EnterpriseCloudDBSecurityGroupRule struct {
	Id     string                                   `json:"id"`
	Source string                                   `json:"source"`
	Status EnterpriseCloudDBSecurityGroupRuleStatus `json:"status"`
	TaskId string                                   `json:"taskId"`
}

func (opts *EnterpriseCloudDBSecurityGroupRuleCreateUpdateOpts) FromResource(d *schema.ResourceData) *EnterpriseCloudDBSecurityGroupRuleCreateUpdateOpts {
	source := helpers.GetNilStringPointerFromData(d, "source")
	opts.Source = *source
	return opts
}
