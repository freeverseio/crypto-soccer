package authproxy

import (
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func hashPrefixedMessage(data []byte) []byte {
	prefixed := []byte{0x19}
	prefixed = append(prefixed, []byte(fmt.Sprintf("Ethereum Signed Message:\n%v", len(data)))...)
	prefixed = append(prefixed, data...)
	return prefixed
}
func SignToken(privateKey *ecdsa.PrivateKey, t time.Time) (string, error) {
	ts := fmt.Sprintf("%v", t.Unix())
	hash := crypto.Keccak256(hashPrefixedMessage([]byte(ts)))
	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to sign token %w", err)
	}
	token := ts + ":" + base64.StdEncoding.EncodeToString(sig)
	return token, nil
}
func VerifyToken(token string, grace time.Duration) (common.Address, time.Time, error) {
	// token have: unix_time:signature
	tokenfields := strings.Split(token, ":")
	if len(tokenfields) != 2 {
		return common.Address{}, time.Time{}, errors.New("malformed token:" + token)
	}
	// check date
	tsunix, err := strconv.Atoi(tokenfields[0])
	if err != nil {
		return common.Address{}, time.Time{}, fmt.Errorf("malformed token timestamp %w", err)
	}
	ts, now := time.Unix(int64(tsunix), 0), time.Now()
	if math.Abs(now.Sub(ts).Seconds()) > grace.Seconds() {
		return common.Address{}, time.Time{}, errors.New("token out of time")
	}
	// retrieve public key
	hash := crypto.Keccak256(hashPrefixedMessage([]byte(tokenfields[0])))
	sig, err := base64.StdEncoding.DecodeString(tokenfields[1])
	if err != nil {
		return common.Address{}, time.Time{}, fmt.Errorf("malformed token signature encoding : %w", err)
	}
	pbk, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return common.Address{}, time.Time{}, fmt.Errorf("malformed token signature: %w", err)
	}
	addr := crypto.PubkeyToAddress(*pbk)
	return addr, ts, nil
}

func MatchTransferFirstBotMutation(data string) (bool, error) {
	ex := `mutation(\s*).*(\s*){(\s*)transferFirstBotToAddr(\s*)\((\s*)timezone(\s*):(\s*)\d{1,2}(\s*),(\s*)countryIdxInTimezone(\s*):(\s*)[0-9]+(\s*),(\s*)address(\s*):(\s*)"[a-zA-Z0-9]+"(\s*)\)(\s*)}`
	return regexp.MatchString(ex, data)
}
