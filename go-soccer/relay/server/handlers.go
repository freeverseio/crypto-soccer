package relay

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

var userNotFound = gin.H{"message": "User not found"}

// NonceGET - get user nonce (/relay/v1/:useraddr/nonce)
func NonceGET(c *gin.Context) {
	user := c.Param("useraddr")
	if entry := GetUserEntry(user); entry != nil {
		c.JSON(http.StatusOK, gin.H{"useraddr": user, "nonceHex": nonce2Hex(entry.Nonce)})
		return
	}
	c.JSON(http.StatusOK, userNotFound)
}

// ActionPOST - post action from user (/relay/v1/:useraddr/action)
func ActionPOST(c *gin.Context) {
	useraddr := c.Params.ByName("useraddr")

	var body struct {
		From  string `json:"from"` // TODO: redundant and could be removed in production
		Type  string `json:"type"`
		Value string `json:"value"`
		Msg   string `json:"msg"` // TODO: redundant, currently just for debugging, can be removed in production
		R     string `json:"r"`
		S     string `json:"s"`
		V     int    `json:"v"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "failed to parse action request"})
		log.Println("Error: failed to parse action request")
		return
	}

	typeBytes, err := hex.DecodeString(body.Type)
	if err != nil {
		log.Println("Error decoding type: ", err)
		return
	}
	valueBytes, err := hex.DecodeString(body.Value)
	if err != nil {
		log.Println("Error decoding value: ", err)
		return
	}
	msgBytes, err := str2bytes(body.Msg)
	if err != nil {
		log.Println("Error decoding msg: ", err)
		return
	}

	typeStr := string(typeBytes)
	valueStr := string(valueBytes)

	nonce := uint64(0)
	if entry := GetUserEntry(useraddr); entry != nil {
		nonce = entry.Nonce
	}

	hashMsg := ethHashMsg(typeStr + valueStr + nonce2Hex(nonce))

	if !bytes.Equal(hashMsg.Bytes(), msgBytes[:]) {
		log.Println("Hash message differs")
		return
	}

	ok := verifyEthMsg(useraddr, hashMsg, body.R, body.S, body.V)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "failed to verify message"})
		return
	}

	if entry := GetUserEntry(useraddr); entry != nil {
		entry.Action = Action{typeStr, valueStr}
		entry.Nonce++
		c.JSON(http.StatusOK, gin.H{"useraddr": useraddr, "action": entry.Action, "verified": ok})
		processAction(typeStr, valueStr)
	}
}

// CreateUserPOST - adds user to db (/relay/createuser)
func CreateUserPOST(c *gin.Context) {

	var body struct {
		User string `json:"useraddr"`
	}

	err := c.ShouldBindJSON(&body)
	user := body.User

	if err != nil {
		log.Println("Error binding to json:", err)
		return
	}

	if entry := GetUserEntry(user); entry != nil {
		c.JSON(http.StatusOK, gin.H{"message": "User already exists"})
		return
	}
	AddUserEntry(user)
	c.JSON(http.StatusCreated, gin.H{"user": user, "message": "User created"})
}

func processAction(actionType string, actionValue string) {
	// TODO: interact with lionel
}
func dbGET(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"db": db})
}
func nonce2Hex(nonce uint64) string {
	return fmt.Sprintf("%064x", nonce)
}
func str2bytes(s string) ([32]byte, error) {
	bytes, err := hex.DecodeString(s)
	var result [32]byte
	if err != nil {
		return result, err
	}
	copy(result[:], bytes)
	return result, nil
}
func ethHashMsg(msg string) common.Hash {
	const web3SignaturePrefix = "\x19Ethereum Signed Message:\n"
	data := []byte(msg)
	header := fmt.Sprintf("%s%d", web3SignaturePrefix, len(msg))
	return crypto.Keccak256Hash([]byte(header), data)
}
func verifyEthMsg(from string, hash common.Hash, r string, s string, v int) bool {
	rBytes, err := str2bytes(r)
	if err != nil {
		log.Println(err)
		return false
	}
	sBytes, err := str2bytes(s)
	if err != nil {
		log.Println(err)
		return false
	}
	sig := append(rBytes[:], sBytes[:]...)
	// recovery id is either 27 or 28. Remove 27 so it becomes either 0 or 1.
	sig = append(sig[:], byte(v-27))

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		log.Println(err)
		return false
	}

	pubKey, _ := crypto.UnmarshalPubkey(sigPublicKey)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(common.HexToAddress(from).Bytes(), recoveredAddr.Bytes()) {
		log.Println("recovered address: ", recoveredAddr.Hex(), " does not match ", from)
		return false
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		log.Println(err)
		return false
	}
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	return crypto.VerifySignature(sigPublicKeyBytes, hash.Bytes(), sig[:len(sig)-1])
}
