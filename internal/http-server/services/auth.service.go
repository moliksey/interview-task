package services

import (
	"awesomeProject/internal/repo"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

//func Auth(guid string, ip string) (access string, refresh string, err error) {

//}

type AuthService struct {
	repo repo.Authorization
}

func NewAuthService(repo repo.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
func (r *AuthService) AuthUser(guid string, ip string) (access string, refresh string, err error) {
	t, rt, err := GenerateTokenPair(guid, ip)
	if err = r.repo.AuthUser(guid, rt); err != nil {
		return "", "", err
	}
	return t, rt, nil
}
func GenerateTokenPair(guid string, ip string) (access string, refresh string, err error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if len(secretKey) == 0 {
		return "", "", fmt.Errorf("SECRET_KEY is not set")
	}
	accessToken := createAccessToken(guid, ip)
	t, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	refreshToken := createRefreshToken(ip)
	rt, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	return t, rt, nil
}
func createAccessToken(guid string, ip string) *jwt.Token {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": guid,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"ip":  ip,
	})
	fmt.Printf("Token claims added: %+v\n", claims)
	return claims
}
func createRefreshToken(ip string) *jwt.Token {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
		"ip":  ip,
	})
	fmt.Printf("Token refresh claims added: %+v\n", claims)
	return claims
}
