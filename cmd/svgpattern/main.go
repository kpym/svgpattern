// Package main provides the command line interface to use svgpattern.
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kpym/svgpattern"
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
}

// generatorFromParameters provides a new Generator using the CLI parameters.
func generatorFromParameters() svgpattern.Generator {
	// The parameter global variables
	var (
		color      string
		hue        float64
		saturation float64
		lightness  float64
		model      string
		rotate     string
		scale      string
	)

	// preserve the declaration order of the flags
	flag.Usage = help
	flag.CommandLine.SortFlags = false
	// declare the flags
	flag.StringVarP(&color, "color", "c", "", "The background color in hex, like '#a17', or 'no' for transparent background.")
	flag.Float64VarP(&hue, "hue", "u", 0, "The hue variation, like -u=10 (using HSL colors). Valid only if color is set.")
	flag.Lookup("hue").NoOptDefVal = "180"
	flag.Float64VarP(&saturation, "saturation", "a", 0, "The saturation variation, like -u=0.2 (using HSL colors). Valid only if color is set.")
	flag.Lookup("saturation").NoOptDefVal = "0.1"
	flag.Float64VarP(&lightness, "lightness", "l", 0, "The lightness variation, like -l=0.3 (using HSL colors). Valid only if color is set.")
	flag.Lookup("lightness").NoOptDefVal = "0.1"
	flag.StringVarP(&model, "model", "m", "", "The pattern model. If multiple choices separate by comma.")
	flag.StringVarP(&rotate, "rotate", "r", "0", "Rotation angle in degree, or interval like '45-90'.")
	flag.Lookup("rotate").NoOptDefVal = "0-360"
	flag.StringVarP(&scale, "scale", "s", "1", "Scale factor, or interval like '.7-1'.")
	flag.Lookup("scale").NoOptDefVal = "0.7-1.4"
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
	// set the color
	color = strings.TrimSpace(color)
	if color != "" {
		if color == "no" {
			g.Options(svgpattern.WithoutColor())
		} else {
			g.Options(svgpattern.WithColor(color))
			if hue != 0 {
				g.Options(svgpattern.RandomizeHue(hue))
			}
			if saturation != 0 {
				g.Options(svgpattern.RandomizeSaturation(saturation))
			}
			if lightness != 0 {
				g.Options(svgpattern.RandomizeLightness(lightness))
			}
		}
	}
	// set the model
	if model != "" {
		set := strings.Split(model, ",")
		for i, m := range set {
			set[i] = strings.TrimSpace(m)
		}
		g.Options(svgpattern.WithModel(set...))
	}
	// set the rotation angle
	if rotate != "0" {
		anglestr := strings.Split(rotate, "-")
		if len(anglestr) > 0 {
			a, _ := strconv.ParseFloat(anglestr[0], 64)
			if len(anglestr) > 1 {
				b, _ := strconv.ParseFloat(anglestr[1], 64)
				g.Options(svgpattern.WithRotationBetween(a, b))
			} else {
				g.Options(svgpattern.WithRotation(a))
			}
		}
	}
	// set the scale factor
	if scale != "1" {
		factorstr := strings.Split(scale, "-")
		if len(factorstr) > 0 {
			a, _ := strconv.ParseFloat(factorstr[0], 64)
			if len(factorstr) > 1 {
				b, _ := strconv.ParseFloat(factorstr[1], 64)
				g.Options(svgpattern.WithScaleBetween(a, b))
			} else {
				g.Options(svgpattern.WithScale(a))
			}
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
