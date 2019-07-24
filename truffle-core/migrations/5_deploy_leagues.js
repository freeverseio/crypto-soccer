const States = artifacts.require('LeagueState');
const Engine = artifacts.require('Engine');
const GameController = artifacts.require("GameController");
const Leagues = artifacts.require('Leagues');

module.exports = function (deployer) {
  // deployer.then(async () => {
  //     const leagues = await deployer.deploy(Leagues, Engine.address, States.address);
  //     const gameController = await deployer.deploy(GameController, leagues.address);
  //     await leagues.setStakersContract(gameController.address);
  //   })
  //   .catch(console.error);
};

