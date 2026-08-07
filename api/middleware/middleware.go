package middleware

func Wrap(next handleFunc) handleFunc {
	return auth(next)
}
