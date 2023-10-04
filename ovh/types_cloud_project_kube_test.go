package ovh

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCloudProjectKubeResponse_ToMap(t *testing.T) {
	t.Skipf("Skipped as we need to create a *schema.ResourceData to test this function.")

	type fields struct {
		ControlPlaneIsUpToDate bool
		Id                     string
		IsUpToDate             bool
		LoadBalancersSubnetId  string
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

	type args struct {
		d *schema.ResourceData
	}

	pointerArray := func(s []string) *[]string { return &s }

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "No customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "",
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
			name: "Deprecated expected apiserver customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "",
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
				"customization": []map[string]interface{}{
					{
						"apiserver": []map[string]interface{}{
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
			},
		},
		{
			name: "Expected apiserver customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "",
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
				LoadBalancersSubnetId:  "",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "",
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
				LoadBalancersSubnetId:  "",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "",
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
				LoadBalancersSubnetId:  "",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "",
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
				LoadBalancersSubnetId:  "",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "",
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
		{
			name: "loadBalancersSubnetId",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "123e4567-e89b-12d3-a456-426614174000",
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
			args: args{},
			want: map[string]interface{}{
				"control_plane_is_up_to_date": false,
				"id":                          "",
				"is_up_to_date":               false,
				"load_balancers_subnet_id":    "123e4567-e89b-12d3-a456-426614174000",
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
				LoadBalancersSubnetId:  tt.fields.LoadBalancersSubnetId,
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

			got := v.ToMap(tt.args.d)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ToMap() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_readApiServerAdmissionPlugins(t *testing.T) {
	type args struct {
		admissionPlugins map[string]interface{}
		apiServerOutput  *APIServer
	}

	pointerArray := func(s []string) *[]string { return &s }

	tests := []struct {
		name string
		args args
		want *APIServer
	}{
		{
			name: "expected admission plugins",
			args: args{
				admissionPlugins: map[string]interface{}{
					"enabled":  []interface{}{"foo", "bar"},
					"disabled": []interface{}{"baz"},
				},
				apiServerOutput: &APIServer{
					AdmissionPlugins: &AdmissionPlugins{},
				},
			},
			want: &APIServer{
				AdmissionPlugins: &AdmissionPlugins{
					Enabled:  pointerArray([]string{"foo", "bar"}),
					Disabled: pointerArray([]string{"baz"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readApiServerAdmissionPlugins(tt.args.admissionPlugins, tt.args.apiServerOutput)
			if diff := cmp.Diff(tt.want, tt.args.apiServerOutput); diff != "" {
				t.Errorf("readApiServerAdmissionPlugins() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
