package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RelayUserEntry - ...
type RelayUserEntry struct {
	ID    int64
	Nonce uint64
}

var db = make(map[string]*RelayUserEntry) // TODO: use storage.go

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
	entry, ok := db[account]
	if ok {
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
	value, ok := db[user]
	if ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "value": value, "nonce": time.Now().Unix()})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value", "nonce": time.Now().Unix()})
	}
}

// ActionPOST - post action from user (http://localhost:8080/relay/v1/:useraccount/action?type=xyz&value=123)
func ActionPOST(c *gin.Context) {
	user := c.Params.ByName("useraccount")
	action := c.Query("type")
	value := c.Query("value")

	entry, ok := db[user]
	_ = entry
	if ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "action": action, "value": value})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"user": user, "message": "user not found"})
	}
}

// CreateUserPOST - adds user to db (http://localhost:8080/relay/createuser?user=xyz)
func CreateUserPOST(c *gin.Context) {
	user := c.Query("user")

	entry := GetUserEntry(user)
	if entry != nil {
		c.JSON(http.StatusBadRequest, gin.H{"user": user, "message": "user already in exists"})
	} else {
		AddUserEntry(user)
		c.JSON(http.StatusCreated, gin.H{"user": user, "message": "user created"})
	}
}

func main() {

	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// GET
	router.GET("/ping", pingGET)
	router.GET("/relay/v1/:useraccount/nonce", NonceGET)

	// POST
	router.POST("/relay/createuser", CreateUserPOST) // TODO: just for debugging
	router.POST("/relay/v1/:useraccount/action", ActionPOST)

	// Listen and Server in localhost:8080
	router.Run(":8080")
}
