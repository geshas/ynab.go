package payee_test

import (
	"fmt"
	"reflect"

	"github.com/geshas/ynab.go/api"

	"github.com/geshas/ynab.go"
	"github.com/geshas/ynab.go/api/payee"
)

func ExampleService_GetPayee() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	p, _ := c.Payee().GetPayee("<valid_plan_id>", "<valid_payee_id>")
	fmt.Println(reflect.TypeOf(p))

	// Output: *payee.Payee
}

func ExampleService_CreatePayee() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	p, _ := c.Payee().CreatePayee("<valid_plan_id>", payee.PayloadPayee{Name: "Coffee Shop"})
	fmt.Println(reflect.TypeOf(p))

	// Output: *payee.Payee
}

func ExampleService_GetPayees() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	f := &api.Filter{LastKnowledgeOfServer: 10}
	payees, _ := c.Payee().GetPayees("<valid_plan_id>", f)
	fmt.Println(reflect.TypeOf(payees))

	// Output: *payee.SearchResultSnapshot
}

func ExampleService_GetPayeeLocation() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	l, _ := c.Payee().GetPayeeLocation("<valid_plan_id>", "<valid_payee_location_id>")
	fmt.Println(reflect.TypeOf(l))

	// Output: *payee.Location
}

func ExampleService_GetPayeeLocations() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	locations, _ := c.Payee().GetPayeeLocations("<valid_plan_id>")
	fmt.Println(reflect.TypeOf(locations))

	// Output: []*payee.Location
}

func ExampleService_GetPayeeLocationsByPayee() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	locations, _ := c.Payee().GetPayeeLocationsByPayee("<valid_plan_id>", "<valid_payee_id>")
	fmt.Println(reflect.TypeOf(locations))

	// Output: []*payee.Location
}
