package ovh

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type CloudProjectRegionLoadbalancerLogSubscriptionResponse struct {
	CreatedAt      string                                                        `json:"createdAt"`
	Kind           string                                                        `json:"kind"`
	OperationID    string                                                        `json:"operationId"`
	Resource       CloudProjectRegionLoadbalancerLogSubscriptionResourceResponse `json:"resource"`
	LDPServiceName string                                                        `json:"serviceName"`
	StreamID       string                                                        `json:"streamId"`
	SubscriptionID string                                                        `json:"subscriptionId"`
	UpdatedAt      string                                                        `json:"updatedAt"`
}

func (v CloudProjectRegionLoadbalancerLogSubscriptionResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["createdAt"] = v.CreatedAt
	obj["kind"] = v.Kind
	obj["resource"] = v.Resource
	obj["serviceName"] = v.LDPServiceName
	obj["streamId"] = v.StreamID
	obj["subscriptionId"] = v.SubscriptionID
	obj["updatedAt"] = v.UpdatedAt
	obj["operation_id"] = v.OperationID

	return obj
}

type CloudProjectRegionLoadbalancerLogSubscriptionResourceResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (v CloudProjectRegionLoadbalancerLogSubscriptionResourceResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["name"] = v.Name
	obj["type"] = v.Type

	return obj
}

type CloudProjectRegionLoadbalancerLogSubscriptionResourceCreateOpts struct {
	Kind     string `json:"kind"`
	StreamId string `json:"streamId"`
}

func (opts *CloudProjectRegionLoadbalancerLogSubscriptionResourceCreateOpts) fromResource(d *schema.ResourceData) *CloudProjectRegionLoadbalancerLogSubscriptionResourceCreateOpts {
	opts.StreamId = d.Get("stream_id").(string)
	opts.Kind = d.Get("kind").(string)
	return opts
}
