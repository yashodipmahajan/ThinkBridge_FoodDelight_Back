package authorization

import (
	"FoodDelight/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var secrateKey = []byte("axsddhgsadghsa")

type Claims struct {
	Id uuid.UUID
	jwt.StandardClaims
}

func (c *Claims) Coder() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secrateKey)
}

func CheckToken(tokenString string) (*jwt.Token, *Claims, error) {
	var claim = &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
		return secrateKey, nil
	})
	fmt.Println("claim", claim)
	return token, claim, err
}
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("HELOO...")
		_, err := ValidateUserToken(w, r)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateUserToken(w http.ResponseWriter, r *http.Request) (*model.Admin, error) {
	authCookie, err := r.Cookie("auth")
	tokenString := authCookie.Value
	// fmt.Println("tokenstring ", tokenString)
	if err != nil {
		return nil, err
	}
	token, claim, err := CheckToken(tokenString)

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}
	fmt.Println("token.Claims", token.Claims)
	fmt.Println("claim", claim)
	// err, foundUser := model.FindUserByIDandCheckIsAdmin(claim.Id)

	err, foundAdmin := model.FindAdminByID(claim.Id)
	if err != nil {
		return nil, err
	}
	fmt.Println("foundUser", foundAdmin)
	return foundAdmin, nil

}
