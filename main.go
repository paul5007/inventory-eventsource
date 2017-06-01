package main

import (
	"context"
	"fmt"
	es "github.com/altairsix/eventsource"
	i "github.com/pfa5007/inventory-eventsource/inventory"
	"log"
	"os"
)

func main() {
	fmt.Println("Welcome to the generic inventory system using event sourcing!")

	serializer := es.NewJSONSerializer(
		i.ProductCreated{},
		i.ProductBought{},
		i.ProductSold{},
		i.ProductSupplierChanged{},
	)
	repo := es.New(&i.Product{},
		es.WithSerializer(serializer),
		es.WithDebug(os.Stdout),
	)

	productId := "123"
	ctx := context.Background()

	err := repo.Dispatch(ctx,
		&i.CreateProduct{CommandModel: es.CommandModel{ID: productId},
			SupplierId:  "ABC",
			Description: "Very cool product. Must be soap",
			Quantity:    100,
			BuyPrice:    50.00,
			SellPrice:   60.00,
		},
	)
	if err != nil {
		log.Println(err)
	}

	err = repo.Dispatch(ctx,
		&i.SellProduct{CommandModel: es.CommandModel{ID: productId},
			QuantitySold: 20,
			SellPrice:    60.00,
		},
	)
	if err != nil {
		log.Println(err)
	}

	err = repo.Dispatch(ctx,
		&i.SellProduct{CommandModel: es.CommandModel{ID: productId},
			QuantitySold: 20,
			SellPrice:    50.00,
		},
	)
	if err != nil {
		log.Println(err)
	}

	err = repo.Dispatch(ctx,
		&i.BuyProduct{CommandModel: es.CommandModel{ID: productId},
			QuantityBought: 50,
			BuyPrice:       60.00,
		},
	)
	if err != nil {
		log.Println(err)
	}

	err = repo.Dispatch(ctx,
		&i.BuyProduct{CommandModel: es.CommandModel{ID: productId},
			QuantityBought: 50,
			BuyPrice:       40.00,
		},
	)
	if err != nil {
		log.Println(err)
	}

	aggregate, err := repo.Load(ctx, productId)
	found := aggregate.(*i.Product)
	fmt.Printf("%+v\n", found)
}
