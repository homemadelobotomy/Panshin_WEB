package ds

import (
	"lab/internal/app/role"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserId      uint
	IsModerator role.Role
}
