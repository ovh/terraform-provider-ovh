package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		NodesSubnetId          string
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
			name: "No customization (with K8s major version)",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "",
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesSubnetId:          "",
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1",
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
				"nodes_subnet_id":             "",
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
			name: "No customization (with K8s minor version)",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "",
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesSubnetId:          "",
				NodesUrl:               "",
				PrivateNetworkId:       "",
				Region:                 "",
				Status:                 "",
				UpdatePolicy:           "",
				Url:                    "",
				Version:                "1.0",
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
				"nodes_subnet_id":             "",
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
			name: "No customization",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "",
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
				NodesSubnetId:          "",
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
				"nodes_subnet_id":             "",
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
			name: "nodessSubnetId",
			fields: fields{
				ControlPlaneIsUpToDate: false,
				Id:                     "",
				IsUpToDate:             false,
				LoadBalancersSubnetId:  "123e4567-e89b-12d3-a456-426614174000",
				Name:                   "",
				NextUpgradeVersions:    nil,
				NodesSubnetId:          "777e4567-e89b-12d3-a456-426614174000",
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
				"nodes_subnet_id":             "777e4567-e89b-12d3-a456-426614174000",
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
				NodesSubnetId:          tt.fields.NodesSubnetId,
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

func TestCloudProjectKubeCreateOpts_FromResource_DisablePublicNetworking(t *testing.T) {
	boolPtr := func(b bool) *bool { return &b }

	tests := []struct {
		name string
		raw  map[string]interface{}
		want *bool
	}{
		{name: "unset is sent as null", raw: map[string]interface{}{}},
		{name: "explicit false is sent", raw: map[string]interface{}{kubeClusterDisablePublicNetworkingKey: false}, want: boolPtr(false)},
		{name: "explicit true is sent", raw: map[string]interface{}{kubeClusterDisablePublicNetworkingKey: true}, want: boolPtr(true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.raw[kubeRegionKey] = "GRA11"
			d := schema.TestResourceDataRaw(t, resourceCloudProjectKube().Schema, tt.raw)

			opts := &CloudProjectKubeCreateOpts{}
			opts.FromResource(d)

			if diff := cmp.Diff(tt.want, opts.DisablePublicNetworking); diff != "" {
				t.Errorf("DisablePublicNetworking mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_loadPrivateNetworkConfiguration(t *testing.T) {
	pncType := cty.Object(map[string]cty.Type{
		kubeClusterDefaultVrackGatewayKey:         cty.String,
		kubeClusterPrivateNetworkRoutingAsDefault: cty.Bool,
		kubeClusterUseUserGatewayKey:              cty.Bool,
	})
	rawConfig := func(pnc cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{kubeClusterPrivateNetworkConfigurationKey: pnc})
	}
	pnc := func(defaultVrackGateway, routingAsDefault, useUserGateway cty.Value) cty.Value {
		return cty.SetVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			kubeClusterDefaultVrackGatewayKey:         defaultVrackGateway,
			kubeClusterPrivateNetworkRoutingAsDefault: routingAsDefault,
			kubeClusterUseUserGatewayKey:              useUserGateway,
		})})
	}
	nullString, nullBool := cty.NullVal(cty.String), cty.NullVal(cty.Bool)

	tests := []struct {
		name      string
		rawConfig cty.Value
		wantJSON  string
	}{
		{name: "no raw config", rawConfig: cty.NullVal(cty.EmptyObject), wantJSON: "null"},
		{name: "no private_network_configuration block", rawConfig: rawConfig(cty.NullVal(cty.Set(pncType))), wantJSON: "null"},
		{name: "empty block", rawConfig: rawConfig(pnc(nullString, nullBool, nullBool)), wantJSON: `{}`},
		{name: "only use_user_gateway", rawConfig: rawConfig(pnc(nullString, nullBool, cty.True)), wantJSON: `{"useUserGateway":true}`},
		{
			name:      "explicit zero values are sent",
			rawConfig: rawConfig(pnc(cty.StringVal(""), cty.False, cty.False)),
			wantJSON:  `{"defaultVrackGateway":"","privateNetworkRoutingAsDefault":false,"useUserGateway":false}`,
		},
		{
			name:      "all attributes set",
			rawConfig: rawConfig(pnc(cty.StringVal("10.0.0.1"), cty.True, cty.True)),
			wantJSON:  `{"defaultVrackGateway":"10.0.0.1","privateNetworkRoutingAsDefault":true,"useUserGateway":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := resourceCloudProjectKube().Data(&terraform.InstanceState{RawConfig: tt.rawConfig})

			got, err := json.Marshal(loadPrivateNetworkConfiguration(d))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.wantJSON, string(got)); diff != "" {
				t.Errorf("request body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCloudProjectKube_PrivateNetworkConfigurationOptionalAttributes(t *testing.T) {
	raw := map[string]interface{}{
		kubeServiceNameKey: "xxx",
		kubeRegionKey:      "GRA11",
		kubeClusterPrivateNetworkConfigurationKey: []interface{}{
			map[string]interface{}{kubeClusterUseUserGatewayKey: true},
		},
	}

	if diags := resourceCloudProjectKube().Validate(terraform.NewResourceConfigRaw(raw)); diags.HasError() {
		t.Fatalf("unexpected validation error: %v", diags)
	}
}

func TestCloudProjectKube_UseUserGatewayDiff(t *testing.T) {
	pnc := func(useUserGateway bool, defaultVrackGateway string) []interface{} {
		return []interface{}{map[string]interface{}{
			kubeClusterUseUserGatewayKey:      useUserGateway,
			kubeClusterDefaultVrackGatewayKey: defaultVrackGateway,
		}}
	}

	for _, tc := range []struct {
		name        string
		statePNC    []interface{}
		configPNC   []interface{}
		requiresNew bool
	}{
		{name: "no change", statePNC: pnc(true, "10.0.0.1"), configPNC: pnc(true, "10.0.0.1"), requiresNew: false},
		{name: "toggle false to true", statePNC: pnc(false, ""), configPNC: pnc(true, ""), requiresNew: true},
		{name: "toggle true to false", statePNC: pnc(true, ""), configPNC: pnc(false, ""), requiresNew: true},
		{name: "add block with true", statePNC: nil, configPNC: pnc(true, ""), requiresNew: true},
		{name: "remove block with true", statePNC: pnc(true, ""), configPNC: nil, requiresNew: true},
		{name: "update default_vrack_gateway", statePNC: pnc(true, "10.0.0.1"), configPNC: pnc(true, "10.0.0.2"), requiresNew: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceCloudProjectKube()

			d := r.TestResourceData()
			d.SetId("kube-id")
			for k, v := range map[string]interface{}{
				kubeServiceNameKey:                        "xxx",
				kubeRegionKey:                             "GRA11",
				kubeClusterDisablePublicNetworkingKey:     true,
				kubeClusterPrivateNetworkConfigurationKey: tc.statePNC,
			} {
				if err := d.Set(k, v); err != nil {
					t.Fatal(err)
				}
			}

			raw := map[string]interface{}{
				kubeServiceNameKey:                    "xxx",
				kubeRegionKey:                         "GRA11",
				kubeClusterDisablePublicNetworkingKey: true,
			}
			if tc.configPNC != nil {
				raw[kubeClusterPrivateNetworkConfigurationKey] = tc.configPNC
			}

			diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(raw), nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff.RequiresNew() != tc.requiresNew {
				t.Errorf("RequiresNew() = %v, want %v (diff: %v)", diff.RequiresNew(), tc.requiresNew, diff)
			}
		})
	}
}

func TestCloudProjectKubeResponse_ToMap_NetworkingFlags(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantDPN bool
		wantPNC interface{}
	}{
		{
			name:    "flags returned by the API",
			body:    `{"version":"1.32.1","disablePublicNetworking":true,"privateNetworkConfiguration":{"defaultVrackGateway":"","privateNetworkRoutingAsDefault":false,"useUserGateway":true}}`,
			wantDPN: true,
			wantPNC: []interface{}{map[string]interface{}{
				kubeClusterDefaultVrackGatewayKey:         "",
				kubeClusterPrivateNetworkRoutingAsDefault: false,
				kubeClusterUseUserGatewayKey:              true,
			}},
		},
		{
			name:    "null useUserGateway is read as false",
			body:    `{"version":"1.32.1","disablePublicNetworking":false,"privateNetworkConfiguration":{"defaultVrackGateway":"10.0.0.1","privateNetworkRoutingAsDefault":true,"useUserGateway":null}}`,
			wantDPN: false,
			wantPNC: []interface{}{map[string]interface{}{
				kubeClusterDefaultVrackGatewayKey:         "10.0.0.1",
				kubeClusterPrivateNetworkRoutingAsDefault: true,
				kubeClusterUseUserGatewayKey:              false,
			}},
		},
		{
			name: "no private network configuration",
			body: `{"version":"1.32.1"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &CloudProjectKubeResponse{}
			if err := json.Unmarshal([]byte(tt.body), res); err != nil {
				t.Fatal(err)
			}

			obj := res.ToMap(nil)

			if diff := cmp.Diff(tt.wantDPN, obj[kubeClusterDisablePublicNetworkingKey]); diff != "" {
				t.Errorf("disable_public_networking mismatch (-want +got):\n%s", diff)
			}
			var gotPNC interface{}
			if res.PrivateNetworkConfiguration != nil {
				gotPNC = res.PrivateNetworkConfiguration.ToMap()
			}
			if diff := cmp.Diff(tt.wantPNC, gotPNC); diff != "" {
				t.Errorf("private_network_configuration mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
