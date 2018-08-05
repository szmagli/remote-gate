package users

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type (
	User struct {
		Role     int
		Login    string `form:"login" json:"login" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	Users struct {
		Users []User
	}
)

func usr() {
	fmt.Println("Gob Example")
	users := Users{}

	user := User{1, "test", "test"}
	users.Users = append(users.Users, user)
	fmt.Println(users)
	err := writeGob("./data/users.gob", users)
	if err != nil {
		fmt.Println(err)
	}

	var studentRead = new(Users)
	err = readGob("./data/users.gob", studentRead)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("RD")
		fmt.Println(studentRead)
	}
}

// CreateUser with role in gob file
func CreateUser(login string, password string, role int) {
	var usersRead = new(Users)
	err := readGob("./data/users.gob", usersRead)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(usersRead)
	}

	user := User{role, login, password}
	usersRead.Users = append(usersRead.Users, user)

	err = writeGob("./data/users.gob", usersRead)
	if err != nil {
		fmt.Println(err)
	}
}

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBind(&u); err == nil {
		if checkLogin(u) {
			session := sessions.Default(c)
			session.Set("logged", true)
			session.Save()
			c.JSON(http.StatusOK, gin.H{"message": "successfully logged in"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"messge": "Invalid cerdinails"})
		}
	} else {
		session := sessions.Default(c)
		session.Set("logged", false)
		session.Save()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func checkLogin(user User) bool {
	var usersRead = new(Users)
	err := readGob("./data/users.gob", usersRead)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, usr := range usersRead.Users {
		if usr.Login == user.Login {
			if user.Password == usr.Password {
				return true
			}
		}
	}

	return false
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
