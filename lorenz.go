package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"
	"time"
)

// Point3D represents a 3D point
type Point3D struct {
	X, Y, Z float64
}

// LorenzSystem represents the Lorenz attractor system
type LorenzSystem struct {
	sigma, rho, beta float64
	dt               float64
	x, y, z          float64
	trail            []Point3D
	maxTrailLength   int
}

// NewLorenzSystem creates a new Lorenz system
func NewLorenzSystem(maxTrail int) *LorenzSystem {
	return &LorenzSystem{
		sigma:          10.0,
		rho:            28.0,
		beta:           8.0 / 3.0,
		dt:             0.01,
		x:              1.0,
		y:              1.0,
		z:              1.0,
		trail:          make([]Point3D, 0, maxTrail),
		maxTrailLength: maxTrail,
	}
}

// Step advances the system by one time step
func (l *LorenzSystem) Step() Point3D {
	// Lorenz equations
	dx := l.sigma * (l.y - l.x)
	dy := l.x*(l.rho-l.z) - l.y
	dz := l.x*l.y - l.beta*l.z

	l.x += dx * l.dt
	l.y += dy * l.dt
	l.z += dz * l.dt

	point := Point3D{X: l.x, Y: l.y, Z: l.z}

	// Add to trail
	l.trail = append(l.trail, point)
	if len(l.trail) > l.maxTrailLength {
		l.trail = l.trail[1:]
	}

	return point
}

// GetTrail returns the current trail
func (l *LorenzSystem) GetTrail() []Point3D {
	return l.trail
}

// Project3DTo2D projects a 3D point to 2D screen coordinates
func Project3DTo2D(p Point3D, width, height int, rotationX, rotationY float64) (int, int) {
	// No rotation - fixed camera view showing X-Z plane
	// This gives the classic butterfly view of the Lorenz attractor

	// Scale factors to fit the full attractor in the view
	// Lorenz attractor typically ranges: X(-20,20), Y(-30,30), Z(0,50)
	scaleX := float64(width) / 50.0  // X range roughly -25 to 25
	scaleZ := float64(height) / 60.0 // Z range roughly 0 to 50

	// Center the attractor in the image
	centerX := float64(width) / 2.0
	centerZ := float64(height) * 0.8 // Position towards bottom since Z starts at ~0

	// Project X and Z coordinates (classic butterfly view)
	projX := p.X*scaleX + centerX
	projZ := centerZ - p.Z*scaleZ // Flip Z so higher values are at top

	return int(projX), int(projZ)
}

// CreateFrame creates a single frame of the animation
func CreateFrame(lorenz *LorenzSystem, width, height int, rotationX, rotationY float64, frameNum int) *image.Paletted {
	// Create a paletted image for GIF
	palette := make(color.Palette, 256)

	// Create a gradient palette from black to bright colors
	for i := 0; i < 256; i++ {
		if i == 0 {
			palette[i] = color.RGBA{0, 0, 0, 255} // Black background
		} else {
			t := float64(i) / 255.0
			r := uint8(math.Sin(t*math.Pi*2)*127 + 128)
			g := uint8(math.Sin(t*math.Pi*2+math.Pi*2/3)*127 + 128)
			b := uint8(math.Sin(t*math.Pi*2+math.Pi*4/3)*127 + 128)
			palette[i] = color.RGBA{r, g, b, 255}
		}
	}

	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	// Fill with black background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetColorIndex(x, y, 0)
		}
	}

	trail := lorenz.GetTrail()
	if len(trail) < 2 {
		return img
	}

	// Draw the trail with fading effect
	for i := 1; i < len(trail); i++ {
		x1, y1 := Project3DTo2D(trail[i-1], width, height, rotationX, rotationY)
		x2, y2 := Project3DTo2D(trail[i], width, height, rotationX, rotationY)

		// Calculate color intensity based on age (newer points are brighter)
		intensity := float64(i) / float64(len(trail))
		colorIndex := uint8(intensity*254 + 1) // Avoid index 0 (black)

		// Draw line between consecutive points
		drawLine(img, x1, y1, x2, y2, colorIndex)
	}

	// Draw current position as a bright dot
	if len(trail) > 0 {
		current := trail[len(trail)-1]
		x, y := Project3DTo2D(current, width, height, rotationX, rotationY)
		drawCircle(img, x, y, 3, 255) // Bright white dot
	}

	// Add frame information
	drawText(img, 10, 20, fmt.Sprintf("Frame: %d", frameNum), 200)
	drawText(img, 10, 35, "Lorenz Attractor Animation", 150)

	return img
}

// drawLine draws a line between two points using Bresenham's algorithm
func drawLine(img *image.Paletted, x1, y1, x2, y2 int, colorIndex uint8) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	var sx, sy int
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		if x1 >= 0 && x1 < img.Bounds().Dx() && y1 >= 0 && y1 < img.Bounds().Dy() {
			img.SetColorIndex(x1, y1, colorIndex)
		}

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// drawCircle draws a filled circle
func drawCircle(img *image.Paletted, centerX, centerY, radius int, colorIndex uint8) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				px, py := centerX+x, centerY+y
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.SetColorIndex(px, py, colorIndex)
				}
			}
		}
	}
}

// drawText draws simple text (very basic implementation)
func drawText(img *image.Paletted, x, y int, text string, colorIndex uint8) {
	// Simple 3x5 bitmap font for basic characters
	for i, char := range text {
		drawChar(img, x+i*6, y, char, colorIndex)
	}
}

