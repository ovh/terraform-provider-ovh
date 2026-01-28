package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Account Logs

type AccountLogsResource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AccountLogsResponse struct {
	CreatedAt      string              `json:"createdAt"`
	Kind           string              `json:"kind"`
	LogType        string              `json:"logType"`
	LDPServiceName string              `json:"serviceName"`
	StreamID       string              `json:"streamId"`
	SubscriptionID string              `json:"subscriptionId"`
	UpdatedAt      string              `json:"updatedAt"`
}

func (r AccountLogsResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = r.CreatedAt
	obj["id"] = r.SubscriptionID
	obj["kind"] = r.Kind
	obj["ldp_service_name"] = r.LDPServiceName
	obj["log_type"] = r.LogType
	obj["stream_id"] = r.StreamID
	obj["updated_at"] = r.UpdatedAt

	return obj
}

type AccountLogsCreateOpts struct {
	StreamID string `json:"streamId"`
	Kind     string `json:"kind"`
}

func (opts *AccountLogsCreateOpts) fromResource(d *schema.ResourceData) *AccountLogsCreateOpts {
	opts.StreamID = d.Get("stream_id").(string)
	opts.Kind = d.Get("kind").(string)
	return opts
}

type AccountLogsListResponse struct {
	Data []AccountLogsResponse `json:"data"`
}

func (r AccountLogsListResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})
	logs := make([]map[string]interface{}, len(r.Data))
	for i, log := range r.Data {
		logs[i] = log.toMap()
	}
	obj["logs"] = logs
	return obj
}
