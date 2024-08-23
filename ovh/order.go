package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/types"
)

var (
	reTerminateEmailToken = regexp.MustCompile(`.*https://www.ovh.com/manager/#/billing/confirmTerminate\?id=[[:alnum:]]+&token=([[:alnum:]]+).*`)
	terminateEmailMatch   = "https://www.ovh.com/manager/#/billing/confirmTerminate"
)

func genericOrderSchema(withOptions bool) map[string]*schema.Schema {
	var planOptionsMaxItems int
	if !withOptions {
		planOptionsMaxItems = 0
	}

	orderSchema := map[string]*schema.Schema{
		"ovh_subsidiary": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Ovh Subsidiary",
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				err := helpers.ValidateSubsidiary(v.(string))
				if err != nil {
					errors = append(errors, err)
				}
				return
			},
		},
		"payment_mean": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Deprecated:  "This field is not anymore used since the API has been deprecated in favor of /payment/mean. Now, the default payment mean is used.",
			Description: "Ovh payment mode",
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				err := helpers.ValidateStringEnum(strings.ToLower(v.(string)), []string{
					"default-payment-mean",
					"fidelity",
					"ovh-account",
				})
				if err != nil {
					errors = append(errors, err)
				}
				return
			},
		},

		"plan": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			Description: "Product Plan to order",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"duration": {
						Type:        schema.TypeString,
						Description: "duration",
						Required:    true,
					},
					"plan_code": {
						Type:        schema.TypeString,
						Description: "Plan code",
						Required:    true,
					},
					"pricing_mode": {
						Type:        schema.TypeString,
						Description: "Pricing model identifier",
						Required:    true,
					},
					"catalog_name": {
						Type:        schema.TypeString,
						Description: "Catalog name",
						Optional:    true,
					},
					"configuration": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Representation of a configuration item for personalizing product",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"label": {
									Type:        schema.TypeString,
									Description: "Identifier of the resource",
									Required:    true,
								},
								"value": {
									Type:        schema.TypeString,
									Description: "Path to the resource in API.OVH.COM",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},

		"plan_option": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    planOptionsMaxItems,
			Description: "Product Plan to order",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"duration": {
						Type:        schema.TypeString,
						Description: "duration",
						Required:    true,
					},
					"plan_code": {
						Type:        schema.TypeString,
						Description: "Plan code",
						Required:    true,
					},
					"pricing_mode": {
						Type:        schema.TypeString,
						Description: "Pricing model identifier",
						Required:    true,
					},
					"catalog_name": {
						Type:        schema.TypeString,
						Description: "Catalog name",
						Optional:    true,
					},
					"configuration": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Representation of a configuration item for personalizing product",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"label": {
									Type:        schema.TypeString,
									Description: "Identifier of the resource",
									Required:    true,
								},
								"value": {
									Type:        schema.TypeString,
									Description: "Path to the resource in API.OVH.COM",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},

		"order": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Details about an Order",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"date": {
						Type:        schema.TypeString,
						Description: "date",
						Computed:    true,
					},
					"order_id": {
						Type:        schema.TypeInt,
						Description: "order id",
						Computed:    true,
					},
					"expiration_date": {
						Type:        schema.TypeString,
						Description: "expiration date",
						Computed:    true,
					},
					"details": {
						Type:        schema.TypeList,
						Computed:    true,
						Description: "Information about a Bill entry",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"description": {
									Type:        schema.TypeString,
									Description: "description",
									Computed:    true,
								},
								"order_detail_id": {
									Type:        schema.TypeInt,
									Description: "order detail id",
									Computed:    true,
								},
								"domain": {
									Type:        schema.TypeString,
									Description: "expiration date",
									Computed:    true,
								},
								"quantity": {
									Type:        schema.TypeString,
									Description: "quantity",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}

	return orderSchema
}

func orderCreateFromResource(d *schema.ResourceData, meta interface{}, product string, waitForCompletion bool) error {
	config := meta.(*Config)
	order := (&OrderModel{}).FromResource(d)

	err := orderCreate(order, config, product, waitForCompletion)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(order.Order.OrderId.ValueInt64()))

	return nil
}

