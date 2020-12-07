package main

import (
	"errors"
	"net/http"
	"strconv"

	. "github.com/Vindexus/userchat-api"
)

var (
	ErrChatSelf = errors.New("you cannot chat with yourself")
)

func ChatRoutes(s *Server) {
	s.R.AddRoute(NewGET("/chat/:otherUserId", func(c *CustomRouteContext) {
		idS := c.GetParam("otherUserId")
		id, err := strconv.Atoi(idS)
		if c.HandledError(err) {
			return
		}

		currentUser, err := c.GetCurrentUser()
		if c.HandledError(err) {
			return
		}

		if currentUser.Id == id {
			c.HandleError(ErrChatSelf)
			return
		}

		otherUser, err := GetUserById(id)
		if c.HandledError(err) {
			return
		}

		messages, err := GetChatMessages(currentUser, otherUser)
		if c.HandledError(err) {
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"OtherUser": otherUser,
			"Messages":  messages,
		})
	}))
}
