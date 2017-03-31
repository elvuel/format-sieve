// +build ignore
package format_sieve

import (
	"testing"
	"reflect"
	"encoding/json"
)

var benchSeed =`{"Foo":{"Name":"FooName","Tags":["FooTagA,FooTagB"]},"Foos":[{"Name":"FooName","Tags":["FooTagA,FooTagB"]},{"Name":"FooName1","Tags":["Foo1TagA,Foo1TagB"]}],"Name":"BarName"}`

func BenchmarkNormal(b *testing.B) {
	type Foo struct {
		Name string
		Tags []string
	}
	type Bar struct {
		Foo Foo
		Foos []Foo
		Name string
	}

	//var data []byte

	for n := 0; n < b.N; n++ {
		bar := Bar{Foos: make([]Foo, 0)}
		json.Unmarshal([]byte(benchSeed), &bar)
		//data ,_ = json.Marshal(&bar)
		//if len(data) != len(benchSeed) {
		//	fmt.Println("Breaked")
		//	break
		//}

	}
}

func BenchmarkSieve(b *testing.B) {
	assembler := NewAssembler()
	box := NewBox(nil, "Foo")
	box.AddElement(Element{Name: "Name", Type: "string"})
	box.AddElement(Element{Name: "Tags", Type: "string",IsSlice:true})
	assembler.Add(box)

	box1 := NewBox(nil, "Bar")
	box1.AddElement(Element{Name: "Name", Type: "string"})
	box1.AddElement(Element{Name: "Foo", Type: "Foo"})
	box1.AddElement(Element{Name: "Foos", Type: "Foo", IsSlice:true})
	assembler.Add(box1)

	//var data []byte

	for n := 0; n < b.N; n++ {
		bar := assembler.NewType("Bar")
		json.Unmarshal([]byte(benchSeed), &bar)
		//data ,_ = json.Marshal(&bar)
		//if len(data) != len(benchSeed) {
		//	fmt.Println("Breaked")
		//	break
		//}
	}
}

func basicTypes() []string{
	return []string{
		"string", "bool",
		"int", "int8", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
	}
}

func TestElement_Valid(t *testing.T) {
	ele := Element{}
	if ele.Valid() == true {
		t.Error("Should not be valid")
	}

	ele.Name = "foo"
	if ele.Valid() == true {
		t.Error("Should not be valid")
	}

	ele.Name = ""
	ele.Type = "bar"
	if ele.Valid() == true {
		t.Error("Should not be valid")
	}

	ele.Name = "foo"
	ele.Type = "bar"
	if !ele.Valid() {
		t.Error("Should be valid")
	}
}

func TestNewAssembler(t *testing.T) {
	assembler := NewAssembler()
	if len(assembler.Boxes) != 0 {
		t.Error("Assembler should be initilized with zero box")
	}

	if len(assembler.Slices) != 0 {
		t.Error("Assembler should be initilized with zero slice")
	}
}

func TestAssembler_Add(t *testing.T) {
	box := NewBox(nil, "Foo")
	box.AddElement(Element{Name: "Foo", Tag: "bar", Type: "string"})

	if box.Assembler != nil {
		t.Error("The assembler should be nil")
	}

	assembler := NewAssembler()

	assembler.Add(box)

	if box.Assembler == nil {
		t.Error("The assembler should not be nil")
	}

	if len(assembler.Boxes) != 1 {
		t.Error("Boxes length should be 1")
	}

	if len(assembler.Slices) != 1 {
		t.Error("Slices length should be 1")
	}

	i := reflect.New(assembler.Boxes["Foo"]).Interface()
	if k := reflect.ValueOf(i).Elem().Field(0).Type().Kind().String(); k != "string" {
		t.Errorf("The type of field wants string, but %s", k)
	}

	inst := reflect.New(assembler.Slices["Foo"]).Interface()
	instEle := reflect.ValueOf(inst).Elem()
	if k := instEle.Type().Kind().String(); k != "slice" {
		t.Error("The type of field should be slice")
	}

	// slice of what
	// instEle.Type().Elem() -> struct{...}
	if sow := instEle.Type().Elem().Field(0).Type.Name(); sow != "string" {
		t.Errorf("The type of field is slice of  -> wants string, but %s", sow)
	}
}

