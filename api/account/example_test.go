package account_test

import (
	"fmt"
	"reflect"

	"github.com/geshas/ynab.go"
	"github.com/geshas/ynab.go/api"
)

func ExampleService_GetAccount() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	account, _ := c.Account().GetAccount("<valid_plan_id>", "<valid_account_id>")
	fmt.Println(reflect.TypeOf(account))

	// Output: *account.Account
}

func ExampleService_GetAccounts() {
	c := ynab.NewClient("<valid_ynab_access_token>")
	f := &api.Filter{LastKnowledgeOfServer: 10}
	snapshot, _ := c.Account().GetAccounts("<valid_plan_id>", f)
	fmt.Println(reflect.TypeOf(snapshot))

	// Output: *account.SearchResultSnapshot
}
