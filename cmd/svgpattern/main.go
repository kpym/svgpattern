// Package main provides the command line interface to use svgpattern.
package main

import (
	"fmt"
	"math"
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

// parse values
func parseValues(s string) (v, d float64, ok bool) {
	var (
		splt []string
		err  error
	)
	// if the string contains ~
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
	// if the string contains :
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
	flag.StringVarP(&color, "color", "c", "", "The background color in hex, like '#a17', or 'no' for transparent background.")
	flag.StringVarP(&hue, "hue", "u", "", "The hue variation in degree (using HSL colors). Value format is '[value][~deviation]' or 'min:max'.")
	flag.StringVarP(&saturation, "saturation", "a", "", "The saturation variation (using HSL colors). Value format is '[value][~deviation]' or 'min:max'.")
	flag.StringVarP(&lightness, "lightness", "l", "", "The lightness variation (using HSL colors). Value format is '[value][~deviation]' or 'min:max'.")
	flag.StringVarP(&model, "model", "m", "", "The pattern model. If multiple choices separate by comma.")
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
	// set the color
	color = strings.TrimSpace(color)
	if color == "no" {
		g.Options(svgpattern.WithoutColor())
	} else {
		if color != "" {
			g.Options(svgpattern.WithColor(color))
		}
		// set/randomize the hue ?
		if hue != "" {
			hv, hd, ok := parseValues(hue)
			if ok {
				if !math.IsNaN(hv) {
					g.Options(svgpattern.WithHue(hv))
				}
				if !math.IsNaN(hd) && hd != 0 {
					g.Options(svgpattern.RandomizeHue(hd))
				}
			} else {
				log("Error parsing the hue parameter '%s'.\n", hue)
				os.Exit(1)
			}
		}
		// set/randomize the saturation ?
		if saturation != "" {
			sv, sd, ok := parseValues(saturation)
			if ok {
				if !math.IsNaN(sv) {
					g.Options(svgpattern.WithSaturation(sv))
				}
				if !math.IsNaN(sd) && sd != 0 {
					g.Options(svgpattern.RandomizeSaturation(sd))
				}
			} else {
				log("Error parsing the saturation parameter '%s'.\n", saturation)
				os.Exit(1)
			}
		}
		// set/randomize the lightness ?
		if lightness != "" {
			lv, ld, ok := parseValues(lightness)
			if ok {
				if !math.IsNaN(lv) {
					g.Options(svgpattern.WithLightness(lv))
				}
				if !math.IsNaN(ld) && ld != 0 {
					g.Options(svgpattern.RandomizeLightness(ld))
				}
			} else {
				log("Error parsing the lightness parameter '%s'.\n", lightness)
				os.Exit(1)
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
	if rotate != "" {
		av, ad, ok := parseValues(rotate)
		if ok {
			if !math.IsNaN(av) {
				g.Options(svgpattern.WithRotation(av))
			}
			if !math.IsNaN(ad) && ad != 0 {
				g.Options(svgpattern.RandomizeRotation(ad))
			}
		} else {
			log("Error parsing the rotation parameter '%s'.\n", rotate)
			os.Exit(1)
		}
	}
	// set the scale factor
	if scale != "" {
		scv, scd, ok := parseValues(scale)
		if ok {
			if !math.IsNaN(scv) {
				g.Options(svgpattern.WithScale(scv))
			}
			if !math.IsNaN(scd) && scd != 0 {
				g.Options(svgpattern.RandomizeScale(scd))
			}
		} else {
			log("Error parsing the rotation parameter '%s'.\n", scale)
			os.Exit(1)
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
