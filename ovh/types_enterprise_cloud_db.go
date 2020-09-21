package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"
)

type CloudDBEnterpriseStatus string
type CloudDBEnterpriseSecurityGroupStatus string
type CloudDBEnterpriseSecurityGroupRuleStatus string

const (
	CloudDBEnterpriseStatusCreated    CloudDBEnterpriseStatus = "created"
	CloudDBEnterpriseStatusCreating   CloudDBEnterpriseStatus = "creating"
	CloudDBEnterpriseStatusDeleting   CloudDBEnterpriseStatus = "deleting"
	CloudDBEnterpriseStatusReopening  CloudDBEnterpriseStatus = "reopening"
	CloudDBEnterpriseStatusRestarting CloudDBEnterpriseStatus = "restarting"
	CloudDBEnterpriseStatusScaling    CloudDBEnterpriseStatus = "scaling"
	CloudDBEnterpriseStatusSuspended  CloudDBEnterpriseStatus = "suspended"
	CloudDBEnterpriseStatusSuspending CloudDBEnterpriseStatus = "suspending"
	CloudDBEnterpriseStatusUpdating   CloudDBEnterpriseStatus = "updating"

	CloudDBEnterpriseSecurityGroupStatusCreated  CloudDBEnterpriseSecurityGroupStatus = "created"
	CloudDBEnterpriseSecurityGroupStatusCreating CloudDBEnterpriseSecurityGroupStatus = "creating"
	CloudDBEnterpriseSecurityGroupStatusDeleting CloudDBEnterpriseSecurityGroupStatus = "deleting"
	CloudDBEnterpriseSecurityGroupStatusUpdated  CloudDBEnterpriseSecurityGroupStatus = "updated"
	CloudDBEnterpriseSecurityGroupStatusUpdating CloudDBEnterpriseSecurityGroupStatus = "updating"

	CloudDBEnterpriseSecurityGroupRuleStatusCreated  CloudDBEnterpriseSecurityGroupRuleStatus = "created"
	CloudDBEnterpriseSecurityGroupRuleStatusCreating CloudDBEnterpriseSecurityGroupRuleStatus = "creating"
	CloudDBEnterpriseSecurityGroupRuleStatusDeleting CloudDBEnterpriseSecurityGroupRuleStatus = "deleting"
	CloudDBEnterpriseSecurityGroupRuleStatusUpdated  CloudDBEnterpriseSecurityGroupRuleStatus = "updated"
	CloudDBEnterpriseSecurityGroupRuleStatusUpdating CloudDBEnterpriseSecurityGroupRuleStatus = "updating"
)

type CloudDBEnterprise struct {
	Id         string                  `json:"id"`
	Status     CloudDBEnterpriseStatus `json:"status"`
	RegionName string                  `json:"regionName"`
}

type CloudDBEnterpriseSecurityGroupCreateUpdateOpts struct {
	Name      string `json:"name"`
	ClusterId string `json:"clusterId"`
}

type CloudDBEnterpriseSecurityGroup struct {
	Id     string                  `json:"id"`
	Name   string                  `json:"name"`
	Status CloudDBEnterpriseStatus `json:"status"`
	TaskId string                  `json:"taskId"`
}

func (opts *CloudDBEnterpriseSecurityGroupCreateUpdateOpts) FromResource(d *schema.ResourceData) *CloudDBEnterpriseSecurityGroupCreateUpdateOpts {
	name := helpers.GetNilStringPointerFromData(d, "name")
	opts.Name = *name
	clusterId := helpers.GetNilStringPointerFromData(d, "cluster_id")
	opts.ClusterId = *clusterId
	return opts
}

type CloudDBEnterpriseSecurityGroupRuleCreateUpdateOpts struct {
	Source string `json:"source"`
}

type CloudDBEnterpriseSecurityGroupRule struct {
	Id     string                                   `json:"id"`
	Source string                                   `json:"source"`
	Status CloudDBEnterpriseSecurityGroupRuleStatus `json:"status"`
	TaskId string                                   `json:"taskId"`
}

func (opts *CloudDBEnterpriseSecurityGroupRuleCreateUpdateOpts) FromResource(d *schema.ResourceData) *CloudDBEnterpriseSecurityGroupRuleCreateUpdateOpts {
	source := helpers.GetNilStringPointerFromData(d, "source")
	opts.Source = *source
	return opts
}
