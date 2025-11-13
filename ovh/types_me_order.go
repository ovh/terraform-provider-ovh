package ovh

type MeOrder struct {
	Date           string `json:"date"`
	ExpirationDate string `json:"expirationDate"`
	OrderId        int64  `json:"orderId"`
}

func (v MeOrder) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["date"] = v.Date
	obj["expiration_date"] = v.ExpirationDate
	obj["order_id"] = v.OrderId
	return obj
}

type MeOrderDetail struct {
	Description   string                  `json:"description"`
	Domain        string                  `json:"domain"`
	OrderDetailId int64                   `json:"orderDetailId"`
	Quantity      string                  `json:"quantity"`
	Extension     *MeOrderDetailExtension `json:"-"`
}

func (v MeOrderDetail) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["description"] = v.Description
	obj["domain"] = v.Domain
	obj["order_detail_id"] = v.OrderDetailId
	obj["quantity"] = v.Quantity
	return obj
}

type MeOrderDetailExtension struct {
	Order struct {
		Plan struct {
			Code string `json:"code"`
		} `json:"plan"`
		Configurations []MeOrderDetailExtensionConfiguration `json:"configurations"`
	} `json:"order"`
}

type MeOrderDetailExtensionConfiguration struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (v MeOrderDetailExtensionConfiguration) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"value": v.Value,
		"label": v.Label,
	}
}

type MeOrderDetailOperation struct {
	Status   string                         `json:"status"`
	ID       int                            `json:"id"`
	Type     string                         `json:"type"`
	Resource MeOrderDetailOperationResource `json:"resource"`
	Quantity int                            `json:"quantity"`
}

type MeOrderDetailOperationResource struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	DisplayName string `json:"displayName"`
}

type MeOrderPaymentOpts struct {
	PaymentMean   string `json:"paymentMean"`
	PaymentMeanId *int64 `json:"paymentMeanId,omitempty"`
}
type MeOrderPaymentMethodOpts struct {
	PaymentMethod PaymentMethod `json:"paymentMethod"`
}
type PaymentMethod struct {
	Id int64 `json:"id"`
}
