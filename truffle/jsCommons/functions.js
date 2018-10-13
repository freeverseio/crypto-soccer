/*
    File for JS functions used by multiple tests, or even by the UI
    One restriction: functions in this file cannot use other functions in this file :-(
*/

function catchPlayerIdxFromEvent(logs) {
    var playerIdx = -1;
    for (var i = 0; i < logs.length; i++) {
        var log = logs[i];
        if (log.event == "PlayerCreation") {
            //console.log("...created player with idx = " + log.args.playerIdx.toNumber());
            playerIdx = log.args.playerIdx.toNumber();
        }
    }
    return playerIdx;
}

function unixMonthToAge(unixMonthOfBirth) {
    // in July 2018, we are at month 582 after 1970.
    age = (582 - unixMonthOfBirth)/12;
    return parseInt(age*10)/10;
}
  
async function createTeam(instance, teamName, playerBasename, maxPlayersPerTeam, playerRoles ) {
    var newTeamIdx = await instance.test_getNCreatedTeams.call(); 
    console.log("creating team: " + teamName);
    await instance.test_createTeam(teamName);
    const userChoice=1;
  
    for (var p=0; p<maxPlayersPerTeam; p++) {
        thisName = playerBasename + p.toString();
        var tx = await instance.test_createBalancedPlayer(
            thisName,
            newTeamIdx,
            userChoice,
            p,
            playerRoles[p]
        );
    }
    nCreatedPlayers = await instance.test_getNCreatedPlayers.call();
    console.log('Final nPlayers in the entire game = ' + nCreatedPlayers);
    return newTeamIdx;
}

function createAlineacion(nDef,nMid,nAtt) {
    alineacion = [0];
    for (var p = 0; p<nDef; p++) {
        alineacion.push(k.RoleDef);
    }
    for (var p = 0; p<nMid; p++) {
        alineacion.push(k.RoleMid);
    }
    for (var p = 0; p<nAtt; p++) {
        alineacion.push(k.RoleAtt);
    }
    return alineacion;
}

async function getRandomNumbers(instance, nRounds, rndSeed)
{
  var result = []
  bits = 10
  var hash = await instance.test_computeKeccak256ForNumber(rndSeed);
  var rndNums1= await instance.test_decode(nRounds, hash , bits);
  hash = await instance.test_computeKeccak256ForNumber(rndSeed+1);
  var rndNums2= await instance.test_decode(nRounds, hash, bits);
  hash = await instance.test_computeKeccak256ForNumber(rndSeed+2);
  var rndNums3= await instance.test_decode(nRounds, hash, bits);
  hash = await instance.test_computeKeccak256ForNumber(rndSeed+3);
  var rndNums4= await instance.test_decode(nRounds, hash, bits);
  result.push(rndNums1);
  result.push(rndNums2);
  result.push(rndNums3);
  result.push(rndNums4);
  return result;
}



  functions2export = {
    createTeam : createTeam,
    catchPlayerIdxFromEvent : catchPlayerIdxFromEvent,
    createAlineacion : createAlineacion,
    getRandomNumbers : getRandomNumbers,
    unixMonthToAge : unixMonthToAge      
}

module.exports = functions2export;