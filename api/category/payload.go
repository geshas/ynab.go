package category

// PayloadMonthCategory is the payload contract for updating a category for a month
type PayloadMonthCategory struct {
	Budgeted int64 `json:"budgeted"`
}

// PayloadCategory is the payload contract for updating a category
type PayloadCategory struct {
	Name            *string `json:"name,omitempty"`
	Note            *string `json:"note,omitempty"`
	CategoryGroupID *string `json:"category_group_id,omitempty"`
	GoalTarget      *int64  `json:"goal_target,omitempty"`
}

// PayloadCreateCategory is the payload contract for creating a category
type PayloadCreateCategory struct {
	Name            string  `json:"name"`
	CategoryGroupID string  `json:"category_group_id"`
	Note            *string `json:"note,omitempty"`
}

// PayloadCreateCategoryGroup is the payload contract for creating a category group
type PayloadCreateCategoryGroup struct {
	Name string `json:"name"`
}

// PayloadUpdateCategoryGroup is the payload contract for updating a category group
type PayloadUpdateCategoryGroup struct {
	Name *string `json:"name,omitempty"`
	Note *string `json:"note,omitempty"`
}
