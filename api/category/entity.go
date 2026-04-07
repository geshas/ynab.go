// Package category implements category entities and services
package category // import "github.com/geshas/ynab.go/api/category"

import "github.com/geshas/ynab.go/api"

// Category represents a category for a budget
type Category struct {
	ID                string `json:"id"`
	CategoryGroupID   string `json:"category_group_id"`
	CategoryGroupName string `json:"category_group_name"`
	Name              string `json:"name"`
	Hidden            bool   `json:"hidden"`
	// Budgeted Budgeted amount in current month in milliunits format
	Budgeted int64 `json:"budgeted"`
	// BudgetedFormatted Budgeted amount formatted in the plan's currency format
	BudgetedFormatted *string `json:"budgeted_formatted"`
	// BudgetedCurrency Budgeted amount as a decimal currency amount
	BudgetedCurrency *float64 `json:"budgeted_currency"`
	// Activity Activity amount in current month in milliunits format
	Activity int64 `json:"activity"`
	// ActivityFormatted Activity amount formatted in the plan's currency format
	ActivityFormatted *string `json:"activity_formatted"`
	// ActivityCurrency Activity amount as a decimal currency amount
	ActivityCurrency *float64 `json:"activity_currency"`
	// Balance Balance in current month in milliunits format
	Balance int64 `json:"balance"`
	// BalanceFormatted Available balance formatted in the plan's currency format
	BalanceFormatted *string `json:"balance_formatted"`
	// BalanceCurrency Available balance as a decimal currency amount
	BalanceCurrency *float64 `json:"balance_currency"`
	// Deleted Deleted category groups will only be included in delta requests
	Deleted bool `json:"deleted"`

	Note *string `json:"note"`
	// OriginalCategoryGroupID If category is hidden this is the ID of the category
	// group it originally belonged to before it was hidden
	OriginalCategoryGroupID *string `json:"original_category_group_id"`

	GoalType *Goal `json:"goal_type"`
	// GoalCreationMonth the month a goal was created
	GoalCreationMonth *api.Date `json:"goal_creation_month"`
	// GoalTarget the goal target amount in milliunits
	GoalTarget *int64 `json:"goal_target"`
	// GoalTargetFormatted Goal target amount formatted in the plan's currency format
	GoalTargetFormatted *string `json:"goal_target_formatted"`
	// GoalTargetCurrency Goal target amount as a decimal currency amount
	GoalTargetCurrency *float64 `json:"goal_target_currency"`
	// GoalTargetMonth if the goal type is GoalTargetCategoryBalanceByDate,
	// this is the target month for the goal to be completed
	GoalTargetMonth *api.Date `json:"goal_target_month"`
	// GoalPercentageComplete the percentage completion of the goal
	GoalPercentageComplete *int32 `json:"goal_percentage_complete"`
	// GoalNeedsWholeAmount indicates monthly rollover behavior for "NEED"-type goals
	// When true: goal asks for target amount in new month ("Set Aside")
	// When false: uses previous month category funding ("Refill")
	GoalNeedsWholeAmount *bool `json:"goal_needs_whole_amount"`
	// GoalDay day offset modifier for goal's due date
	// For weekly goals (cadence=2): specifies day of week (0=Sunday, 6=Saturday)
	// For other goals: specifies day of month (1=1st, 31=31st, null=last day)
	GoalDay *int32 `json:"goal_day"`
	// GoalCadence the goal cadence (0-14 range)
	// Values 0,1,2,13: repeats every goal_cadence * goal_cadence_frequency
	// Values 3-12,14: repeats every goal_cadence (frequency ignored)
	GoalCadence *int32 `json:"goal_cadence"`
	// GoalCadenceFrequency goal cadence frequency multiplier
	// Used with cadence values 0,1,2,13 to determine repeat frequency
	GoalCadenceFrequency *int32 `json:"goal_cadence_frequency"`
	// GoalMonthsToBudget number of months left in current goal period (including current month)
	GoalMonthsToBudget *int32 `json:"goal_months_to_budget"`
	// GoalUnderFunded amount of funding still needed in current month to stay on track (milliunits)
	GoalUnderFunded *int64 `json:"goal_under_funded"`
	// GoalUnderFundedFormatted Goal underfunded amount formatted in the plan's currency format
	GoalUnderFundedFormatted *string `json:"goal_under_funded_formatted"`
	// GoalUnderFundedCurrency Goal underfunded amount as a decimal currency amount
	GoalUnderFundedCurrency *float64 `json:"goal_under_funded_currency"`
	// GoalOverallFunded total amount funded towards goal within current goal period (milliunits)
	GoalOverallFunded *int64 `json:"goal_overall_funded"`
	// GoalOverallFundedFormatted Goal funded amount formatted in the plan's currency format
	GoalOverallFundedFormatted *string `json:"goal_overall_funded_formatted"`
	// GoalOverallFundedCurrency Goal funded amount as a decimal currency amount
	GoalOverallFundedCurrency *float64 `json:"goal_overall_funded_currency"`
	// GoalOverallLeft amount of funding still needed to complete goal within current goal period (milliunits)
	GoalOverallLeft *int64 `json:"goal_overall_left"`
	// GoalOverallLeftFormatted Goal remaining amount formatted in the plan's currency format
	GoalOverallLeftFormatted *string `json:"goal_overall_left_formatted"`
	// GoalOverallLeftCurrency Goal remaining amount as a decimal currency amount
	GoalOverallLeftCurrency *float64 `json:"goal_overall_left_currency"`
}

// Group represents a resumed category group for a budget
type Group struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
	// Deleted Deleted category groups will only be included in delta requests
	Deleted bool `json:"deleted"`
}

// GroupWithCategories represents a category group for a budget
type GroupWithCategories struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
	// Deleted Deleted category groups will only be included in delta requests
	Deleted bool `json:"deleted"`

	Categories []*Category `json:"categories"`
}

// SearchResultSnapshot represents a versioned snapshot for an account search
type SearchResultSnapshot struct {
	GroupWithCategories []*GroupWithCategories
	ServerKnowledge     uint64
}
