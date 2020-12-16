package main

import (
	"fmt"
	"html/template"
	"log"

	"github.com/sholiday/landingpage"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	viper.SetEnvPrefix("LANDINGPAGE")
	viper.SetConfigName("landingpage")
	viper.AddConfigPath("$HOME/.config/landingpage/")
	viper.AddConfigPath("/etc/landingpage/")
	viper.AddConfigPath("/config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}
	var config landingpage.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}

	lp := landingpage.New(config)

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"startRow": func(i int, j int) bool { return i%j == 0 },
		"endRow":   func(i int, j int) bool { return (i+1)%j == 0 },
	})
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", lp.Main)
	hostPort := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Println("Serving at", hostPort)
	r.Run(hostPort)
}
