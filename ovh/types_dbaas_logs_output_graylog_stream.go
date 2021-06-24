package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type DbaasLogsOutputGraylogStream struct {
	CanAlert                 bool    `json:"canAlert"`
	ColdStorageCompression   *string `json:"coldStorageCompression"`
	ColdStorageContent       *string `json:"coldStorageContent"`
	ColdStorageEnabled       *bool   `json:"coldStorageEnabled"`
	ColdStorageNotifyEnabled *bool   `json:"coldStorageNotifyEnabled"`
	ColdStorageRetention     *int64  `json:"coldStorageRetention"`
	ColdStorageTarget        *string `json:"coldStorageTarget"`
	CreatedAt                string  `json:"createdAt"`
	Description              string  `json:"description"`
	IndexingEnabled          *bool   `json:"indexingEnabled"`
	IndexingMaxSize          *int64  `json:"indexingMaxSize"`
	IndexingNotifyEnabled    *bool   `json:"indexingNotifyEnabled"`
	IsEditable               bool    `json:"isEditable"`
	IsShareable              bool    `json:"isShareable"`
	NbAlertCondition         int64   `json:"nbAlertCondition"`
	NbArchive                int64   `json:"nbArchive"`
	ParentStreamId           *string `json:"parentStreamId"`
	PauseIndexingOnMaxSize   *bool   `json:"pauseIndexingOnMaxSize"`
	RetentionId              string  `json:"retentionId"`
	StreamId                 string  `json:"streamId"`
	Title                    string  `json:"title"`
	UpdatedAt                string  `json:"updatedAt"`
	WebSocketEnabled         *bool   `json:"webSocketEnabled"`
}

func (v DbaasLogsOutputGraylogStream) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["can_alert"] = v.CanAlert
	obj["created_at"] = v.CreatedAt
	obj["description"] = v.Description
	obj["is_editable"] = v.IsEditable
	obj["is_shareable"] = v.IsShareable
	obj["nb_alert_condition"] = v.NbAlertCondition
	obj["nb_archive"] = v.NbArchive
	obj["retention_id"] = v.RetentionId
	obj["stream_id"] = v.StreamId
	obj["title"] = v.Title
	obj["updated_at"] = v.UpdatedAt

	if v.ColdStorageCompression != nil {
		obj["cold_storage_compression"] = *v.ColdStorageCompression
	}
	if v.ColdStorageContent != nil {
		obj["cold_storage_content"] = *v.ColdStorageContent
	}
	if v.ColdStorageEnabled != nil {
		obj["cold_storage_enabled"] = *v.ColdStorageEnabled
	}
	if v.ColdStorageNotifyEnabled != nil {
		obj["cold_storage_notify_enabled"] = *v.ColdStorageNotifyEnabled
	}
	if v.ColdStorageRetention != nil {
		obj["cold_storage_retention"] = *v.ColdStorageRetention
	}
	if v.ColdStorageTarget != nil {
		obj["cold_storage_target"] = *v.ColdStorageTarget
	}
	if v.IndexingEnabled != nil {
		obj["indexing_enabled"] = *v.IndexingEnabled
	}
	if v.IndexingMaxSize != nil {
		obj["indexing_max_size"] = *v.IndexingMaxSize
	}
	if v.IndexingNotifyEnabled != nil {
		obj["indexing_notify_enabled"] = *v.IndexingNotifyEnabled
	}
	if v.ParentStreamId != nil {
		obj["parent_stream_id"] = *v.ParentStreamId
	}
	if v.WebSocketEnabled != nil {
		obj["web_socket_enabled"] = *v.WebSocketEnabled
	}

	return obj
}

type DbaasLogsOutputGraylogStreamCreateOpts struct {
	ColdStorageCompression   *string `json:"coldStorageCompression,omitempty"`
	ColdStorageContent       *string `json:"coldStorageContent,omitempty"`
	ColdStorageEnabled       *bool   `json:"coldStorageEnabled,omitempty"`
	ColdStorageNotifyEnabled *bool   `json:"coldStorageNotifyEnabled,omitempty"`
	ColdStorageRetention     *int64  `json:"coldStorageRetention,omitempty"`
	ColdStorageTarget        *string `json:"coldStorageTarget,omitempty"`
	Description              string  `json:"description"`
	IndexingEnabled          *bool   `json:"indexingEnabled,omitempty"`
	IndexingMaxSize          *int64  `json:"indexingMaxSize,omitempty"`
	IndexingNotifyEnabled    *bool   `json:"indexingNotifyEnabled,omitempty"`
	ParentStreamId           *string `json:"parentStreamId,omitempty"`
	PauseIndexingOnMaxSize   *bool   `json:"pauseIndexingOnMaxSize,omitempty"`
	RetentionId              *string `json:"retentionId,omitempty"`
	Title                    string  `json:"title"`
	WebSocketEnabled         *bool   `json:"webSocketEnabled,omitempty"`
}

