const NULL_BYTES32 = web3.eth.abi.encodeParameter('bytes32','0x0');

const merkleUtils = require('../utils/merkleUtils.js');
const teamsUtils = require('../utils/teamsUtils.js');
const fs = require('fs');

const nLeafs = 1024;
const nMatchdays = 14;
const nMatchesPerDay = 4;
const nTeamsInLeague = 8;
const nMatchesPerLeague = nMatchesPerDay * nMatchdays;
const nPlayersInTeam = 25;

function readTeam442(filepath = 'test/testdata/team442.json'){
  return JSON.parse(fs.readFileSync(filepath, 'utf8'));
}

function sortDec(array) { 
  arr = [...array];
  return arr.sort(function(a, b){return b-a});
}

function uint256(x) { return web3.eth.abi.encodeParameter('uint256', x); }

function joinTacticsAndTPs(tactics, TPs) {
  n = tactics.length;
  assert.equal(n, TPs.length, "tactics and TPs should have equal length");
  joined = Array.from(new Array(n), (x,i) => 0);
  for (i = 0; i < tactics.length; i++) {
    joined [2*i] = tactics[i];
    joined [2*i + 1] = TPs[i];
  }
  return joined;
}

// returns, for a given TZ, [nActiveCountry0,... 1023]
function createUniverse(nLevels) {
  nCountries = 1024;
  nTeamsInLeague = 8;
  activeLeagues = Array.from(new Array(nCountries), (x,i) => 2 * ((i+1) % 3));
  activeTeams = Array.from(activeLeagues, (x,i) => x * nTeamsInLeague);
  orgMap = [];
  userActions = [];
  for (c = 0; c < 1024; c++) {
    teamIdxs = Array.from(new Array(activeTeams[c]), (x,i) => uint256(i));
    tactics = Array.from(new Array(activeTeams[c]), (x,i) => uint256(-i));
    TPs = Array.from(new Array(activeTeams[c]), (x,i) => uint256(-2*i-1));
    if (c % 2 == 0) { 
      orgMap = orgMap.concat(teamIdxs); 
      userActions = userActions.concat(joinTacticsAndTPs(tactics, TPs));
      
    } else  { 
      orgMap = orgMap.concat(sortDec(teamIdxs)); 
      userActions = userActions.concat(sortDec(joinTacticsAndTPs(tactics, TPs)));
    }
  }
  return [activeTeams, orgMap, userActions];
}


function zeroPadToLength(x, desiredLength) {
  return x.concat(Array.from(new Array(desiredLength - x.length), (x,i) => 0))
}

// var leagueData = {
//     seeds: [], // [2 * nMatchDays]
//     teamIds: [], // [nTeamsInLeague]
//     startTimes: [], // [2 * nMatchDays]
//     teamStates: [], // [2 * nMatchdays + 1][nTeamsInLeague][PLAYERS_PER_TEAM_MAX]
//     matchLogs: [], // [2 * nMatchdays+ 1][nTeamsInLeague]
//     results: [], // [nMatchesPerLeague][2]  -> goals per team per match
//     points: [], // [2 * nMatchdays][nTeamsInLeague]
//     tactics: [], // [2 * nMatchdays + 1][nTeamsInLeague]
//     trainings: [] // [2 * nMatchdays + 1][nTeamsInLeague]
// }
// - Data[1024] = [League[512], Team$_{i, aft}$[32], Team$_{i, bef}$[32]]
// League[128] = leafsLeague[128] = [Points[team=0,..,7], Goals[56][2]]
// 
// returns leafs AFTER having played the matches at matchday = day, half = half.
//  - sorting results:
//      - idx = day * nMatchesPerDay * 2 + matchInDay * 2 + teamHomeOrAway
function buildLeafs(leagueDataIn, day, half) {
  // oprate on cloned input for safety:
  lData = clone(leagueDataIn)
  var isNoPointsYet = (half == 0) && (day == 0);
  if (isNoPointsYet) { 
      leafs =  Array.from(new Array(nTeamsInLeague), (x,i) => 0);
  } else {
      lastDayToCount = (half == 0) ? day - 1 : day;
      leafs = lData.points[lastDayToCount]; 
      for (d = 0; d <= lastDayToCount; d++) {
          for (m = 0; m < nMatchesPerDay; m++) {
              leafs.push(lData.results[d * nMatchesPerDay + m][0]);
              leafs.push(lData.results[d * nMatchesPerDay + m][1]);
          }
      }
  }
  leafs = zeroPadToLength(leafs, 128);
  for (team = 0; team < nTeamsInLeague; team++) {
      for (extraHalf = 0; extraHalf < 2; extraHalf++) {
          teamData = [];
          for (p = 0; p < nPlayersInTeam; p++) {
              teamData.push(lData.teamStates[2*day + half + extraHalf][team][p])
          }
          teamData.push(lData.tactics[2*day + half + extraHalf][team]);
          teamData.push(lData.trainings[2*day + half + extraHalf][team]);
          teamData.push(lData.matchLogs[2*day + half + extraHalf][team]);
          leafs = leafs.concat(zeroPadToLength(teamData, 32));
      }
  }
  return zeroPadToLength(leafs, nLeafs);
}

