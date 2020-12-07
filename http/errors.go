package main

import (
	"errors"
	"net/http"

	"github.com/Vindexus/userchat-api"
)

var (
	ErrNotAdmin      = errors.New("you are not an admin")
	ErrNotLoggedIn   = errors.New("you are not logged in")
	ErrNotAuthorized = errors.New("you do not have access")
	ErrNotFound      = errors.New("not found")
)

func (c *CustomRouteContext) HandleError(err error) {
	resp := map[string]interface{}{
		"message": err.Error(),
	}
	status := http.StatusInternalServerError
	switch err {
	case userchat.ErrUserNotFound:
		status = http.StatusNotFound
	case ErrInvalidLogin:
		status = http.StatusBadRequest
	}
	c.JSON(status, resp)
}

func (c *CustomRouteContext) HandledError(err error) bool {
	if err != nil {
		c.HandleError(err)
		return true
	}

	return false
}
