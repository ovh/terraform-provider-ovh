package ovh

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Filters sit on a real mailbox, so the sweeper only removes the ones these
// tests created, recognised by this prefix.
const testAccEmailDomainFilterPrefix = "tf-acc-test-filter"

func init() {
	resource.AddTestSweepers("ovh_email_domain_filter", &resource.Sweeper{
		Name: "ovh_email_domain_filter",
		F:    testSweepEmailDomainFilter,
	})
}

func testSweepEmailDomainFilter(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")
	account := os.Getenv("OVH_EMAIL_DOMAIN_ACCOUNT_TEST")
	if domain == "" || account == "" {
		log.Print("[DEBUG] OVH_EMAIL_DOMAIN_TEST or OVH_EMAIL_DOMAIN_ACCOUNT_TEST is not set. No email domain filters to sweep")
		return nil
	}

	var names []string
	if err := client.Get(emailFilterBase(domain, account), &names); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", emailFilterBase(domain, account), err)
	}

	r := &emailDomainFilterResource{config: &Config{OVHClient: client}}

	var tasks []int64
	for _, name := range names {
		// OVHcloud returns names lowercased, whatever they were created as.
		if !strings.HasPrefix(strings.ToLower(name), testAccEmailDomainFilterPrefix) {
			continue
		}

		log.Printf("[DEBUG] Deleting email domain filter %q on %s@%s", name, account, domain)
		ids, err := r.del(emailFilterPath(domain, account, name))
		if err != nil && !emailFilterNotFound(err) {
			return fmt.Errorf("Error calling delete on %s:\n\t %q", emailFilterPath(domain, account, name), err)
		}
		tasks = append(tasks, ids...)
	}

	// Deletes are queued tasks. Wait for them like Delete does, so the sweep is
	// really over when it returns.
	if err := r.waitForTasks(context.Background(), domain, tasks); err != nil {
		return fmt.Errorf("error waiting for the swept filters to be deleted: %w", err)
	}

	return nil
}

// Checks that the environment variables needed for the filter acceptance tests
// are set. They need an existing mailbox, since a filter cannot exist without one.
func testAccPreCheckEmailDomainAccount(t *testing.T) {
	testAccPreCheckEmailDomain(t)
	checkEnvOrSkip(t, "OVH_EMAIL_DOMAIN_ACCOUNT_TEST")
}

