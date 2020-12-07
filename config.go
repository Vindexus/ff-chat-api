package userchat

type Config struct {
	CookieDomain     string
	CorsAllowOrigins []string
	JWTSecret        string
	Port             int
}
