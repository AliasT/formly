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
	"github.com/gorilla/mux"
)

// var src = `package mypackage

// type (
// 	myStruct struct{
// 		field1 string
// 		field2 int
// 	}
// )

// `

func FormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("request for model: %v\n", vars["model"])

	fset := token.NewFileSet()
	// Create the AST by parsing src.
	f, err := parser.ParseFile(fset, fmt.Sprintf("./model/%s.go", vars["model"]), nil, 0)
	if err != nil {
		panic(err)
	}

	var fields []FormField

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
	tpl.Execute(w, Payload{Fields: fields})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{model}.tsx", FormHandler)
	http.ListenAndServe(":8080", r)

}
