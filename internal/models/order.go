package models

import (
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
	"time"
)

const (
	OrderStatusPlaced       = "Order Placed"
	OrderStatusPreparing    = "Preparing"
	OrderStatusCooking      = "Cooking"
	OrderStatusQualityCheck = "Quality Check"
	OrderStatusReady        = "Ready"
)

var (
	OrderStatuses = []string{
		OrderStatusPlaced,
		OrderStatusPreparing,
		OrderStatusCooking,
		OrderStatusQualityCheck,
		OrderStatusReady,
	}
	PizzaSauces = []string{
		"Tomato Sauce", "BBQ Sauce", "Buffalo Sauce", "Garlic Oil", "Truffle Oil", "Pesto",
	}
	Cheeses = []string{
		"Mozzarella", "Vegan Cheese", "Extra Cheese", "Parmesan", "Gorgonzola", "Ricotta", "Feta",
	}
	PizzaTypes = []string{
		"Margherita",
		"Pepperoni",
		"Vegetarian",
		"Hawaiian",
		"BBQ Chicken",
		"Meat Lovers",
		"Buffalo Chicken",
		"Supreme",
		"Truffle Mushroom",
		"Four Cheese",
		"Vegan Pizza",
		"Vegan Meat Lovers",
		"Vegan Garden",
		"Make Your Own",
		"Garlic",
	}
	PizzaCrust = []string{"Thin", "Regular", "Deep Dish", "Cheesy Crust", "Vegan Cheesy Crust"}
	PizzaSizes = []string{"Small", "Medium", "Large", "X-Large", "Family"}
	// map[string][]string allows me to create a junction table
	ToppingCategories = map[string][]string{
		"Meats": {
			"Pepperoni", "Sausage", "Chicken", "Bacon", "Ham", "Mince", "Anchovies",
		},
		"Vegan Meats": {
			"Vegan Pepperoni", "Vegan Chicken", "Vegan Bacon", "Vegan Ham", "Vegan Mince",
		},
		"Vegetables": {
			"Mushroom", "Red Onion", "White Onion", "Capsicum", "Zucchini",
			"Olives", "Jalapenos", "Pumpkin", "Spinach", "Pineapple",
			"Tomatoes", "Basil", "Corn", "Rocket", "Garlic Slices",
		},
		"Cheeses": {
			"Mozzarella", "Extra Cheese", "Vegan Cheese", "Parmesan", "Gorgonzola", "Ricotta", "Feta",
		},
		"Sauces": {
			"Tomato Sauce", "BBQ Sauce", "Buffalo Sauce", "Garlic Oil", "Truffle Oil", "Pesto",
		},
	}
	PizzaDefaultToppings = map[string][]string{
		"Margherita":        {"Tomato Sauce", "Mozzarella", "Tomatoes", "Basil"},
		"Pepperoni":         {"Tomato Sauce", "Mozzarella", "Pepperoni"},
		"Vegetarian":        {"Tomato Sauce", "Mozzarella", "Mushroom", "Capsicum", "Red Onion", "Olives", "Zucchini"},
		"Hawaiian":          {"Tomato Sauce", "Mozzarella", "Ham", "Pineapple"},
		"BBQ Chicken":       {"BBQ Sauce", "Mozzarella", "Chicken", "Red Onion", "Capsicum", "Bacon"},
		"Meat Lovers":       {"Tomato Sauce", "Mozzarella", "Pepperoni", "Sausage", "Bacon", "Ham", "Mince"},
		"Buffalo Chicken":   {"Buffalo Sauce", "Mozzarella", "Chicken", "Red Onion", "Jalapenos"},
		"Supreme":           {"Tomato Sauce", "Mozzarella", "Pepperoni", "Sausage", "Mushroom", "Capsicum", "Red Onion", "Olives"},
		"Truffle Mushroom":  {"Truffle Oil", "Mozzarella", "Mushroom", "Garlic Oil", "Rocket"},
		"Four Cheese":       {"Tomato Sauce", "Mozzarella", "Parmesan", "Gorgonzola", "Ricotta"},
		"Vegan Pizza":       {"Tomato Sauce", "Vegan Cheese", "White Onion", "Basil", "Capsicum", "Garlic Slices", "Corn", "Zucchini"},
		"Vegan Meat Lovers": {"BBQ Sauce", "Vegan Cheese", "White Onion", "Vegan Ham", "Vegan Mince", "Vegan Chicken", "Vegan Bacon", "Vegan Pepperoni"},
		"Vegan Garden":      {"Pesto", "Vegan Cheese", "Red Onion", "Mushroom", "Pumpkin", "Capsicum", "Zucchini", "Spinach", "Corn", "Basil", "Rocket"},
		"Garlic":            {"Garlic", "Mozzarella", "Garlic Oil"},
		"Make Your Own":     {},
	}
	DietaryRequirements = []string{"Vegetarian", "Vegan", "Gluten-Free", "Dairy-Free", "Nut-Free", "Halal", "Kosher"}

	Allergies = []string{"Gluten", "Dairy", "Nuts", "Peanuts", "Shellfish", "Soy", "Eggs", "Fish", "Sesame"}
)

