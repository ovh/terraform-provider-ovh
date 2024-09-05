package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const certificate = `
-----BEGIN CERTIFICATE-----
MIIDnzCCAoegAwIBAgIUchdtmNBNsdO0rJFBZEr14/5zAe4wDQYJKoZIhvcNAQEL
BQAwXzELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEYMBYGA1UEAwwPd3d3LmV4YW1wbGUu
Y29tMB4XDTI0MDkwNTEzNTkxNVoXDTI1MDkwNTEzNTkxNVowXzELMAkGA1UEBhMC
QVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdp
dHMgUHR5IEx0ZDEYMBYGA1UEAwwPd3d3LmV4YW1wbGUuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2MNhJknJe62vUehYuzcRFUC3Xh3L4t5H2ZqD
xWloB6SUpzIhAmXtjmAvdwTfdB8Ne5+paOjki1/9HfV4OU3s5RAU21yGENEaZ7x/
DJJnTyqwrFJY+iV+szvxRww5qt3SKVkRsDDwtk2yG7cEw4d2mlYoJAEAfL3EQeTt
/m/y+39Kq1ScnjcCslm9BdQUXiq8gqx+z6GyaZakVO8Kks0Sb45mkcN0KA71m9is
giLWVIrMy2Sb0pUprVZm8kd87pUBVM3bhbdv0iQTZRKirLJKLKnnq5wGFnb0yRer
8L7f91sMnZmohkWGJZ+oirzbjmX/H3h3wRkn1ns0I5PTkpQcUQIDAQABo1MwUTAd
BgNVHQ4EFgQUKrQqUfXiw0OaOD4d5GtNmKEe+CYwHwYDVR0jBBgwFoAUKrQqUfXi
w0OaOD4d5GtNmKEe+CYwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AQEAVH1KHduXXKVeFXKLUnSpFMyjw/yw7WKiDH2myen8dA6UDUMj4W01PyoWx3Ic
U1D37zNgIXFr8+ViEi5Bpfz1FAGsRCNfP42Wi52xygkm7tezMT+VQN7WLRqPs4n9
ogwNAm9HrYFuEBAeQZlCL1VIJhvk02JHELyZh+kUz77JjO/RcRS/qNIsHbf3TxT0
6mNVQsoYWeBgKR894kGbCjsHFEKJT2Hf5IpsRy8fD70qSx/dE3paXHFzAXDybVxu
e5fge2/fk/rSbUm5CUsOwoxjNx100eRbH0BQTpdtgV1SzF3G127XTldZVkcdai6G
TOP9quQjYN/Q8Q+sMud9sDFeKA==
-----END CERTIFICATE-----
`

