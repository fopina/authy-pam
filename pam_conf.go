package main

import (
	"os"
	"github.com/pelletier/go-toml"
)

type AuthyConfig struct {
	URL 	string	`toml:"url,omitempty"`
	Token 	string	`toml:"token,omitempty"`
}

type Config struct {
    Authy	AuthyConfig			`toml:"authy,omitempty"`
    Users 	map[string]string	`toml:"users,omitempty"`
}

func (c *Config) LoadFromFile(path string) error {
	t, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	err = t.Unmarshal(c)
	if err != nil {
		return err
	}
	// TODO once https://github.com/pelletier/go-toml/issues/252 is fixed, change this
	if c.Authy.URL == "" {
		c.Authy.URL = "https://api.authy.com"
	}
	return nil
}

func (c *Config) SaveToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

/*
func main() {
	config := Config{}
	err := toml.Unmarshal([]byte(`
		[authy]
			url = "x"
		[users]
			name = "roger"
			authy = "123"
	`), &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Authy.URL)
	fmt.Println(config.Users)
	x, _ :=  toml.Marshal(&config)
	fmt.Println("----")
	fmt.Println(string(x))

	err = config.LoadFromFile("data.conf")
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Authy)
	if config.Users == nil {
    	config.Users = map[string]string{}
	}
	config.Users["newone"] = "test"
	config.SaveToFile("data2.conf")
}
*/