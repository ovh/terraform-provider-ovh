package ovh

import (
	"fmt"

	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

var (
	MePaymentMeanOvhAccountPaymentOpts      = &MeOrderPaymentOpts{PaymentMean: "ovhAccount"}
	MePaymentMeanFidelityAccountPaymentOpts = &MeOrderPaymentOpts{PaymentMean: "fidelityAccount"}
)

func MePaymentMeanBankAccounts(c *ovhwrap.Client) ([]*MePaymentMeanBankAccount, error) {
	ids := &[]int64{}
	endpoint := fmt.Sprintf("/me/paymentMean/bankAccount")
	if err := c.Get(endpoint, ids); err != nil {
		return nil, fmt.Errorf("Error calling GET %s", endpoint)
	}

	results := []*MePaymentMeanBankAccount{}
	for _, id := range *ids {
		paymentMean := &MePaymentMeanBankAccount{}
		endpoint := fmt.Sprintf(
			"/me/paymentMean/bankAccount/%d",
			id,
		)

		if err := c.Get(endpoint, paymentMean); err != nil {
			return nil, fmt.Errorf("Error calling GET %s", endpoint)
		}

		results = append(results, paymentMean)
	}

	return results, nil
}

func MePaymentMeanCreditCards(c *ovhwrap.Client) ([]*MePaymentMeanCreditCard, error) {
	ids := &[]int64{}
	endpoint := fmt.Sprintf("/me/paymentMean/creditCard")
	if err := c.Get(endpoint, ids); err != nil {
		return nil, fmt.Errorf("Error calling GET %s", endpoint)
	}

	results := []*MePaymentMeanCreditCard{}
	for _, id := range *ids {
		paymentMean := &MePaymentMeanCreditCard{}
		endpoint := fmt.Sprintf(
			"/me/paymentMean/creditCard/%d",
			id,
		)

		if err := c.Get(endpoint, paymentMean); err != nil {
			return nil, fmt.Errorf("Error calling GET %s", endpoint)
		}

		results = append(results, paymentMean)
	}

	return results, nil
}

func MePaymentMeanPaypals(c *ovhwrap.Client) ([]*MePaymentMeanPaypal, error) {
	ids := &[]int64{}
	endpoint := fmt.Sprintf("/me/paymentMean/paypal")
	if err := c.Get(endpoint, ids); err != nil {
		return nil, fmt.Errorf("Error calling GET %s", endpoint)
	}

	results := []*MePaymentMeanPaypal{}
	for _, id := range *ids {
		paymentMean := &MePaymentMeanPaypal{}
		endpoint := fmt.Sprintf(
			"/me/paymentMean/paypal/%d",
			id,
		)

		if err := c.Get(endpoint, paymentMean); err != nil {
			return nil, fmt.Errorf("Error calling GET %s", endpoint)
		}

		results = append(results, paymentMean)
	}

	return results, nil
}

func MePaymentMeanDefaultPaymentOpts(c *ovhwrap.Client) (*MeOrderPaymentOpts, error) {
	payment := &MeOrderPaymentOpts{}

	bankAccounts, err := MePaymentMeanBankAccounts(c)
	if err != nil {
		return nil, fmt.Errorf("could not find default payment mean: %v", err)
	}

	for _, ba := range bankAccounts {
		if ba.DefaultPaymentMean {
			payment.PaymentMean = "bankAccount"
			payment.PaymentMeanId = helpers.GetNilInt64Pointer(ba.Id)
			return payment, nil
		}
	}

	creditCards, err := MePaymentMeanCreditCards(c)
	if err != nil {
		return nil, fmt.Errorf("could not find default payment mean: %v", err)
	}

	for _, ba := range creditCards {
		if ba.DefaultPaymentMean {
			payment.PaymentMean = "creditCard"
			payment.PaymentMeanId = helpers.GetNilInt64Pointer(ba.Id)
			return payment, nil
		}
	}

	paypals, err := MePaymentMeanPaypals(c)
	if err != nil {
		return nil, fmt.Errorf("could not find default payment mean: %v", err)
	}

	for _, ba := range paypals {
		if ba.DefaultPaymentMean {
			payment.PaymentMean = "paypal"
			payment.PaymentMeanId = helpers.GetNilInt64Pointer(ba.Id)
			return payment, nil
		}
	}

	return nil, nil
}
