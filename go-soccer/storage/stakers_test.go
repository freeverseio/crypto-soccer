package storage

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func TestStaker(t *testing.T) {
	s, err := New("./test");
	if err != nil{
		t.Error(err)
	}

	defer os.RemoveAll("./test")

	address := common.HexToAddress("0x44")
	_, err = s.Staker(address)
	if err != errors.ErrNotFound {
		t.Error(err)
	}
}

func TestHasStaker(t* testing.T) {
	s, err := New("./test");
	if err != nil{
		t.Error(err)
	}

	defer os.RemoveAll("./test")

	address := common.HexToAddress("0x44")
	has, err := s.HasStaker(address)
	if err != nil {
		t.Error(err)
	}

	if has == true {
		t.Error("Has staker in empty db")
	}

}