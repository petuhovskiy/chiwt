package web

import "net/http"

type ctxKey string

const (
	authKey ctxKey = "auth"
)

func GetAuth(r *http.Request) AuthContext {
	ctx := r.Context()
	return ctx.Value(authKey).(AuthContext)
}
