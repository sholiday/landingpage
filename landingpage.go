package landingpage

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func New(c Config) *App {
	return &App{
		config:    c,
		debugMode: os.Getenv("LANDINGPAGE_MODE") == "debug",
	}
}

type App struct {
	config    Config
	debugMode bool
}

func (a *App) GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"startRow": func(i int, j int) bool { return i%j == 0 },
		"endRow":   func(i int, j int) bool { return (i+1)%j == 0 },
	}
}

func (a *App) determineUser(c *gin.Context) string {
	user := c.Request.Header.Get(a.config.Server.UserHeader)
	if a.debugMode && len(c.Query("user")) != 0 {
		userOverride := c.Query("user")
		log.Printf("Overriding user '%s' with '%s'\n", user, userOverride)
		user = userOverride
	}
	return user
}

func (a *App) Main(c *gin.Context) {
	userString := a.determineUser(c)
	log.Printf("User requested: %s", userString)
	appKeys, err := a.config.AppsForUsername(userString)
	if err != nil {
		c.String(http.StatusForbidden, "Username not registered.")
		return
	}
	user, err := a.config.ConfigForUsername(userString)
	if err != nil {
		c.String(http.StatusForbidden, "Username not registered.")
		return
	}

	var apps []AppConfig
	for _, appKey := range appKeys {
		apps = append(apps, a.config.Apps[appKey])
	}
	c.HTML(http.StatusOK, "main.tmpl", gin.H{
		"title": a.config.Server.Title,
		"user":  user,
		"apps":  apps,
	})
}

func (a *App) render(c *gin.Context) {
	c.HTML(http.StatusOK, "main.tmpl", gin.H{
		"title":   a.config.Server.Title,
		"message": "foo",
	})
}
