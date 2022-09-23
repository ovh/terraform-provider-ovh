package ovh

import "fmt"

// getKubeconfig call the kubeconfig endpoint to retreive the kube config file
func getKubeconfig(config *Config, serviceName string, kubeID string) (*string, error) {
	kubeconfigRaw := CloudProjectKubeKubeConfigResponse{}
	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/kubeconfig", serviceName, kubeID)
	err := config.OVHClient.Post(endpoint, nil, &kubeconfigRaw)
	if err != nil {
		return nil, err
	}
	return &kubeconfigRaw.Content, nil
}
