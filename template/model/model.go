// Package model provides the pattern templates for svgpattern.
package model

import (
	"embed"
	"io/fs"
	"strings"
)

// Model represent a go-template model with name and svg code.
type Model struct {
	Name string
	Code string
}

// Models is the list of all available models.
var Models []Model

// GetModelIndex provides index in Models list of the desired model.
// If the model is not found the ok is false and indes is -1
func GetModelIndex(name string) (index int, ok bool) {
	for i, m := range Models {
		if m.Name == name {
			return i, true
		}
	}

	return -1, false
}

// SetModel append or replace an existing model.
func SetModel(name string, code string) {
	i, ok := GetModelIndex(name)
	if ok {
		Models[i].Code = code
	} else {
		Models = append(Models, Model{name, code})
	}
}

//go:embed svgmodels/*.template.svg
var files embed.FS

// Init the Models list with the embeded svg templates.
func init() {
	svgdir, _ := fs.Sub(files, "svgmodels")
	svgmodels, _ := fs.ReadDir(svgdir, ".")
	for _, svgmodfile := range svgmodels {
		fname := svgmodfile.Name()
		fdata, _ := fs.ReadFile(svgdir, fname)
		name := strings.TrimSuffix(fname, ".template.svg")
		code := string(fdata)
		SetModel(name, code)
	}
}
