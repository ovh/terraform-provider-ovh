package ovh

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TestAccIPLoadbalancingRouteHTTPRulePlan = [][]map[string]interface{}{
	{
		{"DisplayName": "Test rule", "Field": "header", "Match": "is", "Negate": false, "Pattern": "example.com", "SubField": "Host"},
	},
}

type TestAccIPLoadbalancingRouteHTTPRule struct {
	ServiceName string
	RuleID      int    `json:"ruleId"`
	RouteID     int    `json:"routeId"`
	DisplayName string `json:"displayName"`
	Field       string `json:"field"`
	Match       string `json:"match"`
	Negate      bool   `json:"negate"`
	Pattern     string `json:"pattern"`
	SubField    string `json:"subField"`
}

type TestAccIPLoadbalancingRouteHTTPRuleWrapper struct {
	Expected *TestAccIPLoadbalancingRouteHTTPRule
}

func (w *TestAccIPLoadbalancingRouteHTTPRuleWrapper) Config() string {
	var config bytes.Buffer
	config.WriteString(fmt.Sprintf(`
		resource "ovh_iploadbalancing_http_route" "testroute" {
			service_name = "%s"
			display_name = "Test route"
			weight = 0

			action {
				status = 302
				target = "http://example.com"
				type = "redirect"
			}
		}

		resource "ovh_iploadbalancing_http_route_rule" "testrule" {
			service_name = "%s"
			route_id  = "${ovh_iploadbalancing_http_route.testroute.id}"
			display_name = "%s"
			field = "%s"
			match = "%s"
			negate = %t
			pattern = "%s"
			sub_field = "%s"
		}
	`, w.Expected.ServiceName,
		w.Expected.ServiceName,
		w.Expected.DisplayName,
		w.Expected.Field,
		w.Expected.Match,
		w.Expected.Negate,
		w.Expected.Pattern,
		w.Expected.SubField))

	return config.String()
}

func (rule *TestAccIPLoadbalancingRouteHTTPRule) MustEqual(compared *TestAccIPLoadbalancingRouteHTTPRule) error {
	if !reflect.DeepEqual(rule.DisplayName, compared.DisplayName) {
		return fmt.Errorf("DisplayName differs")
	}
	if !reflect.DeepEqual(rule.Field, compared.Field) {
		return fmt.Errorf("Field differs")
	}
	if !reflect.DeepEqual(rule.Match, compared.Match) {
		return fmt.Errorf("Match differs")
	}
	if !reflect.DeepEqual(rule.Negate, compared.Negate) {
		return fmt.Errorf("Negate differs")
	}
	if !reflect.DeepEqual(rule.Pattern, compared.Pattern) {
		return fmt.Errorf("Pattern differs")
	}
	if !reflect.DeepEqual(rule.SubField, compared.SubField) {
		return fmt.Errorf("SubField differs")
	}

	return nil
}

type TestAccIPLoadbalancingRouteHTTPRuleStep struct {
	Response *TestAccIPLoadbalancingRouteHTTPRule
	Expected *TestAccIPLoadbalancingRouteHTTPRule
}

func (w *TestAccIPLoadbalancingRouteHTTPRuleWrapper) TestStep(c map[string]interface{}) resource.TestStep {
	if val, ok := c["DisplayName"]; ok {
		w.Expected.DisplayName = val.(string)
	}
	if val, ok := c["Field"]; ok {
		w.Expected.Field = val.(string)
	}
	if val, ok := c["Match"]; ok {
		w.Expected.Match = val.(string)
	}
	if val, ok := c["Negate"]; ok {
		w.Expected.Negate = val.(bool)
	}
	if val, ok := c["Pattern"]; ok {
		w.Expected.Pattern = val.(string)
	}
	if val, ok := c["SubField"]; ok {
		w.Expected.SubField = val.(string)
	}
	expected := *w.Expected

	return resource.TestStep{
		Config: w.Config(),
		Check: resource.ComposeTestCheckFunc(
			w.TestCheck(expected),
		),
	}
}

func (w *TestAccIPLoadbalancingRouteHTTPRuleWrapper) TestCheck(expected TestAccIPLoadbalancingRouteHTTPRule) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		response := &TestAccIPLoadbalancingRouteHTTPRule{}
		name := "ovh_iploadbalancing_http_route_rule.testrule"
		resource, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s/rule/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.Attributes["route_id"], resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, response)
		if err != nil {
			return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
		}

		err = expected.MustEqual(response)
		if err != nil {
			return fmt.Errorf("%s %s state differs from expected : %s", name, resource.Primary.ID, err.Error())
		}
		return nil
	}
}

func (w *TestAccIPLoadbalancingRouteHTTPRuleWrapper) TestDestroy(state *terraform.State) error {
	leftovers := false
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_route_rule" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%d/rule/%s", os.Getenv("OVH_IPLB_SERVICE"), w.Expected.RouteID, resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			leftovers = true
		}
	}
	if leftovers {
		return fmt.Errorf("IpLoadbalancing http route rule still exists")
	}
	return nil
}

func newTestAccIPLoadbalancingRouteHTTPRuleWrapper() *TestAccIPLoadbalancingRouteHTTPRuleWrapper {
	return &TestAccIPLoadbalancingRouteHTTPRuleWrapper{
		Expected: &TestAccIPLoadbalancingRouteHTTPRule{ServiceName: os.Getenv("OVH_IPLB_SERVICE")},
	}
}

func TestAccIpLoadbalancingRouteHTTPRuleBasicCreate(t *testing.T) {
	for _, plan := range TestAccIPLoadbalancingRouteHTTPRulePlan {
		w := newTestAccIPLoadbalancingRouteHTTPRuleWrapper()
		var steps []resource.TestStep
		for _, tcase := range plan {
			steps = append(steps, w.TestStep(tcase))
		}
		resource.Test(t, resource.TestCase{
			Providers:    testAccProviders,
			CheckDestroy: w.TestDestroy,
			Steps:        steps,
		})
	}
}
