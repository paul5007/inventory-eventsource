package inventory

import (
	"context"
	"fmt"
	es "github.com/altairsix/eventsource"
	"time"
)

type Product struct {
	ProductId   string
	Version     int
	UpdatedAt   time.Time
	SupplierId  string
	Description string  // short description
	Quantity    int     // number of units
	BuyPrice    float64 // price bought from supplier
	SellPrice   float64 // price sold to customer
}

func (p *Product) On(event es.Event) error {
	switch e := event.(type) {
	case *ProductCreated:
		p.SupplierId = e.SupplierId
		p.Description = e.Description
		p.Quantity = e.Quantity
		p.BuyPrice = e.BuyPrice
		p.SellPrice = e.SellPrice
	case *ProductBought:
		p.Quantity += e.QuantityBought
	case *ProductSold:
		p.Quantity -= e.QuantitySold
	case *ProductSupplierChanged:
		p.SupplierId = e.SupplierId
	default:
		return fmt.Errorf("unhandled event, %v", e)
	}

	p.ProductId = event.AggregateID()
	p.Version = event.EventVersion()
	p.UpdatedAt = event.EventAt()

	return nil
}

func (p *Product) Apply(ctx context.Context, command es.Command) ([]es.Event, error) {
	switch c := command.(type) {
	case *CreateProduct:
		productCreated := &ProductCreated{
			Model:       es.Model{ID: c.AggregateID(), Version: p.Version + 1, At: time.Now()},
			SupplierId:  c.SupplierId,
			Description: c.Description,
			Quantity:    c.Quantity,
			BuyPrice:    c.BuyPrice,
			SellPrice:   c.SellPrice,
		}
		return []es.Event{productCreated}, nil

	case *BuyProduct:
		if c.BuyPrice > p.BuyPrice {
			return nil, fmt.Errorf("Can't buy product(%+v) at price(%+v), which is above max price(%+v)\n", c.AggregateID(), c.BuyPrice, p.BuyPrice)
		}
		productBought := &ProductBought{
			Model:          es.Model{ID: c.AggregateID(), Version: p.Version + 1, At: time.Now()},
			QuantityBought: c.QuantityBought,
			BuyPrice:       c.BuyPrice,
		}
		return []es.Event{productBought}, nil

	case *SellProduct:
		if p.Quantity-c.QuantitySold <= 0 {
			return nil, fmt.Errorf("Unable to sell quantity(%+v) of product(%+v)\n", c.QuantitySold, c.AggregateID())
		}
		if c.SellPrice < p.SellPrice {
			return nil, fmt.Errorf("Can't sell product(%+v) at price(%+v), which is below bottom line price(%+v)\n", c.AggregateID(), c.SellPrice, p.SellPrice)
		}
		productSold := &ProductSold{
			Model:        es.Model{ID: c.AggregateID(), Version: p.Version + 1, At: time.Now()},
			QuantitySold: c.QuantitySold,
			SellPrice:    c.SellPrice,
		}
		return []es.Event{productSold}, nil

	case *ChangeProductSupplier:
		productSupplierChanged := &ProductSupplierChanged{
			Model: es.Model{ID: c.AggregateID(), Version: p.Version + 1, At: time.Now()},
		}
		return []es.Event{productSupplierChanged}, nil

	default:
		return nil, fmt.Errorf("unhandled command, %v", c)
	}
}
