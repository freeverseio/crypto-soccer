package relay

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	//common3 "github.com/iden3/go-iden3/utils"
)

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

// ActionPOST - post action from user (http://localhost:8080/relay/v1/:useraddr/action)
func ActionPOST(c *gin.Context) {
	user := c.Params.ByName("useraddr")

	var body struct {
		From  string `json:"from"`
		Type  string `json:"type"`
		Value string `json:"value"`
		R     string `json:"r"`
		S     string `json:"s"`
		V     int    `json:"v"`
	}
	c.ShouldBindJSON(&body)
	fmt.Println("body", body)

	// assert that user equals body.From
	// TODO: Get the last user nonce from db, and verify if recieved nonce is the same -> ask adria: which received nonce? ethclient.PendingNonceAt?

	ok := verifyAction(body.From, body.Type, body.Value, body.R, body.S, body.V)
	// TODO: If success, increment the user nonce and store in the database -> ask adria do we need to do any transaction first?

	if entry := GetUserEntry(user); entry != nil {
		entry.Action = Action{body}
		c.JSON(http.StatusOK, gin.H{"user": user, "action": entry.Action, "verified": ok})
		return
	}
	c.JSON(http.StatusOK, userNotFound)
}

func verifyAction(userAddr string, atype string, value string, r string, s string, v int) bool {
	hexUserAddr := common.HexToAddress(userAddr)

	typeBytes, err := str2byte(atype)
	if err != nil {
		fmt.Println("Error decoding type: ", err)
		return false
	}

	valueBytes, err := str2byte(value)
	if err != nil {
		fmt.Println("Error decoding type: ", err)
		return false
	}

	rBytes, err := str2byte(r)
	if err != nil {
		fmt.Println("Error decoding R: ", err)
		return false
	}

	sBytes, err := str2byte(s)
	if err != nil {
		fmt.Println("Error decoding S: ", err)
		return false
	}

	msg := append(typeBytes[:], valueBytes[:]...)
	sig := append(rBytes[:], sBytes[:]...)
	sig = append(sig[:], byte(v))

	_ = hexUserAddr
	_ = msg
	// TODO: ask adria about msg. I also get an error when installing iden3/go-iden3/utils (error loading module requirements)
	//return common3.VerifySigEthMsg(hexUserAddr, sig, msgBytes)
	return true
}

func str2byte(s string) ([32]byte, error) {
	bytes, err := hex.DecodeString(s)
	var result [32]byte
	if err != nil {
		return result, err
	}
	copy(result[:], bytes)
	return result, nil
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

// Server - ...
type Server struct {
}

// Start - starts the server
func (s *Server) Start(port string) {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// GET
	router.GET("/relay/v1/:useraddr/nonce", NonceGET)

	// POST
	router.POST("/relay/v1/:useraddr/action", ActionPOST)
	router.POST("/relay/createuser", CreateUserPOST)

	// DEBUG
	router.GET("/relay/db", dbGET)

	// Listen and Server in localhost:8080
	router.Run(port)
}

func dbGET(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"db": db})
}
