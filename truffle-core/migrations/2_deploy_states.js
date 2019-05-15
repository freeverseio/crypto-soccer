const States = artifacts.require('LeagueState');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(States);
    })
    .catch(console.error);
};

