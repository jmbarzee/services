package painter

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/repeatable"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Bounce is a Painter which provides produces colors bouncing between ColorStart and ColorEnd,
// starting at p.ColorStart and shifting in the direction specified by Up
type Bounce struct {
	ColorStart *color.HSLA
	ColorEnd   *color.HSLA
	Up         *bool
	Shifter    ifaces.Shifter
}

// Paint returns a color based on t
func (p Bounce) Paint(t time.Time) color.HSLA {
	if *p.Up {
		if p.ColorStart.H < p.ColorEnd.H {
			hDistance := p.ColorEnd.H - p.ColorStart.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(hShift)
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		} else {
			hDistance := p.ColorStart.H - p.ColorEnd.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(-hShift) // shifting past 0
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		}
	} else {
		if p.ColorStart.H > p.ColorEnd.H {
			hDistance := p.ColorStart.H - p.ColorEnd.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(hShift)
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		} else {
			hDistance := (1 - p.ColorStart.H) + p.ColorEnd.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(-hShift) // shifting past 0
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		}
	}
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Bounce) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if p.ColorStart == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.ColorStart = pa.SelectColor()
		})
	}
	if p.ColorEnd == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.ColorEnd = pa.SelectColor()
		})
	}
	if p.Up == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			b := repeatable.Chance(pa.Start(), .5)
			p.Up = &b
		})
	}
	if p.Shifter == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.Shifter = pa.SelectShifter()
		})
	} else {
		sFuncs = append(sFuncs, p.Shifter.GetStabilizeFuncs()...)
	}
	return sFuncs
}