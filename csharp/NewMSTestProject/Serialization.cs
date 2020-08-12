using System.Numerics;
public class Serialization {  
    public double Add(double num1, double num2) {  return num1 + num2;   }  
    public BigInteger AddBN(BigInteger x, BigInteger y) {  return BigInteger.Add(x, y);   }

    private uint rightShiftAndMask(BigInteger encoded, int bitsToDisplace, int mask) { return (uint) ((encoded >> bitsToDisplace) & mask); }

    public uint getCurrentShirtNum(BigInteger playerState) { return  rightShiftAndMask(playerState, 43, 31); }

    public uint getSkill(BigInteger encodedSkills, int skillIdx) { return  rightShiftAndMask(encodedSkills, skillIdx * 20, 1048575); } // 1048575 = 2**20 - 1

    public uint getBirthDay(BigInteger encodedSkills) { return  rightShiftAndMask(encodedSkills, 100, 65535); }

    public bool getIsSpecial(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 204, 1) == 1; }

    public BigInteger getPlayerIdFromSkills(BigInteger encodedSkills) {if (getIsSpecial(encodedSkills)) {    return encodedSkills;}return getInternalPlayerId(encodedSkills); }

    public BigInteger getInternalPlayerId(BigInteger encodedSkills) { return ((encodedSkills >> 129) & 8796093022207); } // 2**43 - 1 = 8796093022207

    public uint getPotential(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 116, 15); }

    public uint getForwardness(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 120, 7); }

    public uint getLeftishness(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 123, 7); }

    public uint getAggressiveness(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 126, 7); }

    public bool getAlignedEndOfFirstHalf(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 172, 1) == 1; }

    public bool getRedCardLastGame(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 173, 1) == 1; }

    public uint getGamesNonStopping(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 174, 7); }

    public uint getInjuryWeeksLeft(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 177, 7); }

    public bool getSubstitutedFirstHalf(BigInteger encodedSkills) {	return rightShiftAndMask(encodedSkills, 180, 1) == 1; }

    public uint getSumOfSkills(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 181, 524287); } // 2**19-1

    public uint getGeneration(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 205, 255); }

    public bool getOutOfGameFirstHalf(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 213, 1) == 1; }

    public bool getYellowCardFirstHalf(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 214, 1) == 1; }

    // // MATCH LOG functions:

    public uint getAssister(BigInteger log, int pos) { return rightShiftAndMask(log, 4 + 4 * pos, 15); }

    public uint getShooter(BigInteger log, int pos) { return rightShiftAndMask(log, 52 + 4 * pos, 15); }

    public uint getForwardPos(BigInteger log, int pos) { return rightShiftAndMask(log, 100 + 2 * pos, 3); }

    public bool getPenalty(BigInteger log, int pos) { return rightShiftAndMask(log, 124+pos, 1) == 1; }

    public bool getIsHomeStadiumGo(BigInteger log) { return rightShiftAndMask(log, 248, 1) == 1; }

    /// recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
    public uint getHalfTimeSubs(BigInteger log, int pos) { return rightShiftAndMask(log, 179 + 5 * pos, 31); }

    public uint getNGKAndDefs(BigInteger log, bool is2ndHalf) {
        int offset = is2ndHalf ? 194 : 194 + 4;
        return rightShiftAndMask(log, offset, 15); 
    }

    public uint getNTot(BigInteger log, bool is2ndHalf) {
        int offset = is2ndHalf ? 202 : 202 + 4;
        return rightShiftAndMask(log, offset, 15); 
    }

    public uint getWinner(BigInteger log) { return rightShiftAndMask(log, 210, 3); }

    public uint getTeamSumSkills(BigInteger log) { return rightShiftAndMask(log, 212, 16777215); } // 2^24 - 1 


    public uint getTrainingPoints(BigInteger log) { return rightShiftAndMask(log, 236, 4095); } // 2^12-1 

    public uint getNGoals(BigInteger log) { return rightShiftAndMask(log, 0, 15); }

    public uint getOutOfGamePlayer(BigInteger log, bool is2ndHalf) {
        int offset = is2ndHalf ? 141 : 131;
        return rightShiftAndMask(log, offset, 15); 
    }

    public uint getOutOfGameType(BigInteger log, bool is2ndHalf) {
        int offset = is2ndHalf ? 141 : 131;
        return rightShiftAndMask(log, offset+4, 3); 
    }

    public uint getOutOfGameRound(BigInteger log, bool is2ndHalf) {
        int offset = is2ndHalf ? 141 : 131;
        return rightShiftAndMask(log, offset+6, 15); 
    }

    public uint getYellowCard(BigInteger log, int posInHalf, bool is2ndHalf) {
        int offset = 4 * posInHalf + (is2ndHalf ? 159 : 151);
        return rightShiftAndMask(log, offset, 15); 
    }

    public uint getInGameSubsHappened(BigInteger log, int posInHalf, bool is2ndHalf) {
        int offset = 167 + 2 * posInHalf + (is2ndHalf ? 6 : 0);
        return rightShiftAndMask(log, offset, 3); 
    }

    // TACTICS

    public uint getTacticsId(BigInteger tactics) { return rightShiftAndMask(tactics, 0, 63); }

    public bool getExtraAttack(BigInteger tactics, int p) { return rightShiftAndMask(tactics, 6+p, 1) == 1; }

    public bool[] getFullExtraAttack(BigInteger tactics) { 
        bool[] extraAttack = new bool[10];
        for (int p = 0; p < 10; p++) {
            extraAttack[p] = getExtraAttack(tactics, p);
        }
        return extraAttack; 
    }

    public uint getSubstitution(BigInteger tactics, int p) { return rightShiftAndMask(tactics, 86 + 4 * p, 15); }

    public uint getSubsRound(BigInteger tactics, int p) { return rightShiftAndMask(tactics, 98 + 4 * p, 15); }

    public uint getLinedUp(BigInteger tactics, int p) { return rightShiftAndMask(tactics, 16 + 5 * p, 31); }

    public uint[] getFullLineUp(BigInteger tactics) {
        uint[] lineup = new uint[14];
        for (int p = 0; p < 14; p++) {
            lineup[p] = getLinedUp(tactics, p);
        }
        return lineup; 
    }
}  