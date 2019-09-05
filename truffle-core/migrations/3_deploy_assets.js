const States = artifacts.require('LeagueState');
const Assets = artifacts.require('Assets');

module.exports = function (deployer) {
  deployer.then(async () => {
      deployReceipt = await deployer.deploy(Assets, States.address);
    })
    .catch(console.error);
};

