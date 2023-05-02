package middlewares

import (
	"context"
	"github.com/sirupsen/logrus"
	"jayant/database/dbHelper"
	"jayant/models"
	"net/http"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		user, err := dbHelper.GetUserBySession(apiKey)
		if err != nil || user == nil {
			logrus.WithError(err).Errorf("failed to get user with token: %s", apiKey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userContext, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContext).(*models.User)
	if ok && user != nil {
		return user
	}
	return nil
}
