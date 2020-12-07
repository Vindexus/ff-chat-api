package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/monstercat/golib/request"

	. "github.com/Vindexus/userchat-api"

	grtest "github.com/Vindexus/go-router-test"
)

func testURL(conf *Config, test *grtest.RouteTest) string {
	return fmt.Sprintf("http://localhost:%d%s", conf.Port, test.Path)
}

func runTest(test *grtest.RouteTest) error {
	conf := MustGetTestConfig()
	if test.URL == "" {
		test.URL = testURL(conf, test)
	}
	return test.Run()
}

func runTests(tests []*grtest.RouteTest) error {
	conf := MustGetTestConfig()
	for i, v := range tests {
		tests[i].URL = testURL(conf, v)
	}

	return grtest.RunTests(tests)
}

var testServerRunning = false

func MustGetTestConfig() *Config {
	return &Config{
		Port:      4747,
		JWTSecret: "2tqc9y80qncy98tqcy9cq2y89cqt",
	}
}

func MustSetupTestServer() *Config {
	conf := MustGetTestConfig()
	if testServerRunning {
		return conf
	}
	testServerRunning = true
	server := NewServer(conf)
	testServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: server,
	}

	go func() {
		if err := testServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	return conf
}

// Checks that a value is a string, and then updates a pointer to have
// the value that was found
func ExpectedReturnedString(str *string) func(val interface{}) error {
	return func(val interface{}) error {
		newStr, ok := val.(string)
		if !ok || newStr == "" {
			return errors.New("couldn't convert Id to string")
		}
		*str = newStr
		return nil
	}
}

func modParamJWT(userId int) func(params *request.Params) error {
	conf := MustGetTestConfig()
	return func(params *request.Params) error {
		jwt, err := SignJWT(conf.JWTSecret, userId)
		if err != nil {
			return err
		}
		params.Headers[JWT_HEADER] = jwt
		return nil
	}

}
