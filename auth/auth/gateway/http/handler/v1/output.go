package v1

import "github.com/regismelgaco/go-sdks/auth/auth/entity"

type UserOutput struct {
	UserName string `json:"username"`
}

func ToUserOutput(u entity.User) UserOutput {
	return UserOutput{UserName: u.UserName}
}

type LoginOutput struct {
	Token string `json:"token"`
}

func ToLoginOutput(t entity.Token) LoginOutput {
	return LoginOutput{Token: string(t)}
}
