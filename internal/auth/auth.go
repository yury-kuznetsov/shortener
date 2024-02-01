package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenExp = time.Hour
const SecretKey = "SECRET_KEY"

// Handle проверяет наличие и подлинность куки.
// В случае неудачи создает новую куку, если `create` = true.
func Handle(handler http.HandlerFunc, create bool) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		cookie, _ := req.Cookie("token")

		if cookie == nil && create {
			if token, _ := buildToken(); token != "" {
				cookie = &http.Cookie{
					Name:    "token",
					Value:   token,
					Expires: time.Now().Add(TokenExp),
				}
			}
		}

		if cookie != nil {
			userID := getUserID(cookie.Value)
			req.Header.Set("Content-User-ID", strconv.Itoa(userID))
			res.Header().Set("Authorization", cookie.Value)
			http.SetCookie(res, cookie)
		}

		handler(res, req)
	}

	return handlerFunc
}

func getUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return 0
	}

	return claims.UserID
}

func buildToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: rand.Intn(1000),
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
