package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Signup(c echo.Context) error {
	var user User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		c.String(400, "Invalid request body")
	}

	// check if phone number and password are valid
	if len(user.PhoneNumber) != 11 {
		return c.String(400, "Invalid phone number. Phone number must be 11 digits")
	}
	if len(user.Password) < 8 || len(user.Password) > 64 {
		return c.String(400, "Invalid password. Password must be between 8 and 64 characters")
	}

	// check if user already exists
	exists, err := redisClient.KeyExists(user.PhoneNumber)
	if err != nil {
		log.Println(err)
		return c.String(500, "Failed to create user")
	}
	if exists {
		return c.String(409, "Phone number already exists")
	}

	// set user role
	user.Role = "user"

	// create user
	_, err = redisClient.CreateUser(user)
	if err != nil {
		log.Println(err)
		return c.String(500, "Failed to create user")
	}

	// create JWT
	jwt, err := CreateJWT(user, config.JWTSecret)
	if err != nil {
		log.Println(err)
		return c.String(500, "Failed to create jwt")
	}

	c.SetCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/",
		// Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
		Domain:   config.Domain,
	})

	return c.NoContent(201)
}

func Signin(c echo.Context) error {
	var user User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		c.String(400, "Invalid request body")
	}

	// check if phone number and password are valid
	if len(user.PhoneNumber) != 11 {
		return c.String(400, "Invalid phone number. Phone number must be 11 digits")
	}
	if len(user.Password) < 8 || len(user.Password) > 64 {
		return c.String(400, "Invalid password. Password must be between 8 and 64 characters")
	}

	// find user
	u, err := redisClient.FindUser(user.PhoneNumber)
	if err != nil {
		if err.Error() == "user not found" {
			return c.String(400, "Wrong phone number or password")
		} else {
			log.Println(err)
			return c.String(500, "Failed to find user")
		}
	}

	// check password
	if !CheckPasswordHash(user.Password, u.Password) {
		return c.String(400, "Wrong phone number or password")
	}

	// create JWT
	jwt, err := CreateJWT(*u, config.JWTSecret)
	if err != nil {
		log.Println(err)
		return c.String(500, "Failed to create jwt")
	}

	// set JWT cookie
	c.SetCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/",
		// Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
		Domain:   config.Domain,
	})

	return c.NoContent(200)
}

func Signout(c echo.Context) error {
	// set JWT cookie
	c.SetCookie(&http.Cookie{Name: "jwt", Value: "nil", Path: "/",
		// Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		Domain:   config.Domain,
	})

	return c.NoContent(200)
}

func Account(c echo.Context) error {
	phoneNumber := c.Get("user").(*User).PhoneNumber

	// find user
	u, err := redisClient.FindUser(phoneNumber)
	if err != nil {
		if err.Error() == "user not found" {
			return c.String(400, "Wrong phone number or password")
		} else {
			log.Println(err)
			return c.String(500, "Failed to find user")
		}
	}

	u.Password = ""

	return c.JSON(200, u)
}
