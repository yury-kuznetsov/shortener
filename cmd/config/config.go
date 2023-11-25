package config

import (
	"flag"
	"os"
)

var Options struct {
	HostAddr string
	BaseAddr string
	FilePath string
	Database string
}

func Init() {
	initFlags()
	initEnv()
}

func initFlags() {
	// dnsForExample := "host=localhost user=shortener password=shortener dbname=postgres sslmode=disable"
	flag.StringVar(&Options.HostAddr, "a", ":8080", "TCP network address")
	flag.StringVar(&Options.BaseAddr, "b", "http://localhost:8080", "base address")
	flag.StringVar(&Options.FilePath, "f", "/tmp/short-url-db.json", "storage path")
	flag.StringVar(&Options.Database, "d", "", "database dsn")
	flag.Parse()
}

func initEnv() {
	if envHostAddr := os.Getenv("SERVER_ADDRESS"); envHostAddr != "" {
		Options.HostAddr = envHostAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		Options.BaseAddr = envBaseAddr
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		Options.FilePath = envFilePath
	}
	if envDatabase := os.Getenv("DATABASE_DSN"); envDatabase != "" {
		Options.Database = envDatabase
	}
}
