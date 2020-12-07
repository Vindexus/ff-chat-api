package main

import (
	"net/http"

	"github.com/Vindexus/userchat-api"
)

func UserRoutes(s *Server) {
	s.R.AddRoute(NewGET("/users", func(c *CustomRouteContext) {
		users := userchat.GetUsers()

		c.JSON(http.StatusOK, map[string]interface{}{
			"Users": users,
		})
	}))
}
