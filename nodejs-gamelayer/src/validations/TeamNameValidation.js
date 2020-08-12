
const Web3 = require('web3')

class TeamNameValidation {
  hash() {
    const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
    return web3.utils.soliditySha3( {t: 'string', v: 'ciao'});
  }

}

module.exports = new TeamNameValidation();
