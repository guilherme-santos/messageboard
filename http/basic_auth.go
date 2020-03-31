package http

import (
	"fmt"
	"net/http"

	"github.com/guilherme-santos/messageboard"
)

// Copied from https://github.com/go-chi/chi/blob/master/middleware/basic_auth.go
func BasicAuth(realm string, creds map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w, realm)
				return
			}

			credPass, credUserOk := creds[user]
			if !credUserOk || pass != credPass {
				basicAuthFailed(w, realm)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	err := messageboard.NewError("unauthorized", "user is not authorized to access this resource")
	responseError(w, err)
}
