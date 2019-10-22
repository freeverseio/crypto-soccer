const UNDEF = undefined;
    
async function encodeLog(encoding, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
    outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
    isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
    halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills
) {
    log = 0;
    log = await encoding.addNGoals(log, nGoals).should.be.fulfilled;
    for (p = 0; p <assistersIdx.length; p++) log = await encoding.addAssister(log, assistersIdx[p], p).should.be.fulfilled;
    for (p = 0; p <shootersIdx.length; p++) log = await encoding.addShooter(log, shootersIdx[p], p).should.be.fulfilled;
    for (p = 0; p <shooterForwardPos.length; p++) log = await encoding.addForwardPos(log, shooterForwardPos[p], p).should.be.fulfilled;
    for (p = 0; p <penalties.length; p++) {
      if (penalties[p]) log = await encoding.addScoredPenalty(log, p).should.be.fulfilled;
    }
    is2ndHalfs = [false, true];
    for (p = 0; p <outOfGames.length; p++) {
        log = await encoding.addOutOfGame(log, outOfGames[p], outOfGameRounds[p], typesOutOfGames[p], is2ndHalfs[p]).should.be.fulfilled;
    }
    for (p = 0; p <yellowCardedDidNotFinish1stHalf.length; p++) {
        if ( yellowCardedDidNotFinish1stHalf[p]) {
            log = await encoding.setYellowedDidNotFinished1stHalf(log, p).should.be.fulfilled;
        }
    }
    if (isHomeStadium) log = await encoding.addIsHomeStadium(log).should.be.fulfilled;
    for (p = 0; p <ingameSubs1.length; p++) log = await encoding.setInGameSubsHappened(log, ingameSubs1[p], p, is2nd = false).should.be.fulfilled;
    for (p = 0; p <ingameSubs2.length; p++) log = await encoding.setInGameSubsHappened(log, ingameSubs2[p], p, is2nd = true).should.be.fulfilled;
    for (p = 0; p <yellowCards1.length; p++) log = await encoding.addYellowCard(log, yellowCards1[p], p, is2nd = false).should.be.fulfilled;
    for (p = 0; p <yellowCards2.length; p++) log = await encoding.addYellowCard(log, yellowCards2[p], p, is2nd = true).should.be.fulfilled;
    for (p = 0; p <halfTimeSubstitutions.length; p++) log = await encoding.addHalfTimeSubs(log, halfTimeSubstitutions[p], p).should.be.fulfilled;
    
    log = await encoding.addNDefs(log, nDefs1, is2nd = false).should.be.fulfilled;
    log = await encoding.addNDefs(log, nDefs2, is2nd = true).should.be.fulfilled;
    log = await encoding.addNTot2ndHalf(log, nTot).should.be.fulfilled;
    log = await encoding.addWinner(log, winner).should.be.fulfilled;
    log = await encoding.addTeamSumSkills(log, teamSumSkills).should.be.fulfilled;
    return log;
}

async function checkExpectedLog(encoding, log, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
    outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
    isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
    halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills
) 
{
    if (nGoals != UNDEF) {
        result = await encoding.getNGoals(log);
        result.toNumber().should.be.equal(nGoals);
    }        
    if (assistersIdx != UNDEF) {
        for (p = 0; p <assistersIdx.length; p++) {
            result = await encoding.getAssister(log, p).should.be.fulfilled;
            result.toNumber().should.be.equal(assistersIdx[p]);
        }
    }        
    if (shootersIdx != UNDEF) {
        for (p = 0; p <shootersIdx.length; p++) {
            result = await encoding.getShooter(log, p).should.be.fulfilled;
            result.toNumber().should.be.equal(shootersIdx[p]);
        }
    }        
    if (shooterForwardPos != UNDEF) {
        for (p = 0; p <shooterForwardPos.length; p++) {
            result = await encoding.getForwardPos(log, p).should.be.fulfilled;
            result.toNumber().should.be.equal(shooterForwardPos[p]);
        }
    }        
    if (penalties != UNDEF) {
        for (p = 0; p <penalties.length; p++) {
            result = await encoding.getPenalty(log, p).should.be.fulfilled;
            result.should.be.equal(penalties[p]);
        }
    }        
    if (outOfGames != UNDEF) {
        is2ndHalfs = [false, true];
        for (p = 0; p <outOfGames.length; p++) {
            result = await encoding.getOutOfGame(log, is2ndHalfs[p]).should.be.fulfilled;
            result.player.toNumber().should.be.equal(outOfGames[p]);
            result.round.toNumber().should.be.equal(outOfGameRounds[p]);
            result.typeOfOutOfGame.toNumber().should.be.equal(typesOutOfGames[p]);
        }
    }        
    if (yellowCardedDidNotFinish1stHalf != UNDEF) {
        for (p = 0; p <yellowCardedDidNotFinish1stHalf.length; p++) {
            result = await encoding.getYellowedDidNotFinished1stHalf(log, p);
            result.should.be.equal(yellowCardedDidNotFinish1stHalf[p]);
        }
    }
    if (isHomeStadium != UNDEF) {
        result = await encoding.getIsHomeStadium(log).should.be.fulfilled;
        result.should.be.equal(isHomeStadium);
    }        
    if (ingameSubs1 != UNDEF) {
        for (p = 0; p <ingameSubs1.length; p++) {
            result = await encoding.getInGameSubsHappened(log, p, is2nd = false).should.be.fulfilled;
            result.toNumber().should.be.equal(ingameSubs1[p]);
        }
    }        
    if (ingameSubs2 != UNDEF) {
        for (p = 0; p <ingameSubs2.length; p++) {
            result = await encoding.getInGameSubsHappened(log, p, is2nd = true).should.be.fulfilled;
            result.toNumber().should.be.equal(ingameSubs2[p]);
        }
    }        
    if (yellowCards1 != UNDEF) {
        for (p = 0; p <yellowCards1.length; p++) {
            result = await encoding.getYellowCard(log, p, is2nd = false).should.be.fulfilled;
            result.toNumber().should.be.equal(yellowCards1[p]);
        }
    }        
    if (yellowCards2 != UNDEF) {
        for (p = 0; p <yellowCards2.length; p++) {
            result = await encoding.getYellowCard(log, p, is2nd = true).should.be.fulfilled;
            result.toNumber().should.be.equal(yellowCards2[p]);
        }
    }        
    if (halfTimeSubstitutions != UNDEF) {
        for (p = 0; p <halfTimeSubstitutions.length; p++) {
            result = await encoding.getHalfTimeSubs(log, p).should.be.fulfilled;
            result.toNumber().should.be.equal(halfTimeSubstitutions[p]);
        }
    }        
    if (nDefs1 != UNDEF) {
        result = await encoding.getNDefs(log, is2nd = false).should.be.fulfilled;
        result.toNumber().should.be.equal(nDefs1);
    }        
    if (nDefs2 != UNDEF) {
        result = await encoding.getNDefs(log, is2nd = true).should.be.fulfilled;
        result.toNumber().should.be.equal(nDefs2);
    }        
    if (nTot != UNDEF) {
        result = await encoding.getNTot2ndHalf(log).should.be.fulfilled;
        result.toNumber().should.be.equal(nTot);
    }        
    if (winner != UNDEF) {
        result = await encoding.getWinner(log).should.be.fulfilled;
        result.toNumber().should.be.equal(winner);
    }        
    if (teamSumSkills != UNDEF) {
        result = await encoding.getTeamSumSkills(log).should.be.fulfilled;
        result.toNumber().should.be.equal(teamSumSkills);
    }        
}    
  
  module.exports = {
    encodeLog,
    checkExpectedLog
  }