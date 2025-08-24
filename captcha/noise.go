package captcha

import "fmt"

// NoiseGenerator generates visual noise elements for captchas
type NoiseGenerator struct{}

// NewNoiseGenerator creates a new noise generator
func NewNoiseGenerator() *NoiseGenerator {
	return &NoiseGenerator{}
}

// GenerateLines creates random curved lines for visual noise (now returns PathElements instead of LineElements)
func (ng *NoiseGenerator) GenerateLines(count, width, height int, colorMgr *ColorManager) []*PathElement {
	curves := make([]*PathElement, 0, count)

	for i := 0; i < count; i++ {
		// Generate start and end points
		startX, err1 := secureRandomFloat(0, float64(width))
		startY, err2 := secureRandomFloat(0, float64(height))
		endX, err3 := secureRandomFloat(0, float64(width))
		endY, err4 := secureRandomFloat(0, float64(height))

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			continue // skip this curve if random generation fails
		}

		// Random stroke width
		strokeWidth, err := secureRandomFloat(0.5, 2.0)
		if err != nil {
			strokeWidth = 1.0
		}

		// Generate curve path with random control points
		pathData := ng.generateCurvePath(startX, startY, endX, endY, float64(width), float64(height))

		curve := &PathElement{
			D:           pathData,
			Fill:        "none",
			Stroke:      colorMgr.GetRandomNoiseColor(),
			StrokeWidth: fmt.Sprintf("%.5g", strokeWidth),
		}

		curves = append(curves, curve)
	}

	return curves
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

// generateCurvePath creates a curved path between two points with random control points
func (ng *NoiseGenerator) generateCurvePath(startX, startY, endX, endY, width, height float64) string {
	// Choose curve type randomly
	curveType, _ := secureRandomInt(3)
	
	switch curveType {
	case 0:
		// Quadratic Bezier curve with single control point
		controlX := (startX + endX) / 2
		controlY := (startY + endY) / 2
		
		// Add random offset to control point
		offsetX, _ := secureRandomFloat(-width*0.3, width*0.3)
		offsetY, _ := secureRandomFloat(-height*0.3, height*0.3)
		controlX += offsetX
		controlY += offsetY
		
		return fmt.Sprintf("M%.2f,%.2f Q%.2f,%.2f %.2f,%.2f",
			startX, startY, controlX, controlY, endX, endY)
		
	case 1:
		// Cubic Bezier curve with two control points
		control1X := startX + (endX-startX)*0.33
		control1Y := startY + (endY-startY)*0.33
		control2X := startX + (endX-startX)*0.67
		control2Y := startY + (endY-startY)*0.67
		
		// Add random offsets
		offset1X, _ := secureRandomFloat(-width*0.2, width*0.2)
		offset1Y, _ := secureRandomFloat(-height*0.2, height*0.2)
		offset2X, _ := secureRandomFloat(-width*0.2, width*0.2)
		offset2Y, _ := secureRandomFloat(-height*0.2, height*0.2)
		
		control1X += offset1X
		control1Y += offset1Y
		control2X += offset2X
		control2Y += offset2Y
		
		return fmt.Sprintf("M%.2f,%.2f C%.2f,%.2f %.2f,%.2f %.2f,%.2f",
			startX, startY, control1X, control1Y, control2X, control2Y, endX, endY)
		
	default:
		// Sinusoidal curve using multiple quadratic segments
		numSegments := 3
		path := fmt.Sprintf("M%.2f,%.2f", startX, startY)
		
		for i := 1; i <= numSegments; i++ {
			t := float64(i) / float64(numSegments)
			segmentX := startX + (endX-startX)*t
			segmentY := startY + (endY-startY)*t
			
			// Add sinusoidal variation
			amplitude, _ := secureRandomFloat(10, 30)
			offset := amplitude * (0.5 - 0.5*float64(i%2)) // Alternating pattern
			
			// Perpendicular offset
			dx := endX - startX
			dy := endY - startY
			length := (dx*dx + dy*dy)
			if length > 0 {
				length = 1.0 / length
				perpX := -dy * length * offset
				perpY := dx * length * offset
				segmentX += perpX
				segmentY += perpY
			}
			
			if i == 1 {
				path += fmt.Sprintf(" Q%.2f,%.2f %.2f,%.2f",
					segmentX, segmentY, startX+(endX-startX)*0.5, startY+(endY-startY)*0.5)
			} else {
				path += fmt.Sprintf(" T%.2f,%.2f", segmentX, segmentY)
			}
		}
		
		return path
	}
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

		// Generate more sophisticated curve
		pathData := ng.generateCurvePath(startX, startY, endX, endY, float64(width), float64(height))

		// Random stroke width
		strokeWidth, err := secureRandomFloat(0.3, 1.5)
		if err != nil {
			strokeWidth = 0.8
		}

		arc := &PathElement{
			D:           pathData,
			Fill:        "none",
			Stroke:      colorMgr.GetRandomNoiseColor(),
			StrokeWidth: fmt.Sprintf("%.5g", strokeWidth),
		}

		arcs = append(arcs, arc)
	}

	return arcs
}
