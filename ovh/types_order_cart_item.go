package ovh

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type OrderCartItem struct {
	CartId         string                      `json:"cartId"`
	Configurations []int64                     `json:"configurations"`
	Duration       string                      `json:"duration"`
	ItemId         int64                       `json:"itemId"`
	OfferId        string                      `json:"offerId"`
	Options        []int64                     `json:"options"`
	ParentItemId   int64                       `json:"parentItemId"`
	ProductId      string                      `json:"productId"`
	Settings       OrderCartItemDomainSettings `json:"settings"`
}

type OrderCartItemDomainSettings struct {
	Domain string `json:"domain"`
}

type OrderCartItemConfiguration struct {
	Id    int64  `json:id`
	Label string `json:"label"`
	Value string `json:"value"`
}

type OrderCartItemConfigurationOpts struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (opts *OrderCartItemConfigurationOpts) FromResourceWithPath(d *schema.ResourceData, path string) *OrderCartItemConfigurationOpts {
	opts.Label = d.Get(fmt.Sprintf("%s.label", path)).(string)
	opts.Value = d.Get(fmt.Sprintf("%s.value", path)).(string)
	return opts
}
