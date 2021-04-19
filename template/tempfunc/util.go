// Package tempfunc provide template functions.
// The description is in random.go.

package tempfunc

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

// UtilFunctions provides a set of utility template functions.
func UtilFunctions() template.FuncMap {
	return map[string]interface{}{
		"float":  toFloat64,
		"number": toFloat64,
		"times":  times,
		"plus":   plus,
		"round":  round,
		"isodd":  isodd,
		"iseven": iseven,
		"fromto": fromTo,
		"upto":   upTo,
		"grid":   gridTo,
		"var":    newVar,
		"set":    setVar,
		"list":   list,
	}
}

// toFloat64 convert 'any' interface to float64.
// If the provided parameter n do not looks like a number
// its value is considered zero.
func toFloat64(n interface{}) float64 {
	// if numeric type
	switch n.(type) {
	case float64:
		return n.(float64)
	case uint8:
		return float64(n.(uint8))
	case uint16:
		return float64(n.(uint16))
	case uint32:
		return float64(n.(uint32))
	case uint64:
		return float64(n.(uint64))
	case int8:
		return float64(n.(int8))
	case int16:
		return float64(n.(int16))
	case int32:
		return float64(n.(int32))
	case int64:
		return float64(n.(int64))
	case float32:
		return float64(n.(float32))
	}

	f, _ := strconv.ParseFloat(fmt.Sprintf("%v", n), 64)
	return f
}

// times multiplies the two parameters.
func times(a, b interface{}) float64 {
	return toFloat64(a) * toFloat64(b)
}

// plus adds the two parameters.
func plus(a, b interface{}) float64 {
	return toFloat64(a) + toFloat64(b)
}

// round prints the f parameter with the provided precision.
// The non significant digits (0) are removed.
func round(precision interface{}, f interface{}) string {
	p := int(toFloat64(precision))
	if p < 0 {
		p = 0
	}
	sp := fmt.Sprintf("%d", p)
	sf := fmt.Sprintf("%."+sp+"f", toFloat64(f))
	if strings.Contains(sf, ".") {
		sf = strings.TrimRight(sf, "0")
		sf = strings.TrimRight(sf, ".")
	}
	if sf == "" {
		sf = "0"
	}
	return sf
}

// isodd verifies if n is odd.
func isodd(n interface{}) bool {
	return int(toFloat64(n))%2 == 1
}

// iseven verifies if n is even.
func iseven(n interface{}) bool {
	return int(toFloat64(n))%2 == 0
}

// fromTo generate a list of integers starting from
// the (integer part of) first to the (integer part of) last included.
func fromTo(first interface{}, last interface{}) (r []int) {
	fst, lst := int(toFloat64(first)), int(toFloat64(last))
	switch {
	case fst+1 <= lst:
		len := lst - fst + 1
		r = make([]int, len, len)
		for i, e := 0, fst; e <= lst; i, e = i+1, e+1 {
			r[i] = e
		}
	case fst-1 >= lst:
		len := fst - lst + 1
		r = make([]int, len, len)
		for i, e := 0, fst; e >= lst; i, e = i+1, e-1 {
			r[i] = e
		}
	default:
		r = []int{fst}
	}

	return r
}

// upTo generate a list of integers [0 1 ... n-1].
// If n is 0 or less, an empty list is provided.
func upTo(n interface{}) (r []int) {
	num := int(toFloat64(n))
	if num < 1 {
		return []int{}
	}

	return fromTo(0, num-1)
}

// gridTo provide a result like [[0 n] [1] ... [n-1]].
// This is very useful for the pattern construction where
// the first and the last element should be identical.
// So we fist draw the elements 0 and n, then 1,2,...,n-1
func gridTo(num interface{}) (r [][]int) {
	n := int(toFloat64(num))
	if n < 1 {
		n = 1
	}
	r = make([][]int, n, n)
	r[0] = []int{0, n}
	for i := 1; i < n; i++ {
		r[i] = []int{i}
	}

	return r
}

// newVar defines a new variable that can be modified
// in a sub scope and the modification will be preserved.
// Usage : {{ $a := var }}, then {{ 3 | set $a }}
func newVar(v interface{}) *interface{} {
	x := interface{}(v)
	return &x
}

// setVar modifies a variable initialised by newVar.
func setVar(x *interface{}, v interface{}) string {
	*x = v
	return ""
}

// list convert the set of parameters to a list.
func list(args ...interface{}) []interface{} {
	return args
}
