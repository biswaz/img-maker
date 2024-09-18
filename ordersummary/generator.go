package ordersummary

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

// OrderSummary represents the structure of an order summary
type OrderSummary struct {
	Items    []Item
	Subtotal float64
	Shipping float64
	Taxes    float64
	Total    float64
	Currency string
}

// Item represents a single item in the order
type Item struct {
	Name     string
	Quantity int
	Price    float64
}

// Layout defines the layout parameters for the order summary image
type Layout struct {
	Width          int
	Margin         int
	HeaderHeight   int
	ItemSpacing    int
	SectionSpacing int
	FontSizes      map[string]float64
}

// TextContent defines the text content used in the order summary
type TextContent struct {
    HeaderText     string
    ItemsText      string
    SubtotalText   string
    ShippingText   string
    TaxesText      string
    TotalText      string
}

// GenerateOrderSummary creates an image of the order summary and writes it to the provided file
func GenerateOrderSummary(order OrderSummary, outputFile *os.File, layout Layout, textContent TextContent) error {
	// Load fonts
	headerFont := loadFont(gobold.TTF, layout.FontSizes["header"])
	itemFont := loadFont(goregular.TTF, layout.FontSizes["item"])
	subheaderFont := loadFont(gobold.TTF, layout.FontSizes["subheader"])
	totalFont := loadFont(gobold.TTF, layout.FontSizes["total"])

	// Initialize current Y position
	y := layout.Margin

	// Calculate content width
	contentWidth := layout.Width - layout.Margin*4
	priceColumnWidth := 100
	itemColumnWidth := contentWidth - priceColumnWidth

	// Header
	headerHeight := int(math.Ceil(getTextHeight(headerFont)))
	y += headerHeight + layout.SectionSpacing

	// Horizontal line after header
	y += 1 + layout.SectionSpacing

	// Items subheader
	subheaderHeight := int(math.Ceil(getTextHeight(subheaderFont)))
	y += subheaderHeight + layout.ItemSpacing

	// Items (with wrapped text)
	itemHeight := int(math.Ceil(getTextHeight(itemFont)))
	for _, item := range order.Items {
		itemText := formatItem(item)
		wrappedText := wrapText(itemText, itemColumnWidth, layout.FontSizes["item"])
		y += len(wrappedText) * itemHeight + (len(wrappedText) - 1) * layout.ItemSpacing
		y += layout.ItemSpacing // Extra spacing between items
	}

	// Space before totals
	y += layout.SectionSpacing

	// Horizontal line before totals
	y += 1 + layout.SectionSpacing

	// Total sections (subtotal, shipping, taxes, total)
	totalSectionHeight := int(math.Ceil(getTextHeight(totalFont)))
	y += 5 * totalSectionHeight + 5 * layout.SectionSpacing

	// Bottom margin
	y += layout.Margin

	// Create the image with the calculated height
	img := image.NewRGBA(image.Rect(0, 0, layout.Width, y))

	// Set background color
	bgColor := color.RGBA{245, 245, 245, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw main content area with rounded corners
	contentColor := color.RGBA{255, 255, 255, 255}
	drawRoundedRect(img, layout.Margin, layout.Margin, layout.Width-layout.Margin, y-layout.Margin, 10, contentColor)

	textColor := color.RGBA{60, 60, 60, 255}

	// Reset Y position for drawing
	y = layout.Margin + layout.HeaderHeight/2

	// Draw header
	drawCenteredText(img, textContent.HeaderText, layout.Width/2, y, layout.FontSizes["header"], textColor, true)
	drawHorizontalLine(img, layout.Margin*2, layout.Width-layout.Margin*2, y+layout.HeaderHeight/2, color.RGBA{220, 220, 220, 255})

	// Draw items
	y += layout.HeaderHeight/2 + layout.SectionSpacing
	drawLeftAlignedText(img, textContent.ItemsText, layout.Margin*2, y, layout.FontSizes["subheader"], textColor, true)
	y += int(layout.FontSizes["subheader"]) + layout.ItemSpacing

	for _, item := range order.Items {
		itemText := formatItem(item)
		wrappedText := wrapText(itemText, itemColumnWidth, layout.FontSizes["item"])
		
		for i, line := range wrappedText {
			drawLeftAlignedText(img, line, layout.Margin*3, y, layout.FontSizes["item"], textColor, false)
			
			if i == 0 {
				priceStr := fmt.Sprintf("%s %.2f", order.Currency, item.Price*float64(item.Quantity))
				drawRightAlignedText(img, priceStr, layout.Width-layout.Margin*2, y, layout.FontSizes["item"], textColor, false)
			}
			
			y += itemHeight
		}
		
		y += layout.ItemSpacing
	}

	y += layout.SectionSpacing
	drawHorizontalLine(img, layout.Margin*2, layout.Width-layout.Margin*2, y, color.RGBA{220, 220, 220, 255})
	y += layout.SectionSpacing

	// Draw totals
	drawTotalLine(img, textContent.SubtotalText, order.Subtotal, y, layout, textColor, order.Currency)
	y += totalSectionHeight + layout.SectionSpacing
	drawTotalLine(img, textContent.ShippingText, order.Shipping, y, layout, textColor, order.Currency)
	y += totalSectionHeight + layout.SectionSpacing
	drawTotalLine(img, textContent.TaxesText, order.Taxes, y, layout, textColor, order.Currency)
	y += totalSectionHeight + layout.SectionSpacing
	drawHorizontalLine(img, layout.Margin*2, layout.Width-layout.Margin*2, y, color.RGBA{220, 220, 220, 255})
	y += layout.SectionSpacing 
	drawTotalLine(img, textContent.TotalText, order.Total, y, layout, textColor, order.Currency)

	// Save the image
	return png.Encode(outputFile, img)
}

func formatItem(item Item) string {
	return fmt.Sprintf("%dx %s", item.Quantity, item.Name)
}

// Replace the truncateText function with this new wrapText function
func wrapText(text string, maxWidth int, fontSize float64) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else {
			testLine := currentLine + " " + word
			if measureTextWidth(testLine, fontSize) <= maxWidth {
				currentLine = testLine
			} else {
				lines = append(lines, currentLine)
				currentLine = word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func drawRoundedRect(img *image.RGBA, x1, y1, x2, y2, radius int, c color.Color) {
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			dx, dy := float64(x-x1), float64(y-y1)
			if dx < float64(radius) && dy < float64(radius) {
				if math.Pow(dx-float64(radius), 2)+math.Pow(dy-float64(radius), 2) > math.Pow(float64(radius), 2) {
					continue
				}
			}
			dx, dy = float64(x-x2), float64(y-y1)
			if dx > -float64(radius) && dy < float64(radius) {
				if math.Pow(dx+float64(radius), 2)+math.Pow(dy-float64(radius), 2) > math.Pow(float64(radius), 2) {
					continue
				}
			}
			dx, dy = float64(x-x1), float64(y-y2)
			if dx < float64(radius) && dy > -float64(radius) {
				if math.Pow(dx-float64(radius), 2)+math.Pow(dy+float64(radius), 2) > math.Pow(float64(radius), 2) {
					continue
				}
			}
			dx, dy = float64(x-x2), float64(y-y2)
			if dx > -float64(radius) && dy > -float64(radius) {
				if math.Pow(dx+float64(radius), 2)+math.Pow(dy+float64(radius), 2) > math.Pow(float64(radius), 2) {
					continue
				}
			}
			img.Set(x, y, c)
		}
	}
}

func drawHorizontalLine(img *image.RGBA, x1, x2, y int, c color.Color) {
	for x := x1; x <= x2; x++ {
		img.Set(x, y, c)
	}
}

func addLabel(img *image.RGBA, x, y int, label string, size float64, c color.Color) {
	fontBytes, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	f := truetype.NewFace(fontBytes, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
	})
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: f,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)},
	}
	d.DrawString(label)
}

