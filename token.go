package core

import (
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

type (
	UserClaim struct {
		UserID   uint
		Username string
		Rank     int
		jwt.StandardClaims
	}
)

// BuildJWTToken create jwt token from user info
func BuildJWTToken(user User, signingKey string) (string, error) {
	user.Password = ""
	user.PasswordConfirm = ""

	claim := UserClaim{
		UserID:   user.ID,
		Username: user.Username,
		Rank:     user.Rank,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(signingKey))
}
