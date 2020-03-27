const NULL_BYTES32 = web3.eth.abi.encodeParameter('bytes32','0x0');

const merkleUtils = require('../utils/merkleUtils.js');
const teamsUtils = require('../utils/teamsUtils.js');

function readTeam442(filepath = 'test/testdata/team442.json'){
  fs.readFile(filepath, 'utf8', function (err, data) {
      if (err) throw err;
      try {
        result = JSON.parse(data);
      } catch (e) {
          console.error( e );
      }
  });
  return result;
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
function createLeague() {
  teamState = readTeam442();
  console.log(teamState);
  return teamState;
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

  module.exports = {
    createUniverse,
    createLeague,
  }