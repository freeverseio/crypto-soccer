package relay

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
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

func verifyAction(userAddr string, atype string, value string, r string, s string, v int) bool {
	hexUserAddr := common.HexToAddress(userAddr)

	typeBytes, err := str2bytes(atype)
	if err != nil {
		fmt.Println("Error decoding type: ", err)
		return false
	}

	valueBytes, err := str2bytes(value)
	if err != nil {
		fmt.Println("Error decoding type: ", err)
		return false
	}

	rBytes, err := str2bytes(r)
	if err != nil {
		fmt.Println("Error decoding R: ", err)
		return false
	}

	sBytes, err := str2bytes(s)
	if err != nil {
		fmt.Println("Error decoding S: ", err)
		return false
	}

	msg := append(typeBytes[:], valueBytes[:]...)
	sig := append(rBytes[:], sBytes[:]...)
	sig = append(sig[:], byte(v))

	return verifySigEthMsg(hexUserAddr, sig, msg)
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

func bytes2str(data []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(data[:]))
}

func hashBytes(b ...[]byte) (hash [32]byte) {
	h := crypto.Keccak256(b...)
	copy(hash[:], h)
	return hash
}

// Signature is a secp256k1 ecdsa signature.
type Signature [65]byte

// SignatureEthMsg is a secp256k1 ecdsa signature of an ethereum message:
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sig://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
type SignatureEthMsg [65]byte

// Hash - alias for [32]byte
type Hash [32]byte

// ethHash is the hashing function used before signing ethereum messages.
func ethHash(b []byte) Hash {
	const web3SignaturePrefix = "\x19Ethereum Signed Message:\n"
	header := fmt.Sprintf("%s%d", web3SignaturePrefix, len(b))
	return hashBytes([]byte(header), b)
}

func verifySigEthMsg(addr common.Address, sig []byte, msg []byte) bool {
	hash := ethHash(msg)
	var _sig SignatureEthMsg
	copy(_sig[:], sig[:])
	_sig[64] -= 27
	return verifySig(addr, (*Signature)(&_sig), hash[:])
}

// verifySig verifies a given signature and the msgHash with the expected address
func verifySig(addr common.Address, sig *Signature, msgHash []byte) bool {
	recoveredPub, err := crypto.Ecrecover(msgHash, sig[:])
	if err != nil {
		fmt.Printf("ECRecover error: %s\n", err)
		return false
	}
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return bytes.Equal(addr.Bytes(), recoveredAddr.Bytes())
}
