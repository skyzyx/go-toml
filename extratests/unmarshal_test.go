package extratests

import (
	"github.com/pelletier/go-toml"
	"reflect"
)
import "testing"

func TestUnmarshalStructDefault(t *testing.T) {
	type S struct {
		Field string `default:"hello"`
		Field2 string // no default
		Field3 string `default:"use default"`
		Field4 string // no default
	}
	s := S{
		Field3: "overwritten",
		Field4: "kept",
	}
	data := `
	`
	err := toml.Unmarshal([]byte(data), &s)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if s.Field != "hello" {
		t.Error("Field should be \"hello\", got", s.Field, "instead")
	}
	if s.Field2 != "" {
		t.Error("Field2 should be initialized as its type default (empty string)")
	}
	if s.Field3 != "use default" {
		t.Error("Field3 should be overwritten by the `default` tag's value")
	}
	if s.Field4 != "kept" {
		t.Error("Field4 should not be overwritten")
	}
}

func TestUnmarshalStructBasicFields(t *testing.T) {
	type S struct {
		String string
		Int int
		Int64 int64
	}
	s := S{}
	data := `
String = "str"
Int = 10
Int64 = 20
`
	err := toml.Unmarshal([]byte(data), &s)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if s.String != "str" {
		t.Error("String field should be \"str\", got", s.String, "instead")
	}
	if s.Int != 10 {
		t.Error("Int field should be 10, got", s.Int, "instead")
	}
	if s.Int64 != 20 {
		t.Error("Int64 field should be 20, got", s.Int64, "instead")
	}
}

func TestUnmarshalStructNestedEmpty(t *testing.T) {
	type Root struct {
		Main struct {
			MainDeep struct {
				Field int
			}
		}
	}
	r := Root{}
	err := toml.Unmarshal([]byte(``), &r)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if r.Main.MainDeep.Field != 0 {
		t.Fatal("nested structs are not initialized")
		// in reality, this is expected to panic if the are not initialized
	}
}

func TestUnmarshalStructNested(t *testing.T) {
	type Root struct {
		Main struct {
			MainDeep struct {
				Field int
			}
		}
	}
	r := Root{}
	data := `
[Main.MainDeep]
Field = 42
`
	err := toml.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if r.Main.MainDeep.Field != 42 {
		t.Error("nested structs are not parsed")
	}
}

func TestUnmarshalStructNestedPtrStruct(t *testing.T) {
	type Root struct {
		Main struct {
			MainDeep *struct {
				Field int
			}
		}
	}
	r := Root{}
	data := `
[Main.MainDeep]
Field = 42
`
	err := toml.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if r.Main.MainDeep.Field != 42 {
		t.Error("nested structs are not parsed")
	}
}


func TestUnmarshalMap(t *testing.T) {
	testToml := []byte(`
		a = 1
		b = 2
		c = 3
		`)
	var result map[string]int
	err := toml.Unmarshal(testToml, &result)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
		return
	}

	expected := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	}
}

func TestUnmarshalMapNoEraseRoot(t *testing.T) {
	testToml := []byte(`
		a = 1
		b = 2
		`)
	result := map[string]int{
		"a": 9999,
		"c": 3,
	}
	err := toml.Unmarshal(testToml, &result)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
		return
	}

	expected := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	}
}

func TestUnmarshalMapWithTypedKey(t *testing.T) {
	testToml := []byte(`
		a = 1
		b = 2
		c = 3
		`)

	type letter string
	var result map[letter]int
	err := toml.Unmarshal(testToml, &result)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
		return
	}

	expected := map[letter]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	}
}

func TestUnmarshalMapInStruct(t *testing.T) {
	type S struct {
		Things map[string]interface{}
	}

	data := []byte(`
[Things]
a=1
b="string"
`)
	s := S{}
	err := toml.Unmarshal(data, &s)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if s.Things["a"] != int64(1) {
		t.Errorf("a should be equal to 1, not %s", s.Things["a"])
	}

	if s.Things["b"] != "string" {
		t.Errorf("b should be equal to \"string\", not %s", s.Things["b"])
	}
}


func TestUnmarshal274(t *testing.T) {
	type configTypeInner struct {
		V2 int `default:"456"`
	}

	type configType struct {
		V1    int `default:"123"`
		Inner configTypeInner
	}


	var config configType
	err := toml.Unmarshal([]byte(""), &config)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if config.V1 != 123 {
		t.Error("V1 should be 123, not", config.V1)
	}
	if config.Inner.V2 != 456 {
		t.Error("V2 should be 456, not", config.Inner.V2)
	}
}

func TestUnmarshalMapInMap(t *testing.T) {
	var v map[string]map[string]interface{}
	data := []byte(`
[a]
aa = 1
[b.c]
bca = 1
bcb = "two"
	`)
	err := toml.Unmarshal(data, &v)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	expected := map[string]map[string]interface{}{
		"a": {
			"aa": int64(1),
		},
		"b": {
			"c": map[string]interface{}{
				"bca": int64(1),
				"bcb": "two",
			},
		},
	}

	if !reflect.DeepEqual(v, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, v)
	}
}

func TestUnmarshalArrayInMap(t *testing.T) {
	m := map[string]interface{}{}
	data := []byte(`
[hello]
world = [1,2,3]`)
	err := toml.Unmarshal(data, &m)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	expected := map[string]interface{}{
		"hello": map[string]interface{}{
			"world": []interface{}{int64(1),int64(2),int64(3)},
		},
	}
	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Bad unmarshal: expected\n%v\ngot\n%v", expected, m)
	}
}

func TestUnmarshalArrayInMapRoot(t *testing.T) {
	m := map[string]interface{}{}
	data := []byte(`world = [1,2,3]`)
	err := toml.Unmarshal(data, &m)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	expected := map[string]interface{}{
		"world": []interface{}{int64(1),int64(2),int64(3)},
	}
	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Bad unmarshal: expected\n%v\ngot\n%v", expected, m)
	}
}

func TestUnmarshalStructArray(t *testing.T) {
	type S struct{
		Ints []int64
	}
	s := S{}
	data := []byte(`Ints = [1,2,3]`)
	err := toml.Unmarshal(data, &s)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if len(s.Ints) != 3 {
		t.Fatal("Ints is supposed to have 3 elements, not", len(s.Ints))
	}
}
