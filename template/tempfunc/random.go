// Package tempfunc provide template functions.
// The provided functions are of two kinds :
// for random generation, and
// for general (numerical) manipulation.
package tempfunc

import (
	"math"
	"math/rand"
	"text/template"
)

// RandomFunctions provides template functions for random generation/selection.
func RandomFunctions(seed int64) template.FuncMap {
	r := rand.New(rand.NewSource(seed))
	return map[string]interface{}{
		"randf": func(min interface{}, max interface{}) float64 { return randomFloatMinMax(r, min, max) },
		"randi": func(min interface{}, max interface{}) int { return randomIntMinMax(r, min, max) },
		"pick":  func(values ...interface{}) interface{} { return randomPick(r, values) },
	}
}

func randomFloatMinMax(r *rand.Rand, min interface{}, max interface{}) float64 {
	minf, maxf := toFloat64(min), toFloat64(max)
	return minf + r.Float64()*(maxf-minf)
}

func randomIntMinMax(r *rand.Rand, min interface{}, max interface{}) int {
	minn, maxn := int(math.Round(toFloat64(min))), int(math.Round(toFloat64(max)))
	if minn > maxn {
		minn, maxn = maxn, minn
	}

	return minn + r.Intn(maxn-minn+1)
}

func randomPick(r *rand.Rand, values []interface{}) interface{} {
	n := len(values)
	if n <= 0 {
		return nil
	}

	return values[r.Intn(n)]
}
