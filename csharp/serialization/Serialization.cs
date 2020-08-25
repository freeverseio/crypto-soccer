using System;
using System.Numerics;
using System.Security.Cryptography;
using System.IO;
using System.Linq;
using System.Collections;
using System.Collections.Generic;
using System.Text;  

public class Serialization {  

    // CONSTANTS
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
    const uint NIL_PLAYERID = 0;
    const int EVNT_NULL = -1;
    const uint EVNT_NOONE = 14;
    const uint EVNT_PENALTY = 100;
    const uint EVNT_ATTACK = 0;
    const uint EVNT_YELLOW = 1;
	const uint EVNT_RED  = 2;
	const uint EVNT_SOFT = 3;
	const uint EVNT_HARD = 4;
	const uint EVNT_SUBST = 5;
    BigInteger NIL_PLAYER_ID = new BigInteger(0);


    // ENCODING OF PLAYER STATE => contains info about current team, current shirt number, isInTransit...
    public ulong getCurrentTeamId(BigInteger state) { return  rightShiftAndMask64b(state, 0, MASK_43b); }
    public uint getCurrentShirtNum(BigInteger state) { return  rightShiftAndMask(state, 43, 31); }
    public ulong getPrevPlayerTeamId(BigInteger state){ return  rightShiftAndMask64b(state, 48, MASK_43b); }
    public ulong getLastSaleBlock(BigInteger state){ return  rightShiftAndMask64b(state, 91, MASK_35b); }
    public bool getIsInTransit(BigInteger state) { return  getCurrentShirtNum(state) == IN_TRANSIT_SHIRTNUM; }


    // ENCODING OF PLAYER SKILLS => contains info about (shoot, pass, defence...), birthday, leftishness, etc.
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


    // ENCODING OF MATCH LOG => info about stuff that happened in a match
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
    public uint[] fullDecodeMatchLog(BigInteger log, bool is2ndHalf) {
        uint[] decodedLog = new uint[15];
        decodedLog[0] = getTeamSumSkills(log);
        decodedLog[1] = getWinner(log);
        decodedLog[2] = getNGoals(log);
        if (is2ndHalf) decodedLog[3] = getTrainingPoints(log);
        
        decodedLog[4] = getOutOfGamePlayer(log, is2ndHalf);
        decodedLog[5] = getOutOfGameType(log, is2ndHalf);
        decodedLog[6] = getOutOfGameRound(log, is2ndHalf);
    
        decodedLog[7] = getYellowCard(log, 0, is2ndHalf);
        decodedLog[8] = getYellowCard(log, 1, is2ndHalf);
        
        decodedLog[9]  = getInGameSubsHappened(log, 0, is2ndHalf);
        decodedLog[10] = getInGameSubsHappened(log, 1, is2ndHalf);
        decodedLog[11] = getInGameSubsHappened(log, 2, is2ndHalf);

        /// getHalfTimeSubs: recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
        if (is2ndHalf) {
            decodedLog[12]  = getHalfTimeSubs(log, 0);
            decodedLog[13]  = getHalfTimeSubs(log, 1);
            decodedLog[14]  = getHalfTimeSubs(log, 2);
        }
        return decodedLog;
    }


    // ENCODING OF TACTICS => lineup, extraAttack, and substitutions
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
    public uint[] getFullSubstitutions(BigInteger tactics) {
        uint[] subs = new uint[3];
        for (int p = 0; p < 3; p++) {
            subs[p] = getSubstitution(tactics, p);
        }
        return subs; 
    }
    public uint[] getFullSubsRounds(BigInteger tactics) {
        uint[] subs = new uint[3];
        for (int p = 0; p < 3; p++) {
            subs[p] = getSubsRound(tactics, p);
        }
        return subs; 
    }

