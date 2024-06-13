package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/youngpto/funs_tool/datapack/format"
	"io"
	"os"
	"strings"
	"text/template"
)

const (
	NilType int = iota
	BoolType
	IntType
	FloatType
	StringType
	ArrayType
	MapType
)

var GoTypes = []string{"interface{}", "bool", "int", "float64", "string", "*fs_json.Object", "*fs_json.Array"}

type Object struct {
	content map[string]interface{}
}

func NewObject() *Object {
	return &Object{
		content: make(map[string]interface{}),
	}
}

func (o *Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.content)
}

func (o *Object) UnmarshalJSON(bytes []byte) error {
	o.content = make(map[string]interface{})
	if err := json.Unmarshal(bytes, &o.content); err != nil {
		return err
	}
	return nil
}

func (o *Object) String() string {
	return fmt.Sprintln(o.content)
}

func (o *Object) Set(key string, value interface{}) {
	o.content[key] = value
}

func (o *Object) SetDefault(key string, value interface{}) interface{} {
	if v, ok := o.content[key]; ok {
		return v
	}
	o.content[key] = value
	return value
}

func (o *Object) Get(key string) interface{} {
	return o.content[key]
}

func (o *Object) GetInt(key string) (r int, b bool) {
	in := o.Get(key)
	r, b = in.(int)
	return
}

func (o *Object) GetString(key string) (r string, b bool) {
	in := o.Get(key)
	r, b = in.(string)
	return
}

func (o *Object) GetBool(key string) (r bool, b bool) {
	in := o.Get(key)
	r, b = in.(bool)
	return
}

func (o *Object) GetFloat(key string) (r float64, b bool) {
	in := o.Get(key)
	r, b = in.(float64)
	return
}

func (o *Object) GetObject(key string) (r *Object, b bool) {
	in := o.Get(key)
	m, ok := in.(map[string]interface{})
	if !ok {
		return
	}
	return &Object{content: m}, true
}

func (o *Object) GetArray(key string) (r *Array, b bool) {
	in := o.Get(key)
	m, ok := in.([]interface{})
	if !ok {
		return
	}
	return &Array{content: m}, true
}

type Array struct {
	content []interface{}
}

func NewArray() *Array {
	return &Array{
		content: make([]interface{}, 0),
	}
}

func (a *Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.content)
}

func (a *Array) UnmarshalJSON(bytes []byte) error {
	a.content = make([]interface{}, 0)
	if err := json.Unmarshal(bytes, &a.content); err != nil {
		return err
	}
	return nil
}

func (a *Array) Append(in interface{}) {
	a.content = append(a.content, in)
}

func (a *Array) inRange(idx int) bool {
	return idx >= 0 && idx < len(a.content)
}

func (a *Array) Get(idx int) interface{} {
	if !a.inRange(idx) {
		return nil
	}
	return a.content[idx]
}

func (a *Array) GetInt(idx int) (r int, b bool) {
	in := a.Get(idx)
	r, b = in.(int)
	return
}

func (a *Array) GetString(idx int) (r string, b bool) {
	in := a.Get(idx)
	r, b = in.(string)
	return
}

func (a *Array) GetBool(idx int) (r bool, b bool) {
	in := a.Get(idx)
	r, b = in.(bool)
	return
}

func (a *Array) GetFloat(idx int) (r float64, b bool) {
	in := a.Get(idx)
	r, b = in.(float64)
	return
}

func (a *Array) GetObject(idx int) (r *Object, b bool) {
	in := a.Get(idx)
	m, ok := in.(map[string]interface{})
	if !ok {
		return
	}
	return &Object{content: m}, true
}

func (a *Array) GetArray(idx int) (r *Array, b bool) {
	in := a.Get(idx)
	m, ok := in.([]interface{})
	if !ok {
		return
	}
	return &Array{content: m}, true
}

func (a *Array) String() string {
	return fmt.Sprintln(a.content)
}

func (a *Array) Len() int {
	return len(a.content)
}

func CheckArray(array *Array) (typ int, ok bool) {
	for _, in := range array.content {
		t := whatType(in)
		if typ == NilType {
			typ = t
		} else {
			if typ != t {
				return NilType, false
			}
		}
	}
	return typ, true
}

func whatType(in interface{}) int {
	if in == nil {
		return NilType
	} else if isBool(in) {
		return BoolType
	} else if isInt(in) {
		return IntType
	} else if isFloat(in) {
		return FloatType
	} else if isString(in) {
		return StringType
	} else if isArray(in) {
		return ArrayType
	} else if isMap(in) {
		return MapType
	}
	return NilType
}

func isBool(in interface{}) bool {
	_, ok := in.(bool)
	return ok
}

func isInt(in interface{}) bool {
	_, ok := in.(int)
	return ok
}

func isFloat(in interface{}) bool {
	_, ok := in.(float64)
	return ok
}

func isString(in interface{}) bool {
	_, ok := in.(string)
	return ok
}

func isArray(in interface{}) bool {
	_, ok := in.([]interface{})
	return ok
}

func isMap(in interface{}) bool {
	_, ok := in.(map[string]interface{})
	return ok
}

var privateStructTmpl = `struct {
{{ range . }}	// {{ .Comment }}
	{{ .Key }} {{ .Type }} {{ .Tag }}
{{ end }}}`

func genPrivateStructTyp(gens []GenSpec) string {
	tmpl, err := template.New("private.struct").Parse(privateStructTmpl)
	if err != nil {
		panic(err)
	}
	var sb strings.Builder
	err = tmpl.Execute(&sb, gens)
	if err != nil {
		panic(err)
	}
	return sb.String()
}

type GenSpec struct {
	Key     string
	Type    string
	Comment string
	Tag     string
}

func ParseJSONObject(obj *Object) []GenSpec {
	var result []GenSpec
	for key, value := range obj.content {
		genSpec := GenSpec{
			Key:     format.Title(key),
			Comment: format.Title(key),
			Tag:     format.Tag(key),
		}

		var typStr string

		typ := whatType(value)
		if typ == MapType {
			object, _ := obj.GetObject(key)
			typStr = genPrivateStructTyp(ParseJSONObject(object))
		} else if typ == ArrayType {
			array, _ := obj.GetArray(key)
			arrTyp, _ := CheckArray(array)
			if arrTyp == MapType {
				object, _ := obj.GetObject(key)
				typStr = genPrivateStructTyp(ParseJSONObject(object))
				typStr = fmt.Sprintf("[]%s", typStr)
			} else {
				typStr = fmt.Sprintf("[]%s", GoTypes[arrTyp])
			}
		} else {
			typStr = GoTypes[typ]
		}

		genSpec.Type = typStr
		result = append(result, genSpec)
	}
	return result
}

func LoadJSONArray(path string) *Array {
	var array = NewArray()
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(file).Decode(&array)
	if err != nil {
		panic(err)
	}
	return array
}

func LoadJSONObject(path string) *Object {
	var obj = NewObject()
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(file).Decode(&obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func ValidJSONFile(path string) int {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	return validJSONObj(file)
}

func validJSONObj(r io.Reader) int {
	reader := bufio.NewReader(r)

	for {
		ru, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				panic("File is empty or does not contain valid JSON.")
			}
			panic(err)
		}

		// 忽略空白字符
		if ru == ' ' || ru == '\t' || ru == '\n' || ru == '\r' {
			continue
		}

		// 检查第一个非空白字符
		if ru == '[' {
			return ArrayType
		} else if ru == '{' {
			return MapType
		} else {
			panic("File does not contain valid JSON.")
		}
	}
}
