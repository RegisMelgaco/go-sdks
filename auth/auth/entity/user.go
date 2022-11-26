package entity

type User struct {
	ID int

	UserName string
	Secret   []byte
}
