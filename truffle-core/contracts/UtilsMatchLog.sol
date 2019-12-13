pragma solidity >=0.4.21 <0.6.0;
import "./EncodingMatchLog.sol";

// MatchLog Contains:
        // uint8 nGoals, // 4b, offset 0
        // uint8 assistersIdx[12], 4b each, offset 4
        // uint8 shootersIdx[12], 4b each, offset 52
        // uint8 shooterFwdPos[12], 2b each, offset 100
        // bool[7] memory penalties, // 1b each, offset 128
        // uint8[2] memory outOfGames 4b each
        // uint8[4] memory yellowCards, 4b each, // first 2 for first half, other for half 2
        // uint8[2] memory outOfGameRounds,  
        // uint8[2] memory typesOutOfGames, 
        // bool[3] memory yellowCardedDidNotFinish1stHalf, 1b each,
        // bool isHomeStadium, // 1b each,
        // uint8[3] memory halfTimeSubstitutions, // 4b each, offset 201, the first 3 for half 1, the other for half 2
        // uint8[6] memory ingameSubs, // 2b each, offset 189
        //                             //  ...the first 3 for half 1, the other for half 2.
        //                             // ...0: no change required, 1: change happened, 2: change could not happen  
        // uint8[4] memory numDefTotWinner
        //                             // [4b, 4b, 4b, 2b], offset 213
        //                             // [nDefsHalf1, nDefsHalf2, nTotHalf2, winner]
        //                             // winner: 0 = home, 1 = away, 2 = draw
        // teamSumSkills: 24b // offset 227
        
// The main function in this contact exports data useful for creating game visualization:
// 
//  teamSumSkills 
//  winner: 0 = home, 1 = away, 2 = draw
//  nGoals
//  trainingPoints
//  uint8 memory outOfGames
//  uint8 memory typesOutOfGames, 
//  uint8 memory outOfGameRounds
//  uint8[2] memory yellowCards
//  uint8[3] memory ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen  
//  uint8[3] memory halfTimeSubstitutions: 0...10 the player in the starting 11 that was changed during half time
//  tot = 4 + 3 + 2 + 2*3  = 15

contract UtilsMatchLog is EncodingMatchLog{

    uint256 private constant ONE256       = 1; 
    uint256 private constant CHG_HAPPENED        = uint256(1); 
    uint256 private constant CHG_CANCELLED       = uint256(2); 

    function fullDecodeMatchLog(uint256 log, bool is2ndHalf) public pure returns (uint32[31] memory decodedLog) {
        decodedLog[0] = uint32(getTeamSumSkills(log));
        decodedLog[1] = uint32(getWinner(log));
        decodedLog[2] = uint32(getNGoals(log));
        decodedLog[3] = uint32(getTrainingPoints(log));
        
        (uint8 player, uint8 round, uint8 typeOfOutOfGame) = getOutOfGame(log, is2ndHalf);
        decodedLog[4] = uint32(player);
        decodedLog[5] = uint32(typeOfOutOfGame);
        decodedLog[6] = uint32(round);
    
        decodedLog[7] = uint32(getYellowCard(log, 0, is2ndHalf));
        decodedLog[8] = uint32(getYellowCard(log, 1, is2ndHalf));
        
        decodedLog[9]  = uint32(getInGameSubsHappened(log, 0, is2ndHalf));
        decodedLog[10] = uint32(getInGameSubsHappened(log, 1, is2ndHalf));
        decodedLog[11] = uint32(getInGameSubsHappened(log, 2, is2ndHalf));

        if (is2ndHalf) {
            decodedLog[12]  = uint32(getHalfTimeSubs(log, 0));
            decodedLog[13]  = uint32(getHalfTimeSubs(log, 1));
            decodedLog[14]  = uint32(getHalfTimeSubs(log, 2));
        }
    }

}