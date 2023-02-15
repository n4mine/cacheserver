package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var C *Config

type Config struct {
	Web   WebConfig   `toml:"web"`
	Rpc   RpcConfig   `toml:"rpc"`
	Cache CacheConfig `toml:"cache"`
	GC    GcConfig    `toml:"gc"`
}

type WebConfig struct {
	Enable bool   `toml:"enable"`
	Port   string `toml:"port"`
}

type RpcConfig struct {
	Enable bool   `toml:"enable"`
	Port   string `toml:"port"`
}

type CacheConfig struct {
	SpanInSeconds int `toml:"span_in_seconds"`
	NumOfChunks   int `toml:"num_of_chunks"`
}

type GcConfig struct {
	ExpiresInMinutes    int `toml:"expires_in_minutes"`
	GcIntervalInMinutes int `toml:"gc_interval_in_minutes"`
}

func LoadConfig(path string) *Config {
	var config Config
	var bs []byte
	var err error

	if bs, err = os.ReadFile(path); err != nil {
		fmt.Fprintf(os.Stderr, "read config file failed: %s\n", err.Error())
		os.Exit(1)
	}

	if _, err = toml.Decode(string(bs), &config); err != nil {
		fmt.Fprintf(os.Stderr, "decode config file failed: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("load config from %s:\n%+v\n", path, config)

	C = &config

	return &config
}
