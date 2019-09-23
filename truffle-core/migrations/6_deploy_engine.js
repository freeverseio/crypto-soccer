const Engine = artifacts.require('Engine');
const EngineHL = artifacts.require('EngineHL');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Engine);
      await deployer.deploy(EngineHL);
    })
    .catch(console.error);
};

