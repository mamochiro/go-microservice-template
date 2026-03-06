package middleware

import (
	"net/http"

	"github.com/unrolled/secure"
)

func SecureHeaders(env string) func(http.Handler) http.Handler {
	isDevelopment := env != "production"

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IsDevelopment:         isDevelopment,
		// HSTS
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := secureMiddleware.Process(w, r)
			if err != nil {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