    // ENCODING OF TEAMID and PLAYERID => info about (timezone, country idx in that timezone, idx in that country)
    // - Teams always remain in the same (timezone, country), players
    // - For players, (timezone, country) refer to where were they originally created. 
    //  - To query about the current (timezone, country) for a player => use playerState to find currentTeamId
    public uint getTimezone(BigInteger encodedId) { return rightShiftAndMask(encodedId, 38, 31); }
    public uint getCountryIdxInTZ(BigInteger encodedId) { return rightShiftAndMask(encodedId, 28, 1023); }
    public uint getValInCountry(BigInteger encodedId) { return rightShiftAndMask(encodedId, 0, MASK_28b); }


    // ENCODING OF TRAINING POINTS ASSIGNMENT => encode and decode functions
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


    // ENCODING OF MATCH EVENTS => Creates the events that happen in ONE HALF of a match
    // From all inputs, only the last one is computed by the blockchain. The others are ready as soon as the user actions are submitted.
    // So, the frontend can prepare everything, and only wait for the backend to provide the last input.
    // - INPUTS
    //      - verseSeed used in that halfMatch
    //      - teamIds for home/visitor teams 
    //      - encodedTactics for home/visitor teams for that halfMatch
    //      - playerIds[25] for home/visitor teams for that halfMatch (note: 2nd half could differ from 1st half if trading happened)
    //      - matchLogsAndEvents[2 + 5 * ROUNDS_PER_MATCH] => the only piece computed by the blockchain
    //          - first 2 entries are matchLog[homeTeam], matchLog[visitorTeam]
    //          - next, we have packs of 5 numbers, one for each round of the halfMatch
    public (MatchEvent[] events, string err) processMatchEvents(
        bool is2ndHalf,
        string verseSeed, 
        BigInteger[] teamIds, 
        BigInteger[] tactics, 
        BigInteger[][] playerIds, 
        BigInteger[] matchLogsAndEvents
    ) 
    {
        MatchEvent[] nilEvents = new MatchEvent[0];
        // Test on inputs:
        if (teamIds.Length != 2) { return (nilEvents, "length of teamIds must be 2"); }
        if (tactics.Length != 2) { return (nilEvents, "length of tactics must be 2"); }
        if (matchLogsAndEvents.Length != (2 + 5 * ROUNDS_PER_MATCH)) { return (nilEvents, "length of matchLogsAndEvents must be 2 + 5 * ROUNDS_PER_MATCH"); }
        if (playerIds.Length != 2) { return (nilEvents, "length of playerIds must be 2"); }
        if (playerIds[0].Length != PLAYERS_PER_TEAM_MAX) { return (nilEvents, "length of playerIds[0] must be PLAYERS_PER_TEAM_MAX"); }
        if (playerIds[0].Length != PLAYERS_PER_TEAM_MAX) { return (nilEvents, "length of playerIds[1] must be PLAYERS_PER_TEAM_MAX"); }

        // Deserialize inputs:
        uint[][] lineup = new uint[2][];
        uint[][] purgedLineup = new uint[2][];
        uint[][] substitutions = new uint[2][];
        uint[][] subsRounds = new uint[2][];
        uint[][] decodedLogs = new uint[2][];
        
        for (int team = 0; team < 2; team++) {
            lineup[team] = getFullLineUp(tactics[team]);
            substitutions[team] = getFullSubstitutions(tactics[team]);
            subsRounds[team] = getFullSubsRounds(tactics[team]);
            purgedLineup[team] = RemoveFreeShirtsFromLineUp(lineup[team], playerIds[team]);
            decodedLogs[team] = fullDecodeMatchLog(matchLogsAndEvents[team], is2ndHalf);
        }

        // Actual computation:
        (MatchEvent[] events, string err) = GenerateEvents(
            verseSeed,
            teamIds[0],
            teamIds[1],
            decodedLogs[0],
            decodedLogs[1],
            matchLogsAndEvents,
            purgedLineup[0],
            purgedLineup[1],
            substitutions[0],
            substitutions[1],
            subsRounds[0],
            subsRounds[1],
            is2ndHalf         
        );
        if (!(err == "")) { return (events, err); }

        (events, err) = populateWithPlayerID(events, playerIds[0], playerIds[1]);
        if (!(err == "")) { return (events, err); }

        return (events, "");
    }
    public struct MatchEvent {
        public MatchEvent(uint minute, uint type, uint team, bool managesToShoot, bool isGoal, int primaryPlayer, int secondaryPlayer, BigInteger primaryPlayerID, BigInteger secondaryPlayerID) {
            Minute = minute;
            Type = type;
            Team = team;
            ManagesToShoot = managesToShoot;
            IsGoal = isGoal;
            PrimaryPlayer = primaryPlayer;
            SecondaryPlayer = secondaryPlayer;
            PrimaryPlayerID = primaryPlayerID;
            SecondaryPlayerID = secondaryPlayerID;
        }
        public uint Minute;
        public uint Type;
        public uint Team;
        public bool ManagesToShoot;
        public bool IsGoal;
        public int PrimaryPlayer;
        public int SecondaryPlayer;
        public BigInteger PrimaryPlayerID;
        public BigInteger SecondaryPlayerID;
    }

