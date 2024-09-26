package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const crt1 = `
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

const key1 = `
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

const crt2 = `
-----BEGIN CERTIFICATE-----
MIIDnzCCAoegAwIBAgIUfFQedAHR39RYsu38ubdwgQFVYpEwDQYJKoZIhvcNAQEL
BQAwXzELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEYMBYGA1UEAwwPd3d3LmV4YW1wbGUu
Y29tMB4XDTI0MDkyNzA5MDY1MVoXDTI1MDkyNzA5MDY1MVowXzELMAkGA1UEBhMC
QVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdp
dHMgUHR5IEx0ZDEYMBYGA1UEAwwPd3d3LmV4YW1wbGUuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwWKsvwGKSMLYBU0MDPN+I7lUcFZowMh7XL2Y
fYKypDt23bzBcBzY0vgJJFYA74WG28Djv2VRTTOpho0eYJ1vt2ap6A3uID8xyQc8
Bocwaaws0b24cn0KShagW71XW2ZXI1gZIwRvaJK18hxqkdS5+Nv/CfE46+wrW9GG
9RAF1E/C3+mJq79hjH0bdHDwic91x/cEw8K+CQlZXbtQyPdl9SV398ccfcQ602iv
VcXIOs5OD0YAXkyoeJKpe182Z6BrsqNJanEa0kMXSgyCcGHIaB5HUpCrRW9v4Xgy
vQwoim1iyOonpqUWFtXKM06OiwypJaj6HJAZTvYvbkO9ibURFQIDAQABo1MwUTAd
BgNVHQ4EFgQUvQQRWOJazR1pWoIQZnPPe7o7d+IwHwYDVR0jBBgwFoAUvQQRWOJa
zR1pWoIQZnPPe7o7d+IwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AQEAJ9Uf31Ttcfw7GfsZunYRwt3hMROqp5PFJ/RfYfQ4jUfhuvRG/Dl3OUuiDhCc
ny+jtcitJvghNBwtWRZCLIozsjYz/Ve6khgzU/s4OeR0dxtQAOMynxp6j6ENFjmW
35lelGHT70ClCvksRpZsIIaVIxahPi3T8FAUuyOZhfLSzAQXvLXH8xz1KmllzDew
H9lNBcmUPM7NwFCsBbH5JkbghCiD7cWPltMgZG3fCvDjAvhTkbpieeElrUD1t0CZ
vnmJN5F9yb1QPeeOCsRYVffnAvFuh88x9Xbe/h+G8S1KW+AR8PM64U/V7GiXKcsW
jeTEd26AdLKSb9blbQRooXK8+A==
-----END CERTIFICATE-----
`

const key2 = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDBYqy/AYpIwtgF
TQwM834juVRwVmjAyHtcvZh9grKkO3bdvMFwHNjS+AkkVgDvhYbbwOO/ZVFNM6mG
jR5gnW+3ZqnoDe4gPzHJBzwGhzBprCzRvbhyfQpKFqBbvVdbZlcjWBkjBG9okrXy
HGqR1Ln42/8J8Tjr7Ctb0Yb1EAXUT8Lf6Ymrv2GMfRt0cPCJz3XH9wTDwr4JCVld
u1DI92X1JXf3xxx9xDrTaK9Vxcg6zk4PRgBeTKh4kql7XzZnoGuyo0lqcRrSQxdK
DIJwYchoHkdSkKtFb2/heDK9DCiKbWLI6iempRYW1cozTo6LDKklqPockBlO9i9u
Q72JtREVAgMBAAECggEAUaLoLdHmoISwULyPw7/elhIclTfFDO0VNrdK8mEELNR1
f3G5byyeU6ElqtMrnfGOVqJ9AkUxJxgSDpzwH/UWPiP2weFvjuleiPCm5FKZm2J1
oS0n+hzTBSe8Fj497HWtf1wncGLk7Q5SBJz+WAWnZGjqpzXPw7h1LVOsVzOmYlN1
OIYiN9TJx69GjQZ7q8mPhxmbQGznX3yMKHzpeVJrKs6GaOkFJPsJpNXaqnw/ukZB
IOmkX3Lh89RHZtWOl2J2TZO75d5Eo0UptX+xanDSWSiHGFiFltl//EPIF61gaOUb
6zv5MdPLGBz2DEIKPt9DHL+Y6pO42d5yDyZVIaSCmQKBgQDt0tLCzLyPQsj9pnuq
+XRmT5X5uosKZBm9j6JcH6AN/uA0OffbMnkU2mM8W18bboU5N02G3UlIegyciFIu
7Q56CqtizU8WEkbVNMT9x0njE08o6b13VfJpe317kiJkxDpaarDG35rlBLn5NbMX
Lz2+Pvi93Ghq5MlnFBtF6KkuOwKBgQDQKmQEF0LptnTSdj8RsBB8VnI7fFOOS2vA
2G7FQ6YtOB6TKG4+Xr3+Y3mjKnQzhS8MYBjZW30r7c9mGnE6gsu0nxW89ZxcMoFY
/O4t07FrDyk+Mgylxk0o11CsvCgmNSZsGKQLnNFB22JvgXvMCSm14UcoodOE6JSk
P/MEnJI47wKBgQDjms3scsbPEKqM3sUiV0UIGYi+DMK+uhhMa+eF5Rpn6pKmSQgl
twNcarlobUXzWe2UWZIPzT4yZ+Qj9u84S9V8BTQLBdR3mhxCNhOFkTHsaXFsBW83
g4DRZMct+SiDaa0zFUKDwkJahhB6KeNw/9xGc7vY8NTZ4IXD6eFEIauwOwKBgQC7
8Snvz5ioLgV3Zy82JwIbYGkub+d4l3un10DbdWJ5fSuNrXkrcpqTLPjLai8TvPFn
ePO2erWejc0s4D7OlvyYDPGxcwdC7o59646XbYPHpx6TetiUa1+KuiuqaQ7OwDn6
apzhpyn/DbXn+r3sz3ELN2PpuYGhg+OAEAS3ay1RnQKBgGpx5KLTpbbv4OKSStkA
JgoHXUxZ8ZkDKHd/HikI8q5inplE1v6hjxzdbOLfNCZmsmLN3Lr7zY3V2+KrVNJw
W81qUpgbMTA6TozkgsaS6EuXqPG/iaS6iAu3KBqthzDXeTJqH6kkDFuYuzQK6Egg
AvulFF6dizi0pVFg/H84rN2n
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
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "aaa", crt1, key1, crt1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "display_name", "aaa"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "certificate", fmt.Sprintf("%s\n", crt1)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "key", fmt.Sprintf("%s\n", key1)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "chain", fmt.Sprintf("%s\n", crt1)),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "bbb", crt1, key1, crt1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "display_name", "bbb"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "certificate", fmt.Sprintf("%s\n", crt1)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "key", fmt.Sprintf("%s\n", key1)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "chain", fmt.Sprintf("%s\n", crt1)),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "ccc", crt2, key2, crt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "display_name", "ccc"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "certificate", fmt.Sprintf("%s\n", crt2)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "key", fmt.Sprintf("%s\n", key2)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "chain", fmt.Sprintf("%s\n", crt2)),
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
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "aaa", crt1, key1, crt1),
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
