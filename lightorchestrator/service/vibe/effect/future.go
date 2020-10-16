package effect

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Future is an Effect which displays each consecutive light
// as the "future" of the previous light
type Future struct {
	BasicEffect
	Painter      ifaces.Painter
	TimePerLight *time.Duration
}

// Render will produce a slice of lights based on the time and properties of lights
func (e Future) Render(t time.Time, lights []light.Light) []light.Light {
	for i := range lights {
		distanceInFuture := *e.TimePerLight * time.Duration(i)
		c := e.Painter.Paint(t.Add(distanceInFuture))
		lights[i].SetColor(c)
	}
	return lights
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (e *Future) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if e.Painter == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			e.Painter = pa.SelectPainter()
		})
	} else {
		sFuncs = append(sFuncs, e.Painter.GetStabilizeFuncs()...)
	}
	if e.TimePerLight == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			e.TimePerLight = pa.SelectDuration()
		})
	}

	return sFuncs
}