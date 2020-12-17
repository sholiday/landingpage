package landingpage

import (
	"fmt"
	"sort"

	"github.com/spf13/viper"
)

func LoadConfig() (Config, error) {
	var config Config
	viper.SetEnvPrefix("LANDINGPAGE")
	viper.SetConfigName("landingpage")
	viper.AddConfigPath("$HOME/.config/landingpage/")
	viper.AddConfigPath("/etc/landingpage/")
	viper.AddConfigPath("/config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if viper.Get("server.port") == nil {
		viper.Set("server.port", 8081)
	}
	if viper.Get("server.title") == nil {
		viper.Set("server.title", "My Landing Page")
	}
	if viper.Get("server.userheader") == nil {
		viper.Set("server.userheader", "X-Forwarded-User")
	}
	if err != nil {
		return config, fmt.Errorf("Failed to read config file: %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("Failed to unmarshal config file: %s", err)
	}
	return config, err
}

type ServerConfig struct {
	Host       string
	Port       int
	Title      string
	UserHeader string
}

type UserConfig struct {
	Name      string
	Usernames []string
	Groups    []string
}

type AppConfig struct {
	Title       string
	Description string
	Url         string
	Logo        string
	Groups      []string
}

type Config struct {
	Server ServerConfig
	Users  []UserConfig
	Apps   map[string]AppConfig
}

func (a *AppConfig) IsVisibleToGroup(group string) bool {
	for _, g := range a.Groups {
		if g == group {
			return true
		}
	}
	return false
}

// This function is terrible.
func (c *Config) AppsForUsername(username string) ([]string, error) {
	u, err := c.ConfigForUsername(username)
	if err != nil {
		return nil, err
	}
	appSet := make(map[string]bool)
	for _, group := range u.Groups {
		for appKey, appConfig := range c.Apps {
			if appConfig.IsVisibleToGroup(group) {
				appSet[appKey] = true
			}
		}
	}
	appList := make([]string, len(appSet))
	i := 0
	for k, _ := range appSet {
		appList[i] = k
		i++
	}
	sort.Strings(appList)
	return appList, nil
}

// Retrieve the config data for a given username.
func (c *Config) ConfigForUsername(username string) (UserConfig, error) {
	for _, u := range c.Users {
		for _, n := range u.Usernames {
			if n == username {
				return u, nil
			}
		}
	}
	return UserConfig{}, fmt.Errorf("Username '%s' not found.", username)
}
