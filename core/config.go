package core

import (
	"github.com/TukTuk/common"
)

var CF *Config

type Config struct {
	Maps GMaps
}

type GMaps struct {
	Hostname string
	ApiKey   string
}

func InitConfig() *Config {
	CF = &Config{
		Maps: GMaps{
			Hostname: GMAPS_SERVICE_HOSTNAME,
			ApiKey:   common.API_KEY,
		},
	}
	return CF
}
