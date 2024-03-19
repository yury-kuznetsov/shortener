package config

import (
	"encoding/json"
	"flag"
	"os"
)

// Options represents the configuration options for the application.
// It contains the following fields:
// - HostAddr: TCP network address
// - BaseAddr: base address
// - FilePath: storage path
// - Database: database DSN
// - Secure: enable HTTPS
// - CfgFile: config file
// - TrustedNet: trusted subnet
var Options struct {
	HostAddr   string
	BaseAddr   string
	FilePath   string
	Database   string
	Secure     bool
	CfgFile    string
	TrustedNet string
}

// Init initializes the application by calling the initFlags and initEnv functions.
func Init() {
	initFlags()
	initEnv()
	initFile()
}

func initFlags() {
	// dnsForExample := "host=localhost user=shortener password=shortener dbname=postgres sslmode=disable"
	flag.StringVar(&Options.HostAddr, "a", ":8080", "TCP network address")
	flag.StringVar(&Options.BaseAddr, "b", "http://localhost:8080", "base address")
	flag.StringVar(&Options.FilePath, "f", "/tmp/short-url-db.json", "storage path")
	flag.StringVar(&Options.Database, "d", "", "database dsn")
	flag.BoolVar(&Options.Secure, "s", false, "enable HTTPS")
	flag.StringVar(&Options.CfgFile, "c", "", "config file")
	flag.StringVar(&Options.TrustedNet, "t", "", "trusted subnet")
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
	if envSecure := os.Getenv("ENABLE_HTTPS"); envSecure != "" {
		Options.Secure = true
	}
	if envCfgFile := os.Getenv("CONFIG"); envCfgFile != "" {
		Options.CfgFile = envCfgFile
	}
	if envTrustNet := os.Getenv("TRUSTED_SUBNET"); envTrustNet != "" {
		Options.TrustedNet = envTrustNet
	}
}

func initFile() {
	if Options.CfgFile == "" {
		return
	}

	file, err := os.ReadFile(Options.CfgFile)
	if err != nil {
		return
	}

	var options struct {
		HostAddr   string `json:"server_address"`
		BaseAddr   string `json:"base_url"`
		FilePath   string `json:"file_storage_path"`
		Database   string `json:"database_dsn"`
		Secure     bool   `json:"enable_https"`
		TrustedNet string `json:"trusted_subnet"`
	}

	err = json.Unmarshal(file, &options)
	if err != nil {
		return
	}

	if Options.HostAddr == "" {
		Options.HostAddr = options.HostAddr
	}
	if Options.BaseAddr == "" {
		Options.BaseAddr = options.BaseAddr
	}
	if Options.FilePath == "" {
		Options.FilePath = options.FilePath
	}
	if Options.Database == "" {
		Options.Database = options.Database
	}
	if !Options.Secure {
		Options.Secure = options.Secure
	}
	if Options.TrustedNet == "" {
		Options.TrustedNet = options.TrustedNet
	}
}
