package ovh

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const (
	test_prefix = "testacc-terraform"
)

// Custom type used to check that a plan is empty in acceptance tests
var _ plancheck.PlanCheck = expectEmptyPlan{}

type expectEmptyPlan struct{}

func (e expectEmptyPlan) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result error

	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			result = errors.Join(result, fmt.Errorf("expected empty plan, but %s has planned action(s): %v", rc.Address, rc.Change.Actions))
		}
	}

	resp.Error = result
}

func ExpectEmptyPlan() plancheck.PlanCheck {
	return expectEmptyPlan{}
}
