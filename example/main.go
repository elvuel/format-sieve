package main

import (
	"encoding/json"
	"fmt"

	ss "github.com/elvuel/format-sieve"
)

func main() {

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
