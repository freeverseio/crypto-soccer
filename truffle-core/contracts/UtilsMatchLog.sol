pragma solidity >=0.4.21 <0.6.0;
import "./EncodingMatchLog.sol";

contract UtilsMatchLog is EncodingMatchLog{

    function fullDecodeMatchLog(uint256 log, bool is2ndHalf) public pure returns (uint32[15] memory decodedLog) {
        decodedLog[0] = uint32(getTeamSumSkills(log));
        decodedLog[1] = uint32(getWinner(log));
        decodedLog[2] = uint32(getNGoals(log));
        if (is2ndHalf) decodedLog[3] = uint32(getTrainingPoints(log));
        
        decodedLog[4] = uint32(getOutOfGamePlayer(log, is2ndHalf));
        decodedLog[5] = uint32(getOutOfGameType(log, is2ndHalf));
        decodedLog[6] = uint32(getOutOfGameRound(log, is2ndHalf));
    
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
        return decodedLog;
    }

}