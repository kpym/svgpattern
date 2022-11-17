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

// Models is just a list of models.
type Models []Model

// EmbeddedModels is the list of all available models.
var EmbeddedModels Models

// ModelsString provide the list of all available models as string.
func (models Models) ModelsString() string {
	n := len(models)
	s := make([]string, n, n)
	for i, m := range models {
		s[i] = m.Name
	}

	return strings.Join(s, ", ")
}

// GetModelIndex provides index in Models list of the desired model.
// If the model is not found the ok is false and indes is -1
func (models Models) GetModelIndex(name string) (index int, ok bool) {
	for i, m := range models {
		if m.Name == name {
			return i, true
		}
	}

	return -1, false
}

// SelectModels provides a list of models from the given list of names.
// If a name is not found in the list of available models it is added to the invalid list.
func (models Models) SelectModels(names ...string) (newModels Models, invalid []string) {
	for _, name := range names {
		i, ok := models.GetModelIndex(name)
		if ok {
			newModels = append(newModels, models[i])
		} else {
			invalid = append(invalid, name)
		}
	}

	return
}

// SetModel append or replace an existing model.
func (models *Models) SetModel(name string, code string) {
	i, ok := models.GetModelIndex(name)
	if ok {
		(*models)[i].Code = code
	} else {
		*models = append(*models, Model{name, code})
	}
}

//go:embed svgmodels/*.template.svg
var files embed.FS

// Init the Models list with the embedded svg templates.
func init() {
	svgdir, _ := fs.Sub(files, "svgmodels")
	svgmodels, _ := fs.ReadDir(svgdir, ".")
	for _, svgmodfile := range svgmodels {
		fname := svgmodfile.Name()
		fdata, _ := fs.ReadFile(svgdir, fname)
		name := strings.TrimSuffix(fname, ".template.svg")
		code := string(fdata)
		EmbeddedModels.SetModel(name, code)
	}
}
