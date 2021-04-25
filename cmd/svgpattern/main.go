// Package main provides the command line interface to use svgpattern.
package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/kpym/svgpattern"
	"github.com/kpym/svgpattern/template/model"
	flag "github.com/spf13/pflag"
)

// The version that is set by goreleaser
// quelques variables globales
var version = "dev"

// Aide affiche l'aide d'utilisation
func help() {
	var out = os.Stderr
	fmt.Fprintf(out, "svgpattern (version: %s)\n\n", version)
	fmt.Fprintf(out, "Usage: svgpattern 'phrase' [parapeters].\nThe available parameters are:\n\n")
	flag.PrintDefaults()
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "The available pattern models are: %s\n\n", model.ModelsString())
}

// parseValues provide mean value and deviation from string
// "v~d" → v,d,true
// "v" → v,0,true
// "~d" → NaN,d,true
// "m:n" → (m+n)/2,|m-n|/2,true
// else → NaN,NaN,false
func parseValues(s string) (v, d float64, ok bool) {
	var (
		splt []string
		err  error
	)
	// if the string contains '~'
	splt = strings.Split(s, "~")
	if len(splt) > 1 {
		vstr := strings.TrimSpace(splt[0])
		if vstr == "" {
			v = math.NaN()
		} else {
			v, err = strconv.ParseFloat(vstr, 10)
			if err != nil {
				return math.NaN(), math.NaN(), false
			}
		}
		d, err = strconv.ParseFloat(strings.TrimSpace(splt[1]), 10)
		if err != nil {
			return math.NaN(), math.NaN(), false
		}
		return v, d, true
	}
	// if the string contains ':'
	splt = strings.Split(s, ":")
	if len(splt) > 1 {
		min, err := strconv.ParseFloat(strings.TrimSpace(splt[0]), 10)
		if err != nil {
			return math.NaN(), math.NaN(), false
		}
		max, err := strconv.ParseFloat(strings.TrimSpace(splt[1]), 10)
		if err != nil {
			return math.NaN(), math.NaN(), false
		}
		return (min + max) / 2, math.Abs(min-max) / 2, true
	}
	// if single value
	v, err = strconv.ParseFloat(strings.TrimSpace(splt[0]), 10)
	if err != nil {
		return math.NaN(), math.NaN(), false
	}

	return v, 0, true
}

// getOptions provide set/randomize options based on the strval paramater
// this function is to transform the values of the flags : u, a, l, r, s
func getOptions(strval string, with, rand func(float64) svgpattern.Option) (o []svgpattern.Option, ok bool) {
	val, dev, ok := parseValues(strval)
	if !ok {
		return nil, false
	}
	if !math.IsNaN(val) {
		o = append(o, with(val))
	}
	if !math.IsNaN(dev) && dev != 0 {
		o = append(o, rand(dev))
	}
	return
}

// generatorFromParameters provides a new Generator using the CLI parameters.
func generatorFromParameters() svgpattern.Generator {
	// The parameter global variables
	var (
		color      string
		hue        string
		saturation string
		lightness  string
		model      string
		rotate     string
		scale      string
	)

	// preserve the declaration order of the flags
	flag.Usage = help
	flag.CommandLine.SortFlags = false
	// declare the flags
	flag.StringVarP(&model, "model", "m", "", "The pattern model. If multiple choices separate by comma.")
	flag.StringVarP(&color, "color", "c", "", "The background color in hex, like '#a17', or 'no' for transparent background.")
	flag.StringVarP(&hue, "hue", "u", "", "The hue variation in degree (using HSL colors). Value format is '[value][~deviation]' or 'min:max'.")
	flag.StringVarP(&saturation, "saturation", "a", "", "The saturation variation (using HSL colors). Value format is '[value][~deviation]' or 'min:max'.")
	flag.StringVarP(&lightness, "lightness", "l", "", "The lightness variation (using HSL colors). Value format is '[value][~deviation]' or 'min:max'.")
	flag.StringVarP(&rotate, "rotate", "r", "", "Rotation angle in degree. Value format is '[value][~deviation]' or 'min:max'.")
	flag.StringVarP(&scale, "scale", "s", "", "Scale factor. Value format is '[value][~deviation]' or 'min:max'.")
	//parse the flags
	flag.Parse()
	// chack if parameters were provided
	if len(os.Args) == 1 {
		help()
		os.Exit(1)
	}
	// check the positional parameters
	if flag.NArg() != 1 {
		log("Exactly one positional parameter is expexted.\n")
		if flag.NArg() == 0 {
			log("No positional parameters were provided.\n")
		} else {
			log("Provided %d positional parameters: '%s'.\n", flag.NArg(), strings.Join(flag.Args(), "', '"))
		}
		os.Exit(1)
	}
	// seed the generator
	g := svgpattern.New(flag.Arg(0))
	// set the model
	if model != "" {
		set := strings.Split(model, ",")
		for i, m := range set {
			set[i] = strings.TrimSpace(m)
		}
		g.Options(svgpattern.WithModel(set...))
	}
	// set the color + opacity
	color = strings.TrimSpace(color)
	if color != "" {
		if color == "no" {
			g.Options(svgpattern.WithoutColor())
		} else {
			g.Options(svgpattern.WithColor(color))
		}
	}
	// set/randomize the hue, saturation, lightness, rotate and scale
	for _, par := range []struct {
		name, value string
		with, rand  func(float64) svgpattern.Option
	}{
		{"hue", hue, svgpattern.WithHue, svgpattern.RandomizeHue},
		{"saturation", saturation, svgpattern.WithSaturation, svgpattern.RandomizeSaturation},
		{"lightness", lightness, svgpattern.WithLightness, svgpattern.RandomizeLightness},
		{"rotate", rotate, svgpattern.WithRotation, svgpattern.RandomizeRotation},
		{"scale", scale, svgpattern.WithScale, svgpattern.RandomizeScale},
	} {
		if par.value != "" {
			o, ok := getOptions(par.value, par.with, par.rand)
			if !ok {
				log("Error parsing the %s parameter '%s'.\n", par.name, par.value)
				os.Exit(1)
			}
			g.Options(o...)
		}
	}

	return g
}

func log(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

// Prints pattern's SVG string with a specific background color
func main() {
	g := generatorFromParameters()
	svg, ok := g.Generate()
	if !ok {
		log("There are some errors : %v", g.Errors())
	}
	os.Stdout.Write(svg)
}
