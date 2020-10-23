package web

import (
	"context"
	"net/http"
	"strconv"
	"time"

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
	tokenString := CreateJwtToken(credentials, expirationTime)

	http.SetCookie(w, &http.Cookie{
		Name:    tokenName,
		Value:   tokenString,
		Expires: expirationTime,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func CreateJwtToken(credentials dto.Credentials, expirationTime time.Time) string {
	claims := &Claims{
		UserID: credentials.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func VerifyJwtCookie(c *Container) mux.MiddlewareFunc {
	return func (next http.Handler) http.Handler {
	    return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			apiKey := getAuthKey(request)
			apiToken := getAuthToken(request)

			isAPIRequest := len(apiKey) > 0 && len(apiToken) > 0
			if isAPIRequest {
				handleAPIRequest(c, next, response, request)

				return
			}

			handleUIRequest(c, next, response, request)
		})
	}
}

func getAuthKey(request *http.Request) string {
	var apiKey string
	if key, ok := request.Header["X-Auth-Key"]; ok {
		apiKey = key[0]
	}

	return apiKey
}

func getAuthToken(request *http.Request) string {
	var apiToken string
	if key, ok := request.Header["X-Auth-Token"]; ok {
		apiToken = key[0]
	}

	return apiToken
}


func handleAPIRequest(c *Container, next http.Handler, response http.ResponseWriter, request *http.Request) {
	tokenString := getAuthToken(request)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Can't read access token"))

		return
	}

	if !token.Valid {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Invalid token"))

		return
	}

	vars := mux.Vars(request)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Invalid URL"))

		return
	}

	isNotAuthorized := uint(userID) != claims.UserID
	if isNotAuthorized {
		response.WriteHeader(http.StatusUnauthorized)

		return
	}

	// Check is API key enabled
	// @TODO move the API keys to redis
	isAPIKeyValid := c.UserService.IsAPIKeyValidForUser(int(userID), getAuthKey(request))
	if !isAPIKeyValid {
		response.WriteHeader(http.StatusUnauthorized)

		return
	}

	var userToken ContextKey = "user"
	ctx := context.WithValue(request.Context(), userToken, token)
	next.ServeHTTP(response, request.WithContext(ctx))
}


func handleUIRequest(c *Container, next http.Handler, w http.ResponseWriter, r *http.Request) {
	signinURL, err := c.Router.Get("signin").URL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Can't generate 'signin' URL"))

		return
	}
	signinURLStr := signinURL.String()

	cookie, err := r.Cookie(tokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, signinURLStr, http.StatusFound)
			return
		}
		http.Redirect(w, r, signinURLStr, http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Redirect(w, r, signinURLStr, http.StatusFound)
			return
		}
		http.Redirect(w, r, signinURLStr, http.StatusBadRequest)
		return
	}
	if !token.Valid {
		http.Redirect(w, r, signinURLStr, http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Redirect(w, r, signinURLStr, http.StatusFound)
		return
	}

	isNotAuthorized := uint(userID) != claims.UserID
	if isNotAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not authorized!"))

		return
	}

	var userToken ContextKey = "user"
	ctx := context.WithValue(r.Context(), userToken, token)
	next.ServeHTTP(w, r.WithContext(ctx))
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