// drawChar draws a single character using a simple bitmap
func drawChar(img *image.Paletted, x, y int, char rune, colorIndex uint8) {
	// Very simple character rendering - just a few basic characters
	var pattern [][]bool

	switch char {
	case 'L':
		pattern = [][]bool{
			{true, false, false},
			{true, false, false},
			{true, false, false},
			{true, false, false},
			{true, true, true},
		}
	case 'o':
		pattern = [][]bool{
			{false, true, false},
			{true, false, true},
			{true, false, true},
			{true, false, true},
			{false, true, false},
		}
	case 'r':
		pattern = [][]bool{
			{false, false, false},
			{true, false, false},
			{true, true, false},
			{true, false, false},
			{true, false, false},
		}
	case 'e':
		pattern = [][]bool{
			{false, true, false},
			{true, false, true},
			{true, true, true},
			{true, false, false},
			{false, true, true},
		}
	case 'n':
		pattern = [][]bool{
			{false, false, false},
			{true, true, false},
			{true, false, true},
			{true, false, true},
			{true, false, true},
		}
	case 'z':
		pattern = [][]bool{
			{false, false, false},
			{true, true, true},
			{false, false, true},
			{false, true, false},
			{true, true, true},
		}
	case ' ':
		pattern = [][]bool{
			{false, false, false},
			{false, false, false},
			{false, false, false},
			{false, false, false},
			{false, false, false},
		}
	default:
		// Default to a simple dot for unknown characters
		pattern = [][]bool{
			{false, false, false},
			{false, false, false},
			{false, true, false},
			{false, false, false},
			{false, false, false},
		}
	}

	for row, line := range pattern {
		for col, pixel := range line {
			if pixel {
				px, py := x+col, y+row
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.SetColorIndex(px, py, colorIndex)
				}
			}
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CreateLorenzAnimation creates an animated GIF of the Lorenz attractor
func CreateLorenzAnimation(width, height, frames int, filename string) error {
	fmt.Printf("Creating Lorenz attractor animation with %d frames...\n", frames)

	lorenz := NewLorenzSystem(2000) // Keep trail of last 2000 points

	// Skip initial transient behavior
	for i := 0; i < 1000; i++ {
		lorenz.Step()
	}

	var images []*image.Paletted
	var delays []int

	for frame := 0; frame < frames; frame++ {
		// No rotation - fixed camera view
		rotationX := 0.0
		rotationY := 0.0

		// Step the system forward multiple times per frame for smoother trail
		for i := 0; i < 10; i++ {
			lorenz.Step()
		}

		// Create frame
		img := CreateFrame(lorenz, width, height, rotationX, rotationY, frame)
		images = append(images, img)
		delays = append(delays, 1) // 50ms delay between frames

		// Progress indicator
		if frame%10 == 0 {
			fmt.Printf("Progress: %d/%d frames\n", frame+1, frames)
		}
	}

	// Create GIF
	fmt.Printf("Encoding GIF...\n")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return gif.EncodeAll(file, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

// CreateRealTimeAnimation creates a simple real-time ASCII animation in the terminal
func CreateRealTimeAnimation() {

	lorenz := NewLorenzSystem(100)

	// Skip transient
	for i := 0; i < 1000; i++ {
		lorenz.Step()
	}

	width, height := 80, 24

	for frame := 0; frame < 1000; frame++ {
		// Clear screen
		fmt.Print("\033[2J\033[H")

		// Create ASCII art frame
		canvas := make([][]rune, height)
		for i := range canvas {
			canvas[i] = make([]rune, width)
			for j := range canvas[i] {
				canvas[i][j] = ' '
			}
		}

		// Step system and get trail
		for i := 0; i < 5; i++ {
			lorenz.Step()
		}

		trail := lorenz.GetTrail()
		// No rotation for terminal version either
		rotation := 0.0

		// Plot trail points
		for i, point := range trail {
			x, y := Project3DTo2D(point, width, height, 0, rotation)
			if x >= 0 && x < width && y >= 0 && y < height {
				// Different characters for different trail ages
				intensity := float64(i) / float64(len(trail))
				var char rune
				switch {
				case intensity > 0.9:
					char = '●'
				case intensity > 0.7:
					char = '◆'
				case intensity > 0.5:
					char = '▲'
				case intensity > 0.3:
					char = '♦'
				default:
					char = '·'
				}
				canvas[y][x] = char
			}
		}

		// Draw frame info
		frameInfo := fmt.Sprintf("Frame: %d | Lorenz Attractor", frame)
		for i, char := range frameInfo {
			if i < width {
				canvas[0][i] = char
			}
		}

		// Print canvas
		for _, row := range canvas {
			fmt.Println(string(row))
		}

		time.Sleep(30 * time.Millisecond)
	}
}

func main() {
	fmt.Println("🌀 LORENZ ATTRACTOR ANIMATION GENERATOR 🌀")
	fmt.Println("==========================================")

	// Create animated GIF with larger dimensions for better view
	fmt.Println("1. Creating animated GIF...")
	err := CreateLorenzAnimation(800, 600, 3600, "lorenz_animation.gif")
	if err != nil {
		log.Printf("Error creating animation: %v", err)
	}

	// Run real-time animation for a short demo
	go func() {
		time.Sleep(13 * time.Second)
		os.Exit(0)
	}()

	CreateRealTimeAnimation()
}
