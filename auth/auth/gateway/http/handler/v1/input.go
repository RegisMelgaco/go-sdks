package v1

import (
	"github.com/regismelgaco/go-sdks/auth/auth/usecase"
)

type UserInput struct {
	UserName string `json:"username"`
	Pass     string `json:"pass"`
}

func (i UserInput) ToEntity() usecase.CreateUserInput {
	return usecase.CreateUserInput{
		UserName: i.UserName,
		Pass:     i.Pass,
	}
}

func (i UserInput) ToLoginEntity() usecase.LoginInput {
	return usecase.LoginInput{
		UserName: i.UserName,
		Pass:     i.Pass,
	}
}
