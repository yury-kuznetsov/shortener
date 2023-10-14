package config

import (
	"flag"
	"os"
)

var Options struct {
	HostAddr string
	BaseAddr string
}

func Init() {
	initFlags()
	initEnv()
}

func initFlags() {
	flag.StringVar(&Options.HostAddr, "a", ":8080", "TCP network address")
	flag.StringVar(&Options.BaseAddr, "b", "http://localhost:8080", "base address")
	flag.Parse()
}

func initEnv() {
	if envHostAddr := os.Getenv("SERVER_ADDRESS"); envHostAddr != "" {
		Options.HostAddr = envHostAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		Options.BaseAddr = envBaseAddr
	}
}