type OrderModel struct {
	DB *gorm.DB
}

type Order struct {
	ID           string      `gorm:"primaryKey; size:14" json:"id"`
	Status       string      `gorm:"not null" json:"status"`
	CustomerName string      `gorm:"not null" json:"customerName"`
	Phone        string      `gorm:"not null" json:"phone"`
	Address      string      `gorm:"not null" json:"adress"`
	Items        []OrderItem `gorm:"foreignKey:OrderID" json:"pizzas"`
	CreatedAt    time.Time   `json:"createdAt"`
}

type OrderItem struct {
	ID                 string                        `gorm:"primaryKey; size:14" json:"id"`
	OrderID            string                        `gorm:"index;not null" json:"orderId"`
	Size               string                        `gorm:"not null" json:"size"`
	Pizza              string                        `gorm:"not null" json:"pizza"`
	Instructions       string                        `json:"instruction"`
	DietaryRequirement []OrderItemDietaryRequirement `gorm:"foreignKey:OrderItemID" json:"dietaryRequirement"`
	Toppings           []OrderItemTopping            `gorm:"foreignKey:OrderItemID" json:"toppings"`
	Allergies          []OrderItemAllergy            `gorm:"foreignKey:OrderItemID" json:"allergies"`
}

// Junction table for toppings
type OrderItemTopping struct {
	ID          string `gorm:"primaryKey;size:14" json:"id"`
	OrderItemID string `gorm:"index;not null" json:"orderItemId"`
	Topping     string `gorm:"not null" json:"topping"`
	IsExtra     bool   `gorm:"default:false" json:"isExtra"`
}

// Junction table for dietary requirements
type OrderItemDietaryRequirement struct {
	ID                 string `gorm:"primaryKey;size:14" json:"id"`
	OrderItemID        string `gorm:"index;not null" json:"orderItemId"`
	DietaryRequirement string `gorm:"not null" json:"dietaryRequirement"`
}

// Junction table for allergies
type OrderItemAllergy struct {
	ID          string `gorm:"primaryKey;size:14" json:"id"`
	OrderItemID string `gorm:"index;not null" json:"orderItemId"`
	Allergy     string `gorm:"not null" json:"allergy"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = shortid.MustGenerate()
	}

	return nil
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = shortid.MustGenerate()
	}

	return nil
}

func (t *OrderItemTopping) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = shortid.MustGenerate()
	}
	return nil
}

func (d *OrderItemDietaryRequirement) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = shortid.MustGenerate()
	}
	return nil
}

func (a *OrderItemAllergy) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = shortid.MustGenerate()
	}
	return nil
}
func (o *OrderModel) CreateOrder(order *Order) error {
	return o.DB.Create(order).Error
}

func (o *OrderModel) GetOrder(id string) (*Order, error) {
	var order Order

	err := o.DB.
		Preload("Items.Toppings").Preload("Items.DietaryRequirement").Preload("Items.Allergies").First(&order, "id = ?", id).Error
	return &order, err
}
