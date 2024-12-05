package ovh

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func kmsServiceKeyIamChecks(resName string) []statecheck.StateCheck {
	urnRe := regexp.MustCompile("urn:v1:eu:resource:okms:.*")
	return []statecheck.StateCheck{
		statecheck.CompareValuePairs(
			resName,
			tfjsonpath.New("name"),
			resName,
			tfjsonpath.New("iam").AtMapKey("display_name"),
			compare.ValuesSame(),
		),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("iam").AtMapKey("urn"),
			knownvalue.StringRegexp(urnRe)),
		statecheck.CompareValuePairs(
			resName,
			tfjsonpath.New("type"),
			resName,
			tfjsonpath.New("iam").AtMapKey("tags").AtMapKey("ovh:type"),
			compare.ValuesSame(),
		),
		statecheck.CompareValuePairs(
			resName,
			tfjsonpath.New("state"),
			resName,
			tfjsonpath.New("iam").AtMapKey("tags").AtMapKey("ovh:state"),
			compare.ValuesSame(),
		),
	}
}

func kmsServiceKeyStateCommonChecks(resName string, keyName string) []statecheck.StateCheck {
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("created_at"),
			knownvalue.NotNull()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("deactivation_reason"),
			knownvalue.Null()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("id"),
			knownvalue.NotNull()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("name"),
			knownvalue.StringExact(keyName)),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("state"),
			knownvalue.StringExact("ACTIVE")),
	}
}

func kmsServiceKeyStateSymmetricChecks(resName string) []statecheck.StateCheck {
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("context"),
			knownvalue.NotNull()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("curve"),
			knownvalue.Null()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("size"),
			knownvalue.Int64Exact(256)),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("type"),
			knownvalue.StringExact("oct")),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("operations"),
			knownvalue.SetExact(
				[]knownvalue.Check{
					knownvalue.StringExact("encrypt"),
					knownvalue.StringExact("decrypt"),
				},
			)),
	}
}

func kmsServiceKeyStateRsaChecks(resName string) []statecheck.StateCheck {
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("context"),
			knownvalue.Null()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("curve"),
			knownvalue.Null()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("size"),
			knownvalue.Int64Exact(2048)),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("type"),
			knownvalue.StringExact("RSA")),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("operations"),
			knownvalue.SetExact(
				[]knownvalue.Check{
					knownvalue.StringExact("sign"),
					knownvalue.StringExact("verify"),
				},
			)),
	}
}

func kmsServiceKeyStateECChecks(resName string) []statecheck.StateCheck {
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("context"),
			knownvalue.Null()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("curve"),
			knownvalue.StringExact("P-256")),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("size"),
			knownvalue.Null()),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("type"),
			knownvalue.StringExact("EC")),
		statecheck.ExpectKnownValue(
			resName,
			tfjsonpath.New("operations"),
			knownvalue.SetExact(
				[]knownvalue.Check{
					knownvalue.StringExact("sign"),
					knownvalue.StringExact("verify"),
				},
			)),
	}
}

