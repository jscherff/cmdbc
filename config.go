package main

import (
	"encoding/json"
	"path/filepath"
	"fmt"
	"os"
)

type Config struct {
	AppPath		string `json:"app_path"`
	LogFile		string `json:"log_file"`
	OutFile		string `json:"out_file"`
	InvFile		string `json:"state_file"`
	LegacyLogFile	string `json:"legacy_log_file"`
	LegacyOutFile	string `json:"legacy_out_file"`
	IncludeDefault	bool `json:"include_default"`
	IncludeVID	map[string]bool `json:"include_vid"`
	IncludePID	map[string]map[string]bool `json:"include_pid"`
}

func GetConfig(cf string) (*Config, error) {

	var c *Config

	fn := fmt.Sprintf("%s%c%s", filepath.Dir(os.Args[0]), filepath.Separator, cf)

	fh, e := os.Open(fn)
	defer fh.Close()

	if e == nil {
		jd := json.NewDecoder(fh)
		e = jd.Decode(&c)
	}

	return c, e
}
