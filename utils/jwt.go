package utils

import (
	"os"
	"time"

	"bepuas/app/model"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "change_this_secret_in_production_32chars"
	}
	jwtSecret = []byte(s)
}

type JWTClaim struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`        // MAHASISWA, ADMIN, DOSEN
	Permissions []string `json:"permissions"` // achievement:create, dll
	jwt.RegisteredClaims
}

func ValidateToken(tokenStr string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func GenerateToken(user model.User, perms []string) (string, error) {
	claims := model.JWTClaims{
		UserID:      user.ID,
		Username:    user.Username,
		RoleID:      user.RoleID,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*model.JWTClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tok.Claims.(*model.JWTClaims); ok && tok.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
