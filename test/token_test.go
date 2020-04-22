package test

import (
	"testing"

	"github.com/ledinhbao/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

func TestBuildJWTToken(t *testing.T) {
	user := core.User{
		Username: "ledinhbao",
		Password: "passwrod",
		Rank:     core.RankSuperAdmin,
		ID:       10,
	}
	userClaim := core.UserClaim{
		UserID:   user.ID,
		Username: user.Username,
		Rank:     user.Rank,
	}
	signingKey := "Test*Signing#Key#For@JWT"
	signature, err := core.BuildJWTToken(user, signingKey)
	require.Nil(t, err)

	var claim core.UserClaim
	_, err = jwt.ParseWithClaims(signature, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	assert.Equal(t, claim.UserID, userClaim.UserID)
	assert.Equal(t, claim.Username, userClaim.Username)
	assert.Equal(t, claim.Rank, userClaim.Rank)
}
