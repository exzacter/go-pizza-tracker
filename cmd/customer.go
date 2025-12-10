package main

import (
	"go-pizza-tracker/internal/models"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CustomerData struct {
	Title    string
	Order    models.Order
	Statuses []string
}

// in future i want to change so that a customer can login and then input customer data, this allows a discount and also helps with learnign authetnication
// so in future the form to create will include customer name, phone, address, and then the items which is the pizza/s that they will order
type OrderFormData struct {
	PizzaTypes           []string
	PizzaSizes           []string
	PizzaCrust           []string
	PizzaCheese          []string
	PizzaDefaultToppings map[string][]string
	ToppingCategories    map[string][]string
	DietaryRequirements  []string
	Allergies            []string
}

type OrderRequest struct {
	Name                string   `form:"name" binding:"required,min=2,max=100"`
	Phone               string   `form:"phone" binding:"required,min=9,max=20"`
	Address             string   `form:"address" binding:"required,min=5,max=200"`
	Sizes               []string `form:"size" binding:"required,min=1,dive,valid_pizza_size"`
	PizzaTypes          []string `form:"pizza" binding:"required,min=1,dive,valid_pizza_type"`
	Crusts              []string `form:"crust"`
	Instructions        []string `form:"instruction"`
	Toppings            []string `form:"topping"`
	DietaryRequirements []string `form:"dietary"`
	Allergies           []string `form:"allergy"`
}

func (h *Handler) ServeNewOrder(c *gin.Context) {
	c.HTML(http.StatusOK, "order.tmpl", OrderFormData{
		PizzaTypes:           models.PizzaTypes,
		PizzaSizes:           models.PizzaSizes,
		PizzaCrust:           models.PizzaCrust,
		PizzaCheese:          models.Cheeses,
		PizzaDefaultToppings: models.PizzaDefaultToppings,
		ToppingCategories:    models.ToppingCategories,
		DietaryRequirements:  models.DietaryRequirements,
		Allergies:            models.Allergies,
	})
}

func (h *Handler) HandleNewOrder(c *gin.Context) {
	var form OrderRequest
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse toppings per pizza: format is "pizzaIndex:topping"
	toppingsMap := make(map[int][]string)
	for _, t := range form.Toppings {
		parts := strings.SplitN(t, ":", 2)
		if len(parts) != 2 {
			continue
		}
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		toppingsMap[index] = append(toppingsMap[index], parts[1])
	}

	orderItems := make([]models.OrderItem, len(form.Sizes))
	for i := range orderItems {
		crust := "Regular"
		if i < len(form.Crusts) && form.Crusts[i] != "" {
			crust = form.Crusts[i]
		}

		instruction := ""
		if i < len(form.Instructions) {
			instruction = form.Instructions[i]
		}

		orderItems[i] = models.OrderItem{
			Size:         form.Sizes[i],
			Pizza:        form.PizzaTypes[i],
			Crust:        crust,
			Instructions: instruction,
		}

		// Get toppings for this pizza
		pizzaToppings := toppingsMap[i]

		// If no toppings selected (customize not opened), use defaults
		if len(pizzaToppings) == 0 {
			pizzaToppings = models.PizzaDefaultToppings[form.PizzaTypes[i]]
		}

		// Determine which are default vs extra
		defaultToppings := models.PizzaDefaultToppings[form.PizzaTypes[i]]

		for _, topping := range pizzaToppings {
			isExtra := true
			for _, def := range defaultToppings {
				if def == topping {
					isExtra = false
					break
				}
			}

			orderItems[i].Toppings = append(orderItems[i].Toppings, models.OrderItemTopping{
				Topping: topping,
				IsExtra: isExtra,
			})
		}

		// Dietary requirements
		for _, dietary := range form.DietaryRequirements {
			orderItems[i].DietaryRequirement = append(orderItems[i].DietaryRequirement, models.OrderItemDietaryRequirement{
				DietaryRequirement: dietary,
			})
		}

		// Allergies
		for _, allergy := range form.Allergies {
			orderItems[i].Allergies = append(orderItems[i].Allergies, models.OrderItemAllergy{
				Allergy: allergy,
			})
		}
	}

	order := models.Order{
		CustomerName: form.Name,
		Phone:        form.Phone,
		Address:      form.Address,
		Status:       models.OrderStatusPlaced,
		Items:        orderItems,
	}

	if err := h.orders.CreateOrder(&order); err != nil {
		slog.Error("Failed to create order", "Error", err)
		c.String(http.StatusInternalServerError, "Something went wrong")
		return
	}

	slog.Info("Order Created", "orderID", order.ID, "Customer", order.CustomerName)
	c.Redirect(http.StatusSeeOther, "/customer/"+order.ID)
}

func (h *Handler) serveCustomer(c *gin.Context) {
	orderID := c.Param("id")

	if orderID == "" {
		c.String(http.StatusBadRequest, "Order ID is required")
	}

	order, err := h.orders.GetOrder(orderID)
	if err != nil {
		c.String(http.StatusNotFound, "Order Not Found")
		return
	}

	c.HTML(http.StatusOK, "customer.tmpl", CustomerData{
		Title:    "Pizza Order Status" + orderID,
		Order:    *order,
		Statuses: models.OrderStatuses,
	})
}
