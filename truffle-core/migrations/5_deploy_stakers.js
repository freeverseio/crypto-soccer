const Stakers = artifacts.require("Stakers");
const GameController = artifacts.require("GameController");

module.exports = function (deployer) {
  deployer.then(async () => {
      const gameController = await deployer.deploy(GameController);
      const staker = await deployer.deploy(Stakers, gameController.address);
      await gameController.setStakersContractAddress(staker.address);
    })
    .catch(console.error);
};

