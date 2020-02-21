pragma solidity >=0.5.12 <=0.6.3;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

import "./EncodingTacticsPart1.sol";
import "./EncodingTacticsPart2.sol";
 
/**
    * @dev Tactics serializes a total of 157 bits.
    *               The first 110: 3 * 4 + 3 * 4 + 14*5 + 10 + 6:
    *               The following ones: 50 + 13 + 32
    *      substitutions[3]    = 4 bit each = [3 different nums from 0 to 10], with 11 = no subs
    *      subsRounds[3]       = 4 bit each = [3 different nums from 0 to 11], round at which subs are to happen
    *      lineup[14]          = 5 bit each = [playerIdxInTeam1, ..., ]
    *      extraAttack[10]     = 1 bit each, 0: normal, 1: player has extra attack duties
    *      tacticsId           = 6 bit (0 = 442, 1 = 541, ...
    *
    *      ...added by shop: offests: 110 -> 159, 160 -> 172, 173 -> 204...
    *      staminaRecovery[25] = 2 bit each => 50b ( 0 = none, 1 = 2 games, 2 = 4 games, 3 = full recovery)
    *      itemId              = 13 bit 
    *      itemEncodedBoost    = 32 bit
**/
contract EncodingTactics is EncodingTacticsPart1, EncodingTacticsPart2 {

}
