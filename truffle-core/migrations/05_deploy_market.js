const Market = artifacts.require('Market');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Market);
    })
    .catch(console.error);
};

