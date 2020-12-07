package main

import (
	"net/http"
	"testing"

	grtest "github.com/Vindexus/go-router-test"
)

func TestPingRoute(t *testing.T) {
	MustSetupTestServer()

	test := &grtest.RouteTest{
		Path:           "/ping",
		ExpectedStatus: http.StatusOK,
	}

	if err := runTest(test); err != nil {
		t.Error(err)
	}
}