func TestAssembler_Find(t *testing.T) {
	missAssemblerFunc := func(assembler *Assembler, box string) (result string) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					result = x
				case error:
					result = x.Error()
				default:
					result = "w>>>"
				}
			}
		}()
		assembler.Find(box)
		result = "succeed"
		return
	}

	assembler := NewAssembler()
	fooo := missAssemblerFunc(assembler, "foobar")
	if fooo != BoxMissing {
		t.Error("Should get box missing message")
	}

	box := NewBox(nil, "Foo")
	box.AddElement(Element{Name:"FooBar", Type:"string"})
	assembler.Add(box)
	fooo = missAssemblerFunc(assembler, "Foo")
	if fooo != "succeed" {
		t.Error("Should get succeed message")
	}
}

func TestAssembler_NewType(t *testing.T) {
	assembler := NewAssembler()
	box := NewBox(nil, "Foo")
	box.AddElement(Element{Name:"Title", Type:"string"})
	assembler.Add(box)

	v := assembler.NewType("Foo")
	data, _ := json.Marshal(&v)
	if v := string(data); v != `{"Title":""}` {
		t.Errorf("Wrong Type, %s", v)
	}
}

func TestAssembler_NewSlice(t *testing.T) {
	assembler := NewAssembler()
	box := NewBox(nil, "Foo")
	box.AddElement(Element{Name:"Tags", Tag: "tags", Type:"string"})
	assembler.Add(box)

	foos := assembler.NewSlice("Foo")

	seed := `[{"tags":"foo"},{"tags":"bar"}]`
	err := json.Unmarshal([]byte(seed), &foos)
	if err != nil {
		t.Error("Should not raise error but", err.Error())
	}

	instEle := reflect.ValueOf(foos).Elem()
	if k := instEle.Type().Kind().String(); k != "slice" {
		t.Error("The type of field should be slice")
	}

	if sow := instEle.Type().Elem().Field(0).Type.Name(); sow != "string" {
		t.Errorf("The type of field is slice of  -> wants string, but %s", sow)
	}

	if instEle.Len() != 2 {
		t.Error("Slice length should be 2")
	}

	//for i:=0; i< instEle.Len();i++ {
	//	t.Error(instEle.Index(i).Interface())
	//}
	//if instEle.Index(0).Field(0).String() != "foo" {
	//	t.Error("Wrong")
	//}
	//
	//if instEle.Index(1).Field(0).String() != "bar" {
	//	t.Error("Wrong")
	//}
	vSeeds := []string{"foo", "bar"}
	for i:=0; i< instEle.Len();i++ {
		if instEle.Index(i).Field(0).String() != vSeeds[i] {
			t.Error("Wrong")
		}
	}
}

func TestNewBox(t *testing.T) {
	box := NewBox(nil, "Foo")

	if box.Assembler != nil {
		t.Error("The assembler should be nil")
	}

	if box.Name != "Foo" {
		t.Error("The assembler should be 'Foo'")
	}

	box.Assembler = DefaultAssembler

	if box.Assembler == nil {
		t.Error("The assembler should not be nil")
	}
}

func TestBox_AddElement(t *testing.T) {
	ele := Element{Name: "Foo", Tag: "bar", Type: "string", IsSlice: true}
	box := NewBox(nil, "Foo")
	box.AddElement(ele)

	if len(box.Elements) != 1 {
		t.Error("The elements length should be 1")
	}

	if !box.Elements[0].IsSlice {
		t.Error("The element should be slice")
	}

	eleBar := Element{Name:"Bar", Tag:"foo", Type: "int"}
	box.AddElement(eleBar)
	if len(box.Elements) != 2 {
		t.Error("The elements length should be 2")
	}

	eleDup := Element{Name: "Foo", Tag: "bar", Type: "string"}
	box.AddElement(eleDup)
	if len(box.Elements) != 2 {
		t.Error("The elements length should be 2")
	}
}