const key = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDYw2EmScl7ra9R
6Fi7NxEVQLdeHcvi3kfZmoPFaWgHpJSnMiECZe2OYC93BN90Hw17n6lo6OSLX/0d
9Xg5TezlEBTbXIYQ0RpnvH8MkmdPKrCsUlj6JX6zO/FHDDmq3dIpWRGwMPC2TbIb
twTDh3aaVigkAQB8vcRB5O3+b/L7f0qrVJyeNwKyWb0F1BReKryCrH7PobJplqRU
7wqSzRJvjmaRw3QoDvWb2KyCItZUiszLZJvSlSmtVmbyR3zulQFUzduFt2/SJBNl
EqKsskosqeernAYWdvTJF6vwvt/3WwydmaiGRYYln6iKvNuOZf8feHfBGSfWezQj
k9OSlBxRAgMBAAECggEAWlUl1c53NF7/ypsQ60g6CsjTAdpZ/twSRklhs7HHJDQ+
rOSzo+u1UZmc/jUeKCbOuB+j+m/f2oNwmP0UkpD6ccU/Y+FNj5GMtwFzUtpqSjAo
u09//BMHF4uZ87lRCPdzHz8ao3npvpdna6xcRF3eG9he1w5B1TpCIRHV6qxdrter
IS99GIJvA96OUgDe1jBSjJl9IsnaE4dMNyBQkMajo7NFgwPRIJk/XBbUUe7t+iaM
L/9dGKIZ+WIFsXqXc++wdTCo8ECEMT5/Yy/sFpUOvraeBLYm/+F1PQEoTQZ0c+p6
ZI+owxIDCGv/wJivt2qX9+HHJWHKaKfFuFunEStoDwKBgQD6e9+hjrVd9L8kilfx
Jso1aC51SrCE6tuLx5Q4CnklEyYidFBhEcxMUUgQiK2EQP2+GA/e+ma/IZbO+b2E
2nYkoxKHn8afa6JbXx4HO5q6fZDgcwfI4pX64R1IBpFikUySJM0hUBWVIY+2UF0f
NOLGCBiUU/i4yx5HNxufKPILxwKBgQDdiWcFawqgExaNNdpvU8lCL81hW+V3jNzC
CqhFw0n20MObd5qZwSk3voCAhav30KYHjr8Smecp6rbmLpei0rd8CGa30tL9FdtR
P152YlxHYoYHU2KdWD46/rgClvMMGYAc+1dJv58OiHvehQmzGOKbci0u21MxPxe/
KDsAkM8nJwKBgQCUE5ztribL538EBADfH/ZUQkWMs13NBeZKKO8XfiGF6F8X6TkH
WXUz/K0kkRg64gzfTuw6/j61aQ71RrBiFJ/ZIso2gR7zabbuWzmuPu9Gpip6daY5
fLH7QQ+FX9Sct5bToovd0LEhm1iRB8s1Qpd5SJn3PfkAjZtVsF9U5OjKSwKBgQCm
zyYeY0od3CGX9Fvkhb8+MgZAb9Spnww+o42u8exIh0syTe3AJjzl93CE1aH2OEo7
2JUw6WexHUXYrm6JMIbuQtktQvaRkJqSY9e55jg7nAj1jSjs9xvsig1+DbE2hCD+
MZa5NisK42P52kzCaVN/3on9BTJwG2TDEATVWTRR8wKBgBJ5Wv/pjAQK3JconXw9
AaDcfFXW9/LkdRdHBNDcj8hFPuSlQiMyLNztYuzUH3DPK9HpACtKEfoOGA53QeqZ
tLT0VQ+kvGVHp3ff1oFXslRw4USnjD9wTfhCSDtPZUtZQKD33575FISRV8kVToba
s9niPsoEYo3+0dm/OhJymKKD
-----END PRIVATE KEY-----
`

const testAccCheckOvhIpLoadbalancingSslConfig = `
resource "ovh_iploadbalancing_ssl" "testssl" {
	service_name = "%s"
	display_name = "%s"
	certificate  = <<EOT
%s
EOT
	key          = <<EOT
%s
EOT
	chain        = <<EOT
%s
EOT
}
`

func TestAccIpLoadbalancingSsl_basic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "aaa", certificate, key, certificate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "display_name", "aaa"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "certificate", fmt.Sprintf("%s\n", certificate)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "key", fmt.Sprintf("%s\n", key)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "chain", fmt.Sprintf("%s\n", certificate)),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "bbb", certificate, key, certificate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "display_name", "bbb"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "certificate", fmt.Sprintf("%s\n", certificate)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "key", fmt.Sprintf("%s\n", key)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "chain", fmt.Sprintf("%s\n", certificate)),
				),
			},
		},
	})
}

func TestAccIpLoadbalancingSsl_importBasic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "aaa", certificate, key, certificate),
			},
			{
				ResourceName:            "ovh_iploadbalancing_ssl.testssl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "key", "chain"},
				ImportStateIdFunc:       testAccIpLoadbalancingSsl_import("ovh_iploadbalancing_ssl.testssl"),
			},
		},
	})
}

func testAccIpLoadbalancingSsl_import(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testIpLoadbalancingSsl, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_iploadbalancing_ssl not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testIpLoadbalancingSsl.Primary.Attributes["service_name"],
			testIpLoadbalancingSsl.Primary.Attributes["id"],
		), nil
	}
}
