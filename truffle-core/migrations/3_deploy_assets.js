const States = artifacts.require('LeagueState');
const Assets = artifacts.require('Assets');

module.exports = function (deployer) {
  deployer.then(async () => {
      const assets = await deployer.deploy(Assets);
    })
    .catch(console.error);
};

