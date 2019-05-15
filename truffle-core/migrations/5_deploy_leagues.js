const States = artifacts.require('LeagueState');
const Engine = artifacts.require('Engine');
const Leagues = artifacts.require('Leagues');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Leagues, Engine.address, States.address);
    })
    .catch(console.error);
};

