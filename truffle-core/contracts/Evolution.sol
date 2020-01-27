pragma solidity >=0.5.12 <0.6.2;

import "./Assets.sol";
import "./EngineLib.sol";
import "./EncodingMatchLog.sol";
import "./Engine.sol";
import "./EncodingTPAssignment.sol";
import "./EncodingSkillsSetters.sol";

contract Evolution is EncodingMatchLog, EngineLib, EncodingTPAssignment, EncodingSkillsSetters {

    // uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_OUT_OF_GAME_PLAYER  = 14;   // noone saw a card
    uint8 public constant RED_CARD = 3;   // noone saw a card
    uint8 public constant SOFT_INJURY  = 1;  
    uint8 public constant HARD_INJURY  = 2;  
    uint8 public constant WEEKS_HARD_INJ = 5;  // weeks a player is out when suffered a hard injury
    uint8 public constant WEEKS_SOFT_INJ = 2;  // weeks a player is out when suffered a soft injury
    uint8 private constant CHG_HAPPENED        = uint8(1); 
    // uint8 constant public N_SKILLS = 5;

    function updateStatesAfterPlayHalf(
        uint256[PLAYERS_PER_TEAM_MAX] memory states,
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
        // after 2nd Half, update:
        //  - subtDuringFirstHalf = 0, alignedEndOfFirstHalf = 0
        //  - redCards: 
        //      - set all to false unless it happens in 1st or 2nd half
        //  -  injury => 
        //      - decrease by one unless it happens in 1st or 2nd half
        if (!is2ndHalf) {
            writeOutOfGameState(states, tactics, matchLog, false);
            writeFirstHalfLineUp(states, tactics, matchLog);
        }
        else {
            decreaseOutOfGames(states);
            writeOutOfGameState(states, tactics, matchLog, false);
            writeOutOfGameState(states, tactics, matchLog, true);
            resetFirstHalfLineUp(states);
        }
        return states;
    }

    function resetFirstHalfLineUp(uint256[PLAYERS_PER_TEAM_MAX] memory states) private pure {
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if (states[p] != 0) {
                states[p] = setAlignedEndOfFirstHalf(states[p], false);
                states[p] = setSubstitutedFirstHalf(states[p], false);
            }
        }
    }
    
    function writeFirstHalfLineUp(
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
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
            if (linedUp < NO_LINEUP) {
                states[linedUp] = setAlignedEndOfFirstHalf(states[linedUp], true);
                states[linedUp] = setSubstitutedFirstHalf(states[linedUp], false);
            }
        }
        for (uint8 posInHalf = 0; posInHalf < 3; posInHalf++) {
            if (getInGameSubsHappened(matchLog, posInHalf, false) == CHG_HAPPENED) {
                uint8 leavingFieldPlayer    = substitutions[posInHalf];
                uint8 enteringFieldPlayer   = lineUp[10 + posInHalf];
                states[leavingFieldPlayer]  = setAlignedEndOfFirstHalf(states[leavingFieldPlayer], false);
                states[leavingFieldPlayer]  = setSubstitutedFirstHalf(states[leavingFieldPlayer], true);
                states[enteringFieldPlayer] = setAlignedEndOfFirstHalf(states[enteringFieldPlayer], true);
                states[enteringFieldPlayer] = setSubstitutedFirstHalf(states[enteringFieldPlayer], false);
            }
        }
    }    
    
    function writeOutOfGameState(uint256[PLAYERS_PER_TEAM_MAX] memory states, uint256 tactics, uint256 matchLog, bool is2ndHalf) private pure {
        (,,uint8[14] memory lineUp,,) = decodeTactics(tactics);
        // check if there was an out of player event:
        uint8 outOfGamePlayer = uint8(getOutOfGamePlayer(matchLog, is2ndHalf));
        if (outOfGamePlayer == NO_OUT_OF_GAME_PLAYER) return;
        // convert outOfGamePlayer [0...13] to the index that points to the state in the team [0,..24]
        outOfGamePlayer = lineUp[outOfGamePlayer];
        uint8 outOfGameType = lineUp[uint8(getOutOfGameType(matchLog, is2ndHalf))];
        if (outOfGameType == RED_CARD) {
            states[outOfGamePlayer] = setRedCardLastGame(states[outOfGamePlayer], true);
        }
        else if (outOfGameType == HARD_INJURY) {
            states[outOfGamePlayer] = setInjuryWeeksLeft(states[outOfGamePlayer], WEEKS_HARD_INJ);
        }
        else if (outOfGameType == SOFT_INJURY) {
            states[outOfGamePlayer] = setInjuryWeeksLeft(states[outOfGamePlayer], WEEKS_SOFT_INJ);
        }
    }
        
    // at the begining of a match, decrease the weeks left from injury and redcards.
    function decreaseOutOfGames(uint256[PLAYERS_PER_TEAM_MAX] memory states) public pure {
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if (states[p] != 0) {
                states[p] = setRedCardLastGame(states[p], false);
                if (getInjuryWeeksLeft(states[p]) != 0) {
                    states[p] = setInjuryWeeksLeft(states[p], getInjuryWeeksLeft(states[p])-1);
                }
            }
        }
    }


}

