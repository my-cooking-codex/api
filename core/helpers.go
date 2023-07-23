package core

import "github.com/labstack/echo/v4"

// Ease of use method, when binding & validation is needed.
func BindAndValidate(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	} else if err := ctx.Validate(i); err != nil {
		return err
	}
	return nil
}

func ValueOrDefault[V comparable](value *V, def V) V {
	if value == nil {
		return def
	}
	return *value
}

func PopElement[T comparable](toPop T, items []T) []T {
	var newItems = make([]T, 0, len(items))
	for _, item := range items {
		if item != toPop {
			newItems = append(newItems, item)
		}
	}
	return newItems
}
