package ovh

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// KubectlConfig is a struct to store the kubeconfig file
// Same as https://github.com/kubernetes/kops/blob/2e84499741471ba67582aa0ba6fa3f2e3bdbe3e8/pkg/kubeconfig/config.go#L19 but with yaml format
type KubectlConfig struct {
	Kind           string                    `json:"kind" yaml:"kind"`
	ApiVersion     string                    `json:"apiVersion" yaml:"apiVersion"`
	CurrentContext string                    `json:"current-context" yaml:"current-context"`
	Clusters       []*KubectlClusterWithName `json:"clusters" yaml:"clusters"`
	Contexts       []*KubectlContextWithName `json:"contexts" yaml:"contexts"`
	Users          []*KubectlUserWithName    `json:"users" yaml:"users"`
	Raw            *string                   `json:"-" yaml:"-"`
}

type KubectlClusterWithName struct {
	Name    string         `json:"name" yaml:"name"`
	Cluster KubectlCluster `json:"cluster" yaml:"cluster"`
}

type KubectlCluster struct {
	Server                   string `json:"server,omitempty" yaml:"server,omitempty"`
	CertificateAuthorityData string `json:"certificate-authority-data,omitempty" yaml:"certificate-authority-data,omitempty"`
}

type KubectlContextWithName struct {
	Name    string         `json:"name" yaml:"name"`
	Context KubectlContext `json:"context" yaml:"context"`
}

type KubectlContext struct {
	Cluster string `json:"cluster" yaml:"cluster"`
	User    string `json:"user" yaml:"user"`
}

type KubectlUserWithName struct {
	Name string      `json:"name" yaml:"name"`
	User KubectlUser `json:"user" yaml:"user"`
}

type KubectlUser struct {
	ClientCertificateData string `json:"client-certificate-data,omitempty" yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `json:"client-key-data,omitempty" yaml:"client-key-data,omitempty"`
	Password              string `json:"password,omitempty" yaml:"password,omitempty"`
	Username              string `json:"username,omitempty" yaml:"username,omitempty"`
	Token                 string `json:"token,omitempty" yaml:"token,omitempty"`
}

// getKubeconfig call the kubeconfig endpoint to retrieve the kube config file
func getKubeconfig(config *Config, serviceName string, kubeID string) (*KubectlConfig, error) {
	kubeconfigRaw := CloudProjectKubeKubeConfigResponse{}
	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/kubeconfig", serviceName, kubeID)
	err := config.OVHClient.Post(endpoint, nil, &kubeconfigRaw)
	if err != nil {
		return nil, err
	}

	return parseKubeconfig(&kubeconfigRaw)
}

func parseKubeconfig(kubeconfigRaw *CloudProjectKubeKubeConfigResponse) (*KubectlConfig, error) {
	var kubeconfig KubectlConfig
	if err := yaml.Unmarshal([]byte(kubeconfigRaw.Content), &kubeconfig); err != nil {
		return nil, err
	}

	kubeconfig.Raw = &kubeconfigRaw.Content
	return &kubeconfig, nil
}
