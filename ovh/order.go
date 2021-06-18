package ovh

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
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
			Required:    true,
			ForceNew:    true,
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

func orderCreate(d *schema.ResourceData, meta interface{}, product string) error {
	config := meta.(*Config)

	// create Cart
	cartParams := &OrderCartCreateOpts{
		OvhSubsidiary: strings.ToUpper(d.Get("ovh_subsidiary").(string)),
	}

	cart, err := orderCartCreate(meta, cartParams, true)
	if err != nil {
		return fmt.Errorf("calling creating order cart: %q", err)
	}

	// Create Product Item
	item := &OrderCartItem{}
	cartPlanParams := (&OrderCartPlanCreateOpts{
		Quantity: 1,
	}).FromResourceWithPath(d, "plan.0")

	log.Printf("[DEBUG] Will create order item %s for cart: %s", product, cart.CartId)
	endpoint := fmt.Sprintf("/order/cart/%s/%s", url.PathEscape(cart.CartId), product)
	if err := config.OVHClient.Post(endpoint, cartPlanParams, item); err != nil {
		return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, cartPlanParams, err)
	}

	// apply configurations
	nbOfConfigurations := d.Get("plan.0.configuration.#").(int)
	for i := 0; i < nbOfConfigurations; i++ {
		log.Printf("[DEBUG] Will create order cart item configuration for cart item: %s/%d",
			item.CartId,
			item.ItemId,
		)
		itemConfig := &OrderCartItemConfiguration{}
		itemConfigParams := (&OrderCartItemConfigurationOpts{}).FromResourceWithPath(
			d,
			fmt.Sprintf("plan.0.configuration.%d", i),
		)
		endpoint := fmt.Sprintf("/order/cart/%s/item/%d/configuration",
			url.PathEscape(item.CartId),
			item.ItemId,
		)
		if err := config.OVHClient.Post(endpoint, itemConfigParams, itemConfig); err != nil {
			return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, itemConfigParams, err)
		}
	}

	// Create Product Options Items
	nbOfOptions := d.Get("plan_option.#").(int)
	for i := 0; i < nbOfOptions; i++ {
		optionPath := fmt.Sprintf("plan_option.%d", i)
		log.Printf("[DEBUG] Will create order item options %s for cart: %s", product, cart.CartId)
		productOptionsItem := &OrderCartItem{}
		cartPlanOptionsParams := (&OrderCartPlanOptionsCreateOpts{
			ItemId:   item.ItemId,
			Quantity: 1,
		}).FromResourceWithPath(
			d,
			optionPath,
		)
		endpoint := fmt.Sprintf("/order/cart/%s/%s/options", url.PathEscape(cart.CartId), product)
		if err := config.OVHClient.Post(endpoint, cartPlanOptionsParams, productOptionsItem); err != nil {
			return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, cartPlanParams, err)
		}

		// apply configurations
		nbOfConfigurations := d.Get(fmt.Sprintf("%s.configuration.#", optionPath)).(int)
		for j := 0; j < nbOfConfigurations; j++ {
			log.Printf("[DEBUG] Will create order cart item configuration for cart item: %s/%d",
				item.CartId,
				item.ItemId,
			)
			itemConfig := &OrderCartItemConfiguration{}
			itemConfigParams := (&OrderCartItemConfigurationOpts{}).FromResourceWithPath(
				d,
				fmt.Sprintf("%s.configuration.%d", optionPath, j),
			)
			endpoint := fmt.Sprintf("/order/cart/%s/item/%d/configuration",
				url.PathEscape(item.CartId),
				item.ItemId,
			)
			if err := config.OVHClient.Post(endpoint, itemConfigParams, itemConfig); err != nil {
				return fmt.Errorf("calling Post %s with params %v:\n\t %q", endpoint, itemConfigParams, err)
			}
		}
	}

	// Create Order
	log.Printf("[DEBUG] Will create order %s for cart: %s", product, cart.CartId)
	order := &MeOrder{}

	endpoint = fmt.Sprintf("/order/cart/%s/checkout", url.PathEscape(cart.CartId))
	if err := config.OVHClient.Post(endpoint, nil, order); err != nil {
		return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
	}

	// Pay Order
	log.Printf("[DEBUG] Will pay order %d", order.OrderId)

	var paymentMeanOpts *MeOrderPaymentOpts
	paymentMean := d.Get("payment_mean").(string)

	switch strings.ToLower(paymentMean) {
	case "default-payment-mean":
		paymentMeanOpts, err = MePaymentMeanDefaultPaymentOpts(config.OVHClient)
		if err != nil {
			return fmt.Errorf("Could not order product: %v.", err)
		}
		if paymentMeanOpts == nil {
			return fmt.Errorf("Could not find any default payment mean to order product.")
		}
	case "ovh-account":
		paymentMeanOpts = MePaymentMeanOvhAccountPaymentOpts
	case "fidelity":
		paymentMeanOpts = MePaymentMeanFidelityAccountPaymentOpts
	default:
		return fmt.Errorf("Unsupported payment mean. This is a bug with the provider.")
	}

	endpoint = fmt.Sprintf(
		"/me/order/%d/payWithRegisteredPaymentMean",
		order.OrderId,
	)
	if err := config.OVHClient.Post(endpoint, paymentMeanOpts, nil); err != nil {
		return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
	}

	// Wait for order status
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"checking", "delivering", "ignoreerror"},
		Target:     []string{"delivered"},
		Refresh:    waitForOrder(config.OVHClient, order.OrderId),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for order (%d): %s", order.OrderId, err)
	}

	d.SetId(fmt.Sprint(order.OrderId))

	return nil
}

func orderRead(d *schema.ResourceData, meta interface{}) (*MeOrder, []*MeOrderDetail, error) {
	config := meta.(*Config)
	orderId := d.Id()

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
		return nil, nil, fmt.Errorf("There is no order details for id %s. This shouldn't happen. This is a bug with the API.", orderId)
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

type TerminateFunc func() (string, error)
type ConfirmTerminationFunc func(token string) error

func orderDelete(d *schema.ResourceData, meta interface{}, terminate TerminateFunc, confirm ConfirmTerminationFunc) error {
	oldEmailsIds, err := notificationEmailSortedIds(meta)
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
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		email, err = getNewNotificationEmail(matches, oldEmailsIds, meta)
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

	if d != nil {
		d.SetId("")
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

func waitForOrder(c *ovh.Client, id int64) resource.StateRefreshFunc {
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
