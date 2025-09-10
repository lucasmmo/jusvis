package auth

import (
	"jusvis/pkg/token"
	"net/http"
	"strings"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawToken := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(rawToken, "Bearer ")
		if err := token.Validate(tokenString); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, err := token.Parse(tokenString, []byte("banana"))
		if err != nil {
			http.Error(w, "cannot parse token", http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-User-ID", claims["id"].(string))
		next.ServeHTTP(w, r)
	}
}
