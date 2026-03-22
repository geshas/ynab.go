package plan_test

import (
	"fmt"
	"reflect"

	"github.com/geshas/ynab.go/api"

	"github.com/geshas/ynab.go"
)

func ExampleService_GetPlan() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	b, _ := c.Plan().GetPlan("<valid_plan_id>", nil)
	fmt.Println(reflect.TypeOf(b))

	// Output: *plan.Snapshot
}

func ExampleService_GetLastUsedPlan() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	b, _ := c.Plan().GetLastUsedPlan(nil)
	fmt.Println(reflect.TypeOf(b))

	// Output: *plan.Snapshot
}

func ExampleService_GetPlan_filtered() {
	c := ynab.NewClient("<valid_ynab_access_token>")

	f := api.Filter{LastKnowledgeOfServer: 10}
	b, _ := c.Plan().GetPlan("<valid_plan_id>", &f)
	fmt.Println(reflect.TypeOf(b))

	// Output: *plan.Snapshot
}

func ExampleService_GetPlans() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	plans, _ := c.Plan().GetPlans()
	fmt.Println(reflect.TypeOf(plans))

	// Output: []*plan.Summary
}

func ExampleService_GetPlanSettings() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	s, _ := c.Plan().GetPlanSettings("<valid_plan_id>")
	fmt.Println(reflect.TypeOf(s))

	// Output: *plan.Settings
}
