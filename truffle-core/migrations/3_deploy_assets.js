const States = artifacts.require('LeagueState');
const Assets = artifacts.require('FreezableAssets');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Assets, States.address);
    })
    .catch(console.error);
};

