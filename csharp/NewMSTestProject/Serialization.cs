using System.Numerics;
public class Serialization {  
    public double Add(double num1, double num2) {  
        return num1 + num2;  
    }  
    public BigInteger AddBN(BigInteger x, BigInteger y) {  
        return BigInteger.Add(x, y);  
    }

    private uint rightShiftAndMask(BigInteger encoded, int bitsToDisplace, int mask) {
        return (uint) ((encoded >> bitsToDisplace) & mask);
    }

    public uint getCurrentShirtNum(BigInteger playerState) { return  rightShiftAndMask(playerState, 43, 31); }

    public uint getSkill(BigInteger playerState, int skillIdx) { return  rightShiftAndMask(playerState, skillIdx * 20, 1048575); } // 1048575 = 2**20 - 1

    public uint getBirthDay(BigInteger playerState) { return  rightShiftAndMask(playerState, 100, 65535); }


// public uint getPlayerIdFromSkills(BigInteger encodedSkills) {
// 	if getIsSpecial(encodedSkills) {
// 		return encodedSkills
// 	}
// 	return getInternalPlayerId(encodedSkills)
// }

// public uint getInternalPlayerId(BigInteger encodedSkills) {
// 	return rightShiftAndMask(encodedSkills, 129), 8796093022207) // 2**43 - 1 = 8796093022207
// }

    public uint getPotential(BigInteger encodedSkills) {
        return rightShiftAndMask(encodedSkills, 116, 15);
    }

    public uint getForwardness(BigInteger encodedSkills) {
        return rightShiftAndMask(encodedSkills, 120, 7);
    }

    public uint getLeftishness(BigInteger encodedSkills) {
        return rightShiftAndMask(encodedSkills, 123, 7);
    }

    public uint getAggressiveness(BigInteger encodedSkills) {
        return rightShiftAndMask(encodedSkills, 126, 7);
    }

    public bool getAlignedEndOfFirstHalf(BigInteger encodedSkills) {
        return rightShiftAndMask(encodedSkills, 172, 1) == 1;
    }

    public bool getRedCardLastGame(BigInteger encodedSkills) {
        return rightShiftAndMask(encodedSkills, 173, 1) == 1;
    }

