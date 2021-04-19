package tempfunc

import (
	"fmt"
	"os"
	"testing"
	"text/template"
)

func TestToFloat(t *testing.T) {
	data := []struct {
		in  interface{}
		out float64
	}{
		{float64(1), 1},
		{uint8(2), 2},
		{uint16(3), 3},
		{uint32(4), 4},
		{uint64(5), 5},
		{int8(6), 6},
		{int16(7), 7},
		{int32(8), 8},
		{int64(9), 9},
		{float32(10), 10},
		{"1", 1},
		{"1.1", 1.1},
		{"bingo", 0},
	}
	for _, tt := range data {
		res := toFloat64(tt.in)
		if res != tt.out {
			t.Errorf("got %f, want %f", res, tt.out)
		}
	}
}

func TestFromTo(t *testing.T) {
	data := []struct {
		min, max int
		out      string
		msg      string
	}{
		{2, 4, "[2 3 4]", "three element list generation"},
		{4, 2, "[4 3 2]", "three element list generation"},
	}
	for _, tt := range data {
		res := fmt.Sprintf("%v", fromTo(tt.min, tt.max))
		if res != tt.out {
			t.Errorf(tt.msg+" got %s, want %s", res, tt.out)
		}
	}
}

func TestGridTo(t *testing.T) {
	data := []struct {
		n   int
		out string
		msg string
	}{
		{2, "[[0 2] [1]]", "problem with gridTo(2)"},
		{5, "[[0 5] [1] [2] [3] [4]]", "problem with gridTo(5)"},
		{-1, "[[0 1]]", "problem with gridTo(-1)"},
	}
	for _, tt := range data {
		res := fmt.Sprintf("%v", gridTo(tt.n))
		if res != tt.out {
			t.Errorf(tt.msg+" got %s, want %s", res, tt.out)
		}
	}
}

// The functions `number` and `float` are the same.
func ExampleUtilFunctions_number() {
	const hello string = `Hello, {{ "1.2e3" | number }} and {{ "not a number" | number }} !`
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(UtilFunctions()).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 1200 and 0 !
}

// The function `plus` can be used in the same way as `times`.
func ExampleUtilFunctions_times() {
	const hello string = `Hello, {{ "1.2e3" | times .01 }} and {{ times -1 "3" }} !`
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(UtilFunctions()).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 12 and -3 !
}

func ExampleUtilFunctions_round() {
	const hello string = `Hello, {{ "1.7e-3" | round "3" }} and {{ 1.00001 | round 1 }} !`
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(UtilFunctions()).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 0.002 and 1 !
}

// The parameters of `fromto` are rounded.
// The generated sequence include the second parameter.
func ExampleUtilFunctions_fromto() {
	const hello string = `{{ $r := fromto "1.1" -2 }} Hello{{ range $r }}, {{.}}{{ end }} !`
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(UtilFunctions()).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 1, 0, -1, -2 !
}

// The `upto` starts from 0 and do not include the (rounded) parameter value.
// If the parameter is 0 or less, an empty list is provided.
func ExampleUtilFunctions_upto() {
	const hello string = `{{ $r := upto 3 }} Hello{{ range $r }}, {{.}}{{ end }} !`
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(UtilFunctions()).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 0, 1, 2 !
}

func ExampleUtilFunctions_varset() {
	const hello string = `{{ $what := var "initial" -}}
	{{- if eq ("1.001" | round 1) "1" }}
		{{- "from inside" | set $what }}
	{{- end }}
	Hello {{ $what }} !`
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(UtilFunctions()).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello from inside !
}

func ExampleUtilFunctions_list() {
	const hello string = `Hello{{ range (list 1 2 4) }}, {{ . | times 2 }}{{ end }} !`
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(UtilFunctions()).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 2, 4, 8 !
}
