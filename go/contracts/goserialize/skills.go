package native

import "math/big"

// func getSkill(encodedSkills *big.Int, uint8 skillIdx) *big.Int {
// 	return (encodedSkills >> (uint256(skillIdx) * 20)) & 1048575; /// 1048575 = 2**20 - 1
// }

// func getBirthDay(encodedSkills *big.Int) *big.Int {
// 	return (encodedSkills >> 100) & 65535;
// }

// func getPlayerIdFromSkills(encodedSkills *big.Int) *big.Int {
// 	return (getIsSpecial(encodedSkills)) ? encodedSkills : getInternalPlayerId(encodedSkills);
// }

// func getInternalPlayerId(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills >> 129 & 8796093022207); /// 2**43 - 1 = 8796093022207
// }

// func getPotential(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills >> 116 & 15);
// }

func getForwardnessNat(encodedSkills *big.Int) *big.Int {
	return big.NewInt(encodedSkills >> 120 & 7)
}

// func getLeftishness(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills >> 123 & 7);
// }

// func getAggressiveness(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills >> 126 & 7);
// }

// func getAlignedEndOfFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return (encodedSkills >> 172 & 1) == 1;
// }

// func getRedCardLastGame(encodedSkills *big.Int) public pure returns (bool) {
// 	return (encodedSkills >> 173 & 1) == 1;
// }

// func getGamesNonStopping(encodedSkills *big.Int) public pure returns (uint8) {
// 	return uint8(encodedSkills >> 174 & 7);
// }

// func getInjuryWeeksLeft(encodedSkills *big.Int) public pure returns (uint8) {
// 	return uint8(encodedSkills >> 177 & 7);
// }

// func getSubstitutedFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return (encodedSkills >> 180 & 1) == 1;
// }

// func getSumOfSkills(encodedSkills *big.Int) *big.Int {
// 	return big.NewInt(encodedSkills >> 181 & 524287); /// 2**19-1
// }

// func getIsSpecial(encodedSkills *big.Int) public pure returns (bool) {
// 	return big.NewInt(encodedSkills >> 204 & 1) == 1;
// }

// func getGeneration(encodedSkills *big.Int) *big.Int {
// 	return (encodedSkills >> 205) & 255;
// }

// func getOutOfGameFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return big.NewInt(encodedSkills >> 213 & 1) == 1;
// }

// func getYellowCardFirstHalf(encodedSkills *big.Int) public pure returns (bool) {
// 	return big.NewInt(encodedSkills >> 214 & 1) == 1;
// }
