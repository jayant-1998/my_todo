package middlewares

import (
	"context"
	"jayant/database/dbHelper"
	"jayant/models"
	"jayant/utils"
	"net/http"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-api-key")
		user, err := dbHelper.GetUserBySession(token)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, err, "session is expired")
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

//func validateToken(w http.ResponseWriter, r *http.Request) (err error) {
//
//	if r.Header["X-Api-Key"] == nil {
//		fmt.Fprintf(w, "can not find token in header")
//		return
//	}
//
//	token, err := jwt.Parse(r.Header["X-Api-Key"][0], func(token *jwt.Token) (interface{}, error) {
//		_, ok := token.Method.(*jwt.SigningMethodHMAC)
//		if !ok {
//			return nil, fmt.Errorf("There was an error in parsing")
//		}
//		return models.JwtKey, nil
//	})
//
//	if token == nil {
//		fmt.Println(w, "invalid token")
//	}
//
//	return nil
//}
