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

    console.log(deployer)

    config = {};
    config.assetsContractAddress = Assets.address;
    config.statesContractAddress = States.address;
    config.engineContractAddress = Engine.address;
    config.gameControllerContractAddress = GameController.address;
    config.leaguesContractAddress = Leagues.address;
    config.stakersContractAddress = Stakers.address;
    await fs.writeFileSync("../migrate_" + deployer.network + ".json",JSON.stringify(config, null, 4));

    console.log(log);
  })
    .catch(console.error);
};

