const fs = require('fs');
const Assets = artifacts.require('Assets');
const States = artifacts.require('LeagueState');
const Engine = artifacts.require('Engine');
const GameController = artifacts.require("GameController");
const Leagues = artifacts.require('Leagues');
const Stakers = artifacts.require('Stakers');

module.exports = function (deployer) {
  deployer.then(async () => {
    console.log("");
    console.log("🚀  Deployed on:", deployer.network)
    console.log("------------------------");
    config = {};
    config.assetsContractAddress = Assets.address;
    config.statesContractAddress = States.address;
    config.engineContractAddress = Engine.address;
    config.gameControllerContractAddress = GameController.address;
    config.leaguesContractAddress = Leagues.address;
    config.stakersContractAddress = Stakers.address;
    console.log(JSON.stringify(config, null, 4));
  })
    .catch(console.error);
};

