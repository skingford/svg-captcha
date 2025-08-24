package captcha

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// SVGElement represents the root SVG element
type SVGElement struct {
	XMLName    xml.Name         `xml:"svg"`
	Width      int              `xml:"width,attr"`
	Height     int              `xml:"height,attr"`
	ViewBox    string           `xml:"viewBox,attr"`
	Xmlns      string           `xml:"xmlns,attr"`
	Background *RectElement     `xml:"rect,omitempty"`
	Texts      []*TextElement   `xml:"text,omitempty"`
	Paths      []*PathElement   `xml:"path,omitempty"`
	Lines      []*LineElement   `xml:"line,omitempty"`
	Circles    []*CircleElement `xml:"circle,omitempty"`
}

// TextElement represents an SVG text element
type TextElement struct {
	XMLName    xml.Name `xml:"text"`
	X          float64  `xml:"x,attr"`
	Y          float64  `xml:"y,attr"`
	Fill       string   `xml:"fill,attr"`
	FontSize   int      `xml:"font-size,attr"`
	FontFamily string   `xml:"font-family,attr"`
	Transform  string   `xml:"transform,attr,omitempty"`
	Content    string   `xml:",chardata"`
}

// RectElement represents an SVG rectangle
type RectElement struct {
	XMLName xml.Name `xml:"rect"`
	X       int      `xml:"x,attr"`
	Y       int      `xml:"y,attr"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	Fill    string   `xml:"fill,attr"`
}

// PathElement represents an SVG path (for text rendering)
type PathElement struct {
	XMLName     xml.Name `xml:"path"`
	D           string   `xml:"d,attr"`
	Fill        string   `xml:"fill,attr"`
	Stroke      string   `xml:"stroke,attr,omitempty"`
	StrokeWidth string   `xml:"stroke-width,attr,omitempty"`
}

// LineElement represents an SVG line (for noise)
type LineElement struct {
	XMLName xml.Name `xml:"line"`
	X1      float64  `xml:"x1,attr"`
	Y1      float64  `xml:"y1,attr"`
	X2      float64  `xml:"x2,attr"`
	Y2      float64  `xml:"y2,attr"`
	Stroke  string   `xml:"stroke,attr"`
	Width   float64  `xml:"stroke-width,attr"`
}

// CircleElement represents an SVG circle (for noise)
type CircleElement struct {
	XMLName xml.Name `xml:"circle"`
	CX      float64  `xml:"cx,attr"`
	CY      float64  `xml:"cy,attr"`
	R       float64  `xml:"r,attr"`
	Fill    string   `xml:"fill,attr"`
}

// SVGRenderer handles the generation of SVG content
type SVGRenderer struct {
	width    int
	height   int
	fontSize int
	colorMgr *ColorManager
}

// NewSVGRenderer creates a new SVG renderer
func NewSVGRenderer(config *Config) *SVGRenderer {
	return &SVGRenderer{
		width:    config.Width,
		height:   config.Height,
		fontSize: config.FontSize,
		colorMgr: NewColorManager(config),
	}
}

// RenderMathExpression converts a math expression into SVG format
func (sr *SVGRenderer) RenderMathExpression(expr *MathExpression, config *Config) (string, error) {
	// Create SVG container
	svg := sr.createSVGContainer(config)

	// Generate text paths for the expression
	questionText := strings.Replace(expr.Question, " = ?", " = ", 1)
	err := sr.addTextToSVG(svg, questionText, config)
	if err != nil {
		return "", NewError(ErrSVGGeneration, "failed to add text to SVG: "+err.Error(), 500)
	}

	// Add noise elements
	sr.addNoiseToSVG(svg, config)

	// Convert to XML
	xmlData, err := xml.MarshalIndent(svg, "", "  ")
	if err != nil {
		return "", NewError(ErrSVGGeneration, "failed to marshal SVG to XML: "+err.Error(), 500)
	}

	return xml.Header + string(xmlData), nil
}

// createSVGContainer creates the base SVG element with background
func (sr *SVGRenderer) createSVGContainer(config *Config) *SVGElement {
	svg := &SVGElement{
		Width:   sr.width,
		Height:  sr.height,
		ViewBox: fmt.Sprintf("0 0 %d %d", sr.width, sr.height),
		Xmlns:   "http://www.w3.org/2000/svg",
		Background: &RectElement{
			X:      0,
			Y:      0,
			Width:  sr.width,
			Height: sr.height,
			Fill:   config.Background,
		},
	}
	return svg
}

// addTextToSVG adds the math expression text as SVG text elements
func (sr *SVGRenderer) addTextToSVG(svg *SVGElement, text string, config *Config) error {
	// Calculate positioning
	textLen := len(text)
	charWidth := float64(sr.fontSize) * 0.6 // Approximate character width
	totalWidth := float64(textLen) * charWidth
	startX := (float64(sr.width) - totalWidth) / 2
	baseY := float64(sr.height)/2 + float64(sr.fontSize)/3 // Adjust for text baseline

	// Add some randomness to positioning
	yOffset, err := secureRandomFloat(-5, 5)
	if err != nil {
		yOffset = 0
	}

	// Render each character as a text element
	for i, char := range text {
		if char == ' ' {
			continue // Skip spaces
		}

		charX := startX + float64(i)*charWidth
		charY := baseY + yOffset

		// Add small random offset for each character
		xJitter, _ := secureRandomFloat(-3, 3)
		yJitter, _ := secureRandomFloat(-3, 3)
		charX += xJitter
		charY += yJitter

		// Add random rotation
		rotation, _ := secureRandomFloat(-15, 15)
		transform := fmt.Sprintf("rotate(%.1f %.2f %.2f)", rotation, charX, charY)

		// Create text element
		textElement := &TextElement{
			X:          charX,
			Y:          charY,
			Fill:       sr.colorMgr.GetRandomTextColor(),
			FontSize:   sr.fontSize,
			FontFamily: "Arial, sans-serif",
			Transform:  transform,
			Content:    string(char),
		}

		svg.Texts = append(svg.Texts, textElement)
	}

	return nil
}

// Character path generators are no longer needed since we use SVG text elements
// These functions are kept for backward compatibility but not used
func (sr *SVGRenderer) generateCharPath(char string, x, y float64) string {
	// Simplified fallback - just return empty path since we use text elements now
	return ""
}

// addNoiseToSVG adds visual noise to make OCR more difficult
func (sr *SVGRenderer) addNoiseToSVG(svg *SVGElement, config *Config) {
	if config.Noise <= 0 {
		return
	}

	noiseGen := NewNoiseGenerator()

	// Add random lines
	lines := noiseGen.GenerateLines(config.Noise*2, sr.width, sr.height, sr.colorMgr)
	svg.Lines = append(svg.Lines, lines...)

	// Add random dots
	circles := noiseGen.GenerateDots(config.Noise*3, sr.width, sr.height, sr.colorMgr)
	svg.Circles = append(svg.Circles, circles...)
}

// secureRandomFloat generates a secure random float between min and max
func secureRandomFloat(min, max float64) (float64, error) {
	if min >= max {
		return min, nil
	}

	range_ := max - min
	randInt, err := secureRandomInt(10000)
	if err != nil {
		return min, err
	}

	return min + (float64(randInt)/10000.0)*range_, nil
}
