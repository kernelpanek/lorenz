package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// Point3D represents a 3D point
type Point3D struct {
	X, Y, Z float64
}

// LorenzAttractor generates points following the Lorenz equations
type LorenzAttractor struct {
	sigma, rho, beta float64
	dt               float64
	x, y, z          float64
}

// NewLorenzAttractor creates a new Lorenz attractor with classic parameters
func NewLorenzAttractor() *LorenzAttractor {
	return &LorenzAttractor{
		sigma: 10.0,
		rho:   28.0,
		beta:  8.0 / 3.0,
		dt:    0.01,
		x:     1.0,
		y:     1.0,
		z:     1.0,
	}
}

// NextPoint calculates the next point in the Lorenz attractor
func (l *LorenzAttractor) NextPoint() Point3D {
	dx := l.sigma * (l.y - l.x)
	dy := l.x*(l.rho-l.z) - l.y
	dz := l.x*l.y - l.beta*l.z

	l.x += dx * l.dt
	l.y += dy * l.dt
	l.z += dz * l.dt

	return Point3D{X: l.x, Y: l.y, Z: l.z}
}

// LogisticMap demonstrates the logistic map equation
func LogisticMap(r, x float64) float64 {
	return r * x * (1 - x)
}

// GenerateLogisticSequence generates a sequence using the logistic map
func GenerateLogisticSequence(r, x0 float64, iterations int) []float64 {
	sequence := make([]float64, iterations)
	x := x0
	for i := 0; i < iterations; i++ {
		x = LogisticMap(r, x)
		sequence[i] = x
	}
	return sequence
}

// MandelbrotSet calculates if a point is in the Mandelbrot set
func MandelbrotSet(c complex128, maxIter int) int {
	z := complex(0, 0)
	for i := 0; i < maxIter; i++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			return i
		}
		z = z*z + c
	}
	return maxIter
}

// CreateLorenzImage generates a 2D projection of the Lorenz attractor
func CreateLorenzImage(width, height int, iterations int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with black background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	lorenz := NewLorenzAttractor()

	// Skip transient behavior
	for i := 0; i < 1000; i++ {
		lorenz.NextPoint()
	}

	// Scale factors for mapping to image coordinates
	scaleX := float64(width) / 60.0
	scaleY := float64(height) / 60.0
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	// Generate and plot points
	for i := 0; i < iterations; i++ {
		point := lorenz.NextPoint()

		// Project 3D to 2D (using X and Z coordinates)
		x := int(point.X*scaleX + centerX)
		y := int(point.Z*scaleY + centerY)

		if x >= 0 && x < width && y >= 0 && y < height {
			// Create a gradient effect based on Y coordinate
			intensity := uint8(math.Max(0, math.Min(255, (point.Y+30)*3)))
			img.Set(x, y, color.RGBA{intensity, intensity / 2, 255 - intensity, 255})
		}
	}

	return img
}

// CreateMandelbrotImage generates an image of the Mandelbrot set
func CreateMandelbrotImage(width, height int, maxIter int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	xmin, xmax := -2.5, 1.0
	ymin, ymax := -1.25, 1.25

	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			// Map pixel coordinates to complex plane
			x := xmin + (xmax-xmin)*float64(px)/float64(width)
			y := ymin + (ymax-ymin)*float64(py)/float64(height)
			c := complex(x, y)

			iter := MandelbrotSet(c, maxIter)

			if iter == maxIter {
				img.Set(px, py, color.RGBA{0, 0, 0, 255})
			} else {
				// Color based on iteration count
				r := uint8(iter * 255 / maxIter)
				g := uint8((iter * 7) % 255)
				b := uint8((iter * 13) % 255)
				img.Set(px, py, color.RGBA{r, g, b, 255})
			}
		}
	}

	return img
}

