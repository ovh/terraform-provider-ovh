package ovh

import ()

type OrderCartGenericProduct struct {
	PlanCode    string                         `json:"planCode"`
	Prices      []OrderCartGenericProductPrice `json:"prices"`
	ProductName string                         `json:"productName"`
	ProductType string                         `json:"productType"`
}

func (v OrderCartGenericProduct) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["plan_code"] = v.PlanCode
	obj["product_name"] = v.ProductName
	obj["product_type"] = v.ProductType

	if v.Prices != nil {
		var prices []map[string]interface{}
		for _, price := range v.Prices {
			prices = append(prices, price.ToMap())
		}
		obj["prices"] = prices
	}
	return obj
}

type OrderCartGenericProductPrice struct {
	Capacities      []string                          `json:"capacities"`
	Description     string                            `json:"description"`
	Duration        string                            `json:"duration"`
	Interval        int                               `json:"interval"`
	MaximumQuantity int                               `json:"maximumQuantity"`
	MaximumRepeat   int                               `json:"maximumRepeat"`
	MinimumQuantity int                               `json:"minimumQuantity"`
	MinimumRepeat   int                               `json:"minimumRepeat"`
	Price           OrderCartGenericProductPricePrice `json:"price"`
	PriceInUcents   int64                             `json:"priceInUcents"`
	PricingMode     string                            `json:"pricingMode"`
	PricingType     string                            `json:"pricingType"`
}

func (v OrderCartGenericProductPrice) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["capacities"] = v.Capacities
	obj["description"] = v.Description
	obj["duration"] = v.Duration
	obj["interval"] = v.Interval
	obj["maximum_quantity"] = v.MaximumQuantity
	obj["maximum_repeat"] = v.MaximumRepeat
	obj["minimum_quantity"] = v.MinimumQuantity
	obj["minimum_repeat"] = v.MinimumRepeat
	obj["price"] = []interface{}{v.Price.ToMap()}
	obj["price_in_ucents"] = v.PriceInUcents
	obj["pricing_mode"] = v.PricingMode
	obj["pricing_type"] = v.PricingType
	return obj
}

type OrderCartGenericProductPricePrice struct {
	CurrencyCode string  `json:"currencyCode"`
	Text         string  `json:"text"`
	Value        float64 `json:"value"`
}

func (v OrderCartGenericProductPricePrice) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["currency_code"] = v.CurrencyCode
	obj["text"] = v.Text
	obj["value"] = v.Value
	return obj
}

type OrderCartGenericOptions struct {
	OrderCartGenericProduct

	Exclusive bool   `json:"exclusive"`
	Family    string `json:"family"`
	Mandatory bool   `json:"mandatory"`
}

func (v OrderCartGenericOptions) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["exclusive"] = v.Exclusive
	obj["family"] = v.Family
	obj["mandatory"] = v.Mandatory
	obj["plan_code"] = v.PlanCode
	obj["product_name"] = v.ProductName
	obj["product_type"] = v.ProductType

	if v.Prices != nil {
		var prices []map[string]interface{}
		for _, price := range v.Prices {
			prices = append(prices, price.ToMap())
		}
		obj["prices"] = prices
	}
	return obj
}
