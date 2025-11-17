package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type IpService struct {
	CanBeTerminated bool               `json:"canBeTerminated"`
	Country         *string            `json:"country"`
	Description     *string            `json:"description"`
	Ip              string             `json:"ip"`
	OrganisationId  *string            `json:"organisationId"`
	RoutedTo        *IpServiceRoutedTo `json:"routedTo"`
	Type            string             `json:"type"`
}

func (v IpService) ToMap() map[string]interface{} {
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

type IpServiceRoutedTo struct {
	ServiceName string `json:"serviceName"`
}

func (v IpServiceRoutedTo) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["service_name"] = v.ServiceName
	return obj
}

type IpServiceUpdateOpts struct {
	Description *string `json:"description,omitempty"`
}

func (opts *IpServiceUpdateOpts) FromResource(d *schema.ResourceData) *IpServiceUpdateOpts {
	opts.Description = helpers.GetNilStringPointerFromData(d, "description")
	return opts
}

// see https://api.ovh.com/console/#/ip/%7Bip%7D/move~POST
type IpMoveOpts struct {
	NextHop *string `json:"nexthop,omitempty"`
	To      *string `json:"to"`
}

func (opts *IpMoveOpts) FromResource(d *schema.ResourceData) (*IpMoveOpts, error) {
	opts.To = GetRoutedToServiceName(d)
	return opts, nil
}

func GetRoutedToServiceName(d *schema.ResourceData) *string {
	routedTo := (d.Get("routed_to")).([]interface{})
	if len(routedTo) == 0 || routedTo[0] == nil {
		return nil
	}
	serviceName := (routedTo[0].(map[string]interface{}))["service_name"].(string)
	if serviceName == "" {
		return nil
	}
	return &serviceName
}

type IpServiceConfirmTerminationOpts struct {
	Token string `json:"token"`
}
