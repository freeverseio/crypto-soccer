const fs = require('fs');
const States = artifacts.require('LeagueState');
const Assets = artifacts.require('Assets');

module.exports = function (deployer) {
  deployer.then(async () => {
      const assets = await deployer.deploy(Assets, States.address);
      fs.appendFileSync('deploy_addresses.txt', "assets : " + assets.address + "\n");
    })
    .catch(console.error);
};

