package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"testing"
)

func fileFilter(file fs.FileInfo) bool {
	return true
}

func Test(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse src but stop after processing the imports.
	f, err := parser.ParseDir(fset, "./model", fileFilter, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the imports from the file's AST.
	for _, s := range f {
		t.Log(s)
	}

}
