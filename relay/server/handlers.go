package relay

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

var userNotFound = gin.H{"message": "User not found"}

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
		Msg   string `json:"msg"`
		R     string `json:"r"`
		S     string `json:"s"`
		V     int    `json:"v"`
	}
	c.ShouldBindJSON(&body)
	fmt.Println("body", body)

	// assert that user equals body.From
	// TODO: Get the last user nonce from db, and verify if recieved nonce is the same -> ask adria: which received nonce? ethclient.PendingNonceAt?
	// TODO: make sure V is 27 (just as a double check, but we actually do not need it)

	//ok := verifyEthMsg(user, ethHashMsg(body.Type+body.Value), body.R, body.S)
	ok := verifyEthMsg(user, common.HexToHash(body.Msg), body.R, body.S)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "failed to verify message"})
		return
	}
	// TODO: If success, increment the user nonce and store in the database -> ask adria do we need to do any transaction first?
	typeBytes, _ := hex.DecodeString(body.Type)
	valueBytes, _ := hex.DecodeString(body.Value)
	message := string(typeBytes) + " " + string(valueBytes)

	if entry := GetUserEntry(user); entry != nil {
		entry.Action = Action{message}
		c.JSON(http.StatusOK, gin.H{"user": user, "action": entry.Action, "verified": ok})
	}
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

func dbGET(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"db": db})
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
	header := fmt.Sprintf("%s%d", web3SignaturePrefix, len(data))
	return crypto.Keccak256Hash([]byte(header), data)
}
func verifyEthMsg(from string, hash common.Hash, r string, s string) bool {
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
	sig = append(sig[:], byte(0))

	// remove ethereum quirk
	//ethPostfix := byte(27)
	//sig[64] -= ethPostfix

	log.Println("verifying message: ", hash.Hex())
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		log.Println(err)
		return false
	}

	pubKey, _ := crypto.UnmarshalPubkey(sigPublicKey)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(common.HexToAddress(from).Bytes(), recoveredAddr.Bytes()) {
		log.Println("reovered address: ", recoveredAddr.Hex(), " does not match ", from)
		return false
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		log.Println(err)
		return false
	}
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	fmt.Println("public key: ", hexutil.Encode(sigPublicKeyBytes[1:]))
	return crypto.VerifySignature(sigPublicKeyBytes, hash.Bytes(), sig[:len(sig)-1])
}
