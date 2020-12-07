package main

import (
	"net/http"
	"testing"

	"github.com/Vindexus/userchat-api/pkg/fixtures"

	"github.com/monstercat/golib/request"

	"github.com/monstercat/golib/expectm"

	. "github.com/Vindexus/userchat-api"

	grtest "github.com/Vindexus/go-router-test"
)

func TestSessionRoutes(t *testing.T) {
	MustSetupTestServer()

	test := &grtest.RouteTest{
		Path:           "/login",
		Method:         http.MethodPost,
		ExpectedStatus: http.StatusBadRequest,
	}

	var jwt string
	tests := test.Apply([]*grtest.RouteTest{
		{
			Body: M{
				"Username": "",
			},
		},
		{
			Body: M{
				"Username": fixtures.User1Username,
			},
		},
		{
			Body: M{
				"Username": fixtures.User1Username,
				"Password": "fdsafkldajslkfjas",
			},
		},
		{
			Body: M{
				"Username": fixtures.User1Username,
				"Password": "password",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedM: &expectm.ExpectedM{
				"jwt": ExpectedReturnedString(&jwt),
			},
		},
	})

	if err := runTests(tests); err != nil {
		t.Fatal(err)
	}

	if jwt == "" {
		t.Error("JWT should not be blank")
	}

	t.Log("jwt", jwt)

	test = &grtest.RouteTest{
		Method:         http.MethodGet,
		ExpectedStatus: http.StatusOK,
		Path:           "/session",
		ModifyParams: func(params *request.Params) error {
			params.Headers[JWT_HEADER] = jwt
			return nil
		},
		ExpectedM: &expectm.ExpectedM{
			"LoggedIn": true,
			"UserId":   fixtures.User1Id,
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}
