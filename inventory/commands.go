package inventory

import (
	es "github.com/altairsix/eventsource"
)

// CreateProduct is a command used to create a product in inventory
type CreateProduct struct {
	es.CommandModel
	SupplierId  string
	Description string
	Quantity    int
	BuyPrice    float64
	SellPrice   float64
}

// RemoveProduct is a command used to remove a product in inventory
type RemoveProduct struct {
	es.CommandModel
}

// BuyProduct is a command used to buy a product from a supplier
type BuyProduct struct {
	es.CommandModel
	QuantityBought int
	BuyPrice       float64
}

// SellProduct is a command used to sell a product to a customer
type SellProduct struct {
	es.CommandModel
	QuantitySold int
	SellPrice    float64
}

// ChangeProductSupplier is a command used to change the supplier of the product
type ChangeProductSupplier struct {
	es.CommandModel
	SupplierId string
}
