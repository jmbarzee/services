package neopixel

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/shared"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe"
)

// ChandelierMedium is a Medium Chandelier (4 rings)
type ChandelierMedium struct {
	device.BasicDevice
	SmallRings []*Ring
	LargeRings []*Ring
	Top        space.Vector // Mounting location for the chandilier
}

// NewChandelierMedium returns a new Medium Chandelier
func NewChandelierMedium(top space.Vector, theta float32) ChandelierMedium {
	smallRings := make([]*Ring, 2)
	largeRings := make([]*Ring, 2)

	center := space.Vector{
		X: top.X,
		Y: top.Y,
		Z: top.Z - .6,
	}
	smallRings[0] = NewRing(center, 0.7, 0.0+theta, math.Pi/6)
	largeRings[0] = NewRing(center, 1.3, math.Pi/2+theta, math.Pi/6)

	center = space.Vector{
		X: top.X,
		Y: top.Y,
		Z: top.Z - 1.0,
	}
	smallRings[1] = NewRing(center, 0.7, math.Pi/2+theta, math.Pi/6)
	largeRings[1] = NewRing(center, 1.3, math.Pi+theta, math.Pi/6)

	return ChandelierMedium{
		SmallRings: smallRings,
		LargeRings: largeRings,
		Top:        top,
	}
}

// Allocate takes Vibes and Distributes them to the rings
func (c ChandelierMedium) Allocate(feeling vibe.Vibe) {
	newVibe := feeling.Stabalize()
	for _, smallRing := range c.SmallRings {
		smallRing.Allocate(newVibe)
	}
	for _, largeRing := range c.LargeRings {
		largeRing.Allocate(newVibe)
	}

}

// Render calls render on each of the rings and then appends all the lights
func (c ChandelierMedium) Render(t time.Time) []shared.Light {
	allLights := []shared.Light{}
	for i := 0; i < 3; i++ {
		smallLights := c.SmallRings[i].Render(t)
		allLights = append(allLights, smallLights...)

		largeLights := c.LargeRings[i].Render(t)
		allLights = append(allLights, largeLights...)
	}
	return allLights
}

// PruneEffects removes all effects from the rigns which have ended before a time t
func (c ChandelierMedium) PruneEffects(t time.Time) {
	for _, smallRing := range c.SmallRings {
		smallRing.PruneEffects(t)
	}
	for _, largeRing := range c.LargeRings {
		largeRing.PruneEffects(t)
	}
}
