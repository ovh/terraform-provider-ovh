package ovh

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type GetCloudProjectRegionLoadbalancerLogSubscriptionResponse struct {
	CreatedAt      string                                                        `json:"createdAt"`
	Kind           string                                                        `json:"kind"`
	Resource       CloudProjectRegionLoadbalancerLogSubscriptionResourceResponse `json:"resource"`
	LDPServiceName string                                                        `json:"serviceName"`
	StreamID       string                                                        `json:"streamId"`
	SubscriptionID string                                                        `json:"subscriptionId"`
	UpdatedAt      string                                                        `json:"updatedAt"`
}

type CreateCloudProjectRegionLoadbalancerLogSubscriptionResponse struct {
	ServiceName string `json:"serviceName"`
	OperationID string `json:"operationId"`
}

func (v CreateCloudProjectRegionLoadbalancerLogSubscriptionResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["ldp_service_name"] = v.ServiceName
	obj["operation_id"] = v.OperationID

	return obj
}

func (v GetCloudProjectRegionLoadbalancerLogSubscriptionResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = v.CreatedAt
	obj["kind"] = v.Kind
	obj["resource_type"] = v.Resource.ToMap()
	obj["ldp_service_name"] = v.LDPServiceName
	obj["stream_id"] = v.StreamID
	obj["subscription_id"] = v.SubscriptionID
	obj["updated_at"] = v.UpdatedAt

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
