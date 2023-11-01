package config

import (
	"flag"
	"os"
)

var Options struct {
	HostAddr string
	BaseAddr string
	FilePath string
}

func Init() {
	initFlags()
	initEnv()
}

func initFlags() {
	flag.StringVar(&Options.HostAddr, "a", ":8080", "TCP network address")
	flag.StringVar(&Options.BaseAddr, "b", "http://localhost:8080", "base address")
	flag.StringVar(&Options.FilePath, "f", "/tmp/short-url-db.json", "storage path")
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
}
