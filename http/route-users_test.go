package main

import (
	"net/http"
	"testing"

	grtest "github.com/Vindexus/go-router-test"
	"github.com/monstercat/golib/expectm"
)

func TestGetUsersRoute(t *testing.T) {
	MustSetupTestServer()
	test := &grtest.RouteTest{
		Path:           "/users",
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Users.#":          2,
			"Users.1.Username": "Storm",
			"Users.0.Id":       "23",
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}
