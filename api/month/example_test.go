package month_test

import (
	"fmt"
	"reflect"

	"github.com/geshas/ynab.go"
	"github.com/geshas/ynab.go/api"
)

//nolint:govet
func Example() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	d, _ := api.DateFromString("2010-01-01")
	m, _ := c.Month().GetMonth("<valid_plan_id>", d)
	fmt.Println(reflect.TypeOf(m))

	f := &api.Filter{LastKnowledgeOfServer: 10}
	months, _ := c.Month().GetMonths("<valid_plan_id>", f)
	fmt.Println(reflect.TypeOf(months))

	// Output: *month.Month
	// *month.SearchResultSnapshot
}
