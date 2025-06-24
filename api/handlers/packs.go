package handlers

import (
	"errors"
	"net/http"

	"github.com/corel-frim/item-packer-inc/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type Packs struct {
	storage *storage.PackStorage
}

func NewPacks(storage *storage.PackStorage) *Packs {
	return &Packs{
		storage: storage,
	}
}

func (p *Packs) RegisterRoutes(app *fiber.App) {
	group := app.Group("/packs")
	group.Get("", p.GetPacks)
	group.Post("/:amount", p.AddPack)
	group.Put("/:oldAmount/:newAmount", p.UpdatePack)
	group.Delete("/:amount", p.DeletePack)
}

// GetPacks handles GET /packs
// @Summary Get all available packs
// @Description Get a list of all available packs
// @Tags packs
// @Produce json
// @Success 200 {array} models.Pack
// @Router /packs [get]
func (p *Packs) GetPacks(c *fiber.Ctx) error {
	packs := p.storage.GetPacks()
	return c.Status(http.StatusOK).JSON(packs)
}

// AddPack handles POST /packs/{amount}
// @Summary Add a new pack
// @Description Add a new pack with the specified amount
// @Tags packs
// @Produce json
// @Param amount path int true "Pack amount"
// @Success 201 {object} models.Pack
// @Failure 400 {object} map[string]string "Invalid amount"
// @Failure 409 {object} map[string]string "Limit for packs reached"
// @Router /packs/{amount} [post]
func (p *Packs) AddPack(c *fiber.Ctx) error {
	amount, err := c.ParamsInt("amount")
	if err != nil || amount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid amount"})
	}

	err = p.storage.AddPack(amount)
	if err != nil {
		return c.Status(http.StatusConflict).JSON(map[string]string{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(map[string]int{"amount": amount})
}

// UpdatePack handles PUT /packs/{oldAmount}/{newAmount}
// @Summary Update a pack
// @Description Update a pack's amount
// @Tags packs
// @Produce json
// @Param oldAmount path int true "Current pack amount"
// @Param newAmount path int true "New pack amount"
// @Success 200 {object} models.Pack
// @Failure 400 {object} map[string]string "Invalid amount"
// @Failure 404 {object} map[string]string "Pack not found"
// @Failure 409 {object} map[string]string "Pack with new amount already exists"
// @Router /packs/{oldAmount}/{newAmount} [put]
func (p *Packs) UpdatePack(c *fiber.Ctx) error {
	oldAmount, err := c.ParamsInt("oldAmount")
	if err != nil || oldAmount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid old amount"})
	}
	newAmount, err := c.ParamsInt("newAmount")
	if err != nil || newAmount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid new amount"})
	}

	err = p.storage.UpdatePack(oldAmount, newAmount)
	if err == nil {
		return c.Status(http.StatusOK).JSON(map[string]int{"oldAmount": oldAmount, "amount": newAmount})
	}

	// if err != nil
	switch {
	case errors.Is(err, storage.ErrPackNotFound):
		return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "Pack not found"})
	case errors.Is(err, storage.ErrPackExists):
		return c.Status(http.StatusConflict).JSON(map[string]string{"error": "Pack with new amount already exists"})
	default:
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to update pack"})
	}
}

// DeletePack handles DELETE /packs/{amount}
// @Summary Delete a pack
// @Description Delete a pack with the specified amount
// @Tags packs
// @Produce json
// @Param amount path int true "Pack amount"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid amount"
// @Router /packs/{amount} [delete]
func (p *Packs) DeletePack(c *fiber.Ctx) error {
	amount, err := c.ParamsInt("amount")
	if err != nil || amount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid amount"})
	}

	err = p.storage.DeletePack(amount)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
