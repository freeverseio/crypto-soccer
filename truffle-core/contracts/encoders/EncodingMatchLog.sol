pragma solidity >= 0.6.3;

import "../encoders/EncodingMatchLogBase1.sol";
import "./EncodingMatchLogBase2.sol";
import "../encoders/EncodingMatchLogBase3.sol";

/**
 @title Library of pure functions to serialize matchLogs
 @author Freeverse.io, www.freeverse.io
 @dev Matchlog contain all the info determined in a given match
 @dev Due to contract-too-large-to-deploy, these functions had to be split Base1,..., Base4
 @dev so that each contract only inherits the subset truly required
 @dev note the difference in naming of functions addX vs setX
 @dev addX: adds to previous value, setX: replaces previous value
*/

/**
 Here is the full spec of this serialization:

 MAX_GOALS = ROUNDS_PER_MATCH = 12,  NO_OUT_OF_GAME_PLAYER = 14
 
 uint8 nGoals,  4b, offset 0 
 uint8 assistersIdx[MAX_GOALS], 4b each, offset 4
 uint8 shootersIdx[MAX_GOALS], 4b each, offset 52
 uint8 shooterFwdPos[MAX_GOALS], 2b each, offset 100
 bool[7] memory penaltyScored,  1b each, offset 128
 uint8[2] memory outOfGamePlayer 4b each, offset 155:135
    an int between [0, 13], 14 = NOTHING HAPPENED       
 uint8[2] memory outOfGameRounds, 4b each, offset 155:135 +4 
 uint8[2] memory outOfGameType, 2b each, offset 155:135 +8
    null:  0 (cannot return null if outOfGamePlayer != NO_OUT_OF_GAME_PLAYER)
    injuryHard:  1
    injuryLow:   2
    redCard:     3
 uint8[4] memory yellowCards, 4b each,  first 2 for first half, other for half 2, offset 165:145
 bool[2] memory yellowCardedDidNotFinish1stHalf,  1b each, offset 163  
 uint8[6] memory ingameSubsHappened,  2b each, offset 173
    the first 3 for half 1, the other for half 2.
    0: no change required, 1: change happened, 2: change could not happen, cancelled.  
 uint8[3] memory halfTimeSubstitutions,  5b each, offset 185
    0 means no subs, and we store here lineUp[p]+1 (where lineUp[p] = player shirt in the 25 that was substituted)
 nDefsHalf[2], 4b each, offset 200
 NTot2ndHalf, 4b offset 208
 winner, 2b, winner: 0 = home, 1 = away, 2 = draw, offset 212      
 teamSumSkills: 24b  offset 214
 trainingPoints, 12b, offset 238
 bool isHomeStadium,  1b each, offset 250
*/

contract EncodingMatchLog is EncodingMatchLogBase1, EncodingMatchLogBase2, EncodingMatchLogBase3 {

}
