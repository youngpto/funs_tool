package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/youngpto/funs_tool/string_utils"
	"os"
	"strconv"
	"strings"
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

var GoTypes = []string{"interface{}", "bool", "int", "float64", "string", "fs_csv.Slice", "fs_csv.Map"}

type Map map[interface{}]interface{}
type Slice []interface{}

type Reader struct {
	Column   int
	Keys     []string
	KeyTypes []int
	Defs     []string
	Comments []string
	Content  [][]string
}

func NewReader(name string) *Reader {
	if strings.HasSuffix(name, ".csv") {
		return NewCsvReader(name)
	}
	panic("unknow file type")
}

func checkAllKeyTypes(keys []string, defs []string, content [][]string, name string) []int {
	ret := make([]int, len(keys))

	for i, key := range keys {
		//if strings.Index(key, "_") > -1 {
		//	continue
		//}
		var typ int
		if i == 0 {
			key = "ID"
			typ = IntType
		} else {
			typ = whatType(defs[i])
		}

		fieldType := typ
		for j := 0; j < len(content); j++ {
			val := content[j][i]
			if len(val) != 0 {
				t := whatType(val)
				if t > fieldType {
					fieldType = t
				}
			}
		}

		if typ != NilType {
			t1, t2 := typ, fieldType
			if t1 != t2 {
				if t1+t2 == IntType+FloatType && t1*t2 == IntType*FloatType {
					typ = FloatType
				} else {
					panic(fmt.Sprintf("%s field %s type %d not same as default %d", name, key, fieldType, typ))
				}
			}
		} else {
			typ = fieldType
		}
		ret[i] = typ
	}
	return ret
}

func ConvType(v string) interface{} {
	typ := whatType(v)
	switch typ {
	case NilType:
		return nil
	case BoolType:
		return boolSet[v]
	case IntType:
		atoi, _ := strconv.Atoi(strings.TrimSpace(v))
		return atoi
	case FloatType:
		float, _ := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return float
	case StringType:
		if isString(v) {
			return v[1 : len(v)-1]
		} else {
			return v
		}
	case ArrayType:
		var array Slice
		arrayStr := v[1 : len(v)-1]
		elems := strings.Split(strings.TrimSpace(arrayStr), ";")
		for _, elem := range elems {
			elemStr := strings.TrimSpace(elem)
			array = append(array, ConvType(elemStr))
		}
		return array
	case MapType:
		var obj = make(Map)
		mapStr := v[1 : len(v)-1]
		elems := strings.Split(strings.TrimSpace(mapStr), ";")
		for _, elem := range elems {
			elemStr := strings.TrimSpace(elem)
			kv := strings.Split(elemStr, "=")
			if len(kv) != 2 {
				panic(fmt.Sprintf("invalid map %s", elemStr))
			}
			obj[ConvType(kv[0])] = ConvType(kv[1])
		}
		return obj
	}
	return nil
}

func whatType(v string) int {
	v = strings.TrimSpace(v)
	if len(v) == 0 {
		return NilType
	} else if isInt(v) {
		return IntType
	} else if isFloat(v) {
		return FloatType
	} else if isBool(v) {
		return BoolType
	} else if isString(v) {
		return StringType
	} else if isArray(v) {
		return ArrayType
	} else if isMap(v) {
		return MapType
	} else if isCSV(v) {
		return StringType
	} else if isJSON(v) {
		return StringType
	}
	return StringType
}

var boolSet = map[string]bool{
	"true":  true,
	"false": false,
	"True":  true,
	"False": false,
	"TRUE":  true,
	"FALSE": false,
}

func isBool(v string) bool {
	_, ok := boolSet[v]
	return ok
}

func isInt(v string) bool {
	_, err := strconv.Atoi(strings.TrimSpace(v))
	return err == nil
}

func isFloat(v string) bool {
	_, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	return err == nil
}

func isString(v string) bool {
	return v[0] == '"' && v[len(v)-1] == '"'
}

func isArray(v string) bool {
	return v[0] == '<' && v[len(v)-1] == '>'
}

func isMap(v string) bool {
	return v[0] == '{' && v[len(v)-1] == '}'
}

func isCSV(v string) bool {
	return strings.HasSuffix(v, ".csv")
}

func isJSON(v string) bool {
	return strings.HasSuffix(v, ".json")
}

func NewCsvReader(name string) *Reader {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("error in file %s\n", name)
			panic(err)
		}
	}()

	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bom := true
	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		panic(err)
	}
	if r != '\uFEFF' {
		br.UnreadRune()
		bom = false
	}

	content, err := csv.NewReader(br).ReadAll()
	if err != nil {
		panic(err)
	}

	valid := make([]int, 0, len(content[0]))
	for i, key := range content[0] {
		if key == "" {
			break
		}
		if key[0] != '_' {
			valid = append(valid, i)
		}
	}

	validFunc := func(col []string) []string {
		values := make([]string, 0, len(valid))
		for _, i := range valid {
			var str string
			if bom {
				str = col[i]
			} else {
				ss, _ := string_utils.GBK2UTF8([]byte(col[i]))
				str = string(ss)
			}
			values = append(values, strings.TrimSpace(str))
		}
		return values
	}

	reader := &Reader{
		Column: len(valid),
	}
	reader.Keys = validFunc(content[0])
	reader.Defs = validFunc(content[1])
	reader.Comments = validFunc(content[2])
	reader.Content = make([][]string, 0, len(content)-3)
	for _, col := range content[3:] {
		if col[0] == "" {
			continue
		}
		reader.Content = append(reader.Content, validFunc(col))
	}
	reader.KeyTypes = checkAllKeyTypes(reader.Keys, reader.Defs, reader.Content, name)

	return reader
}