func (opts *DbaasLogsOutputGraylogStreamCreateOpts) FromResource(d *schema.ResourceData) *DbaasLogsOutputGraylogStreamCreateOpts {
	opts.ColdStorageCompression = helpers.GetNilStringPointerFromData(d, "cold_storage_compression")
	opts.ColdStorageContent = helpers.GetNilStringPointerFromData(d, "cold_storage_content")
	opts.ColdStorageEnabled = helpers.GetNilBoolPointerFromData(d, "cold_storage_enabled")
	opts.ColdStorageNotifyEnabled = helpers.GetNilBoolPointerFromData(d, "cold_storage_notify_enabled")
	opts.ColdStorageRetention = helpers.GetNilInt64PointerFromData(d, "cold_storage_retention")
	opts.ColdStorageTarget = helpers.GetNilStringPointerFromData(d, "cold_storage_target")
	opts.Description = d.Get("description").(string)
	opts.IndexingEnabled = helpers.GetNilBoolPointerFromData(d, "indexing_enabled")
	opts.IndexingMaxSize = helpers.GetNilInt64PointerFromData(d, "indexing_max_size")
	opts.IndexingNotifyEnabled = helpers.GetNilBoolPointerFromData(d, "indexing_notify_enabled")
	opts.ParentStreamId = helpers.GetNilStringPointerFromData(d, "parent_stream_id")
	opts.PauseIndexingOnMaxSize = helpers.GetNilBoolPointerFromData(d, "pause_indexing_on_max_size")
	opts.RetentionId = helpers.GetNilStringPointerFromData(d, "retention_id")
	opts.Title = d.Get("title").(string)
	opts.WebSocketEnabled = helpers.GetNilBoolPointerFromData(d, "web_socket_enabled")

	return opts
}

type DbaasLogsOutputGraylogStreamUpdateOpts struct {
	ColdStorageCompression   *string `json:"coldStorageCompression,omitempty"`
	ColdStorageContent       *string `json:"coldStorageContent,omitempty"`
	ColdStorageEnabled       *bool   `json:"coldStorageEnabled,omitempty"`
	ColdStorageNotifyEnabled *bool   `json:"coldStorageNotifyEnabled,omitempty"`
	ColdStorageRetention     *int64  `json:"coldStorageRetention,omitempty"`
	ColdStorageTarget        *string `json:"coldStorageTarget,omitempty"`
	Description              string  `json:"description"`
	IndexingEnabled          *bool   `json:"indexingEnabled,omitempty"`
	IndexingMaxSize          *int64  `json:"indexingMaxSize,omitempty"`
	IndexingNotifyEnabled    *bool   `json:"indexingNotifyEnabled,omitempty"`
	ParentStreamId           *string `json:"parentStreamId,omitempty"`
	PauseIndexingOnMaxSize   *bool   `json:"pauseIndexingOnMaxSize,omitempty"`
	RetentionId              *string `json:"retentionId"`
	Title                    string  `json:"title"`
	WebSocketEnabled         *bool   `json:"webSocketEnabled,omitempty"`
}

func (opts *DbaasLogsOutputGraylogStreamUpdateOpts) FromResource(d *schema.ResourceData) *DbaasLogsOutputGraylogStreamUpdateOpts {
	opts.ColdStorageCompression = helpers.GetNilStringPointerFromData(d, "cold_storage_compression")
	opts.ColdStorageContent = helpers.GetNilStringPointerFromData(d, "cold_storage_content")
	opts.ColdStorageEnabled = helpers.GetNilBoolPointerFromData(d, "cold_storage_enabled")
	opts.ColdStorageNotifyEnabled = helpers.GetNilBoolPointerFromData(d, "cold_storage_notify_enabled")
	opts.ColdStorageRetention = helpers.GetNilInt64PointerFromData(d, "cold_storage_retention")
	opts.ColdStorageTarget = helpers.GetNilStringPointerFromData(d, "cold_storage_target")
	opts.Description = d.Get("description").(string)
	opts.IndexingEnabled = helpers.GetNilBoolPointerFromData(d, "indexing_enabled")
	opts.IndexingMaxSize = helpers.GetNilInt64PointerFromData(d, "indexing_max_size")
	opts.IndexingNotifyEnabled = helpers.GetNilBoolPointerFromData(d, "indexing_notify_enabled")
	opts.PauseIndexingOnMaxSize = helpers.GetNilBoolPointerFromData(d, "pause_indexing_on_max_size")
	opts.Title = d.Get("title").(string)
	opts.WebSocketEnabled = helpers.GetNilBoolPointerFromData(d, "web_socket_enabled")

	return opts
}
