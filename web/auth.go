package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"palindromex/web/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type ContextKey string
const userContext ContextKey = "user"

const tokenName = "access-token"
const cookieValidityPeriod = time.Duration(7 * 24) * time.Hour
const maxDuration = 1<<63 - 1
const permanentTokenPeriod = time.Duration(maxDuration)

type Claims struct {
	UserID uint	`json:user_id`
	APIKey string `json:api_key`
	jwt.StandardClaims
}

func SetJwtCookie(c *Container, w http.ResponseWriter, user model.User) error {
	expirationTime := time.Now().Add(cookieValidityPeriod)
	tokenString := CreateJwtToken(user.ID, expirationTime, "")

	http.SetCookie(w, &http.Cookie{
		Name:    tokenName,
		Value:   tokenString,
		Expires: expirationTime,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func GetAPICredentials(c *Container, w http.ResponseWriter, user model.User) (string, string, error) {
	expirationTime := time.Now().Add(permanentTokenPeriod)
	apiKey := fmt.Sprintf("key-%d", expirationTime.Unix())
	err := c.ApiKeyService.CreateNew(user, apiKey)
	if err != nil {
		return "", "", err
	}

	tokenString := CreateJwtToken(user.ID, expirationTime, apiKey)

	return apiKey, tokenString, nil
}

func CreateJwtToken(userID uint, expirationTime time.Time, apiKey string) string {
	claims := &Claims{
		UserID: userID,
		APIKey: apiKey,
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
			apiToken := getAuthToken(request)

			isAPIRequest := len(apiToken) > 0
			if isAPIRequest {
				handleAPIRequest(c, next, response, request)

				return
			}

			handleUIRequest(c, next, response, request)
		})
	}
}

func getAuthToken(request *http.Request) string {
	var apiToken string
	if key, ok := request.Header["Authorization"]; ok {
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
		response.Write([]byte("Unauthorized"))

		return
	}

	// Check is API key enabled
	// @TODO move the API keys to redis
	isAPIKeyValid := c.UserService.IsAPIKeyValidForUser(int(userID), claims.APIKey)
	if !isAPIKeyValid {
		response.WriteHeader(http.StatusUnauthorized)

		return
	}

	ctx := context.WithValue(request.Context(), userContext, token)
	next.ServeHTTP(response, request.WithContext(ctx))
}


func handleUIRequest(c *Container, next http.Handler, w http.ResponseWriter, r *http.Request) {
	signinURL, err := c.Router.Get("signin").URL()
	if err != nil {
		msg := "Server error: can't generate 'signin' URL"
		log.Println(msg)
		w.Write([]byte(msg))

		return
	}

	signinURLStr := signinURL.String()
	cookie, err := r.Cookie(tokenName)
	if err != nil {
		c.Flash.AddWarning(w, r, "Not logged in")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, signinURLStr, http.StatusFound)

			return
		}
		RemoveAccessToken(w)
		http.Redirect(w, r, signinURLStr, http.StatusFound)

		return
	}

	tokenString := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		c.Flash.AddWarning(w, r, "Error while parsing access token")
		RemoveAccessToken(w)
		http.Redirect(w, r, signinURLStr, http.StatusFound)

		return
	}
	if !token.Valid {
		c.Flash.AddWarning(w, r, "Invalid access token")
		RemoveAccessToken(w)
		http.Redirect(w, r, signinURLStr, http.StatusFound)

		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		c.Flash.AddError(w, r, "Invalid URL")
		http.Redirect(w, r, signinURLStr, http.StatusFound)

		return
	}

	isNotAuthorized := uint(userID) != claims.UserID
	if isNotAuthorized {
		c.Flash.AddError(w, r, "Not authorized")
		http.Redirect(w, r, signinURLStr, http.StatusFound)

		return
	}

	ctx := context.WithValue(r.Context(), userContext, token)
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