// https: //gist.github.com/imantung/60d0c82b8b1641c0aa1c071e1cf77adf
package main

import (
	"fmt"
	. "formly/form"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/structtag"
)

// var src = `package mypackage

// type (
// 	myStruct struct{
// 		field1 string
// 		field2 int
// 	}
// )

// `

func main() {
	// if len(os.Args) < 2 {
	// 	log.Fatalln("expect a model")
	// }

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "model/user.go", nil, 0)
	if err != nil {
		panic(err)
	}

	var fields []FormField

	log.Println(fields)

	for _, node := range f.Decls {
		switch node.(type) {

		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)

					fmt.Printf("Struct: name=%s\n", typeSpec.Name.Name)

					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structType := typeSpec.Type.(*ast.StructType)
						for _, field := range structType.Fields.List {

							switch field.Type.(type) {
							case *ast.Ident:

								// i := field.Type.(*ast.Ident)
								// fieldType := i.Name

								for _, name := range field.Names {
									fields = append(fields, FormField{Name: name.Name, Kind: "text", Options: nil})
									// fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
								}

							case *ast.ArrayType:

								log.Println(field.Tag.Value)
								tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
								if err != nil {
									log.Fatalln(err)
								}

								if optionsTag, err := tags.Get("options"); err == nil {
									var options []Option
									for _, option := range optionsTag.Options {
										pair := strings.Split(option, "=")
										options = append(options, Option{Label: pair[0], Value: pair[1]})
									}
									fields = append(fields, FormField{Name: field.Names[0].Name, Kind: `select`, Options: options})
								}

							}

						}

					}
				}
			}
		}
	}

	tpl, err := template.ParseFiles("model.html")

	if err != nil {
		log.Fatalln("parse template failed", err)
	}

	log.Println(fields)

	// Test
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, Payload{Fields: fields})
	})

	http.ListenAndServe(":8080", nil)

}