    // FROM THIS POINT: Private functions required by ProcessMatchEvents
    private (MatchEvent[] events, string err) GenerateEvents(
            string verseSeed,
            BigInteger teamId0,
            BigInteger teamId1,
            uint[] matchLog0,
            uint[] matchLog1,
            BigInteger[] blockchainEvents,
            uint[] lineup0,
            uint[] lineup1,
            uint[] substitutions0,
            uint[] substitutions1,
            uint[] subsRounds0,
            uint[] subsRounds1,
            bool is2ndHalf
    ) 
    {
        MatchEvent[] dummyEvents = new MatchEvent[0];
        // toni
        // Minimal input checks
        if ((blockchainEvents.Length-2) % 5 != 0) { return (dummyEvents, "the length of blockchainEvents should be 2 + a multiple of"); }
        if (!isOutOfGameDataOK(matchLog0)) { return (dummyEvents, "incorrect matchlog entry"); }
        if (!isOutOfGameDataOK(matchLog1)) { return (dummyEvents, "incorrect matchlog entry"); }

        ulong seed0 = intHash(verseSeed + "_0_" + teamId0.ToString() + "_" + teamId1.ToString());
        ulong seed1 = intHash(verseSeed + "_1_" + teamId0.ToString() + "_" + teamId1.ToString());
        ulong seed2 = intHash(verseSeed + "_2_" + teamId0.ToString() + "_" + teamId1.ToString());

        // There are mainly 3 types of events to reports, which are in different parts of the inputs:
        // - per-round (always 12 per half)
        // - cards & injuries
        // - substitutions
        (MatchEvent[] events, uint[] rounds2mins) = addEventsInRound(seed0, blockchainEvents, lineup0, lineup1);
        string err;
    	(events, err) = addCardsAndInjuries(0, events, seed1, matchLog0, rounds2mins, lineup0);
        if (err != "") { return (events, err); }
    	(events, err) = addCardsAndInjuries(1, events, seed2, matchLog1, rounds2mins, lineup1);
        if (err != "") { return (events, err); }

        events = addSubstitutions(0, events, matchLog0, rounds2mins, lineup0, substitutions0, subsRounds0);
        events = addSubstitutions(1, events, matchLog1, rounds2mins, lineup1, substitutions1, subsRounds1);

        if (is2ndHalf) {
            for (uint e = 0; e < events.Length; e++) {
                events[e].Minute += 145;
            }
        }
        return (events, "");
    }

