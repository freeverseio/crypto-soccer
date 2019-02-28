package stakers

import (
	"fmt"
	"sync"
	"time"

	"github.com/freeverseio/go-soccer/config"
	eth "github.com/freeverseio/go-soccer/eth"

	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/freeverseio/go-soccer/storage"
	log "github.com/sirupsen/logrus"
)

const (
	StateUnenrolled     = 0
	StateEnrolling      = 1
	StateUnenrolling    = 2
	StateUnenrollable   = 3
	StateEnrolled       = 4
	StateChallageTT     = 5
	StateChallangeLIRES = 6
	StateChallangeTTRES = 7
	StateSlashable      = 8
	StateSlashed        = 9
)

const (
	HashOnionLayers = 100
)

type Staker struct {
	Address   common.Address
	Client    *eth.Web3Client
	onionhash common.Hash
}

func (s *Staker) OnionAt(level uint64) common.Hash {
	hash := s.onionhash
	for level > 0 {
		hash = crypto.Keccak256Hash(hash.Bytes())
		level--
	}
	return hash
}

type Stakers struct {
	sync.Mutex
	stks     map[common.Address]*Staker
	storage  *storage.Storage
	contract *eth.Contract
}

func New(stakers []*Staker, sto *storage.Storage) (*Stakers, error) {

	stks := make(map[common.Address]*Staker)
	for i, stk := range stakers {
		stks[stk.Address] = stakers[i]
	}

	stakersAbi, err := abi.JSON(strings.NewReader(stakersAbiJson))
	if err != nil {
		return nil, err
	}
	stakersAddress := common.HexToAddress(config.C.Contracts.StakersAddress)

	contract, err := eth.NewContract(stakers[0].Client, &stakersAbi, nil, &stakersAddress)
	if err != nil {
		return nil, err
	}

	return &Stakers{
		stks:     stks,
		storage:  sto,
		contract: contract,
	}, nil
}

func (s *Stakers) Info() string {
	var info string
	for _, stk := range s.stks {
		info = info + fmt.Sprintf("%s", stk.Address.Hex())

		state, err := s.GetState(stk.Address, 0)
		if err == nil {
			info = info + fmt.Sprintf(" %s", s.State2Str(state))
		} else {
			info = info + fmt.Sprintf(" %s", err.Error())
		}

		hasStacker, err := s.storage.HasStaker(stk.Address)
		if err == nil {
			if hasStacker {
				stackerentry, err := s.storage.Staker(stk.Address)
				if err == nil {
					info = info + fmt.Sprintf(" H%v", stackerentry.HashIndex)
				}
			}
		}
		if err != nil {
			info = info + fmt.Sprintf("%v", err)
		}

		milliether := new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil)

		balance, err := stk.Client.BalanceInfo()
		if err == nil {
			balance = balance.Div(balance, milliether)
			info = info + fmt.Sprintf(" %smΞ", balance.String())
		} else {
			info = info + fmt.Sprintf(" %s", err.Error())
		}

		info = info + "\n"
	}
	return info
}

func (s *Stakers) Get(staker common.Address) *Staker {
	return s.stks[staker]
}

func (s *Stakers) Members() []*Staker {
	members := []*Staker{}
	for _, s := range s.stks {
		members = append(members, s)
	}
	return members
}

func (s *Stakers) NextFreeStaker() (*common.Address, error) {
	s.Lock()
	defer s.Unlock()

	for addr, v := range s.stks {
		state, err := s.GetState(addr, 0)
		if err != nil {
			return nil, err
		}
		if state == StateEnrolled {
			return &v.Address, nil
		}
	}
	return nil, nil
}

func (s *Stakers) NeedsTouch(staker common.Address) (bool, error) {
	state, err := s.GetState(staker, 60*5)
	if err != nil {
		return false, err
	}
	return state == StateSlashable, nil
}

func (s *Stakers) NeedsResolve(staker common.Address) (bool, error) {
	state, err := s.GetState(staker, 0)
	if err != nil {
		return false, err
	}
	return state == StateChallangeLIRES || state == StateChallangeTTRES, nil
}

