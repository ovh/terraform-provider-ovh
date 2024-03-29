package ovh

import (
	"fmt"
	"net/url"
	"path"

	"github.com/ovh/go-ovh/ovh"
)

func serviceFromServiceName(c *ovh.Client, serviceType, serviceName string) (*Service, error) {
	var (
		service      Service
		serviceInfos ServiceInfos
		endpoint     = path.Join("/", serviceType, url.PathEscape(serviceName), "/serviceInfos")
	)

	if err := c.Get(endpoint, &serviceInfos); err != nil {
		return nil, fmt.Errorf("failed to get service infos: %w", err)
	}

	if err := c.Get(fmt.Sprintf("/services/%d", serviceInfos.ServiceID), &service); err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return &service, nil
}