func TestBox_Structured(t *testing.T) {
	// element missing
	missBoxFunc := func() (result string) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					result = x
				case error:
					result = x.Error()
				default:
					result = "w>>>"
				}
			}
		}()
		box := NewBox(nil, "Doo")
		box.Structured()
		result = "succeed"
		return
	}

	fooo := missBoxFunc()
	if fooo != ElementMissing {
		t.Error("Should got element missing panic message")
	}

	// == Basic Types
	// 	Validate non-slice
	for _, typ := range basicTypes() {
		box := NewBox(nil, "Doo")
		ele := Element{Name: "Foo", Tag: "bar", Type: typ}
		box.AddElement(ele)
		doo := box.Structured()

		inst := reflect.New(doo).Interface()
		instEle := reflect.ValueOf(inst).Elem()
		if instEle.NumField() !=1 {
			t.Error("The struct length of fields should be 1")
		}

		field := instEle.Type().Field(0)
		if field.Name != "Foo" {
			t.Error("The name of field should be Foo")
		}

		if k := field.Type.Kind().String(); k != typ {
			t.Errorf("The type of field wants %s, but %s", typ, k)
		}
	}

	// 	validate slice validate
	for _, typ := range basicTypes() {
		box := NewBox(nil, "Doo")
		ele := Element{Name: "Foo", Tag: "bar", Type: typ, IsSlice: true}
		box.AddElement(ele)
		doo := box.Structured()

		inst := reflect.New(doo).Interface()
		instEle := reflect.ValueOf(inst).Elem()
		if instEle.NumField() !=1 {
			t.Error("The struct length of fields should be 1")
		}

		field := instEle.Type().Field(0)
		if field.Name != "Foo" {
			t.Error("The name of field should be Foo")
		}

		if k := field.Type.Kind().String(); k != "slice" {
			t.Error("The type of field should be slice")
		}


		if sow := field.Type.Elem().Name(); sow != typ {
			t.Errorf("The type of field is slice of  -> wants %s, but %s", typ, sow)
		}
	}

	// == Non Basic Types
	//	Non box missing
	assembler := NewAssembler()
	boxBar := NewBox(assembler, "Bar")
	boxBar.AddElement(Element{Name: "B3a", Tag: "b3a", Type: "string"})
	assembler.Add(boxBar)

	boxFoo := NewBox(assembler, "Foo")
	boxFoo.AddElement(Element{Name: "F2a", Tag: "f2a", Type: "string"})
	boxFoo.AddElement(Element{Name: "Baref", Tag: "baref", Type: "Bar"})
	assembler.Add(boxFoo)

	Foo := assembler.Find("Foo")
	foo := reflect.New(Foo).Interface()
	fooEle := reflect.ValueOf(foo).Elem()
	sField, existed := fooEle.Type().FieldByName("Baref")

	if !existed {
		t.Error("Baref filed should existed")
	}

	if sField.Tag != "baref" {
		t.Error("Field should be 'baref'")
	}

	Bar := sField.Type
	bar := reflect.New(Bar).Interface()
	barEle := reflect.ValueOf(bar).Elem()
	if barEle.Type().Field(0).Name != "B3a" {
		t.Error("Bar field name should be 'B3a'")
	}

	if barEle.Type().Field(0).Tag != "b3a" {
		t.Error("Bar field tag should be 'b3a'")
	}

	if barEle.Type().Field(0).Type.Kind().String() != "string" {
		t.Error("Bar field type should be 'string'")
	}

	// repeat test
	missFunc := func() (result string) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					result = x
				case error:
					result = x.Error()
				default:
					result = "w>>>"
				}
			}
		}()
		box :=  NewBox(assembler, "Missssssss")
		box.AddElement(Element{Name: "F2a", Tag: "f2a", Type: "FOOOOOBAAAAAA"})
		assembler.Add(box)
		result = "succeed"
		return
	}

	fooo = missFunc()
	if fooo != BoxMissing {
		t.Error("Should get box missing message")
	}
}
