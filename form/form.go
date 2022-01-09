package formly

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
