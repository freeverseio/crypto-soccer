package storage

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func TestStaker(t *testing.T) {
	s, err := New("./test");
	defer os.RemoveAll("./test")
	if err != nil{
		t.Error(err)
	}

	address := common.HexToAddress("0x44")
	_, err = s.Staker(address)
	if err != errors.ErrNotFound {
		t.Error(err)
	}
}

func TestHasStaker(t* testing.T) {
	s, err := New("./test");
	defer os.RemoveAll("./test")
	if err != nil{
		t.Error(err)
	}

	address := common.HexToAddress("0x44")
	has, err := s.HasStaker(address)
	if err != nil {
		t.Error(err)
	}

	if has == true {
		t.Error("Has staker in empty db")
	}

}

func TestSetStaker(t* testing.T) {
	s, err := New("./test")
	defer os.RemoveAll("./test")
	if err != nil{
		t.Error(err)
	}

	entry := StakerEntry{45}
	address := common.HexToAddress("0x44")
	err = s.SetStaker(address, &entry)
	if err != nil {
		t.Error(err)
	}

	result, err := s.Staker(address)
	if err != nil {
		t.Error(err)
	}
	if *result != entry {
		t.Error("Expected is ", entry, " and result is ", *result)
	}

	entry2 := StakerEntry{1}
	err = s.SetStaker(address, &entry2)
	if err != nil {
		t.Error(err)
	}

	result, err = s.Staker(address)
	if err != nil {
		t.Error(err)
	}
	if *result != entry2 {
		t.Error("Expected is ", entry, " and result is ", *result)
	}


}