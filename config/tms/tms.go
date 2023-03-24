package tms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Path string
	Data map[string]interface{}
}

func (c *Config) Setup() error {
	jsonFile, err := os.Open(c.Path)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	fmt.Println("Config : ", result)
	c.Data = result

	return nil
}

func (c *Config) RepoConfig() map[string]interface{} {
	return map[string]interface{}{
		"connection_string": fmt.Sprint(c.Data["database_url"]),
	}
}

func (c *Config) GisConfig() map[string]interface{} {
	return map[string]interface{}{
		"connection_string": fmt.Sprint(c.Data["gis_database_url"]),
	}
}
