package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

// Action - ...
type Action struct {
	Type  string
	Value string
}

// RelayUserEntry - ...
type RelayUserEntry struct {
	ID     int64
	Nonce  uint64
	Action Action
}

var db = make(map[string]*RelayUserEntry) // TODO: use storage.go
var userNotFound = gin.H{"message": "User not found"}

// AddUserEntry - adds user to db
func AddUserEntry(account string) error {
	_, ok := db[account]
	if ok {
		return errors.New("User already exist")
	}
	db[account] = &RelayUserEntry{ID: time.Now().Unix(), Nonce: 0}
	return nil
}

// GetUserEntry - adds user to db
func GetUserEntry(account string) *RelayUserEntry {
	if entry, ok := db[account]; ok == true {
		return entry
	}
	return nil
}

func pingGET(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// NonceGET - get user nonce (http://localhost:8080/relay/v1/1234/nonce)
func NonceGET(c *gin.Context) {
	user := c.Param("useraccount")
	if entry := GetUserEntry(user); entry != nil {
		c.JSON(http.StatusOK, gin.H{"useraccount": user, "nonce": entry.Nonce})
		return
	}
	c.JSON(http.StatusBadRequest, userNotFound)
}

// ActionPOST - post action from user (http://localhost:8080/relay/v1/:useraccount/action?type=xyz&value=123)
func ActionPOST(c *gin.Context) {
	user := c.Params.ByName("useraccount")

	var body struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println("Error binding to json:", err)
		return
	}

	if entry := GetUserEntry(user); entry != nil {
		entry.Action = Action{body.Type, body.Value}
		c.JSON(http.StatusOK, gin.H{"user": user, "action": entry.Action})
		return
	}
	c.JSON(http.StatusBadRequest, userNotFound)
}

// CreateUserPOST - adds user to db (http://localhost:8080/relay/createuser?user=xyz)
func CreateUserPOST(c *gin.Context) {
	//user := c.Query("user")

	var body struct {
		User string `json:"user"`
	}
	err := c.ShouldBindJSON(&body)
	user := body.User

	if err != nil {
		fmt.Println("Error binding to json:", err)
		return
	}

	if entry := GetUserEntry(user); entry != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already in exists"})
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
	router.GET("/relay/v1/:useraccount/nonce", NonceGET)

	// POST
	router.POST("/relay/v1/:useraccount/action", ActionPOST)
	router.POST("/relay/createuser", CreateUserPOST)

	// DEBUG
	router.GET("/ping", pingGET)
	router.GET("/relay/db", dbGET)

	// Listen and Server in localhost:8080
	router.Run(":8080")
}

func dbGET(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"db": db})
}
