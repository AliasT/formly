package model

// User
type User struct {
	Id    int      `validate:"-"`
	Name  string   `validate:"presence,min=2,max=32"`
	Email string   `validate:"email,required"`
	Area  []string `validate:"-" options:",a=1,b=2,c=3"`
}
