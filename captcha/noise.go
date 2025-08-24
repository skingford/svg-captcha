package captcha

import "fmt"

// NoiseGenerator generates visual noise elements for captchas
type NoiseGenerator struct{}

// NewNoiseGenerator creates a new noise generator
func NewNoiseGenerator() *NoiseGenerator {
	return &NoiseGenerator{}
}

// GenerateLines creates random lines for visual noise
func (ng *NoiseGenerator) GenerateLines(count, width, height int, colorMgr *ColorManager) []*LineElement {
	lines := make([]*LineElement, 0, count)

	for i := 0; i < count; i++ {
		x1, err1 := secureRandomFloat(0, float64(width))
		y1, err2 := secureRandomFloat(0, float64(height))
		x2, err3 := secureRandomFloat(0, float64(width))
		y2, err4 := secureRandomFloat(0, float64(height))

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			continue // skip this line if random generation fails
		}

		// Random line width
		lineWidth, err := secureRandomFloat(0.5, 2.0)
		if err != nil {
			lineWidth = 1.0
		}

		line := &LineElement{
			X1:     x1,
			Y1:     y1,
			X2:     x2,
			Y2:     y2,
			Stroke: colorMgr.GetRandomNoiseColor(),
			Width:  lineWidth,
		}

		lines = append(lines, line)
	}

	return lines
}

// GenerateDots creates random circles for visual noise
func (ng *NoiseGenerator) GenerateDots(count, width, height int, colorMgr *ColorManager) []*CircleElement {
	circles := make([]*CircleElement, 0, count)

	for i := 0; i < count; i++ {
		cx, err1 := secureRandomFloat(0, float64(width))
		cy, err2 := secureRandomFloat(0, float64(height))

		if err1 != nil || err2 != nil {
			continue // skip this circle if random generation fails
		}

		// Random radius
		radius, err := secureRandomFloat(1.0, 4.0)
		if err != nil {
			radius = 2.0
		}

		circle := &CircleElement{
			CX:   cx,
			CY:   cy,
			R:    radius,
			Fill: colorMgr.GetRandomNoiseColor(),
		}

		circles = append(circles, circle)
	}

	return circles
}

// GenerateArcs creates random arc segments for more sophisticated noise
func (ng *NoiseGenerator) GenerateArcs(count, width, height int, colorMgr *ColorManager) []*PathElement {
	arcs := make([]*PathElement, 0, count)

	for i := 0; i < count; i++ {
		startX, err1 := secureRandomFloat(0, float64(width))
		startY, err2 := secureRandomFloat(0, float64(height))
		endX, err3 := secureRandomFloat(0, float64(width))
		endY, err4 := secureRandomFloat(0, float64(height))

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			continue
		}

		// Create a simple curved path
		midX := (startX + endX) / 2
		midY := (startY + endY) / 2

		// Add some curvature
		curveOffset, _ := secureRandomFloat(-20, 20)
		midY += curveOffset

		pathData := fmt.Sprintf("M%.2f,%.2f Q%.2f,%.2f %.2f,%.2f",
			startX, startY, midX, midY, endX, endY)

		arc := &PathElement{
			D:      pathData,
			Fill:   "none",
			Stroke: colorMgr.GetRandomNoiseColor(),
		}

		arcs = append(arcs, arc)
	}

	return arcs
}
