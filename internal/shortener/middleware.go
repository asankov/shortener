package shortener

import (
	"errors"
	"net/http"

	"github.com/asankov/shortener/internal/apis"
	"github.com/asankov/shortener/internal/auth"
	"github.com/asankov/shortener/internal/users"
)

func (h *handler) authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		roles, ok := ctx.Value(apis.JWTScopes).([]string)
		if !ok {
			// no roles required for this endpoint, no token required
			next.ServeHTTP(w, r)
			return
		}

		jwtToken := r.Header.Get("Authorization")
		if jwtToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization header not provided"))
			return
		}

		user, err := h.authenticator.DecodeToken(jwtToken)
		if err != nil {
			if errors.Is(err, auth.ErrTokenExpired) {
				// TODO: return body that indicates that the UI should try to get a new token
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if errors.Is(err, auth.ErrInvalidSignature) || errors.Is(err, auth.ErrInvalidFormat) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Warn("unknown error while decoding token", "error", err)
			return
		}

		for _, rr := range roles {
			role, err := users.RoleFrom(rr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if user.HasRole(role) {
				next.ServeHTTP(w, r)
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	})
}
