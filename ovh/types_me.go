package ovh

import (
	"fmt"
)

type MeResponse struct {
	Address                             *string     `json:"address"`
	Area                                *string     `json:"area"`
	BirthCity                           *string     `json:"birthCity"`
	BirthDay                            *string     `json:"birthDay"`
	City                                *string     `json:"city"`
	CompanyNationalIdentificationNumber *string     `json:"companyNationalIdentificationNumber"`
	CorporationType                     *string     `json:"corporationType"`
	Country                             string      `json:"country"`
	Currency                            *MeCurrency `json:"currency"`
	CustomerCode                        *string     `json:"customerCode"`
	Email                               string      `json:"email"`
	Fax                                 *string     `json:"fax"`
	Firstname                           *string     `json:"firstname"`
	ItalianSDI                          *string     `json:"italianSDI"`
	Language                            *string     `json:"language"`
	Legalform                           string      `json:"legalform"`
	Name                                *string     `json:"name"`
	NationalIdentificationNumber        *string     `json:"nationalIdentificationNumber"`
	Nichandle                           string      `json:"nichandle"`
	Organisation                        *string     `json:"organisation"`
	OvhCompany                          string      `json:"ovhCompany"`
	OvhSubsidiary                       string      `json:"ovhSubsidiary"`
	Phone                               *string     `json:"phone"`
	PhoneCountry                        *string     `json:"phoneCountry"`
	Sex                                 *string     `json:"sex"`
	SpareEmail                          *string     `json:"spareEmail"`
	State                               string      `json:"state"`
	Vat                                 *string     `json:"vat"`
	Zip                                 *string     `json:"zip"`
}

func (m MeResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	// Non-nullable values
	obj["country"] = m.Country
	obj["email"] = m.Email
	obj["legalform"] = m.Legalform
	obj["nichandle"] = m.Nichandle
	obj["ovh_company"] = m.OvhCompany
	obj["ovh_subsidiary"] = m.OvhSubsidiary
	obj["state"] = m.State

	if m.Currency != nil {
		obj["currency"] = []interface{}{m.Currency.ToMap()}
	}

	// Nullable values
	if m.Address != nil {
		obj["address"] = *m.Address
	}
	if m.Area != nil {
		obj["area"] = *m.Area
	}
	if m.BirthCity != nil {
		obj["birth_city"] = *m.BirthCity
	}
	if m.BirthDay != nil {
		obj["birth_day"] = *m.BirthDay
	}
	if m.City != nil {
		obj["city"] = *m.City
	}
	if m.CompanyNationalIdentificationNumber != nil {
		obj["company_national_identification_number"] = *m.CompanyNationalIdentificationNumber
	}
	if m.CorporationType != nil {
		obj["corporation_type"] = *m.CorporationType
	}
	if m.CustomerCode != nil {
		obj["customer_code"] = *m.CustomerCode
	}
	if m.Fax != nil {
		obj["fax"] = *m.Fax
	}
	if m.Firstname != nil {
		obj["firstname"] = *m.Firstname
	}
	if m.ItalianSDI != nil {
		obj["italian_sdi"] = *m.ItalianSDI
	}
	if m.Language != nil {
		obj["language"] = *m.Language
	}
	if m.Name != nil {
		obj["name"] = *m.Name
	}
	if m.NationalIdentificationNumber != nil {
		obj["national_identification_number"] = *m.NationalIdentificationNumber
	}
	if m.Organisation != nil {
		obj["organisation"] = *m.Organisation
	}
	if m.Phone != nil {
		obj["phone"] = *m.Phone
	}
	if m.PhoneCountry != nil {
		obj["phone_country"] = *m.PhoneCountry
	}
	if m.Sex != nil {
		obj["sex"] = *m.Sex
	}
	if m.SpareEmail != nil {
		obj["spare_email"] = *m.SpareEmail
	}
	if m.Vat != nil {
		obj["vat"] = *m.Vat
	}
	if m.Zip != nil {
		obj["zip"] = *m.Zip
	}

	return obj
}

