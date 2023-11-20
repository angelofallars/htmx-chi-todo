// jwt.go provides files for working with JWT
package auth

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	// Reference to user.ID
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Attempt to extract JwtClaims from request,
// otherwise return an error
func JwtClaimsFromRequest(r *http.Request) (*JwtClaims, error) {
	_, claimsMap, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil, errors.New("error fetching JWT from context")
	}

	id, ok := claimsMap["id"].(string)
	if !ok {
		return nil, errors.New("error fetching ID from JWT")
	}
	username, ok := claimsMap["username"].(string)
	if !ok {
		return nil, errors.New("error fetching username from JWT")
	}

	claims := &JwtClaims{
		ID:       id,
		Username: username,
	}

	return claims, nil
}