    public uint getGamesNonStopping(BigInteger encodedSkills) {
        return rightShiftAndMask(encodedSkills, 174, 7);
    }

public uint getInjuryWeeksLeft(BigInteger encodedSkills) {
	return rightShiftAndMask(encodedSkills, 177, 7);
}

public bool getSubstitutedFirstHalf(BigInteger encodedSkills) {
	return rightShiftAndMask(encodedSkills, 180, 1) == 1;
}

public uint getSumOfSkills(BigInteger encodedSkills) {
	return rightShiftAndMask(encodedSkills, 181, 524287); // 2**19-1
}

public bool getIsSpecial(BigInteger encodedSkills) {
	return rightShiftAndMask(encodedSkills, 204, 1) == 1;
}

public uint getGeneration(BigInteger encodedSkills) {
	return rightShiftAndMask(encodedSkills, 205, 255);
}

public bool getOutOfGameFirstHalf(BigInteger encodedSkills) {
	return rightShiftAndMask(encodedSkills, 213, 1) == 1;
}

public bool getYellowCardFirstHalf(BigInteger encodedSkills) {
	return rightShiftAndMask(encodedSkills, 214, 1) == 1;
}

// // MATCH LOG functions:

// public uint getAssisterGo(log , pos uint8) {
// 	val := rightShiftAndMask(log, (4+4*uint(pos))), 15)
// 	return uint8(val.Uint64())
// }

// public uint getShooterGo(log , pos uint8) {
// 	val := rightShiftAndMask(log, (52+4*uint(pos))), 15)
// 	return uint8(val.Uint64())
// }

// public uint getForwardPosGo(log , pos uint8) {
// 	val := rightShiftAndMask(log, (100+2*uint(pos))), 3)
// 	return uint8(val.Uint64())
// }

// public uint getPenaltyGo(log , pos uint8) {
// 	return equals(rightShiftAndMask(log, 124+uint(pos)), 1), 1)
// }

// public uint getIsHomeStadiumGo(log) {
// 	return equals(rightShiftAndMask(log, 248), 1), 1)
// }

// /// recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
// public uint getHalfTimeSubsGo(log , pos uint8) {
// 	val := rightShiftAndMask(log, (179+5*uint(pos))), 31)
// 	return uint8(val.Uint64())
// }

// public uint getNGKAndDefsGo(log , is2ndHalf bool) {
// 	offset := uint(194)
// 	if is2ndHalf {
// 		offset += 4
// 	}
// 	val := rightShiftAndMask(log, offset), 15)
// 	return uint8(val.Uint64())
// }

// public uint getNTotGo(log , is2ndHalf bool) {
// 	offset := uint(202)
// 	if is2ndHalf {
// 		offset += 4
// 	}
// 	val := rightShiftAndMask(log, offset), 15)
// 	return uint8(val.Uint64())
// }

// public uint getWinnerGo(log) {
// 	val := rightShiftAndMask(log, 210), 3)
// 	return uint8(val.Uint64())
// }

// public uint getTeamSumSkillsGo(log) {
// 	return rightShiftAndMask(log, 212), 16777215) // 2^24 - 1
// }

// public uint addTrainingPointsGo(log , points) {
// 	return orBN(log, left(points, 236))
// }

// public uint getTrainingPointsGo(log) uint16 {
// 	val := rightShiftAndMask(log, 236), 4095) // 2^12-1
// 	return uint16(val.Uint64())
// }

// public uint getNGoalsGo(log) {
// 	val := and(log, 15)
// 	return uint8(val.Uint64())
// }

// public uint getOutOfGamePlayerGo(log , is2ndHalf bool) {
// 	var offset uint
// 	if is2ndHalf {
// 		offset = 141
// 	} else {
// 		offset = 131
// 	}
// 	return rightShiftAndMask(log, offset), 15)
// }

// public uint getOutOfGameTypeGo(log , is2ndHalf bool) {
// 	var offset uint
// 	if is2ndHalf {
// 		offset = 141
// 	} else {
// 		offset = 131
// 	}
// 	return rightShiftAndMask(log, offset+4), 3)
// }

// public uint getOutOfGameRoundGo(log , is2ndHalf bool) {
// 	var offset uint
// 	if is2ndHalf {
// 		offset = 141
// 	} else {
// 		offset = 131
// 	}
// 	return rightShiftAndMask(log, offset+6), 15)
// }

// public uint getYellowCardGo(log , posInHaf uint8, is2ndHalf bool) {
// 	offset := uint(posInHaf) * 4
// 	if is2ndHalf {
// 		offset += 159
// 	} else {
// 		offset += 151
// 	}
// 	val := rightShiftAndMask(log, offset), 15)
// 	return uint8(val.Uint64())
// }

// public uint getInGameSubsHappenedGo(log , posInHalf uint8, is2ndHalf bool) {
// 	offset := 167 + 2*uint(posInHalf)
// 	if is2ndHalf {
// 		offset += 6
// 	}
// 	val := rightShiftAndMask(log, offset), 3)
// 	return uint8(val.Uint64())
// }

// // TACTICS

// public uint getTacticsIdGo(tactics) {
// 	val := and(tactics, 63)
// 	return uint8(val.Uint64())
// }

// public uint getExtraAttackGo(tactics , p uint8) {
// 	return equals(rightShiftAndMask(tactics, 6+uint(p)), 1), 1)
// }

// public uint getFullExtraAttackGo(tactics) [10]{
// 	var extraAttack [10]bool
// 	for p := uint8(0); p < 10; p++ {
// 		extraAttack[p] = getExtraAttackGo(tactics, p)
// 	}
// 	return extraAttack
// }

// public uint getSubstitutionGo(tactics , p uint8) {
// 	val := rightShiftAndMask(tactics, 86+4*uint(p)), 15)
// 	return uint8(val.Uint64())
// }

// public uint getSubsRoundGo(tactics , p uint8) {
// 	val := rightShiftAndMask(tactics, 98+4*uint(p)), 15)
// 	return uint8(val.Uint64())
// }

// public uint getLinedUpGo(tactics , p uint8) {
// 	val := rightShiftAndMask(tactics, 16+5*uint(p)), 31)
// 	return uint8(val.Uint64())
// }

// public uint getFullLineUpGo(tactics) [14]{
// 	var lineup [14]uint8
// 	for p := uint8(0); p < 14; p++ {
// 		lineup[p] = getLinedUpGo(tactics, p)
// 	}
// 	return lineup
// }

}  