package inventory

import (
	es "github.com/altairsix/eventsource"
)

// ProductCreated is an event in which a specific product was created
type ProductCreated struct {
	es.Model
	SupplierId  string
	Description string
	Quantity    int
	BuyPrice    float64
	SellPrice   float64
}

// ProductBought is an event in which a specific product was bought from supplier
type ProductBought struct {
	es.Model
	QuantityBought int
	BuyPrice       float64
}

// ProductSold is an event in which a specific product was sold to a customer
type ProductSold struct {
	es.Model
	QuantitySold int
	SellPrice    float64
}

type ProductSupplierChanged struct {
	es.Model
	SupplierId string
}
