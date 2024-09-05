package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const certificate = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArtzN8YipTfmxwLtfj0K5
3bK22jOxuINkkBsOcjRzIvk5HeVoacq12ri+qhX/+82AA70Pq5wjKSEld9MDtvrf
tzSzysdxGqAO2ZzrP/vBtljA+/fOQjZn6H5DpAjfz64TKd9aISgtXVGQm98PmTgD
FqIXtY5jMFEIqM/tY7qVjX/CHxQGYYXvKXp2FMfwxQXUe4oK431GOChvIOm0LYO6
64hPphS1rZSboNkxgyni3fQSgrktQhYSKSdnp5693NX9yBysVjtSIOuVytBR2gsz
8o5uqCb8s8LO0YSwLxhnRwfO7sFVBQLQC5cwMJSAsOPCyaAha7MxCrYSbDJGbZXY
swIDAQAB
-----END PUBLIC KEY-----
`

const key = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCu3M3xiKlN+bHA
u1+PQrndsrbaM7G4g2SQGw5yNHMi+Tkd5WhpyrXauL6qFf/7zYADvQ+rnCMpISV3
0wO2+t+3NLPKx3EaoA7ZnOs/+8G2WMD7985CNmfofkOkCN/PrhMp31ohKC1dUZCb
3w+ZOAMWohe1jmMwUQioz+1jupWNf8IfFAZhhe8penYUx/DFBdR7igrjfUY4KG8g
6bQtg7rriE+mFLWtlJug2TGDKeLd9BKCuS1CFhIpJ2ennr3c1f3IHKxWO1Ig65XK
0FHaCzPyjm6oJvyzws7RhLAvGGdHB87uwVUFAtALlzAwlICw48LJoCFrszEKthJs
MkZtldizAgMBAAECggEAOV3WMJnZlWvH/YdbQ9gruwhhBbf046gzizVdKWl8pmol
62FyPlbTD3URlCJQj07tYwyZAf9g56LotGVlcBfg4i1nmKySth3xkUHySdTxyb1X
XrJ/F+jacQNPdJa2ul8NCW0tl/gi3d0e+IwXBXjDRp8Z8SXL87N6CEviwRea7ccC
9RIl937DVSxXZyLYvAF6CgVCljBGD4g+vuRBwukbtDtbU/psouapffbIi0rO2HeP
CpRItQxH+VMqVbIl0oXu1flMinEHz/kLFsftrBSHP6gZDsXG6x4JFgBCsap2dOOc
oe+bxgrArzRuXnp/iFFWlYZLH9AO64+pfsmlpEP7/QKBgQDs7Z411UO+wFqO2bWT
N6TDTQQELmoiui0UJOsN6/oddmrqdUVMLFmU5+PnydPvRZplLokw9IYP0WEF6NU0
2pD3qERImN4d/HScntCuom2JQ0rxQixnDpoukxm5S3x++0ZNoz0zFHqhKwa0Rlot
cONTqKOLs5o8yYCjvhJiP2v89QKBgQC88DLGBUmUSm7bDOscT66dazt0CQZM3T/q
tWBFSS1g4lD/igexmJN4cqVUFPcsLa3nq4f3ehBgAXMfUEBpVcMnPPAMTMwbCYD0
0/jega3bHm6vToxlzidtSeXQvzJLJCKxeAtZOHjQmDGwHvbU+ccLg1lv/dIo+eYF
GqP6qUV2BwKBgQCr5QnwDLaF4pDRK2rtUGWdvHa5geNHJsQl5VMUWqywS4XubP7F
8TddLZDQTkIRSvJljonCluXc/A/kdbSaECk1RUOlWCNupgcEysSkrvvBpqzstRH6
A0IhoF/9a6L7jdrH4TM5qBYAcHMAwDoU1d5Yh/WAGeJBUACgy/oShekXMQKBgFWL
h+GiuTrcLK8ffUA3T79UsvmJsIGS14LElo8oX9RzA+t/qpYdl/+8IOEeEP++uvOe
9ZP2f32IioBODKvkudSFQca/6tX/CpVPeGn+WyJP+BuFvAnIOo/AGr7WIsZk2RRz
XugJqqH/ltfAXU/2u8mZsiAD02jcJOqAsmgmxh13AoGBAN5G6MLQ5VNKuB30kpao
7EK1p3fIhoEUY/WPaB1+0eBK+0XxZKNGi+4ANKH7HTuGn5mulx5MqLQv04xOHyYP
Jj5JoQmLfOgmSXN184jhNauGiu4KMzCFlKjoMBhDAhs43WiXHmSGLxq1HficWsMr
cj+wIDlKII65E3OW8V0sfFK4
-----END PRIVATE KEY-----
`

const chain = ``

const testAccCheckOvhIpLoadbalancingSslConfig = `
resource "ovh_iploadbalancing_ssl" "testssl" {
	service_name = "%s"
	display_name = "%s"
	certificate  = "%s"
	key          = "%s"
	chain		 = "%s
}
`

func TestAccIpLoadbalancingSsl_basic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "aaa", certificate, key, chain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "display_name", "aaa"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "certificate", certificate),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "key", key),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "chain", chain),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingSslConfig, iplb, "bbb", certificate, key, chain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "display_name", "bbb"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "certificate", certificate),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "key", key),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_ssl.testssl", "chain", chain),
				),
			},
		},
	})
}
