package svgpattern

import (
	"math"
	"math/rand"
	"testing"

	"github.com/kpym/svgpattern/template/model"
	"github.com/lucasb-eyer/go-colorful"
)

func TestDefalutNew(t *testing.T) {
	g := New("").(*generator)
	// Random generator
	if len(g.phrase) > 0 {
		t.Error("By default phrase shoud be empty.")
	}
	if g.rand == nil {
		t.Error("The random generator should be set.")
	}
	// template model
	if g.name == "" {
		t.Error("The generator template name is empty.")
	}
	if g.code == nil {
		t.Error("The generator template is not set.")
	}
	if !g.templateOk {
		t.Error("The templateOk is false.")
	}
	// template parameters
	colorCheck := (g.color.R > 1 || g.color.R < 0 ||
		g.color.G > 1 || g.color.G < 0 ||
		g.color.B > 1 || g.color.B < 0)
	if colorCheck {
		t.Errorf("The color values are outside [0,1] : (%v, %v, %v)", g.color.R, g.color.G, g.color.B)
	}
	if g.opacity != 1 {
		t.Errorf("The default opacity should be 1, but it is %v.", g.opacity)
	}
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if g.rotate != 0 {
		t.Errorf("The default rotate should be 0, but it is %v.", g.rotate)
	}
	if g.scale != 1 {
		t.Errorf("The default scale should be 0, but it is %v.", g.scale)
	}
	// status
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}
}

func TestAddError(t *testing.T) {
	g := New("").(*generator)

	g.addError("1")
	g.addError("2")

	if g.errors[0] != "1" || g.errors[1] != "2" {
		t.Error("Problem adding errors in the default generator.", g.errors)
	}
}

func TestPhraseSeed(t *testing.T) {
	g := New("Test").(*generator)

	// Random generator
	if string(g.phrase) != "Test" {
		t.Error("The phrase is not 'Test' but", g.phrase)
	}
	if g.seed != 7234017283807667300 {
		t.Error("The random seed for 'Test' should be 7234017283807667300, but it is", g.seed)
	}
	if g.rand == nil {
		t.Error("The random generator should be set.")
	}
	// template model
	if g.name != "hexagons" {
		t.Error("The generator template name shoud be 'hexagons' but it is.", g.name)
	}
	if g.code == nil {
		t.Error("The generator template is not set.")
	}
	if !g.templateOk {
		t.Error("The templateOk is false.")
	}
	// template parameters
	colorCheck := (g.color.R > 1 || g.color.R < 0 ||
		g.color.G > 1 || g.color.G < 0 ||
		g.color.B > 1 || g.color.B < 0)
	if colorCheck {
		t.Errorf("The color values are outside [0,1] : (%v, %v, %v)", g.color.R, g.color.G, g.color.B)
	}
	if g.opacity != 1 {
		t.Errorf("The default opacity should be 1, but it is %v.", g.opacity)
	}
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if g.rotate != 0 {
		t.Errorf("The default rotate should be 0, but it is %v.", g.rotate)
	}
	if g.scale != 1 {
		t.Errorf("The default scale should be 0, but it is %v.", g.scale)
	}
	// status
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}
}

func TestWithModel(t *testing.T) {
	modelName := model.Models[rand.Intn(len(model.Models))].Name
	g := New("", WithModel(modelName)).(*generator)
	// check the name
	if g.name != modelName {
		t.Errorf("The model is not set as desired, want: %s, got: %s", modelName, g.name)
	}
	all := make([]string, len(model.Models))
	for i, m := range model.Models {
		all[i] = m.Name
	}
	g.Options(WithModel(all...))
	// status
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}
}

func TestWithColor(t *testing.T) {
	color := "#010203"
	g := New("", WithColor(color)).(*generator)
	// check the name
	got := g.color.Hex()
	if got != color {
		t.Errorf("The color is not set as desired, want: %s, got: %s", color, got)
	}
	// status
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}

	color = "#778899"
	g.Options(WithColor("#777"), WithColor("#789"))
	// check the name
	got = g.color.Hex()
	if got != color {
		t.Errorf("The color is not set as desired, want: %s, got: %s", color, got)
	}
	// status
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}

	g.Options(WithColor("bidon"))
	// status
	if len(g.errors) != 1 {
		t.Error("There should be an error (bad color).", g.errors)
	}
}

