package datapack

import (
	"bufio"
	stdjson "encoding/json"
	"fmt"
	"github.com/youngpto/funs_tool/algorithm"
	"github.com/youngpto/funs_tool/coll/sets/hashset"
	"github.com/youngpto/funs_tool/coll_utils"
	"github.com/youngpto/funs_tool/datapack/csv"
	"github.com/youngpto/funs_tool/datapack/format"
	"github.com/youngpto/funs_tool/datapack/json"
	utils "github.com/youngpto/funs_tool/os"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

var allowFileType = []string{
	".json",
	//".yaml",
	".csv",
}

type inode struct {
	name  string
	isdir bool
	ext   string
	path  string
	nodes []*inode
	prev  *inode

	structname  string
	variatename string

	isJsonArray   bool
	jsonArrayType int

	binary *json.Object
	sync.Mutex
}

func (i *inode) gen() (structSpec, bool) {
	spec := structSpec{
		Name:    i.structname,
		VName:   i.variatename,
		Comment: prettycomment(i.path),
	}
	if i.isdir {
		if i.prev != nil {
			i.prev.binary.Set(i.name, i.binary)
		}

		spec.Fields = make([]fieldSpec, 0, len(i.nodes))
		mergeJson := make(map[string]struct{})
		for _, nn := range i.nodes {
			typ := fmt.Sprintf("*%s", nn.structname)
			if !nn.isdir {
				switch nn.ext {
				case ".csv":
					typ = fmt.Sprintf("map[int]%s", typ)
				case ".json":
					idx := strings.LastIndex(typ, "_")
					if idx >= 0 && idx+1 < len(typ) {
						_, err := strconv.ParseInt(typ[idx+1:], 10, 64)
						if err == nil {
							typ = typ[:idx]
						}
					}

					if coll_utils.Inmap(typ, mergeJson) {
						continue
					}

					mergeJson[typ] = struct{}{}

					jsonType := json.ValidJSONFile(nn.path)
					if jsonType == json.ArrayType {
						jsonArray := json.LoadJSONArray(nn.path)
						eleTyp, _ := json.CheckArray(jsonArray)
						if eleTyp != json.MapType {
							typ = fmt.Sprintf("[]%s", json.GoTypes[eleTyp])
						} else {
							typ = fmt.Sprintf("[]%s", typ)
						}
						nn.isJsonArray = true
						nn.jsonArrayType = eleTyp
					}

					typ = fmt.Sprintf("map[string]%s", typ)
				}
			}
			spec.Fields = append(spec.Fields, fieldSpec{
				Name:    nn.variatename,
				Type:    typ,
				Comment: prettycomment(nn.path),
				Tag:     format.Tag(nn.name),
			})
		}
	} else {
		switch i.ext {
		case ".json":
			var structName = i.name
			var jsonObj interface{}
			if i.isJsonArray {
				jsonObj = json.LoadJSONArray(i.path)
			} else {
				jsonObj = json.LoadJSONObject(i.path)
			}
			defer func() {
				if i.prev != nil {
					obj := i.prev.binary.SetDefault(structName, json.NewObject()).(*json.Object)
					obj.Set(i.name, jsonObj)
				}
			}()

			idx := strings.LastIndex(i.structname, "_")
			if idx >= 0 && idx+1 < len(i.structname) {
				_, err := strconv.ParseInt(i.structname[idx+1:], 10, 64)
				if err == nil {
					structName = i.name[:idx]
					return spec, false
				}
			}

			if i.isJsonArray && i.jsonArrayType != json.MapType {
				return spec, false
			}
		}
		spec.Fields = i.parseFields()
	}

	return spec, true
}

func (i *inode) parseFields() []fieldSpec {
	switch i.ext {
	case ".csv":
		return i.parseCSVFields()
	case ".json":
		return i.parseJSONFields()
	default:
		var specs = []fieldSpec{
			{
				Name:    "F1",
				Type:    "int",
				Comment: "is f1",
				Tag:     format.Tag("F1"),
			},
			{
				Name:    "F2",
				Type:    "string",
				Comment: "is f2",
				Tag:     format.Tag("F2"),
			},
		}
		return specs
	}
}

func (i *inode) parseJSONFields() []fieldSpec {
	var specs []fieldSpec

	var obj = json.NewObject()
	if i.isJsonArray {
		jsonArray := json.LoadJSONArray(i.path)
		obj, _ = jsonArray.GetObject(0)
	} else {
		obj = json.LoadJSONObject(i.path)
	}
	gens := json.ParseJSONObject(obj)
	for _, gen := range gens {
		specs = append(specs, fieldSpec{
			Name:    gen.Key,
			Type:    gen.Type,
			Comment: gen.Comment,
			Tag:     gen.Tag,
		})
	}
	return specs
}

func (i *inode) parseCSVFields() []fieldSpec {
	var specs []fieldSpec
	reader := csv.NewReader(i.path)
	keys := reader.Keys
	keyTypes := reader.KeyTypes
	specs = make([]fieldSpec, 0, len(keys))
	for j, key := range keys {
		//if strings.Index(key, "_") > -1 {
		//	continue
		//}
		var typ int
		if j == 0 {
			key = "ID"
			typ = csv.IntType
		} else {
			typ = keyTypes[j]
		}

		spec := fieldSpec{
			Name:    format.Title(key),
			Type:    csv.GoTypes[typ],
			Comment: reader.Comments[j],
			Tag:     format.Tag(key),
		}
		specs = append(specs, spec)
	}

	if i.prev != nil {
		csvMap := make(map[int]*json.Object)
		for _, values := range reader.Content {
			record := json.NewObject()
			for idx, value := range values {
				if idx == 0 {
					id, _ := strconv.Atoi(value)
					record.Set("ID", id)
				} else {
					conVal := csv.ConvType(value)
					if conVal == nil {
						conVal = csv.ConvType(reader.Defs[idx])
					}
					record.Set(keys[idx], conVal)
				}
			}
			id, _ := record.GetInt("ID")
			csvMap[id] = record
		}
		i.prev.binary.Set(i.name, csvMap)
	}
	return specs
}

/*
var rootPath = "./bin/game_conf/"
var genFile = "./conf/generated.go"
var msgpackFile = "./conf/msgpack.json"
*/

func Conf2Src(rootPath string, genFile string, msgpackFile string) error {
	prefix = filepath.ToSlash(filepath.Join(rootPath, ""))
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		fmt.Printf("cost time %ds:\n", end-start)
	}()

	root := &inode{
		name:       "config",
		isdir:      true,
		path:       rootPath,
		structname: "gameConfig",
		binary:     json.NewObject(),
	}
	exist := hashset.New[string]()
	visit(rootPath, root, exist)
	conf2go(root, genFile)
	_ = utils.SysRun(os.Stdout, os.Stderr, "gofmt", "-l", "-w", "-e", genFile)

	bytes, err := stdjson.Marshal(root.binary)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(msgpackFile, bytes, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}

func visit(path string, parent *inode, exist *hashset.Set[string]) {
	files, _ := ioutil.ReadDir(path)
	parent.nodes = make([]*inode, 0, len(files))

	for _, file := range files {
		name := file.Name()
		if name[0] == '.' {
			continue
		}

		if !file.IsDir() {
			if !coll_utils.In(filepath.Ext(name), allowFileType) {
				fmt.Printf("not register file type from %s \n", name)
				continue
			}
		}

		fpath := filepath.Join(path, name)
		fmt.Println("visit", fpath)

		n := strings.Split(filepath.Base(name), ".")[0]
		subpath := strings.Split(fpath, ".")[0]
		subpath = prettycomment(subpath)
		structname := strings.Join(strings.Split(subpath, "/"), "_")
		structname = format.Title(structname)
		if exist.Contains(structname) {
			panic(fmt.Sprintf("filename %s is exist \n", structname))
		}
		exist.Add(structname)
		node := &inode{
			name:        n,
			isdir:       file.IsDir(),
			path:        fpath,
			structname:  structname,
			variatename: format.Title(n),
			binary:      json.NewObject(),
		}
		if !file.IsDir() {
			node.ext = filepath.Ext(name)
		}
		node.prev = parent
		parent.nodes = append(parent.nodes, node)

		if file.IsDir() {
			visit(fpath, node, exist)
		}
	}
}

var prefix string

func prettycomment(s string) string {
	comment := filepath.ToSlash(s)
	if strings.HasPrefix(comment, prefix) {
		idx := len(prefix)
		if idx >= len(comment) {
			comment = ""
		} else {
			if comment[idx] == '/' {
				idx += 1
			}
			comment = comment[idx:]
		}
	}
	return comment
}

var goTmpl = `// Code generated - DO NOT EDIT.
package conf

import (
	fs_csv "github.com/youngpto/funs_tool/datapack/csv"
	fs_json "github.com/youngpto/funs_tool/datapack/json"
)

var assertCsv fs_csv.Reader
var assertJson fs_json.Object

{{range .}}// {{ .Comment }}
type {{ .Name }} struct {
{{ range .Fields }}	// {{ .Comment }}
	{{ .Name }} {{ .Type }} {{ .Tag }}
{{ end }}}

{{ end }}
`

type structSpec struct {
	Name    string
	VName   string
	Comment string

	Fields []fieldSpec
}

type fieldSpec struct {
	Name    string
	Type    string
	Comment string
	Tag     string
}

func conf2go(root *inode, out string) {
	var structSpecs = make([]structSpec, 0)
	algorithm.DFS(root, func(pop *inode) []*inode {
		if gen, ok := pop.gen(); ok {
			structSpecs = append(structSpecs, gen)
		}
		return pop.nodes
	})

	write2File(out, func(writer io.Writer) {
		tmpl, err := template.New("go.tmpl").Parse(goTmpl)
		if err != nil {
			panic(err)
		}

		bw := bufio.NewWriter(writer)
		err = tmpl.Execute(bw, structSpecs)
		if err != nil {
			bw.Flush()
			panic(err)
		}
		err = bw.Flush()
		if err != nil {
			panic(err)
		}
	})
}

func write2File(name string, do func(writer io.Writer)) {
	os.Remove(name)
	file, err := os.Create(name)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	do(file)
}
