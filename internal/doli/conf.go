package doli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	baseURL string
	doliKey string
)

func init() {
	baseURL = os.Getenv("DOLAPIURL")
	doliKey = os.Getenv("DOLAPIKEY")
	if baseURL == "" || doliKey == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Errorf("cannot determine home directory: %w", err))
		}
		conf, err := readConf(filepath.Join(home, ".doliconf"))
		if err != nil {
			panic(fmt.Errorf("cannot read ~/.doliconf: %w", err))
		}
		if baseURL == "" {
			baseURL = conf["url"]
		}
		if doliKey == "" {
			doliKey = conf["apikey"]
		}
	}
	if baseURL == "" {
		panic("dolibarr URL not set: provide url in ~/.doliconf or DOLAPIURL env var")
	}
	if doliKey == "" {
		panic("dolibarr API key not set: provide apikey in ~/.doliconf or DOLAPIKEY env var")
	}
}

func readConf(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := make(map[string]string)
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		conf[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return conf, nil
}
