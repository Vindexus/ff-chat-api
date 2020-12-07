package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	. "github.com/cyc-ttn/gorouter"

	. "github.com/Vindexus/userchat-api"
)

var conf string

var (
	ErrBadParam = errors.New("empty or malformed param found")
)

const JWT_HEADER = "Authorization" // name of header in requests where the JWT is

type CustomHandlerFunc func(ctx *CustomRouteContext)

type Server struct {
	C *Config
	R *RouterNode
}

type CustomRouteContext struct {
	RouteContext
	C *Config
}

func (c *CustomRouteContext) GetParam(param string) string {
	s, ok := c.Params[param]
	if ok {
		return s
	}

	return ""
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, v := range s.C.CorsAllowOrigins {
		if v == origin {
			w.Header().Set("Access-Control-Allow-Origin", v)
			break
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT, PATCH")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := &RouteContext{
		W:      w,
		R:      r,
		Method: r.Method,
		Path:   r.URL.Path,
		Query:  r.URL.Query(),
	}
	route, err := s.R.Match(r.Method, r.URL.Path, ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	custom := &CustomRouteContext{
		RouteContext: *ctx,
		C:            s.C,
	}
	route.GetHandler()(custom)
}

func NewRoute(method string, path string, handler CustomHandlerFunc) *DefaultRoute {
	return &DefaultRoute{
		Method: method,
		Path:   path,
		HandlerFunc: func(ctx interface{}) {
			handler(ctx.(*CustomRouteContext))
		},
	}
}

func (c *CustomRouteContext) GetCurrentUser() (*User, error) {
	tokenString := c.R.Header.Get(JWT_HEADER)
	if tokenString == "" {
		return nil, ErrBlankJWT
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.C.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	// Convert its claims to our custom, wich has the user id
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Find that user and verify that they are an admin
		return GetUserById(claims.UserId)
	}
	return nil, ErrInvalidJWT
}

/*
func GetCtxUser(crc *CustomRouteContext) bool {
	cook, err := crc.R.Cookie("jwt")
	if err == http.ErrNoCookie {
		crc.HandleError(ErrNotLoggedIn)
		return false
	}
	if crc.HandledError(err) {
		return false
	}
	tokenString := cook.Value

	// We validate the provided JWT string
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(crc.C.JWTSecret), nil
	})
	if crc.HandledError(err) {
		return false
	}

	// Convert its claims to our custom, wich has the user id
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Find that user and verify that they are an admin
		user, err := GetUserById(crc.S.DB, claims.UserId)
		if crc.HandledError(err) {
			return false
		}
		crc.User = user

		if crc.C.Unsafe {
			fmt.Println("UNSAFE MODE: No authentication is done")
			return true
		}

		if crc.User == nil {
			crc.HandleError(ErrNotLoggedIn)
			return false
		}

		return true
	}
	crc.HandleError(errors.New("token wasnt or ok or wasn't valid"))
	return false
}*/

func NewGET(path string, handler CustomHandlerFunc) *DefaultRoute {
	return NewRoute(http.MethodGet, path, handler)
}

func NewPOST(path string, handler CustomHandlerFunc) *DefaultRoute {
	return NewRoute(http.MethodPost, path, handler)
}

func NewServer(c *Config) *Server {
	s := &Server{
		R: NewRouter(),
		C: c,
	}

	ChatRoutes(s)
	SessionRoutes(s)
	UserRoutes(s)

	s.R.AddRoute(NewGET("/ping", getPing))

	return s
}

func main() {
	conf := &Config{
		CorsAllowOrigins: []string{
			"http://localhost:8080",
		},
		JWTSecret: "y348t9tyhewoahge89whag8eawhgewa",
		Port:      4949,
	}
	s := NewServer(conf)
	fmt.Printf("Listening on port %d\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), s)
}

func getPing(ctx *CustomRouteContext) {
	ctx.String(http.StatusOK, "Pong!")
}
