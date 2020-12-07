package main

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/Vindexus/userchat-api"
)

var (
	ErrInvalidLogin = errors.New("invalid login info")
	ErrBlankJWT     = errors.New("Auth header is blank")
	ErrInvalidJWT   = errors.New("invalid JWT provided")
)

func SessionRoutes(s *Server) {
	s.R.AddRoute(NewGET("/session", getSession))
	s.R.AddRoute(NewPOST("/login", postLogin))
}

func postLogin(c *CustomRouteContext) {
	type Payload struct {
		Username string
		Password string
	}

	var pl Payload
	if err := c.ShouldBindJSON(&pl); c.HandledError(err) {
		return
	}

	user, err := GetUserByUsername(pl.Username)
	if err == ErrUserNotFound {
		c.HandleError(ErrInvalidLogin)
		return
	}

	if !ComparePassword(user.Password, pl.Password) {
		c.HandleError(ErrInvalidLogin)
		return
	}

	jwt, err := SignJWT(c.C.JWTSecret, user.Id)
	/*	// We will also supply the JWT as a cookie
		// so that it can be used to authentication requests
		// from the browser that won't set the Authorization header
		// EG: iframes that preview campaigns
		expire := time.Now().Add(JWTDuration)
		// Set the cookie for the requesting domain (the app where the login is coming from)
		cookie2 := http.Cookie{
			Name:     "jwt",
			Value:    jwt,
			Expires:  expire,
			Domain:   c.C.CookieDomain,
			HttpOnly: true,
		}
		http.SetCookie(c.W, &cookie2)*/

	c.JSON(http.StatusOK, M{
		"jwt": jwt,
	})
}

func getSession(c *CustomRouteContext) {
	user, err := c.GetCurrentUser()
	if err != nil {
		fmt.Println("GetCtxUser Error:", err)
		c.JSON(http.StatusOK, M{
			"LoggedIn": false,
		})
		return
	}
	c.JSON(http.StatusOK, M{
		"LoggedIn": true,
		"UserId":   user.Id,
	})
}
