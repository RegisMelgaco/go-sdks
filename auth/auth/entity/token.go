package entity

import "time"

type Token string

type TokenClaims struct {
	UserName   string
	Expiration time.Time
}