func TestRandomizeHue(t *testing.T) {
	h, s, l := 180.0, 0.5, 0.5
	delta, d := 70.0, 0.002
	color := colorful.Hsl(h, s, l).Hex()
	g := New("", WithColor(color), RandomizeHue(delta)).(*generator)
	// check the color
	nh, ns, nl := g.color.Hsl()
	deltaH := math.Abs(h - nh)
	deltaS := math.Abs(s - ns)
	deltaL := math.Abs(l - nl)
	if deltaH > delta || deltaS > d || deltaL > d {
		t.Errorf("The color is not randomized as desired, delta should be less then (%f,%f,%f), but got: (%f,%f,%f)", delta, d, d, deltaH, deltaS, deltaL)
	}
	// status
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}
}

func TestRandomizeSaturation(t *testing.T) {
	h, s, l := 180.0, 0.5, 0.5
	delta, d := 0.1, 0.002
	color := colorful.Hsl(h, s, l).Hex()
	g := New("", WithColor(color), RandomizeSaturation(delta)).(*generator)
	// check the color
	nh, ns, nl := g.color.Hsl()
	deltaH := math.Abs(h - nh)
	deltaS := math.Abs(s - ns)
	deltaL := math.Abs(l - nl)
	if deltaH > d || deltaS > delta || deltaL > d {
		t.Errorf("The color is not randomized as desired, delta should be less then (%f,%f,%f), but got: (%f,%f,%f)", delta, d, d, deltaH, deltaS, deltaL)
	}
	// status
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}
}

func TestRandomizeLightness(t *testing.T) {
	h, s, l := 180.0, 0.5, 0.5
	delta, d := 0.1, 0.002
	color := colorful.Hsl(h, s, l).Hex()
	g := New("", WithColor(color), RandomizeLightness(delta)).(*generator)
	// check the color
	nh, ns, nl := g.color.Hsl()
	deltaH := math.Abs(h - nh)
	deltaS := math.Abs(s - ns)
	deltaL := math.Abs(l - nl)
	if deltaH > d || deltaS > d || deltaL > delta {
		t.Errorf("The color is not randomized as desired, delta should be less then (%f,%f,%f), but got: (%f,%f,%f)", delta, d, d, deltaH, deltaS, deltaL)
	}
	// status
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}
}

func TestWithoutColor(t *testing.T) {
	g := New("", WithoutColor()).(*generator)
	if g.opacity != 0 {
		t.Errorf("The color is not transparent as desired, got opacity: %f", g.opacity)
	}
	// status
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}

	g.Options(WithColor("#fff"))
	if g.opacity != 1 {
		t.Errorf("The color is not opaque as desired, got opacity: %f", g.opacity)
	}
	g.Options(WithoutColor())
	if g.opacity != 0 {
		t.Errorf("The color is not transparent as desired, got opacity: %f", g.opacity)
	}
	// status
	if !g.colorOk {
		t.Error("The colorOk is false.")
	}
	if len(g.errors) > 0 {
		t.Error("There are errors in the default generator.", g.errors)
	}
}

func TestWithRotation(t *testing.T) {
	angle := 7.0
	g := New("", WithRotation(angle)).(*generator)
	if g.rotate != angle {
		t.Errorf("The angle is not as desired, want: %f, got: %f.", angle, g.rotate)
	}
}

func TestWithRotationBetween(t *testing.T) {
	min, max := 7.0, 7.0
	g := New("", WithRotationBetween(min, max)).(*generator)
	if g.rotate != min {
		t.Errorf("The randomized rotate is not as desired, we should have %f = %f = %f.", min, g.rotate, max)
	}

	min, max = -70, 70
	g.Options(WithRotationBetween(max, min))
	if min > g.rotate || max < g.rotate {
		t.Errorf("The randomized rotate is not as desired, we should have %f <= %f <= %f.", min, g.rotate, max)
	}
}

func TestWithScale(t *testing.T) {
	factor := 1.4
	g := New("", WithScale(factor)).(*generator)
	if g.scale != factor {
		t.Errorf("The scale is not as desired, want: %f, got: %f.", factor, g.scale)
	}
}

func TestWithScaleBetween(t *testing.T) {
	min, max := 0.7, 0.7
	g := New("", WithScaleBetween(min, max)).(*generator)
	if min > g.scale || max < g.scale {
		t.Errorf("The randomized scale is not as desired, we should have %f = %f = %f.", min, g.scale, max)
	}

	min, max = 2.1, 2.8
	g.Options(WithScaleBetween(max, min))
	if min > g.scale || max < g.scale {
		t.Errorf("The randomized scale is not as desired, we should have %f <= %f <= %f.", min, g.scale, max)
	}
}
