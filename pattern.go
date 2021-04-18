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
	"math"
	"math/rand"
	"text/template"
	"time"

	"github.com/kpym/svgpattern/template/model"
	"github.com/kpym/svgpattern/template/tempfunc"
	"github.com/lucasb-eyer/go-colorful"
)

// A generator is the internal implementation of Generator interface.
// A generator is a context for the Generate function.
//
// Random generator
//
// A phrase is converted to sha1 hash, from which an int64 is constructed.
// this integer is used to seed a random generator.
//
// Template model
//
// The template models are go-templates available in the module package
//
// Template parameters
//
// These parameters are (transformed and) passed to the template to generate
// the pattern. They are randomly selected if not determined by an option.
//
// Errors
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
	name       string
	code       *template.Template
	templateOk bool
	// template parameters
	color   colorful.Color
	opacity float64
	colorOk bool
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

	err := g.code.Execute(&result, data)
	if err != nil {
		g.addError("Error executing the template " + g.name)
		g.addError(err.Error())
	}

	return result.Bytes(), len(g.errors) == 0
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
	g := new(generator)
	g.phraseSeed(phrase)
	g.scale = 1

	g.Options(options...)

	if !g.templateOk {
		g.randomModel()
	}

	if !g.colorOk {
		g.randomColor()
	}

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

// setModel set the go-template used to produce the svg pattern.
// This is a low level function : no validity check is provided here.
// If the index is not in the desired range, the program will panic.
// This function is the main relation with the packages model and tempfunc.
func (g *generator) setModel(index int) {
	m := model.Models[index]
	g.name = m.Name
	rf := tempfunc.RandomFunctions(g.seed)
	uf := tempfunc.UtilFunctions()
	code, err := template.New(g.name).Funcs(rf).Funcs(uf).Parse(m.Code)
	if err != nil {
		g.addError("Error parsing template " + g.name + ": " + err.Error())
	}
	g.code = code
	g.templateOk = true
}

// randomModel pick a random model from all available models.
func (g *generator) randomModel() {
	g.setModel(g.rand.Intn(len(model.Models)))
}

// WithModel is a Generator option that select the model
// from the list of 'valid' models.
// If only one valid model name is provided it is used.
// If no such model is provided, then a random one is
// chosen among all models.
func WithModel(models ...string) Option {
	if len(models) == 0 {
		return func(g *generator) {
			g.addError("Empty set of models. Use random model.")
			g.randomModel()
		}
	}

	var invalid string
	valid := make([]int, 0, len(models))
	for _, modelName := range models {
		i, ok := model.GetModelIndex(modelName)
		if ok {
			valid = append(valid, i)
		} else {
			invalid += " " + modelName
		}
	}

	if len(valid) == 0 {
		return func(g *generator) {
			g.addError("Invalid models:" + invalid + ". Use random models.")
			g.randomModel()
		}
	}

	return func(g *generator) {
		g.setModel(valid[g.rand.Intn(len(valid))])
		if len(invalid) > 0 {
			g.addError("Invalid models:" + invalid + ". Choose from the others.")
		}
	}
}

// setColor set the background color for the svg pattern.
func (g *generator) setColor(color colorful.Color) {
	g.color = color
	g.opacity = 1
	g.colorOk = true
}

// setOpacity set the background opacity.
// Actually only 0 and 1 can be provided from the public interface.
func (g *generator) setOpacity(opacity float64) {
	g.opacity = opacity
	g.colorOk = true
}

// randomColor generate a random background color.
func (g *generator) randomColor() {
	randCol := colorful.Hsl(360*g.rand.Float64(), 0.4+0.3*g.rand.Float64(), 0.3+0.4*g.rand.Float64())
	g.setColor(randCol)
}

// Set the background color.
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

// RandomizeHue is a Generator option that randomize the color hue (of the HSL representation).
// The (absolute value of) delta parameter is the maximal deviation in degree of the already provided color.
// This option should be used after WithColor.
func RandomizeHue(delta float64) Option {
	return func(g *generator) {
		if !g.colorOk {
			g.addError("Can't randomize hue before to set color.")
			return
		}
		h, c, l := g.color.Hsl()
		rh := h + g.rd(delta)
		// no need to normalize rh, it is defined mod 360.
		g.color = colorful.Hsl(rh, c, l)
	}
}

// RandomizeSaturation is a Generator option that randomize the color saturation (of the HSL representation).
// The (absolute value of) delta parameter is the maximal deviation of the already provided color with saturation in [0,1].
// This option should be used after WithColor.
func RandomizeSaturation(delta float64) Option {
	return func(g *generator) {
		if !g.colorOk {
			g.addError("Can't randomize saturation before to set color.")
			return
		}
		h, s, l := g.color.Hsl()
		rs := s + g.rd(delta)
		rs = math.Min(math.Max(rs, 0), 1)
		g.color = colorful.Hsl(h, rs, l)
	}
}

// RandomizeLightness is a Generator option that randomize the color lightness (of the HSL representation).
// The (absolute value of) delta parameter is the maximal deviation of the already provided color with lightness in [0,1].
// This option should be used after WithColor.
func RandomizeLightness(delta float64) Option {
	return func(g *generator) {
		if !g.colorOk {
			g.addError("Can't randomize lightness before to set color.")
			return
		}
		h, s, l := g.color.Hsl()
		rl := l + g.rd(delta)
		rl = math.Min(math.Max(rl, 0), 1)
		g.color = colorful.Hsl(h, s, rl)
	}
}

// RandomizeLightness is a Generator option that remove the background color.
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

// WithRotation is a Generator option that provide an interval [min, max]
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

// WithScaleBetween is a Generator option that provide an interval [min, max]
// in witch the scale factor of the patter should be randomly chosen.
func WithScaleBetween(min, max float64) Option {
	mid := (max + min) / 2
	delta := (max - min) / 2
	return func(g *generator) {
		g.scale = mid + g.rd(delta)
	}
}
