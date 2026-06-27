package configs

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DSN        string

	RedisAddr   string

	JWTSecret   string
	CSRFSecrete string
}
