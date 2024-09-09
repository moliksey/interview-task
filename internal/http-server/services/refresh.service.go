package services

import (
	"awesomeProject/internal/repo"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

type RefreshService struct {
	repo repo.Refresh
}

func NewRefreshService(repo repo.Refresh) *RefreshService {
	return &RefreshService{repo: repo}
}
func (r *RefreshService) RefreshToken(accessOld string, refreshOld string, ip string) (access string, refresh string, err error) {

	token, err := ParseToken(accessOld)
	if err != nil {
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		clGuid := claims["sub"].(string)
		email, dbRefresh, errBd := r.repo.GetUser(clGuid)
		if errBd != nil {
			return "", "", errBd
		}
		clIp := claims["ip"].(string)
		clIat := strconv.FormatFloat(claims["iat"].(float64), 'f', -1, 64)
		if int64(claims["exp"].(float64)) != time.Unix(int64(claims["iat"].(float64)), 0).Add(time.Hour).Unix() {
			return "", "", fmt.Errorf("differnt token time")
		}
		newAccess, newRefresh, errBd := updateTokens(refreshOld, dbRefresh, email, clIp, clGuid, clIat, ip)
		if errBd != nil {
			logrus.Error(errBd.Error())
			return "", "", errBd
		}
		errBd = r.repo.UpdateRefreshToken(clGuid, newRefresh)
		if errBd != nil {
			return "", "", errBd
		}
		return newAccess, newRefresh, nil
	}
	return "", "", fmt.Errorf("error taking claims")
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("SECRET_KEY is not set")
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func updateTokens(refreshOld string, dbRefresh string, email string, ip string, guid string, iat string, realIp string) (access string, refresh string, err error) {
	if dbRefresh != refreshOld {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	oldRefToken, err := ParseToken(refreshOld)
	if err != nil {
		return "", "", err
	}
	oldIp := ""
	if oldClaims, bdOk := oldRefToken.Claims.(jwt.MapClaims); bdOk && oldRefToken.Valid {
		if strconv.FormatFloat(oldClaims["iat"].(float64), 'f', -1, 64) != iat {
			return "", "", fmt.Errorf("differnt token time")
		}
		if int64(oldClaims["exp"].(float64)) != time.Unix(int64(oldClaims["iat"].(float64)), 0).Add(time.Hour*24).Unix() {
			return "", "", fmt.Errorf("differnt token time")
		}
		oldIp = oldClaims["ip"].(string)
	} else {
		return "", "", fmt.Errorf("invalid token")
	}
	dbToken, err := ParseToken(dbRefresh)
	if err != nil {
		return "", "", err
	}
	if dbClaims, bdOk := dbToken.Claims.(jwt.MapClaims); bdOk && dbToken.Valid {
		if realIp != dbClaims["ip"].(string) || realIp != ip || oldIp != ip {
			if badEmail := SendMail(email); badEmail != nil {
				return "", "", fmt.Errorf("error sending email")
			}
			return "", "", fmt.Errorf("invalid ip")
		}
	} else {
		return "", "", fmt.Errorf("invalid token")
	}
	newAccess, newRefresh, err := GenerateTokenPair(guid, ip)
	if err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}

func SendMail(email string) error {
	pass := os.Getenv("EMAIL_PASS")
	if pass == "" {
		return fmt.Errorf("EMAIL_PASS is not set")
	}
	auth := smtp.PlainAuth(
		"",
		"moliksey@gmail.com",
		pass,
		"smtp.gmail.com",
	)
	message := "Subject: Your account is under threat\n Someone tried to log into your account from another ip."
	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"moliksey@gmail.com",
		[]string{email},
		[]byte(message),
	)

}
