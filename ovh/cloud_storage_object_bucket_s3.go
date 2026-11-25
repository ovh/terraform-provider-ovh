package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

// storageObjectBucketS3Target is where the S3 API serves a bucket: its virtual host, which carries
// the endpoint of its storage tier, and the lowercase region that signs the requests.
type storageObjectBucketS3Target struct {
	Host   string
	Region string
}

func storageObjectBucketEndpoint(serviceName, bucketId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/object/bucket/" + url.PathEscape(bucketId)
}

// storageObjectBucketS3TargetRead locates a bucket on the S3 API from its publicCloud identifier.
// The errors of the publicCloud API are returned as is, so a 404 means the bucket is gone.
func storageObjectBucketS3TargetRead(ctx context.Context, client *ovhwrap.Client, serviceName, bucketId string) (*storageObjectBucketS3Target, error) {
	var bucket CloudStorageObjectBucketAPIResponse
	if err := client.GetWithContext(ctx, storageObjectBucketEndpoint(serviceName, bucketId), &bucket); err != nil {
		return nil, err
	}

	state := bucket.CurrentState
	if state == nil || state.VirtualHost == "" || state.Location == nil {
		return nil, fmt.Errorf("bucket %s has no virtual host yet (status %s)", bucketId, bucket.ResourceStatus)
	}

	return &storageObjectBucketS3Target{
		Host:   state.VirtualHost,
		Region: strings.ToLower(state.Location.Region),
	}, nil
}