func addLabelBold(img *image.RGBA, x, y int, label string, size float64, c color.Color) {
	fontBytes, err := truetype.Parse(gobold.TTF)
	if err != nil {
		panic(err)
	}
	f := truetype.NewFace(fontBytes, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
	})
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: f,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)},
	}
	d.DrawString(label)
}

func measureTextWidth(text string, size float64) int {
	fontBytes, _ := truetype.Parse(goregular.TTF)
	f := truetype.NewFace(fontBytes, &truetype.Options{Size: size})
	return font.MeasureString(f, text).Round()
}

func drawTotalLine(img *image.RGBA, label string, value float64, y int, layout Layout, textColor color.Color, currency string) {
	drawLeftAlignedText(img, label, layout.Margin*3, y, layout.FontSizes["item"], textColor, false)
	valueStr := fmt.Sprintf("%s %.2f", currency, value)
	drawRightAlignedText(img, valueStr, layout.Width-layout.Margin*2, y, layout.FontSizes["item"], textColor, false)
}

func drawCenteredText(img *image.RGBA, text string, x, y int, size float64, c color.Color, bold bool) {
	width := measureTextWidth(text, size)
	x -= width / 2
	if bold {
		addLabelBold(img, x, y, text, size, c)
	} else {
		addLabel(img, x, y, text, size, c)
	}
}

func drawLeftAlignedText(img *image.RGBA, text string, x, y int, size float64, c color.Color, bold bool) {
	if bold {
		addLabelBold(img, x, y, text, size, c)
	} else {
		addLabel(img, x, y, text, size, c)
	}
}

func drawRightAlignedText(img *image.RGBA, text string, x, y int, size float64, c color.Color, bold bool) {
	width := measureTextWidth(text, size)
	x -= width
	if bold {
		addLabelBold(img, x, y, text, size, c)
	} else {
		addLabel(img, x, y, text, size, c)
	}
}

func loadFont(fontData []byte, size float64) font.Face {
	f, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatalf("Failed to parse font: %v", err)
	}
	return truetype.NewFace(f, &truetype.Options{Size: size})
}

func getTextHeight(face font.Face) float64 {
	metrics := face.Metrics()
	return float64(metrics.Height) / 64
}