const confOkmsServiceKeyTest = `
data "ovh_me" "current_account" {
}

resource "ovh_okms" "kms" {
  ovh_subsidiary = "FR"
  display_name = "%[1]s"
  region = "eu-west-sbg"
}

resource "ovh_okms_service_key" "key_oct" {
  okms_id    = ovh_okms.kms.id
  name       = "%[1]s-sk-oct"
  type       = "oct"
  size       = 256
  operations = ["encrypt", "decrypt"]
  context    = "%[2]s"
}

resource "ovh_okms_service_key" "key_rsa" {
  okms_id    = ovh_okms.kms.id
  name       = "%[1]s-sk-rsa"
  type       = "RSA"
  size       = 2048
  operations = ["sign", "verify"]
}

resource "ovh_okms_service_key" "key_ecdsa" {
  okms_id    = ovh_okms.kms.id
  name       = "%[1]s-sk-ecdsa"
  type       = "EC"
  curve      = "P-256"
  operations = ["sign", "verify"]
}

resource "ovh_okms_service_key_jwk" "key_oct_jwk" {
  okms_id    = ovh_okms.kms.id
  name       = "%[1]s-sk-oct-jwk"
  context    = "%[2]s"
  keys       = [{
      k       = "Wc7IgEZzWicZf1LTJUtA0w"
      key_ops = ["encrypt", "decrypt"]
      kty     = "oct"
  }]
}

resource "ovh_okms_service_key_jwk" "key_rsa_jwk" {
  okms_id    = ovh_okms.kms.id
  name       = "%[1]s-sk-rsa-jwk"
  context    = "%[2]s"
  keys       = [jsondecode(<<EOT
      {
        "key_ops": ["sign", "verify"],
        "kty": "RSA",
        "d": "UgaHfwn09SnrhqNXZC-bD8-1hW-02aAW0IovKrZMpKrx2Fbs1WeHISTfbdQLn6mFEB5fEO7tas6oXANwb2NHnaBgK5Bv15xkI6QqK0ieZ-_--8cvqzO8-3wBnotd3-SRmpbrgpCOx9NN97lc7Iuj2alsVfuFilDF8XptQFS5bbMlLUsEi48Y3al_m9Y8nVr1SJQl5-I8noEcfzup2fG1a7zvdCKUWqsZ3rO2MNjs6xSZ2JeOdgWVQutNM0cTqg1NtDHdhKRIRd0cRATn5r5elrBjtA-7J5Tas-XN9S9UKgQSQw9xGfXg4SSq2ass01mpVNbiUTHiGzm8kexcfG3dYw",
        "dp": "SIrzLAa0L3LqmDh9BbgEzOjosZprbSLJYCcmt9aodgkG-mmiO2HFFpLCHU4YNdDia_plivVrQ4APf9SGhYh4MECxIieSZ5VDLI8cmdXOJnKMwlOUijh3WMccCPDZ-C7RCT7-966GslxuhvLnRNKtegm8Z7YoeWMMY5cD-TUfVMU",
        "dq": "aQv6Kg0Lw13PB0H2VHjBHxm0L6UliHDoE_-GafpGm9A_6i9b_M2wgQT9K1EO7rWeCfZ3t6Q65DG6AaGvztaCEGy_vHHs9xDajFjJxZweD3kbddLb4ZZa5QdauFJNpTH5QcZobjFuq_Ylw1mZxcWjIEDBEk6-CDg2eVYs-AtVCOE",
        "e": "AQAB",
        "n": "qrFKVDudlle41HlYSzK-xE4o55RMLSd6_NmiZ98v_tN80FgG1KxQ0hO_yg3OX0TiVN3juyJJFSfLk1nFw84Qzy7ptOd7ssNrS2Zd48Jh2WZ3F1fDxQ_4Yp_9OXcoSrS8mCgTOmzgywnUdo8ovG1NG1v2rgJUPGotvfUnd3ywCar-UAm0pHdwmzzagmMdGJ7kz4W_4KgS6tsLZRTcFRbG9tA-4eJsqUsrpHNhAIemT5oLuMYIA9wlxd_h1KxZPRje4de8yHKBIqYFMOcpdF9SKqpCc0Xo21MNS8a6g18BYMNSDuyFGaZ8HsKizQfzSAzRV3lOTF-LuJJtWbSWaDdmNQ",
        "p": "7O4PCo_cWzuoD8HO7447T4Dml8G0HD-rHvUXIUO_jlISXT4WkaTqIyinyu4H25PQ6KnaQQBCC-7KIn9npVriO3aJ32ZZ2xmTlMLAWaRl6c9oZXSEJM5-a89gk_Ov_uMoPT0C0L0AFn5qpewGHswArCupxmsEAAuBOFtXStoQnAM",
        "q": "uG5pDYvV-eu3ZL0obLxx2AXWlhCaz_uBL5CyDBM07n2AdKNSDTvk40nV-bu1Ycz0CVui04gVIDzHLtM9c8709v53AqUZMJhtDWjxnUiOgkxycOaN6U0dtIVuDpLcD1wNhcIuqrKebqQO4PEqrj9D9pb_XNSMnKCivXsM75jxi2c",
        "qi": "B5z00bGrZOSNmuQB8KpfkRilLgemroz9Oc-beYSOkwO4RZ4YL9z_c5B66msQWytcF_BFrUpDNYAoWiKNTX3EhvrI5UXFo8UVuO7PzPlO2HZ_1qMrUdeyOybKtVRApdQy02O45zdbwBNzn4T8kviLJHFZcQrUpy97AVOEXYf-8ak"
      }
  EOT
  )]
}

resource "ovh_okms_service_key_jwk" "key_ecdsa_jwk" {
  okms_id    = ovh_okms.kms.id
  name       = "%[1]s-sk-ecdsa-jwk"
  context    = "%[2]s"
  keys       = [jsondecode(<<EOT
      {
        "key_ops": ["sign", "verify"],
        "kty": "EC",
        "crv": "P-256",
        "d": "SIy6AYrv5nGBLQsM7bg7WCbAPxHyUIVTaDyTxrCWPks",
        "x": "V7a79Iv0RdykDIzhJhu5OvkCFJ8rFkFm5r11qwR9QeY",
        "y": "RSUYb-RPSkF5al1D2fnxerahFpHCHtmJRAlUeS1Ehtw"
      }
  EOT
  )]
}
`

