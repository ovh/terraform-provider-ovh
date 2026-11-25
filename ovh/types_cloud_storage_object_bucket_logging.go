package ovh

import (
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudStorageObjectBucketLoggingModel struct {
	ServiceName  ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketId     ovhtypes.TfStringValue `tfsdk:"bucket_id"`
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
	TargetBucket ovhtypes.TfStringValue `tfsdk:"target_bucket"`
	TargetPrefix ovhtypes.TfStringValue `tfsdk:"target_prefix"`
}

func (m *CloudStorageObjectBucketLoggingModel) ToS3() *s3LoggingEnabled {
	return &s3LoggingEnabled{
		TargetBucket: m.TargetBucket.ValueString(),
		TargetPrefix: m.TargetPrefix.ValueString(),
	}
}

func (m *CloudStorageObjectBucketLoggingModel) MergeWith(logging *s3LoggingEnabled) {
	m.Id = m.BucketId
	m.TargetBucket = ovhtypes.NewTfStringValue(logging.TargetBucket)
	m.TargetPrefix = ovhtypes.NewTfStringValue(logging.TargetPrefix)
}
