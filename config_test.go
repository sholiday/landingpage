package landingpage

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var toml = []byte(`
[server]
host="localhost"
port=8081
title="My Landing Page"

[[users]]
name="Alice"
usernames=["alice@gmail.com", "aperson@example.com"]
groups=["group1", "group2"]

[[users]]
name="Bob"
usernames=["bob@gmail.com"]
groups=["group2", "group3"]

[apps]
  [apps.web1]
  title="Web 1"
  description="My web service"
	url="http://web1.example.com"
	groups=["group1", "group3"]

  [apps.web3]
  title="Web 3"
  description="My web service"
	url="http://web3.example.com"
	groups=["group2", "group3"]
`)

func parseConfig(a *assert.Assertions) (Config, bool) {
	var c Config
	fmt.Printf("%s", toml)
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer(toml))
	if !a.NoError(err) {
		return Config{}, false
	}
	err = viper.Unmarshal(&c)
	if !a.NoError(err) {
		return Config{}, false
	}
	fmt.Printf("%+v", c)
	return c, true
}

func TestParse(t *testing.T) {
	assert := assert.New(t)
	c, ok := parseConfig(assert)
	if !ok {
		return
	}

	// Server
	assert.Equal("localhost", c.Server.Host)
	assert.Equal(8081, c.Server.Port)
	assert.Equal("My Landing Page", c.Server.Title)

	// Users
	if !assert.Equal(2, len(c.Users)) {
		return
	}
	assert.Contains(c.Users[0].Groups, "group1")
	assert.Contains(c.Users[0].Groups, "group2")
	assert.Contains(c.Users[1].Groups, "group2")
	assert.Contains(c.Users[1].Groups, "group3")
}

func TestConfigForUsername(t *testing.T) {
	assert := assert.New(t)
	c, ok := parseConfig(assert)
	if !ok {
		return
	}
	user, err := c.ConfigForUsername("aperson@example.com")
	if !assert.NoError(err) {
		return
	}
	assert.Equal("Alice", user.Name)
}

func TestAppsForUsername(t *testing.T) {
	assert := assert.New(t)
	c, ok := parseConfig(assert)
	if !ok {
		return
	}
	apps, err := c.AppsForUsername("aperson@example.com")
	if !assert.NoError(err) {
		return
	}
	assert.Equal([]string{"web1", "web3"}, apps)
}