// SaveImage saves an image to a file
func SaveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// DemonstrateLogisticMap shows the behavior of the logistic map for different r values
func DemonstrateLogisticMap() {
	fmt.Println("=== LOGISTIC MAP DEMONSTRATION ===")
	fmt.Println("The logistic map: x(n+1) = r * x(n) * (1 - x(n))")
	fmt.Println()

	testValues := []float64{2.5, 3.2, 3.5, 3.8, 4.0}
	x0 := 0.5
	iterations := 20

	for _, r := range testValues {
		fmt.Printf("r = %.1f (x0 = %.1f):\n", r, x0)
		sequence := GenerateLogisticSequence(r, x0, iterations)

		fmt.Print("Sequence: ")
		for i, val := range sequence {
			if i > 10 {
				fmt.Print("...")
				break
			}
			fmt.Printf("%.4f ", val)
		}
		fmt.Println()

		// Analyze behavior
		if r < 3.0 {
			fmt.Println("Behavior: Converges to fixed point")
		} else if r < 3.449 {
			fmt.Println("Behavior: Oscillates between two values")
		} else if r < 3.544 {
			fmt.Println("Behavior: Period-4 cycle")
		} else if r < 3.569 {
			fmt.Println("Behavior: More complex periodic behavior")
		} else {
			fmt.Println("Behavior: Chaotic (sensitive to initial conditions)")
		}
		fmt.Println()
	}
}

// DemonstrateSensitivityToInitialConditions shows butterfly effect
func DemonstrateSensitivityToInitialConditions() {
	fmt.Println("=== SENSITIVITY TO INITIAL CONDITIONS ===")
	fmt.Println("Two Lorenz attractors with nearly identical starting conditions:")
	fmt.Println()

	lorenz1 := &LorenzAttractor{
		sigma: 10.0, rho: 28.0, beta: 8.0 / 3.0, dt: 0.01,
		x: 1.0, y: 1.0, z: 1.0,
	}

	lorenz2 := &LorenzAttractor{
		sigma: 10.0, rho: 28.0, beta: 8.0 / 3.0, dt: 0.01,
		x: 1.0001, y: 1.0, z: 1.0, // Tiny difference
	}

	fmt.Printf("Initial difference: %.6f\n", lorenz2.x-lorenz1.x)
	fmt.Println("Time\tSystem1_X\tSystem2_X\tDifference")
	fmt.Println("----\t---------\t---------\t----------")

	for i := 0; i < 20; i++ {
		p1 := lorenz1.NextPoint()
		p2 := lorenz2.NextPoint()
		diff := math.Abs(p1.X - p2.X)

		fmt.Printf("%.2f\t%9.4f\t%9.4f\t%10.6f\n",
			float64(i)*lorenz1.dt*100, p1.X, p2.X, diff)

		// Skip some iterations to show progression
		for j := 0; j < 100; j++ {
			lorenz1.NextPoint()
			lorenz2.NextPoint()
		}
	}
}

func main() {
	fmt.Println("🌀 CHAOS THEORY SHOWCASE 🌀")
	fmt.Println("==============================")
	fmt.Println()

	// Demonstrate the logistic map
	DemonstrateLogisticMap()

	// Demonstrate sensitivity to initial conditions
	DemonstrateSensitivityToInitialConditions()

	// Generate Lorenz attractor visualization
	fmt.Println("=== GENERATING LORENZ ATTRACTOR ===")
	fmt.Println("Creating lorenz_attractor.png...")
	lorenzImg := CreateLorenzImage(800, 600, 50000)
	if err := SaveImage(lorenzImg, "lorenz_attractor.png"); err != nil {
		log.Printf("Error saving Lorenz image: %v", err)
	} else {
		fmt.Println("✓ Lorenz attractor saved as lorenz_attractor.png")
	}

	// Generate Mandelbrot set
	fmt.Println("\n=== GENERATING MANDELBROT SET ===")
	fmt.Println("Creating mandelbrot_set.png...")
	mandelbrotImg := CreateMandelbrotImage(800, 600, 100)
	if err := SaveImage(mandelbrotImg, "mandelbrot_set.png"); err != nil {
		log.Printf("Error saving Mandelbrot image: %v", err)
	} else {
		fmt.Println("✓ Mandelbrot set saved as mandelbrot_set.png")
	}
}
