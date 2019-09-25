const Engine = artifacts.require('Engine');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Engine);
    })
    .catch(console.error);
};

