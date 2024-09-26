package main

import (
	"log"
	"os"
	"time"

	"github.com/biswaz/img-maker/ordersummary"
)

func main() {
	log.Println("Starting the generator")

	order := ordersummary.OrderSummary{
		Items: []ordersummary.Item{
			{Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			{Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
			// {Name: "Phalaenopsis Amabilis 'Moth Orchid' - Large White Blooms, Ceramic Pot, 2-3 Flower Spikes", Quantity: 2, Price: 349.99},
		},
		Subtotal: 2357.97,
		Discount: 100.00,
		Shipping: 50.00,
		Taxes:    235.80,
		Total:    2643.77,
		Currency: "INR",
	}

	order.Items = append(order.Items, Items...)

	layout := ordersummary.Layout{
		Width:          700,
		Margin:         20,
		HeaderHeight:   80,
		ItemSpacing:    8,
		SectionSpacing: 20,
		FontSizes: ordersummary.FontSizes{
			Header:    24,
			Subheader: 18,
			Item:      12,
			Total:     14,
		},
	}

	start := time.Now()

	outputFile, err := os.OpenFile("order_summary.png", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error creating/opening output file: %v", err)
	}
	defer outputFile.Close()

	textContent := ordersummary.TextContent{
		HeaderText:   "Order Summary",
		ItemsText:    "Items",
		SubtotalText: "Subtotal:",
		ShippingText: "Shipping:",
		TaxesText:    "Taxes:",
		TotalText:    "Total:",
		DiscountText: "Discount:",
	}

	err = ordersummary.GenerateOrderSummary(order, outputFile, layout, textContent)
	if err != nil {
		log.Fatalf("Failed to generate order summary: %v", err)
	}

	log.Printf("Time taken: %v", time.Since(start))
}
