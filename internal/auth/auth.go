package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the custom claims for a JWT token, which includes the standard RegisteredClaims and an additional UserID field.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// TokenExp represents the expiration time for a token.
const TokenExp = time.Hour

// SecretKey is a constant value used for signing and validating JWT tokens.
// It should be kept secret and not shared publicly.
// Example usage:
//   - In getUserID function, SecretKey is used as the key to validate the token and extract the user ID from it.
//   - In buildToken function, SecretKey is used as the key to sign the token.
//
// Type: string
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
