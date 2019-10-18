const Assets = artifacts.require('Assets');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Assets);
    })
    .catch(console.error);
};

