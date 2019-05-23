const fs = require('fs');
const Assets = artifacts.require('Assets');
const States = artifacts.require('LeagueState');
const Engine = artifacts.require('Engine');
const GameController = artifacts.require("GameController");
const Leagues = artifacts.require('Leagues');
const Stakers = artifacts.require('Stakers');

module.exports = function (deployer) {
  deployer.then(async () => {
      let log = "--------------------------------" + "\n";
      log += "Assets:         " + Assets.address + "\n";
      log += "States:         " + States.address + "\n";
      log += "Engine:         " + Engine.address + "\n";
      log += "GameController: " + GameController.address + "\n";
      log += "Leagues:        " + Leagues.address + "\n";
      log += "Stakers:        " + Stakers.address + "\n";
      log += "--------------------------------";

      fs.writeFileSync('deploy_addresses.txt', log);
      console.log(log);
    })
    .catch(console.error);
};

