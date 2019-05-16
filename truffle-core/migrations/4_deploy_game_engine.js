const fs = require('fs');
const Engine = artifacts.require('Engine');

module.exports = function (deployer) {
  deployer.then(async () => {
      const engine = await deployer.deploy(Engine);
      fs.appendFileSync('deploy_addresses.txt', "engine : " + engine.address + "\n");
    })
    .catch(console.error);
};

