const Cups = artifacts.require('Cups');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Cups);
    })
    .catch(console.error);
};

