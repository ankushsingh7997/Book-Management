// implement the middleware to protect the routes

// token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 	return []byte(secret), nil
// })
//  above is the example of we are parsing the token with the claims and the secret key

package middlewares

import (
	"net/http"
	"strings"

	"github.com/ankush/bookstore/pkg/models"
)

func Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get the Authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header is in the format `Bearer token`
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Validate the token
		tokenString := headerParts[1]
		_, err := models.ValidateToken(tokenString)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(res, req)
	})
}
