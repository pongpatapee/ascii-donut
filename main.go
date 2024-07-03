package main

import (
	"fmt"
	"math"
	"time"
)

const (
	screen_height int     = 50
	screen_width  int     = 50
	R1            float64 = 1.0
	R2            float64 = 2.0
	K2            float64 = 5.0
)

var (
	output  [][]string
	zBuffer [][]float64
)

func renderFrame(A, B float64) {
	for i := range screen_height {
		for j := range screen_width {
			output[i][j] = " "
			zBuffer[i][j] = 0.0
		}
	}

	thetaSpacing := 0.07
	phiSpacing := 0.02

	// R1 := 1.0
	// R2 := 2.0
	// K2 := 5.0
	// Calculate K1 based on screen size: the maximum x-distance occurs
	// roughly at the edge of the torus, which is at x=R1+R2, z=0.  we
	// want that to be displaced 3/8ths of the width of the screen, which
	// is 3/4th of the way from the center to the side of the screen.
	// screen_width*3/8 = K1*(R1+R2)/(K2+0)
	// screen_width*K2*3/(8*(R1+R2)) = K1
	K1 := float64(screen_width) * K2 * 3 / (8 * (R1 + R2))

	cosA := math.Cos(A)
	cosB := math.Cos(B)
	sinA := math.Sin(A)
	sinB := math.Sin(B)

	var x float64
	var y float64
	var z float64

	var xp int // x prime
	var yp int // y prime

	for theta := 0.0; theta < 2*math.Pi; theta += thetaSpacing {
		cosTheta, sinTheta := math.Cos(theta), math.Sin(theta)
		for phi := 0.0; phi < 2*math.Pi; phi += phiSpacing {
			cosPhi, sinPhi := math.Cos(phi), math.Sin(phi)

			x = (R2+R1*cosTheta)*(cosB*cosPhi+sinA*sinB*sinPhi) - (R1 * cosA * sinB * sinTheta)
			y = (R2+R1*cosTheta)*(sinB*cosPhi-cosB*sinA*sinPhi) + (R1 * cosA * cosB * sinTheta)
			z = cosA*(R2+R1*cosTheta)*sinPhi + (R1 * sinA * sinTheta) + K2

			ooz := 1 / z // one over z

			xp = int(float64(screen_width/2) + K1*ooz*x)
			yp = int(float64(screen_width/2) - K1*ooz*y)

			// Calculated from L = (Nx, Ny, Nz) dot (0, 1, -1) <- pre-chosen light vector
			luminance := cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi - sinA*sinTheta + cosB*(cosA*sinTheta-cosTheta*sinA*sinPhi)

			if luminance > 0 {
				// larger 1/z means pixel is closer so it should override for current x', y'
				if ooz > zBuffer[xp][yp] {
					luminance_index := int(luminance * 8)
					zBuffer[xp][yp] = ooz
					output[xp][yp] = string(".,-~:;=!*#$@"[luminance_index])
				}
			}

		}
	}
}

func main() {
	output = make([][]string, screen_height)
	for i := range output {
		output[i] = make([]string, screen_width)
	}

	zBuffer = make([][]float64, screen_height)
	for i := range zBuffer {
		zBuffer[i] = make([]float64, screen_width)
	}

	A := 0.0
	B := 0.0

	for {
		A += 0.04
		B += 0.02

		renderFrame(A, B)
		fmt.Print("\033[H\033[2J")
		for i := range screen_height {
			for j := range screen_width {
				fmt.Printf("%v", output[i][j])
			}
			fmt.Print("\n")
		}

		time.Sleep(25 * time.Millisecond)

	}
}
