package config

import "flag"

var Options struct {
	HostAddr string
	BaseAddr string
}

func Init() {
	initFlags()
}

func initFlags() {
	flag.StringVar(&Options.HostAddr, "a", ":8080", "TCP network address")
	flag.StringVar(&Options.BaseAddr, "b", "http://localhost:8080", "base address")
	flag.Parse()
}
