package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	ss "github.com/elvuel/format-sieve"
	"github.com/jinzhu/gorm"
)

type MyPrototype string
type MyFunctype func(string) string
type MyMaptype map[string]string
type MyProtostruct Fav

type User struct {
	fmt.GoStringer
	*gorm.Model
	A *string
	// MyPrototype
	// MyFunctype
	// MyMaptype
	// MyProtostruct
	Name     string
	LockedAt *time.Time
	Fav      Fav
	// Favs     []Fav
}

type Fav struct {
	Thing string
}

func (fav Fav) GoString() string {
	return "fav"
}

func main() {
	u := new(User)
	u.GoStringer = new(Fav)
	xxx := tt(u)
	u.Model = &gorm.Model{
		ID:        1,
		CreatedAt: time.Now(),
	}
	data, _ := json.Marshal(xxx)
	fmt.Println(string(data))
}

type Infor struct {
	// Type same as reflect.StructField , for Struct hold `struct name`
	Name string
	// PkgPath same as reflect.StructField , for Struct hold `struct package path`
	PkgPath string
	// Kind is the kind of struct for structfield
	Kind string
	// Type same as reflect.StructField
	Type string `json:",omitempty"`
	// Tag same as reflect.StructField
	Tag string `json:",omitempty"`
	// Anonymous same as reflect.StructField
	Anonymous bool // is an embedded field
	Children  map[string]*Infor
}

func tt(val interface{}) *Infor {
	// t := reflect.TypeOf(val)

	// for t.Kind() == reflect.Ptr {
	// 	t = t.Elem()
	// }
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Struct {
		v = reflect.Indirect(v)
	}

	infor := new(Infor)
	infor.Name = v.Type().Name()
	infor.PkgPath = v.Type().PkgPath()
	infor.Kind = v.Kind().String()

	vType := v.Type()
	if v.NumField() > 0 {
		infor.Children = make(map[string]*Infor)
	}

	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		typeField := vType.Field(i)

		// 如果不是嵌入结构体

		if !typeField.Anonymous {

			// n := new(Infor)
			// n.Name = typeField.Name
			// n.PkgPath = valueField.Type().PkgPath()
			// n.Type = valueField.Type().Name()
			// n.Tag = string(typeField.Tag)
			// infor.Children[typeField.Name] = n
			// fmt.Println("value field type:", valueField.Type())
			// fmt.Println("value field kind:", valueField.Kind())
			// fmt.Println("value field type - pkgpath:", valueField.Type().PkgPath())

			// fmt.Println("field name", typeField.Name)
			switch valueField.Kind() {
			case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:

				n := new(Infor)
				n.Name = typeField.Name
				n.PkgPath = valueField.Type().PkgPath()
				n.Kind = valueField.Kind().String()
				n.Type = valueField.Type().Name()
				n.Tag = string(typeField.Tag)
				infor.Children[typeField.Name] = n

			case reflect.Ptr:
				ft := valueField.Type()
				if ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
			case reflect.Interface:
			case reflect.Struct:
			case reflect.Map:
			case reflect.Slice:

			}

		} else {
			// fieldtype 嵌入式指针处理
			ft := valueField.Type()
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}

			fmt.Println(ft.Kind())
			fmt.Println(valueField.Type())
			fmt.Println(ft.Name())
			fmt.Println(ft.PkgPath())
			// fmt.Println(valueField.IsNil())
			// fmt.Println(valueField.IsZero())

			switch ft.Kind() {
			case reflect.Struct: // 结构体
				v := reflect.New(ft).Interface()
				infor.Children[typeField.Name] = tt(v)
			case reflect.Interface: // 接口
				// ignore the zero value interface
				if !valueField.IsZero() && !valueField.IsNil() { // 接口非空处理
					vfCurrentValue := valueField.Interface() // value of Interface
					vfv := reflect.ValueOf(vfCurrentValue)   // 反射value
					if vfv.Kind() == reflect.Ptr {           // 找到值
						vfv = vfv.Elem()
					}
					if vfv.Kind() == reflect.Struct {
						v := reflect.New(vfv.Type()).Interface()
						infor.Children[typeField.Name] = tt(v)
					}
				}
			}

			// vvv := reflect.ValueOf(reflect.New(ft).Interface())

			// infor.Children[typeField.Name] = tt(vvv)
			// infor.Children[typeField.Name] =
		}
		fmt.Println("----------------------")
	}
	return infor
}

func normal() {
	user := &ss.Box{Name: "User", Assembler: ss.DefaultAssembler}
	user.AddElement(ss.Element{Name: "Name", Tag: `json:"name"`, Type: "string"})
	user.AddElement(ss.Element{Name: "Age", Tag: `json:"Age"`, Type: "int"})
	ss.DefaultAssembler.Add(user)

	label := &ss.Box{Name: "Label", Assembler: ss.DefaultAssembler}
	label.AddElement(ss.Element{Name: "Name", Tag: `json:"name"`, Type: "string"})
	label.AddElement(ss.Element{Name: "Popular", Tag: `json:"popular"`, Type: "int"})
	ss.DefaultAssembler.Add(label)

	post := &ss.Box{Name: "Post", Assembler: ss.DefaultAssembler}
	post.AddElement(ss.Element{Name: "Title", Tag: `json:"title"`, Type: "string"})
	post.AddElement(ss.Element{Name: "User", Tag: `json:"user"`, Type: "User"})
	post.AddElement(ss.Element{Name: "Labels", Tag: `json:"labels,omitempty"`, Type: "Label", IsSlice: true})
	ss.DefaultAssembler.Add(post)

	p1 := ss.DefaultAssembler.NewType("Post")

	data := `{
		"title": "tttt",
		"user": {
			"name": "foo",
			"Age": 20,
			"god": "is girl"
		}
	}`

	json.Unmarshal([]byte(data), &p1)
	fmt.Printf("%#v\n", p1)
	fmt.Println("")

	posts := ss.DefaultAssembler.NewSlice("Post")
	data = `[{
		"title": "tttt",
		"user": {
			"name": "foo",
			"Age": 20
		},
		"labels": [{
			"name": "hot",
			"popular": 30
		}, {
			"name": "download",
			"popular": 20
		}]
	}, {
		"title": "tttt",
		"user": {
			"name": "foo",
			"Age": 21
		}
	}]`
	json.Unmarshal([]byte(data), &posts)
	fmt.Printf("%#v\n", posts)
}
