package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	Camera  string `yaml:"camera" json:"camera"`
	Default int    `yaml:"default" json:"default"`
}

type Timer struct {
	Time int `json:"duration"  binding:"required"`
}

var Gate bool
var Time int
var Runnign bool

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func GetData(c *gin.Context) {
	session := sessions.Default(c)
	var logged bool
	v := session.Get("logged")
	if v != nil {
		logged = v.(bool)
	}

	if !logged {
		c.JSON(http.StatusUnauthorized, gin.H{"messge": "Please login"})
	} else {
		var cnf conf
		cnf.getConf()
		c.JSON(http.StatusOK, gin.H{"message": "Data", "data": gin.H{"camera": cnf.Camera, "default": cnf.Default, "gate": Gate, "running": Runnign, "time": Time}})
	}

	fmt.Println(logged)
}

func Timing(c *gin.Context) {
	var t Timer
	if err := c.ShouldBind(&t); err == nil {
		go runTimer(t.Time)
		c.JSON(http.StatusOK, gin.H{"message": "Running timer"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func runTimer(t int) {
	Time = t
	Runnign = true
	Gate = true
	timeout := time.After(time.Duration(t) * time.Second)
	tick := time.Tick(1 * time.Second)
	for {
		if Runnign {
			select {
			case <-timeout:
				Gate = false
				Runnign = false
				SetGate()
			case <-tick:
				Time -= 1
			}
		}
	}

}

func Manual(c *gin.Context) {
	Gate = !Gate
	// var cnf conf
	// cnf.getConf()
	c.JSON(http.StatusOK, gin.H{"message": "Data", "data": gin.H{"gate": Gate}})
}

func SetGate() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "8")

	work := func() {
		gobot.Every(1*time.Second, func() {
			if Gate {
				led.On()
			} else {
				led.Off()
			}
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
