package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io/engineio/transport/polling"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"

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
func (s *Server) Start() error {
	wsServer, err := s.GetWSServer()
	if err != nil {
		return err
	}
	go wsServer.Serve()
	defer wsServer.Close()

	http.Handle("/socket.io/", wsServer)
	http.Handle("/", s)

	log.Println("Serving at localhost:" + fmt.Sprintf(":%d", s.C.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.C.Port), nil))
	return nil
}

// I couldn't actually get the websocket server to work
// The package I was attempting to use that implements it in Golang
// had a few different versions, and their examples were using an older version
// It also appears to have some weird encoding issue with its JSON
// Anyway this doesn't work
func (s *Server) GetWSServer() (*socketio.Server, error) {
	allowOrigin := func(r *http.Request) bool {
		return true
	}
	server, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				Client: &http.Client{
					Timeout: time.Minute,
				},
				CheckOrigin: allowOrigin,
			},
			&websocket.Transport{
				CheckOrigin: allowOrigin,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// All of this server code is taken from their examples, which is a basic
	// chat example without room
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		fmt.Println("/chat msg", msg)
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	return server, nil
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
	if err := s.Start(); err != nil {
		fmt.Println("Error starting server", err)
	}
}

func getPing(ctx *CustomRouteContext) {
	ctx.String(http.StatusOK, "Pong!")
}
