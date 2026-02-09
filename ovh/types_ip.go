package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

type IpReverse struct {
	IpReverse string `json:"ipReverse"`
	Reverse   string `json:"reverse"`
}

type IpTaskFunctionEnum string

const (
	IpTaskFunctionEnumArinBlockReassign                 IpTaskFunctionEnum = "arinBlockReassign"
	IpTaskFunctionEnumChangeRipeOrg                     IpTaskFunctionEnum = "changeRipeOrg"
	IpTaskFunctionEnumCheckAndReleaseIp                 IpTaskFunctionEnum = "checkAndReleaseIp"
	IpTaskFunctionEnumGenericMoveFloatingIp             IpTaskFunctionEnum = "genericMoveFloatingIp"
	IpTaskFunctionEnumSupernetByoipFailoverPartitioning IpTaskFunctionEnum = "supernetByoipFailoverPartitioning"
)

type IpTaskStatusEnum string

const (
	IpTaskStatusCancelled     IpTaskStatusEnum = "cancelled"
	IpTaskStatusCustomerError IpTaskStatusEnum = "customerError"
	IpTaskStatusDoing         IpTaskStatusEnum = "doing"
	IpTaskStatusDone          IpTaskStatusEnum = "done"
	IpTaskStatusInit          IpTaskStatusEnum = "init"
	IpTaskStatusOvhError      IpTaskStatusEnum = "ovhError"
	IpTaskStatusTodo          IpTaskStatusEnum = "todo"
)

type IpTask struct {
	Comment     *string            `json:"comment,omitempty"`
	Destination *IpServiceRoutedTo `json:"routedTo,omitempty"`
	DoneDate    *time.Time         `json:"doneDate,omitempty"`
	Function    IpTaskFunctionEnum `json:"function"`
	LastUpdate  *time.Time         `json:"lastUpdate,omitempty"`
	StartDate   time.Time          `json:"startDate"`
	Status      IpTaskStatusEnum   `json:"status"`
	TaskId      int64              `json:"taskId"`
}

func (v IpReverse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["ip_reverse"] = v.IpReverse
	obj["reverse"] = v.Reverse
	return obj
}

type IpReverseCreateOpts struct {
	IpReverse string `json:"ipReverse"`
	Reverse   string `json:"reverse"`
}

func (opts *IpReverseCreateOpts) FromResource(d *schema.ResourceData) *IpReverseCreateOpts {
	opts.IpReverse = d.Get("ip_reverse").(string)
	opts.Reverse = d.Get("reverse").(string)

	return opts
}

// Ip represents the response from GET /ip/{ip}
// This is different from IpService which comes from /ip/service/{serviceName}
type Ip struct {
	CanBeTerminated bool               `json:"canBeTerminated"`
	Country         *string            `json:"country"`
	Description     *string            `json:"description"`
	Ip              string             `json:"ip"`
	OrganisationId  *string            `json:"organisationId"`
	RoutedTo        *IpServiceRoutedTo `json:"routedTo"`
	Type            string             `json:"type"`
}

func (v Ip) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["can_be_terminated"] = v.CanBeTerminated
	obj["ip"] = v.Ip
	obj["type"] = v.Type

	if v.Country != nil {
		obj["country"] = *v.Country
	}

	if v.Description != nil {
		obj["description"] = *v.Description
	}

	if v.OrganisationId != nil {
		obj["organisation_id"] = *v.OrganisationId
	}

	if v.RoutedTo != nil {
		obj["routed_to"] = []interface{}{v.RoutedTo.ToMap()}
	}

	return obj
}
