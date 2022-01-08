// https: //gist.github.com/imantung/60d0c82b8b1641c0aa1c071e1cf77adf
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var src = `package mypackage

type (
	myStruct struct{
		field1 string
		field2 int
	}
)

`

func main() {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "model/user.go", nil, 0)
	if err != nil {
		panic(err)
	}

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
							i := field.Type.(*ast.Ident)
							fieldType := i.Name

							for _, name := range field.Names {
								fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
							}

						}

					}
				}
			}
		}
	}
}
