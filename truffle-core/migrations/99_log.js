const fs = require('fs');
const Assets = artifacts.require('Market');
const AssetsLib = artifacts.require('AssetsLib');

module.exports = function (deployer) {
  deployer.then(async () => {
    console.log("");
    console.log("🚀  Deployed on:", deployer.network)
    console.log("------------------------");
    config = {};
    config.assetsContractAddress = Assets.address;
    config.assetsLibContractAddress = AssetsLib.address;
    console.log(JSON.stringify(config, null, 4));
  })
    .catch(console.error);
};

