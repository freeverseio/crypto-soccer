const NULL_BYTES32 = web3.eth.abi.encodeParameter('bytes32','0x0');

const merkleUtils = require('../utils/merkleUtils.js');

// returns, for a given TZ, [nActiveCountry0,... 1023]
function createActiveTeams(nLevels) {
  // activeTeams
  activeTeams = Array.from(new Array(1024), (x,i) => 0);

}

  module.exports = {
    test,
  }