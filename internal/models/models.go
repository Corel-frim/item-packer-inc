package models

// Pack represents a package with a specific amount of items
type Pack struct {
	Amount int `json:"amount"`
}

// OrderPack represents a pack used in an order with its quantity
type OrderPack struct {
	Quantity int   `json:"quantity"`
	Pack     *Pack `json:"pack"`
}

// Order represents a customer order with requested items and packing details
type Order struct {
	RequestedItems  int         `json:"requestedItems"`
	OverpackedItems int         `json:"overpackedItems"`
	TotalItems      int         `json:"totalItems"`
	Packs           []OrderPack `json:"packs"`
}
