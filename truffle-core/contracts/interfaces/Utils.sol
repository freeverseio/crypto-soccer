pragma solidity >= 0.6.3;
import "../encoders/EncodingMatchLog.sol";
import "../storage/AssetsView.sol";

contract Utils is EncodingMatchLog, AssetsView{

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

        // getHalfTimeSubs: recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
        if (is2ndHalf) {
            decodedLog[12]  = uint32(getHalfTimeSubs(log, 0));
            decodedLog[13]  = uint32(getHalfTimeSubs(log, 1));
            decodedLog[14]  = uint32(getHalfTimeSubs(log, 2));
        }
        return decodedLog;
    }
    
    function fullDecodeSkills(uint256 encodedSkills) public pure returns(
        uint32[N_SKILLS] memory skills,
        uint16 dayOfBirth,
        uint8[4] memory birthTraits,
        uint256 playerId, 
        bool[3] memory aligned1stSubst1stRedCardLastGame,
        uint8[3] memory generationGamesNonStopInjuryWeeks
    ) {
        for (uint8 sk = 0; sk < N_SKILLS; sk++) skills[sk] = uint32(getSkill(encodedSkills, sk));

        dayOfBirth = uint16(getBirthDay(encodedSkills));

        birthTraits[0] = uint8(getPotential(encodedSkills));
        birthTraits[1] = uint8(getForwardness(encodedSkills));
        birthTraits[2] = uint8(getLeftishness(encodedSkills));
        birthTraits[3] = uint8(getAggressiveness(encodedSkills));
        
        playerId = getPlayerIdFromSkills(encodedSkills);
        
        aligned1stSubst1stRedCardLastGame[0] = getAlignedEndOfFirstHalf(encodedSkills);
        aligned1stSubst1stRedCardLastGame[1] = getSubstitutedFirstHalf(encodedSkills);
        aligned1stSubst1stRedCardLastGame[2] = getRedCardLastGame(encodedSkills);
        
        generationGamesNonStopInjuryWeeks[0] = uint8(getGeneration(encodedSkills));
        generationGamesNonStopInjuryWeeks[1] = getGamesNonStopping(encodedSkills);
        generationGamesNonStopInjuryWeeks[2] = getInjuryWeeksLeft(encodedSkills);
    }
    

}
