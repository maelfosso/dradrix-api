package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey struct {
	name string
}

var JwtMemberKey *contextKey
var JwtClaimsKey *contextKey
var JwtTokenKey *contextKey
var JwtErrorKey *contextKey
var jwtSecretKey []byte // (os.Getenv("JWT_SECRET"))

type TWAJWTClaims struct {
	*jwt.RegisteredClaims
	User interface{}
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}

func init() {
	JwtMemberKey = &contextKey{"Member"}
	JwtClaimsKey = &contextKey{"Claims"}
	JwtTokenKey = &contextKey{"Token"}
	JwtErrorKey = &contextKey{"Error"}
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))
	// TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), "s-tschwaa")
}

func GenerateJWTToken(data map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(60 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	// claims["authorized"] = true
	claims["user"] = data

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	var tokenString string

	authorizationHeader := r.Header.Get("Authorization")
	log.Println("Authorization Bearer : ", authorizationHeader)
	if strings.HasPrefix(authorizationHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authorizationHeader, "Bearer ")
	} else {
		return "", fmt.Errorf("no authorization found in header")
	}

	if tokenString == "" {
		return "", fmt.Errorf("empty authorization found")
	}

	return tokenString, nil
}

func Verifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ExtractTokenFromRequest(r)
		log.Println("Verifier : ", tokenString, err)

		ctx := r.Context()
		ctx = context.WithValue(ctx, JwtTokenKey, tokenString)
		ctx = context.WithValue(ctx, JwtErrorKey, err)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseJWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString, _ := ctx.Value(JwtTokenKey).(string)
		err, _ := ctx.Value(JwtErrorKey).(error)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		log.Println("ParseJWTToken starts with: ", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecretKey, nil
		})

		if err != nil {
			ctx = context.WithValue(ctx, JwtErrorKey, fmt.Errorf("invalidate token: %v", err))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if token == nil || !token.Valid {
			ctx = context.WithValue(ctx, JwtErrorKey, errors.New("invalid token"))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx = context.WithValue(ctx, JwtErrorKey, errors.New("invalid token claims"))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx = context.WithValue(ctx, JwtClaimsKey, claims["user"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err, _ := ctx.Value(JwtErrorKey).(error)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
