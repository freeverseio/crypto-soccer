package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
)

func DecodeSkills(assets *assets.Assets, encodedSkills *big.Int) (
	defence *big.Int,
	speed *big.Int,
	pass *big.Int,
	shoot *big.Int,
	endurance *big.Int,
	potencial *big.Int,
	err error,
) {
	opts := &bind.CallOpts{}
	if defence, err = assets.GetDefence(opts, encodedSkills); err != nil {
		return defence, speed, pass, shoot, endurance, potencial, err
	} else if speed, err = assets.GetSpeed(opts, encodedSkills); err != nil {
		return defence, speed, pass, shoot, endurance, potencial, err
	} else if pass, err = assets.GetPass(opts, encodedSkills); err != nil {
		return defence, speed, pass, shoot, endurance, potencial, err
	} else if shoot, err = assets.GetShoot(opts, encodedSkills); err != nil {
		return defence, speed, pass, shoot, endurance, potencial, err
	} else if endurance, err = assets.GetEndurance(opts, encodedSkills); err != nil {
		return defence, speed, pass, shoot, endurance, potencial, err
	} else if potencial, err = assets.GetPotential(opts, encodedSkills); err != nil {
		return defence, speed, pass, shoot, endurance, potencial, err
	}
	return defence, speed, pass, shoot, endurance, potencial, err
}
