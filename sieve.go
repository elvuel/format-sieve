package formatsieve

// Conventions:
//
// Element.Name must be First Character uppercase
// If Element.Type is not in the basic types, the name must same as Box.Name(casesensitive)
// Nested struct must be declared separately

import (
	"reflect"
)

const (
	boxMissing     = "box missing"
	elementMissing = "element missing"
)

var (
	tString  = reflect.TypeOf("")
	tBool    = reflect.TypeOf(true)
	tInt     = reflect.TypeOf(int(0))
	tInt8    = reflect.TypeOf(int8(0))
	tInt32   = reflect.TypeOf(int32(0))
	tInt64   = reflect.TypeOf(int64(0))
	tUint    = reflect.TypeOf(uint(0))
	tUint8   = reflect.TypeOf(uint8(0))
	tUint16  = reflect.TypeOf(uint16(0))
	tUint32  = reflect.TypeOf(uint32(0))
	tUint64  = reflect.TypeOf(uint64(0))
	tFloat32 = reflect.TypeOf(float32(0.0))
	tFloat64 = reflect.TypeOf(float64(0.0))
	//== SKIPS:
	//Uintptr
	//Complex64
	//Complex128
	//Array
	//Chan
	//Func
	//Interface
	//Map
	//Ptr
	//Slice
	//Structured
	//UnsafePointer
)

// DefaultAssembler is the default assembler
var DefaultAssembler *Assembler

func init() {
	DefaultAssembler = NewAssembler()
}

// Box represents `struct` box
type Box struct {
	Assembler *Assembler
	Name      string
	Elements  []Element
}

// Element represents `structfiled`
type Element struct {
	Name    string
	Tag     string
	Type    string
	IsSlice bool
}

// Valid simply checks the element contains non-empty name and type
func (ele Element) Valid() bool {
	return ele.Name != "" && ele.Type != ""
}

// NewBox returns new box
func NewBox(assembler *Assembler, name string) *Box {
	return &Box{Assembler: assembler, Name: name, Elements: make([]Element, 0)}
}

// AddElement adds an elemnt to the box
func (b *Box) AddElement(e Element) {
	if !e.Valid() {
		return
	}

	found := false
	for _, ele := range b.Elements {
		if ele.Name == e.Name {
			found = true
			break
		}
	}
	if !found {
		b.Elements = append(b.Elements, e)
	}
}

// Structured constructs a `struct` type
func (b *Box) Structured() reflect.Type {
	if len(b.Elements) == 0 {
		panic(elementMissing)
	}

	fields := make([]reflect.StructField, 0)
	for _, ele := range b.Elements {
		field := reflect.StructField{
			Name: ele.Name,
			Tag:  reflect.StructTag(ele.Tag),
		}

		var t reflect.Type

		switch ele.Type {
		case "string":
			t = tString
		case "bool":
			t = tBool
		case "int":
			t = tInt
		case "int8":
			t = tInt8
		case "int32":
			t = tInt32
		case "int64":
			t = tInt64
		case "uint":
			t = tUint
		case "uint8":
			t = tUint8
		case "uint16":
			t = tUint16
		case "uint32":
			t = tUint32
		case "uint64":
			t = tUint64
		case "float32":
			t = tFloat32
		case "float64":
			t = tFloat64
		default:
			t = b.Assembler.Find(ele.Type)
		}

		if ele.IsSlice {
			t = reflect.SliceOf(t)
		}
		field.Type = t

		fields = append(fields, field)
	}

	return reflect.StructOf(fields)
}

// Assembler holds all boxes and corresponding slices
type Assembler struct {
	Boxes  map[string]reflect.Type
	Slices map[string]reflect.Type
}

// NewAssembler returns a new assembler
func NewAssembler() *Assembler {
	return &Assembler{Boxes: make(map[string]reflect.Type), Slices: make(map[string]reflect.Type)}
}

// Add adds box to the assembler
func (a *Assembler) Add(b *Box) {
	b.Assembler = a
	a.Boxes[b.Name] = b.Structured()
	a.Slices[b.Name] = reflect.SliceOf(a.Boxes[b.Name])
}

// Find returns box with the given name
func (a *Assembler) Find(name string) reflect.Type {
	var box reflect.Type
	var e bool
	if box, e = a.Boxes[name]; e == false {
		panic(boxMissing)
	}
	return box
}

// NewSlice returns a new zero value slice of the given name box struct
func (a *Assembler) NewSlice(name string) interface{} {
	return reflect.New(a.Slices[name]).Interface()
}

// NewType returns a new zero value of the given name box struct
func (a *Assembler) NewType(name string) interface{} {
	return reflect.New(a.Boxes[name]).Interface()
}
