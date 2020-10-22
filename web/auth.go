package web

import (
	"time"
	"context"
	"net/http"
	"strconv"

	"palindromex/web/dto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type ContextKey string

const tokenName = "access-token"
const tokenValidityPeriod = time.Duration(7 * 24) * time.Hour

type Claims struct {
	UserID uint	`json:user_id`
	jwt.StandardClaims
}

func SetJwtCookie(c *Container, w http.ResponseWriter, credentials dto.Credentials) error {
	expirationTime := time.Now().Add(tokenValidityPeriod)
	claims := &Claims{
		UserID: credentials.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    tokenName,
		Value:   tokenString,
		Expires: expirationTime,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func VerifyJwtCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie(tokenName)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(response, request, "/signin", http.StatusFound)
				return
			}
			http.Redirect(response, request, "/signin", http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(response, request, "/signin", http.StatusFound)
				return
			}
			http.Redirect(response, request, "/signin", http.StatusBadRequest)
			return
		}
		if !token.Valid {
			http.Redirect(response, request, "/signin", http.StatusFound)
			return
		}

		vars := mux.Vars(request)
		userID, err := strconv.Atoi(vars["userID"])
		if err != nil {
			http.Redirect(response, request, "/signin", http.StatusFound)
			return
		}

		isNotAuthorized := uint(userID) != claims.UserID
		if isNotAuthorized {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Not authorized!"))

			return
		}		

		var userToken ContextKey = "user"
		ctx := context.WithValue(request.Context(), userToken, token)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

func RemoveAccessToken(w http.ResponseWriter) {
	newCookie := http.Cookie {
		Name: tokenName,
		Value: "",
		Path: "/",
		Expires: time.Unix(0,0),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &newCookie)
}