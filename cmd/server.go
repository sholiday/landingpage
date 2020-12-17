package main

import (
	"fmt"
	"log"

	"github.com/sholiday/landingpage"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config, err := landingpage.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	lp := landingpage.New(config)

	r := gin.Default()
	r.SetFuncMap(lp.GetFuncMap())
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", lp.Main)
	hostPort := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Println("Serving at", hostPort)
	r.Run(hostPort)
}