    // This function makes sure that all players in lineUp exist in the Universe.
    // To avoid, for example, selling a player after setting the lineUp.
    // It sets to NOONE all lineUp entries pointing to playerIds that are larger than 2.
    // (recall that playerID = 0 if not set, or = 1 if sold)
    private uint[] RemoveFreeShirtsFromLineUp(uint[] lineUp, BigInteger[] playerIDs) {
        BigInteger MIN_PLAYERID = new BigInteger(2);
        for (int l = 0; l < lineUp.Length; l++) {
            if (lineUp[l] < NO_LINEUP) {
                BigInteger playerId = playerIDs[lineUp[l]];
                if ((playerId == NIL_PLAYERID) || (playerId <= MIN_PLAYERID)) {
                    lineUp[l] = NO_LINEUP;
                }
            }
        }
        return lineUp;
    }

    // INPUTS.MATCHLOG:
    //  	0 teamSumSkills
    //  	1 winner: 0 = home, 1 = away, 2 = draw
    //  	2 nGoals
    //  	3 trainingPoints
    //  	4 uint8 memory outOfGamePlayer
    //  	5 uint8 memory typesOutOfGames,
    //     		 injuryHard:  1
    //     		 injuryLow:   2
    //     		 redCard:     3
    //  	6 uint8 memory outOfGameRounds
    //  	7,8 uint8[2] memory yellowCards
    //  	9,10,11 uint8[3] memory ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen
    //  	12,13,14 uint8[3] memory halfTimeSubstitutions: 0...10 the player in the starting 11 that was changed during half time
    // OUTPUTS:
    //		an array of variable size, where each entry is a Matchevent struct
    //			0: minute
    // 			1: eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
    // 				see: getInGameSubsHappened
    // 			2: team: 0, 1
    // 			3: managesToShoot
    // 			4: isGoal
    // 			5: primary player (0...11):
    // 				(type == 0, 1) && managesToShoot 	: shooter
    // 				(type == 0, 1) && !managesToShoot 	: tackler
    // 				(type == 2) 						: yellowCarded
    // 				(type == 3) 						: redCarded
    // 				(type == 4,5) 						: injured
    // 				(type == 6) 						: getting out of field
    // 			6: secondary player (0...11):
    // 				(type == 0, 1) && managesToShoot 	: assister
    // 				(type == 0, 1) && !managesToShoot 	: null
    // 				(type == 2) 						: null
    // 				(type == 3) 						: null
    // 				(type == 4,5) 						: null
    // 				(type == 6) 						: getting inside field
    private bool isOutOfGameDataOK(uint[] matchLog) {
        uint outOfGamePlayer = matchLog[4];
        bool thereWasAnOutOfGame = (outOfGamePlayer < EVNT_NOONE);
        if (thereWasAnOutOfGame && (matchLog[5] > 3 || matchLog[5] == 0)) {
            return false;
        }
        return true;
    }

   public ulong intHash(string x) {
        SHA256 sha256Hash = SHA256.Create();
        byte[] bytes = sha256Hash.ComputeHash(Encoding.UTF8.GetBytes(x));  
        // retain only the first 8 bytes and convert to uint64
        BigInteger biggy = new BigInteger(bytes, true, true);
        return (ulong) (biggy >> 192);
    }

    public uint GenerateRnd(BigInteger seed, string salt, uint max_val) {
        ulong result = intHash(seed.ToString() + salt);
        return (uint) (result % ((ulong) max_val));
    }

