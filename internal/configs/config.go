package configs

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DSN        string

	RedisURL string
	NatsURL  string

	JWTSecret          string
	CSRFSecrete        string
	CacheEncryptionKey string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string

	S3Bucket          string
	S3Region          string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Endpoint        string 
	S3PublicEndpoint  string
	S3UsePathStyle    bool
}
