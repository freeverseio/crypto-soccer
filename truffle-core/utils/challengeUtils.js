const NULL_BYTES32 = web3.eth.abi.encodeParameter('bytes32','0x0');

const merkleUtils = require('../utils/merkleUtils.js');

function sortDec(array) { 
  arr = [...array];
  return arr.sort(function(a, b){return b-a});
}

// returns, for a given TZ, [nActiveCountry0,... 1023]
function createUniverse(nLevels) {
  nCountries = 1024;
  nTeamsInLeague = 8;
  activeLeagues = Array.from(new Array(nCountries), (x,i) => 2 * ((i+1) % 3));
  activeTeams = Array.from(activeLeagues, (x,i) => x * nTeamsInLeague);
  orgMap = [];
  for (c = 0; c < 1024; c++) {
    teamIdxs = Array.from(new Array(activeTeams[c]), (x,i) => i);
    if (c % 2 == 0) { orgMap = orgMap.concat(teamIdxs); }
    else  { orgMap = orgMap.concat(sortDec(teamIdxs)); }
  }
  return [activeTeams, orgMap];
}

  module.exports = {
    createUniverse,
  }