package relay

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server - ...
type Server struct {
}

// Action - ...
type Action struct {
	Value interface{}
}

// UserEntry - ...
type UserEntry struct {
	ID     int64
	Nonce  uint64
	Action Action
}

var db = make(map[string]*UserEntry) // TODO: use storage.go
var userNotFound = gin.H{"message": "User not found"}

// AddUserEntry - adds user to db
func AddUserEntry(account string) error {
	_, ok := db[account]
	if ok {
		return errors.New("User already exist")
	}
	db[account] = &UserEntry{ID: time.Now().Unix(), Nonce: 0}
	return nil
}

// GetUserEntry - adds user to db
func GetUserEntry(account string) *UserEntry {
	if entry, ok := db[account]; ok == true {
		return entry
	}
	return nil
}

// NonceGET - get user nonce (http://localhost:8080/relay/v1/1234/nonce)
func NonceGET(c *gin.Context) {
	user := c.Param("useraddr")
	if entry := GetUserEntry(user); entry != nil {
		c.JSON(http.StatusOK, gin.H{"useraddr": user, "nonce": entry.Nonce})
		return
	}
	c.JSON(http.StatusOK, userNotFound)
}

// ActionPOST - post action from user (http://localhost:8080/relay/v1/:useraddr/action?value=123)
func ActionPOST(c *gin.Context) {
	user := c.Params.ByName("useraddr")

	var body struct {
		//Type  string `json:"type"`
		Value string `json:"value"`
	}
	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println("Error binding to json:", err)
		return
	}

	// TODO: verify action message
	// https://golang.org/pkg/crypto/ecdsa/

	if entry := GetUserEntry(user); entry != nil {
		entry.Action = Action{body.Value}
		c.JSON(http.StatusOK, gin.H{"user": user, "action": entry.Action})
		return
	}
	c.JSON(http.StatusOK, userNotFound)
}

// CreateUserPOST - adds user to db (http://localhost:8080/relay/createuser?useraddr=xyz)
func CreateUserPOST(c *gin.Context) {
	var body struct {
		User string `json:"useraddr"`
	}
	err := c.ShouldBindJSON(&body)
	user := body.User

	if err != nil {
		fmt.Println("Error binding to json:", err)
		return
	}

	if entry := GetUserEntry(user); entry != nil {
		c.JSON(http.StatusOK, gin.H{"message": "User already exists"})
		return
	}
	AddUserEntry(user)
	c.JSON(http.StatusCreated, gin.H{"user": user, "message": "User created"})
}

// Start - starts the server
func (s Server) Start() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// GET
	router.GET("/relay/v1/:useraddr/nonce", NonceGET)

	// POST
	router.POST("/relay/v1/:useraddr/action", ActionPOST)
	router.POST("/relay/createuser", CreateUserPOST)

	// DEBUG
	router.GET("/ping", pingGET)
	router.GET("/relay/db", dbGET)

	// Listen and Server in localhost:8080
	router.Run(":8080")
}

func pingGET(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func dbGET(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"db": db})
}
