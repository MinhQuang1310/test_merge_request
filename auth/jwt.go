package auth

import (
	"fmt"
	"os"
	"test-echo/common"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserID uint     `json:"ID"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type JwtResetPassword struct {
	UserID uint `json:"ID"`
	jwt.RegisteredClaims
}

func CreateToken(userID uint, roles ...string) (string, error) {
	common.LoadEnvFile()
	sKey := os.Getenv("JWT_KEY")

	// Set custom claims
	claims := &JwtCustomClaims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(sKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func CreateResetPassowrdToken(userID uint) (string, error) {
	common.LoadEnvFile()
	sKey := os.Getenv("JWT_KEY")

	// Set custom claims
	claims := &JwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(sKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ValidateResetPasswordToken(tokenString string) (*JwtResetPassword, error) {

	sKey := os.Getenv("JWT_KEY")

	token, err := jwt.ParseWithClaims(tokenString, &JwtResetPassword{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
		}
		return []byte(sKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtResetPassword)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GetClaimsFromToken(c echo.Context) (*JwtCustomClaims, error) {
	claims := c.Get("user").(*JwtCustomClaims)
	return claims, nil
}
