package main

// in future i want to change so that a customer can login and then input customer data, this allows a discount and also helps with learnign authetnication
// so in future the form to create will include customer name, phone, address, and then the items which is the pizza/s that they will order
type OrderFormData struct {
	PizzaTypes           []string
	PizzaSizes           []string
	PizzaDefaultToppings map[string][]string
	ToppingCategories    map[string][]string
	DietaryRequirements  []string
	Allergies            []string
}
