// Package svgpattern generates 'random' svg patterns from text.
// The text is used to seed a sequence of random numbers that
// are used as parameters to generate a 'random' svg pattern.
//
// Some of the parameters can be fixed or randomly assigned
// depending on the generator options.
package svgpattern

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/kpym/svgpattern/template/model"
	"github.com/kpym/svgpattern/template/tempfunc"
	"github.com/lucasb-eyer/go-colorful"
)

// A generator is the internal implementation of Generator interface.
// A generator is a context for the Generate function.
//
// # Random generator
//
// A phrase is converted to sha1 hash, from which an int64 is constructed.
// this integer is used to seed a random generator.
//
// # Template model
//
// # The template models are go-templates available in the module package
//
// # Template parameters
//
// These parameters are (transformed and) passed to the template to generate
// the pattern. They are randomly selected if not determined by an option.
//
// # Errors
//
// No go errors are generated during the process.
// Only messages are saved in the errors field.
// If some error occurs, in place to stop the process, some random value is used.
type generator struct {
	// Random generator
	phrase string
	seed   int64
	rand   *rand.Rand
	// template model
	models model.Models
	name   string
	code   *template.Template
	// template parameters
	color   colorful.Color
	opacity float64
	rotate  float64
	scale   float64
	// status
	errors []string
}

// A Generator produce the svg pattern based on the provided Options.
// If some errors occurs during the initialization and/or the generation
// they are available through Errors() method.
type Generator interface {
	Options(...Option)
	Generate() (svg []byte, ok bool)
	Errors() []string
	Color() string
}

// Errors provide the error messages generated during the initialization/generation process.
func (g *generator) Errors() []string {
	return g.errors
}

// Generate provides the svg pattern as first parameter.
// The second parameter is true if no errors are present.
func (g *generator) Generate() (svg []byte, ok bool) {
	var result bytes.Buffer

	data := struct {
		Color   string
		Opacity float64
		Rotate  float64
		Scale   float64
	}{
		g.color.Hex(),
		g.opacity,
		g.rotate,
		g.scale,
	}

	if g.name == "" || g.code == nil {
		g.addError("Missing template.")
		return nil, false
	}
	err := g.code.Execute(&result, data)
	if err != nil {
		g.addError("Error executing the template " + g.name)
		g.addError(err.Error())
		return nil, false
	}

	return result.Bytes(), len(g.errors) == 0
}

// Color returns the color used to generate the pattern as hex string.
func (g *generator) Color() string {
	return g.color.Hex()
}

// An Option is a function that customize the Generator.
type Option func(*generator)

// Options applies all provided options to the Generator.
func (g *generator) Options(options ...Option) {
	for _, opt := range options {
		opt(g)
	}
}

// New initialize a new Generator from the provided phrase and options.
func New(phrase string, options ...Option) Generator {
	// create the pattern generator
	g := new(generator)
	// set the models to default
	g.models = model.EmbeddedModels
	// init the pattern generator
	g.phraseSeed(phrase)
	g.scale = 1
	g.randomColor()
	g.randomModel()
	// apply the provided options
	g.Options(options...)

	return g
}

// addError add an error message to the generator.
func (g *generator) addError(msg string) {
	if g.errors == nil {
		g.errors = make([]string, 0, 1)
	}

	g.errors = append(g.errors, msg)
}

// setSeed seed the random generator.
func (g *generator) setSeed(seed int64) {
	g.seed = seed
	g.rand = rand.New(rand.NewSource(seed))
}

// timeSeed is used if no phrase is provided.
// The resulting patter is not reproducible.
func (g *generator) timeSeed() {
	g.setSeed(time.Now().UTC().UnixNano())
}

// phraseSeed use the sha1 sum of the provided phrase to seed the generator.
// In this way the produced 'random' svg pattern will be reproducible.
func (g *generator) phraseSeed(phrase string) {
	if phrase == "" {
		g.timeSeed()
		return
	}

	g.phrase = phrase
	h := sha1.New()
	h.Write([]byte(g.phrase))
	hs := h.Sum(nil)
	// convert the first 8 bytes to single int64
	// this gives us eighteen quintillion four hundred forty-six quadrillion seven hundred forty-four trillion
	// seventy-three billion seven hundred nine million five hundred fifty-one thousand six hundred sixteen possibilities
	var seed int64
	for i := 0; i < 8; i++ {
		seed += int64(hs[0]) << (i * 8)
	}
	g.setSeed(seed)
}

// randomModel pick a random model from all available models.
func (g *generator) randomModel() {
	numModels := len(g.models)
	if numModels == 0 {
		g.addError("Can't pick a model from empty list.")
		return
	}

	index := g.rand.Intn(numModels)
	m := g.models[index]
	rf := tempfunc.RandomFunctions(g.seed)
	uf := tempfunc.UtilFunctions()
	code, err := template.New(m.Name).Funcs(rf).Funcs(uf).Parse(m.Code)
	if err != nil {
		g.addError("Error parsing template " + m.Name + ": " + err.Error())
		return
	}

	g.name = m.Name
	g.code = code
}

// WithModel is a Generator option that select the model
// from the list of 'valid' models.
// If only one valid model name is provided, it is used.
// If no such model is provided, then a random one is
// chosen among all models.
func WithModel(models ...string) Option {
	return func(g *generator) {
		var invalid []string
		g.models, invalid = g.models.SelectModels(models...)
		if len(invalid) > 0 {
			g.addError(fmt.Sprintf("The following %d models are invalid: %s.", len(invalid), strings.Join(invalid, ", ")))
		}

		if len(g.models) == 0 {
			g.addError("Empty set of models. Use all builtin models.")
			g.models = model.EmbeddedModels
		}

		g.randomModel()
		return
	}
}

