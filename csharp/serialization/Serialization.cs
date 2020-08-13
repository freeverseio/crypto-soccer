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
    const uint ERR_TRAINING_SPLAYER = 2;
    const uint ERR_TRAINING_SINGLESKILL = 3;
    const uint ERR_TRAINING_SUMSKILLS = 4;
    const uint ERR_TRAINING_PREVMATCH = 5;
    const uint ERR_TRAINING_STAMINA = 6;
    const uint PLAYERS_PER_TEAM_MAX = 25;
    const uint ROUNDS_PER_HALF = 12;
    const uint NO_SUBST = 11;
    const uint NO_LINEUP = 25;
    const uint MAX_PERCENT = 60;
    const uint ROUNDS_PER_MATCH = 12;


    private uint rightShiftAndMask(BigInteger encoded, int bitsToDisplace, int mask) { return (uint) ((encoded >> bitsToDisplace) & mask); }

    private ulong rightShiftAndMask64b(BigInteger encoded, int bitsToDisplace, ulong mask) { return (ulong) ((encoded >> bitsToDisplace) & mask); }

    private BigInteger OrWithLeftShift(BigInteger original, uint val, int bitsToDisplace) { 
        return original | ((new BigInteger(val)) << bitsToDisplace);
    }

    // PLAYER STATE => contains info about current team, current shirt number, isInTransit...
    public ulong getCurrentTeamId(BigInteger state) { return  rightShiftAndMask64b(state, 0, MASK_43b); }
    public uint getCurrentShirtNum(BigInteger state) { return  rightShiftAndMask(state, 43, 31); }

    public ulong getPrevPlayerTeamId(BigInteger state){ return  rightShiftAndMask64b(state, 48, MASK_43b); }

    public ulong getLastSaleBlock(BigInteger state){ return  rightShiftAndMask64b(state, 91, MASK_35b); }

    public bool getIsInTransit(BigInteger state) { return  getCurrentShirtNum(state) == IN_TRANSIT_SHIRTNUM; }


    // PLAYER SKILLS => contains info about (shoot, pass, defence...), birthday, leftishness, etc.
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


    // MATCH LOG => info about stuff that happened in a match

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

    // TACTICS => lineup, extraAttack, and substituions
    public (BigInteger encoded, string err) encodeTactics(
        uint[] substitutions, 
        uint[] subsRounds, 
        uint[] lineup, 
        bool[] extraAttack, 
        uint tacticsId
    ) 
    {
        // Test on inputs:
        if (substitutions.Length != 3) { return (new BigInteger(0), "length of substitutions must be 3"); }
        if (subsRounds.Length != 3) { return (new BigInteger(0), "length of subsRounds must be 3"); }
        if (lineup.Length != 14) { return (new BigInteger(0), "length of lineup must be 14"); }
        if (extraAttack.Length != 10) { return (new BigInteger(0), "length of extraAttack must be 10"); }
        if (!(tacticsId < 64)) { return (new BigInteger(0), "tacticsId must be less than 64"); }
        for (int p = 0; p < substitutions.Length; p++) {
            if (!(substitutions[p] < 12)) { return (new BigInteger(0), "substitutions entries must be lesss than 12"); }
            if (!(subsRounds[p] < ROUNDS_PER_HALF)) { return (new BigInteger(0), "subsRounds entries must be lesss than ROUNDS_PER_HALF"); }
        }        
        for (int p = 0; p < lineup.Length; p++) {
            if (!(lineup[p] <= PLAYERS_PER_TEAM_MAX)) { return (new BigInteger(0), "lineup entries must be lesss than PLAYERS_PER_TEAM_MAX"); }
        }
        // Start encoding:
        BigInteger encoded = new BigInteger(tacticsId);
        for (int p = 0; p < 10; p++) {
            encoded = OrWithLeftShift(encoded, (uint) (extraAttack[p] ? 1 : 0), 6 + p);
        }          
        for (int p = 0; p < 11; p++) {
            encoded = OrWithLeftShift(encoded, lineup[p], 16 + 5 * p);
        }          
        for (int p = 0; p < 3; p++) {
            /// requirement: if there is no subst at "i", lineup[i + 11] = 25 + p (so that all lineups are different, and sortable)
            if (substitutions[p] == NO_SUBST) {
                if (!(lineup[p + 11] == NO_LINEUP)) return (new BigInteger(0), "incorrect lineup entry for no substituted player");
            }
            encoded = OrWithLeftShift(encoded, lineup[p+11], 16 + 5 * (p + 11));
            encoded = OrWithLeftShift(encoded, substitutions[p], 86 + 4 * p);
            encoded = OrWithLeftShift(encoded, subsRounds[p], 98 + 4 * p);
        }          
        return (encoded, "");
    }


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

    // TEAMID and PLAYERID => info about (timezone, country idx in that timezone, idx in that country)
    // - Teams always remain in the same (timezone, country), players
    // - For players, (timezone, country) refer to where were they originally created. 
    //  - To query about the current (timezone, country) for a player => use playerState to find currentTeamId
    public uint getTimezone(BigInteger encodedId) { return rightShiftAndMask(encodedId, 38, 31); }
    public uint getCountryIdxInTZ(BigInteger encodedId) { return rightShiftAndMask(encodedId, 28, 1023); }
    public uint getValInCountry(BigInteger encodedId) { return rightShiftAndMask(encodedId, 0, MASK_28b); }


    // TRAINING POINTS ASSIGNMENT => encode and decode functions
    public (uint[] TPperSkill, uint specialPlayer, uint TP, uint err) decodeTP(BigInteger encoded) {
        uint[] TPperSkill = new uint[25];
        uint specialPlayer = rightShiftAndMask(encoded, 234, 31);
        uint err = 0;
        uint TP = rightShiftAndMask(encoded, 225, 511);

        if (specialPlayer > PLAYERS_PER_TEAM_MAX) return (TPperSkill, specialPlayer, TP, ERR_TRAINING_SPLAYER); // specialPlayer value too large

        uint TPtemp = TP;
        uint maxRHS = (TPtemp < 4) ? 100 * TPtemp : MAX_PERCENT * TPtemp;
        for (int bucket = 0; bucket < 5; bucket++) {
            if (bucket == 4) {
                TPtemp = (TPtemp * 11)/10;
                maxRHS = (TPtemp < 4) ? 100 * TPtemp : MAX_PERCENT * TPtemp;
            }
            uint sum = 0;
            for (int sk = 5 * bucket; sk < 5 * (bucket+1); sk++) {
                uint skill = rightShiftAndMask(encoded, 9 * sk, 511);
                if (100*skill > maxRHS) return (TPperSkill, specialPlayer, TP, ERR_TRAINING_SINGLESKILL); // one of the assigned TPs is too large or too small
                TPperSkill[sk] = skill;
                sum += skill;
            }
            if (sum > TPtemp) return (TPperSkill, specialPlayer, TP, ERR_TRAINING_SUMSKILLS); // sum of Traning Points is too large"
        }
        return (TPperSkill, specialPlayer, TP, err);
    }


    public (BigInteger encoded, string err) encodeTP(uint TP, uint[] TPperSkill, uint specialPlayer) {
        // Test on inputs:
        if (!(TP < 65536)) { return (new BigInteger(0), "TP value too large"); }
        if (!(specialPlayer <= PLAYERS_PER_TEAM_MAX)) { return (new BigInteger(0), "specialPlayer value too large"); }
        if (TPperSkill.Length != 25) { return (new BigInteger(0), "length of TPperSkill must be 25"); }
        for (int p = 0; p < TPperSkill.Length; p++) {
            if (!(TPperSkill[p] < 65536)) { return (new BigInteger(0), "TPperSkill entries too large"); }
        }
        // Start encoding:
        BigInteger encoded = 0;
        encoded = OrWithLeftShift(encoded, TP, 225);
        encoded = OrWithLeftShift(encoded, specialPlayer, 234);
        uint maxRHS = (TP < 4) ? 100 * TP : MAX_PERCENT * TP;
        int lastBucket = (specialPlayer == PLAYERS_PER_TEAM_MAX ? 4 : 5);
        for (int bucket = 0; bucket < lastBucket; bucket++) {
            if (bucket == 4) {
                TP = (TP * 11)/10;
                maxRHS = (TP < 4) ? 100 * TP : MAX_PERCENT * TP;
            }
            uint sum = 0;
            for (int sk = 5 * bucket; sk < 5 * (bucket+1); sk++) {
                if (!(100*TPperSkill[sk] <= maxRHS)) { return (new BigInteger(0), "one of the assigned TPs is too large"); }
                sum += TPperSkill[sk];
                encoded = OrWithLeftShift(encoded, TPperSkill[sk], 9 * sk);
            }
            if (!(sum <= TP)) { return (new BigInteger(0), "sum of Traning Points is too large"); }
        }
        return (encoded, "");
    }


    // MATCH EVENTS => Creates the events that happen in ONE HALF of a match
    // From all inputs, only the last one is computed by the blockchain. The others are ready as soon as the user actions are submitted.
    // So, the frontend can prepare everything, and only wait for the backend to provide the last input.
    // - INPUTS
    //      - seed used in that halfMatch
    //      - teamIds for home/visitor teams 
    //      - encodedTactics for home/visitor teams for that halfMatch
    //      - playerIds[25] for home/visitor teams for that halfMatch (note: 2nd half could differ from 1st half if trading happened)
    //      - matchLogsAndEvents[2 + 5 * ROUNDS_PER_MATCH] => the only piece computed by the blockchain
    //          - first 2 entries are matchLog[homeTeam], matchLog[visitorTeam]
    //          - next, we have packs of 5 numbers, one for each round of the halfMatch
    public (MatchEvent[] events, string err) processMatchEvents(
        BigInteger seed, 
        BigInteger[] teamIds, 
        BigInteger[] tactics, 
        BigInteger[][] playerIds, 
        BigInteger[] matchLogsAndEvents
    ) 
    {
        MatchEvent[] events = new MatchEvent[0];
        // Test on inputs:
        if (teamIds.Length != 2) { return (events, "length of teamIds must be 2"); }
        if (tactics.Length != 2) { return (events, "length of tactics must be 2"); }
        if (matchLogsAndEvents.Length != (2 + 5 * ROUNDS_PER_MATCH)) { return (events, "length of matchLogsAndEvents must be 2 + 5 * ROUNDS_PER_MATCH"); }
        if (playerIds.Length != 2) { return (events, "length of playerIds must be 2"); }
        if (playerIds[0].Length != PLAYERS_PER_TEAM_MAX) { return (events, "length of playerIds[0] must be PLAYERS_PER_TEAM_MAX"); }
        if (playerIds[0].Length != PLAYERS_PER_TEAM_MAX) { return (events, "length of playerIds[1] must be PLAYERS_PER_TEAM_MAX"); }
        return (events, "");
    }

    public struct MatchEvent {
        public int Minute { get; }
        public int Type { get; }
        public int Team { get; }
        public bool ManagesToShoot { get; }
        public bool IsGoal { get; }
        public int PrimaryPlayer { get; }
        public int SecondaryPlayer { get; }
        public BigInteger PrimaryPlayerID { get; }
        public BigInteger SecondaryPlayerID { get; }
    }
}