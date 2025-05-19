package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

type savingsPlansSubscribable struct {
	OfferID string `json:"offerId"`
}

type savingsPlansSimulateRequest struct {
	DisplayName string `json:"displayName"`
	OfferID     string `json:"offerId"`
	Size        int    `json:"size"`
}

type savingsPlansSimulateResponse struct {
	DisplayName string `json:"displayName"`
	OfferID     string `json:"offerId"`
	Size        int    `json:"size"`

	Flavor          string `json:"flavor"`
	ID              string `json:"id"`
	Period          string `json:"period"`
	PeriodEndAction string `json:"periodEndAction"`
	PeriodStartDate string `json:"periodStartDate"`
	PeriodEndDate   string `json:"periodEndDate"`
	Status          string `json:"status"`

	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	TerminationDate string `json:"terminationDate"`
}

type savingsPlanPeriodEndActionRequest struct {
	PeriodEndAction string `json:"periodEndAction"`
}

func serviceIdFromResourceName(c *ovhwrap.Client, resourceName string) (int, error) {
	var serviceIds []int
	endpoint := fmt.Sprintf("/services?resourceName=%s", url.QueryEscape(resourceName))

	if err := c.Get(endpoint, &serviceIds); err != nil {
		return 0, fmt.Errorf("failed to get service infos: %w", err)
	}

	return serviceIds[0], nil
}

func serviceIdFromRouteAndResourceName(c *ovhwrap.Client, route, resourceName string) (int, error) {
	var serviceIds []int
	endpoint := fmt.Sprintf("/services?resourceName=%s&routes=%s", url.QueryEscape(resourceName), url.QueryEscape(route))

	if err := c.Get(endpoint, &serviceIds); err != nil {
		return 0, fmt.Errorf("failed to get service infos: %w", err)
	}

	return serviceIds[0], nil
}

func serviceInfoFromServiceName(c *ovhwrap.Client, serviceType, serviceName string) (*ServiceInfos, error) {
	var (
		serviceInfos ServiceInfos
		endpoint     = path.Join("/", serviceType, url.PathEscape(serviceName), "/serviceInfos")
	)

	if err := c.Get(endpoint, &serviceInfos); err != nil {
		return nil, fmt.Errorf("failed to get service infos: %w", err)
	}

	return &serviceInfos, nil
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

func serviceUpdateDisplayNameAPIv2(config *Config, serviceName string, displayName string, diagnostics *diag.Diagnostics) error {
	serviceId, err := serviceIdFromResourceName(config.OVHClient, serviceName)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("Error locating service %s", serviceName),
			err.Error(),
		)
		return err
	}

	endpoint := fmt.Sprintf("/services/%d", serviceId)
	if err := config.OVHClient.Put(endpoint, &ServiceUpdatePayload{
		DisplayName: displayName,
	}, nil); err != nil {
		log.Printf("[WARN] update failed : %v", err)
		diagnostics.AddError(
			fmt.Sprintf("Failed to update display name for service %d", serviceId),
			err.Error(),
		)
		return err
	}

	return nil
}
