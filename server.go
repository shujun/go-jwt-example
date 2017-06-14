package main

import (
	"net/http"
	"time"

	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/dgrijalva/jwt-go"
)

// User model
type User struct {
	UserId   string `form:"userid" json:"userid" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Field validator
func (u *User) Validate(errors *binding.Errors, req *http.Request) {

	if len(u.UserId) < 4 {
		errors.Fields["userid"] = "Too short; minimum 4 characters"
	}
}

const (
	ValidUser = "shujun"
	ValidPass = "123456"
	SecretKey = "WOW,MuchShibe,ToDogge"
)

type MyCustomClaim struct {
	Userid string `json:"userid"`
	jwt.StandardClaims
}

func main() {

	m := martini.Classic()

	m.Use(martini.Static("static"))
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(201, "index", nil)
	})

	// Authenticate user
	m.Post("/auth", binding.Bind(User{}), func(user User, r render.Render) {

		if user.UserId == ValidUser && user.Password == ValidPass {
			claim := &MyCustomClaim{
				user.UserId,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
				},
			}
			// Create JWT token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			tokenString, err := token.SignedString([]byte(SecretKey))
			if err != nil {
				fmt.Println("sign failed, error = ", err.Error())
				r.HTML(201, "error", nil)
				return
			}

			data := map[string]string{
				"token": tokenString,
			}
			fmt.Println("token: ", tokenString)
			r.HTML(201, "success", data)
		} else {
			fmt.Println("not match: ", user)
			r.HTML(201, "error", nil)
		}

	})
	// Check Key is ok
	m.Get("/debug/:token", func(params martini.Params, r render.Render) string {
		token, err := jwt.Parse(params["token"], func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err == nil && token.Valid {
			return "User id: " + token.Claims.(MyCustomClaim).Userid
		} else {
			return "Invalid"
		}
	})

	// Only accesible if authenticated
	m.Post("/secret", func() {

	})

	m.Run()
}
