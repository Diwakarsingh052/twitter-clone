package middleware

import (
	"encoding/json"
	"net/http"
	"twitter-clone/graph/model"
)

type RequireUser struct {
	Us *model.UserService
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//cookie, err := r.Cookie("remember_token")

		if mw.Us.Token == "" {
			//if err != nil {
			data, _ := json.Marshal("You are not logged in")

			w.Write(data)
			return
		}
		_, err := mw.Us.ByRemember(mw.Us.Token)
		if err != nil {
			data, _ := json.Marshal("User Token not Found")
			w.Write(data)
			return
		}

		next(w, r)
		return
	})
}
