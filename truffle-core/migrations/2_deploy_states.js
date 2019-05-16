const fs = require('fs');
const States = artifacts.require('LeagueState');

module.exports = function (deployer) {
  deployer.then(async () => {
      const states = await deployer.deploy(States);
      fs.writeFileSync('deploy_addresses.txt', "stateLib : " + states.address + "\n");
    })
    .catch(console.error);
};