type MeCurrency struct {
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

func (c MeCurrency) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["code"] = c.Code
	obj["symbol"] = c.Symbol
	return obj
}

type MeIdentityGroupResponse struct {
	Name         string `json:"name"`
	DefaultGroup bool   `json:"defaultGroup"`
	Role         string `json:"role"`
	Creation     string `json:"creation"`
	Description  string `json:"description"`
	LastUpdate   string `json:"lastUpdate"`
}

type MeIdentityGroupCreateOpts struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

type MeIdentityGroupUpdateOpts struct {
	Description string `json:"description"`
	Role        string `json:"role"`
}

type MeIdentityUserResponse struct {
	Creation           string `json:"creation"`
	Description        string `json:"description"`
	Email              string `json:"email"`
	Group              string `json:"group"`
	LastUpdate         string `json:"lastUpdate"`
	Login              string `json:"login"`
	PasswordLastUpdate string `json:"passwordLastUpdate"`
	Status             string `json:"status"`
}

// MeIdentityUser Opts
type MeIdentityUserCreateOpts struct {
	Description string `json:"description"`
	Email       string `json:"email"`
	Group       string `json:"group"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

type MeIdentityUserUpdateOpts struct {
	Description string `json:"description"`
	Email       string `json:"email"`
	Group       string `json:"group"`
}

// MeSshKey Opts
type MeSshKeyCreateOpts struct {
	KeyName string `json:"keyName"`
	Key     string `json:"key"`
}

type MeSshKeyResponse struct {
	KeyName string `json:"keyName"`
	Key     string `json:"key"`
	Default bool   `json:"default"`
}

func (s *MeSshKeyResponse) String() string {
	return fmt.Sprintf("SSH Key: %s, key:%s, default:%t",
		s.Key, s.KeyName, s.Default)
}

type MeSshKeyUpdateOpts struct {
	Default bool `json:"default"`
}

// MeIpxeScript Opts
type MeIpxeScriptCreateOpts struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Script      string `json:"script"`
}

type MeIpxeScriptResponse struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

func (s *MeIpxeScriptResponse) String() string {
	return fmt.Sprintf("IpxeScript: %s", s.Name)
}

type MePaymentMeanBankAccount struct {
	Bic                    string             `json:"bic"`
	CreationDate           string             `json:"creationDate"`
	DefaultPaymentMean     bool               `json:"defaultPaymentMean"`
	Description            *string            `json:"description"`
	Iban                   string             `json:"iban"`
	Icon                   *MePaymentMeanIcon `json:"icon"`
	Id                     int64              `json:"id"`
	MandateSignatureDate   *string            `json:"mandateSignatureDate"`
	OwnerAddress           string             `json:"ownerAddress"`
	OwnerName              string             `json:"ownerName"`
	State                  string             `json:"state"`
	UniqueReference        string             `json:"uniqueReference"`
	ValidationDocumentLink *string            `json:"validationDocumentLink"`
}

type MePaymentMeanCreditCard struct {
	DefaultPaymentMean bool               `json:"defaultPaymentMean"`
	Description        *string            `json:"description"`
	ExpirationDate     string             `json:"expirationDate"`
	Icon               *MePaymentMeanIcon `json:"icon"`
	Id                 int64              `json:"id"`
	Number             string             `json:"number"`
	State              string             `json:"state"`
	ThreeDsValidated   string             `json:"threeDsValidated"`
	Type               string             `json:"type"`
}

type MePaymentMeanPaypal struct {
	AgreementId        string             `json:"agreementId"`
	CreationDate       string             `json:"creationDate"`
	DefaultPaymentMean bool               `json:"defaultPaymentMean"`
	Description        *string            `json:"description"`
	Email              string             `json:"email"`
	Icon               *MePaymentMeanIcon `json:"icon"`
	Id                 int64              `json:"id"`
	State              string             `json:"state"`
}

type MePaymentMeanIcon struct {
	Data *string `json:"data"`
	Name *string `json:"name"`
}
