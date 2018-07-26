// +build linux darwin/arm

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/szmagli/remote-gate/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/szmagli/remote-gate/users"
)

func flags() {
	var login string
	var password string
	var serve bool
	flag.StringVar(&login, "login", "", "Set login for administator")
	flag.StringVar(&password, "password", "", "Set password for administator")
	flag.BoolVar(&serve, "serve", false, "Serve app on web")
	flag.Parse()

	if login == "" && password == "" && !serve {
		flag.PrintDefaults()
	} else if login != "" && password != "" {
		fmt.Println("Creating user")
		users.CreateUser(login, password, 0)
	} else if serve {
		fmt.Println("Starting application")
		serveApp()
	} else {
		flag.PrintDefaults()
	}
}

func main() {
	flags()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
}

func serveApp() {
	r := gin.Default()
	// api := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("remoteGate", store))
	r.Static("/", "./web")

	v1 := r.Group("/v1")
	{
		v1.POST("/login", users.Login)
		v1.POST("/manual", controller.Manual)
		v1.POST("/timing", controller.Timing)

		v1.POST("/status", controller.GetData)
	}

	go controller.SetGate()
	r.Use(CORSMiddleware())
	r.Run()

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
