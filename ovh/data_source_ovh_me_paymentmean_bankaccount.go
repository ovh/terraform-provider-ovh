package ovh

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type BankAccount struct {
	Description            string `json:"description"`
	Default                bool   `json:"defaultPaymentMean"`
	State                  string `json:"state"`
	Id                     int    `json:"id"`
	ValidationDocumentLink string `json:"validationDocumentLink"`
	UniqueReference        string `json:"uniqueReference"`
	CreationDate           string `json:"creationDate"`
	MandateSignatureDate   string `json:"mandateSignatureDate"`
	OwnerName              string `json:"ownerName"`
	OwnerAddress           string `json:"ownerAddress"`
	Iban                   string `json:"iban"`
	Bic                    string `json:"bic"`
}

func dataSourceMePaymentmeanBankaccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMePaymentmeanBankaccountRead,
		Schema: map[string]*schema.Schema{
			"description_regexp": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  ".*",
			},
			"use_default": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"use_oldest": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceMePaymentmeanBankaccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	state, state_ok := d.GetOk("state")
	description_regexp := regexp.MustCompile(d.Get("description_regexp").(string))
	use_oldest := d.Get("use_oldest").(bool)
	use_default := d.Get("use_default").(bool)
	var the_bank_account *BankAccount
	var bank_account_ids []int
	endpoint := "/me/paymentMean/bankAccount"
	if state_ok {
		endpoint = fmt.Sprintf("%s?state=%s", endpoint, state)
	}
	err := config.OVHClient.Get(
		endpoint,
		&bank_account_ids,
	)

	if err != nil {
		return fmt.Errorf("Error getting Bank Account list:\n\t %q", err)
	}
	filtered_bank_accounts := []*BankAccount{}
	for _, account_id := range bank_account_ids {
		bank_account := BankAccount{}
		err = config.OVHClient.Get(
			fmt.Sprintf("/me/paymentMean/bankAccount/%d", account_id),
			&bank_account,
		)
		if err != nil {
			return fmt.Errorf("Error getting Bank Account %d:\n\t %q", account_id, err)
		}
		if use_default && bank_account.Default == false {
			continue
		}
		if !description_regexp.MatchString(bank_account.Description) {
			continue
		}
		filtered_bank_accounts = append(filtered_bank_accounts, &bank_account)
	}
	if len(filtered_bank_accounts) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}
	if len(filtered_bank_accounts) > 1 {
		if use_oldest {
			sort.Slice(filtered_bank_accounts, func(i, j int) bool {
				return (*filtered_bank_accounts[i]).CreationDate < (*filtered_bank_accounts[j]).CreationDate
			})
			the_bank_account = filtered_bank_accounts[0]
		}
		if use_default {
			match := false
			for _, bank_account := range filtered_bank_accounts {
				if (*bank_account).Default {
					match = true
					the_bank_account = bank_account
					break
				}
			}
			if match == false {
				return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
			}
		}
	}
	if len(filtered_bank_accounts) == 1 {
		the_bank_account = filtered_bank_accounts[0]
	}
	// Set data
	d.Set("description", (*the_bank_account).Description)
	d.Set("state", (*the_bank_account).State)
	d.Set("default", (*the_bank_account).Default)

	d.SetId(fmt.Sprintf("%d", (*the_bank_account).Id))
	return nil
}
