package main

import (
	"fmt"
)

type config struct {
	dir        string
	dbfilename string
}

var configValues config

func initConfigValues(dir *string, dbfilename *string) config {

	return config{
		dir:        *dir,
		dbfilename: *dbfilename,
	}

}

func (cn *config) encodeConfigValues(configType string) string {

	var configTypeValue string

	if configType == "dir" {
		configTypeValue = cn.dir
	} else {
		configTypeValue = cn.dbfilename
	}

	return fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(configType), configType, len(configTypeValue), configTypeValue)
}
