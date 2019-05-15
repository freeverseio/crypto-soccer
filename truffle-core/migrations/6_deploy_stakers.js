const Stakers = artifacts.require("Stakers");
const Leagues = artifacts.require('Leagues');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Stakers, Leagues.address);
    })
    .catch(console.error);
};

