using System;
using System.Numerics;
public class Serialization {  

    const int IN_TRANSIT_SHIRTNUM = 26;
    const int MASK_12b = 4095;
    const int MASK_19b = 524287;
    const int MASK_20b = 1048575;
    const int MASK_24b = 16777215;
    const int MASK_28b = 268435455;
    const ulong MASK_35b = 34359738367;
    const ulong MASK_43b = 8796093022207;

    private uint rightShiftAndMask(BigInteger encoded, int bitsToDisplace, int mask) { return (uint) ((encoded >> bitsToDisplace) & mask); }

    private ulong rightShiftAndMask64b(BigInteger encoded, int bitsToDisplace, ulong mask) { return (ulong) ((encoded >> bitsToDisplace) & mask); }

    // STATE
    public ulong getCurrentTeamId(BigInteger state) { return  rightShiftAndMask64b(state, 0, MASK_43b); }
    public uint getCurrentShirtNum(BigInteger state) { return  rightShiftAndMask(state, 43, 31); }

    public ulong getPrevPlayerTeamId(BigInteger state){ return  rightShiftAndMask64b(state, 48, MASK_43b); }

    public ulong getLastSaleBlock(BigInteger state){ return  rightShiftAndMask64b(state, 91, MASK_35b); }

    public bool getIsInTransit(BigInteger state) { return  getCurrentShirtNum(state) == IN_TRANSIT_SHIRTNUM; }

    // SKILLS
    public uint getSkill(BigInteger encodedSkills, int skillIdx) { return  rightShiftAndMask(encodedSkills, skillIdx * 20, MASK_20b); } 

    public uint getBirthDay(BigInteger encodedSkills) { return  rightShiftAndMask(encodedSkills, 100, 65535); }

    public bool getIsSpecial(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 204, 1) == 1; }

    public BigInteger getPlayerIdFromSkills(BigInteger encodedSkills) {
        if (getIsSpecial(encodedSkills)) {    
            return encodedSkills;
        }
        return getInternalPlayerId(encodedSkills); 
    }

    public BigInteger getInternalPlayerId(BigInteger encodedSkills) { return ((encodedSkills >> 129) & MASK_43b); }

    public uint getPotential(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 116, 15); }

    public uint getForwardness(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 120, 7); }

    public uint getLeftishness(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 123, 7); }

    public uint getAggressiveness(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 126, 7); }

    public bool getAlignedEndOfFirstHalf(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 172, 1) == 1; }

    public bool getRedCardLastGame(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 173, 1) == 1; }

    public uint getGamesNonStopping(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 174, 7); }

    public uint getInjuryWeeksLeft(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 177, 7); }

    public bool getSubstitutedFirstHalf(BigInteger encodedSkills) {	return rightShiftAndMask(encodedSkills, 180, 1) == 1; }

    public uint getSumOfSkills(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 181, MASK_19b); }

    public uint getGeneration(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 205, 255); }

    public bool getOutOfGameFirstHalf(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 213, 1) == 1; }

    public bool getYellowCardFirstHalf(BigInteger encodedSkills) { return rightShiftAndMask(encodedSkills, 214, 1) == 1; }


    // // MATCH LOG public uints:

    public uint getAssister(BigInteger log, int pos) { return rightShiftAndMask(log, 4 + 4 * pos, 15); }

    public uint getShooter(BigInteger log, int pos) { return rightShiftAndMask(log, 52 + 4 * pos, 15); }

    public uint getForwardPos(BigInteger log, int pos) { return rightShiftAndMask(log, 100 + 2 * pos, 3); }

    public bool getPenalty(BigInteger log, int pos) { return rightShiftAndMask(log, 124+pos, 1) == 1; }

    public bool getIsHomeStadium(BigInteger log) { return rightShiftAndMask(log, 248, 1) == 1; }

    /// recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
    public uint getHalfTimeSubs(BigInteger log, int pos) { return rightShiftAndMask(log, 179 + 5 * pos, 31); }

    public uint getNGKAndDefs(BigInteger log, bool is2ndHalf) {
        int offset = 194 + 4 * (is2ndHalf ? 1 : 0);
        return rightShiftAndMask(log, offset, 15); 
    }

    public uint getNTot(BigInteger log, bool is2ndHalf) {
        int offset = 202 + (is2ndHalf ? 4 : 0);
        return rightShiftAndMask(log, offset, 15); 
    }

    public uint getWinner(BigInteger log) { return rightShiftAndMask(log, 210, 3); }

    public uint getTeamSumSkills(BigInteger log) { return rightShiftAndMask(log, 212, MASK_24b); } 


    public uint getTrainingPoints(BigInteger log) { return rightShiftAndMask(log, 236, MASK_12b); } 

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
    public uint getChangesAtHalfTime(BigInteger log) {
        return rightShiftAndMask(log, 249, 3); 
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

    // TeamId and PlayerId
    public uint getTimezone(BigInteger encodedId) { return rightShiftAndMask(encodedId, 38, 31); }
    public uint getCountryIdxInTZ(BigInteger encodedId) { return rightShiftAndMask(encodedId, 28, 1023); }
    public uint getValInCountry(BigInteger encodedId) { return rightShiftAndMask(encodedId, 0, MASK_28b); }

}  