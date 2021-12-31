package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"package/model"
	"package/utility"

	"github.com/dgrijalva/jwt-go"
)

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err model.Error
			err = utility.SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}
		secretkey := os.Getenv("SECRET_KEY")
		var mySigningKey = []byte(secretkey)
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err model.Error
			err = utility.SetError(err, "Your Token has been expired")
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return
			} else if claims["role"] == "user" {
				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		var reserr model.Error
		reserr = utility.SetError(reserr, "Not Authorized")
		json.NewEncoder(w).Encode(err)
	}
}
