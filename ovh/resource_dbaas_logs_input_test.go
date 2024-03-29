package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccResourceDbaasLogsInput_basic = `
data "ovh_dbaas_logs_input_engine" "logstash" {	
	service_name  = "%s"
	name          = "%s"
	version       = "%s"
}

resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
 service_name = "%s"
 title        = "%s"
 description  = "%s"
}

resource "ovh_dbaas_logs_input" "input" {
 service_name = ovh_dbaas_logs_output_graylog_stream.stream.service_name
 description  = ovh_dbaas_logs_output_graylog_stream.stream.description
 title        = ovh_dbaas_logs_output_graylog_stream.stream.title
 engine_id    = data.ovh_dbaas_logs_input_engine.logstash.id
 stream_id    = ovh_dbaas_logs_output_graylog_stream.stream.id

 allowed_networks = ["10.0.0.0/16"]
 exposed_port     = "6154"
 nb_instance      = 2

 configuration {
   logstash {
       input_section = <<EOF
beats {
  port => 6514
  ssl => true
  ssl_certificate => "/etc/ssl/private/server.crt"
  ssl_key => "/etc/ssl/private/server.key"
}
EOF

   }
 }
}
`

func init() {
	resource.AddTestSweepers("ovh_dbaas_logs_input", &resource.Sweeper{
		Name: "ovh_dbaas_logs_input",
		F:    testSweepDbaasInput,
	})
}

func testSweepDbaasInput(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_DBAAS_LOGS_SERVICE_TEST is not set. No ovh_dbaas_input to sweep")
		return nil
	}

	res := []string{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/input",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	if len(res) == 0 {
		log.Print("[DEBUG] No ovh_dbaas_input to sweep")
		return nil
	}

	for _, id := range res {
		log.Printf("[DEBUG] Will read dbaas logs input id : %s/%s", serviceName, id)

		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/input/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		input := &DbaasLogsInput{}
		if err := config.OVHClient.Get(endpoint, &input); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		if !strings.HasPrefix(input.Title, test_prefix) {
			continue
		}

		res := &DbaasLogsOperation{}
		ctx := context.Background()
		err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
			if input.Status == "PROCESSING" {
				err := fmt.Errorf("[WARN] stop: input %s/%s already has an ongoing action",
					serviceName,
					id,
				)
				return retry.RetryableError(err)
			}

			if input.Status == "RUNNING" {
				log.Printf("[INFO] Will end dbaas logs input for: %s/%s", serviceName, id)
				res := &DbaasLogsOperation{}
				endpoint := fmt.Sprintf(
					"/dbaas/logs/%s/input/%s/end",
					url.PathEscape(serviceName),
					url.PathEscape(id),
				)
				if err := config.OVHClient.Post(endpoint, nil, res); err != nil {
					return retry.RetryableError(err)
				}

				// Wait for operation status
				if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
					return retry.RetryableError(err)
				}
			}
			log.Printf("[INFO] Will delete dbaas logs input : %s/%s", serviceName, id)
			if err := config.OVHClient.Delete(endpoint, res); err != nil {
				return retry.RetryableError(err)
			}

			// Wait for operation status
			if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
				return retry.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccResourceDbaasLogsInput_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	name := "LOGSTASH"
	version := os.Getenv("OVH_DBAAS_LOGS_LOGSTASH_VERSION_TEST")
	title := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccResourceDbaasLogsInput_basic,
		serviceName,
		name,
		version,
		serviceName,
		title,
		desc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogsInput(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_input.input",
						"description",
						desc,
					),
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_input.input",
						"title",
						title,
					),
				),
			},
		},
	})
}
