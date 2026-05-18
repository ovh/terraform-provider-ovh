package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectKubeLogSubscriptionCreateOpts struct {
	Kind     string `json:"kind"`
	StreamId string `json:"streamId"`
}

type CloudProjectKubeLogSubscriptionResource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Returned by POST (create)
type CloudProjectKubeLogSubscriptionCreateResponse struct {
	OperationId string `json:"operationId"`
	ServiceName string `json:"serviceName"`
}

// Returned by GET (read)
type CloudProjectKubeLogSubscriptionResponse struct {
	CreatedAt      string                                   `json:"createdAt"`
	Kind           string                                   `json:"kind"`
	Resource       *CloudProjectKubeLogSubscriptionResource `json:"resource"`
	ServiceName    string                                   `json:"serviceName"`
	StreamId       string                                   `json:"streamId"`
	SubscriptionId string                                   `json:"subscriptionId"`
	UpdatedAt      string                                   `json:"updatedAt"`
}

func (opts *CloudProjectKubeLogSubscriptionCreateOpts) fromResource(d *schema.ResourceData) *CloudProjectKubeLogSubscriptionCreateOpts {
	opts.Kind = d.Get(kubeLogSubscriptionKindKey).(string)
	opts.StreamId = d.Get(kubeLogSubscriptionStreamIdKey).(string)
	return opts
}

func (v *CloudProjectKubeLogSubscriptionResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj[kubeLogSubscriptionIdKey] = v.SubscriptionId
	obj[kubeLogSubscriptionKindKey] = v.Kind
	obj[kubeLogSubscriptionStreamIdKey] = v.StreamId
	obj[kubeCreatedAtKey] = v.CreatedAt
	obj[kubeUpdatedAtKey] = v.UpdatedAt

	if v.Resource != nil {
		obj[kubeLogSubscriptionResourceKey] = []map[string]interface{}{
			{
				kubeLogSubscriptionResourceNameKey: v.Resource.Name,
				kubeLogSubscriptionResourceTypeKey: v.Resource.Type,
			},
		}
	} else {
		obj[kubeLogSubscriptionResourceKey] = []map[string]interface{}{}
	}

	return obj
}
