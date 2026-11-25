package ovh

import (
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudStorageObjectBucketACLModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketId    ovhtypes.TfStringValue `tfsdk:"bucket_id"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	ACL         ovhtypes.TfStringValue `tfsdk:"acl"`
}