    private (MatchEvent[] events, uint[] rounds2mins) addEventsInRound(BigInteger seed, BigInteger[] blockchainEvents, uint[] lineup0, uint[] lineup1) {
        uint nEvents = (uint) (blockchainEvents.Length - 2) / 5;
        MatchEvent[] events = new MatchEvent[nEvents];
        uint[] rounds2mins = new uint[nEvents];
        uint deltaMinutes = 45 / (nEvents - 1);
        Console.WriteLine(deltaMinutes);
        uint lastMinute = 0;
        uint[][] lineUps = new uint[2][]{lineup0, lineup1};
        for (uint e = 0; e < nEvents; e++) {
            // compute minute
            string salt = "a" + e.ToString();
            uint minute = e * deltaMinutes + GenerateRnd(seed, salt, deltaMinutes);
            if (minute <= lastMinute) {
                minute = lastMinute + 1;
            }
            lastMinute = minute;
            rounds2mins[e] = minute;
            // parse type of event and data
            // note that both "shooter" and "assister" referred to the lineup players (0...10)
            // so they need to be converted to shirtNums by using lineUp
            uint teamThatAttacks = (uint) blockchainEvents[2+5*e];
            bool managesToShoot = blockchainEvents[2+5*e+1] == 1;
            uint shooter = (uint) blockchainEvents[2+5*e+2];
            bool isGoal = blockchainEvents[2+5*e+3] == 1;
            uint assister = (uint) blockchainEvents[2+5*e+4];
            MatchEvent thisEvent = new MatchEvent();
            thisEvent.Minute = minute;
            thisEvent.Type = EVNT_ATTACK;
            thisEvent.Team = teamThatAttacks;
            thisEvent.ManagesToShoot = managesToShoot;
            thisEvent.IsGoal = isGoal;
            if (managesToShoot) {
                // // select the players from the team that attacks:
                thisEvent.PrimaryPlayer = toShirtNum(shooter, lineUps[thisEvent.Team]);
                if (assister == EVNT_PENALTY) {
                    thisEvent.SecondaryPlayer = (int) EVNT_PENALTY;
                } else {
                    thisEvent.SecondaryPlayer = toShirtNum(assister, lineUps[thisEvent.Team]);
                }
            } else {
                salt = "b" + e.ToString();
                // select the player from the team that defends:
                thisEvent.PrimaryPlayer = toShirtNum(1+GenerateRnd(seed, salt, 9), lineUps[1-thisEvent.Team]);
                thisEvent.SecondaryPlayer = EVNT_NULL;
            }
            events[e] = thisEvent;
    	}
        return (events, rounds2mins);
    }

    private int toShirtNum(uint posInLineUp, uint[] lineUp) {
        if (posInLineUp < EVNT_NOONE) {
            return preventNoPlayer(lineUp[posInLineUp]);
        } else {
            return EVNT_NULL;
        }
    }

    private int preventNoPlayer(uint inPlayer) {
        if (inPlayer < 25) {
            return (int) inPlayer;
        } else {
            return EVNT_NULL;
        }
    }

