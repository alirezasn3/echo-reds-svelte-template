package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwt, err := c.Cookie("jwt")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.NoContent(401)
			}
			log.Println(err)
			return c.NoContent(500)
		}
		user, err := DecodeJWT(jwt.Value, config.JWTSecret)
		if err != nil {
			return c.NoContent(401)
		}
		c.Set("user", user)

		return next(c)
	}
}

func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user").(*User).Role != "admin" {
			return c.NoContent(403)
		}

		return next(c)
	}
}

func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: config.AllowedOrigins, AllowCredentials: true})
}
