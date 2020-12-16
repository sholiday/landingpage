package landingpage

import (
	"fmt"
	"sort"
)

type ServerConfig struct {
	Host  string
	Port  int
	Title string
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
