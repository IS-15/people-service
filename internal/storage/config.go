package storage

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}
