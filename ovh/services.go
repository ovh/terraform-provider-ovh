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

func serviceIDFromServiceNameInQuery(c *ovh.Client, serviceName string) (int, error) {
	var (
		endpoint = "/services?resourceName=" + url.QueryEscape(serviceName)
		ids      []int
	)

	if err := c.Get(endpoint, &ids); err != nil {
		return 0, fmt.Errorf("failed to retrieve service name: %w", err)
	}

	if len(ids) != 1 {
		return 0, fmt.Errorf("invalid number of services retrieved, expected 1 got %d", len(ids))
	}

	return ids[0], nil
}
