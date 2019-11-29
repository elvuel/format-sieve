package formatsieve

// Conventions:
//
// Element.Name must be First Character uppercase
// If Element.Type is not in basic types, The name must same as Box.Name(casesensitive)
// Nested struct must be declared separately

import (
	"reflect"
)

const (
	BoxMissing     = "box missing"
	ElementMissing = "element missing"
)

var (
	TString  = reflect.TypeOf("")
	TBool    = reflect.TypeOf(true)
	TInt     = reflect.TypeOf(int(0))
	TInt8    = reflect.TypeOf(int8(0))
	TInt32   = reflect.TypeOf(int32(0))
	TInt64   = reflect.TypeOf(int64(0))
	TUint    = reflect.TypeOf(uint(0))
	TUint8   = reflect.TypeOf(uint8(0))
	TUint16  = reflect.TypeOf(uint16(0))
	TUint32  = reflect.TypeOf(uint32(0))
	TUint64  = reflect.TypeOf(uint64(0))
	TFloat32 = reflect.TypeOf(float32(0.0))
	TFloat64 = reflect.TypeOf(float64(0.0))
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

var DefaultAssembler *Assembler

func init() {
	DefaultAssembler = NewAssembler()
}

type Box struct {
	Assembler *Assembler
	Name      string
	Elements  []Element
}

func NewBox(assembler *Assembler, name string) *Box {
	return &Box{Assembler: assembler, Name: name, Elements: make([]Element, 0)}
}

type Element struct {
	Name    string
	Tag     string
	Type    string
	IsSlice bool
}

func (ele Element) Valid() bool {
	return ele.Name != "" && ele.Type != ""
}

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

func (b *Box) Structured() reflect.Type {
	if len(b.Elements) == 0 {
		panic(ElementMissing)
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
			t = TString
		case "bool":
			t = TBool
		case "int":
			t = TInt
		case "int8":
			t = TInt8
		case "int32":
			t = TInt32
		case "int64":
			t = TInt64
		case "uint":
			t = TUint
		case "uint8":
			t = TUint8
		case "uint16":
			t = TUint16
		case "uint32":
			t = TUint32
		case "uint64":
			t = TUint64
		case "float32":
			t = TFloat32
		case "float64":
			t = TFloat64
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

type Assembler struct {
	Boxes  map[string]reflect.Type
	Slices map[string]reflect.Type
}

func NewAssembler() *Assembler {
	return &Assembler{Boxes: make(map[string]reflect.Type), Slices: make(map[string]reflect.Type)}
}

func (a *Assembler) Add(b *Box) {
	b.Assembler = a
	a.Boxes[b.Name] = b.Structured()
	a.Slices[b.Name] = reflect.SliceOf(a.Boxes[b.Name])

}

func (a *Assembler) Find(name string) reflect.Type {
	if _, e := a.Boxes[name]; e == false {
		panic(BoxMissing)
	}
	return a.Boxes[name]
}

func (a *Assembler) NewSlice(name string) interface{} {
	return reflect.New(a.Slices[name]).Interface()
}

func (a *Assembler) NewType(name string) interface{} {
	return reflect.New(a.Boxes[name]).Interface()
}