const confOkmsServiceKeyDatasourceTest = `
data "ovh_okms_service_key" "key_oct" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_oct.id
}

data "ovh_okms_service_key" "key_rsa" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_rsa.id
}

data "ovh_okms_service_key" "key_ecdsa" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_ecdsa.id
}

data "ovh_okms_service_key_jwk" "key_oct" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_oct.id
}

data "ovh_okms_service_key_jwk" "key_rsa" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_rsa.id
}

data "ovh_okms_service_key_jwk" "key_ecdsa" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_ecdsa.id
}

data "ovh_okms_service_key_pem" "key_oct" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_oct.id
}

data "ovh_okms_service_key_pem" "key_rsa" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_rsa.id
}

data "ovh_okms_service_key_pem" "key_ecdsa" {
  okms_id = ovh_okms.kms.id
  id  = ovh_okms_service_key.key_ecdsa.id
}
`

func getAllChecks(resName string) []statecheck.StateCheck {
	checks := []statecheck.StateCheck{
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("id"),
			"ovh_okms_service_key.key_oct",
			tfjsonpath.New("okms_id"),
			compare.ValuesSame(),
		),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("id"),
			"ovh_okms_service_key.key_rsa",
			tfjsonpath.New("okms_id"),
			compare.ValuesSame(),
		),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("id"),
			"ovh_okms_service_key.key_ecdsa",
			tfjsonpath.New("okms_id"),
			compare.ValuesSame(),
		),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("id"),
			"ovh_okms_service_key_jwk.key_oct_jwk",
			tfjsonpath.New("okms_id"),
			compare.ValuesSame(),
		),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("id"),
			"ovh_okms_service_key_jwk.key_rsa_jwk",
			tfjsonpath.New("okms_id"),
			compare.ValuesSame(),
		),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("id"),
			"ovh_okms_service_key_jwk.key_ecdsa_jwk",
			tfjsonpath.New("okms_id"),
			compare.ValuesSame(),
		),
	}
	checks = append(checks, kmsServiceKeyStateCommonChecks("ovh_okms_service_key.key_oct", resName+"-sk-oct")...)
	checks = append(checks, kmsServiceKeyStateCommonChecks("ovh_okms_service_key.key_rsa", resName+"-sk-rsa")...)
	checks = append(checks, kmsServiceKeyStateCommonChecks("ovh_okms_service_key.key_ecdsa", resName+"-sk-ecdsa")...)
	checks = append(
		checks,
		kmsServiceKeyStateCommonChecks("ovh_okms_service_key_jwk.key_oct_jwk", resName+"-sk-oct-jwk")...)
	checks = append(
		checks,
		kmsServiceKeyStateCommonChecks("ovh_okms_service_key_jwk.key_rsa_jwk", resName+"-sk-rsa-jwk")...)
	checks = append(
		checks,
		kmsServiceKeyStateCommonChecks("ovh_okms_service_key_jwk.key_ecdsa_jwk", resName+"-sk-ecdsa-jwk")...)
	checks = append(checks, kmsServiceKeyStateSymmetricChecks("ovh_okms_service_key.key_oct")...)
	checks = append(checks, kmsServiceKeyStateRsaChecks("ovh_okms_service_key.key_rsa")...)
	checks = append(checks, kmsServiceKeyStateECChecks("ovh_okms_service_key.key_ecdsa")...)
	checks = append(checks, kmsJwkOctChecks("ovh_okms_service_key_jwk.key_oct_jwk")...)
	checks = append(checks, kmsJwkRsaChecks("ovh_okms_service_key_jwk.key_rsa_jwk")...)
	checks = append(checks, kmsJwkEcdsaChecks("ovh_okms_service_key_jwk.key_ecdsa_jwk")...)

	return checks
}

