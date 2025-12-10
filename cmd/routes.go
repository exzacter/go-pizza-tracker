package main

import "github.com/gin-gonic/gin"

func setupRotues(router *gin.Engine, h *Handler) {
	// this is the root path, eventually might want to add a landing pag with a button to make a new order
	router.GET("/", h.ServeNewOrder)
	// this is the post request from the form of creating the order
	router.POST("/new-order", h.HandleNewOrder)
	// the page of monitoring the pizza order status
	router.GET("/customer/:id", h.serveCustomer)

	// gets the router to look in the templates/static folder for the files
	router.Static("/static", "/templates/static")
}
