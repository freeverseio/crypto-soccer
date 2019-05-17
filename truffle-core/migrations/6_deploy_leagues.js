const fs = require('fs');
const States = artifacts.require('LeagueState');
const Engine = artifacts.require('Engine');
const GameController = artifacts.require("GameController");
const Leagues = artifacts.require('Leagues');

module.exports = function (deployer) {
  deployer.then(async () => {
      const leagues = await deployer.deploy(Leagues, Engine.address, States.address);
      fs.appendFileSync('deploy_addresses.txt', "leagues : " + leagues.address);
      await leagues.setStakersContract(GameController.address);
    })
    .catch(console.error);
};

