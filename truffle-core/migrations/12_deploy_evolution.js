const Evolution = artifacts.require('Evolution');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Evolution);
    })
    .catch(console.error);
};