func (s *Stakers) GetState(staker common.Address, lookup uint64) (uint64, error) {

	now := big.NewInt(time.Now().UTC().Unix() + int64(lookup))

	var state byte
	if err := s.contract.Call(&state, "state", staker, now); err != nil {
		return 0, err
	}
	return uint64(state), nil
}

func (s *Stakers) State2Str(state uint64) string {
	switch state {
	case StateUnenrolled:
		return "Unenrolled"
	case StateEnrolling:
		return "Enrolling"
	case StateUnenrolling:
		return "Unenrolling"
	case StateUnenrollable:
		return "Unenrollable"
	case StateEnrolled:
		return "Enrolled"
	case StateChallageTT:
		return "ChallageTT"
	case StateChallangeLIRES:
		return "ChallangeLIRES"
	case StateChallangeTTRES:
		return "ChallangeTTRES"
	case StateSlashable:
		return "Slashable"
	case StateSlashed:
		return "Slashed"
	}
	panic("UNKNOWNSTATE")
}

func (s *Stakers) Enroll(staker common.Address) error {
	log.Info("Enrolling ", staker.Hex())
	stk := s.stks[staker]

	hasStaker, err := s.storage.HasStaker(stk.Address)
	if err != nil {
		return err
	}
	if hasStaker {
		return fmt.Errorf("Stacker already enrolled")
	}

	// get stake and check if there's enough balance

	var requieredStake *big.Int
	err = s.contract.Call(&requieredStake, "REQUIRED_STAKE")
	if err != nil {
		return err
	}

	balance, err := stk.Client.BalanceInfo()
	if err != nil {
		return err
	}

	if requieredStake.Cmp(balance) > 0 {
		return fmt.Errorf("Not enough balance to stake (requiered=%v avaliable=%v)", requieredStake.String(), balance.String())
	}

	// enroll it
	_, _, err = s.contract.SendTransactionSyncWithClient(
		stk.Client, requieredStake, 0,
		"enroll", stk.OnionAt(HashOnionLayers))

	// write into db

	if err = s.storage.SetStaker(stk.Address, &storage.StakerEntry{
		HashIndex: HashOnionLayers,
	}); err != nil {
		return err
	}

	return err
}

func (s *Stakers) Touch(staker common.Address) error {
	log.Info("Touch ", staker.Hex())
	stk := s.stks[staker]
	_, _, err := s.contract.SendTransactionSyncWithClient(stk.Client, nil, 0,
		"touch")
	return err
}

func (s *Stakers) Slash(staker common.Address, toslash common.Address) error {
	log.Info("Slashing ", toslash.Hex())
	stk := s.stks[staker]
	_, _, err := s.contract.SendTransactionSyncWithClient(stk.Client, nil, 0,
		"slash", toslash)
	return err
}

func (s *Stakers) QueryUnenroll(staker common.Address) error {
	log.Info("QueryUnenroll ", staker.Hex())
	stk := s.stks[staker]
	_, _, err := s.contract.SendTransactionSyncWithClient(stk.Client, nil, 0,
		"touch")
	return err
}

func (s *Stakers) Unenroll(staker common.Address) error {
	log.Info("Unenroll ", staker.Hex())
	stk := s.stks[staker]
	_, _, err := s.contract.SendTransactionSyncWithClient(stk.Client, nil, 0,
		"unenroll")
	return err
}

func (s *Stakers) IsTrueTeller(staker common.Address) (bool, error) {
	stakerEntry, err := s.storage.Staker(staker)
	if err != nil {
		return false, err
	}

	stk := s.stks[staker]

	stakerEntry.HashIndex--
	hash := stk.OnionAt(stakerEntry.HashIndex)
	IsLier := hash[31]&0xf != 0

	return IsLier, nil
}

func (s *Stakers) Resolve(staker common.Address) error {
	log.Info("ResolveChallange ", staker.Hex())
	stk := s.stks[staker]

	stakerEntry, err := s.storage.Staker(staker)
	if err != nil {
		return err
	}
	stakerEntry.HashIndex--
	hash := stk.OnionAt(stakerEntry.HashIndex)
	_, _, err = s.contract.SendTransactionSyncWithClient(stk.Client, nil, 0,
		"resolveChallenge", hash)
	if err != nil {
		return err
	}

	err = s.storage.SetStaker(staker, stakerEntry)
	return err
}
