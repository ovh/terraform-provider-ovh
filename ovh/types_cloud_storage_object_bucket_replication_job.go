package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudStorageObjectBucketReplicationJobAPIResponse struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
	Failed    int64  `json:"failed"`
	Stuck     int64  `json:"stuck"`
	CreatedAt string `json:"createdAt"`
}

type CloudStorageObjectBucketReplicationJobModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketId    ovhtypes.TfStringValue `tfsdk:"bucket_id"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	Status      ovhtypes.TfStringValue `tfsdk:"status"`
	Total       types.Int64            `tfsdk:"total"`
	Completed   types.Int64            `tfsdk:"completed"`
	Failed      types.Int64            `tfsdk:"failed"`
	Stuck       types.Int64            `tfsdk:"stuck"`
	CreatedAt   ovhtypes.TfStringValue `tfsdk:"created_at"`
}

func (m *CloudStorageObjectBucketReplicationJobModel) MergeWith(response *CloudStorageObjectBucketReplicationJobAPIResponse) {
	m.Id = ovhtypes.NewTfStringValue(response.Id)
	m.Status = ovhtypes.NewTfStringValue(response.Status)
	m.Total = types.Int64Value(response.Total)
	m.Completed = types.Int64Value(response.Completed)
	m.Failed = types.Int64Value(response.Failed)
	m.Stuck = types.Int64Value(response.Stuck)
	m.CreatedAt = ovhtypes.NewTfStringValue(response.CreatedAt)
}

type storageObjectBucketReplicationJobItemModel struct {
	Id        ovhtypes.TfStringValue `tfsdk:"id"`
	Status    ovhtypes.TfStringValue `tfsdk:"status"`
	Total     types.Int64            `tfsdk:"total"`
	Completed types.Int64            `tfsdk:"completed"`
	Failed    types.Int64            `tfsdk:"failed"`
	Stuck     types.Int64            `tfsdk:"stuck"`
	CreatedAt ovhtypes.TfStringValue `tfsdk:"created_at"`
}

func newStorageObjectBucketReplicationJobItem(response *CloudStorageObjectBucketReplicationJobAPIResponse) storageObjectBucketReplicationJobItemModel {
	return storageObjectBucketReplicationJobItemModel{
		Id:        ovhtypes.NewTfStringValue(response.Id),
		Status:    ovhtypes.NewTfStringValue(response.Status),
		Total:     types.Int64Value(response.Total),
		Completed: types.Int64Value(response.Completed),
		Failed:    types.Int64Value(response.Failed),
		Stuck:     types.Int64Value(response.Stuck),
		CreatedAt: ovhtypes.NewTfStringValue(response.CreatedAt),
	}
}
