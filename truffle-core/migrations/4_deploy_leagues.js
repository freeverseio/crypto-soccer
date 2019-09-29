const Championships = artifacts.require('Championships');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Championships);
    })
    .catch(console.error);
};

