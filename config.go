package main

import (
	"encoding/json"
	"path/filepath"
	"fmt"
	"os"
)

type Config struct {
	AppPath		string	`json:"app_path"`
	LogFile		string	`json:"log_file"`
	OutFile		string	`json:"out_file"`
	//stateFile	string	`json:"state_file"`
	LegacyLogFile	string	`json:"legacy_log_file"`
	LegacyOutFile	string	`json:"legacy_out_file"`
}

func GetConfig(cf string) (*Config, error) {

	var c *Config

	fn := fmt.Sprintf("%s%c%s", filepath.Dir(os.Args[0]), filepath.Separator, cf)

	fh, e := os.Open(fn)
	defer fh.Close()

	if e == nil {
		d := json.NewDecoder(fh)
		e = d.Decode(&c)
	}

/*
	j, e := ioutil.ReadFile(fn)
	fmt.Println(j)
	fmt.Println(string(j))

	if e == nil {
		e = json.Unmarshal(j, *c)
	}
*/

	return c, e
}
