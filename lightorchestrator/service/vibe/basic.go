package vibe

import (
	"fmt"
	"time"

	"github.com/jmbarzee/color"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/repeatable"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/painter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/shifter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
)

// Basic is a vibe which can produce most Effects
type Basic struct {
	span.Span
	count   int // incremented by StartAdvance()
	Effects []ifaces.Effect
}

var _ ifaces.Vibe = (*Basic)(nil)

// Duplicate creates a copy of a vibe and insures that
// the dupliacted vibe will stabalize/materialize differently
func (v *Basic) Duplicate() ifaces.Vibe {
	newVibe := *v
	(&newVibe).count++
	return &newVibe
}

// Stabilize locks in part of the visual representation of a vibe.
func (v *Basic) Stabilize() ifaces.Vibe {
	newVibe := *v
	sFuncs := newVibe.GetStabilizeFuncs()
	if len(sFuncs) == 0 {
		return &newVibe
	}
	option := repeatable.Option(newVibe.randSeed(), len(sFuncs))
	sFuncs[option](&newVibe)
	return &newVibe
}

// Materialize locks all remaining unlocked visuals of a vibe
// then returns the resulting effects
func (v *Basic) Materialize() []ifaces.Effect {
	for {
		sFuncs := v.GetStabilizeFuncs()
		if len(sFuncs) == 0 {
			break
		}
		for _, sFunc := range sFuncs {
			sFunc(v)
		}
	}
	return v.Effects
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (v *Basic) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	for _, e := range v.Effects {
		sFuncs = append(sFuncs, e.GetStabilizeFuncs()...)
	}
	if len(v.Effects) == 0 {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			v.Effects = append(v.Effects, v.SelectEffect())
		})
	}
	return sFuncs
}

func (v Basic) String() string {
	s := fmt.Sprintf("vibe.Basic{StartTime:%v, EndTime:%v, Effects:[", v.StartTime, v.EndTime)
	for i, e := range v.Effects {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", e)
	}
	s += "]}"
	return s
}

func (v *Basic) randSeed() time.Time {
	v.count++
	return v.Start().Add(time.Second * time.Duration(v.count))
}

// ======== ifaces.Palette implementation ========

// SelectBender returns a Bender
func (v *Basic) SelectBender() ifaces.Bender {
	options := []ifaces.Bender{
		//&bender.Static{},
		&bender.Linear{},
		&bender.Exponential{},
		&bender.Sinusoidal{},
		&bender.Combo{},
	}
	length := len(options)
	option := repeatable.Option(v.randSeed(), length)

	return options[option]
}

// SelectColor returns a Color
func (v *Basic) SelectColor() color.Color {
	length := len(color.AllColors)
	option := repeatable.Option(v.randSeed(), length)
	c := color.AllColors[option]
	return c.HSL()
}

// SelectDuration returns a Duration
func (v *Basic) SelectDuration() *time.Duration {
	min := time.Second / 10
	max := time.Second * 10
	d := repeatable.RandDuration(v.randSeed(), min, max)
	return &d
}

// SelectShift returns a Shift
func (v *Basic) SelectShift() *float64 {
	min := 0.005
	max := .50
	s := repeatable.RandShift(v.randSeed(), min, max, 0.001)
	return &s

}

// SelectShifter returns a Shifter
func (v *Basic) SelectShifter() ifaces.Shifter {
	options := []ifaces.Shifter{
		//&shifter.Static{},
		&shifter.Positional{},
		&shifter.Locational{},
		&shifter.Directional{},
		&shifter.Temporal{},
		&shifter.Combo{},
	}
	length := len(options)
	option := repeatable.Option(v.randSeed(), length)

	return options[option]
}

// SelectPainter returns a Painter
func (v *Basic) SelectPainter() ifaces.Painter {
	options := []ifaces.Painter{
		//&painter.Static{},
		&painter.Move{},
		&painter.Bounce{},
	}
	length := len(options)
	option := repeatable.Option(v.randSeed(), length)

	return options[option]
}

// SelectEffect returns a Effect
func (v *Basic) SelectEffect() ifaces.Effect {
	options := []ifaces.Effect{
		&effect.Solid{
			BasicEffect: effect.BasicEffect{
				Span: v.Span,
			},
		},
		&effect.Future{
			BasicEffect: effect.BasicEffect{
				Span: v.Span,
			},
		},
	}
	length := len(options)
	option := repeatable.Option(v.randSeed(), length)

	return options[option]
}