func orderCreate(d *OrderModel, config *Config, product string, waitForCompletion bool) error {
	// create Cart
	cartParams := &OrderCartCreateOpts{
		OvhSubsidiary: strings.ToUpper(d.OvhSubsidiary.ValueString()),
	}

	cart, err := orderCartCreate(config, cartParams, true)
	if err != nil {
		return fmt.Errorf("calling creating order cart: %q", err)
	}

	// Create Product Item
	item := &OrderCartItem{}
	cartPlanParamsList := d.Plan.Elements()
	cartPlanParams := cartPlanParamsList[0].(PlanValue)
	cartPlanParams.Quantity = types.TfInt64Value{Int64Value: basetypes.NewInt64Value(1)}

	log.Printf("[DEBUG] Will create order item %s for cart: %s", product, cart.CartId)
	endpoint := fmt.Sprintf("/order/cart/%s/%s", url.PathEscape(cart.CartId), product)
	if err := config.OVHClient.Post(endpoint, cartPlanParams, item); err != nil {
		return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, cartPlanParams, err)
	}

	// apply configurations
	configs := cartPlanParams.Configuration.Elements()

	for _, cfg := range configs {
		log.Printf("[DEBUG] Will create order cart item configuration for cart item: %s/%d",
			item.CartId,
			item.ItemId,
		)
		itemConfig := &OrderCartItemConfiguration{}
		endpoint := fmt.Sprintf("/order/cart/%s/item/%d/configuration",
			url.PathEscape(item.CartId),
			item.ItemId,
		)
		if err := config.OVHClient.Post(endpoint, cfg, itemConfig); err != nil {
			return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, cfg, err)
		}
	}

	planOptionValue := d.PlanOption.Elements()

	// Create Product Options Items
	for _, option := range planOptionValue {
		opt := option.(PlanOptionValue)

		log.Printf("[DEBUG] Will create order item options %s for cart: %s", product, cart.CartId)
		productOptionsItem := &OrderCartItem{}

		opt.ItemId = types.TfInt64Value{Int64Value: basetypes.NewInt64Value(item.ItemId)}
		opt.Quantity = types.TfInt64Value{Int64Value: basetypes.NewInt64Value(1)}

		endpoint := fmt.Sprintf("/order/cart/%s/%s/options", url.PathEscape(cart.CartId), product)
		if err := config.OVHClient.Post(endpoint, opt, productOptionsItem); err != nil {
			return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, cartPlanParams, err)
		}

		optionConfigs := opt.Configuration.Elements()
		for _, cfg := range optionConfigs {
			log.Printf("[DEBUG] Will create order cart item configuration for cart item: %s/%d",
				item.CartId,
				item.ItemId,
			)
			itemConfig := &OrderCartItemConfiguration{}
			endpoint := fmt.Sprintf("/order/cart/%s/item/%d/configuration",
				url.PathEscape(item.CartId),
				item.ItemId,
			)
			if err := config.OVHClient.Post(endpoint, cfg, itemConfig); err != nil {
				return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, cfg, err)
			}
		}
	}

	// get defaultPayment
	paymentIds := []int64{}
	endpoint = "/me/payment/method?default=true"
	if err := config.OVHClient.Get(endpoint, &paymentIds); err != nil {
		return fmt.Errorf("calling Get %s \n\t %q", endpoint, err)
	}

	fallbackToFidelityAccount := false
	if len(paymentIds) == 0 {
		// check if amount of the order is 0€, we can try to use fidelity account to pay order
		log.Printf("[DEBUG] Will create order %s for cart: %s", product, cart.CartId)
		checkout := &OrderCartCheckout{}

		endpoint = fmt.Sprintf("/order/cart/%s/checkout", url.PathEscape(cart.CartId))
		if err := config.OVHClient.Get(endpoint, checkout); err != nil {
			return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
		}

		if checkout.Prices.WithoutTax.Value == 0 {
			fallbackToFidelityAccount = true

		} else {
			return fmt.Errorf("no default payment found")
		}
	}

	// Create Order
	log.Printf("[DEBUG] Will create order %s for cart: %s", product, cart.CartId)
	checkout := &OrderCartCheckout{}

	endpoint = fmt.Sprintf("/order/cart/%s/checkout", url.PathEscape(cart.CartId))
	if err := config.OVHClient.Post(endpoint, nil, checkout); err != nil {
		return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
	}

	// Pay Order
	if !fallbackToFidelityAccount {
		log.Printf("[DEBUG] Will pay order %d with PaymentId %d", checkout.OrderID, paymentIds[0])

		endpoint = fmt.Sprintf(
			"/me/order/%d/pay",
			checkout.OrderID,
		)
		var paymentMethodOpts = &MeOrderPaymentMethodOpts{
			PaymentMethod: PaymentMethod{
				Id: paymentIds[0],
			},
		}
		if err := config.OVHClient.Post(endpoint, paymentMethodOpts, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}
	} else {
		log.Printf("[DEBUG] Will pay free order %d with fidelityAccount", checkout.OrderID)
		endpoint = fmt.Sprintf(
			"/me/order/%d/payWithRegisteredPaymentMean",
			checkout.OrderID,
		)
		var paymentMethodOpts = &MeOrderPaymentOpts{
			PaymentMean: "fidelityAccount",
		}
		if err := config.OVHClient.Post(endpoint, paymentMethodOpts, nil); err != nil {
			return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
		}

	}

	// Wait for order to be completed
	if waitForCompletion {
		if err := waitOrderCompletion(config, checkout.OrderID); err != nil {
			return fmt.Errorf("waiting for order (%d): %s", checkout.OrderID, err)
		}
	}

	d.Order.OrderId = types.TfInt64Value{Int64Value: basetypes.NewInt64Value(checkout.OrderID)}

	return nil
}

