package ovh

import (
	"reflect"
	"testing"
)

func Test_parseKubeconfig(t *testing.T) {
	type args struct {
		kubeconfigRaw *CloudProjectKubeKubeConfigResponse
	}

	expectedKubeconfigRaw := `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: Zm9vCg==
    server: https://foo.bar
  name: foo
contexts:
- context:
    cluster: foo
    user: kubernetes-admin-foo
  name: kubernetes-admin@foo
current-context: kubernetes-admin@foo
kind: Config
preferences: {}
users:
- name: kubernetes-admin-foo
  user:
    client-certificate-data: Zm9vCg==
    client-key-data: Zm9vCg==

`

	tests := []struct {
		name    string
		args    args
		want    *KubectlConfig
		wantErr bool
	}{
		{
			name: "expected kubeconfig content",
			args: args{
				kubeconfigRaw: &CloudProjectKubeKubeConfigResponse{
					Content: expectedKubeconfigRaw,
				},
			},
			want: &KubectlConfig{
				Kind:           "Config",
				ApiVersion:     "v1",
				CurrentContext: "kubernetes-admin@foo",
				Clusters: []*KubectlClusterWithName{
					{
						Name:    "foo",
						Cluster: KubectlCluster{Server: "https://foo.bar", CertificateAuthorityData: "Zm9vCg=="},
					},
				},
				Contexts: []*KubectlContextWithName{
					{
						Name:    "kubernetes-admin@foo",
						Context: KubectlContext{Cluster: "foo", User: "kubernetes-admin-foo"},
					},
				},
				Users: []*KubectlUserWithName{
					{
						Name: "kubernetes-admin-foo",
						User: KubectlUser{ClientCertificateData: "Zm9vCg==", ClientKeyData: "Zm9vCg=="}},
				},
				Raw: &expectedKubeconfigRaw,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseKubeconfig(tt.args.kubeconfigRaw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseKubeconfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseKubeconfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
