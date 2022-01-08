// experiment
package main

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// Name of the struct tag used in examples
const (
	ValidateTagName   = "validate"
	SearchableTagName = "searchable"
	OptionsTagName    = "options"
)

type User struct {
	Id    int      `validate:"-",searchable`
	Name  string   `validate:"presence,min=2,max=32",all"`
	Email string   `validate:"email,required"`
	Area  []string `validate:"-" options:"a=1,b=2,c=3"`
}

type FormField struct {
	Name    string
	Kind    string
	Options []Option
}

type Payload struct {
	Fields []FormField
}

type Option struct {
	Label string
	Value string
}

func main() {
	tpl, err := template.ParseFiles("model.html")

	fields := make([]FormField, 0)

	if err != nil {
		panic(err)
	}

	user := User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example",
	}

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(user)

	for i := 0; i < t.NumField(); i++ {
		var kind string
		var options []Option

		field := t.Field(i)

		log.Println(field.Type.Kind())

		// from go type to html type
		switch field.Type.Kind() {
		case reflect.String:
			kind = "text"
			break
		case reflect.Slice:
			kind = "select"
			for _, pair := range strings.Split(field.Tag.Get(OptionsTagName), ",") {
				log.Println(pair)
				option := strings.Split(pair, "=")
				log.Println("option", option)
				options = append(options, Option{Label: option[0], Value: option[1]})
			}
			break
		default:
			kind = "text"
			break
		}

		fields = append(fields, FormField{field.Name, kind, options})
	}

	// Test
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, Payload{fields})
	})

	http.ListenAndServe(":8080", nil)

}
