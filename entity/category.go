package entity

type Category string

const (
	SportCategory Category = "sport"
)

func (c Category) IsValid() bool {
	// TODO: better approach instead of each time creation of this array for validation.
	validCategorySet := map[Category]bool{
		"sport": true,
	}

	if _, ok := validCategorySet[c]; ok {
		return true
	}

	return false
}
