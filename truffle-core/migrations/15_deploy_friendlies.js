const Friendlies = artifacts.require('Friendlies');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Friendlies);
    })
    .catch(console.error);
};

