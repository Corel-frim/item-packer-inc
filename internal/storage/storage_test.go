package storage

import (
	"testing"

	"github.com/corel-frim/item-packer-inc/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewPackStorage(t *testing.T) {
	storage := NewPackStorage()
	assert.NotNil(t, storage)
	assert.Empty(t, storage.packs)
	assert.Empty(t, storage.orders)
}

func TestGetPacks(t *testing.T) {
	storage := NewPackStorage()

	// Test with empty storage
	packs := storage.GetPacks()
	assert.Empty(t, packs)

	// Add some packs and test again
	_ = storage.AddPack(100)
	_ = storage.AddPack(200)

	packs = storage.GetPacks()
	assert.Len(t, packs, 2)
	assert.Equal(t, 200, packs[0].Amount) // Packs should be sorted in descending order
	assert.Equal(t, 100, packs[1].Amount)
}

func TestAddPack(t *testing.T) {
	storage := NewPackStorage()

	// Test normal case
	err := storage.AddPack(100)
	assert.NoError(t, err)
	assert.Len(t, storage.packs, 1)
	assert.Equal(t, 100, storage.packs[0].Amount)

	// Test adding duplicate pack
	err = storage.AddPack(100)
	assert.NoError(t, err)
	assert.Len(t, storage.packs, 1) // Should still have only one pack

	// Test soft limit
	// Temporarily reduce the soft limit for testing
	originalLimit := SoftLimit
	SoftLimit = 2
	defer func() { SoftLimit = originalLimit }() // Restore original limit after test

	err = storage.AddPack(200)
	assert.NoError(t, err)
	assert.Len(t, storage.packs, 2)

	// Adding one more should hit the soft limit
	err = storage.AddPack(300)
	assert.Equal(t, ErrSoftLimitReached, err)
	assert.Len(t, storage.packs, 2) // Should still have only two packs
}

func TestUpdatePack(t *testing.T) {
	storage := NewPackStorage()

	// Add a pack
	_ = storage.AddPack(100)

	// Test normal case
	err := storage.UpdatePack(100, 150)
	assert.NoError(t, err)
	assert.Len(t, storage.packs, 1)
	assert.Equal(t, 150, storage.packs[0].Amount)

	// Test updating non-existent pack
	err = storage.UpdatePack(200, 250)
	assert.Equal(t, ErrPackNotFound, err)

	// Test updating to an amount that already exists
	_ = storage.AddPack(200)
	err = storage.UpdatePack(150, 200)
	assert.Equal(t, ErrPackExists, err)
}

func TestDeletePack(t *testing.T) {
	storage := NewPackStorage()

	_ = storage.AddPack(100)
	_ = storage.AddPack(200)

	err := storage.DeletePack(100)
	assert.NoError(t, err)
	assert.Len(t, storage.packs, 1)
	assert.Equal(t, 200, storage.packs[0].Amount)
}

func TestGetOrders(t *testing.T) {
	storage := NewPackStorage()

	// Test with empty storage
	orders := storage.GetOrders()
	assert.Empty(t, orders)

	// Add a pack and create an order
	_ = storage.AddPack(100)
	_, err := storage.CalculateOrder(100)
	assert.NoError(t, err)

	// Test getting orders
	orders = storage.GetOrders()
	assert.Len(t, orders, 1)
	assert.Equal(t, 100, orders[0].RequestedItems)
}

