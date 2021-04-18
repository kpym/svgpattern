// Package model provides the pattern templates for svgpattern.
package model

import (
	_ "embed"
)

//go:embed chevrons.svg
var Chevrons string

//go:embed  concentric-circles.svg
var ConcentricCircles string

//go:embed  diamonds.svg
var Diamonds string

//go:embed  hexagons.svg
var Hexagons string

//go:embed plaid.svg
var Plaid string

//go:embed squares.svg
var Squares string

var Models = []struct{ Name, Code string }{
	{"chevrons", Chevrons},
	{"concentric-circles", ConcentricCircles},
	{"diamonds", Diamonds},
	{"hexagons", Hexagons},
	{"plaid", Plaid},
	{"squares", Squares},
}

func GetModelCode(name string) (code string, ok bool) {
	for _, m := range Models {
		if m.Name == name {
			return m.Code, true
		}
	}

	return "", false
}

func GetModelIndex(name string) (index int, ok bool) {
	for i, m := range Models {
		if m.Name == name {
			return i, true
		}
	}

	return -1, false
}