function vec2str(y) {
  yStr = [...y];
  for (i = 0; i < y.length; i++) yStr[i] = y[i].toString();
  return yStr;
}

function clone(a) {
  return JSON.parse(JSON.stringify(a));
}

async function createLeagueData(champs, play, encodeLog, now, teamState442, teamId) {
  let secsBetweenMatches = 12*3600;
  var leagueData = {
      seeds: [], // [2 * nMatchDays]
      teamIds: [], // [nTeamsInLeague]
      startTimes: [], // [2 * nMatchDays]
      teamStates: [], // [1 + 2 * nMatchdays][nTeamsInLeague][PLAYERS_PER_TEAM_MAX]
      matchLogs: [], // [1 + 2 * nMatchdays][nTeamsInLeague]
      results: [], // [nMatchesPerLeague][2]  ->  per team per match
      points: [], // [2 * nMatchdays][nTeamsInLeague]
      tactics: [], // [2 * nMatchdays + 1][nTeamsInLeague]
      trainings: [] // [2 * nMatchdays + 1][nTeamsInLeague]
  }
  // on starting points: if we query computeLeagueLeaderBoard, I would get 
  // a non-null value, sorting because of all tied, which would depend on a seed.
  // we don't have that seed before a match starts, so we set all points to 0.

  leagueData.seeds = Array.from(new Array(2 * nMatchdays), (x,i) => web3.utils.keccak256(i.toString()).toString());
  leagueData.startTimes = Array.from(new Array(2 * nMatchdays), (x,i) => now + i * secsBetweenMatches);
  allMatchLogs = Array.from(new Array(nTeamsInLeague), (x,i) => 0);
  leagueData.matchLogs.push([...allMatchLogs]);
  teamState442 = vec2str(teamState442);
  allTeamsSkills = Array.from(new Array(nTeamsInLeague), (x,i) => teamState442);
  leagueData.teamStates.push([...allTeamsSkills]);
  // nosub = [NO_SUBST, NO_SUBST, NO_SUBST];
  // tact = await engine.encodeTactics(nosub , ro = [0, 0, 0], setNoSubstInLineUp(lineupConsecutive, nosub), extraAttackNull, tacticsId = 0).should.be.fulfilled;
  leagueData.teamIds = Array.from(new Array(nTeamsInLeague), (x,i) => teamId.toNumber() + i);
  leagueData.results = Array.from(new Array(nMatchesPerLeague), (x,i) => [0,0]);

  // tactics and trainings start at all 0 (undefined until we play the first match)
  leagueData.tactics.push(Array.from(new Array(nTeamsInLeague), (x,i) => 0));
  leagueData.trainings.push(Array.from(new Array(nTeamsInLeague), (x,i) => 0));
  // same tactics and trainings for all matchdays:
  tact = tactics442NoChanges.toString();
  for (day = 0; day < 2 * nMatchdays; day++) {
      leagueData.tactics.push(Array.from(new Array(nTeamsInLeague), (x,i) => tact));
  }
  // trainings after 2nd half are required to be 0
  for (day = 0; day < nMatchdays; day++) {
      leagueData.trainings.push(Array.from(new Array(nTeamsInLeague), (x,i) => almostNullTraning.toString()));
      leagueData.trainings.push(Array.from(new Array(nTeamsInLeague), (x,i) => 0));
  }

  // we just need to build, across the league: teamStates, points, teamIds
  // for (day = 0; day < 1; day++) {
  for (day = 0; day < nMatchdays; day++) {
      // 1st half
      for (matchIdxInDay = 0; matchIdxInDay < nMatchesPerDay; matchIdxInDay++) {
          var {0: t0, 1: t1} = await champs.getTeamsInLeagueMatch(day, matchIdxInDay).should.be.fulfilled;
          t0 = t0.toNumber();
          t1 = t1.toNumber();
          console.log("day:", day, ", matchIdxInDay:", matchIdxInDay, ", half 0,  teams:", t0, t1);
          var {0: newSkills, 1: newLogs} =  await play.play1stHalfAndEvolve(
              leagueData.seeds[2 * day], leagueData.startTimes[2 * day], 
              [allTeamsSkills[t0], allTeamsSkills[t1]], 
              [leagueData.teamIds[t0], leagueData.teamIds[t1]], 
              [leagueData.tactics[2 * day + 1][t0], leagueData.tactics[2 * day + 1][t1]], 
              [allMatchLogs[t0], allMatchLogs[t1]],
              [is2nd = false, isHom = true, isPlay = false],
              [tp = 0, tp = 0]
          ).should.be.fulfilled;
          allTeamsSkills[t0] = vec2str(newSkills[0]);
          allTeamsSkills[t1] = vec2str(newSkills[1]);
          allMatchLogs[t0] = newLogs[0].toString();
          allMatchLogs[t1] = newLogs[1].toString();
      }
      leagueData.teamStates.push([...allTeamsSkills]);        
      leagueData.matchLogs.push([...allMatchLogs]);        
      // 2nd half
      for (matchIdxInDay = 0; matchIdxInDay < nMatchesPerDay; matchIdxInDay++) {
          var {0: t0, 1: t1} = await champs.getTeamsInLeagueMatch(day, matchIdxInDay).should.be.fulfilled;
          t0 = t0.toNumber();
          t1 = t1.toNumber();
          console.log("day:", day, ", matchIdxInDay:", matchIdxInDay, ", half 1,  teams:", t0, t1);
          var {0: newSkills, 1: newLogs} =  await play.play2ndHalfAndEvolve(
              leagueData.seeds[2*day + 1], leagueData.startTimes[2*day + 1], 
              [allTeamsSkills[t0], allTeamsSkills[t1]], 
              [leagueData.teamIds[t0], leagueData.teamIds[t1]], 
              [leagueData.tactics[2 * day + 1][t0], leagueData.tactics[2 * day + 1][t1]], 
              [allMatchLogs[t0], allMatchLogs[t1]],
              [is2nd = true, isHom = true, isPlay = false]
          ).should.be.fulfilled;
          allTeamsSkills[t0] = vec2str(newSkills[0]);
          allTeamsSkills[t1] = vec2str(newSkills[1]);
          allMatchLogs[t0] = newLogs[0].toString();
          allMatchLogs[t1] = newLogs[1].toString(); 
          goals0 = await encodeLog.getNGoals(newLogs[0]).should.be.fulfilled;
          goals1 = await encodeLog.getNGoals(newLogs[1]).should.be.fulfilled;
          leagueData.results[nMatchesPerDay * day + matchIdxInDay] = [goals0.toNumber(), goals1.toNumber()];
      }
      leagueData.teamStates.push([...allTeamsSkills]);        
      leagueData.matchLogs.push([...allMatchLogs]);   
      var {0: rnking, 1: lPoints} = await champs.computeLeagueLeaderBoard([...leagueData.results], day, leagueData.seeds[2*day + 1]).should.be.fulfilled;
      leagueData.points.push(vec2str(lPoints));   
  }
  return leagueData;
}

function readCreatedLeagueData() {
  return JSON.parse(fs.readFileSync('test/testdata/fullleague.json', 'utf8'));
}

function readCreatedLeagueLeafs() {
  return JSON.parse(fs.readFileSync('test/testdata/leafsPerHalf.json', 'utf8'));
}

  module.exports = {
    createUniverse,
    createLeagueData,
    buildLeafs,
    readCreatedLeagueData,
    readCreatedLeagueLeafs,
  }