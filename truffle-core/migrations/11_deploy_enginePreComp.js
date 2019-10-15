const EnginePreComp = artifacts.require('EnginePreComp');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(EnginePreComp);
    })
    .catch(console.error);
};

