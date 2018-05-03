package ovh

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
)

type CreditCard struct {
	Description      string `json:"description"`
	Number           string `json:"number"`
	Expiration       string `json:"expirationDate"`
	Default          bool   `json:"defaultPaymentMean"`
	State            string `json:"state"`
	ThreeDSValidated bool   `json:"threeDsValidated"`
	Id               int    `json:"id"`
	Type             string `json:"type"`
}

func dataSourceCreditCard() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCreditCardRead,
		Schema: map[string]*schema.Schema{
			"description_regexp": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  ".*",
			},
			"use_default": &schema.Schema{
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"use_last_to_expire": &schema.Schema{
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"states": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			// Computed
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"threeds_validated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCreditCardRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	states_val, states_ok := d.GetOk("states")
	description_regexp := regexp.MustCompile(d.Get("description_regexp").(string))
	use_last_to_expire := d.Get("use_last_to_expire").(bool)
	use_default := d.Get("use_default").(bool)
	var the_credit_card *CreditCard
	var states []interface{}
	if states_ok {
		states = states_val.(*schema.Set).List()
	}
	var credit_card_ids []int
	err := config.OVHClient.Get(
		"/me/paymentMean/creditCard",
		&credit_card_ids,
	)

	if err != nil {
		return fmt.Errorf("Error getting Credit Cards list:\n\t %q", err)
	}
	filtered_credit_cards := []*CreditCard{}
	for _, card_id := range credit_card_ids {
		credit_card := CreditCard{}
		err = config.OVHClient.Get(
			fmt.Sprintf("/me/paymentMean/creditCard/%d", card_id),
			&credit_card,
		)
		if err != nil {
			return fmt.Errorf("Error getting Credit Card %d:\n\t %q", card_id, err)
		}
		if use_default && credit_card.Default == false {
			continue
		}
		if states_ok {
			match := false
			for _, wanted_state := range states {
				if credit_card.State == wanted_state {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		if !description_regexp.MatchString(credit_card.Description) {
			continue
		}
		filtered_credit_cards = append(filtered_credit_cards, &credit_card)
	}
	if len(filtered_credit_cards) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}
	if len(filtered_credit_cards) > 1 {
		if use_last_to_expire {
			sort.Slice(filtered_credit_cards, func(i, j int) bool {
				return (*filtered_credit_cards[i]).Expiration > (*filtered_credit_cards[j]).Expiration
			})
			the_credit_card = filtered_credit_cards[0]
		}
		if use_default {
			match := false
			for _, credit_card := range filtered_credit_cards {
				if (*credit_card).Default {
					match = true
					the_credit_card = credit_card
					break
				}
			}
			if match == false {
				return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
			}
		}
	}
	if len(filtered_credit_cards) == 1 {
		the_credit_card = filtered_credit_cards[0]
	}
	// Set data
	d.Set("description", (*the_credit_card).Description)
	d.Set("number", (*the_credit_card).Number)
	d.Set("expiration_date", (*the_credit_card).Expiration)
	d.Set("type", (*the_credit_card).Type)
	d.Set("state", (*the_credit_card).State)
	d.Set("id", (*the_credit_card).Id)
	d.Set("default", (*the_credit_card).Default)
	d.Set("threeds_validated", (*the_credit_card).ThreeDSValidated)
	d.SetId(fmt.Sprintf("%d", (*the_credit_card).Id))
	return nil
}
