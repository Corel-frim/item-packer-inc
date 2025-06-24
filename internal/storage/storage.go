package storage

import (
	"errors"
	"sort"
	"sync"

	"github.com/corel-frim/item-packer-inc/internal/models"
)

var (
	ErrPackNotFound     = errors.New("pack not found")
	ErrNoPacksAvailable = errors.New("no packs available")
	ErrPackExists       = errors.New("pack with this amount already exists")
	ErrSoftLimitReached = errors.New("soft limit reached, cannot add more packs")
	SoftLimit           = 20 // Soft limit for arrays. Just for demonstration purposes
)

// PackStorage provides an in-memory storage for packs
type PackStorage struct {
	packs  []*models.Pack
	orders []models.Order
	mu     sync.RWMutex
}

// NewPackStorage creates a new instance of PackStorage
func NewPackStorage() *PackStorage {
	return &PackStorage{
		packs:  make([]*models.Pack, 0),
		orders: make([]models.Order, 0),
	}
}

// GetPacks returns all available packs
func (s *PackStorage) GetPacks() []*models.Pack {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getPacks()
}

// AddPack adds a new pack with the specified amount
func (s *PackStorage) AddPack(amount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If amount already exists - do nothing
	for _, p := range s.packs {
		if p.Amount == amount {
			return nil
		}
	}

	if len(s.packs) >= SoftLimit {
		return ErrSoftLimitReached
	}

	s.packs = append(s.packs, &models.Pack{Amount: amount})

	s.resortPacks()

	return nil
}

// UpdatePack updates a pack's amount
func (s *PackStorage) UpdatePack(oldAmount, newAmount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if new amount already exists
	for _, p := range s.packs {
		if p.Amount == newAmount {
			return ErrPackExists
		}
	}

	// Find and update the pack
	for _, p := range s.packs {
		if p.Amount == oldAmount {
			p.Amount = newAmount

			s.resortPacks()

			return nil
		}
	}

	return ErrPackNotFound
}

// DeletePack removes a pack with the specified amount
func (s *PackStorage) DeletePack(amount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.packs {
		if p.Amount == amount {
			// Remove the pack
			s.packs = append(s.packs[:i], s.packs[i+1:]...)
			return nil
		}
	}
	return ErrPackNotFound
}

func (s *PackStorage) GetOrders() []models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getOrders()
}

// CalculateOrder calculates the optimal packing for the requested items
func (s *PackStorage) CalculateOrder(requestedItems int) (models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.packs) == 0 {
		return models.Order{}, ErrNoPacksAvailable
	}

	s.resortPacks()
	packs := s.getPacks()

	order := &models.Order{
		RequestedItems: requestedItems,
		TotalItems:     0,
		Packs:          make([]models.OrderPack, 0),
	}

	// Use a greedy algorithm to find the optimal packing
	// First try to use the largest packs possible
	remainingItems, order := useFullPacks(packs, order)
	// If we still have remaining items, use the smallest pack
	order = s.addPackForRemainingItems(remainingItems, packs, order)

	order.OverpackedItems = order.TotalItems - requestedItems

	s.mergePacks(packs, order)

	// If we've reached the soft limit, keep only the most recent orders
	if len(s.orders) >= SoftLimit {
		// Keep only the most recent (SoftLimit - 1) orders to make room for the new one
		s.orders = s.orders[len(s.orders)-(SoftLimit-1):]
	}
	// Add the new order to the end of the slice
	s.orders = append(s.orders, *order)

	return *order, nil
}

func (s *PackStorage) addPackForRemainingItems(remainingItems int, packs []*models.Pack, order *models.Order) *models.Order {
	if remainingItems <= 0 {
		return order
	}
	smallestPack := packs[len(packs)-1]
	order.Packs = append(order.Packs, models.OrderPack{
		Quantity: 1,
		Pack:     smallestPack,
	})
	order.TotalItems += smallestPack.Amount

	return order
}

