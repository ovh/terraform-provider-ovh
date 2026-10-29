package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func TestCloudStorageFileShareNetworkToCreate_WithoutNetworkID(t *testing.T) {
	model := CloudStorageFileShareNetworkModel{
		Name:      ovhtypes.NewTfStringValue("share-network"),
		SubnetId:  ovhtypes.NewTfStringValue("subnet-id"),
		Region:    ovhtypes.NewTfStringValue("GRA11"),
		NetworkId: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
	}

	payload := model.ToCreate()
	if payload.TargetSpec.Network != nil {
		t.Fatalf("expected no targetSpec.network when network_id is not configured, got %+v", payload.TargetSpec.Network)
	}
}

func TestCloudStorageFileShareNetworkMergeWith_NetworkID(t *testing.T) {
	tests := []struct {
		name              string
		targetSpecNetwork *CloudStorageFileShareNetworkAPINetworkRef
		expectedNetworkID ovhtypes.TfStringValue
	}{
		{
			name:              "targetSpec without network",
			targetSpecNetwork: nil,
			expectedNetworkID: ovhtypes.TfStringValue{StringValue: types.StringNull()},
		},
		{
			name:              "targetSpec with an empty network id",
			targetSpecNetwork: &CloudStorageFileShareNetworkAPINetworkRef{},
			expectedNetworkID: ovhtypes.TfStringValue{StringValue: types.StringNull()},
		},
		{
			name:              "targetSpec with a network id",
			targetSpecNetwork: &CloudStorageFileShareNetworkAPINetworkRef{Id: "net-id"},
			expectedNetworkID: ovhtypes.NewTfStringValue("net-id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := CloudStorageFileShareNetworkModel{
				NetworkId: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
			}

			model.MergeWith(context.Background(), &CloudStorageFileShareNetworkAPIResponse{
				Id: "share-network-id",
				TargetSpec: &CloudStorageFileShareNetworkAPITargetSpec{
					Name:    "share-network",
					Network: tt.targetSpecNetwork,
					Subnet:  &CloudStorageFileShareNetworkAPISubnetRef{Id: "subnet-id"},
				},
			})

			if !model.NetworkId.Equal(tt.expectedNetworkID) {
				t.Fatalf("expected network_id %s, got %s", tt.expectedNetworkID, model.NetworkId)
			}
		})
	}
}
