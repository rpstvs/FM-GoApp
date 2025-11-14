package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rpstvs/fm-goapp/internal/store"
	"github.com/rpstvs/fm-goapp/internal/tokens"
	"github.com/rpstvs/fm-goapp/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)

	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)

	if !ok {
		return nil
	}

	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headersParts := strings.Split(authHeader, " ")

		if len(headersParts) != 2 || headersParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{})
			return
		}

		token := headersParts[1]

		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)

		if err != nil {
			return
		}
		if user == nil {
			return
		}

		r = SetUser(r, user)

		next.ServeHTTP(w, r)
		return
	})
}

func (um *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)

			if user.IsAnonymous() {
				utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{})
			}

			next.ServeHTTP(w, r)
		})
}