    private (MatchEvent[] events, string err) addCardsAndInjuries(uint team, MatchEvent[] events, BigInteger seed, uint[] matchLog, uint[] rounds2mins, uint[] lineUp)  {
        // matchLog[4,5,6] = outOfGamePlayer, outOfGameType, outOfGameRound
        // note that outofgame is a number from 0 to 13, and that NO OUT OF GAME = 14
        // eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
        if (matchLog[5] > 3) {
            return (new MatchEvent[0], "typeOfEvent larger than 3");
        }
        if (matchLog[6] >= rounds2mins.Length) {
            return (new MatchEvent[0], "outOfGameRound larger than allowed");
        }
        List<MatchEvent> newEvents = events.ToList();


        uint outOfGamePlayer = matchLog[4];
        // convert player in the lineUp to shirtNum before storing it as match event:
        int primaryPlayer = toShirtNum(outOfGamePlayer, lineUp);
        bool thereWasAnOutOfGame = (primaryPlayer != EVNT_NULL);
        uint outOfGameMinute = 0;
        if (thereWasAnOutOfGame) {
            if (matchLog[5] == 0) {
                return (events, "typeOfEvent = 0 is not allowed if thereWasAnOutOfGame");
            }
            uint typeOfEvent;
            if (matchLog[5] == 1) {
                typeOfEvent = EVNT_SOFT;
            } else if (matchLog[5] == 2) {
                typeOfEvent = EVNT_HARD;
            } else {
                if (!(matchLog[5] == 3)) { return (events, "Incorrect value of matchLog[5]"); }
                typeOfEvent = EVNT_RED;
            }
            outOfGameMinute = rounds2mins[matchLog[6]];
            MatchEvent thisEvent = new MatchEvent(outOfGameMinute, typeOfEvent, team, false, false, primaryPlayer, EVNT_NULL, NIL_PLAYER_ID, NIL_PLAYER_ID);
            newEvents.Add(thisEvent);
        }

        // First yellow card:
        uint yellowCardPlayer = matchLog[7];
        // convert player in the lineUp to shirtNum before storing it as match event:
        primaryPlayer = toShirtNum(yellowCardPlayer, lineUp);
        bool thereWasYellowCard = (primaryPlayer != EVNT_NULL);
        bool firstYellowCoincidesWithRed = false;
        if (thereWasYellowCard) {
            uint maxMinute = 45;
            if (yellowCardPlayer == outOfGamePlayer) {
                if (outOfGameMinute > 0) {
                    maxMinute = outOfGameMinute - 1;
                } else {
                    maxMinute = outOfGameMinute;
                }
                firstYellowCoincidesWithRed = true;
            }
            string salt = "c" + yellowCardPlayer.ToString();
            uint minute = GenerateRnd(seed, salt, maxMinute);
            uint typeOfEvent = EVNT_YELLOW;
            MatchEvent thisEvent = new MatchEvent(minute, typeOfEvent, team, false, false, primaryPlayer, EVNT_NULL, NIL_PLAYER_ID, NIL_PLAYER_ID);
            newEvents.Add(thisEvent);
        }

        // Second second yellow card:
        yellowCardPlayer = matchLog[8];
        // convert player in the lineUp to shirtNum before storing it as match event:
        primaryPlayer = toShirtNum(yellowCardPlayer, lineUp);
        thereWasYellowCard = (primaryPlayer != EVNT_NULL);
        if (thereWasYellowCard) {
            uint maxMinute = 45;
            uint typeOfEvent = EVNT_YELLOW;
            if (yellowCardPlayer == outOfGamePlayer) {
                if (firstYellowCoincidesWithRed) {
                    MatchEvent newEvent = new MatchEvent(outOfGameMinute, typeOfEvent, team, false, false, primaryPlayer, EVNT_NULL, NIL_PLAYER_ID, NIL_PLAYER_ID);
                    newEvents.Add(newEvent);
                    return (newEvents.ToArray(), "");
                } else {
                    maxMinute = outOfGamePlayer;
                }
            }
            string salt = "d" + yellowCardPlayer.ToString();
            uint minute = GenerateRnd(seed, salt, maxMinute);
            // convert player in the lineUp to shirtNum before storing it as match event:
            MatchEvent thisEvent = new MatchEvent(minute, typeOfEvent, team, false, false, primaryPlayer, EVNT_NULL, NIL_PLAYER_ID, NIL_PLAYER_ID);
            newEvents.Add(thisEvent);
        }
        return (newEvents.ToArray(), "");
    }

    private MatchEvent[] addSubstitutions(uint team, MatchEvent[] events, uint[] matchLog, uint[] rounds2mins, uint[] lineup, uint[] substitutions, uint[] subsRounds) {
        // matchLog:	9,10,11 ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen
        // halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
        List<MatchEvent> newEvents = events.ToList();

        for (uint i = 0; i < 3; i++) {
            bool subHappened = (matchLog[9+i] == 1);
            if (subHappened) {
                uint minute = rounds2mins[subsRounds[i]];
                int leavingPlayer = toShirtNum(substitutions[i], lineup);
                int enteringPlayer = toShirtNum(11+i, lineup);
                uint typeOfEvent = EVNT_SUBST;
                MatchEvent thisEvent = new MatchEvent(minute, typeOfEvent, team, false, false, leavingPlayer, enteringPlayer, NIL_PLAYER_ID, NIL_PLAYER_ID);
                newEvents.Add(thisEvent);
            }
        }
        return adjustSubstitutions(team, newEvents.ToArray());
    }

