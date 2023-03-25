package tms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
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

func (c *Config) RedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprint(c.Data["redis_address"]),
		Password: fmt.Sprint(c.Data["redis_password"]),
		DB:       0, // use default DB
	}
}

func (c *Config) JwtConfig() jwt.JwtConfig {
	return jwt.JwtConfig{
		Secret: []byte(fmt.Sprint(c.Data["jwt_secret"])),
		Salt:   fmt.Sprint(c.Data["jwt_salt"]),
	}
}