func kmsServiceKeyDatasourceChecks(resName string, datasourceName string) []statecheck.StateCheck {
	checks := []statecheck.StateCheck{}
	for _, key := range []string{
		"created_at",
		"curve",
		"id",
		"name",
		"okms_id",
		"operations",
		"size",
		"state",
		"type",
	} {
		checks = append(checks, statecheck.CompareValuePairs(
			resName,
			tfjsonpath.New(key),
			datasourceName,
			tfjsonpath.New(key),
			compare.ValuesSame()))
	}

	checks = append(checks, kmsServiceKeyIamChecks(datasourceName)...)
	return checks
}

func kmsPemDatasourceChecks(resName string, datasourceName string, hasPem bool) []statecheck.StateCheck {
	// The PEM keys have the same fields are regular keys,
	// plus an additonal keyPEM array for asymetric keys.
	checks := kmsServiceKeyDatasourceChecks(resName, datasourceName)
	if hasPem {
		checks = append(checks, statecheck.ExpectKnownValue(
			datasourceName,
			tfjsonpath.New("keys_pem").AtSliceIndex(0).AtMapKey("pem"),
			knownvalue.NotNull(),
		))
	} else {
		checks = append(checks, statecheck.ExpectKnownValue(
			datasourceName,
			tfjsonpath.New("keys_pem"),
			knownvalue.Null(),
		))
	}

	checks = append(checks, kmsServiceKeyIamChecks(datasourceName)...)
	return checks
}

func kmsJwkDatasourceChecks(resName string, datasourceName string) []statecheck.StateCheck {
	checks := []statecheck.StateCheck{}
	for _, key := range []string{
		"created_at",
		"id",
		"name",
		"okms_id",
		"size",
		"state",
		"type",
	} {
		checks = append(checks, statecheck.CompareValuePairs(
			resName,
			tfjsonpath.New(key),
			datasourceName,
			tfjsonpath.New(key),
			compare.ValuesSame()))
	}

	checks = append(checks, statecheck.CompareValuePairs(
		resName,
		tfjsonpath.New("id"),
		datasourceName,
		tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kid"),
		compare.ValuesSame()))
	checks = append(checks, statecheck.CompareValuePairs(
		resName,
		tfjsonpath.New("type"),
		datasourceName,
		tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kty"),
		compare.ValuesSame()))

	checks = append(checks, kmsServiceKeyIamChecks(datasourceName)...)
	return checks
}

func kmsJwkOctChecks(name string) []statecheck.StateCheck {
	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("key_ops"),
			knownvalue.SetExact(
				[]knownvalue.Check{
					knownvalue.StringExact("encrypt"),
					knownvalue.StringExact("decrypt"),
				})),
		statecheck.CompareValuePairs(
			name,
			tfjsonpath.New("id"),
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kid"),
			compare.ValuesSame()),
		statecheck.ExpectKnownValue(
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kty"),
			knownvalue.StringExact("oct")),
	}

	for _, key := range []string{
		"alg",
		"crv",
		"e",
		"n",
		"use",
		"x",
		"y",
	} {
		checks = append(
			checks,
			statecheck.ExpectKnownValue(
				name,
				tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey(key),
				knownvalue.Null()),
		)
	}

	return checks
}

func kmsJwkRsaChecks(name string) []statecheck.StateCheck {
	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("key_ops"),
			knownvalue.SetExact(
				[]knownvalue.Check{
					knownvalue.StringExact("sign"),
					knownvalue.StringExact("verify"),
				})),
		statecheck.CompareValuePairs(
			name,
			tfjsonpath.New("id"),
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kid"),
			compare.ValuesSame()),
		statecheck.ExpectKnownValue(
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kty"),
			knownvalue.StringExact("RSA")),
	}

	for _, key := range []string{
		"alg",
		"crv",
		"use",
		"x",
		"y",
	} {
		checks = append(
			checks,
			statecheck.ExpectKnownValue(
				name,
				tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey(key),
				knownvalue.Null()),
		)
	}

	for _, key := range []string{
		"e",
		"n",
	} {
		checks = append(
			checks,
			statecheck.ExpectKnownValue(
				name,
				tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey(key),
				knownvalue.NotNull()),
		)
	}

	return checks
}