// setColor set the background color for the svg pattern.
func (g *generator) setColor(color colorful.Color) {
	g.color = color
	g.opacity = 1
}

// setOpacity set the background opacity.
// Actually only 0 and 1 can be provided from the public interface.
func (g *generator) setOpacity(opacity float64) {
	g.opacity = opacity
}

// randomColor generate a random background color.
func (g *generator) randomColor() {
	randCol := colorful.Hsl(360*g.rand.Float64(), 0.3+0.1*g.rand.Float64(), 0.3+0.1*g.rand.Float64())
	g.setColor(randCol)
}

// WithColor sets the background color.
// If the color is not valid a random one is chosen.
func WithColor(hex string) Option {
	color, err := colorful.Hex(hex)
	if err != nil {
		return func(g *generator) {
			g.randomColor()
			g.addError("Error parsing color: " + hex + ".")
		}
	}

	return func(g *generator) {
		g.setColor(color)
	}
}

// rd (random deviation) is a utility function
// that provides a random number in the interval [-|delta|, |delta|].
func (g *generator) rd(delta float64) float64 {
	return delta * (1 - 2*g.rand.Float64())
}

// WithHue set the color hue of the HSL representation.
func WithHue(hue float64) Option {
	// no need to normalize hue, it is defined mod 360.
	return func(g *generator) {
		_, s, l := g.color.Hsl()
		g.color = colorful.Hsl(hue, s, l)
	}
}

// RandomizeHue is a Generator option that randomize the color hue (of the HSL representation).
// The (absolute value of) delta parameter is the maximal deviation in degree of the already provided color.
// This option should be used after WithColor.
func RandomizeHue(delta float64) Option {
	return func(g *generator) {
		h, s, l := g.color.Hsl()
		rh := h + g.rd(delta)
		// no need to normalize rh, it is defined mod 360.
		g.color = colorful.Hsl(rh, s, l)
	}
}

// WithSaturation set the color saturation of the HSL representation.
func WithSaturation(sat float64) Option {
	sat = math.Min(math.Max(sat, 0), 1)
	return func(g *generator) {
		h, _, l := g.color.Hsl()
		g.color = colorful.Hsl(h, sat, l)
	}
}

// RandomizeSaturation is a Generator option that randomize the color saturation (of the HSL representation).
// The (absolute value of) delta parameter is the maximal deviation of the already provided color with saturation in [0,1].
// This option should be used after WithColor.
func RandomizeSaturation(delta float64) Option {
	return func(g *generator) {
		h, s, l := g.color.Hsl()
		rs := s + g.rd(delta)
		rs = math.Min(math.Max(rs, 0), 1)
		g.color = colorful.Hsl(h, rs, l)
	}
}

// WithLightness set the color lightness of the HSL representation.
func WithLightness(light float64) Option {
	light = math.Min(math.Max(light, 0), 1)
	return func(g *generator) {
		h, s, _ := g.color.Hsl()
		g.color = colorful.Hsl(h, s, light)
	}
}

// RandomizeLightness is a Generator option that randomize the color lightness (of the HSL representation).
// The (absolute value of) delta parameter is the maximal deviation of the already provided color with lightness in [0,1].
// This option should be used after WithColor.
func RandomizeLightness(delta float64) Option {
	return func(g *generator) {
		h, s, l := g.color.Hsl()
		rl := l + g.rd(delta)
		rl = math.Min(math.Max(rl, 0), 1)
		g.color = colorful.Hsl(h, s, rl)
	}
}

// WithoutColor is a Generator option that remove the background color.
// If this option is used only the patterns are present, without a background.
// Using this option allows to stack the svg pattern on top of a gradient or image.
func WithoutColor() Option {
	return func(g *generator) {
		g.setOpacity(0)
	}
}

// WithRotation is a Generator option that fix the pattern rotation transformation.
func WithRotation(angle float64) Option {
	return func(g *generator) {
		g.rotate = angle
	}
}

// RandomizeRotation is a Generator option that randomize the rotation angle.
func RandomizeRotation(delta float64) Option {
	return func(g *generator) {
		g.rotate += g.rd(delta)
	}
}

// WithRotationBetween is a Generator option that provide an interval [min, max]
// in witch the rotation angle of the patter should be randomly chosen.
func WithRotationBetween(min, max float64) Option {
	mid := (max + min) / 2
	delta := (max - min) / 2
	return func(g *generator) {
		g.rotate = mid + g.rd(delta)
	}
}

// WithScale is a Generator option that fix the scale factor transformation.
func WithScale(factor float64) Option {
	return func(g *generator) {
		g.scale = factor
	}
}

// RandomizeScale is a Generator option that randomize the scale factor.
func RandomizeScale(delta float64) Option {
	return func(g *generator) {
		g.scale += g.rd(delta)
	}
}

// WithScaleBetween is a Generator option that provide an interval [min, max]
// in witch the scale factor of the patter should be randomly chosen.
func WithScaleBetween(min, max float64) Option {
	mid := (max + min) / 2
	delta := (max - min) / 2
	return func(g *generator) {
		g.scale = mid + g.rd(delta)
	}
}
