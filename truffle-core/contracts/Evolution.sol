pragma solidity >=0.5.12 <=0.6.3;

import "./Assets.sol";
import "./EngineLib.sol";
import "./EncodingMatchLog.sol";
import "./Engine.sol";
import "./EncodingTPAssignment.sol";
import "./EncodingSkillsSetters.sol";
import "./EncodingTacticsPart1.sol";

contract Evolution is EncodingMatchLog, EngineLib, EncodingTPAssignment, EncodingSkillsSetters, EncodingTacticsPart1 {
    uint8 constant private PLAYERS_PER_TEAM_MAX = 25;

    // uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_OUT_OF_GAME_PLAYER  = 14;   // noone saw a card
    uint8 public constant RED_CARD = 3;   // noone saw a card
    uint8 public constant SOFT_INJURY  = 1;  
    uint8 public constant HARD_INJURY  = 2;  
    uint8 public constant WEEKS_HARD_INJ = 5;  // weeks a player is out when suffered a hard injury
    uint8 public constant WEEKS_SOFT_INJ = 2;  // weeks a player is out when suffered a soft injury
    uint8 private constant CHG_HAPPENED        = uint8(1); 
    // uint8 constant public N_SKILLS = 5;

    function updateSkillsAfterPlayHalf(
        uint256[PLAYERS_PER_TEAM_MAX] memory skills,
        uint256 matchLog,
        uint256 tactics,
        bool is2ndHalf
    ) 
        public
        pure
        returns (uint256[PLAYERS_PER_TEAM_MAX] memory)
    {
        // after 1st Half, update:
        //  - subtDuringFirstHalf, alignedEndOfFirstHalf => properly update
        //  - redCards, injury => add if any of this happens
        //  - adds +2 to daysNonStopping to all linedUp players
        // after 2nd Half, update:
        //  - subtDuringFirstHalf = 0, alignedEndOfFirstHalf = 0
        //  - redCards: 
        //      - set all to false unless it happens in 1st or 2nd half
        //  -  injury => 
        //      - decrease by one unless it happens in 1st or 2nd half
        if (!is2ndHalf) {
            writeOutOfGameInSkills(skills, tactics, matchLog, false);
            writeFirstHalfLineUp(skills, tactics, matchLog);
        }
        else {
            decreaseOutOfGames(skills);
            writeOutOfGameInSkills(skills, tactics, matchLog, false);
            writeOutOfGameInSkills(skills, tactics, matchLog, true);
            updateGamesNonStopping2ndHalf(skills, tactics, matchLog); //TODO: rename to "update"
            updatePlayerAtEndOfMatch(skills);
        }
        return skills;
    }

    function updateGamesNonStopping2ndHalf(
        uint256[PLAYERS_PER_TEAM_MAX] memory skills, 
        uint256 tactics, 
        uint256 matchLog 
    ) private pure 
    {
        uint8[3] memory joinedAt2ndHalf;
        uint8 nJoined = 0;
        // first increase +2 the gamesNonStopping for those who joined
        (,,uint8[14] memory lineUp,,) = decodeTactics(tactics);
        for (uint8 posInHalf = 0; posInHalf < 3; posInHalf++) {
            // First: those who joined at half time:
            // note that getHalfTimeSubs: returns lineUp[p]+1 for halftime subs, 0 = NO_SUBS
            uint8 enteringPlayer = getHalfTimeSubs(matchLog, posInHalf); 
            if (enteringPlayer > 0) {
                joinedAt2ndHalf[nJoined] = enteringPlayer-1;
                nJoined += 1;
            }
            if (getInGameSubsHappened(matchLog, posInHalf, true) == CHG_HAPPENED) {
                joinedAt2ndHalf[nJoined]  = lineUp[11 + posInHalf];
                nJoined += 1;
            }
        }
        require(nJoined <= 3, "Too many changes detected in this match!");
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if (skills[p] != 0) {
                if (hasPlayedThisMatch(skills[p], p, joinedAt2ndHalf)) {
                    skills[p] = increaseGamesNonStopping(skills[p]);
                } else {
                    skills[p] = setGamesNonStopping(skills[p], 0); 
                }
            }
        }
    }

    function hasPlayedThisMatch(uint256 skills, uint8 p, uint8[3] memory joinedAt2ndHalf) public pure returns(bool) {
        return (
                    getAlignedEndOfFirstHalf(skills) || 
                    getSubstitutedFirstHalf(skills) ||
                    p == joinedAt2ndHalf[0] ||
                    p == joinedAt2ndHalf[1] ||
                    p == joinedAt2ndHalf[2]
                );
    }

    function updatePlayerAtEndOfMatch(uint256[PLAYERS_PER_TEAM_MAX] memory skills) private pure {
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if (skills[p] != 0) {
                skills[p] = setAlignedEndOfFirstHalf(skills[p], false);
                skills[p] = setSubstitutedFirstHalf(skills[p], false);
            }
        }
    }
    
    function writeFirstHalfLineUp(
        uint256[PLAYERS_PER_TEAM_MAX] memory skills, 
        uint256 tactics, 
        uint256 matchLog 
    ) 
        private 
        pure 
    {
        (uint8[3] memory  substitutions,,uint8[14] memory lineUp,,) = decodeTactics(tactics);
        // NO_LINEUP = 25, NO_SUBS = 11
        for (uint8 p = 0; p < 11; p++) {
            uint8 linedUp = lineUp[p];
            if ((linedUp < NO_LINEUP) && (skills[linedUp] != 0)) {
                skills[linedUp] = setAlignedEndOfFirstHalf(skills[linedUp], true);
                skills[linedUp] = setSubstitutedFirstHalf(skills[linedUp], false);
            }
        }
        for (uint8 posInHalf = 0; posInHalf < 3; posInHalf++) {
            if (getInGameSubsHappened(matchLog, posInHalf, false) == CHG_HAPPENED) {
                uint8 leavingFieldPlayer    = lineUp[substitutions[posInHalf]];
                uint8 enteringFieldPlayer   = lineUp[11 + posInHalf];
                if (skills[leavingFieldPlayer] != 0) {
                    skills[leavingFieldPlayer]  = setAlignedEndOfFirstHalf(skills[leavingFieldPlayer], false);
                    skills[leavingFieldPlayer]  = setSubstitutedFirstHalf(skills[leavingFieldPlayer], true);
                }
                if (skills[enteringFieldPlayer] != 0) {
                    skills[enteringFieldPlayer] = setAlignedEndOfFirstHalf(skills[enteringFieldPlayer], true);
                    skills[enteringFieldPlayer] = setSubstitutedFirstHalf(skills[enteringFieldPlayer], false);
                }
            }
        }
    }    
    
    // we increase 2 units, so that we get a maximum of 7. 
    function increaseGamesNonStopping(uint256 skills) public pure returns (uint256) {
        uint8 gamesNonStopping = getGamesNonStopping(skills);
        if (gamesNonStopping < 7) return setGamesNonStopping(skills, gamesNonStopping + 1); 
        else return skills;
    }

    function writeOutOfGameInSkills(uint256[PLAYERS_PER_TEAM_MAX] memory skills, uint256 tactics, uint256 matchLog, bool is2ndHalf) private pure {
        (,,uint8[14] memory lineUp,,) = decodeTactics(tactics);
        // check if there was an out of player event:
        uint8 outOfGamePlayer = uint8(getOutOfGamePlayer(matchLog, is2ndHalf));
        if (outOfGamePlayer == NO_OUT_OF_GAME_PLAYER) return;
        // convert outOfGamePlayer [0...13] to the index that points to the skills in the team [0,..24]
        outOfGamePlayer = lineUp[outOfGamePlayer];
        if (skills[outOfGamePlayer] == 0) return;
        uint8 outOfGameType = lineUp[uint8(getOutOfGameType(matchLog, is2ndHalf))];
        if (outOfGameType == RED_CARD) {
            skills[outOfGamePlayer] = setRedCardLastGame(skills[outOfGamePlayer], true);
        }
        else if (outOfGameType == HARD_INJURY) {
            skills[outOfGamePlayer] = setInjuryWeeksLeft(skills[outOfGamePlayer], WEEKS_HARD_INJ);
        }
        else if (outOfGameType == SOFT_INJURY) {
            skills[outOfGamePlayer] = setInjuryWeeksLeft(skills[outOfGamePlayer], WEEKS_SOFT_INJ);
        }
    }
        
    // at the end of a match, decrease the weeks left from injury, and set redcards = false.
    // the function called right after this one will add the redcards of this particular game where appropriate
    function decreaseOutOfGames(uint256[PLAYERS_PER_TEAM_MAX] memory skills) public pure {
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if (skills[p] != 0) {
                skills[p] = setRedCardLastGame(skills[p], false);
                if (getInjuryWeeksLeft(skills[p]) != 0) {
                    skills[p] = setInjuryWeeksLeft(skills[p], getInjuryWeeksLeft(skills[p])-1);
                }
            }
        }
    }
}

