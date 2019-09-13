const Leagues = artifacts.require('Leagues');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Leagues);
    })
    .catch(console.error);
};

