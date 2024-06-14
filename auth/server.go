package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var secretKey = []byte("your-secret-key")

type User struct {
	Username string `json:"username" xml:"username" form:"username" query:"username"`
	Password string `json:"password" xml:"password" form:"password" query:"password"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
}

type UserBTS struct {
	Name    string
	Pass    string
	Email   string
	IsAdmin bool
}

type UserDTO struct {
	Name    string
	Email   string
	IsAdmin bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getRole(user UserBTS) string {
	if user.IsAdmin {
		return "senior"
	}
	return "employee"
}

func userValidation(u UserBTS) bool {
	if u.Name == "joe" && u.Pass == "secret" {
		return true
	}
	return false
}

func createToken(user UserBTS) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Name,                        // Subject (user identifier)
		"email": user.Email,                       // Email of user
		"iss":   "todo-app",                       // Issuer
		"aud":   getRole(user),                    // Audience (user role)
		"exp":   time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat":   time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	check(err)
	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Status 200 Ok!")
	})

	e.POST("/auth", func(c echo.Context) error {
		// curl -X POST http://localhost:1323/auth   -H 'Content-Type: application/json'   -d @user.json
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}

		user := UserBTS{
			Name:    u.Username,
			Pass:    u.Password,
			Email:   u.Email,
			IsAdmin: false,
		}

		if userValidation(user) {
			fmt.Print(string(user.Email))
			tokenString, err := createToken(user)
			fmt.Print(string(tokenString))
			check(err)
			return c.JSON(http.StatusCreated, tokenString)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}

	})
	e.Logger.Fatal(e.Start(":1323"))
}