    // make sure that if a player that enters via a substitution appears in any other action (goal, pass, cards & injuries),
    // then the substitution time must take place at least before that minute.
    private MatchEvent[] adjustSubstitutions(uint team, MatchEvent[] events) {
        MatchEvent[] adjustedEvents = (MatchEvent[]) events.Clone();

        for (uint e = 0; e < events.Length; e++) {
            if ((events[e].Type == EVNT_SUBST) && (events[e].Team == team)) {
                int enteringPlayer = events[e].SecondaryPlayer;
                if (enteringPlayer != EVNT_NULL) {
                    uint enteringMin = events[e].Minute;
                    for (uint e2 = 0; e2 < events.Length; e2++) {
                        if ((e != e2) && (events[e2].Team == team) && (enteringPlayer == events[e2].PrimaryPlayer) && (enteringMin >= events[e2].Minute-1)) {
                            adjustedEvents[e].Minute = events[e2].Minute - 1;
                        }
                    }
                }
            }
        }
        return adjustedEvents;
    }

    private (MatchEvent[] events, string err) populateWithPlayerID(MatchEvent[] events, BigInteger[] homeTeamPlayerIDs, BigInteger[] visitorTeamPlayerIDs) {
        for (uint e = 0; e < events.Length; e++) {
            BigInteger[] primaryPlayerTeam;
            BigInteger[] secondaryPlayerTeam;
            BigInteger[] tacklerPlayerTeam;
            if (events[e].Team == 0) {
                primaryPlayerTeam = (BigInteger[]) homeTeamPlayerIDs.Clone();
                secondaryPlayerTeam = (BigInteger[]) homeTeamPlayerIDs.Clone();
                tacklerPlayerTeam = (BigInteger[]) visitorTeamPlayerIDs.Clone();
            } else {
                primaryPlayerTeam = (BigInteger[]) visitorTeamPlayerIDs.Clone();
                secondaryPlayerTeam = (BigInteger[]) visitorTeamPlayerIDs.Clone();
                tacklerPlayerTeam = (BigInteger[]) homeTeamPlayerIDs.Clone();
            }

            if (events[e].PrimaryPlayer != -1) {
                if (events[e].Type == EVNT_ATTACK && !events[e].ManagesToShoot) {
                    if (tacklerPlayerTeam[events[e].PrimaryPlayer] == NIL_PLAYER_ID) {
                        return (events, "inconsistent event: " + e.ToString());
                    }
                    events[e].PrimaryPlayerID = tacklerPlayerTeam[events[e].PrimaryPlayer];
                } else {
                    if (primaryPlayerTeam[events[e].PrimaryPlayer] == NIL_PLAYER_ID) {
                        return (events, "inconsistent event: " + e.ToString());
                    }
                    events[e].PrimaryPlayerID = primaryPlayerTeam[events[e].PrimaryPlayer];
                }
            }

            if (events[e].SecondaryPlayer == -1) { 
                // no secondary player
            } else if (events[e].SecondaryPlayer == EVNT_PENALTY) {
                 events[e].SecondaryPlayerID = new BigInteger(EVNT_PENALTY);
            } else {
                // default:
                if (secondaryPlayerTeam[events[e].SecondaryPlayer] == NIL_PLAYER_ID) {
                    return (events, "inconsistent event: " + e.ToString());
                }
                events[e].SecondaryPlayerID = secondaryPlayerTeam[events[e].SecondaryPlayer];
            }
        }
        return (events, "");
    }

    // INTERNAL UTILITIES

    private uint rightShiftAndMask(BigInteger encoded, int bitsToDisplace, int mask) { return (uint) ((encoded >> bitsToDisplace) & mask); }

    private ulong rightShiftAndMask64b(BigInteger encoded, int bitsToDisplace, ulong mask) { return (ulong) ((encoded >> bitsToDisplace) & mask); }

    private BigInteger OrWithLeftShift(BigInteger original, uint val, int bitsToDisplace) { 
        return original | ((new BigInteger(val)) << bitsToDisplace);
    }
}
