package tests

import (
	"awesomeProject/internal/http-server/handlers/ref"
	"awesomeProject/internal/http-server/handlers/reg"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/url"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8080"
)

func TestHappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/api/v1/reg").
		WithJSON(reg.Request{
			Email: gofakeit.Email(),
		}).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("guid").NotContainsKey("error")
}

//nolint:funlen
func TestUseCase(t *testing.T) {
	testCases := []struct {
		name         string
		email        string
		error        bool `default:"false"`
		fakeGuid     bool `default:"false"`
		fakeAccess   bool `default:"false"`
		fakeRefresh  bool `default:"false"`
		emptyGuid    bool `default:"false"`
		emptyAccess  bool `default:"false"`
		emptyRefresh bool `default:"false"`
	}{
		{
			name:  "Valid email",
			email: gofakeit.Email(),
		},
		{
			name:  "Invalid email",
			email: "iasasas",
			error: true,
		},
		{
			name:  "Empty email",
			email: "",
			error: true,
		},
		{
			name:     "Fake guid",
			email:    gofakeit.Email(),
			fakeGuid: true,
		},
		{
			name:       "Fake access",
			email:      gofakeit.Email(),
			fakeAccess: true,
		},
		{
			name:        "Fake refresh",
			email:       gofakeit.Email(),
			fakeRefresh: true,
		},
		{
			name:      "Empty guid",
			email:     gofakeit.Email(),
			emptyGuid: true,
		},
		{
			name:         "Empty access",
			email:        gofakeit.Email(),
			emptyRefresh: true,
		},
		{
			name:        "Empty refresh",
			email:       gofakeit.Email(),
			emptyAccess: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Reg
			resp := e.POST("/api/v1/reg").
				WithJSON(reg.Request{
					Email: tc.email,
				}).
				Expect().
				JSON().Object()

			if tc.error == true {
				resp.NotContainsKey("guid")
				resp.ContainsKey("message")
				return
			}
			resp.ContainsKey("guid")
			guid := resp.Value("guid").String().Raw()

			// Auth
			if tc.emptyGuid {
				guid = ""
			}
			if tc.fakeGuid {
				i, err := strconv.Atoi(guid)
				if err != nil {
					panic(err)
				}
				guid = strconv.Itoa(i + 1)
			}
			resp = e.GET("/api/v1/auth").WithQuery("guid", guid).Expect().JSON().Object()
			if tc.fakeGuid || tc.emptyGuid {
				resp.ContainsKey("message")
				return
			}
			resp.ContainsKey("access_token")
			resp.ContainsKey("refresh_token")
			access := resp.Value("access_token").String().Raw()
			refresh := resp.Value("refresh_token").String().Raw()
			//Refresh
			if tc.fakeAccess {
				access = fakeToken(access)
			}
			if tc.fakeRefresh {
				refresh = fakeToken(refresh)
			}
			if tc.emptyAccess {
				access = ""
			}
			if tc.emptyRefresh {
				refresh = ""
			}
			resp = e.POST("/api/v1/refresh").WithJSON(ref.Request{
				Access:  access,
				Refresh: refresh,
			}).Expect().JSON().Object()
			if tc.fakeAccess || tc.fakeRefresh || tc.emptyAccess || tc.emptyRefresh {
				resp.ContainsKey("message")
				return
			}
			resp.ContainsKey("access_token")
			resp.ContainsKey("refresh_token")
		})
	}
}
func fakeToken(tokenStr string) string {

	secretKey := "secret_key"
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		panic(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := claims["exp"].(float64)
		claims["exp"] = exp + 200.0
		token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		ans, err := token.SignedString([]byte(secretKey))
		if err != nil {
			panic(err)
		}
		return ans
	}
	panic("can't get claims from token")
}
