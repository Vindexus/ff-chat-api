package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/monstercat/golib/expectm"

	"github.com/Vindexus/userchat-api/pkg/fixtures"

	grtest "github.com/Vindexus/go-router-test"
)

func TestGetChatRoute(t *testing.T) {
	MustSetupTestServer()
	test := grtest.RouteTest{
		ModifyParams: modParamJWT(fixtures.User1Id),
	}

	tests := test.Apply([]*grtest.RouteTest{
		{
			Path:           fmt.Sprintf("/chat/%d", 53245432),
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Path:           fmt.Sprintf("/chat/%d", fixtures.User2Id),
			ExpectedStatus: http.StatusOK,
			ExpectedM: &expectm.ExpectedM{
				"OtherUser.Id":       fixtures.User2Id,
				"Messages.#":         2,
				"Messages.0.Message": "Hello, how are you?",
			},
		},
		{
			Path:           fmt.Sprintf("/chat/%d", fixtures.User1Id),
			ExpectedStatus: http.StatusInternalServerError,
		},
	})
	if err := runTests(tests); err != nil {
		t.Error(err)
	}
}
