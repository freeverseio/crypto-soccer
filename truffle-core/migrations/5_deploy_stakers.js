const fs = require('fs');
const Stakers = artifacts.require("Stakers");
const GameController = artifacts.require("GameController");

module.exports = function (deployer) {
  deployer.then(async () => {
      const gameController = await deployer.deploy(GameController);
      fs.appendFileSync('deploy_addresses.txt', "gameController : " + gameController.address + "\n");
      const stakers = await deployer.deploy(Stakers, gameController.address);
      fs.appendFileSync('deploy_addresses.txt', "stakers : " + stakers.address + "\n");
    })
    .catch(console.error);
};

