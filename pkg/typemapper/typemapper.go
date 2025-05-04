package typemapper

import "github.com/samber/lo"

// ArrayMapper converts a slice of one type to another using a conversion function
func ArrayMapper[I any, O any](inputs []I, converter func(I) O) []O {
	return lo.Map(inputs, func(input I, _ int) O {
		return converter(input)
	})
}