func TestAccEmailDomainFilter_Basic(t *testing.T) {
	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")
	account := os.Getenv("OVH_EMAIL_DOMAIN_ACCOUNT_TEST")
	name := acctest.RandomWithPrefix(testAccEmailDomainFilterPrefix)
	res := "ovh_email_domain_filter.test"

	// Every rule tests for the filter's own random name, which no real mail
	// carries, so the filter never acts on anything even on a mailbox in use.
	subject := [3]string{"Subject", "contains", name}
	from := [3]string{"From", "contains", name}
	custom := [3]string{"X-Tf-Acc-Test", "contains", name}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckEmailDomainAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEmailDomainFilterDestroyed(domain, account, name),
		Steps: []resource.TestStep{
			// A bare local part would resolve against OVHcloud's own MX host.
			{
				Config:      testAccEmailDomainFilterConfig(domain, account, name, 1, true, "account", "postmaster", subject),
				ExpectError: regexp.MustCompile(`must be a full address`),
			},
			// accept and delete take no parameter.
			{
				Config:      testAccEmailDomainFilterConfig(domain, account, name, 1, true, "accept", "someone@example.com", subject),
				ExpectError: regexp.MustCompile(`takes no action_param`),
			},
			// Creation with two rules: made inactive, completed, then activated.
			{
				Config: testAccEmailDomainFilterConfig(domain, account, name, 1, true, "accept", "", subject, from),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", fmt.Sprintf("%s/%s/%s", domain, account, name)),
					resource.TestCheckResourceAttr(res, "priority", "1"),
					resource.TestCheckResourceAttr(res, "active", "true"),
					resource.TestCheckResourceAttr(res, "action", "accept"),
					resource.TestCheckNoResourceAttr(res, "action_param"),
					resource.TestCheckResourceAttr(res, "rules.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(res, "rules.*", map[string]string{
						"header": "Subject", "operand": "contains", "value": name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(res, "rules.*", map[string]string{
						"header": "From", "operand": "contains", "value": name,
					}),
				),
			},
			// Priority, activity and a rule change in place.
			{
				Config: testAccEmailDomainFilterConfig(domain, account, name, 2, false, "accept", "", subject, custom),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate)},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "priority", "2"),
					resource.TestCheckResourceAttr(res, "active", "false"),
					resource.TestCheckResourceAttr(res, "rules.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(res, "rules.*", map[string]string{
						"header": "X-Tf-Acc-Test", "operand": "contains", "value": name,
					}),
				),
			},
			// OVHcloud matches names case-insensitively, so a change of case alone
			// is applied in place rather than replacing the filter.
			{
				Config: testAccEmailDomainFilterConfig(domain, account, strings.ToUpper(name), 2, false, "accept", "", subject, custom),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate)},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "name", strings.ToUpper(name)),
					resource.TestCheckResourceAttr(res, "id", fmt.Sprintf("%s/%s/%s", domain, account, name)),
				),
			},
			// There is no endpoint to change the action, so it forces a new filter.
			{
				Config: testAccEmailDomainFilterConfig(domain, account, strings.ToUpper(name), 2, true, "delete", "", subject, custom),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(res, plancheck.ResourceActionDestroyBeforeCreate)},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "action", "delete"),
					resource.TestCheckResourceAttr(res, "active", "true"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateIdFunc: testAccEmailDomainFilterImportID(res),
				ImportStateVerify: true,
			},
		},
	})
}

// The API schema says a delete returns an array of tasks, but the API returns a
// single task object, which is what the first acceptance run tripped over.
func TestEmailFilterTaskIDs(t *testing.T) {
	cases := map[string][]int64{
		`{"id": 3793799, "action": "del", "account": "postmaster"}`: {3793799},
		`[{"id": 1}, {"id": 2}]`:                                    {1, 2},
		`[]`:                                                        {0},
		`null`:                                                      {0},
		`"unexpected"`:                                              {0},
	}

	for raw, want := range cases {
		got := emailFilterTaskIDs([]byte(raw))
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("emailFilterTaskIDs(%s) = %v, want %v", raw, got, want)
		}
	}
}

// testAccCheckEmailDomainFilterDestroyed checks that destroy waited for the
// filter to be gone, rather than returning while its deletion was queued.
func testAccCheckEmailDomainFilterDestroyed(domain, account, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		var f emailDomainFilterAPI
		err := testAccOVHClient.Get(emailFilterPath(domain, account, name), &f)
		if err == nil {
			return fmt.Errorf("filter %q still exists on %s@%s after destroy", name, account, domain)
		}
		if !emailFilterNotFound(err) {
			return fmt.Errorf("error reading filter %q: %w", name, err)
		}
		return nil
	}
}

func testAccEmailDomainFilterImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}

		a := rs.Primary.Attributes
		return fmt.Sprintf("%s/%s/%s", a["domain"], a["account_name"], a["name"]), nil
	}
}

func testAccEmailDomainFilterConfig(domain, account, name string, priority int, active bool, action, actionParam string, rules ...[3]string) string {
	var param string
	if actionParam != "" {
		param = fmt.Sprintf("  action_param = %q\n", actionParam)
	}

	var body strings.Builder
	for _, r := range rules {
		fmt.Fprintf(&body, "    {\n      header  = %q\n      operand = %q\n      value   = %q\n    },\n", r[0], r[1], r[2])
	}

	return fmt.Sprintf(`
resource "ovh_email_domain_filter" "test" {
  domain       = %q
  account_name = %q
  name         = %q
  priority     = %d
  active       = %t
  action       = %q
%s  rules = [
%s  ]
}`, domain, account, name, priority, active, action, param, body.String())
}