func kmsJwkEcdsaChecks(name string) []statecheck.StateCheck {
	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("key_ops"),
			knownvalue.SetExact(
				[]knownvalue.Check{
					knownvalue.StringExact("sign"),
					knownvalue.StringExact("verify"),
				})),
		statecheck.CompareValuePairs(
			name,
			tfjsonpath.New("id"),
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kid"),
			compare.ValuesSame()),
		statecheck.ExpectKnownValue(
			name,
			tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey("kty"),
			knownvalue.StringExact("EC")),
	}

	for _, key := range []string{
		"alg",
		"e",
		"n",
		"use",
	} {
		checks = append(
			checks,
			statecheck.ExpectKnownValue(
				name,
				tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey(key),
				knownvalue.Null()),
		)
	}

	for _, key := range []string{
		"crv",
		"x",
		"y",
	} {
		checks = append(
			checks,
			statecheck.ExpectKnownValue(
				name,
				tfjsonpath.New("keys").AtSliceIndex(0).AtMapKey(key),
				knownvalue.NotNull()),
		)
	}

	return checks
}

func getAllDatasourceChecks() []statecheck.StateCheck {
	// Check non-jwk datasources
	checks := kmsServiceKeyDatasourceChecks("ovh_okms_service_key.key_oct", "data.ovh_okms_service_key.key_oct")
	checks = append(
		checks,
		kmsServiceKeyDatasourceChecks("ovh_okms_service_key.key_rsa", "data.ovh_okms_service_key.key_rsa")...)
	checks = append(
		checks,
		kmsServiceKeyDatasourceChecks("ovh_okms_service_key.key_ecdsa", "data.ovh_okms_service_key.key_ecdsa")...)
	// Check jwk datasources
	checks = append(
		checks,
		kmsJwkDatasourceChecks("ovh_okms_service_key.key_oct", "data.ovh_okms_service_key_jwk.key_oct")...)
	checks = append(
		checks,
		kmsJwkDatasourceChecks("ovh_okms_service_key.key_rsa", "data.ovh_okms_service_key_jwk.key_rsa")...)
	checks = append(
		checks,
		kmsJwkDatasourceChecks("ovh_okms_service_key.key_ecdsa", "data.ovh_okms_service_key_jwk.key_ecdsa")...)
	checks = append(checks, kmsJwkOctChecks("data.ovh_okms_service_key_jwk.key_oct")...)
	checks = append(checks, kmsJwkRsaChecks("data.ovh_okms_service_key_jwk.key_rsa")...)
	checks = append(checks, kmsJwkEcdsaChecks("data.ovh_okms_service_key_jwk.key_ecdsa")...)
	// Check PEM datasources
	checks = append(
		checks,
		kmsPemDatasourceChecks("ovh_okms_service_key.key_oct", "data.ovh_okms_service_key_pem.key_oct", false)...)
	checks = append(
		checks,
		kmsPemDatasourceChecks("ovh_okms_service_key.key_rsa", "data.ovh_okms_service_key_pem.key_rsa", true)...)
	checks = append(
		checks,
		kmsPemDatasourceChecks("ovh_okms_service_key.key_ecdsa", "data.ovh_okms_service_key_pem.key_ecdsa", true)...)

	return checks
}

func TestAccOkmsServiceKey(t *testing.T) {
	resName := fmt.Sprintf("test-tf-%d", acctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Test key creation
				Config:            fmt.Sprintf(confOkmsServiceKeyTest, resName, "ctx"),
				ConfigStateChecks: getAllChecks(resName),
			},
			{
				// Test name update
				Config: fmt.Sprintf(confOkmsServiceKeyTest, resName+"2", "ctx"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"ovh_okms_service_key.key_oct",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(
							"ovh_okms_service_key.key_rsa",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(
							"ovh_okms_service_key.key_ecdsa",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: getAllChecks(resName + "2"),
			},
			{
				// Test context update
				Config: fmt.Sprintf(confOkmsServiceKeyTest, resName+"2", "newctx"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"ovh_okms_service_key.key_oct",
							plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction(
							"ovh_okms_service_key.key_rsa",
							plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(
							"ovh_okms_service_key.key_ecdsa",
							plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: getAllChecks(resName + "2"),
			},
			{
				// Test datasource
				Config:            fmt.Sprintf(confOkmsServiceKeyTest+confOkmsServiceKeyDatasourceTest, resName+"2", "newctx"),
				ConfigStateChecks: getAllDatasourceChecks(),
			},
		},
	})
}
