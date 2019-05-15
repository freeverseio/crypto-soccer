const States = artifacts.require('LeagueState');
const Assets = artifacts.require('Teams');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Assets, States.address);
    })
    .catch(console.error);
};