func waitOrderCompletion(config *Config, orderID int64) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"checking", "delivering", "ignoreerror"},
		Target:     []string{"delivered"},
		Refresh:    waitForOrder(config.OVHClient, orderID),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func orderReadInResource(d *schema.ResourceData, meta interface{}) (*MeOrder, []*MeOrderDetail, error) {
	config := meta.(*Config)
	orderId := d.Id()

	order, details, err := orderRead(orderId, config)
	if err != nil {
		return nil, nil, err
	}

	detailsData := make([]map[string]interface{}, len(details))
	for i, detail := range details {
		detailsData[i] = detail.ToMap()
	}

	orderData := order.ToMap()
	orderData["details"] = detailsData
	d.Set("order", []interface{}{orderData})

	return order, details, nil
}

func orderRead(orderId string, config *Config) (*MeOrder, []*MeOrderDetail, error) {
	order := &MeOrder{}
	log.Printf("[DEBUG] Will read order %s", orderId)
	endpoint := fmt.Sprintf("/me/order/%s",
		url.PathEscape(orderId),
	)
	if err := config.OVHClient.Get(endpoint, &order); err != nil {
		return nil, nil, fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
	}

	details, err := orderDetails(config.OVHClient, order.OrderId)
	if err != nil {
		return nil, nil, err
	}

	if len(details) < 1 {
		return nil, nil, fmt.Errorf("there is no order details for id %s. This shouldn't happen. This is a bug with the API", orderId)
	}

	return order, details, nil
}

type TerminateFunc func() (string, error)
type ConfirmTerminationFunc func(token string) error

func orderDeleteFromResource(d *schema.ResourceData, meta interface{}, terminate TerminateFunc, confirm ConfirmTerminationFunc) error {
	config := meta.(*Config)

	if err := orderDelete(config, terminate, confirm); err != nil {
		return err
	}

	if d != nil {
		d.SetId("")
	}

	return nil
}

func orderDelete(config *Config, terminate TerminateFunc, confirm ConfirmTerminationFunc) error {
	oldEmailsIds, err := notificationEmailSortedIds(config)
	if err != nil {
		return err
	}

	match, err := terminate()
	if err != nil {
		return err
	}

	if match == "" {
		log.Printf("[INFO] nothing to terminate or service already suspended.")
		return nil
	}

	matches := []string{
		match,
		terminateEmailMatch,
	}

	var email *NotificationEmail
	// wait for email
	err = resource.Retry(30*time.Minute, func() *resource.RetryError {
		email, err = getNewNotificationEmail(matches, oldEmailsIds, config)
		if err != nil {
			log.Printf("[DEBUG] error while getting email notification. retry: %v", err)
			return resource.RetryableError(err)
		}

		if email == nil {
			return resource.RetryableError(fmt.Errorf("email notification not found"))
		}

		// Successful cascade delete
		log.Printf("[DEBUG] successfully found termination email for %s with id %d", match, email.Id)
		return nil
	})

	if err != nil {
		return fmt.Errorf("email notification not found: %v", err)
	}

	tokenMatch := reTerminateEmailToken.FindStringSubmatch(email.Body)
	if len(tokenMatch) != 2 {
		return fmt.Errorf("could not find termination token in email notification: %v", email.Id)
	}

	if err := confirm(tokenMatch[1]); err != nil {
		return err
	}

	return nil
}

