const Stakers = artifacts.require("Stakers");
const GameController = artifacts.require("GameController");

module.exports = function (deployer) {
  deployer.then(async () => {
      const gameController = await deployer.deploy(GameController);
      await deployer.deploy(Stakers, gameController.address);
    })
    .catch(console.error);
};

