package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
)

func serviceIdFromResourceName(c *ovh.Client, resourceName string) (int, error) {
	var serviceIds []int
	endpoint := fmt.Sprintf("/services?resourceName=%s", url.PathEscape(resourceName))

	if err := c.Get(endpoint, &serviceIds); err != nil {
		return 0, fmt.Errorf("failed to get service infos: %w", err)
	}

	return serviceIds[0], nil
}

func serviceInfoFromServiceName(c *ovh.Client, serviceType, serviceName string) (*ServiceInfos, error) {
	var (
		serviceInfos ServiceInfos
		endpoint     = path.Join("/", serviceType, url.PathEscape(serviceName), "/serviceInfos")
	)

	if err := c.Get(endpoint, &serviceInfos); err != nil {
		return nil, fmt.Errorf("failed to get service infos: %w", err)
	}

	return &serviceInfos, nil
}

func serviceFromServiceName(c *ovh.Client, serviceType, serviceName string) (*Service, error) {
	serviceInfo, err := serviceInfoFromServiceName(c, serviceType, serviceName)
	if err != nil {
		return nil, err
	}

	var service Service
	if err := c.Get(fmt.Sprintf("/services/%d", serviceInfo.ServiceID), &service); err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return &service, nil
}

// serviceUpdateDisplayName allows to update the display name of any service.
// It first retrieves the ID of the service using route "/${serviceType}/${serviceName}/serviceInfos", and
// then uses this ID to call PUT /services/${serviceId}.
// It finally calls route "/${serviceType}/${serviceName}" to verify that the display name in field "iam" has been updated.
func serviceUpdateDisplayName(ctx context.Context, config *Config, serviceType, serviceName, displayName string) error {
	serviceInfo, err := serviceInfoFromServiceName(config.OVHClient, serviceType, serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service info: %w", err)
	}

	endpoint := fmt.Sprintf("/services/%d", serviceInfo.ServiceID)
	if err := config.OVHClient.PutWithContext(ctx, endpoint, &ServiceUpdatePayload{
		DisplayName: displayName,
	}, nil); err != nil {
		return fmt.Errorf("failed to update service info: %w", err)
	}

	endpoint = "/" + serviceType + "/" + url.PathEscape(serviceName)

	return retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		resourceObject := GenericServiceWithIAMInjection{}

		if err := config.OVHClient.GetWithContext(ctx, endpoint, &resourceObject); err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to get resource %q: %w", serviceName, err))
		}

		if resourceObject.DisplayName != displayName {
			return retry.RetryableError(errors.New("timeout waiting for displayName to be updated"))
		}

		return nil
	})
}