func orderDetails(c *ovh.Client, orderId int64) ([]*MeOrderDetail, error) {
	log.Printf("[DEBUG] Will read order details %d", orderId)
	detailIds := []int64{}
	endpoint := fmt.Sprintf("/me/order/%d/details", orderId)
	if err := c.Get(endpoint, &detailIds); err != nil {
		return nil, fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
	}

	details := make([]*MeOrderDetail, len(detailIds))
	for i, detailId := range detailIds {
		detail := &MeOrderDetail{}
		log.Printf("[DEBUG] Will read order detail %d/%d", orderId, detailId)
		endpoint := fmt.Sprintf("/me/order/%d/details/%d", orderId, detailId)
		if err := c.Get(endpoint, detail); err != nil {
			return nil, fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
		}

		details[i] = detail
	}
	return details, nil
}

func serviceNameFromOrder(c *ovh.Client, orderId int64, plan string) (string, error) {
	detailIds := []int64{}
	endpoint := fmt.Sprintf("/me/order/%d/details", orderId)
	if err := c.Get(endpoint, &detailIds); err != nil {
		return "", fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
	}

	for _, detailId := range detailIds {
		detailExtension := &MeOrderDetailExtension{}
		log.Printf("[DEBUG] Will read order detail extension %d/%d", orderId, detailId)
		endpoint := fmt.Sprintf("/me/order/%d/details/%d/extension", orderId, detailId)
		if err := c.Get(endpoint, detailExtension); err != nil {
			return "", fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
		}

		if detailExtension.Order.Plan.Code != plan {
			continue
		}

		detail := &MeOrderDetail{}
		log.Printf("[DEBUG] Will read order detail %d/%d", orderId, detailId)
		endpoint = fmt.Sprintf("/me/order/%d/details/%d", orderId, detailId)
		if err := c.Get(endpoint, detail); err != nil {
			return "", fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
		}

		return detail.Domain, nil
	}

	return "", errors.New("serviceName not found")
}

func waitForOrder(c *ovh.Client, id int64) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var r string
		endpoint := fmt.Sprintf("/me/order/%d/status", id)
		if err := c.Get(endpoint, &r); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				log.Printf("[DEBUG] order id %d deleted", id)
				return nil, "deleted", nil
			}

			log.Printf("[WARNING] order id %d ignore error: %v", id, err)
			return nil, "ignoreerror", nil
		}

		log.Printf("[DEBUG] Pending order: %s", r)
		return r, r, nil
	}
}

func waitOrderCompletionV2(ctx context.Context, config *Config, orderID int64) (string, error) {
	endpoint := fmt.Sprintf("/me/order/%d/status", orderID)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"checking", "delivering"},
		Target:  []string{"delivered"},
		Refresh: func() (interface{}, string, error) {
			var status string
			if err := config.OVHClient.Get(endpoint, &status); err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					log.Printf("[DEBUG] order id %d deleted", orderID)
					return nil, "deleted", nil
				}

				log.Printf("[WARNING] order id %d, got error: %v", orderID, err)
				return nil, "error", err
			}

			return status, status, nil
		},
		Timeout:    time.Hour,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "error", err
	}

	return result.(string), err
}

func orderDetailOperations(c *ovh.Client, orderId int64, orderDetailId int64) ([]*MeOrderDetailOperation, error) {
	log.Printf("[DEBUG] Will list order detail operations %d/%d", orderId, orderDetailId)
	operationsIds := []int64{}
	endpoint := fmt.Sprintf("/me/order/%d/details/%d/operations", orderId, orderDetailId)
	if err := c.Get(endpoint, &operationsIds); err != nil {
		return nil, fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
	}

	operations := make([]*MeOrderDetailOperation, len(operationsIds))
	for i, operationId := range operationsIds {
		operation := &MeOrderDetailOperation{}
		log.Printf("[DEBUG] Will read order detail operations %d/%d/%d", orderId, orderDetailId, operationId)
		endpoint := fmt.Sprintf("/me/order/%d/details/%d/operations/%d", orderId, orderDetailId, operationId)
		if err := c.Get(endpoint, operation); err != nil {
			return nil, fmt.Errorf("calling get %s:\n\t %q", endpoint, err)
		}

		operations[i] = operation
	}
	return operations, nil
}
