package proxy2contracts

import "math/big"

func right(x *big.Int, n uint) *big.Int {
	return new(big.Int).Rsh(x, n)
}

func left(x *big.Int, n uint) *big.Int {
	return new(big.Int).Rsh(x, n)
}

// func getSkill(encodedSkills *big.Int, uint8 skillIdx) *big.Int {
// 	return (encodedSkills.Rsh((uint256(skillIdx) * 20)).And(1048575; /// 1048575 = 2**20 - 1
// }

// func getBirthDay(encodedSkills *big.Int) *big.Int {
// 	return (encodedSkills.Rsh(100).And(65535;
// }

// func getPlayerIdFromSkills(encodedSkills *big.Int) *big.Int {
// 	return (getIsSpecial(encodedSkills)) ? encodedSkills : getInternalPlayerId(encodedSkills);
// }

// func getInternalPlayerId(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills.Rsh(129.And(8796093022207); /// 2**43 - 1 = 8796093022207
// }

// func getPotential(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills.Rsh(116.And(15);
// }

func and(x *big.Int, n int64) *big.Int {
	return new(big.Int).And(x, big.NewInt(n))
}

func getForwardnessNat(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 120), 7)
}

func getLeftishnessNat(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 123), 7)
}

func getAggressivenessNat(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 126), 7)
}

// func getAlignedEndOfFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return (encodedSkills.Rsh(172.And(1) == 1;
// }

// func getRedCardLastGame(encodedSkills *big.Int) public pure returns (bool) {
// 	return (encodedSkills.Rsh(173.And(1) == 1;
// }

// func getGamesNonStopping(encodedSkills *big.Int) public pure returns (uint8) {
// 	return uint8(encodedSkills.Rsh(174.And(7);
// }

// func getInjuryWeeksLeft(encodedSkills *big.Int) public pure returns (uint8) {
// 	return uint8(encodedSkills.Rsh(177.And(7);
// }

// func getSubstitutedFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return (encodedSkills.Rsh(180.And(1) == 1;
// }

// func getSumOfSkills(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills.Rsh(181.And(524287); /// 2**19-1
// }

// func getIsSpecial(encodedSkills *big.Int) public pure returns (bool) {
// 	return big.NewInt(encodedSkills.Rsh(204.And(1) == 1;
// }

// func getGeneration(encodedSkills *big.Int) *big.Int {
// 	return (encodedSkills.Rsh(205).And(255;
// }

// func getOutOfGameFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return big.NewInt(encodedSkills.Rsh(213.And(1) == 1;
// }

// func getYellowCardFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return big.NewInt(encodedSkills.Rsh(214.And(1) == 1;
// }
