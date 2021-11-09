package http

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware - adds middleware around endpoints
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logrus.WithFields(logrus.Fields{
			"Method": r.Method,
			"Path":   r.URL.Path,
		}).Info("handled request")

		next.ServeHTTP(w, r)
	})
}

// validateToken - validates an incoming jwt token
func validateToken(accessToken string) bool {
	// replace this by loading in a private RSA cert for more security
	var mySigningKey = []byte(os.Getenv("SIGNING_KEY"))
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

// JWTAuth - middleware that provides JWT authorization for specified endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("jwt auth endpoint hit")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
	}
}
