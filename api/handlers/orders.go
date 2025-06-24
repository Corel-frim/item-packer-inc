package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/corel-frim/item-packer-inc/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type Orders struct {
	storage *storage.PackStorage
}

func NewOrders(storage *storage.PackStorage) *Orders {
	return &Orders{
		storage: storage,
	}
}

func (o *Orders) RegisterRoutes(app *fiber.App) {
	group := app.Group("/orders")
	group.Post("/items/:amount", o.CreateOrder)
	group.Get("", o.GetOrders)
}

// CreateOrder handles POST /order/items/{amount}
// @Summary Create an order
// @Description Create an order with the specified number of items
// @Tags orders
// @Produce json
// @Param amount path int true "Number of items"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Invalid amount"
// @Failure 404 {object} map[string]string "No packs available"
// @Router /order/items/{amount} [post]
func (o *Orders) CreateOrder(c *fiber.Ctx) error {
	path := c.Params("amount")
	amount, err := strconv.Atoi(path)
	if err != nil || amount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid amount"})
	}

	order, err := o.storage.CalculateOrder(amount)
	if err != nil {
		if errors.Is(err, storage.ErrNoPacksAvailable) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "No packs available"})
		}
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Internal server error"})
	}
	c.Set("Content-Type", "application/json")
	return c.Status(http.StatusOK).JSON(order)
}

// GetOrders handles GET /orders
// @Summary Get all orders
// @Description Retrieve a list of all orders
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Router /orders [get]
func (o *Orders) GetOrders(c *fiber.Ctx) error {
	orders := o.storage.GetOrders()

	c.Set("Content-Type", "application/json")
	return c.Status(http.StatusOK).JSON(orders)
}
