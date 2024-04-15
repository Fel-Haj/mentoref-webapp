package middleware

type BoolContextKey bool
type UserMailContextKey string

var (
	AuthContextKey BoolContextKey
	UserContextKey UserMailContextKey
)

const (
	UserMailKey string = "user_ID"
)
