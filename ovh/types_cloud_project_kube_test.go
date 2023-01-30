package ovh

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCloudProjectKubeResponse_ToMap(t *testing.T) {
	type fields struct {
		ControlPlaneIsUpToDate bool
		Id                     string
		IsUpToDate             bool
		Name                   string
		NextUpgradeVersions    []string
		NodesUrl               string
		PrivateNetworkId       string
		Region                 string
		Status                 string
		UpdatePolicy           string
		Url                    string
		Version                string
		Customization          Customization
		KubeProxyMode          string
	}

	pointerArray := func(s []string) *[]string { return &s }

	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "No customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1.0.0",
				Customization:          Customization{},
				KubeProxyMode:          "",
			},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"name":                        "",
				"next_upgrade_versions":       []string(nil),
				"nodes_url":                   "",
				"private_network_id":          "",
				"region":                      "",
				"status":                      "",
				"update_policy":               "",
				"url":                         "",
				"version":                     "1.0",
				kubeClusterProxyModeKey:       "",
			},
		},
		{
			name: "Expected apiserver customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1.0.0",
				Customization: Customization{
					APIServer: &APIServer{
						AdmissionPlugins: &AdmissionPlugins{
							Enabled:  pointerArray([]string{"foo"}),
							Disabled: pointerArray([]string{"bar"}),
						},
					},
					KubeProxy: nil,
				},
				KubeProxyMode: "",
			},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"name":                        "",
				"next_upgrade_versions":       []string(nil),
				"nodes_url":                   "",
				"private_network_id":          "",
				"region":                      "",
				"status":                      "",
				"update_policy":               "",
				"url":                         "",
				"version":                     "1.0",
				kubeClusterProxyModeKey:       "",
				"customization_apiserver": []map[string]interface{}{
					{
						"admissionplugins": []map[string]interface{}{
							{
								"enabled":  pointerArray([]string{"foo"}),
								"disabled": pointerArray([]string{"bar"}),
							},
						},
					},
				},
			},
		},
		{
			name: "IPTables customization with one field",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1.0.0",
				Customization: Customization{
					APIServer: nil,
					KubeProxy: &kubeProxyCustomization{
						IPTables: &kubeProxyCustomizationIPTables{
							MinSyncPeriod: strPtr("PT30S"),
						},
						IPVS: nil,
					},
				},
				KubeProxyMode: "iptables",
			},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"name":                        "",
				"next_upgrade_versions":       []string(nil),
				"nodes_url":                   "",
				"private_network_id":          "",
				"region":                      "",
				"status":                      "",
				"update_policy":               "",
				"url":                         "",
				"version":                     "1.0",
				kubeClusterProxyModeKey:       "iptables",
				"customization_kube_proxy": []map[string]interface{}{
					{
						"iptables": []map[string]interface{}{
							{
								"min_sync_period": strPtr("PT30S"),
							},
						},
					},
				},
			},
		},
		{
			name: "IPTables customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1.0.0",
				Customization: Customization{
					APIServer: nil,
					KubeProxy: &kubeProxyCustomization{
						IPTables: &kubeProxyCustomizationIPTables{
							MinSyncPeriod: strPtr("PT30S"),
							SyncPeriod:    strPtr("PT30S"),
						},
						IPVS: nil,
					},
				},
				KubeProxyMode: "iptables",
			},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"name":                        "",
				"next_upgrade_versions":       []string(nil),
				"nodes_url":                   "",
				"private_network_id":          "",
				"region":                      "",
				"status":                      "",
				"update_policy":               "",
				"url":                         "",
				"version":                     "1.0",
				kubeClusterProxyModeKey:       "iptables",
				"customization_kube_proxy": []map[string]interface{}{
					{
						"iptables": []map[string]interface{}{
							{
								"min_sync_period": strPtr("PT30S"),
								"sync_period":     strPtr("PT30S"),
							},
						},
					},
				},
			},
		},
		{
			name: "IPVS customization with one field",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1.0.0",
				Customization: Customization{
					APIServer: nil,
					KubeProxy: &kubeProxyCustomization{
						IPTables: nil,
						IPVS: &kubeProxyCustomizationIPVS{
							MinSyncPeriod: strPtr("PT30S"),
						},
					},
				},
				KubeProxyMode: "ipvs",
			},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"name":                        "",
				"next_upgrade_versions":       []string(nil),
				"nodes_url":                   "",
				"private_network_id":          "",
				"region":                      "",
				"status":                      "",
				"update_policy":               "",
				"url":                         "",
				"version":                     "1.0",
				kubeClusterProxyModeKey:       "ipvs",
				"customization_kube_proxy": []map[string]interface{}{
					{
						"ipvs": []map[string]interface{}{
							{
								"min_sync_period": strPtr("PT30S"),
							},
						},
					},
				},
			},
		},
		{
			name: "IPVS customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1.0.0",
				Customization: Customization{
					APIServer: nil,
					KubeProxy: &kubeProxyCustomization{
						IPTables: nil,
						IPVS: &kubeProxyCustomizationIPVS{
							MinSyncPeriod: strPtr("PT30S"),
							SyncPeriod:    strPtr("PT30S"),
							Scheduler:     strPtr("rr"),
							TCPFinTimeout: strPtr("PT30S"),
							TCPTimeout:    strPtr("PT30S"),
							UDPTimeout:    strPtr("PT30S"),
						},
					},
				},
				KubeProxyMode: "ipvs",
			},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"name":                        "",
				"next_upgrade_versions":       []string(nil),
				"nodes_url":                   "",
				"private_network_id":          "",
				"region":                      "",
				"status":                      "",
				"update_policy":               "",
				"url":                         "",
				"version":                     "1.0",
				kubeClusterProxyModeKey:       "ipvs",
				"customization_kube_proxy": []map[string]interface{}{
					{
						"ipvs": []map[string]interface{}{
							{
								"min_sync_period": strPtr("PT30S"),
								"sync_period":     strPtr("PT30S"),
								"scheduler":       strPtr("rr"),
								"tcp_fin_timeout": strPtr("PT30S"),
								"tcp_timeout":     strPtr("PT30S"),
								"udp_timeout":     strPtr("PT30S"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CloudProjectKubeResponse{
				ControlPlaneIsUpToDate: tt.fields.ControlPlaneIsUpToDate,
				Id:                     tt.fields.Id,
				IsUpToDate:             tt.fields.IsUpToDate,
				Name:                   tt.fields.Name,
				NextUpgradeVersions:    tt.fields.NextUpgradeVersions,
				NodesUrl:               tt.fields.NodesUrl,
				PrivateNetworkId:       tt.fields.PrivateNetworkId,
				Region:                 tt.fields.Region,
				Status:                 tt.fields.Status,
				UpdatePolicy:           tt.fields.UpdatePolicy,
				Url:                    tt.fields.Url,
				Version:                tt.fields.Version,
				Customization:          tt.fields.Customization,
				KubeProxyMode:          tt.fields.KubeProxyMode,
			}

			got := v.ToMap()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ToMap() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
