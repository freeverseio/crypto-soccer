var Updater = artifacts.require("./Updater.sol");

module.exports = function(deployer) {
  deployer.deploy(Updater);
};

