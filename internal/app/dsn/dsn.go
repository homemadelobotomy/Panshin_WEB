package dsn

import (
	"fmt"
	"os"
)

func FromEnv() string {
	host := os.Getenv("DB_HOST")

	if host == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

func GetMinioURL(filename string) string {
	host := os.Getenv("MINIO_HOST")
	if host == "" {
		return ""
	}

	port := os.Getenv("MINIO_PORT")

	return fmt.Sprintf("http://%s:%s/%s/%s", host, port, "images", filename)
}
