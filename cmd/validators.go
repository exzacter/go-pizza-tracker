package main

import (
	"go-pizza-tracker/internal/models"
	"slices"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.registerValidation("valid_pizza_type", createSliceValidator(models.PizzaTypes))
		v.registerValidation("valid_pizza_size", createSliceValidator(models.PizzaSizes))
		v.registerValidation("valid_topping", createToppingValidator())
		v.registerValidation("valid_pizza_size", createSliceValidator(models.DietaryRequirements))
		v.registerValidation("valid_pizza_size", createSliceValidator(models.Allergies))
	}
}

func createSliceValidator(allowedValues []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return slices.Contains(allowedValues, fl.Field().String())
	}
}

func createToppingValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		topping := fl.Field().String()
		for _, toppings := range models.ToppingCategories {
			if slices.Contains(toppings, topping) {
				return true
			}
		}
		return false
	}
}
