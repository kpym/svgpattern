package tempfunc

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"text/template"
)

func TestRandomFloats(t *testing.T) {

	// basic test
	data := []struct {
		min interface{}
		max interface{}
		out string
		msg string
	}{
		{1, 1.0, "1.0000000", "random float between 1 and 1 should be 1"},

		{"bingo", nil, "0.0000000", "non numeric values should be converted to 0"},
		{1.0, int(2), "1.3730284", "any numeric type should be accepted as min and max"},
		{"1e0", 2, "1.3730284", "strings should be parsed to floats"},
		{2, 1, "1.6269716", "if max < min the the random should be in ]max, min]"},
	}
	for _, tt := range data {
		randf := RandomFunctions(42)["randf"].(func(min interface{}, max interface{}) float64)
		res := fmt.Sprintf("%.7f", randf(tt.min, tt.max))
		if res != tt.out {
			t.Errorf(tt.msg+", got %s, want %s", res, tt.out)
		}
	}

	// same or not ?
	rand1 := RandomFunctions(42)["randf"].(func(min interface{}, max interface{}) float64)(1, 2)
	rand2 := RandomFunctions(42)["randf"].(func(min interface{}, max interface{}) float64)(1, 2)
	randf := RandomFunctions(7)["randf"].(func(min interface{}, max interface{}) float64)
	rand3 := randf(1, 2)
	rand4 := randf(1, 2)
	if rand1 != rand2 {
		t.Errorf("the two results must be the same : %f and %f", rand1, rand2)
	}
	if rand1 == rand3 {
		t.Errorf("the two results must be different : %f and %f", rand1, rand3)
	}
	if rand3 == rand4 {
		t.Errorf("the two results must be different : %f and %f", rand3, rand4)
	}
}

func TestRandomInts(t *testing.T) {

	// basic test
	data := []struct {
		min interface{}
		max interface{}
		out int
		msg string
	}{
		{1, 1.0, 1, "random integer between 1 and 1 should be 1"},
		{"bingo", nil, 0, "non numeric values should be converted to 0"},
		{7.1, 42, 24, "the float numbers should be rounded"},
		{"7e0", int(42), 24, "should be able to parse strings"},
		{42, 7, 24, "random between 7 and 42 should be the same as between 42 and 7"},
	}
	for _, tt := range data {
		randi := RandomFunctions(42)["randi"].(func(min interface{}, max interface{}) int)
		res := randi(tt.min, tt.max)
		if res != tt.out {
			t.Errorf(tt.msg+", got %d, want %d", res, tt.out)
		}
	}

	// same or not ?
	rand1 := RandomFunctions(42)["randi"].(func(min interface{}, max interface{}) int)(1, 42)
	rand2 := RandomFunctions(42)["randi"].(func(min interface{}, max interface{}) int)(1, 42)
	randi := RandomFunctions(7)["randi"].(func(min interface{}, max interface{}) int)
	rand3 := randi(1, 42)
	rand4 := randi(1, 42)
	if rand1 != rand2 {
		t.Errorf("the two results must be the same : %d and %d", rand1, rand2)
	}
	if rand1 == rand3 {
		t.Errorf("the two results must be different : %d and %d", rand1, rand3)
	}
	if rand3 == rand4 {
		t.Errorf("the two results must be different : %d and %d", rand3, rand4)
	}
}

func TestRandomPick(t *testing.T) {

	// basic test
	data := []struct {
		values []interface{}
		out    interface{}
		msg    string
	}{
		{[]interface{}{1}, 1, "pick between 1 element should be non random"},
		{[]interface{}{"apple", "tomato", "banana"}, "banana", "pick a string test"},
		{[]interface{}{"apple", 1, 2.0}, 2.0, "pick a string test"},
	}
	for _, tt := range data {
		pick := RandomFunctions(42)["pick"].(func(values ...interface{}) interface{})
		res := pick(tt.values...)
		if res != tt.out {
			t.Errorf(tt.msg+", got %v, want %v", res, tt.out)
		}
	}

	// same or not ?
	pick1 := RandomFunctions(42)["pick"].(func(values ...interface{}) interface{})(1, 2, 3, 4, 5)
	pick2 := RandomFunctions(42)["pick"].(func(values ...interface{}) interface{})(1, 2, 3, 4, 5)
	pick := RandomFunctions(7)["pick"].(func(values ...interface{}) interface{})
	pick3 := pick(1, 2, 3, 4, 5)
	pick4 := pick(1, 2, 3, 4, 5)
	if pick1 != pick2 {
		t.Errorf("the two results must be the same : %v and %v", pick1, pick2)
	}
	if pick1 == pick3 {
		t.Errorf("the two results must be different : %v and %v", pick1, pick3)
	}
	if pick3 == pick4 {
		t.Errorf("the two results must be different : %v and %v", pick3, pick4)
	}

}

func TestRandomInTemplate(t *testing.T) {
	var result bytes.Buffer

	// basic test
	data := []struct {
		in  string
		out string
		msg string
	}{
		{`Hello, {{ randf 1 1 }}!`, `Hello, 1!`, "random float between 1 and 1 should be 1"},
		{`Hello, {{ randi .9 1.1 }}!`, `Hello, 1!`, "random integer between .9 and 1.1 should be 1"},
		{`Hello, {{ pick 1 }}!`, `Hello, 1!`, "random choice of {1} should be 1"},
		{`Hello, {{ pick 1 1 1 1 }}!`, `Hello, 1!`, "random choice of 1s should be 1"},
		{`Hello, {{ randf 7 42 | printf "%.1f"}}!`, `Hello, 20.1!`, "test reproducible random float"},
		{`Hello, {{ randi 7 42 }}!`, `Hello, 24!`, "test reproducible random int"},
		{`Hello, {{ pick 1 "James" -3 nil }}!`, `Hello, James!`, "test reproducible random choice"},
	}
	for _, tt := range data {
		rf := RandomFunctions(42)
		tmpl, err := template.New("foo").Funcs(rf).Parse(tt.in)
		if err != nil {
			t.Error("error parsing the template", tt.in)
		}
		result.Reset()
		err = tmpl.Execute(&result, nil)
		if err != nil {
			t.Error("error executing the template", tt.in)
		}
		res := result.String()
		if res != tt.out {
			t.Errorf(tt.msg+", got %v, want %v", res, tt.out)
		}
	}
}

// The function `randf` provides a random float between the two parameters (excluding the second one).
// The function `randi` provides a random integer between (including) the two parameters.
func ExampleRandomFunctions_rand() {
	// the parameters of randf are loosely numeric
	const hello string = `Hello, {{ randf "1" 3 }} and {{ randi "1" 3 }} !`
	// the random generator use 42 as seed
	rf := RandomFunctions(42)
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(rf).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 1.7460567220932652 and 3 !
}

// The function `pick` select one of the provided parameters.
func ExampleRandomFunctions_pick() {
	const hello string = `Hello, {{ pick "apple" "banana" 7 }}!`
	// the random generator use 42 as seed
	rf := RandomFunctions(42)
	// compile and execute the template (without error check, very bad idea!)
	t, _ := template.New("hi").Funcs(rf).Parse(hello)
	t.Execute(os.Stdout, nil)
	// Output:
	// Hello, 7!
}