func (s *PackStorage) mergePacks(packs []*models.Pack, order *models.Order) {
	// Create ascending sorted pack sizes for merging
	availablePacks := getSortedPackSizes(packs)

	// Try merging multiple times to handle chain merges (e.g., 250+250=500, then 500+500=1000)
	for range availablePacks {
		if merged := tryMergeSameSizePacks(availablePacks, order); merged {
			continue
		}
		tryMergeDifferentSizePacks(availablePacks, order)
	}
}

func getSortedPackSizes(packs []*models.Pack) []*models.Pack {
	sorted := make([]*models.Pack, len(packs))
	copy(sorted, packs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Amount < sorted[j].Amount
	})
	return sorted
}

func tryMergeSameSizePacks(availablePacks []*models.Pack, order *models.Order) bool {
	// Group packs by size
	sizeGroups := make(map[int]int)
	for _, op := range order.Packs {
		sizeGroups[op.Pack.Amount] += op.Quantity
	}

	// Try to merge each group into larger packs
	for _, targetPack := range availablePacks {
		for size, count := range sizeGroups {
			if targetPack.Amount <= size {
				continue
			}

			if targetPack.Amount%size == 0 && count >= targetPack.Amount/size {
				mergePack(order, size, targetPack.Amount, targetPack.Amount/size)
				return true
			}
		}
	}
	return false
}

func tryMergeDifferentSizePacks(availablePacks []*models.Pack, order *models.Order) {
	for _, targetPack := range availablePacks {
		for _, orderPack := range order.Packs {
			smallSize := orderPack.Pack.Amount
			if targetPack.Amount <= smallSize {
				continue
			}

			if targetPack.Amount%smallSize == 0 && orderPack.Quantity >= targetPack.Amount/smallSize {
				mergePack(order, smallSize, targetPack.Amount, targetPack.Amount/smallSize)
				return
			}
		}
	}
}

func mergePack(order *models.Order, fromSize, toSize int, quantity int) {
	// Remove smaller packs
	newPacks := make([]models.OrderPack, 0, len(order.Packs))
	remainingToRemove := quantity

	for _, p := range order.Packs {
		if p.Pack.Amount == fromSize {
			if p.Quantity > remainingToRemove {
				p.Quantity -= remainingToRemove
				newPacks = append(newPacks, p)
			}
			remainingToRemove -= min(remainingToRemove, p.Quantity)
		} else {
			newPacks = append(newPacks, p)
		}
	}

	// Add or update larger pack
	found := false
	for i := range newPacks {
		if newPacks[i].Pack.Amount == toSize {
			newPacks[i].Quantity++
			found = true
			break
		}
	}

	if !found {
		newPacks = append(newPacks, models.OrderPack{
			Quantity: 1,
			Pack:     &models.Pack{Amount: toSize},
		})
	}

	order.Packs = newPacks
}

// min was added for readability, don't want to deal with math.Min for ints w/o a generics version
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// useFullPacks tries to use full packs for the requested items, but can leave some items unfulfilled if no pack fits exactly
func useFullPacks(packs []*models.Pack, order *models.Order) (int, *models.Order) {
	remainingItems := order.RequestedItems

	for _, pack := range packs {
		if pack.Amount <= remainingItems {
			quantity := remainingItems / pack.Amount
			if quantity > 0 {
				order.Packs = append(order.Packs, models.OrderPack{
					Quantity: quantity,
					Pack:     pack,
				})
				order.TotalItems += quantity * pack.Amount
				remainingItems -= quantity * pack.Amount
			}
		}
	}
	return remainingItems, order
}

// resortPacks sorts the packs in descending order by amount
func (s *PackStorage) resortPacks() {
	sort.Slice(s.packs, func(i, j int) bool {
		return s.packs[i].Amount > s.packs[j].Amount
	})
}

func (s *PackStorage) getPacks() []*models.Pack {
	// Return a deep copy to prevent external modifications. Delete copying if moved to external db
	result := make([]*models.Pack, len(s.packs))
	for i, pack := range s.packs {
		// Create a new Pack with the same amount
		result[i] = &models.Pack{Amount: pack.Amount}
	}

	return result
}

func (s *PackStorage) getOrders() []models.Order {
	// Return a copy to prevent external modifications. Delete copying if moved to external db
	result := make([]models.Order, len(s.orders))
	copy(result, s.orders)

	return result
}
