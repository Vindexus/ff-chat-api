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
