package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Close() error
	DB() *sql.DB
}

type service struct {
	db *sql.DB
}

var (
	database   = getEnv("POSTGRES_DB", "postgres")
	password   = getEnv("POSTGRES_PASSWORD", "postgres")
	username   = getEnv("POSTGRES_USER", "postgres")
	port       = getEnv("POSTGRES_PORT", "5432")
	host       = getEnv("POSTGRES_HOST", "127.0.0.1")
	schema     = getEnv("POSTGRES_SCHEMA", "public")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	sslmode := getEnv("DB_SSLMODE", "disable")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s", username, password, host, port, database, sslmode, schema)
	} else if !strings.Contains(connStr, "sslmode=") {
		sep := "?"
		if strings.Contains(connStr, "?") {
			sep = "&"
		}
		connStr = connStr + sep + "sslmode=" + sslmode
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if logDBConnections() {
		log.Println("Connected to database")
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Close() error {
	if logDBConnections() {
		log.Printf("Disconnected from database: %s", database)
	}
	return s.db.Close()
}

func (s *service) DB() *sql.DB {
	return s.db
}

func logDBConnections() bool {
	val := strings.ToLower(os.Getenv("DB_LOG_CONNECTIONS"))
	return val == "1" || val == "true" || val == "yes" || val == "on"
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
