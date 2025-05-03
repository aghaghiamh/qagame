package entity

import "github.com/samber/lo"

type Category string

const (
	SportCategory Category = "sport"
	// MusicCategory Category = "music"
)

func AllCategories() []Category {
	return []Category{
		SportCategory,
		// Warning: Add more categories here as you develop
		// MusicCategory,
	}
}

func (c Category) IsValid() bool {
	return lo.Contains(AllCategories(), c)
}