func TestCalculateOrder(t *testing.T) {
	storage := NewPackStorage()

	// Test with no packs available
	_, err := storage.CalculateOrder(100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no packs available")

	// Add some packs
	_ = storage.AddPack(250)
	_ = storage.AddPack(500)
	_ = storage.AddPack(1000)
	_ = storage.AddPack(2000)
	_ = storage.AddPack(5000)

	// Test exact match
	order, err := storage.CalculateOrder(500)
	assert.NoError(t, err)
	assert.Equal(t, 500, order.RequestedItems)
	assert.Equal(t, 500, order.TotalItems)
	assert.Equal(t, 0, order.OverpackedItems)
	assert.Len(t, order.Packs, 1)
	assert.Equal(t, 500, order.Packs[0].Pack.Amount)
	assert.Equal(t, 1, order.Packs[0].Quantity)

	// Test using multiple packs
	order, err = storage.CalculateOrder(1750)
	assert.NoError(t, err)
	assert.Equal(t, 1750, order.RequestedItems)
	assert.Equal(t, 1750, order.TotalItems)
	assert.Equal(t, 0, order.OverpackedItems)

	// Test with overpacking
	order, err = storage.CalculateOrder(1001)
	assert.NoError(t, err)
	assert.Equal(t, 1001, order.RequestedItems)
	assert.Equal(t, 1250, order.TotalItems)
	assert.Equal(t, 249, order.OverpackedItems)

	// Test soft limit for orders
	originalLimit := SoftLimit
	SoftLimit = 2
	defer func() { SoftLimit = originalLimit }() // Restore original limit after test

	// Create more orders to hit the soft limit
	_, _ = storage.CalculateOrder(100)
	_, _ = storage.CalculateOrder(200)
	_, _ = storage.CalculateOrder(300)

	// Should only keep the latest orders
	orders := storage.GetOrders()
	assert.Len(t, orders, 2)
	assert.Equal(t, 200, orders[0].RequestedItems)
	assert.Equal(t, 300, orders[1].RequestedItems)
}

func TestAddPackForRemainingItems(t *testing.T) {
	storage := NewPackStorage()

	// Add some packs
	_ = storage.AddPack(250)
	_ = storage.AddPack(500)

	packs := storage.GetPacks()
	order := &models.Order{
		RequestedItems: 600,
		TotalItems:     500,
		Packs: []models.OrderPack{
			{
				Quantity: 1,
				Pack:     &models.Pack{Amount: 500},
			},
		},
	}

	// Test adding a pack for remaining items
	result := storage.addPackForRemainingItems(100, packs, order)
	assert.Equal(t, 750, result.TotalItems)
	assert.Len(t, result.Packs, 2)
	assert.Equal(t, 250, result.Packs[1].Pack.Amount)
	assert.Equal(t, 1, result.Packs[1].Quantity)

	// Test with no remaining items
	order = &models.Order{
		RequestedItems: 500,
		TotalItems:     500,
		Packs: []models.OrderPack{
			{
				Quantity: 1,
				Pack:     &models.Pack{Amount: 500},
			},
		},
	}

	result = storage.addPackForRemainingItems(0, packs, order)
	assert.Equal(t, 500, result.TotalItems)
	assert.Len(t, result.Packs, 1)
}

func TestUseFullPacks(t *testing.T) {
	packs := []*models.Pack{
		{Amount: 5000},
		{Amount: 2000},
		{Amount: 1000},
		{Amount: 500},
		{Amount: 250},
	}

	order := &models.Order{
		RequestedItems: 7750,
	}

	// Test using full packs
	remaining, result := useFullPacks(packs, order)
	assert.Equal(t, 0, remaining)
	assert.Equal(t, 7750, result.TotalItems)
	assert.Len(t, result.Packs, 4)

	// Verify the packs used
	packCounts := make(map[int]int)
	for _, p := range result.Packs {
		packCounts[p.Pack.Amount] = p.Quantity
	}

	assert.Equal(t, 1, packCounts[5000])
	assert.Equal(t, 1, packCounts[2000])
	assert.Equal(t, 0, packCounts[1000]) // Not used
	assert.Equal(t, 1, packCounts[500])
	assert.Equal(t, 1, packCounts[250])

	// Test with remaining items
	order = &models.Order{
		RequestedItems: 7760,
	}

	remaining, result = useFullPacks(packs, order)
	assert.Equal(t, 10, remaining)
	assert.Equal(t, 7750, result.TotalItems)
}

func TestResortPacks(t *testing.T) {
	storage := NewPackStorage()

	// Add packs in random order
	_ = storage.AddPack(250)
	_ = storage.AddPack(1000)
	_ = storage.AddPack(500)

	// Verify they are sorted in descending order
	assert.Equal(t, 1000, storage.packs[0].Amount)
	assert.Equal(t, 500, storage.packs[1].Amount)
	assert.Equal(t, 250, storage.packs[2].Amount)

	// Add another pack and verify sorting is maintained
	_ = storage.AddPack(2000)
	assert.Equal(t, 2000, storage.packs[0].Amount)
	assert.Equal(t, 1000, storage.packs[1].Amount)
	assert.Equal(t, 500, storage.packs[2].Amount)
	assert.Equal(t, 250, storage.packs[3].Amount)
}

func TestMergePacks(t *testing.T) {
	storage := NewPackStorage()

	// Add packs of different sizes
	_ = storage.AddPack(250)
	_ = storage.AddPack(500)
	_ = storage.AddPack(1000)

	packs := storage.GetPacks()

	// Create an order with multiple small packs
	order := &models.Order{
		RequestedItems: 1000,
		TotalItems:     1000,
		Packs: []models.OrderPack{
			{
				Quantity: 2,
				Pack:     &models.Pack{Amount: 250},
			},
			{
				Quantity: 1,
				Pack:     &models.Pack{Amount: 500},
			},
		},
	}

	// Test merging packs
	storage.mergePacks(packs, order)

	// Verify that 2x250 packs were merged into 1x500 pack
	// and 1x500 + 1x500 were merged into 1x1000 pack
	assert.Len(t, order.Packs, 1)
	assert.Equal(t, 1000, order.Packs[0].Pack.Amount)
	assert.Equal(t, 1, order.Packs[0].Quantity)

	// Test with packs that can't be merged
	order = &models.Order{
		RequestedItems: 750,
		TotalItems:     750,
		Packs: []models.OrderPack{
			{
				Quantity: 1,
				Pack:     &models.Pack{Amount: 500},
			},
			{
				Quantity: 1,
				Pack:     &models.Pack{Amount: 250},
			},
		},
	}

	storage.mergePacks(packs, order)

	// Verify that packs remain unchanged (can't merge 500+250 into any available pack)
	assert.Len(t, order.Packs, 2)

	// Test with single quantity packs (should not be merged)
	order = &models.Order{
		RequestedItems: 500,
		TotalItems:     500,
		Packs: []models.OrderPack{
			{
				Quantity: 1,
				Pack:     &models.Pack{Amount: 250},
			},
			{
				Quantity: 1,
				Pack:     &models.Pack{Amount: 250},
			},
		},
	}

	storage.mergePacks(packs, order)

	// Verify that packs were merged (2x250 into 1x500)
	assert.Len(t, order.Packs, 1)
	assert.Equal(t, 500, order.Packs[0].Pack.Amount)
	assert.Equal(t, 1, order.Packs[0].Quantity)
}

func TestGetPacksReturnsCopy(t *testing.T) {
	storage := NewPackStorage()

	// Add some packs
	_ = storage.AddPack(100)
	_ = storage.AddPack(200)

	// Get packs and modify the returned slice
	packs := storage.GetPacks()
	packs[0].Amount = 300

	// Get packs again and verify the original values are unchanged
	packsAgain := storage.GetPacks()
	assert.Equal(t, 200, packsAgain[0].Amount)
	assert.Equal(t, 100, packsAgain[1].Amount)
}

func TestGetOrdersReturnsCopy(t *testing.T) {
	storage := NewPackStorage()

	// Add a pack and create an order
	_ = storage.AddPack(100)
	_, err := storage.CalculateOrder(100)
	assert.NoError(t, err)

	// Get orders and modify the returned slice
	orders := storage.GetOrders()
	orders[0].RequestedItems = 999

	// Get orders again and verify the original values are unchanged
	ordersAgain := storage.GetOrders()
	assert.Equal(t, 100, ordersAgain[0].RequestedItems)
}
