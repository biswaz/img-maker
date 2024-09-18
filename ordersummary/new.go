package ordersummary

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

// GenerateOrderSummaryGG creates an image of the order summary using the gg package and writes it to the provided file
func GenerateOrderSummaryGG(order OrderSummary, outputFile *os.File, layout Layout) error {
	// Initialize the context
	dc := gg.NewContext(layout.Width, calculateHeight(order, layout))

	// Set background color
	dc.SetColor(color.RGBA{245, 245, 245, 255})
	dc.Clear()

	// Draw main content area with rounded corners
	dc.SetColor(color.White)
	dc.DrawRoundedRectangle(float64(layout.Margin), float64(layout.Margin),
		float64(layout.Width-2*layout.Margin), float64(dc.Height()-2*layout.Margin), 10)
	dc.Fill()

	// Load fonts
	headerFont := loadFontGG(gobold.TTF, layout.FontSizes["header"])
	itemFont := loadFontGG(goregular.TTF, layout.FontSizes["item"])
	subheaderFont := loadFontGG(gobold.TTF, layout.FontSizes["subheader"])
	totalFont := loadFontGG(gobold.TTF, layout.FontSizes["total"])

	dc.SetColor(color.Black)

	// Draw header
	y := float64(layout.Margin + layout.HeaderHeight/2)
	dc.SetFontFace(truetype.NewFace(headerFont, &truetype.Options{Size: layout.FontSizes["header"]}))
	dc.DrawStringAnchored("Order Summary", float64(layout.Width/2), y, 0.5, 0.5)

	// Remove or comment out any SetColor calls here that might be changing the color

	// Draw horizontal line
	y += float64(layout.HeaderHeight/2)
	drawHorizontalLineGG(dc, layout.Margin*2, layout.Width-layout.Margin*2, int(y))

	// Ensure color is set to black before drawing items
	dc.SetColor(color.Black)

	// Draw items
	y += float64(layout.SectionSpacing)
	dc.SetFontFace(truetype.NewFace(subheaderFont, &truetype.Options{Size: layout.FontSizes["subheader"]}))
	dc.DrawString("Items", float64(layout.Margin*2), y)

	y += dc.FontHeight() + float64(layout.ItemSpacing)
	dc.SetFontFace(truetype.NewFace(itemFont, &truetype.Options{Size: layout.FontSizes["item"]}))

	for _, item := range order.Items {
		itemText := formatItem(item)
		wrappedText := wrapTextGG(dc, itemText, float64(layout.Width-layout.Margin*5))
		
		for i, line := range wrappedText {
			dc.DrawString(line, float64(layout.Margin*3), y)
			
			if i == 0 {
				priceStr := fmt.Sprintf("%s %.2f", order.Currency, item.Price*float64(item.Quantity))
				dc.DrawStringAnchored(priceStr, float64(layout.Width-layout.Margin*2), y, 1, 0)
			}
			
			y += dc.FontHeight()
		}
		
		y += float64(layout.ItemSpacing)
	}

	// Draw totals
	y += float64(layout.SectionSpacing)
	drawHorizontalLineGG(dc, layout.Margin*2, layout.Width-layout.Margin*2, int(y))
	y += float64(layout.SectionSpacing)

	dc.SetFontFace(truetype.NewFace(totalFont, &truetype.Options{Size: layout.FontSizes["total"]}))
	drawTotalLineGG(dc, "Subtotal:", order.Subtotal, y, layout, order.Currency)
	y += dc.FontHeight() + float64(layout.SectionSpacing)
	drawTotalLineGG(dc, "Shipping:", order.Shipping, y, layout, order.Currency)
	y += dc.FontHeight() + float64(layout.SectionSpacing)
	drawTotalLineGG(dc, "Taxes:", order.Taxes, y, layout, order.Currency)
	y += dc.FontHeight() + float64(layout.SectionSpacing)
	drawHorizontalLineGG(dc, layout.Margin*2, layout.Width-layout.Margin*2, int(y))
	y += float64(layout.SectionSpacing)
	drawTotalLineGG(dc, "Total:", order.Total, y, layout, order.Currency)

	// Save the image
	return dc.EncodePNG(outputFile)
}

func wrapTextGG(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var line string
	for _, word := range words {
		testLine := line
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if width, _ := dc.MeasureString(testLine); width <= maxWidth {
			line = testLine
		} else {
			lines = append(lines, line)
			line = word
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

func drawHorizontalLineGG(dc *gg.Context, x1, x2, y int) {
	dc.SetColor(color.RGBA{220, 220, 220, 255})
	dc.DrawLine(float64(x1), float64(y), float64(x2), float64(y))
	dc.Stroke()
	setTextColor(dc)
}

func setTextColor(dc *gg.Context) {
	dc.SetColor(color.Black)
}

func drawTotalLineGG(dc *gg.Context, label string, value float64, y float64, layout Layout, currency string) {
	dc.DrawString(label, float64(layout.Margin*3), y)
	valueStr := fmt.Sprintf("%s %.2f", currency, value)
	dc.DrawStringAnchored(valueStr, float64(layout.Width-layout.Margin*2), y, 1, 0)
}

func loadFontGG(fontData []byte, size float64) *truetype.Font {
	f, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatalf("Failed to parse font: %v", err)
	}
	return f
}

func calculateHeight(order OrderSummary, layout Layout) int {
	height := layout.Margin * 2 // Top and bottom margins
	height += layout.HeaderHeight
	height += layout.SectionSpacing // Space after header

	// Items section
	height += int(layout.FontSizes["subheader"]) // "Items" subheader
	height += layout.ItemSpacing

	// Calculate height for each item
	dc := gg.NewContext(1, 1) // Temporary context for text measurements
	itemFont := loadFontGG(goregular.TTF, layout.FontSizes["item"])
	dc.SetFontFace(truetype.NewFace(itemFont, &truetype.Options{Size: layout.FontSizes["item"]}))

	for _, item := range order.Items {
		itemText := formatItem(item)
		wrappedText := wrapTextGG(dc, itemText, float64(layout.Width-layout.Margin*5))
		height += (int(dc.FontHeight()) * len(wrappedText)) + layout.ItemSpacing
	}

	// Totals section
	height += layout.SectionSpacing * 2 // Space before and after horizontal line
	height += int(layout.FontSizes["total"]) * 4 // Four lines: Subtotal, Shipping, Taxes, Total
	height += layout.SectionSpacing * 5 // Spacing between total lines

	return height
}

